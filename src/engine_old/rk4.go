package engine_old

import (
	"math"

	"github.com/MathieuMoalic/amumax/src/engine_old/cuda_old"
	"github.com/MathieuMoalic/amumax/src/engine_old/data_old"
	"github.com/MathieuMoalic/amumax/src/engine_old/log_old"
)

// Classical 4th order RK solver.
type rk4 struct {
}

func (rk *rk4) Step() {
	m := NormMag.Buffer()
	size := m.Size()

	if FixDt != 0 {
		Dt_si = FixDt
	}

	t0 := Time
	// backup magnetization
	m0 := cuda_old.Buffer(3, size)
	defer cuda_old.Recycle(m0)
	data_old.Copy(m0, m)

	k1, k2, k3, k4 := cuda_old.Buffer(3, size), cuda_old.Buffer(3, size), cuda_old.Buffer(3, size), cuda_old.Buffer(3, size)

	defer cuda_old.Recycle(k1)
	defer cuda_old.Recycle(k2)
	defer cuda_old.Recycle(k3)
	defer cuda_old.Recycle(k4)

	h := float32(Dt_si * gammaLL) // internal time step = Dt * gammaLL

	// stage 1
	torqueFn(k1)

	// stage 2
	Time = t0 + (1./2.)*Dt_si
	cuda_old.Madd2(m, m, k1, 1, (1./2.)*h) // m = m*1 + k1*h/2
	NormMag.normalize()
	torqueFn(k2)

	// stage 3
	cuda_old.Madd2(m, m0, k2, 1, (1./2.)*h) // m = m0*1 + k2*1/2
	NormMag.normalize()
	torqueFn(k3)

	// stage 4
	Time = t0 + Dt_si
	cuda_old.Madd2(m, m0, k3, 1, 1.*h) // m = m0*1 + k3*1
	NormMag.normalize()
	torqueFn(k4)

	err := cuda_old.MaxVecDiff(k1, k4) * float64(h)

	// adjust next time step
	if err < MaxErr || Dt_si <= MinDt || FixDt != 0 { // mindt check to avoid infinite loop
		// step OK
		// 4th order solution
		cuda_old.Madd5(m, m0, k1, k2, k3, k4, 1, (1./6.)*h, (1./3.)*h, (1./3.)*h, (1./6.)*h)
		NormMag.normalize()
		NSteps++
		adaptDt(math.Pow(MaxErr/err, 1./4.))
		setLastErr(err)
		setMaxTorque(k4)
	} else {
		// undo bad step
		log_old.AssertMsg(FixDt == 0, "Invalid step: cannot undo step when FixDt is set")
		Time = t0
		data_old.Copy(m, m0)
		NUndone++
		adaptDt(math.Pow(MaxErr/err, 1./5.))
	}
}

func (*rk4) Free() {}
