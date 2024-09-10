/*
engine does the simulation bookkeeping, I/O and GUI.

space-dependence:
value: space-independent
param: region-dependent parameter (always input)
field: fully space-dependent field

TODO: godoc everything
*/
package engine

import (
	"os"
	"sync"
	"time"

	"github.com/MathieuMoalic/amumax/script"
	"github.com/MathieuMoalic/amumax/timer"
	"github.com/MathieuMoalic/amumax/util"
)

var VERSION = "dev"

var StartTime = time.Now()

var (
	busyLock sync.Mutex
	busy     bool // are we so busy we can't respond from run loop? (e.g. calc kernel)
)

// We set SetBusy(true) when the simulation is too busy too accept GUI input on Inject channel.
// E.g. during kernel init.
func SetBusy(b bool) {
	busyLock.Lock()
	defer busyLock.Unlock()
	busy = b
}

func GetBusy() bool {
	busyLock.Lock()
	defer busyLock.Unlock()
	return busy
}

// Cleanly exits the simulation, assuring all output is flushed.
func CleanExit() {
	if outputdir == "" {
		return
	}
	drainOutput()
	util.Log.Info("**************** Simulation Ended ****************** //")
	Table.Flush()
	if SyncAndLog {
		timer.Print(os.Stdout)
	}
	script.MMetadata.Add("steps", NSteps)
	script.MMetadata.End()
	util.Log.FlushToFile()
}
