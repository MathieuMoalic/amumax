package engine

import (
	"fmt"
	"math"
	"os"
	"time"

	"github.com/MathieuMoalic/amumax/src/cuda"
	"github.com/MathieuMoalic/amumax/src/data"
	"github.com/MathieuMoalic/amumax/src/log"
	"github.com/MathieuMoalic/amumax/src/mag"
	"github.com/MathieuMoalic/amumax/src/zarr"
)

// Solver globals
var (
	Time                    float64                              // time in seconds
	alarm                   float64                              // alarm clock marks end time of run, dt adaptation must not cross it!
	Pause                   = true                               // set pause at any time to stop running after the current step
	postStep                []func()                             // called on after every full time step
	Inject                                   = make(chan func()) // injects code in between time steps. Used by web interface.
	Dt_si                   float64          = 1e-15             // time step = dt_si (seconds) *dt_mul, which should be nice float32
	MinDt, MaxDt            float64                              // minimum and maximum time step
	MaxErr                  float64          = 1e-5              // maximum error/step
	Headroom                float64          = 0.8               // solver headroom, (Gustafsson, 1992, Control of Error and Convergence in ODE Solvers)
	LastErr, PeakErr        float64                              // error of last step, highest error ever
	LastTorque              float64                              // maxTorque of last time step
	NSteps, NUndone, NEvals int                                  // number of good steps, undone steps
	FixDt                   float64                              // fixed time step?
	stepper                 stepperInterface                     // generic step, can be EulerStep, HeunStep, etc
	Solvertype              int
	ProgressBar             zarr.ProgressBar
	exchangeLenghtWarned    bool
)

func init() {

	setSolver(DORMANDPRINCE)
	_ = newScalarValue("dt", "s", "Time Step", func() float64 { return Dt_si })
	_ = newScalarValue("LastErr", "", "Error of last step", func() float64 { return LastErr })
	_ = newScalarValue("PeakErr", "", "Overall maxium error per step", func() float64 { return PeakErr })
	_ = newScalarValue("NEval", "", "Total number of torque evaluations", func() float64 { return float64(NEvals) })
	exchangeLenghtWarned = false
}

// Time stepper like Euler, Heun, RK23
type stepperInterface interface {
	Step() // take time step using solver globals
	Free() // free resources, if any (e.g.: RK23 previous torque)
}

// Arguments for SetSolver
const (
	BACKWARD_EULER = -1
	EULER          = 1
	HEUN           = 2
	BOGAKISHAMPINE = 3
	RUNGEKUTTA     = 4
	DORMANDPRINCE  = 5
	FEHLBERG       = 6
)

func setSolver(typ int) {
	// free previous solver, if any
	if stepper != nil {
		stepper.Free()
	}
	switch typ {
	default:
		log.Log.ErrAndExit("SetSolver: unknown solver type:  %v", typ)
	case BACKWARD_EULER:
		stepper = new(backwardEuler)
	case EULER:
		stepper = new(euler)
	case HEUN:
		stepper = new(heun)
	case BOGAKISHAMPINE:
		stepper = new(rk23)
	case RUNGEKUTTA:
		stepper = new(rk4)
	case DORMANDPRINCE:
		stepper = new(rk45DP)
	case FEHLBERG:
		stepper = new(rk56)
	}
	Solvertype = typ
}

// write torque to dst and increment NEvals
func torqueFn(dst *data.Slice) {
	setTorque(dst)
	NEvals++
}

// update lastErr and peakErr
func setLastErr(err float64) {
	LastErr = err
	if err > PeakErr {
		PeakErr = err
	}
}

func setMaxTorque(τ *data.Slice) {
	LastTorque = cuda.MaxVecNorm(τ)
}

// adapt time step: dt *= corr, but limited to sensible values.
func adaptDt(corr float64) {
	if FixDt != 0 {
		Dt_si = FixDt
		return
	}

	// corner case triggered by err = 0: just keep time step.
	// see test/regression017.mx3
	if math.IsNaN(corr) {
		corr = 1
	}

	log.AssertMsg(corr != 0, "Time step too small, check if parameters are sensible")
	corr *= Headroom
	if corr > 2 {
		corr = 2
	}
	if corr < 1./2. {
		corr = 1. / 2.
	}
	Dt_si *= corr
	if MinDt != 0 && Dt_si < MinDt {
		Dt_si = MinDt
	}
	if MaxDt != 0 && Dt_si > MaxDt {
		Dt_si = MaxDt
	}
	if Dt_si == 0 {
		log.Log.ErrAndExit("time step too small")
	}

	// do not cross alarm time
	if Time < alarm && Time+Dt_si > alarm {
		Dt_si = alarm - Time
	}

	log.AssertMsg(Dt_si > 0, fmt.Sprint("Time step too small: ", Dt_si))
}

// Run the simulation for a number of seconds.
func runWithoutPrecession(seconds float64) {
	prevPrecess := precess
	run(seconds)
	precess = prevPrecess
}

// run the simulation for a number of seconds.
func run(seconds float64) {
	checkExchangeLenght()
	start := Time
	stop := Time + seconds
	alarm = stop // don't have dt adapt to go over alarm
	sanityCheck()
	Pause = false // may be set by <-Inject
	const output = true
	stepper.Free() // start from a clean state

	saveIfNeeded() // allow t=0 output
	ProgressBar := zarr.NewProgressBar(start, stop, "🧲", ShowProgresBar)

	for (Time < stop) && !Pause {
		select {
		default:
			ProgressBar.Update(Time)
			step(output)
		// accept tasks form Inject channel
		case f := <-Inject:
			f()
		}
	}
	ProgressBar.Finish()
	Pause = true
}

// Run the simulation for a number of steps.
func steps(n int) {
	stop := NSteps + n
	runWhile(func() bool { return NSteps < stop })
}

// Runs as long as condition returns true, saves output.
func runWhile(condition func() bool) {
	checkExchangeLenght()
	sanityCheck()
	Pause = false // may be set by <-Inject
	const output = true
	stepper.Free() // start from a clean state
	runWhileInner(condition, output)
	Pause = true
}

func runWhileInner(condition func() bool, output bool) {
	saveIfNeeded() // allow t=0 output
	for condition() && !Pause {
		select {
		default:
			step(output)
		// accept tasks form Inject channel
		case f := <-Inject:
			f()
		}
	}
}

// Runs indefinitely
func RunInteractive() {
	for {
		f := <-Inject
		f()
		time.Sleep(100 * time.Millisecond)
	}
}

// take one time step
func step(output bool) {
	stepper.Step()
	for _, f := range postStep {
		f()
	}
	if output {
		saveIfNeeded()
	}
}

// Register function f to be called after every time step.
// Typically used, e.g., to manipulate the magnetization.
func PostStep(f func()) {
	postStep = append(postStep, f)
}

func Break() {
	Inject <- func() { Pause = true }
}

// inject code into engine and wait for it to complete.
func InjectAndWait(task func()) {
	ready := make(chan int)
	Inject <- func() { task(); ready <- 1 }
	<-ready
}

func sanityCheck() {
	if Msat.isZero() {
		log.Log.Info("Note: Msat = 0")
	}
	if Aex.isZero() {
		log.Log.Info("Note: Aex = 0")
	}
}

func Exit() {
	CleanExit()
	os.Exit(0)
}

func checkExchangeLenght() {
	if exchangeLenghtWarned {
		return
	}
	// iterate over all of the quantities
	for _, region := range Regions.GetExistingIndices() {
		Msat_r := Msat.GetRegion(region)
		Aex_r := Aex.GetRegion(region)
		lex := math.Sqrt(2 * Aex_r / (mag.Mu0 * Msat_r * Msat_r))
		if !exchangeLenghtWarned {
			if Dx > lex {
				log.Log.Warn("Warning: Exchange length (%.3g nm) smaller than dx (%.3g nm) in region %d", lex*1e9, Dx*1e9, region)
				exchangeLenghtWarned = true
			}
			if Dy > lex {
				log.Log.Warn("Warning: Exchange length (%.3g nm) smaller than dy (%.3g nm) in region %d", lex*1e9, Dy*1e9, region)
				exchangeLenghtWarned = true
			}
			if Dz > lex && Nz > 1 {
				log.Log.Warn("Warning: Exchange length (%.3g nm) smaller than dz (%.3g nm) in region %d", lex*1e9, Dz*1e9, region)
				exchangeLenghtWarned = true
			}
		}

	}

}
