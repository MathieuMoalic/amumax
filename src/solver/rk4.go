package solver

import (
	"math"

	"github.com/MathieuMoalic/amumax/src/cuda"
	"github.com/MathieuMoalic/amumax/src/engine_old/data_old"
	"github.com/MathieuMoalic/amumax/src/engine_old/log_old"
)

// Classical 4th order RK solver.
func (s *Solver) rk4() {
	m := NormMag.Buffer()
	size := m.Size()

	if s.FixDt != 0 {
		s.dt_si = s.FixDt
	}

	t0 := s.Time
	// backup magnetization
	m0 := cuda.Buffer(3, size)
	defer cuda.Recycle(m0)
	data_old.Copy(m0, m)

	k1, k2, k3, k4 := cuda.Buffer(3, size), cuda.Buffer(3, size), cuda.Buffer(3, size), cuda.Buffer(3, size)

	defer cuda.Recycle(k1)
	defer cuda.Recycle(k2)
	defer cuda.Recycle(k3)
	defer cuda.Recycle(k4)

	h := float32(s.dt_si * gammaLL) // internal time step = Dt * gammaLL

	// stage 1
	s.torqueFn(k1)

	// stage 2
	s.Time = t0 + (1./2.)*s.dt_si
	cuda.Madd2(m, m, k1, 1, (1./2.)*h) // m = m*1 + k1*h/2
	NormMag.normalize()
	s.torqueFn(k2)

	// stage 3
	cuda.Madd2(m, m0, k2, 1, (1./2.)*h) // m = m0*1 + k2*1/2
	NormMag.normalize()
	s.torqueFn(k3)

	// stage 4
	s.Time = t0 + s.dt_si
	cuda.Madd2(m, m0, k3, 1, 1.*h) // m = m0*1 + k3*1
	NormMag.normalize()
	s.torqueFn(k4)

	err := cuda.MaxVecDiff(k1, k4) * float64(h)

	// adjust next time step
	if err < s.MaxErr || s.dt_si <= s.MinDt || s.FixDt != 0 { // mindt check to avoid infinite loop
		// step OK
		// 4th order solution
		cuda.Madd5(m, m0, k1, k2, k3, k4, 1, (1./6.)*h, (1./3.)*h, (1./3.)*h, (1./6.)*h)
		NormMag.normalize()
		s.NSteps++
		s.adaptDt(math.Pow(s.MaxErr/err, 1./4.))
		s.setLastErr(err)
		s.setMaxTorque(k4)
	} else {
		// undo bad step
		log_old.AssertMsg(s.FixDt == 0, "Invalid step: cannot undo step when s.fixDt is set")
		s.Time = t0
		data_old.Copy(m, m0)
		s.nUndone++
		s.adaptDt(math.Pow(s.MaxErr/err, 1./5.))
	}
}
