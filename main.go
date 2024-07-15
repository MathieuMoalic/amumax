// mumax3 main command
package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strconv"
	"time"

	"github.com/MathieuMoalic/amumax/api"
	"github.com/MathieuMoalic/amumax/cuda"
	"github.com/MathieuMoalic/amumax/engine"
	"github.com/MathieuMoalic/amumax/script"
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
	log.SetPrefix("")
	log.SetFlags(0)

	cuda.Init(*engine.Flag_gpu)

	cuda.Synchronous = *engine.Flag_sync
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

func doUpdate() error {
	resp, err := http.Get("https://github.com/mathieumoalic/amumax/releases/latest/download/amumax")
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	err = selfupdate.Apply(resp.Body, selfupdate.Options{})
	if err != nil {
		color.Red("Error updating")
		color.Red(fmt.Sprint(err))
	}
	return err
}

func checkUpdate() {
	if engine.VERSION == "dev" {
		return
	}
	resp, err := http.Get("https://api.github.com/repos/mathieumoalic/amumax/releases/latest")
	if err != nil {
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	var release Release
	err = json.Unmarshal(body, &release)
	if err != nil {
		return
	}
	if release.TagName != engine.VERSION {
		color.HiCyan("\rCurrent amumax version          : %s", engine.VERSION)
		color.HiCyan("Latest amumax version on github : %s", release.TagName)
		color.HiCyan("Run the following command to update amumax:")
		color.Blue("amumax --update")
	}
}

func runInteractive() {
	util.Log("No input files: starting interactive session")
	// setup outut dir
	now := time.Now()
	outdir := fmt.Sprintf("/tmp/amumax-%v-%02d-%02d_%02dh%02d.zarr", now.Year(), int(now.Month()), now.Day(), now.Hour(), now.Minute())

	engine.InitIO(outdir, outdir)

	engine.Timeout = 365 * 24 * time.Hour // basically forever

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
	goServeGUI()
	engine.RunInteractive()
}

func runFileAndServe(fname string) {
	if path.Ext(fname) == ".go" {
		runGoFile(fname)
	} else {
		runScript(fname)
	}
}

func runScript(fname string) {
	if _, err := os.Stat(fname); errors.Is(err, os.ErrNotExist) {
		util.Fatal("Error: File `", fname, "` does not exist")
	}
	outDir := util.NoExt(fname) + ".zarr"
	if *engine.Flag_od != "" {
		outDir = *engine.Flag_od
	}
	engine.InitIO(fname, outDir)

	fname = engine.InputFile

	var code *script.BlockStmt
	var err2 error
	if fname != "" {
		// first we compile the entire file into an executable tree
		code, err2 = engine.CompileFile(fname)
		if err2 != nil {
			engine.LogErr("Error while parsing `", fname, "`")
		}
		util.FatalErr(err2)
	}

	// now the parser is not used anymore so it can handle web requests
	goServeGUI()

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
	log.Println("go", flags)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		os.Exit(1)
	}
}

func findAvailablePort(startAddress string) (string, error) {
	// Split the address to extract the host and port
	host, portStr, err := net.SplitHostPort(startAddress)
	if err != nil {
		return "", fmt.Errorf("invalid start address: %v", err)
	}

	// Convert the port to an integer
	startPort, err := strconv.Atoi(portStr)
	if err != nil {
		return "", fmt.Errorf("invalid port number: %v", err)
	}

	// Loop to find the first available port
	for port := startPort; port <= 65535; port++ {
		address := net.JoinHostPort(host, strconv.Itoa(port))
		listener, err := net.Listen("tcp", address)
		if err == nil {
			// Close the listener immediately, we just wanted to check availability
			listener.Close()
			return address, nil
		}
	}
	return "", fmt.Errorf("no available ports found")
}

// start Gui server and return server address
func goServeGUI() {
	if *engine.Flag_webui_addr == "" {
		util.LogWarn(`WebUI is disabled (-http="")`)
		return
	}
	addr, err := findAvailablePort(*engine.Flag_webui_addr)
	if err != nil {
		log.Fatalf("Failed to find available port: %v", err)
	}
	util.Log(fmt.Sprintf("Serving GUI at http://%s", addr))
	go api.Start(addr)
}

// print version to stdout
func printVersion() {
	engine.LogOut(engine.UNAME)
	engine.LogOut(fmt.Sprintf("GPU info: %s, using cc=%d PTX", cuda.GPUInfo, cuda.UseCC))
}
