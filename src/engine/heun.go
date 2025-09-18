package engine

import (
	"math"

	"github.com/MathieuMoalic/amumax/src/cuda"
	"github.com/MathieuMoalic/amumax/src/log"
)

// Adaptive heun solver.
type heun struct{}

// Step Adaptive Heun method, can be used as solver.Step
func (*heun) Step() {
	y := NormMag.Buffer()
	dy0 := cuda.Buffer(VECTOR, y.Size())
	defer cuda.Recycle(dy0)

	if FixDt != 0 {
		DtSi = FixDt
	}

	dt := float32(DtSi * gammaLL)
	log.AssertMsg(dt > 0, "Invalid time step: dt must be positive in Heun Step")

	// stage 1
	torqueFn(dy0)
	cuda.Madd2(y, y, dy0, 1, dt) // y = y + dt * dy

	// stage 2
	dy := cuda.Buffer(3, y.Size())
	defer cuda.Recycle(dy)
	Time += DtSi
	torqueFn(dy)

	err := cuda.MaxVecDiff(dy0, dy) * float64(dt)

	// adjust next time step
	if err < MaxErr || DtSi <= MinDt || FixDt != 0 { // mindt check to avoid infinite loop
		// step OK
		cuda.Madd3(y, y, dy, dy0, 1, 0.5*dt, -0.5*dt)
		NormMag.normalize()
		NSteps++
		adaptDt(math.Pow(MaxErr/err, 1./2.))
		setLastErr(err)
		setMaxTorque(dy)
	} else {
		// undo bad step
		log.AssertMsg(FixDt == 0, "Invalid step: cannot undo step when FixDt is set in Heun Step")
		Time -= DtSi
		cuda.Madd2(y, y, dy0, 1, -dt)
		NUndone++
		adaptDt(math.Pow(MaxErr/err, 1./3.))
	}
}

func (*heun) Free() {}
