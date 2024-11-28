package new_engine

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/MathieuMoalic/amumax/src/flags"
	"github.com/MathieuMoalic/amumax/src/log"
	"github.com/MathieuMoalic/amumax/src/mesh"
	"github.com/MathieuMoalic/amumax/src/metadata"
	"github.com/MathieuMoalic/amumax/src/new_fsutil"
	"github.com/MathieuMoalic/amumax/src/new_log"
	"github.com/MathieuMoalic/amumax/src/timer"
	"github.com/MathieuMoalic/amumax/src/zarr"
	"github.com/fatih/color"
)

type engineState struct {
	zarrPath        string
	script          string
	scriptPath      string
	flags           *flags.FlagsType
	metadata        *metadata.Metadata
	world           *world
	log             *new_log.Logs
	table           *table
	solver          *solver
	mesh            *mesh.Mesh
	magnetization   *magnetization
	geometry        *geometry
	regions         *regions
	savedQuantities *savedQuantities
	utils           *utils
	windowShift     *windowShift
	shape           *shapeList
	grains          *grains
	fs              *new_fsutil.FileSystem
	config          *configList
}

func newEngineState(givenFlags *flags.FlagsType) *engineState {
	return &engineState{flags: givenFlags}
}

func (s *engineState) start(mx3path string) {
	scriptBytes, err := os.ReadFile(mx3path)
	if err != nil {
		color.Red("Error reading script: %v", err)
		os.Exit(1)
	}
	s.script = string(scriptBytes)
	s.scriptPath = mx3path
	s.run()
}

func (s *engineState) startInteractive() {
	log.Log.Info("No input files: starting interactive session")
	s.script = `
	Nx = 128
	Ny = 64
	Nz = 1
	dx = 3e-9
	dy = 3e-9
	dz = 3e-9
	Msat = 1e6
	Aex = 10e-12
	alpha = 1
	m = RandomMag()`
	s.run()
	// setup outut dir
	now := time.Now()
	s.zarrPath = fmt.Sprintf("/tmp/amumax-%v-%02d-%02d_%02dh%02d.zarr", now.Year(), int(now.Month()), now.Day(), now.Hour(), now.Minute())
}

func (s *engineState) run() {
	defer s.cleanExit()
	s.initIO()
	s.log = new_log.NewLogs(s.zarrPath, s.fs, s.flags.Debug)
	s.metadata = metadata.NewMetadata(s.fs, s.log)
	s.world = newWorld(s)
	s.windowShift = newWindowShift(s)
	s.shape = newShape(s)
	s.table = newTable(s)
	s.mesh = &mesh.Mesh{}
	s.solver = newSolver(s)
	s.magnetization = newMagnetization(s)
	s.regions = newRegions(s)
	s.geometry = newGeom(s)
	s.savedQuantities = newSavedQuantities(s)
	s.utils = newUtils(s)
	s.grains = newGrains(s)
	s.config = newConfigList(s.mesh, s.world)
	s.world.register()
	scriptParser := newScriptParser(s)
	err := scriptParser.Parse(s.script)
	if err != nil {
		s.log.ErrAndExit("Error parsing script: %v", err)
	}

	err = scriptParser.execute()
	if err != nil {
		s.log.ErrAndExit("Error executing script: %v", err)
	}
}

func (s *engineState) makeZarrPath() {
	if s.flags.OutputDir != "" {
		s.zarrPath = s.flags.OutputDir
	} else {
		if s.scriptPath == "" {
			now := time.Now()
			s.zarrPath = fmt.Sprintf("/tmp/amumax-%v-%02d-%02d_%02dh%02d.zarr", now.Year(), int(now.Month()), now.Day(), now.Hour(), now.Minute())
		} else {
			s.zarrPath = strings.TrimSuffix(s.scriptPath, ".mx3") + ".zarr"
		}
	}
	if !strings.HasSuffix(s.zarrPath, "/") {
		s.zarrPath += "/"
	}
}

func (s *engineState) initIO() {
	s.makeZarrPath()
	s.fs = new_fsutil.NewFileSystem(s.zarrPath)
	if s.fs.IsDir("") {
		// if directory exists and --skip-exist flag is set, skip the directory
		if s.flags.SkipExists {
			log.Log.Warn("Directory `%s` exists, skipping `%s` because of --skip-exist flag.", s.zarrPath, s.scriptPath)
			os.Exit(0)
			// if directory exists and --force-clean flag is set, remove the directory
		} else if s.flags.ForceClean {
			log.Log.Warn("Cleaning `%s`", s.zarrPath)
			log.Log.PanicIfError(s.fs.Remove(""))
			log.Log.PanicIfError(s.fs.Mkdir(""))
		}
	} else {
		log.Log.PanicIfError(s.fs.Mkdir(""))
	}
	zarr.InitZgroup("", s.zarrPath)
}

func (s *engineState) cleanExit() {
	s.fs.Drain()    // wait for the save queue to finish
	s.table.flush() // flush table to disk
	if s.flags.Sync {
		timer.Print(os.Stdout)
	}
	// s.metadata.Add("steps", NSteps)
	s.metadata.End()
	s.log.Info("**************** Simulation Ended ****************** //")
	s.log.FlushToFile()
}
