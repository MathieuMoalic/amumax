package solver

import (
	"github.com/MathieuMoalic/amumax/src/engine_old/cuda_old"
	"github.com/MathieuMoalic/amumax/src/engine_old/log_old"
)

// Euler method
func (s *Solver) euler() {
	y := NormMag.Buffer()
	dy0 := cuda_old.Buffer(3, y.Size())
	defer cuda_old.Recycle(dy0)

	s.torqueFn(dy0)
	s.setMaxTorque(dy0)

	// Adaptive time stepping: treat s.maxErr as the maximum magnetization delta
	// (proportional to the error, but an overestimation for sure)
	var dt float32
	if s.FixDt != 0 {
		s.dt_si = s.FixDt
		dt = float32(s.dt_si * gammaLL)
	} else {
		dt = float32(s.MaxErr / s.lastTorque)
		s.dt_si = float64(dt) / gammaLL
	}
	log_old.AssertMsg(dt > 0, "Euler solver requires fixed time step > 0")
	s.setLastErr(float64(dt) * s.lastTorque)

	cuda_old.Madd2(y, y, dy0, 1, dt) // y = y + dt * dy
	NormMag.normalize()
	s.Time += s.dt_si
	s.NSteps++
}
