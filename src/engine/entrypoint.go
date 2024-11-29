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

func Entrypoint(cmd *cobra.Command, args []string, givenFlags *flags.FlagsType) {
	if givenFlags.Update {
		update.ShowUpdateMenu()
		return
	}
	go slurm.SetEndTimerIfSlurm()
	cuda.Init(givenFlags.Gpu)

	cuda.Synchronous = givenFlags.Sync
	timer.Enabled = givenFlags.Sync

	// we create the log as early as possible to catch all messages
	log := log.NewLogs(givenFlags.Debug)

	log.PrintVersion(version.VERSION)
	if givenFlags.Version {
		return
	}
	// engine.Insecure = givenFlags.Insecure

	if givenFlags.Vet {
		return
	}
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
