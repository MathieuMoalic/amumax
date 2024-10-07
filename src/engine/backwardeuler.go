package engine

import (
	"github.com/MathieuMoalic/amumax/src/cuda"
	"github.com/MathieuMoalic/amumax/src/data"
	"github.com/MathieuMoalic/amumax/src/log"
)

// Implicit midpoint solver.
type backwardEuler struct {
	dy1 *data.Slice
}

// Euler method, can be used as solver.Step.
func (s *backwardEuler) Step() {
	log.AssertMsg(MaxErr > 0, "Backward euler solver requires MaxErr > 0")

	t0 := Time

	y := normMag.Buffer()

	y0 := cuda.Buffer(VECTOR, y.Size())
	defer cuda.Recycle(y0)
	data.Copy(y0, y)

	dy0 := cuda.Buffer(VECTOR, y.Size())
	defer cuda.Recycle(dy0)
	if s.dy1 == nil {
		s.dy1 = cuda.Buffer(VECTOR, y.Size())
	}
	dy1 := s.dy1

	Dt_si = FixDt
	dt := float32(Dt_si * gammaLL)
	log.AssertMsg(dt > 0, "Backward Euler solver requires fixed time step > 0")

	// Fist guess
	Time = t0 + 0.5*Dt_si // 0.5 dt makes it implicit midpoint method

	// with temperature, previous torque cannot be used as predictor
	if Temp.isZero() {
		cuda.Madd2(y, y0, dy1, 1, dt) // predictor euler step with previous torque
		normMag.normalize()
	}

	torqueFn(dy0)
	cuda.Madd2(y, y0, dy0, 1, dt) // y = y0 + dt * dy
	normMag.normalize()

	// One iteration
	torqueFn(dy1)
	cuda.Madd2(y, y0, dy1, 1, dt) // y = y0 + dt * dy1
	normMag.normalize()

	Time = t0 + Dt_si

	err := cuda.MaxVecDiff(dy0, dy1) * float64(dt)

	NSteps++
	setLastErr(err)
	setMaxTorque(dy1)
}

func (s *backwardEuler) Free() {
	s.dy1.Free()
	s.dy1 = nil
}
