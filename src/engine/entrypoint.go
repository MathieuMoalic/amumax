package engine

import (
	"github.com/MathieuMoalic/amumax/src/cuda"
	"github.com/MathieuMoalic/amumax/src/flags"
	"github.com/MathieuMoalic/amumax/src/log"
	"github.com/MathieuMoalic/amumax/src/queue"
	"github.com/MathieuMoalic/amumax/src/slurm"
	"github.com/MathieuMoalic/amumax/src/timer"
	"github.com/MathieuMoalic/amumax/src/update"
	"github.com/MathieuMoalic/amumax/src/version"
	"github.com/spf13/cobra"
)

func Entrypoint(cmd *cobra.Command, args []string, givenFlags *flags.Flags) {
	// we create the log as early as possible to catch all messages
	log := log.NewLogs(givenFlags.Debug)

	if givenFlags.Update {
		update.ShowUpdateMenu()
		return
	}
	GpuInfo := cuda.Init(givenFlags.Gpu)

	cuda.Synchronous = givenFlags.Sync
	timer.Enabled = givenFlags.Sync

	log.PrintVersion(version.VERSION, GpuInfo)
	if givenFlags.Version {
		return
	}
	// engine.Insecure = givenFlags.Insecure

	if givenFlags.Vet {
		return
	}

	go slurm.SetEndTimerIfSlurm()

	if len(args) == 0 && givenFlags.Interactive {
		engineState := newEngineState(givenFlags, log)
		engineState.start("") // interactive
	} else if len(args) == 1 {
		engineState := newEngineState(givenFlags, log)
		engineState.start(args[0])
	} else if len(args) > 1 {
		queue.RunQueue(args, givenFlags)
	} else {
		_ = cmd.Help()
	}
}
