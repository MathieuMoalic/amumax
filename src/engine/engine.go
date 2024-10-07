package engine

import (
	"os"
	"sync"
	"time"

	"github.com/MathieuMoalic/amumax/src/log"
	"github.com/MathieuMoalic/amumax/src/script"
	"github.com/MathieuMoalic/amumax/src/timer"
)

var VERSION = "dev"

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
	log.Log.Info("**************** Simulation Ended ****************** //")
	Table.Flush()
	if SyncAndLog {
		timer.Print(os.Stdout)
	}
	script.MMetadata.Add("steps", NSteps)
	script.MMetadata.End()
	log.Log.FlushToFile()
}
