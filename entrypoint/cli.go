package entrypoint

import (
	"fmt"
	"os"

	"github.com/MathieuMoalic/amumax/cuda"
	"github.com/MathieuMoalic/amumax/engine"
	"github.com/MathieuMoalic/amumax/timer"
	"github.com/MathieuMoalic/amumax/util"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "amumax [mx3 paths...]",
	Short: "Amumax, a micromagnetic simulator",
	Run:   cliEntrypoint,
	Args:  cobra.ArbitraryArgs,
}

func cliEntrypoint(cmd *cobra.Command, args []string) {
	util.Log.SetDebug(flags.debug)
	if flags.update {
		showUpdateMenu()
		return
	}
	go setEndTimerIfSlurm()
	cuda.Init(flags.gpu)

	cuda.Synchronous = flags.sync
	timer.Enabled = flags.sync

	if flags.version {
		printVersion()
	}
	engine.Insecure = flags.insecure

	// used by bootstrap launcher to test cuda
	// successful exit means cuda was initialized fine
	if flags.test {
		fmt.Println(cuda.GPUInfo)
		os.Exit(0)
	}

	defer engine.CleanExit() // flushes pending output, if any

	if flags.vet {
		vet()
		return
	}
	if len(args) == 0 && flags.interactive {
		runInteractive()
	} else if len(args) == 1 {
		runFileAndServe(args[0])
	} else if len(args) > 1 {
		RunQueue(args)
	} else {
		util.Log.ErrAndExit("No input files")
	}
}
