package engine_old

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
type rk23 struct {
	k1 *data_old.Slice // torque at end of step is kept for beginning of next step
}

func (rk *rk23) Step() {
	m := NormMag.Buffer()
	size := m.Size()

	if FixDt != 0 {
		Dt_si = FixDt
	}

	// upon resize: remove wrongly sized k1
	if rk.k1.Size() != m.Size() {
		rk.Free()
	}

	// first step ever: one-time k1 init and eval
	if rk.k1 == nil {
		rk.k1 = cuda_old.NewSlice(3, size)
		torqueFn(rk.k1)
	}

	// FSAL cannot be used with temperature
	if !Temp.isZero() {
		torqueFn(rk.k1)
	}

	t0 := Time
	// backup magnetization
	m0 := cuda_old.Buffer(3, size)
	defer cuda_old.Recycle(m0)
	data_old.Copy(m0, m)

	k2, k3, k4 := cuda_old.Buffer(3, size), cuda_old.Buffer(3, size), cuda_old.Buffer(3, size)
	defer cuda_old.Recycle(k2)
	defer cuda_old.Recycle(k3)
	defer cuda_old.Recycle(k4)

	h := float32(Dt_si * gammaLL) // internal time step = Dt * gammaLL

	// there is no explicit stage 1: k1 from previous step

	// stage 2
	Time = t0 + (1./2.)*Dt_si
	cuda_old.Madd2(m, m, rk.k1, 1, (1./2.)*h) // m = m*1 + k1*h/2
	NormMag.normalize()
	torqueFn(k2)

	// stage 3
	Time = t0 + (3./4.)*Dt_si
	cuda_old.Madd2(m, m0, k2, 1, (3./4.)*h) // m = m0*1 + k2*3/4
	NormMag.normalize()
	torqueFn(k3)

	// 3rd order solution
	cuda_old.Madd4(m, m0, rk.k1, k2, k3, 1, (2./9.)*h, (1./3.)*h, (4./9.)*h)
	NormMag.normalize()

	// error estimate
	Time = t0 + Dt_si
	torqueFn(k4)
	Err := k2 // re-use k2 as error
	// difference of 3rd and 2nd order torque without explicitly storing them first
	cuda_old.Madd4(Err, rk.k1, k2, k3, k4, (7./24.)-(2./9.), (1./4.)-(1./3.), (1./3.)-(4./9.), (1. / 8.))

	// determine error
	err := cuda_old.MaxVecNorm(Err) * float64(h)

	// adjust next time step
	if err < MaxErr || Dt_si <= MinDt || FixDt != 0 { // mindt check to avoid infinite loop
		// step OK
		setLastErr(err)
		setMaxTorque(k4)
		NSteps++
		Time = t0 + Dt_si
		adaptDt(math.Pow(MaxErr/err, 1./3.))
		data_old.Copy(rk.k1, k4) // FSAL
	} else {
		// undo bad step
		//log.Println("Bad step at t=", t0, ", err=", err)
		log_old.AssertMsg(FixDt == 0, "Invalid step: cannot undo step when FixDt is set")
		Time = t0
		data_old.Copy(m, m0)
		NUndone++
		adaptDt(math.Pow(MaxErr/err, 1./4.))
	}
}

func (rk *rk23) Free() {
	rk.k1.Free()
	rk.k1 = nil
}
