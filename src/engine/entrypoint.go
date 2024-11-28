package engine

import (
	"runtime"

	"github.com/MathieuMoalic/amumax/src/cuda"
	"github.com/MathieuMoalic/amumax/src/cuda/cu"
	"github.com/MathieuMoalic/amumax/src/engine_old"
	"github.com/MathieuMoalic/amumax/src/flags"
	"github.com/MathieuMoalic/amumax/src/log_old"
	"github.com/MathieuMoalic/amumax/src/queue"
	"github.com/MathieuMoalic/amumax/src/slurm"
	"github.com/MathieuMoalic/amumax/src/timer"
	"github.com/MathieuMoalic/amumax/src/update"
	"github.com/spf13/cobra"
)

func Entrypoint(cmd *cobra.Command, args []string, givenFlags *flags.FlagsType) {
	if givenFlags.Update {
		update.ShowUpdateMenu()
		return
	}
	go slurm.SetEndTimerIfSlurm()
	cuda.Init(givenFlags.Gpu)

	cuda.Synchronous = givenFlags.Sync
	timer.Enabled = givenFlags.Sync

	printVersion()
	if givenFlags.Version {
		return
	}
	// engine.Insecure = givenFlags.Insecure

	if givenFlags.Vet {
		return
	}
	if len(args) == 0 && givenFlags.Interactive {
		engineState := newEngineState(givenFlags)
		engineState.startInteractive()

	} else if len(args) == 1 {
		engineState := newEngineState(givenFlags)
		engineState.start(args[0])
	} else if len(args) > 1 {
		queue.RunQueue(args, givenFlags)
	} else {
		_ = cmd.Help()
	}
}

// print version to stdout
func printVersion() {
	log_old.Log.Info("Version:         %s", engine_old.VERSION)
	log_old.Log.Info("Platform:        %s_%s", runtime.GOOS, runtime.GOARCH)
	log_old.Log.Info("Go Version:      %s (%s)", runtime.Version(), runtime.Compiler)
	log_old.Log.Info("CUDA Version:    %d.%d (CC=%d PTX)", cu.CUDA_VERSION/1000, (cu.CUDA_VERSION%1000)/10, cuda.UseCC)
	log_old.Log.Info("GPU Information: %s", cuda.GPUInfo)
}
