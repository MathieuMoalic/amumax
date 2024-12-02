package engine_old

import (
	"math"

	"github.com/MathieuMoalic/amumax/src/cuda_old"
	"github.com/MathieuMoalic/amumax/src/engine_old/data_old"
	"github.com/MathieuMoalic/amumax/src/engine_old/log_old"
)

type rk56 struct {
}

func (rk *rk56) Step() {

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

	k1, k2, k3, k4, k5, k6, k7, k8 := cuda_old.Buffer(3, size), cuda_old.Buffer(3, size), cuda_old.Buffer(3, size), cuda_old.Buffer(3, size), cuda_old.Buffer(3, size), cuda_old.Buffer(3, size), cuda_old.Buffer(3, size), cuda_old.Buffer(3, size)
	defer cuda_old.Recycle(k1)
	defer cuda_old.Recycle(k2)
	defer cuda_old.Recycle(k3)
	defer cuda_old.Recycle(k4)
	defer cuda_old.Recycle(k5)
	defer cuda_old.Recycle(k6)
	defer cuda_old.Recycle(k7)
	defer cuda_old.Recycle(k8)
	//k2 will be recyled as k9

	h := float32(Dt_si * gammaLL) // internal time step = Dt * gammaLL

	// stage 1
	torqueFn(k1)

	// stage 2
	Time = t0 + (1./6.)*Dt_si
	cuda_old.Madd2(m, m, k1, 1, (1./6.)*h) // m = m*1 + k1*h/6
	NormMag.normalize()
	torqueFn(k2)

	// stage 3
	Time = t0 + (4./15.)*Dt_si
	cuda_old.Madd3(m, m0, k1, k2, 1, (4./75.)*h, (16./75.)*h)
	NormMag.normalize()
	torqueFn(k3)

	// stage 4
	Time = t0 + (2./3.)*Dt_si
	cuda_old.Madd4(m, m0, k1, k2, k3, 1, (5./6.)*h, (-8./3.)*h, (5./2.)*h)
	NormMag.normalize()
	torqueFn(k4)

	// stage 5
	Time = t0 + (4./5.)*Dt_si
	cuda_old.Madd5(m, m0, k1, k2, k3, k4, 1, (-8./5.)*h, (144./25.)*h, (-4.)*h, (16./25.)*h)
	NormMag.normalize()
	torqueFn(k5)

	// stage 6
	Time = t0 + (1.)*Dt_si
	cuda_old.Madd6(m, m0, k1, k2, k3, k4, k5, 1, (361./320.)*h, (-18./5.)*h, (407./128.)*h, (-11./80.)*h, (55./128.)*h)
	NormMag.normalize()
	torqueFn(k6)

	// stage 7
	Time = t0
	cuda_old.Madd5(m, m0, k1, k3, k4, k5, 1, (-11./640.)*h, (11./256.)*h, (-11/160.)*h, (11./256.)*h)
	NormMag.normalize()
	torqueFn(k7)

	// stage 8
	Time = t0 + (1.)*Dt_si
	cuda_old.Madd7(m, m0, k1, k2, k3, k4, k5, k7, 1, (93./640.)*h, (-18./5.)*h, (803./256.)*h, (-11./160.)*h, (99./256.)*h, (1.)*h)
	NormMag.normalize()
	torqueFn(k8)

	// stage 9: 6th order solution
	Time = t0 + (1.)*Dt_si
	//madd6(m, m0, k1, k3, k4, k5, k6, 1, (31./384.)*h, (1125./2816.)*h, (9./32.)*h, (125./768.)*h, (5./66.)*h)
	cuda_old.Madd7(m, m0, k1, k3, k4, k5, k7, k8, 1, (7./1408.)*h, (1125./2816.)*h, (9./32.)*h, (125./768.)*h, (5./66.)*h, (5./66.)*h)
	NormMag.normalize()
	torqueFn(k2) // re-use k2

	// error estimate
	Err := cuda_old.Buffer(3, size)
	defer cuda_old.Recycle(Err)
	cuda_old.Madd4(Err, k1, k6, k7, k8, (-5. / 66.), (-5. / 66.), (5. / 66.), (5. / 66.))

	// determine error
	err := cuda_old.MaxVecNorm(Err) * float64(h)

	// adjust next time step
	if err < MaxErr || Dt_si <= MinDt || FixDt != 0 { // mindt check to avoid infinite loop
		// step OK
		setLastErr(err)
		setMaxTorque(k2)
		NSteps++
		Time = t0 + Dt_si
		adaptDt(math.Pow(MaxErr/err, 1./6.))
	} else {
		// undo bad step
		//log.Println("Bad step at t=", t0, ", err=", err)
		log_old.AssertMsg(FixDt == 0, "Invalid step: cannot undo step when FixDt is set")
		Time = t0
		data_old.Copy(m, m0)
		NUndone++
		adaptDt(math.Pow(MaxErr/err, 1./7.))
	}
}

func (rk *rk56) Free() {
}
