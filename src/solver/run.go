package solver

import (
	"fmt"
	"math"
	"time"

	"github.com/MathieuMoalic/amumax/src/cuda"
	"github.com/MathieuMoalic/amumax/src/log"
	"github.com/MathieuMoalic/amumax/src/magnetization"
	"github.com/MathieuMoalic/amumax/src/mesh"
	"github.com/MathieuMoalic/amumax/src/progressbar"
	"github.com/MathieuMoalic/amumax/src/regions"
	"github.com/MathieuMoalic/amumax/src/slice"
	"github.com/MathieuMoalic/amumax/src/torque"
)

// START OF TODO
// TODO: implement temperature
type Temperature interface {
	isZero() bool
	GetRegion(region int) float64
}

var Temp Temperature

var Msat, Aex Temperature

// var Msat = newScalarParam("Msat", "A/m", "Saturation magnetization")
// var Aex = newScalarParam("Msat", "A/m", "Saturation magnetization")

// Aex    = newScalarParam("Aex", "J/m", "Exchange stiffness", &lex2)

// TODO: implement saveIfNeeded
func saveIfNeeded() {}

// END OF TODO

type Solver struct {
	regions              *regions.Regions
	log                  *log.Logs
	mesh                 *mesh.Mesh
	magnetization        *magnetization.Magnetization
	torque               *torque.Torque
	Time                 float64      // Current time in seconds
	alarm                float64      // End time for the run, dt adaptation must not cross it
	pause                bool         // Set to true to stop running after the current step
	postStep             []func()     // Functions to call after every full time step
	inject               chan func()  // Injects code between time steps
	dt_si                float64      //s.time step in seconds
	MinDt, MaxDt         float64      // Minimum and maximum time steps
	MaxErr               float64      // Maximum error per step
	headroom             float64      // Solver headroom, (Gustafsson, 1992)
	lastErr, peakErr     float64      // Error of last step, highest error ever
	lastTorque           float64      // Maximum torque of last time step
	NSteps, nUndone      int          // Number of successful steps and undone steps
	nEvals               int          // Number of evaluations
	FixDt                float64      // Fixed time step (if any)
	solverType           int          // Identifier for the solver type
	exchangeLengthWarned bool         // Whether the exchange length warning has been issued
	previousStepBuffer   *slice.Slice // used by backwardEuler, rk23 and rk45DP
	precess              bool         // Precession of the magnetization
	gammaLL              float64
}

// NewSolver creates a new instance of the solver with default settings.
func (s *Solver) Init(log *log.Logs, regions *regions.Regions, mesh *mesh.Mesh, magnetization *magnetization.Magnetization, torque *torque.Torque) {
	s.log = log
	s.regions = regions
	s.mesh = mesh
	s.magnetization = magnetization
	s.torque = torque
	s.Time = 0
	s.alarm = 0
	s.pause = true
	s.postStep = []func(){}
	s.inject = make(chan func())
	s.dt_si = 1e-15
	s.MinDt = 0
	s.MaxDt = 0
	s.MaxErr = 1e-5
	s.headroom = 0.8
	s.lastErr = 0
	s.peakErr = 0
	s.lastTorque = 0
	s.NSteps = 0
	s.nUndone = 0
	s.nEvals = 0
	s.FixDt = 0
	s.solverType = 1
	s.exchangeLengthWarned = false
	s.previousStepBuffer = nil
	s.precess = true
}

func (s *Solver) SetSolver(solverIndex int) {
	// free previous solver buffer, if any
	if s.previousStepBuffer != nil {
		s.previousStepBuffer.Free()
		s.previousStepBuffer = nil
	}
	// check if solverIndex is valid
	if solverIndex < -1 || solverIndex > 6 {
		s.log.ErrAndExit("SetSolver: unknown solver type:  %v", solverIndex)
	}
	s.solverType = solverIndex
}

// write torque to dst and increment NEvals
func (s *Solver) calculateTorqueAndIncrementEvals(dst *slice.Slice) {
	s.torque.SetTorque(dst)
	s.nEvals++
}

// update lastErr and peakErr
func (s *Solver) setLastErr(err float64) {
	s.lastErr = err
	if err > s.peakErr {
		s.peakErr = err
	}
}

func (s *Solver) setMaxTorque(τ *slice.Slice) {
	s.lastTorque = cuda.MaxVecNorm(τ)
}

// adapt time step: dt *= corr, but limited to sensible values.
func (s *Solver) adaptDt(corr float64) {
	if s.FixDt != 0 {
		s.dt_si = s.FixDt
		return
	}

	// corner case triggered by err = 0: just keep time step.
	// see test/regression017.mx3
	if math.IsNaN(corr) {
		corr = 1
	}

	s.log.AssertMsg(corr != 0, "Time step too small, check if parameters are sensible")
	corr *= s.headroom
	if corr > 2 {
		corr = 2
	}
	if corr < 1./2. {
		corr = 1. / 2.
	}
	s.dt_si *= corr
	if s.MinDt != 0 && s.dt_si < s.MinDt {
		s.dt_si = s.MinDt
	}
	if s.MaxDt != 0 && s.dt_si > s.MaxDt {
		s.dt_si = s.MaxDt
	}
	if s.dt_si == 0 {
		s.log.ErrAndExit("time step too small")
	}

	// do not cross alarm time
	if s.Time < s.alarm && s.Time+s.dt_si > s.alarm {
		s.dt_si = s.alarm - s.Time
	}

	s.log.AssertMsg(s.dt_si > 0, fmt.Sprint("Time step too small: ", s.dt_si))
}

// Run the simulation for a number of seconds.
func (s *Solver) RunWithoutPrecession(seconds float64) {
	prevPrecess := s.precess
	s.Run(seconds)
	s.precess = prevPrecess
}

func (s *Solver) freeBuffer() {
	if s.previousStepBuffer != nil {
		s.previousStepBuffer.Free()
	}
}

// Run the simulation for a number of seconds.
func (s *Solver) Run(seconds float64) {
	// s.checkExchangeLength()
	start := s.Time
	stop := s.Time + seconds
	s.alarm = stop // don't have dt adapt to go over alarm
	s.sanityCheck()
	s.pause = false // may be set by <-Inject
	const output = true
	s.freeBuffer() // start from a clean state

	saveIfNeeded() // allow t=0 output
	ProgressBar := progressbar.NewProgressBar(start, stop, "🧲", true)

	for (s.Time < stop) && !s.pause {
		select {
		default:
			ProgressBar.Update(s.Time)
			s.step(output)
		// accept tasks form Inject channel
		case f := <-s.inject:
			f()
		}
	}
	ProgressBar.Finish()
	s.pause = true
}

// Run the simulation for a number of Steps.
func (s *Solver) Steps(n int) {
	stop := s.NSteps + n
	s.RunWhile(func() bool { return s.NSteps < stop })
}

// Runs as long as condition returns true, saves output.
func (s *Solver) RunWhile(condition func() bool) {
	// s.checkExchangeLength()
	s.sanityCheck()
	s.pause = false // may be set by <-Inject
	const output = true
	s.freeBuffer() // start from a clean state
	s.runWhileInner(condition, output)
	s.pause = true
}

func (s *Solver) runWhileInner(condition func() bool, output bool) {
	saveIfNeeded() // allow t=0 output
	for condition() && !s.pause {
		select {
		default:
			s.step(output)
		// accept tasks form Inject channel
		case f := <-s.inject:
			f()
		}
	}
}

// Runs indefinitely
func (s *Solver) RunInteractive() {
	for {
		f := <-s.inject
		f()
		time.Sleep(100 * time.Millisecond)
	}
}

// take one time step
func (s *Solver) step(output bool) {
	switch s.solverType {
	default:
		s.log.ErrAndExit("Step: unknown solver type:  %v", s.solverType)
	case -1:
		s.backWardEulerStep()
	case 1:
		s.euler()
	case 2:
		s.heun()
	case 3:
		s.rk23()
	case 4:
		s.rk4()
	case 5:
		s.rk45()
	case 6:
		s.rk56()
	}
	for _, f := range s.postStep {
		f()
	}
	if output {
		saveIfNeeded()
	}
}

// Register function f to be called after every time step.
// Typically used, e.g., to manipulate the magnetization.
func (s *Solver) PostStep(f func()) {
	s.postStep = append(s.postStep, f)
}

func (s *Solver) Break() {
	s.inject <- func() { s.pause = true }
}

// inject code into engine and wait for it to complete.
func (s *Solver) InjectAndWait(task func()) {
	ready := make(chan int)
	s.inject <- func() { task(); ready <- 1 }
	<-ready
}

func (s *Solver) sanityCheck() {
	// if Msat.isZero() {
	// 	s.log.Info("Note: Msat = 0")
	// }
	// if Aex.isZero() {
	// 	s.log.Info("Note: Aex = 0")
	// }
}

// func (s *Solver) checkExchangeLength() {
// 	if s.exchangeLengthWarned {
// 		return
// 	}
// 	existingRegions := s.regions.GetExistingIndices()
// 	// iterate over all of the quantities
// 	for _, region := range existingRegions {
// 		Msat_r := Msat.GetRegion(region)
// 		Aex_r := Aex.GetRegion(region)
// 		lex := math.Sqrt(2 * Aex_r / (constants.Mu0 * Msat_r * Msat_r))
// 		if !s.exchangeLengthWarned {
// 			if s.mesh.Dx > lex {
// 				s.log.Warn("Warning: Exchange length (%.3g nm) smaller than dx (%.3g nm) in region %d", lex*1e9, s.mesh.Dx*1e9, region)
// 				s.exchangeLengthWarned = true
// 			}
// 			if s.mesh.Dy > lex {
// 				s.log.Warn("Warning: Exchange length (%.3g nm) smaller than dy (%.3g nm) in region %d", lex*1e9, s.mesh.Dy*1e9, region)
// 				s.exchangeLengthWarned = true
// 			}
// 			if s.mesh.Dz > lex && s.mesh.Nz > 1 {
// 				s.log.Warn("Warning: Exchange length (%.3g nm) smaller than dz (%.3g nm) in region %d", lex*1e9, s.mesh.Dz*1e9, region)
// 				s.exchangeLengthWarned = true
// 			}
// 		}

// 	}

// }
