// mumax3 main command
package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/MathieuMoalic/amumax/api"
	"github.com/MathieuMoalic/amumax/cuda"
	"github.com/MathieuMoalic/amumax/engine"
	"github.com/MathieuMoalic/amumax/script"
	"github.com/MathieuMoalic/amumax/util"
	"github.com/fatih/color"
	"github.com/minio/selfupdate"
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
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

	engine.InitIO(outdir, outdir, flags.cacheDir, flags.skipExists, flags.forceClean, flags.progress, flags.selfTest, flags.sync)

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
		go api.Start(flags.webUIHost, flags.webUIPort)
	}
	engine.RunInteractive()
}

func runFileAndServe(mx3Path string) {
	if _, err := os.Stat(mx3Path); errors.Is(err, os.ErrNotExist) {
		util.Log.ErrAndExit("Error: File `%s` does not exist", mx3Path)
	}

	outputdir := strings.Replace(mx3Path, ".mx3", ".zarr", 1)
	if flags.outputDir != "" {
		outputdir = flags.outputDir
	}
	engine.InitIO(mx3Path, outputdir, flags.cacheDir, flags.skipExists, flags.forceClean, flags.progress, flags.selfTest, flags.sync)
	util.Log.Comment("Input file: %s", mx3Path)
	util.Log.Comment("Output directory: %s", engine.OD())
	util.Log.Init(engine.OD(), engine.VERSION == "dev")
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
		go api.Start(flags.webUIHost, flags.webUIPort)
	}
	// start executing the tree, possibly injecting commands from web gui
	engine.EvalFile(code)

	if flags.interactive {
		engine.RunInteractive()
	}
}

// print version to stdout
func printVersion() {
	util.Log.Comment("%v", engine.UNAME)
	util.Log.Comment("GPU info: %s, using cc=%d PTX", cuda.GPUInfo, cuda.UseCC)
}
