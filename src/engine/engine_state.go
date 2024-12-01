package engine

import (
	"os"
	"strings"
	"time"

	"github.com/MathieuMoalic/amumax/src/flags"
	"github.com/MathieuMoalic/amumax/src/fsutil"
	"github.com/MathieuMoalic/amumax/src/geometry"
	"github.com/MathieuMoalic/amumax/src/grains"
	"github.com/MathieuMoalic/amumax/src/log"
	"github.com/MathieuMoalic/amumax/src/mag_config"
	"github.com/MathieuMoalic/amumax/src/magnetization"
	"github.com/MathieuMoalic/amumax/src/mesh"
	"github.com/MathieuMoalic/amumax/src/metadata"
	"github.com/MathieuMoalic/amumax/src/regions"
	"github.com/MathieuMoalic/amumax/src/saved_quantities"
	"github.com/MathieuMoalic/amumax/src/script"
	"github.com/MathieuMoalic/amumax/src/shape"
	"github.com/MathieuMoalic/amumax/src/solver"
	"github.com/MathieuMoalic/amumax/src/table"
	"github.com/MathieuMoalic/amumax/src/timer"
	"github.com/MathieuMoalic/amumax/src/window_shift"
)

type engineState struct {
	flags           *flags.Flags
	fs              *fsutil.FileSystem
	log             *log.Logs
	metadata        *metadata.Metadata
	table           *table.Table
	solver          *solver.Solver
	mesh            *mesh.Mesh
	magnetization   *magnetization.Magnetization
	geometry        *geometry.Geometry
	regions         *regions.Regions
	savedQuantities *saved_quantities.SavedQuantities
	windowShift     *window_shift.WindowShift
	shape           *shape.ShapeList
	grains          *grains.Grains
	config          *mag_config.ConfigList
	script          *script.ScriptParser

	autoFlushInterval time.Duration
}

func newEngineState(givenFlags *flags.Flags, log *log.Logs) *engineState {
	return &engineState{flags: givenFlags, log: log, autoFlushInterval: 5 * time.Second}
}
func (s *engineState) init(scriptStr string) {
	// initialize empty structs first so we can pass the pointers
	// to the actual init functions

	s.metadata = &metadata.Metadata{}
	s.table = &table.Table{}
	s.solver = &solver.Solver{}
	s.mesh = &mesh.Mesh{}
	s.magnetization = &magnetization.Magnetization{}
	s.geometry = &geometry.Geometry{}
	s.regions = &regions.Regions{}
	s.savedQuantities = &saved_quantities.SavedQuantities{}
	s.windowShift = &window_shift.WindowShift{}
	s.shape = &shape.ShapeList{}
	s.grains = &grains.Grains{}
	s.config = &mag_config.ConfigList{}
	s.script = &script.ScriptParser{}

	s.metadata.Init(s.fs, s.log)
	s.mesh.Init(s.log)
	s.script.Init(&scriptStr, s.log, s.metadata, s.initializeMeshIfReady)
	s.windowShift.Init()
	s.shape.Init(s.mesh, s.log, s.fs, s.grains)
	s.table.Init(s.solver, s.log, s.fs)
	s.solver.Init()
	s.magnetization.Init(s.mesh, s.config, s.geometry)
	s.regions.Init(s.mesh, s.log)
	s.geometry.Init(s.mesh, s.log, s.config, s.magnetization.Value(), s.magnetization.Normalize)
	s.savedQuantities.Init(s.log, s.fs, s.solver)
	s.grains.Init(s.regions.Voronoi)
	s.config.Init(s.mesh)
}

func (s *engineState) start(scriptPath string) {
	// I commented the following line for debugging purposes
	// add it back when the code is stable
	// defer s.cleanExit()
	// The order of the following lines is important
	scriptStr := s.readScript(scriptPath)

	s.initFileSystem(scriptPath)
	s.init(scriptStr)
	s.script.AddToScopeAll(s.fs, s.mesh, s.geometry, s.grains, s.config, s.magnetization, s.metadata, s.regions, s.savedQuantities, s.solver, s.table, s.windowShift, s.shape)

	err := s.script.Parse()
	if err != nil {
		s.log.ErrAndExit("Error parsing script: %v", err)
	}
	// start autosave goroutine before executing the script
	go s.autoFlush()
	err = s.script.Execute()
	if err != nil {
		s.log.ErrAndExit("Error executing script: %v", err)
	}
	s.cleanExit()
}

func (s *engineState) readScript(scriptPath string) string {
	scriptStr := ""
	if scriptPath == "" {
		scriptStr = `
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
	} else {
		scriptBytes, err := os.ReadFile(scriptPath)
		if err != nil {
			s.log.ErrAndExit("Error reading script: %v", err)
		}
		if len(scriptBytes) == 0 {
			s.log.ErrAndExit("Empty input file: %s", scriptPath)
		}
		scriptStr = string(scriptBytes)

	}
	return scriptStr
}

// fs cannot depend on log, so we need to initialize it here
func (s *engineState) initFileSystem(scriptPath string) {
	fs, warn, err := fsutil.NewFileSystem(scriptPath, s.flags.OutputDir, s.flags.SkipExists, s.flags.ForceClean)
	if err != nil {
		s.log.ErrAndExit("Error creating file system: %v", err)
	}
	if warn != "" {
		// this is only for skipping the directory if it already exists with --skip-exist flag
		s.log.Warn("%s", warn)
		// if warn contains "skip-exist", then we must exit
		if strings.Contains(warn, "skip-exist") {
			os.Exit(0)
		}
	}
	s.fs = fs
	if scriptPath == "" {
		s.log.Info("No input files: starting interactive session")
	} else {
		s.log.Info("Input path: %s", scriptPath)
	}
	s.log.Info("Output directory: %s", s.fs.Wd)
}

func (s *engineState) autoFlush() {
	for {
		// sleep first to avoid saving uninitialized data
		time.Sleep(s.autoFlushInterval)
		err := s.metadata.FlushToFile()
		if err != nil {
			s.log.Err("Failed to save metadata to file: %v", err)
		}
		err = s.table.FlushToFile()
		if err != nil {
			s.log.Err("Failed to save table to file: %v", err)
		}
		err = s.log.FlushToFile()
		if err != nil {
			s.log.Err("Failed to save log to file: %v", err)
		}
	}
}

func (s *engineState) cleanExit() {
	s.fs.Drain() // wait for the save queue to finish
	s.table.Close()
	if s.flags.Sync {
		timer.Print(os.Stdout)
	}
	s.metadata.Add("steps", s.solver.NSteps)
	s.metadata.Close()
	s.log.Info("**************** Simulation Ended ****************** //")
	err := s.log.Close()
	if err != nil {
		s.log.Err("Failed to close log file: %v", err)
	}
}

// this is called by the script parser when the mesh is ready to be created
func (s *engineState) initializeMeshIfReady() {
	if s.mesh.ReadyToCreate() {
		s.log.Info("Creating mesh")
		s.mesh.Create()
		s.magnetization.InitializeBuffer()
		s.regions.InitializeBuffer()
		s.metadata.AddMesh(s.mesh)
	}
}
