package engine

// support for running Go files as if they were mx3 files.

import (
	"flag"
	"os"
)

var (
	// These flags are shared between cmd/mumax3 and Go input files.
	Flag_cachedir         = flag.String("cache", os.TempDir(), "Kernel cache directory (empty disables caching)")
	Flag_gpu              = flag.Int("gpu", 0, "Specify GPU")
	Flag_interactive      = flag.Bool("i", false, "Open interactive browser session")
	Flag_od               = flag.String("o", "", "Override output directory")
	Flag_webui_addr       = flag.String("http", "0.0.0.0:35367", "Address to serve web gui")
	Flag_webui_queue_addr = flag.String("qhttp", "0.0.0.0:35366", "Address to serve web gui")
	Flag_selftest         = flag.Bool("paranoid", false, "Enable convolution self-test for cuFFT sanity.")
	Flag_silent           = flag.Bool("s", false, "Silent") // provided for backwards compatibility
	Flag_sync             = flag.Bool("sync", false, "Synchronize all CUDA calls (debug)")
	Flag_forceclean       = flag.Bool("f", false, "Force start, clean existing output directory")
	Flag_skip_exists      = flag.Bool("skip-exist", false, "Don't run the simulation if the output directory exists ( if the simulation has been run before )")
	Flag_magnets          = flag.Bool("magnets", true, "Show magnets progress bar")
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
// func InitAndClose() func() {
// 	// ONLY FOR GO FILES
// 	flag.Parse()

// 	cuda.Init(*Flag_gpu)
// 	cuda.Synchronous = *Flag_sync

// 	od := *Flag_od
// 	if od == "" {
// 		od = path.Base(os.Args[0]) + ".zarr"
// 	}
// 	inFile := util.NoExt(od)
// 	InitIO(inFile, od)

// if *Flag_webui_addr == "" {
// 	util.LogWarn(`WebUI is disabled (-http="")`)
// }
// addr, err := findAvailablePort(*Flag_webui_addr)
// if err != nil {
// 	log.Fatalf("Failed to find available port: %v", err)
// }
// util.Log(fmt.Sprintf("Serving GUI at http://%s", addr))
// go api.Start(addr)

// 	return func() {
// 		CleanExit()
// 	}
// }
