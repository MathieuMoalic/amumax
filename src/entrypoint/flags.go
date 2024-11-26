package entrypoint

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

type flagsType struct {
	debug       bool
	version     bool
	vet         bool
	update      bool
	cacheDir    string
	gpu         int
	interactive bool
	outputDir   string
	selfTest    bool
	silent      bool
	sync        bool
	forceClean  bool
	skipExists  bool
	progress    bool
	tunnel      string
	insecure    bool

	webUIEnabled      bool
	webUIAddress      string
	webUIQueueAddress string
}

func parseFlags(rootCmd *cobra.Command) {
	rootCmd.Flags().BoolVarP(&flags.debug, "debug", "d", false, "Debug mode")
	rootCmd.Flags().BoolVarP(&flags.version, "version", "v", false, "Print version")
	rootCmd.Flags().BoolVar(&flags.vet, "vet", false, "Check input files for errors, but don't run them")
	rootCmd.Flags().BoolVar(&flags.update, "update", false, "Update the amumax binary from the latest github release")
	rootCmd.Flags().StringVarP(&flags.cacheDir, "cache", "c", fmt.Sprintf("%v/amumax_kernels", os.TempDir()), "Kernel cache directory (empty disables caching)")
	rootCmd.Flags().IntVarP(&flags.gpu, "gpu", "g", 0, "Specify GPU")
	rootCmd.Flags().BoolVarP(&flags.interactive, "interactive", "i", false, "Open interactive browser session")
	rootCmd.Flags().StringVarP(&flags.outputDir, "output-dir", "o", "", "Override output directory")
	rootCmd.Flags().BoolVar(&flags.selfTest, "paranoid", false, "Enable convolution self-test for cuFFT sanity.")
	rootCmd.Flags().BoolVarP(&flags.silent, "silent", "s", false, "Silent mode (backwards compatibility)")
	rootCmd.Flags().BoolVar(&flags.sync, "sync", false, "Synchronize all CUDA calls (debug)")
	rootCmd.Flags().BoolVarP(&flags.forceClean, "force-clean", "f", false, "Force start, clean existing output directory")
	rootCmd.Flags().BoolVar(&flags.skipExists, "skip-exist", false, "Don't run the simulation if the output directory exists")
	rootCmd.Flags().BoolVar(&flags.progress, "progress", true, "Show progress bar")
	rootCmd.Flags().StringVarP(&flags.tunnel, "tunnel", "t", "", "Tunnel the web interface through SSH using the given host from your ssh config, empty string disables tunneling")
	rootCmd.Flags().BoolVar(&flags.insecure, "insecure", false, "Allows to run shell commands")

	rootCmd.Flags().BoolVar(&flags.webUIEnabled, "webui-enable", true, "Whether to enable the web interface")
	rootCmd.Flags().StringVar(&flags.webUIAddress, "webui-addr", "localhost:35367", "Address (URI) to serve web GUI (e.g., 0.0.0.0:8080/proxy/worker1)")
	rootCmd.Flags().StringVar(&flags.webUIQueueAddress, "webui-queue-addr", "localhost:35366", "Address (URI) to serve Queue web GUI (e.g., 0.0.0.0:8080/proxy/worker1)")
}
