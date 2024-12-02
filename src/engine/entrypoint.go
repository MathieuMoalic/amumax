package engine

import (
	"github.com/MathieuMoalic/amumax/src/engine_old/cuda_old"
	"github.com/MathieuMoalic/amumax/src/engine_old/timer_old"
	"github.com/MathieuMoalic/amumax/src/flags"
	"github.com/MathieuMoalic/amumax/src/log"
	"github.com/MathieuMoalic/amumax/src/update"
	"github.com/MathieuMoalic/amumax/src/version"
	"github.com/spf13/cobra"
)

// Entrypoint is the entrypoint for the new engine which is not very functional yet
// the cuda code still relies on global variables
func Entrypoint(cmd *cobra.Command, args []string, givenFlags *flags.Flags) {
	// we create the log as early as possible to catch all messages
	log := log.NewLogs(givenFlags.Debug)

	if givenFlags.Update {
		update.ShowUpdateMenu()
		return
	}
	GpuInfo := cuda_old.Init(givenFlags.Gpu)

	cuda_old.Synchronous = givenFlags.Sync
	timer_old.Enabled = givenFlags.Sync

	// engine.Insecure = givenFlags.Insecure

	if givenFlags.Vet {
		log.PrintVersion(version.VERSION, GpuInfo)
		log.Err("vet is not implemented yet with the new engine")
	} else if len(args) == 0 && givenFlags.Interactive {
		log.PrintVersion(version.VERSION, GpuInfo)
		engineState := newEngineState(givenFlags, log)
		engineState.start("") // interactive
	} else if len(args) == 1 {
		log.PrintVersion(version.VERSION, GpuInfo)
		engineState := newEngineState(givenFlags, log)
		engineState.start(args[0])
	} else if len(args) > 1 {
		log.Err("Queue is not implemented yet with the new engine")
	} else if givenFlags.Version {
		log.PrintVersion(version.VERSION, GpuInfo)
	} else {
		_ = cmd.Help()
	}
}
