package engine_old

import (
	"math"

	"github.com/MathieuMoalic/amumax/src/engine_old/cuda_old"
	"github.com/MathieuMoalic/amumax/src/engine_old/log_old"
)

// Adaptive heun solver.
type heun struct{}

// Adaptive Heun method, can be used as solver.Step
func (*heun) Step() {
	y := NormMag.Buffer()
	dy0 := cuda_old.Buffer(VECTOR, y.Size())
	defer cuda_old.Recycle(dy0)

	if FixDt != 0 {
		Dt_si = FixDt
	}

	dt := float32(Dt_si * gammaLL)
	log_old.AssertMsg(dt > 0, "Invalid time step: dt must be positive in Heun Step")

	// stage 1
	torqueFn(dy0)
	cuda_old.Madd2(y, y, dy0, 1, dt) // y = y + dt * dy

	// stage 2
	dy := cuda_old.Buffer(3, y.Size())
	defer cuda_old.Recycle(dy)
	Time += Dt_si
	torqueFn(dy)

	err := cuda_old.MaxVecDiff(dy0, dy) * float64(dt)

	// adjust next time step
	if err < MaxErr || Dt_si <= MinDt || FixDt != 0 { // mindt check to avoid infinite loop
		// step OK
		cuda_old.Madd3(y, y, dy, dy0, 1, 0.5*dt, -0.5*dt)
		NormMag.normalize()
		NSteps++
		adaptDt(math.Pow(MaxErr/err, 1./2.))
		setLastErr(err)
		setMaxTorque(dy)
	} else {
		// undo bad step
		log_old.AssertMsg(FixDt == 0, "Invalid step: cannot undo step when FixDt is set in Heun Step")
		Time -= Dt_si
		cuda_old.Madd2(y, y, dy0, 1, -dt)
		NUndone++
		adaptDt(math.Pow(MaxErr/err, 1./3.))
	}
}

func (*heun) Free() {}
