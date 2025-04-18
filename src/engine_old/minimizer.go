package engine_old

// Minimize follows the steepest descent method as per Exl et al., JAP 115, 17D118 (2014).

import (
	"time"

	"github.com/MathieuMoalic/amumax/src/engine_old/cuda_old"
	"github.com/MathieuMoalic/amumax/src/engine_old/data_old"
	"github.com/MathieuMoalic/amumax/src/engine_old/log_old"
)

var (
	dmSamples              int     = 10   // number of dm to keep for convergence check
	stopMaxDm              float64 = 1e-6 // stop minimizer if sampled dm is smaller than this
	minimizeMaxSteps       int     = 1000000
	minimizeMaxTimeSeconds         = 60 * 60 * 24 * 7 // one week
)

// fixed length FIFO. Items can be added but not removed
type fifoRing struct {
	count int
	tail  int // index to put next item. Will loop to 0 after exceeding length
	data  []float64
}

func createFifoRing(length int) fifoRing {
	return fifoRing{data: make([]float64, length)}
}

func (r *fifoRing) Add(item float64) {
	r.data[r.tail] = item
	r.count++
	r.tail = (r.tail + 1) % len(r.data)
	if r.count > len(r.data) {
		r.count = len(r.data)
	}
}

func (r *fifoRing) Max() float64 {
	max := r.data[0]
	for i := 1; i < r.count; i++ {
		if r.data[i] > max {
			max = r.data[i]
		}
	}
	return max
}

type Minimizer struct {
	k      *data_old.Slice // torque saved to calculate time step
	lastDm fifoRing
	h      float32
}

func (mini *Minimizer) Step() {
	m := NormMag.Buffer()
	size := m.Size()

	if mini.k == nil {
		mini.k = cuda_old.Buffer(3, size)
		torqueFn(mini.k)
	}

	k := mini.k
	h := mini.h

	// save original magnetization
	m0 := cuda_old.Buffer(3, size)
	defer cuda_old.Recycle(m0)
	data_old.Copy(m0, m)

	// make descent
	cuda_old.Minimize(m, m0, k, h)

	// calculate new torque for next step
	k0 := cuda_old.Buffer(3, size)
	defer cuda_old.Recycle(k0)
	data_old.Copy(k0, k)
	torqueFn(k)
	setMaxTorque(k) // report to user

	// just to make the following readable
	dm := m0
	dk := k0

	// calculate step difference of m and k
	cuda_old.Madd2(dm, m, m0, 1., -1.)
	cuda_old.Madd2(dk, k, k0, -1., 1.) // reversed due to LLNoPrecess sign

	// get maxdiff and add to list
	max_dm := cuda_old.MaxVecNorm(dm)
	mini.lastDm.Add(max_dm)
	setLastErr(mini.lastDm.Max()) // report maxDm to user as LastErr

	// adjust next time step
	var nom, div float32
	if NSteps%2 == 0 {
		nom = cuda_old.Dot(dm, dm)
		div = cuda_old.Dot(dm, dk)
	} else {
		nom = cuda_old.Dot(dm, dk)
		div = cuda_old.Dot(dk, dk)
	}
	if div != 0. {
		mini.h = nom / div
	} else { // in case of division by zero
		mini.h = 1e-4
	}

	NormMag.normalize()

	// as a convention, time does not advance during relax
	NSteps++
}

func (mini *Minimizer) Free() {
	mini.k.Free()
}

var (
	MinimizeStartTime   time.Time
	MinimizeTimeoutStep int
)

func minimize() {
	checkExchangeLength()
	MinimizeStartTime = time.Now()
	MinimizeTimeoutStep = NSteps + minimizeMaxSteps
	sanityCheck()
	// Save the settings we are changing...
	prevType := Solvertype
	prevFixDt := FixDt
	prevPrecess := precess
	t0 := Time

	relaxing = true // disable temperature noise

	// ...to restore them later
	defer func() {
		setSolver(prevType)
		FixDt = prevFixDt
		precess = prevPrecess
		Time = t0

		relaxing = false
	}()

	precess = false // disable precession for torque calculation
	// remove previous stepper
	if stepper != nil {
		stepper.Free()
	}

	// set stepper to the minimizer
	mini := Minimizer{
		h:      1e-4,
		k:      nil,
		lastDm: createFifoRing(dmSamples)}
	stepper = &mini

	cond := func() bool {
		maxStepsReached := MinimizeTimeoutStep < NSteps
		maxTimeReached := int(time.Since(MinimizeStartTime).Seconds()) > minimizeMaxTimeSeconds
		maxDmSamplesReached := mini.lastDm.count < dmSamples
		maxDmReached := mini.lastDm.Max() > stopMaxDm
		out := !(maxStepsReached || maxTimeReached || !(maxDmSamplesReached || maxDmReached))
		if maxStepsReached {
			log_old.Log.Info("Stopping `Minimize()`: Maximum time steps reached ( MinimizeMaxSteps= %v steps", minimizeMaxSteps)
		}
		if maxTimeReached {
			log_old.Log.Info("Stopping `Minimize()`: Maximum time reached ( MinimizeMaxTimeSeconds= %vs )", minimizeMaxTimeSeconds)
		}
		return out
	}
	runWhile(cond)
	Pause = true
}
