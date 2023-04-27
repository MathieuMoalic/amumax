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
	"fmt"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/MathieuMoalic/amumax/cuda/cu"
	"github.com/MathieuMoalic/amumax/timer"
)

var VERSION = "NOT_SET"

var UNAME = fmt.Sprintf("Amumax v.%s a fork of mumax 3.10 [%s_%s %s(%s) CUDA-%d.%d]",
	VERSION, runtime.GOOS, runtime.GOARCH, runtime.Version(), runtime.Compiler,
	cu.CUDA_VERSION/1000, (cu.CUDA_VERSION%1000)/10)

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
func Close() {
	if outputdir == "" {
		return
	}
	drainOutput()
	LogOut("**************** Simulation Ended ****************** //")
	ZTables.Flush()
	if logfile != nil {
		logfile.Close()
	}
	// newlogfile.Flush()
	// newlogfilefile.Close()
	// if newlogfilefile != nil {
	// }
	if *Flag_sync {
		timer.Print(os.Stdout)
	}
	EndMetadata()
}
