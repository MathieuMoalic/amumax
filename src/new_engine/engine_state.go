package new_engine

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/MathieuMoalic/amumax/src/cuda"
	"github.com/MathieuMoalic/amumax/src/flags"
	"github.com/MathieuMoalic/amumax/src/log"
	"github.com/MathieuMoalic/amumax/src/mesh"
	"github.com/MathieuMoalic/amumax/src/new_fsutil"
	"github.com/MathieuMoalic/amumax/src/timer"
	"github.com/MathieuMoalic/amumax/src/zarr"
	"github.com/fatih/color"
)

type EngineStateStruct struct {
	zarrPath        string
	script          string
	scriptPath      string
	flags           *flags.FlagsType
	metadata        *zarr.Metadata
	world           *World
	log             *log.Logs
	table           *Table
	solver          *Solver
	mesh            *mesh.Mesh
	magnetization   *Magnetization
	geometry        *Geometry
	regions         *Regions
	savedQuantities *savedQuantities
	utils           *Utils
	windowShift     *WindowShift
	shape           *Shape
	grains          *Grains
	fs              *new_fsutil.FileSystem
}

func NewEngineState(givenFlags *flags.FlagsType) *EngineStateStruct {
	return &EngineStateStruct{flags: givenFlags, metadata: &zarr.Metadata{}}
}

func (s *EngineStateStruct) Start(mx3path string) {
	scriptBytes, err := os.ReadFile(mx3path)
	if err != nil {
		color.Red("Error reading script: %v", err)
		os.Exit(1)
	}
	s.script = string(scriptBytes)
	s.scriptPath = mx3path
	s.run()
}

func (s *EngineStateStruct) StartInteractive() {
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

func (s *EngineStateStruct) run() {
	defer s.CleanExit()
	s.initIO()
	s.initLog()
	s.initMetadata()
	s.world = NewWorld(s)
	s.windowShift = NewWindowShift(s)
	s.shape = NewShape(s)
	s.initTable()
	s.mesh = &mesh.Mesh{}
	s.solver = NewSolver(s)
	s.magnetization = NewMagnetization(s)
	s.regions = NewRegions(s)
	s.geometry = NewGeom(s)
	s.savedQuantities = NewSavedQuantities(s)
	s.utils = NewUtils(s)
	s.grains = NewGrains(s)
	s.world.register()
	scriptParser := NewScriptParser(s)
	err := scriptParser.Parse(s.script)
	if err != nil {
		s.log.ErrAndExit("Error parsing script: %v", err)
	}

	err = scriptParser.Execute()
	if err != nil {
		s.log.ErrAndExit("Error executing script: %v", err)
	}
}

func (s *EngineStateStruct) makeZarrPath() {
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

func (s *EngineStateStruct) initIO() {
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

func (s *EngineStateStruct) initLog() {
	s.log = &log.Logs{}
	s.log.Info("Input file: %s", s.scriptPath)
	s.log.Info("Output directory: %s", s.zarrPath)
	s.log.Init(s.zarrPath)
	s.log.SetDebug(s.flags.Debug)
	go s.log.AutoFlushToFile()
}

func (s *EngineStateStruct) initTable() {
	s.table = &Table{
		EngineState:    s,
		Data:           make(map[string][]float64),
		Step:           -1,
		AutoSavePeriod: 0.0,
		FlushInterval:  5 * time.Second,
	}
	err := s.fs.Remove("table")
	log.Log.PanicIfError(err)
	zarr.InitZgroup("table", s.zarrPath)
	s.table.AddColumn("step", "")
	s.table.AddColumn("t", "s")
	// s.Table.tableAdd(s.Magnetization)
	go s.table.tablesAutoFlush()
}

func (s *EngineStateStruct) initMetadata() {
	s.metadata = &zarr.Metadata{}
	s.metadata.Init(s.zarrPath, time.Now(), cuda.GPUInfo)
}

func (s *EngineStateStruct) CleanExit() {
	s.fs.Drain()
	s.table.Flush()
	if s.flags.Sync {
		timer.Print(os.Stdout)
	}
	// s.Metadata.Add("steps", NSteps)
	s.metadata.End()
	s.log.Info("**************** Simulation Ended ****************** //")
	s.log.FlushToFile()
}
