package solver

import (
	"math"

	"github.com/MathieuMoalic/amumax/src/cuda_old"
	"github.com/MathieuMoalic/amumax/src/engine_old/data_old"
	"github.com/MathieuMoalic/amumax/src/engine_old/log_old"
)

func (s *Solver) rk45() {
	m := NormMag.Buffer()
	size := m.Size()

	if s.FixDt != 0 {
		s.dt_si = s.FixDt
	}

	// upon resize: remove wrongly sized k1
	if s.previousStepBuffer != nil && s.previousStepBuffer.Size() != m.Size() {
		s.previousStepBuffer.Free()
	}

	// first step ever: one-time k1 init and eval
	if s.previousStepBuffer == nil {
		s.previousStepBuffer = cuda_old.NewSlice(3, size)
		s.torqueFn(s.previousStepBuffer)
	}

	// FSAL cannot be used with finite temperature
	if !Temp.isZero() {
		s.torqueFn(s.previousStepBuffer)
	}

	t0 := s.Time
	// backup magnetization
	m0 := cuda_old.Buffer(3, size)
	defer cuda_old.Recycle(m0)
	data_old.Copy(m0, m)

	k2, k3, k4, k5, k6 := cuda_old.Buffer(3, size), cuda_old.Buffer(3, size), cuda_old.Buffer(3, size), cuda_old.Buffer(3, size), cuda_old.Buffer(3, size)
	defer cuda_old.Recycle(k2)
	defer cuda_old.Recycle(k3)
	defer cuda_old.Recycle(k4)
	defer cuda_old.Recycle(k5)
	defer cuda_old.Recycle(k6)
	// k2 will be re-used as k7

	h := float32(s.dt_si * gammaLL) // internal time step = Dt * gammaLL

	// there is no explicit stage 1: k1 from previous step

	// stage 2
	s.Time = t0 + (1./5.)*s.dt_si
	cuda_old.Madd2(m, m, s.previousStepBuffer, 1, (1./5.)*h) // m = m*1 + k1*h/5
	NormMag.normalize()
	s.torqueFn(k2)

	// stage 3
	s.Time = t0 + (3./10.)*s.dt_si
	cuda_old.Madd3(m, m0, s.previousStepBuffer, k2, 1, (3./40.)*h, (9./40.)*h)
	NormMag.normalize()
	s.torqueFn(k3)

	// stage 4
	s.Time = t0 + (4./5.)*s.dt_si
	cuda_old.Madd4(m, m0, s.previousStepBuffer, k2, k3, 1, (44./45.)*h, (-56./15.)*h, (32./9.)*h)
	NormMag.normalize()
	s.torqueFn(k4)

	// stage 5
	s.Time = t0 + (8./9.)*s.dt_si
	cuda_old.Madd5(m, m0, s.previousStepBuffer, k2, k3, k4, 1, (19372./6561.)*h, (-25360./2187.)*h, (64448./6561.)*h, (-212./729.)*h)
	NormMag.normalize()
	s.torqueFn(k5)

	// stage 6
	s.Time = t0 + (1.)*s.dt_si
	cuda_old.Madd6(m, m0, s.previousStepBuffer, k2, k3, k4, k5, 1, (9017./3168.)*h, (-355./33.)*h, (46732./5247.)*h, (49./176.)*h, (-5103./18656.)*h)
	NormMag.normalize()
	s.torqueFn(k6)

	// stage 7: 5th order solution
	s.Time = t0 + (1.)*s.dt_si
	// no k2
	cuda_old.Madd6(m, m0, s.previousStepBuffer, k3, k4, k5, k6, 1, (35./384.)*h, (500./1113.)*h, (125./192.)*h, (-2187./6784.)*h, (11./84.)*h) // 5th
	NormMag.normalize()
	k7 := k2       // re-use k2
	s.torqueFn(k7) // next torque if OK

	// error estimate
	Err := cuda_old.Buffer(3, size) //k3 // re-use k3 as error estimate
	defer cuda_old.Recycle(Err)
	cuda_old.Madd6(Err, s.previousStepBuffer, k3, k4, k5, k6, k7, (35./384.)-(5179./57600.), (500./1113.)-(7571./16695.), (125./192.)-(393./640.), (-2187./6784.)-(-92097./339200.), (11./84.)-(187./2100.), (0.)-(1./40.))

	// determine error
	err := cuda_old.MaxVecNorm(Err) * float64(h)

	// adjust next time step
	if err < s.MaxErr || s.dt_si <= s.MinDt || s.FixDt != 0 { // mindt check to avoid infinite loop
		// step OK
		s.setLastErr(err)
		s.setMaxTorque(k7)
		s.NSteps++
		s.Time = t0 + s.dt_si
		s.adaptDt(math.Pow(s.MaxErr/err, 1./5.))
		data_old.Copy(s.previousStepBuffer, k7) // FSAL
	} else {
		// undo bad step
		//log.Println("Bad step at t=", t0, ", err=", err)
		log_old.AssertMsg(s.FixDt == 0, "Invalid step: cannot undo step when s.fixDt is set")
		s.Time = t0
		data_old.Copy(m, m0)
		s.nUndone++
		s.adaptDt(math.Pow(s.MaxErr/err, 1./6.))
	}
}
