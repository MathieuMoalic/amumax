// Package entrypoint handles the main entry point of the application.
package entrypoint

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/MathieuMoalic/amumax/src/api"
	"github.com/MathieuMoalic/amumax/src/cuda"
	"github.com/MathieuMoalic/amumax/src/cuda/cu"
	"github.com/MathieuMoalic/amumax/src/engine"
	"github.com/MathieuMoalic/amumax/src/flags"
	"github.com/MathieuMoalic/amumax/src/log"
	"github.com/MathieuMoalic/amumax/src/queue"
	"github.com/MathieuMoalic/amumax/src/script"
	"github.com/MathieuMoalic/amumax/src/timer"
	"github.com/MathieuMoalic/amumax/src/update"
	"github.com/MathieuMoalic/amumax/src/url"
	"github.com/MathieuMoalic/amumax/src/version"
)

func Entrypoint(cmd *cobra.Command, args []string, flags *flags.Flags) {
	log.Log.SetDebug(flags.Debug)

	if flags.Update {
		update.ShowUpdateMenu()
		return
	}

	// go slurm.SetEndTimerIfSlurm()
	cuda.Init(flags.Gpu)

	cuda.Synchronous = flags.Sync
	timer.Enabled = flags.Sync

	printVersion()
	if flags.Version {
		return
	}

	engine.Insecure = flags.Insecure

	defer engine.CleanExit() // flushes pending output, if any

	if flags.Vet {
		engine.Vet()
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
	log.Log.Info("No input files: starting interactive session")
	now := time.Now()
	outdir := fmt.Sprintf("/tmp/amumax-%v-%02d-%02d_%02dh%02d.zarr", now.Year(), int(now.Month()), now.Day(), now.Hour(), now.Minute())
	engine.InitIO(outdir, outdir, flags.CacheDir, flags.SkipExists, flags.ForceClean, flags.HideProgressBar, flags.SelfTest, flags.Sync)
	log.Log.Info("Input file: %s", "none")
	log.Log.Info("Output directory: %s", engine.OD())
	log.Log.Init(engine.OD())
	go log.Log.AutoFlushToFile()

	if !flags.WebUIDisabled {
		host, port, path, err := url.ParseAddrPath(flags.WebUIAddress)
		log.Log.PanicIfError(err)
		go api.Start(host, port, path, flags.Tunnel, flags.Debug)
	}

	// Compile the hardcoded script and evaluate it
	script := `
        Nx = 128
        Ny = 64
        Nz = 1
        dx = 3e-9
        dy = 3e-9
        dz = 3e-9
        Msat = 1e6
        Aex = 10e-12
        m = RandomMag()
    `
	code, err := engine.World.Compile(script)
	log.Log.PanicIfError(err)
	engine.EvalFile(code)

	engine.RunInteractive()
}

func runFileAndServe(mx3Path string, flags *flags.Flags) {
	code, err := setupAndServe(flags, mx3Path, false)
	if err != nil {
		log.Log.PanicIfError(err)
	}

	// Execute the compiled input file
	engine.EvalFile(code)

	// Enter interactive mode if requested
	if flags.Interactive {
		engine.RunInteractive()
	}
}

// print version to stdout
func printVersion() {
	log.Log.Info("Version:         %s", version.VERSION)
	log.Log.Info("Platform:        %s_%s", runtime.GOOS, runtime.GOARCH)
	log.Log.Info("Go Version:      %s (%s)", runtime.Version(), runtime.Compiler)
	log.Log.Info("CUDA Version:    %d.%d (CC=%d PTX)", cu.CUDA_VERSION/1000, (cu.CUDA_VERSION%1000)/10, cuda.UseCC)
	log.Log.Info("GPU Information: %s", cuda.GPUInfoOld)
}

func setupAndServe(flags *flags.Flags, mx3Path string, isInteractive bool) (*script.BlockStmt, error) {
	var code *script.BlockStmt
	var err error

	// Initialize input/output directories
	var inputPath, outputPath string
	if isInteractive {
		now := time.Now()
		outputPath = fmt.Sprintf("/tmp/amumax-%v-%02d-%02d_%02dh%02d.zarr", now.Year(), int(now.Month()), now.Day(), now.Hour(), now.Minute())
		inputPath = outputPath // No input file for interactive mode
		log.Log.Info("Input file: %s", "none")
	} else {
		if _, err = os.Stat(mx3Path); errors.Is(err, os.ErrNotExist) {
			log.Log.ErrAndExit("Error: File `%s` does not exist", mx3Path)
		}
		inputPath = mx3Path
		outputPath = flags.OutputDir
		if outputPath == "" {
			outputPath = strings.TrimSuffix(mx3Path, ".mx3") + ".zarr"
		}
		log.Log.Info("Input file: %s", mx3Path)
	}

	engine.InitIO(inputPath, outputPath, flags.CacheDir, flags.SkipExists, flags.ForceClean, flags.HideProgressBar, flags.SelfTest, flags.Sync)
	log.Log.Info("Output directory: %s", engine.OD())
	log.Log.Init(engine.OD())
	go log.Log.AutoFlushToFile()

	// Web UI setup
	if !flags.WebUIDisabled {
		host, port, path, err1 := url.ParseAddrPath(flags.WebUIAddress)
		log.Log.PanicIfError(err1)
		go api.Start(host, port, path, flags.Tunnel, flags.Debug)
	}

	// Compile input file if in non-interactive mode
	if !isInteractive && mx3Path != "" {
		code, err = engine.CompileFile(mx3Path)
		if err != nil {
			log.Log.ErrAndExit("Error while parsing `%s`: %v", mx3Path, err)
		}
	}

	return code, err
}
