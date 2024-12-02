package solver

import (
	"math"

	"github.com/MathieuMoalic/amumax/src/cuda_old"
	"github.com/MathieuMoalic/amumax/src/engine_old/data_old"
	"github.com/MathieuMoalic/amumax/src/engine_old/log_old"
)

// Bogacki-Shampine solver. 3rd order, 3 evaluations per step, adaptive step.
//
//	http://en.wikipedia.org/wiki/Bogacki-Shampine_method
//
//	k1 = f(tn, yn)
//	k2 = f(tn + 1/2 h, yn + 1/2 h k1)
//	k3 = f(tn + 3/4 h, yn + 3/4 h k2)
//	y{n+1}  = yn + 2/9 h k1 + 1/3 h k2 + 4/9 h k3            // 3rd order
//	k4 = f(tn + h, y{n+1})
//	z{n+1} = yn + 7/24 h k1 + 1/4 h k2 + 1/3 h k3 + 1/8 h k4 // 2nd order
func (s *Solver) rk23() {
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

	// FSAL cannot be used with temperature
	if !Temp.isZero() {
		s.torqueFn(s.previousStepBuffer)
	}

	t0 := s.Time
	// backup magnetization
	m0 := cuda_old.Buffer(3, size)
	defer cuda_old.Recycle(m0)
	data_old.Copy(m0, m)

	k2, k3, k4 := cuda_old.Buffer(3, size), cuda_old.Buffer(3, size), cuda_old.Buffer(3, size)
	defer cuda_old.Recycle(k2)
	defer cuda_old.Recycle(k3)
	defer cuda_old.Recycle(k4)

	h := float32(s.dt_si * gammaLL) // internal time step = Dt * gammaLL

	// there is no explicit stage 1: k1 from previous step

	// stage 2
	s.Time = t0 + (1./2.)*s.dt_si
	cuda_old.Madd2(m, m, s.previousStepBuffer, 1, (1./2.)*h) // m = m*1 + k1*h/2
	NormMag.normalize()
	s.torqueFn(k2)

	// stage 3
	s.Time = t0 + (3./4.)*s.dt_si
	cuda_old.Madd2(m, m0, k2, 1, (3./4.)*h) // m = m0*1 + k2*3/4
	NormMag.normalize()
	s.torqueFn(k3)

	// 3rd order solution
	cuda_old.Madd4(m, m0, s.previousStepBuffer, k2, k3, 1, (2./9.)*h, (1./3.)*h, (4./9.)*h)
	NormMag.normalize()

	// error estimate
	s.Time = t0 + s.dt_si
	s.torqueFn(k4)
	Err := k2 // re-use k2 as error
	// difference of 3rd and 2nd order torque without explicitly storing them first
	cuda_old.Madd4(Err, s.previousStepBuffer, k2, k3, k4, (7./24.)-(2./9.), (1./4.)-(1./3.), (1./3.)-(4./9.), (1. / 8.))

	// determine error
	err := cuda_old.MaxVecNorm(Err) * float64(h)

	// adjust next time step
	if err < s.MaxErr || s.dt_si <= s.MinDt || s.FixDt != 0 { // mindt check to avoid infinite loop
		// step OK
		s.setLastErr(err)
		s.setMaxTorque(k4)
		s.NSteps++
		s.Time = t0 + s.dt_si
		s.adaptDt(math.Pow(s.MaxErr/err, 1./3.))
		data_old.Copy(s.previousStepBuffer, k4) // FSAL
	} else {
		// undo bad step
		//log.Println("Bad step at t=", t0, ", err=", err)
		log_old.AssertMsg(s.FixDt == 0, "Invalid step: cannot undo step when s.fixDt is set")
		s.Time = t0
		data_old.Copy(m, m0)
		s.nUndone++
		s.adaptDt(math.Pow(s.MaxErr/err, 1./4.))
	}
}
