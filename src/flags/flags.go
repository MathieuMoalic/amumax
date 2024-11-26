package flags

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var Flags FlagsType

type FlagsType struct {
	Debug       bool
	Version     bool
	Vet         bool
	Update      bool
	CacheDir    string
	Gpu         int
	Interactive bool
	OutputDir   string
	SelfTest    bool
	Silent      bool
	Sync        bool
	ForceClean  bool
	SkipExists  bool
	Progress    bool
	Tunnel      string
	Insecure    bool
	NewParser   bool

	WebUIEnabled      bool
	WebUIAddress      string
	WebUIQueueAddress string
}

func ParseFlags(rootCmd *cobra.Command) {
	rootCmd.Flags().BoolVarP(&Flags.Debug, "debug", "d", false, "Debug mode")
	rootCmd.Flags().BoolVarP(&Flags.Version, "version", "v", false, "Print version")
	rootCmd.Flags().BoolVar(&Flags.Vet, "vet", false, "Check input files for errors, but don't run them")
	rootCmd.Flags().BoolVar(&Flags.Update, "update", false, "Update the amumax binary from the latest github release")
	rootCmd.Flags().StringVarP(&Flags.CacheDir, "cache", "c", fmt.Sprintf("%v/amumax_kernels", os.TempDir()), "Kernel cache directory (empty disables caching)")
	rootCmd.Flags().IntVarP(&Flags.Gpu, "gpu", "g", 0, "Specify GPU")
	rootCmd.Flags().BoolVarP(&Flags.Interactive, "interactive", "i", false, "Open interactive browser session")
	rootCmd.Flags().StringVarP(&Flags.OutputDir, "output-dir", "o", "", "Override output directory")
	rootCmd.Flags().BoolVar(&Flags.SelfTest, "paranoid", false, "Enable convolution self-test for cuFFT sanity.")
	rootCmd.Flags().BoolVarP(&Flags.Silent, "silent", "s", false, "Silent mode (backwards compatibility)")
	rootCmd.Flags().BoolVar(&Flags.Sync, "sync", false, "Synchronize all CUDA calls (debug)")
	rootCmd.Flags().BoolVarP(&Flags.ForceClean, "force-clean", "f", false, "Force start, clean existing output directory")
	rootCmd.Flags().BoolVar(&Flags.SkipExists, "skip-exist", false, "Don't run the simulation if the output directory exists")
	rootCmd.Flags().BoolVar(&Flags.Progress, "progress", true, "Show progress bar")
	rootCmd.Flags().StringVarP(&Flags.Tunnel, "tunnel", "t", "", "Tunnel the web interface through SSH using the given host from your ssh config, empty string disables tunneling")
	rootCmd.Flags().BoolVar(&Flags.Insecure, "insecure", false, "Allows to run shell commands")
	rootCmd.Flags().BoolVarP(&Flags.NewParser, "new-parser", "p", false, "New parser")

	rootCmd.Flags().BoolVar(&Flags.WebUIEnabled, "webui-enable", true, "Whether to enable the web interface")
	rootCmd.Flags().StringVar(&Flags.WebUIAddress, "webui-addr", "localhost:35367", "Address (URI) to serve web GUI (e.g., 0.0.0.0:8080/proxy/worker1)")
	rootCmd.Flags().StringVar(&Flags.WebUIQueueAddress, "webui-queue-addr", "localhost:35366", "Address (URI) to serve Queue web GUI (e.g., 0.0.0.0:8080/proxy/worker1)")
}
