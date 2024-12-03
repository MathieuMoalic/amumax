package solver

import (
	"github.com/MathieuMoalic/amumax/src/cuda"
	"github.com/MathieuMoalic/amumax/src/slice"
)

// Backward Euler method
func (s *Solver) backWardEulerStep() {
	s.log.AssertMsg(s.MaxErr > 0, "Backward euler solver requires s.maxErr > 0")

	t0 := s.Time

	y := NormMag.Buffer()

	y0 := cuda.Buffer(3, y.Size())
	defer cuda.Recycle(y0)
	slice.Copy(y0, y)

	dy0 := cuda.Buffer(3, y.Size())
	defer cuda.Recycle(dy0)
	if s.previousStepBuffer == nil {
		s.previousStepBuffer = cuda.Buffer(3, y.Size())
	}
	dy1 := s.previousStepBuffer

	s.dt_si = s.FixDt
	dt := float32(s.dt_si * gammaLL)
	s.log.AssertMsg(dt > 0, "Backward Euler solver requires fixed time step > 0")

	// Fist guess
	s.Time = t0 + 0.5*s.dt_si // 0.5 dt makes it implicit midpoint method

	// with temperature, previous torque cannot be used as predictor
	if Temp.isZero() {
		cuda.Madd2(y, y0, dy1, 1, dt) // predictor euler step with previous torque
		NormMag.normalize()
	}

	s.torqueFn(dy0)
	cuda.Madd2(y, y0, dy0, 1, dt) // y = y0 + dt * dy
	NormMag.normalize()

	// One iteration
	s.torqueFn(dy1)
	cuda.Madd2(y, y0, dy1, 1, dt) // y = y0 + dt * dy1
	NormMag.normalize()

	s.Time = t0 + s.dt_si

	err := cuda.MaxVecDiff(dy0, dy1) * float64(dt)

	s.NSteps++
	s.setLastErr(err)
	s.setMaxTorque(dy1)
}
