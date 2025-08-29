package engine

import (
	"github.com/MathieuMoalic/amumax/src/engine/cuda"
	"github.com/MathieuMoalic/amumax/src/engine/log"
)

type euler struct{}

// Euler method, can be used as solver.Step.
func (*euler) Step() {
	y := NormMag.Buffer()
	dy0 := cuda.Buffer(VECTOR, y.Size())
	defer cuda.Recycle(dy0)

	torqueFn(dy0)
	setMaxTorque(dy0)

	// Adaptive time stepping: treat MaxErr as the maximum magnetization delta
	// (proportional to the error, but an overestimation for sure)
	var dt float32
	if FixDt != 0 {
		Dt_si = FixDt
		dt = float32(Dt_si * gammaLL)
	} else {
		dt = float32(MaxErr / LastTorque)
		Dt_si = float64(dt) / gammaLL
	}
	log.AssertMsg(dt > 0, "Euler solver requires fixed time step > 0")
	setLastErr(float64(dt) * LastTorque)

	cuda.Madd2(y, y, dy0, 1, dt) // y = y + dt * dy
	NormMag.normalize()
	Time += Dt_si
	NSteps++
}

func (*euler) Free() {}
