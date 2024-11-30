package cli

import (
	"fmt"
	"os"

	"github.com/MathieuMoalic/amumax/src/engine"
	"github.com/MathieuMoalic/amumax/src/entrypoint_old"
	"github.com/MathieuMoalic/amumax/src/flags"
	"github.com/MathieuMoalic/amumax/src/template"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func Entrypoint() {
	// default command
	cmdflags := &flags.Flags{}
	rootCmd := &cobra.Command{
		Use:   "amumax [mx3 paths...]",
		Short: "Amumax, a micromagnetic simulator",
		Args:  cobra.ArbitraryArgs,
		Run: func(cmd *cobra.Command, args []string) {
			cliEntrypoint(cmd, args, cmdflags)
		},
	}
	cmdflags.ParseFlags(rootCmd)

	// Define the template subcommand
	templateFlags := &flags.TemplateFlags{}
	templateCmd := &cobra.Command{
		Use:   "template [template path]",
		Short: "Generate files based on a template",
		Args:  cobra.ExactArgs(1), // expects exactly one argument, the template path
		Run: func(cmd *cobra.Command, args []string) {
			templateEntrypoint(args[0], templateFlags)
		},
	}
	templateFlags.ParseFlags(templateCmd)
	rootCmd.AddCommand(templateCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func cliEntrypoint(cmd *cobra.Command, args []string, flags *flags.Flags) {
	if flags.NewEngine {
		engine.Entrypoint(cmd, args, flags)
		return
	} else {
		entrypoint_old.Entrypoint(cmd, args, flags)
	}
}

func templateEntrypoint(templatePath string, flags *flags.TemplateFlags) {
	err := template.Template(templatePath, flags)
	if err != nil {
		color.Red(fmt.Sprintf("Error processing template: %v", err))
		os.Exit(1)
	}
	color.Green("Template processed successfully")
}
