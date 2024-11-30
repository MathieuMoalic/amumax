package solver

import (
	"math"

	"github.com/MathieuMoalic/amumax/src/cuda"
	"github.com/MathieuMoalic/amumax/src/log_old"
)

// Adaptive Heun method, can be used as solver.Step
func (s *Solver) heun() {
	y := NormMag.Buffer()
	dy0 := cuda.Buffer(3, y.Size())
	defer cuda.Recycle(dy0)

	if s.fixDt != 0 {
		s.dt_si = s.fixDt
	}

	dt := float32(s.dt_si * gammaLL)
	log_old.AssertMsg(dt > 0, "Invalid time step: dt must be positive in Heun Step")

	// stage 1
	s.torqueFn(dy0)
	cuda.Madd2(y, y, dy0, 1, dt) // y = y + dt * dy

	// stage 2
	dy := cuda.Buffer(3, y.Size())
	defer cuda.Recycle(dy)
	s.Time += s.dt_si
	s.torqueFn(dy)

	err := cuda.MaxVecDiff(dy0, dy) * float64(dt)

	// adjust next time step
	if err < s.maxErr || s.dt_si <= s.minDt || s.fixDt != 0 { // mindt check to avoid infinite loop
		// step OK
		cuda.Madd3(y, y, dy, dy0, 1, 0.5*dt, -0.5*dt)
		NormMag.normalize()
		s.NSteps++
		s.adaptDt(math.Pow(s.maxErr/err, 1./2.))
		s.setLastErr(err)
		s.setMaxTorque(dy)
	} else {
		// undo bad step
		log_old.AssertMsg(s.fixDt == 0, "Invalid step: cannot undo step when s.fixDt is set in Heun Step")
		s.Time -= s.dt_si
		cuda.Madd2(y, y, dy0, 1, -dt)
		s.nUndone++
		s.adaptDt(math.Pow(s.maxErr/err, 1./3.))
	}
}
