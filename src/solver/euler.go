package solver

import (
	"github.com/MathieuMoalic/amumax/src/cuda"
)

// Euler method
func (s *Solver) euler() {
	y := s.magnetization.Slice
	dy0 := cuda.Buffer(3, y.Size())
	defer cuda.Recycle(dy0)

	s.torqueFn(dy0)
	s.setMaxTorque(dy0)

	// Adaptive time stepping: treat s.maxErr as the maximum magnetization delta
	// (proportional to the error, but an overestimation for sure)
	var dt float32
	if s.FixDt != 0 {
		s.dt_si = s.FixDt
		dt = float32(s.dt_si * s.gammaLL)
	} else {
		dt = float32(s.MaxErr / s.lastTorque)
		s.dt_si = float64(dt) / s.gammaLL
	}
	s.log.AssertMsg(dt > 0, "Euler solver requires fixed time step > 0")
	s.setLastErr(float64(dt) * s.lastTorque)

	cuda.Madd2(y, y, dy0, 1, dt) // y = y + dt * dy
	s.magnetization.Normalize()
	s.Time += s.dt_si
	s.NSteps++
}
