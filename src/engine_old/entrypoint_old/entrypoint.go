package entrypoint_old

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
	"github.com/MathieuMoalic/amumax/src/engine_old"
	"github.com/MathieuMoalic/amumax/src/engine_old/log_old"
	"github.com/MathieuMoalic/amumax/src/engine_old/queue"
	"github.com/MathieuMoalic/amumax/src/engine_old/script_old"
	"github.com/MathieuMoalic/amumax/src/flags"
	"github.com/MathieuMoalic/amumax/src/slurm"
	"github.com/MathieuMoalic/amumax/src/timer"
	"github.com/MathieuMoalic/amumax/src/update"
	"github.com/MathieuMoalic/amumax/src/url"
	"github.com/MathieuMoalic/amumax/src/version"
	"github.com/spf13/cobra"
)

func Entrypoint(cmd *cobra.Command, args []string, flags *flags.Flags) {

	log_old.Log.SetDebug(flags.Debug)

	if flags.Update {
		update.ShowUpdateMenu()
		return
	}

	go slurm.SetEndTimerIfSlurm()
	cuda.Init(flags.Gpu)

	cuda.Synchronous = flags.Sync
	timer.Enabled = flags.Sync

	printVersion()
	if flags.Version {
		return
	}

	engine_old.Insecure = flags.Insecure

	defer engine_old.CleanExit() // flushes pending output, if any

	if flags.Vet {
		engine_old.Vet()
		return
	}

	if len(args) == 0 && flags.Interactive {
		runInteractive(flags)
	} else if len(args) == 1 {
		runFileAndServe(args[0], flags)
	} else if len(args) > 1 {
		queue.RunQueue(args, flags)
	} else {
		_ = cmd.Help()
	}
}

type Release struct {
	TagName string `json:"tag_name"`
}

func runInteractive(flags *flags.Flags) {
	log_old.Log.Info("No input files: starting interactive session")
	// setup outut dir
	now := time.Now()
	outdir := fmt.Sprintf("/tmp/amumax-%v-%02d-%02d_%02dh%02d.zarr", now.Year(), int(now.Month()), now.Day(), now.Hour(), now.Minute())

	engine_old.InitIO(outdir, outdir, flags.CacheDir, flags.SkipExists, flags.ForceClean, flags.HideProgressBar, flags.SelfTest, flags.Sync)
	log_old.Log.Info("Input file: %s", "none")
	log_old.Log.Info("Output directory: %s", engine_old.OD())
	log_old.Log.Init(engine_old.OD())

	// set up some sensible start configuration
	engine_old.Eval(`
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
	if !flags.WebUIDisabled {
		host, port, path, err := url.ParseAddrPath(flags.WebUIAddress)
		log_old.Log.PanicIfError(err)
		go api.Start(host, port, path, flags.Tunnel, flags.Debug)
	}
	engine_old.RunInteractive()
}

func runFileAndServe(mx3Path string, flags *flags.Flags) {
	if _, err := os.Stat(mx3Path); errors.Is(err, os.ErrNotExist) {
		log_old.Log.ErrAndExit("Error: File `%s` does not exist", mx3Path)
	}
	outputdir := strings.TrimSuffix(mx3Path, ".mx3") + ".zarr"

	if flags.OutputDir != "" {
		outputdir = flags.OutputDir
	}
	engine_old.InitIO(mx3Path, outputdir, flags.CacheDir, flags.SkipExists, flags.ForceClean, flags.HideProgressBar, flags.SelfTest, flags.Sync)
	log_old.Log.Info("Input file: %s", mx3Path)
	log_old.Log.Info("Output directory: %s", engine_old.OD())

	log_old.Log.Init(engine_old.OD())
	go log_old.Log.AutoFlushToFile()

	mx3Path = engine_old.InputFile

	var code *script_old.BlockStmt
	var err2 error
	if mx3Path != "" {
		// first we compile the entire file into an executable tree
		code, err2 = engine_old.CompileFile(mx3Path)
		if err2 != nil {
			log_old.Log.ErrAndExit("Error while parsing `%s`: %v", mx3Path, err2)
		}
		log_old.Log.PanicIfError(err2)
	}
	// now the parser is not used anymore so it can handle web requests
	if !flags.WebUIDisabled {
		host, port, path, err := url.ParseAddrPath(flags.WebUIAddress)
		log_old.Log.PanicIfError(err)
		go api.Start(host, port, path, flags.Tunnel, flags.Debug)
	}
	// start executing the tree, possibly injecting commands from web gui
	engine_old.EvalFile(code)

	if flags.Interactive {
		engine_old.RunInteractive()
	}
}

// print version to stdout
func printVersion() {
	log_old.Log.Info("Version:         %s", version.VERSION)
	log_old.Log.Info("Platform:        %s_%s", runtime.GOOS, runtime.GOARCH)
	log_old.Log.Info("Go Version:      %s (%s)", runtime.Version(), runtime.Compiler)
	log_old.Log.Info("CUDA Version:    %d.%d (CC=%d PTX)", cu.CUDA_VERSION/1000, (cu.CUDA_VERSION%1000)/10, cuda.UseCC)
	log_old.Log.Info("GPU Information: %s", cuda.GPUInfo_old)
}
