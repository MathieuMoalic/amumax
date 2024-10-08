package entrypoint

import (
	"fmt"
	"os"

	"github.com/MathieuMoalic/amumax/src/cuda"
	"github.com/MathieuMoalic/amumax/src/engine"
	"github.com/MathieuMoalic/amumax/src/log"
	"github.com/MathieuMoalic/amumax/src/timer"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var flags flagsType

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
			// Call your template logic here, using args[0] for the template path
			templatePath := args[0]
			err := template(templatePath)
			if err != nil {
				color.Red(fmt.Sprintf("Error processing template: %v", err))
				os.Exit(1)
			}
			color.Green("Template processed successfully")
		},
	}
	rootCmd.AddCommand(templateCmd)
	parseFlags(rootCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func cliEntrypoint(cmd *cobra.Command, args []string) {
	log.Log.SetDebug(flags.debug)
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
		runInteractive(&flags)
	} else if len(args) == 1 {
		runFileAndServe(args[0], &flags)
	} else if len(args) > 1 {
		RunQueue(args, &flags)
	} else {
		log.Log.ErrAndExit("No input files")
	}
}
