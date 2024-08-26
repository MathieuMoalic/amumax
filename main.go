// mumax3 main command
package main

import (
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"

	"github.com/MathieuMoalic/amumax/api"
	"github.com/MathieuMoalic/amumax/cuda"
	"github.com/MathieuMoalic/amumax/engine"
	"github.com/MathieuMoalic/amumax/script"
	"github.com/MathieuMoalic/amumax/timer"
	"github.com/MathieuMoalic/amumax/util"
	"github.com/fatih/color"
	"github.com/minio/selfupdate"
)

var (
	flag_failfast = flag.Bool("failfast", false, "If one simulation fails, stop entire batch immediately")
	flag_test     = flag.Bool("test", false, "Cuda test (internal)")
	flag_version  = flag.Bool("v", true, "Print version")
	flag_vet      = flag.Bool("vet", false, "Check input files for errors, but don't run them")
	flag_update   = flag.Bool("update", false, "Update the amumax binary from the latest github release")
	// more flags in engine/gofiles.go
)

func main() {
	flag.Parse()
	if *flag_update {
		doUpdate()
		return
	}
	// go checkUpdate()

	cuda.Init(*engine.Flag_gpu)

	cuda.Synchronous = *engine.Flag_sync
	timer.Enabled = true //*engine.Flag_sync

	if *flag_version {
		printVersion()
	}

	// used by bootstrap launcher to test cuda
	// successful exit means cuda was initialized fine
	if *flag_test {
		fmt.Println(cuda.GPUInfo)
		os.Exit(0)
	}

	defer engine.CleanExit() // flushes pending output, if any

	if *flag_vet {
		vet()
		return
	}

	switch flag.NArg() {
	case 0:
		if *engine.Flag_interactive {
			runInteractive()
		}
	case 1:
		runFileAndServe(flag.Arg(0))
	default:
		RunQueue(flag.Args())
	}
}

type Release struct {
	TagName string `json:"tag_name"`
}

func doUpdate() {
	resp, err := http.Get("https://github.com/mathieumoalic/amumax/releases/latest/download/amumax")
	if err != nil {
		util.Log.PanicIfError(err)
	}
	defer resp.Body.Close()
	err = selfupdate.Apply(resp.Body, selfupdate.Options{})
	if err != nil {
		color.Red("Error updating")
		color.Red(fmt.Sprint(err))
	}
}

func runInteractive() {
	util.Log.Comment("No input files: starting interactive session")
	// setup outut dir
	now := time.Now()
	outdir := fmt.Sprintf("/tmp/amumax-%v-%02d-%02d_%02dh%02d.zarr", now.Year(), int(now.Month()), now.Day(), now.Hour(), now.Minute())

	engine.InitIO(outdir, outdir)

	// set up some sensible start configuration
	engine.Eval(`
Nx = 128
Ny = 64
Nz = 1
dx = 4e-9
dy = 4e-9
dz = 4e-9
Msat = 1e6
Aex = 10e-12
alpha = 1
m = RandomMag()`)
	if *engine.Flag_webui_enabled {
		go api.Start()
	}
	engine.RunInteractive()
}

func runFileAndServe(fname string) {
	if path.Ext(fname) == ".go" {
		runGoFile(fname)
	} else {
		runMx3File(fname)
	}
}

func runMx3File(mx3Path string) {
	if _, err := os.Stat(mx3Path); errors.Is(err, os.ErrNotExist) {
		util.Log.ErrAndExit("Error: File `%s` does not exist", mx3Path)
	}
	zarrPath := strings.Replace(mx3Path, ".mx3", ".zarr", 1)
	if *engine.Flag_od != "" {
		zarrPath = *engine.Flag_od
	}
	engine.InitIO(mx3Path, zarrPath)
	util.Log.Comment("Input file: %s", mx3Path)
	util.Log.Comment("Output directory: %s", engine.OD())
	util.Log.Init(zarrPath, engine.VERSION == "dev")
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
	if *engine.Flag_webui_enabled {
		go api.Start()
	}
	// start executing the tree, possibly injecting commands from web gui
	engine.EvalFile(code)

	if *engine.Flag_interactive {
		engine.RunInteractive()
	}
}

func runGoFile(fname string) {

	// pass through flags
	flags := []string{"run", fname}
	flag.Visit(func(f *flag.Flag) {
		if f.Name != "o" {
			flags = append(flags, fmt.Sprintf("-%v=%v", f.Name, f.Value))
		}
	})

	if *engine.Flag_od != "" {
		flags = append(flags, fmt.Sprintf("-o=%v", *engine.Flag_od))
	}

	cmd := exec.Command("go", flags...)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		os.Exit(1)
	}
}

// print version to stdout
func printVersion() {
	util.Log.Comment("%v", engine.UNAME)
	util.Log.Comment("GPU info: %s, using cc=%d PTX", cuda.GPUInfo, cuda.UseCC)
}
