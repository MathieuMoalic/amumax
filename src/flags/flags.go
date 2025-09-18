// Package flags provides command-line flag parsing for the application.
package flags

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

type Flags struct {
	Debug           bool
	Version         bool
	Vet             bool
	Update          bool
	CacheDir        string
	Gpu             int
	Interactive     bool
	OutputDir       string
	SelfTest        bool
	Silent          bool
	Sync            bool
	ForceClean      bool
	SkipExists      bool
	HideProgressBar bool
	Tunnel          string
	Insecure        bool
	NewEngine       bool

	WebUIDisabled     bool
	WebUIAddress      string
	WebUIQueueAddress string
}

func (flags *Flags) ParseFlags(rootCmd *cobra.Command) {
	rootCmd.Flags().BoolVarP(&flags.Debug, "debug", "d", false, "Debug mode")
	rootCmd.Flags().BoolVarP(&flags.Version, "version", "v", false, "Print version")
	rootCmd.Flags().BoolVar(&flags.Vet, "vet", false, "Check input files for errors, but don't run them")
	rootCmd.Flags().BoolVarP(&flags.Update, "update", "u", false, "Update the amumax binary from the latest github release")
	rootCmd.Flags().StringVarP(&flags.CacheDir, "cache", "c", fmt.Sprintf("%v/amumax_kernels", os.TempDir()), "Kernel cache directory (empty disables caching)")
	rootCmd.Flags().IntVarP(&flags.Gpu, "gpu", "g", 0, "Specify GPU")
	rootCmd.Flags().BoolVarP(&flags.Interactive, "interactive", "i", false, "Open interactive browser session")
	rootCmd.Flags().StringVarP(&flags.OutputDir, "output-dir", "o", "", "Override output directory")
	rootCmd.Flags().BoolVarP(&flags.SelfTest, "paranoid", "p", false, "Enable convolution self-test for cuFFT sanity.")
	rootCmd.Flags().BoolVarP(&flags.Silent, "silent", "s", false, "Silent mode (backwards compatibility)")
	rootCmd.Flags().BoolVar(&flags.Sync, "sync", false, "Synchronize all CUDA calls (debug)")
	rootCmd.Flags().BoolVarP(&flags.ForceClean, "force-clean", "f", false, "Force start, clean existing output directory")
	rootCmd.Flags().BoolVar(&flags.SkipExists, "skip-exist", false, "Don't run the simulation if the output directory exists")
	rootCmd.Flags().BoolVar(&flags.HideProgressBar, "hide-progress-bar", false, "Hide the progress bar")
	rootCmd.Flags().StringVarP(&flags.Tunnel, "tunnel", "t", "", "Tunnel the web interface through SSH using the given host from your ssh config, empty string disables tunneling")
	rootCmd.Flags().BoolVar(&flags.Insecure, "insecure", false, "Allows to run shell commands")
	rootCmd.Flags().BoolVarP(&flags.NewEngine, "new-engine", "n", false, "New engine, experimental")

	rootCmd.Flags().BoolVar(&flags.WebUIDisabled, "webui-disable", false, "Whether to disable the web interface")
	rootCmd.Flags().StringVar(&flags.WebUIAddress, "webui-addr", "localhost:35367", "Address (URI) to serve web GUI (e.g., 0.0.0.0:8080/proxy/worker1)")
	rootCmd.Flags().StringVar(&flags.WebUIQueueAddress, "webui-queue-addr", "localhost:35366", "Address (URI) to serve Queue web GUI (e.g., 0.0.0.0:8080/proxy/worker1)")
}

type TemplateFlags struct {
	Flat bool
	Run  bool
}

func (flags *TemplateFlags) ParseFlags(templateCmd *cobra.Command) {
	templateCmd.Flags().BoolVar(&flags.Flat, "flat", false, "Generate flat output without subdirectories")
	templateCmd.Flags().BoolVar(&flags.Run, "run", false, "Run the generated script")
}
