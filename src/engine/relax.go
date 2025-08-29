package engine

// Relax tries to find the minimum energy state.

import (
	"math"

	"github.com/MathieuMoalic/amumax/src/engine/cuda"
)

// Stopping relax Maxtorque in T. The user can check MaxTorque for sane values (e.g. 1e-3).
// If set to 0, relax() will stop when the average torque is steady or increasing.
var relaxTorqueThreshold float64 = -1.

// are we relaxing?
var relaxing = false

func relax() {
	checkExchangeLength()
	sanityCheck()
	Pause = false

	// Save the settings we are changing...
	prevType := Solvertype
	prevErr := MaxErr
	prevFixDt := FixDt
	prevPrecess := precess

	// ...to restore them later
	defer func() {
		setSolver(prevType)
		MaxErr = prevErr
		FixDt = prevFixDt
		precess = prevPrecess
		relaxing = false
		//	Temp.upd_reg = prevTemp
		//	Temp.invalidate()
		//	Temp.update()
	}()

	// Set good solver for relax
	setSolver(BOGAKISHAMPINE)
	FixDt = 0
	precess = false
	relaxing = true

	// Minimize energy: take steps as long as energy goes down.
	// This stops when energy reaches the numerical noise floor.
	const N = 3 // evaluate energy (expensive) every N steps
	relaxSteps(N)
	E0 := getTotalEnergy()
	relaxSteps(N)
	E1 := getTotalEnergy()
	for E1 < E0 && !Pause {
		relaxSteps(N)
		E0, E1 = E1, getTotalEnergy()
	}

	// Now we are already close to equilibrium, but energy is too noisy to be used any further.
	// So now we minimize the torque which is less noisy.
	solver := stepper.(*rk23)
	defer stepper.Free() // purge previous rk.k1 because FSAL will be dead wrong.

	maxTorque := func() float64 {
		return cuda.MaxVecNorm(solver.k1)
	}
	avgTorque := func() float32 {
		return cuda.Dot(solver.k1, solver.k1)
	}

	if relaxTorqueThreshold > 0 {
		// run as long as the max torque is above threshold. Then increase the accuracy and step more.
		for !Pause {
			for maxTorque() > relaxTorqueThreshold && !Pause {
				relaxSteps(N)
			}
			MaxErr /= math.Sqrt2
			if MaxErr < 1e-9 {
				break
			}
		}
	} else {
		// previous (<jan2018) behaviour: run as long as torque goes down. Then increase the accuracy and step more.
		// if MaxErr < 1e-9, this code won't run.
		var T0 float32
		var T1 float32 = avgTorque()
		// Step as long as torque goes down. Then increase the accuracy and step more.
		for MaxErr > 1e-9 && !Pause {
			MaxErr /= math.Sqrt2
			relaxSteps(N) // TODO: Play with other values
			T0, T1 = T1, avgTorque()
			for T1 < T0 && !Pause {
				relaxSteps(N) // TODO: Play with other values
				T0, T1 = T1, avgTorque()
			}
		}
	}
	Pause = true
}

// take n steps without setting pause when done or advancing time
func relaxSteps(n int) {
	t0 := Time
	stop := NSteps + n
	cond := func() bool { return NSteps < stop }
	const output = false
	runWhileInner(cond, output)
	Time = t0
}
