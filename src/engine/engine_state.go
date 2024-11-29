package engine

import (
	"os"
	"strings"
	"time"

	"github.com/MathieuMoalic/amumax/src/flags"
	"github.com/MathieuMoalic/amumax/src/fsutil"
	"github.com/MathieuMoalic/amumax/src/log"
	"github.com/MathieuMoalic/amumax/src/mesh"
	"github.com/MathieuMoalic/amumax/src/metadata"
	"github.com/MathieuMoalic/amumax/src/script"
	"github.com/MathieuMoalic/amumax/src/timer"
)

type engineState struct {
	flags    *flags.FlagsType
	fs       *fsutil.FileSystem
	log      *log.Logs
	metadata *metadata.Metadata

	table           *table
	solver          *solver
	mesh            *mesh.Mesh
	magnetization   *magnetization
	geometry        *geometry
	regions         *regions
	savedQuantities *savedQuantities
	windowShift     *windowShift
	shape           *shapeList
	grains          *grains
	config          *configList
	script          *script.ScriptParser

	autoFlushInterval time.Duration
}

func newEngineState(givenFlags *flags.FlagsType, log *log.Logs) *engineState {
	return &engineState{flags: givenFlags, log: log, autoFlushInterval: 5 * time.Second}
}

func (s *engineState) start(scriptPath string) {
	// I commented the following line for debugging purposes
	// add it back when the code is stable
	// defer s.cleanExit()
	// The order of the following lines is important
	scriptStr := s.readScript(scriptPath)
	s.initFileSystem(scriptPath)
	s.metadata = metadata.NewMetadata(s.fs, s.log)
	s.mesh = mesh.NewMesh(s.log)
	s.script = script.NewScriptParser(&scriptStr, s.log, s.metadata, s.initializeMeshIfReady)

	s.script.RegisterMesh(s.mesh)
	s.windowShift = newWindowShift(s)
	s.shape = newShape(s)
	s.table = newTable(s)
	s.solver = newSolver(s)
	s.magnetization = newMagnetization(s)
	s.regions = newRegions(s)
	s.geometry = newGeom(s)
	s.savedQuantities = newSavedQuantities(s)
	s.grains = newGrains(s)
	s.config = newConfigList(s.mesh, s.script)
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
	fs, warn, err := fsutil.NewFileSystem(scriptPath, flags.Flags.OutputDir, flags.Flags.SkipExists, flags.Flags.ForceClean)
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
		s.metadata.FlushToFile()
		s.table.flushToFile()
		s.log.FlushToFile()
		time.Sleep(s.autoFlushInterval)
	}
}

func (s *engineState) cleanExit() {
	s.fs.Drain() // wait for the save queue to finish
	s.table.close()
	if s.flags.Sync {
		timer.Print(os.Stdout)
	}
	s.metadata.Add("steps", s.solver.NSteps)
	s.metadata.Close()
	s.log.Info("**************** Simulation Ended ****************** //")
	s.log.Close()
}

// this is called by the script parser when the mesh is ready to be created
func (s *engineState) initializeMeshIfReady() {
	if s.mesh.ReadyToCreate() {
		s.mesh.Create()
		s.magnetization.initializeBuffer()
		s.regions.initializeBuffer()
		s.metadata.AddMesh(s.mesh)
	}
}
