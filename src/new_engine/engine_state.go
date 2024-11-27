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
	ZarrPath        string
	Script          string
	ScriptPath      string
	Flags           *flags.FlagsType
	Metadata        *zarr.Metadata
	World           *World
	Log             *log.Logs
	Table           *TableStruct
	Solver          *Solver
	Mesh            *mesh.Mesh
	Magnetization   *Magnetization
	Geometry        *Geometry
	Regions         *Regions
	Quantities      *Quantities
	SavedQuantities *savedQuantitiesType
	Utils           *Utils
	Saver           *Saver
	WindowShift     *WindowShift
	Shape           *ShapeStruct
	Grains          *Grains
	fs              *new_fsutil.FileSystem
}

func NewEngineState(givenFlags *flags.FlagsType) *EngineStateStruct {
	return &EngineStateStruct{Flags: givenFlags, Metadata: &zarr.Metadata{}}
}

func (s *EngineStateStruct) Start(mx3path string) {
	scriptBytes, err := os.ReadFile(mx3path)
	if err != nil {
		color.Red("Error reading script: %v", err)
		os.Exit(1)
	}
	s.Script = string(scriptBytes)
	s.ScriptPath = mx3path
	s.run()
}

func (s *EngineStateStruct) StartInteractive() {
	log.Log.Info("No input files: starting interactive session")
	s.Script = `
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
	s.ZarrPath = fmt.Sprintf("/tmp/amumax-%v-%02d-%02d_%02dh%02d.zarr", now.Year(), int(now.Month()), now.Day(), now.Hour(), now.Minute())
}

func (s *EngineStateStruct) run() {
	defer s.CleanExit()
	s.initIO()
	s.initLog()
	s.initMetadata()
	s.World = NewWorld(s)
	s.Saver = NewSaver(100)
	s.WindowShift = NewWindowShift(s)
	s.Shape = NewShape(s)
	s.initTable()
	s.Mesh = &mesh.Mesh{}
	s.Solver = NewSolver(s)
	s.Magnetization = NewMagnetization(s)
	s.Regions = NewRegions(s)
	s.Geometry = NewGeom(s)
	s.SavedQuantities = NewSavedQuantities(s)
	s.Utils = NewUtils(s)
	s.Grains = NewGrains(s)
	s.World.register()
	scriptParser := NewScriptParser(s)
	err := scriptParser.Parse(s.Script)
	if err != nil {
		s.Log.ErrAndExit("Error parsing script: %v", err)
	}

	err = scriptParser.Execute()
	if err != nil {
		s.Log.ErrAndExit("Error executing script: %v", err)
	}
}

func (s *EngineStateStruct) makeZarrPath() {
	if s.Flags.OutputDir != "" {
		s.ZarrPath = s.Flags.OutputDir
	} else {
		if s.ScriptPath == "" {
			now := time.Now()
			s.ZarrPath = fmt.Sprintf("/tmp/amumax-%v-%02d-%02d_%02dh%02d.zarr", now.Year(), int(now.Month()), now.Day(), now.Hour(), now.Minute())
		} else {
			s.ZarrPath = strings.TrimSuffix(s.ScriptPath, ".mx3") + ".zarr"
		}
	}
	if !strings.HasSuffix(s.ZarrPath, "/") {
		s.ZarrPath += "/"
	}
}

func (s *EngineStateStruct) initIO() {
	s.makeZarrPath()
	s.fs = new_fsutil.NewFileSystem(s.ZarrPath)
	if s.fs.IsDir("") {
		// if directory exists and --skip-exist flag is set, skip the directory
		if s.Flags.SkipExists {
			log.Log.Warn("Directory `%s` exists, skipping `%s` because of --skip-exist flag.", s.ZarrPath, s.ScriptPath)
			os.Exit(0)
			// if directory exists and --force-clean flag is set, remove the directory
		} else if s.Flags.ForceClean {
			log.Log.Warn("Cleaning `%s`", s.ZarrPath)
			log.Log.PanicIfError(s.fs.Remove(""))
			log.Log.PanicIfError(s.fs.Mkdir(""))
		}
	} else {
		log.Log.PanicIfError(s.fs.Mkdir(""))
	}
	zarr.InitZgroup("", s.ZarrPath)
}

func (s *EngineStateStruct) initLog() {
	s.Log = &log.Logs{}
	s.Log.Info("Input file: %s", s.ScriptPath)
	s.Log.Info("Output directory: %s", s.ZarrPath)
	s.Log.Init(s.ZarrPath)
	s.Log.SetDebug(s.Flags.Debug)
	go s.Log.AutoFlushToFile()
}

func (s *EngineStateStruct) initTable() {
	s.Table = &TableStruct{
		EngineState:    s,
		Data:           make(map[string][]float64),
		Step:           -1,
		AutoSavePeriod: 0.0,
		FlushInterval:  5 * time.Second,
	}
	err := s.fs.Remove("table")
	log.Log.PanicIfError(err)
	zarr.InitZgroup("table", s.ZarrPath)
	s.Table.AddColumn("step", "")
	s.Table.AddColumn("t", "s")
	// s.Table.tableAdd(s.Magnetization)
	go s.Table.tablesAutoFlush()
}

func (s *EngineStateStruct) initMetadata() {
	s.Metadata = &zarr.Metadata{}
	s.Metadata.Init(s.ZarrPath, time.Now(), cuda.GPUInfo)
}

func (s *EngineStateStruct) CleanExit() {
	s.Saver.Drain()
	s.Table.Flush()
	if s.Flags.Sync {
		timer.Print(os.Stdout)
	}
	// s.Metadata.Add("steps", NSteps)
	s.Metadata.End()
	s.Log.Info("**************** Simulation Ended ****************** //")
	s.Log.FlushToFile()
}
