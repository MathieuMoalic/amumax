package entrypoint

import (
	"fmt"
	"os"

	"github.com/MathieuMoalic/amumax/src/cuda"
	"github.com/MathieuMoalic/amumax/src/engine"
	"github.com/MathieuMoalic/amumax/src/engine_old"
	"github.com/MathieuMoalic/amumax/src/flags"
	"github.com/MathieuMoalic/amumax/src/log_old"
	"github.com/MathieuMoalic/amumax/src/queue"
	"github.com/MathieuMoalic/amumax/src/slurm"
	"github.com/MathieuMoalic/amumax/src/template"
	"github.com/MathieuMoalic/amumax/src/timer"
	"github.com/MathieuMoalic/amumax/src/update"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var flat bool

func Entrypoint() {
	rootCmd := &cobra.Command{
		Use:   "amumax [mx3 paths...]",
		Short: "Amumax, a micromagnetic simulator",
		Run:   cliEntrypoint,
		Args:  cobra.ArbitraryArgs,
	}

	// Define the template subcommand
	templateCmd := &cobra.Command{
		Use:   "template [template path]",
		Short: "Generate files based on a template",
		Args:  cobra.ExactArgs(1), // expects exactly one argument, the template path
		Run: func(cmd *cobra.Command, args []string) {
			// Get the value of the "flat" flag
			flat, _ = cmd.Flags().GetBool("flat")

			// Call your template logic here, using args[0] for the template path
			templatePath := args[0]
			err := template.Template(templatePath, flat)
			if err != nil {
				color.Red(fmt.Sprintf("Error processing template: %v", err))
				os.Exit(1)
			}
			color.Green("Template processed successfully")
		},
	}

	// Add the "flat" flag to the template command
	templateCmd.Flags().BoolVar(&flat, "flat", false, "Generate flat output without subdirectories")

	rootCmd.AddCommand(templateCmd)
	flags.ParseFlags(rootCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func cliEntrypoint(cmd *cobra.Command, args []string) {
	if flags.Flags.NewEngine {
		engine.Entrypoint(cmd, args, &flags.Flags)
		return
	}
	log_old.Log.SetDebug(flags.Flags.Debug)
	if flags.Flags.Update {
		update.ShowUpdateMenu()
		return
	}
	go slurm.SetEndTimerIfSlurm()
	cuda.Init(flags.Flags.Gpu)

	cuda.Synchronous = flags.Flags.Sync
	timer.Enabled = flags.Flags.Sync

	printVersion()
	if flags.Flags.Version {
		return
	}
	engine_old.Insecure = flags.Flags.Insecure

	defer engine_old.CleanExit() // flushes pending output, if any

	if flags.Flags.Vet {
		vet()
		return
	}
	if len(args) == 0 && flags.Flags.Interactive {
		runInteractive(&flags.Flags)
	} else if len(args) == 1 {
		runFileAndServe(args[0], &flags.Flags)
	} else if len(args) > 1 {
		queue.RunQueue(args, &flags.Flags)
	} else {
		_ = cmd.Help()
	}
}
