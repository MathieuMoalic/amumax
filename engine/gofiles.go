package engine

// support for running Go files as if they were mx3 files.

import (
	"flag"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/MathieuMoalic/amumax/cuda"
)

var (
	// These flags are shared between cmd/mumax3 and Go input files.
	Flag_cachedir    = flag.String("cache", fmt.Sprintf("%v/amumax_kernels", os.TempDir()), "Kernel cache directory (empty disables caching)")
	Flag_gpu         = flag.Int("gpu", 0, "Specify GPU")
	Flag_interactive = flag.Bool("i", false, "Open interactive browser session")
	Flag_od          = flag.String("o", "", "Override output directory")
	Flag_selftest    = flag.Bool("paranoid", false, "Enable convolution self-test for cuFFT sanity.")
	Flag_silent      = flag.Bool("s", false, "Silent") // provided for backwards compatibility
	Flag_sync        = flag.Bool("sync", false, "Synchronize all CUDA calls (debug)")
	Flag_forceclean  = flag.Bool("f", false, "Force start, clean existing output directory")
	Flag_skip_exists = flag.Bool("skip-exist", false, "Don't run the simulation if the output directory exists")
	Flag_magnets     = flag.Bool("magnets", true, "Show magnets progress bar")

	Flag_webui_enabled    = flag.Bool("webui-enable", true, "Whether to enable the web interface")
	Flag_webui_host       = flag.String("webui-host", "localhost", "Host to serve web gui i.e. 0.0.0.0")
	Flag_webui_port       = flag.Int("webui-port", 35367, "Port to serve web gui")
	Flag_webui_queue_host = flag.String("webui-queue-host", "localhost", "Host to serve the queue web gui i.e. 0.0.0.0")
	Flag_webui_queue_port = flag.Int("webui-queue-port", 35366, "Port to serve queue web gui")
)

// Usage: in every Go input file, write:
//
//	func main(){
//		defer InitAndClose()()
//		// ...
//	}
//
// This initialises the GPU, output directory, etc,
// and makes sure pending output will get flushed.
func InitAndClose() func() {
	// ONLY FOR GO FILES
	flag.Parse()

	cuda.Init(*Flag_gpu)
	cuda.Synchronous = *Flag_sync

	od := *Flag_od
	if od == "" {
		od = path.Base(os.Args[0]) + ".zarr"
	}
	inFile := strings.Replace(od, ".zarr", ".mx3", 1)
	InitIO(inFile, od)

	return func() {
		CleanExit()
	}
}
