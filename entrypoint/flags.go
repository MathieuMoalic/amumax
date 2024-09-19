package entrypoint

import (
	"fmt"
	"os"
)

type Flags struct {
	debug       bool
	test        bool
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
	tunnel      bool

	webUIEnabled   bool
	webUIHost      string
	webUIPort      int
	webUIQueueHost string
	webUIQueuePort int
}

var flags Flags

func init() {
	rootCmd.Flags().BoolVarP(&flags.debug, "debug", "d", false, "Debug mode")
	rootCmd.Flags().BoolVar(&flags.test, "test", false, "Cuda test (internal)")
	rootCmd.Flags().BoolVarP(&flags.version, "version", "v", true, "Print version")
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
	rootCmd.Flags().BoolVarP(&flags.tunnel, "tunnel", "t", false, "Tunnel the web interface through SSH")

	rootCmd.Flags().BoolVar(&flags.webUIEnabled, "webui-enable", true, "Whether to enable the web interface")
	rootCmd.Flags().StringVar(&flags.webUIHost, "webui-host", "localhost", "Host to serve web GUI (e.g., 0.0.0.0)")
	rootCmd.Flags().IntVar(&flags.webUIPort, "webui-port", 35367, "Port to serve web GUI")
	rootCmd.Flags().StringVar(&flags.webUIQueueHost, "webui-queue-host", "localhost", "Host to serve the queue web GUI (e.g., 0.0.0.0)")
	rootCmd.Flags().IntVar(&flags.webUIQueuePort, "webui-queue-port", 35366, "Port to serve queue web GUI")
}