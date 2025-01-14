package entrypoint_old

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/MathieuMoalic/amumax/src/api"
	"github.com/MathieuMoalic/amumax/src/engine_old"
	"github.com/MathieuMoalic/amumax/src/engine_old/cuda_old"
	"github.com/MathieuMoalic/amumax/src/engine_old/cuda_old/cu"
	"github.com/MathieuMoalic/amumax/src/engine_old/log_old"
	"github.com/MathieuMoalic/amumax/src/engine_old/queue"
	"github.com/MathieuMoalic/amumax/src/engine_old/script_old"
	"github.com/MathieuMoalic/amumax/src/engine_old/timer_old"
	"github.com/MathieuMoalic/amumax/src/flags"
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

	// go slurm_old.SetEndTimerIfSlurm()
	cuda_old.Init(flags.Gpu)

	cuda_old.Synchronous = flags.Sync
	timer_old.Enabled = flags.Sync

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
	now := time.Now()
	outdir := fmt.Sprintf("/tmp/amumax-%v-%02d-%02d_%02dh%02d.zarr", now.Year(), int(now.Month()), now.Day(), now.Hour(), now.Minute())
	engine_old.InitIO(outdir, outdir, flags.CacheDir, flags.SkipExists, flags.ForceClean, flags.HideProgressBar, flags.SelfTest, flags.Sync)
	log_old.Log.Info("Input file: %s", "none")
	log_old.Log.Info("Output directory: %s", engine_old.OD())
	log_old.Log.Init(engine_old.OD())
	go log_old.Log.AutoFlushToFile()

	if !flags.WebUIDisabled {
		host, port, path, err := url.ParseAddrPath(flags.WebUIAddress)
		log_old.Log.PanicIfError(err)
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
	code, err := engine_old.World.Compile(script)
	log_old.Log.PanicIfError(err)
	engine_old.EvalFile(code)

	engine_old.RunInteractive()
}

func runFileAndServe(mx3Path string, flags *flags.Flags) {
	code, err := setupAndServe(flags, mx3Path, false)
	if err != nil {
		log_old.Log.PanicIfError(err)
	}

	// Execute the compiled input file
	engine_old.EvalFile(code)

	// Enter interactive mode if requested
	if flags.Interactive {
		engine_old.RunInteractive()
	}
}

// print version to stdout
func printVersion() {
	log_old.Log.Info("Version:         %s", version.VERSION)
	log_old.Log.Info("Platform:        %s_%s", runtime.GOOS, runtime.GOARCH)
	log_old.Log.Info("Go Version:      %s (%s)", runtime.Version(), runtime.Compiler)
	log_old.Log.Info("CUDA Version:    %d.%d (CC=%d PTX)", cu.CUDA_VERSION/1000, (cu.CUDA_VERSION%1000)/10, cuda_old.UseCC)
	log_old.Log.Info("GPU Information: %s", cuda_old.GPUInfo_old)
}

func setupAndServe(flags *flags.Flags, mx3Path string, isInteractive bool) (*script_old.BlockStmt, error) {
	var code *script_old.BlockStmt
	var err error

	// Initialize input/output directories
	var inputPath, outputPath string
	if isInteractive {
		now := time.Now()
		outputPath = fmt.Sprintf("/tmp/amumax-%v-%02d-%02d_%02dh%02d.zarr", now.Year(), int(now.Month()), now.Day(), now.Hour(), now.Minute())
		inputPath = outputPath // No input file for interactive mode
		log_old.Log.Info("Input file: %s", "none")
	} else {
		if _, err = os.Stat(mx3Path); errors.Is(err, os.ErrNotExist) {
			log_old.Log.ErrAndExit("Error: File `%s` does not exist", mx3Path)
		}
		inputPath = mx3Path
		outputPath = flags.OutputDir
		if outputPath == "" {
			outputPath = strings.TrimSuffix(mx3Path, ".mx3") + ".zarr"
		}
		log_old.Log.Info("Input file: %s", mx3Path)
	}

	engine_old.InitIO(inputPath, outputPath, flags.CacheDir, flags.SkipExists, flags.ForceClean, flags.HideProgressBar, flags.SelfTest, flags.Sync)
	log_old.Log.Info("Output directory: %s", engine_old.OD())
	log_old.Log.Init(engine_old.OD())
	go log_old.Log.AutoFlushToFile()

	// Web UI setup
	if !flags.WebUIDisabled {
		host, port, path, err1 := url.ParseAddrPath(flags.WebUIAddress)
		log_old.Log.PanicIfError(err1)
		go api.Start(host, port, path, flags.Tunnel, flags.Debug)
	}

	// Compile input file if in non-interactive mode
	if !isInteractive && mx3Path != "" {
		code, err = engine_old.CompileFile(mx3Path)
		if err != nil {
			log_old.Log.ErrAndExit("Error while parsing `%s`: %v", mx3Path, err)
		}
	}

	return code, err
}
