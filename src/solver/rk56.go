package solver

import (
	"math"

	"github.com/MathieuMoalic/amumax/src/cuda"
	"github.com/MathieuMoalic/amumax/src/data"
	"github.com/MathieuMoalic/amumax/src/log_old"
)

func (s *Solver) rk56() {

	m := NormMag.Buffer()
	size := m.Size()

	if s.fixDt != 0 {
		s.dt_si = s.fixDt
	}

	t0 := s.time
	// backup magnetization
	m0 := cuda.Buffer(3, size)
	defer cuda.Recycle(m0)
	data.Copy(m0, m)

	k1, k2, k3, k4, k5, k6, k7, k8 := cuda.Buffer(3, size), cuda.Buffer(3, size), cuda.Buffer(3, size), cuda.Buffer(3, size), cuda.Buffer(3, size), cuda.Buffer(3, size), cuda.Buffer(3, size), cuda.Buffer(3, size)
	defer cuda.Recycle(k1)
	defer cuda.Recycle(k2)
	defer cuda.Recycle(k3)
	defer cuda.Recycle(k4)
	defer cuda.Recycle(k5)
	defer cuda.Recycle(k6)
	defer cuda.Recycle(k7)
	defer cuda.Recycle(k8)
	//k2 will be recyled as k9

	h := float32(s.dt_si * gammaLL) // internal time step = Dt * gammaLL

	// stage 1
	s.torqueFn(k1)

	// stage 2
	s.time = t0 + (1./6.)*s.dt_si
	cuda.Madd2(m, m, k1, 1, (1./6.)*h) // m = m*1 + k1*h/6
	NormMag.normalize()
	s.torqueFn(k2)

	// stage 3
	s.time = t0 + (4./15.)*s.dt_si
	cuda.Madd3(m, m0, k1, k2, 1, (4./75.)*h, (16./75.)*h)
	NormMag.normalize()
	s.torqueFn(k3)

	// stage 4
	s.time = t0 + (2./3.)*s.dt_si
	cuda.Madd4(m, m0, k1, k2, k3, 1, (5./6.)*h, (-8./3.)*h, (5./2.)*h)
	NormMag.normalize()
	s.torqueFn(k4)

	// stage 5
	s.time = t0 + (4./5.)*s.dt_si
	cuda.Madd5(m, m0, k1, k2, k3, k4, 1, (-8./5.)*h, (144./25.)*h, (-4.)*h, (16./25.)*h)
	NormMag.normalize()
	s.torqueFn(k5)

	// stage 6
	s.time = t0 + (1.)*s.dt_si
	cuda.Madd6(m, m0, k1, k2, k3, k4, k5, 1, (361./320.)*h, (-18./5.)*h, (407./128.)*h, (-11./80.)*h, (55./128.)*h)
	NormMag.normalize()
	s.torqueFn(k6)

	// stage 7
	s.time = t0
	cuda.Madd5(m, m0, k1, k3, k4, k5, 1, (-11./640.)*h, (11./256.)*h, (-11/160.)*h, (11./256.)*h)
	NormMag.normalize()
	s.torqueFn(k7)

	// stage 8
	s.time = t0 + (1.)*s.dt_si
	cuda.Madd7(m, m0, k1, k2, k3, k4, k5, k7, 1, (93./640.)*h, (-18./5.)*h, (803./256.)*h, (-11./160.)*h, (99./256.)*h, (1.)*h)
	NormMag.normalize()
	s.torqueFn(k8)

	// stage 9: 6th order solution
	s.time = t0 + (1.)*s.dt_si
	//madd6(m, m0, k1, k3, k4, k5, k6, 1, (31./384.)*h, (1125./2816.)*h, (9./32.)*h, (125./768.)*h, (5./66.)*h)
	cuda.Madd7(m, m0, k1, k3, k4, k5, k7, k8, 1, (7./1408.)*h, (1125./2816.)*h, (9./32.)*h, (125./768.)*h, (5./66.)*h, (5./66.)*h)
	NormMag.normalize()
	s.torqueFn(k2) // re-use k2

	// error estimate
	Err := cuda.Buffer(3, size)
	defer cuda.Recycle(Err)
	cuda.Madd4(Err, k1, k6, k7, k8, (-5. / 66.), (-5. / 66.), (5. / 66.), (5. / 66.))

	// determine error
	err := cuda.MaxVecNorm(Err) * float64(h)

	// adjust next time step
	if err < s.maxErr || s.dt_si <= s.minDt || s.fixDt != 0 { // mindt check to avoid infinite loop
		// step OK
		s.setLastErr(err)
		s.setMaxTorque(k2)
		s.nSteps++
		s.time = t0 + s.dt_si
		s.adaptDt(math.Pow(s.maxErr/err, 1./6.))
	} else {
		// undo bad step
		//log.Println("Bad step at t=", t0, ", err=", err)
		log_old.AssertMsg(s.fixDt == 0, "Invalid step: cannot undo step when s.fixDt is set")
		s.time = t0
		data.Copy(m, m0)
		s.nUndone++
		s.adaptDt(math.Pow(s.maxErr/err, 1./7.))
	}
}
