// mumax3 main command
package entrypoint

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/MathieuMoalic/amumax/src/api"
	"github.com/MathieuMoalic/amumax/src/cuda"
	"github.com/MathieuMoalic/amumax/src/cuda/cu"
	"github.com/MathieuMoalic/amumax/src/engine"
	"github.com/MathieuMoalic/amumax/src/log"
	"github.com/MathieuMoalic/amumax/src/script"
)

func Entrypoint() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

type Release struct {
	TagName string `json:"tag_name"`
}

func runInteractive() {
	log.Log.Info("No input files: starting interactive session")
	// setup outut dir
	now := time.Now()
	outdir := fmt.Sprintf("/tmp/amumax-%v-%02d-%02d_%02dh%02d.zarr", now.Year(), int(now.Month()), now.Day(), now.Hour(), now.Minute())

	engine.InitIO(outdir, outdir, flags.cacheDir, flags.skipExists, flags.forceClean, flags.progress, flags.selfTest, flags.sync)
	log.Log.Info("Input file: %s", "none")
	log.Log.Info("Output directory: %s", engine.OD())
	log.Log.Init(engine.OD())

	// set up some sensible start configuration
	engine.Eval(`
Nx = 128
Ny = 64
Nz = 1
dx = 3e-9
dy = 3e-9
dz = 3e-9
Msat = 1e6
Aex = 10e-12
alpha = 1
m = RandomMag()`)
	if flags.webUIEnabled {
		go api.Start(flags.webUIHost, flags.webUIPort, flags.tunnel, flags.debug)
	}
	engine.RunInteractive()
}

func runFileAndServe(mx3Path string) {
	if _, err := os.Stat(mx3Path); errors.Is(err, os.ErrNotExist) {
		log.Log.ErrAndExit("Error: File `%s` does not exist", mx3Path)
	}
	outputdir := strings.TrimSuffix(mx3Path, ".mx3") + ".zarr"

	if flags.outputDir != "" {
		outputdir = flags.outputDir
	}
	engine.InitIO(mx3Path, outputdir, flags.cacheDir, flags.skipExists, flags.forceClean, flags.progress, flags.selfTest, flags.sync)
	log.Log.Info("Input file: %s", mx3Path)
	log.Log.Info("Output directory: %s", engine.OD())

	log.Log.Init(engine.OD())
	go log.Log.AutoFlushToFile()

	mx3Path = engine.InputFile

	var code *script.BlockStmt
	var err2 error
	if mx3Path != "" {
		// first we compile the entire file into an executable tree
		code, err2 = engine.CompileFile(mx3Path)
		if err2 != nil {
			log.Log.ErrAndExit("Error while parsing `%s`: %v", mx3Path, err2)
		}
		log.Log.PanicIfError(err2)
	}

	// now the parser is not used anymore so it can handle web requests
	if flags.webUIEnabled {
		go api.Start(flags.webUIHost, flags.webUIPort, flags.tunnel, flags.debug)
	}
	// start executing the tree, possibly injecting commands from web gui
	engine.EvalFile(code)

	if flags.interactive {
		engine.RunInteractive()
	}
}

// print version to stdout
func printVersion() {
	log.Log.Info("Version:         %s", engine.VERSION)
	log.Log.Info("Platform:        %s_%s", runtime.GOOS, runtime.GOARCH)
	log.Log.Info("Go Version:      %s (%s)", runtime.Version(), runtime.Compiler)
	log.Log.Info("CUDA Version:    %d.%d (CC=%d PTX)", cu.CUDA_VERSION/1000, (cu.CUDA_VERSION%1000)/10, cuda.UseCC)
	log.Log.Info("GPU Information: %s", cuda.GPUInfo)
}
