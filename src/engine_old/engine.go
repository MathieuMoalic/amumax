package engine_old

import (
	"os"
	"sync"
	"time"

	"github.com/MathieuMoalic/amumax/src/engine_old/log_old"
	"github.com/MathieuMoalic/amumax/src/engine_old/timer_old"
)

var StartTime = time.Now()

var (
	busyLock sync.Mutex
)

// We set setBusy(true) when the simulation is too busy too accept GUI input on Inject channel.
// E.g. during kernel init.
func setBusy(_b bool) {
	// TODO is it needed?
	_ = _b
	busyLock.Lock()
	defer busyLock.Unlock()
}

// Cleanly exits the simulation, assuring all output is flushed.
func CleanExit() {
	if outputdir == "" {
		return
	}
	drainOutput()
	log_old.Log.Info("**************** Simulation Ended ****************** //")
	Table.Flush()
	if SyncAndLog {
		timer_old.Print(os.Stdout)
	}
	EngineState.Metadata.Add("steps", NSteps)
	EngineState.Metadata.End()
	log_old.Log.FlushToFile()
}
