// mumax3 main command
package entrypoint

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/MathieuMoalic/amumax/api"
	"github.com/MathieuMoalic/amumax/cuda"
	"github.com/MathieuMoalic/amumax/cuda/cu"
	"github.com/MathieuMoalic/amumax/engine"
	"github.com/MathieuMoalic/amumax/script"
	"github.com/MathieuMoalic/amumax/util"
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
	util.Log.Info("No input files: starting interactive session")
	// setup outut dir
	now := time.Now()
	outdir := fmt.Sprintf("/tmp/amumax-%v-%02d-%02d_%02dh%02d.zarr", now.Year(), int(now.Month()), now.Day(), now.Hour(), now.Minute())

	engine.InitIO(outdir, outdir, flags.cacheDir, flags.skipExists, flags.forceClean, flags.progress, flags.selfTest, flags.sync)
	util.Log.Info("Input file: %s", "none")
	util.Log.Info("Output directory: %s", engine.OD())
	util.Log.Init(engine.OD())

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
		go api.Start(flags.webUIHost, flags.webUIPort, flags.tunnel)
	}
	engine.RunInteractive()
}

func runFileAndServe(mx3Path string) {
	if _, err := os.Stat(mx3Path); errors.Is(err, os.ErrNotExist) {
		util.Log.ErrAndExit("Error: File `%s` does not exist", mx3Path)
	}
	outputdir := strings.TrimSuffix(mx3Path, ".mx3") + ".zarr"

	if flags.outputDir != "" {
		outputdir = flags.outputDir
	}
	engine.InitIO(mx3Path, outputdir, flags.cacheDir, flags.skipExists, flags.forceClean, flags.progress, flags.selfTest, flags.sync)
	util.Log.Info("Input file: %s", mx3Path)
	util.Log.Info("Output directory: %s", engine.OD())

	util.Log.Init(engine.OD())
	go util.Log.AutoFlushToFile()

	mx3Path = engine.InputFile

	var code *script.BlockStmt
	var err2 error
	if mx3Path != "" {
		// first we compile the entire file into an executable tree
		code, err2 = engine.CompileFile(mx3Path)
		if err2 != nil {
			util.Log.Err("Error while parsing `%s`", mx3Path)
		}
		util.Log.PanicIfError(err2)
	}

	// now the parser is not used anymore so it can handle web requests
	if flags.webUIEnabled {
		go api.Start(flags.webUIHost, flags.webUIPort, flags.tunnel)
	}
	// start executing the tree, possibly injecting commands from web gui
	engine.EvalFile(code)

	if flags.interactive {
		engine.RunInteractive()
	}
}

// print version to stdout
func printVersion() {
	util.Log.Info("Version:         %s", engine.VERSION)
	util.Log.Info("Platform:        %s_%s", runtime.GOOS, runtime.GOARCH)
	util.Log.Info("Go Version:      %s (%s)", runtime.Version(), runtime.Compiler)
	util.Log.Info("CUDA Version:    %d.%d (CC=%d PTX)", cu.CUDA_VERSION/1000, (cu.CUDA_VERSION%1000)/10, cuda.UseCC)
	util.Log.Info("GPU Information: %s", cuda.GPUInfo)
}
