package engine

import (
	"fmt"
	"os"
	"time"

	"github.com/MathieuMoalic/amumax/src/flags"
	"github.com/MathieuMoalic/amumax/src/fsutil"
	"github.com/MathieuMoalic/amumax/src/log"
	"github.com/MathieuMoalic/amumax/src/log_old"
	"github.com/MathieuMoalic/amumax/src/mesh"
	"github.com/MathieuMoalic/amumax/src/metadata"
	"github.com/MathieuMoalic/amumax/src/script"
	"github.com/MathieuMoalic/amumax/src/timer"
	"github.com/fatih/color"
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
	utils           *utils
	windowShift     *windowShift
	shape           *shapeList
	grains          *grains
	config          *configList
	script          *script.ScriptParser
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
	s.run(mx3path, string(scriptBytes))
}

func (s *engineState) startInteractive() {
	log_old.Log.Info("No input files: starting interactive session")
	scriptStr := `
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
	now := time.Now()
	fakeMx3Path := fmt.Sprintf("/tmp/amumax-%v-%02d-%02d_%02dh%02d.zarr", now.Year(), int(now.Month()), now.Day(), now.Hour(), now.Minute())
	s.run(fakeMx3Path, scriptStr)
}

func (s *engineState) run(scriptPath, scriptStr string) {
	// defer s.cleanExit()
	// The order of the following lines is important
	s.fs = fsutil.NewFileSystem(scriptPath, flags.Flags.OutputDir, flags.Flags.SkipExists, flags.Flags.ForceClean)
	s.log = log.NewLogs(s.fs, s.flags.Debug)
	s.metadata = metadata.NewMetadata(s.fs, s.log)
	s.mesh = &mesh.Mesh{}
	s.script = script.NewScriptParser(&scriptStr, s.log, s.metadata, s.initializeMeshIfReady)

	s.windowShift = newWindowShift(s)
	s.shape = newShape(s)
	s.table = newTable(s)
	s.solver = newSolver(s)
	s.magnetization = newMagnetization(s)
	s.regions = newRegions(s)
	s.geometry = newGeom(s)
	s.savedQuantities = newSavedQuantities(s)
	s.utils = newUtils(s)
	s.grains = newGrains(s)
	s.config = newConfigList(s.mesh, s.script)
	s.script.RegisterMesh(s.mesh)
	err := s.script.Parse()
	if err != nil {
		s.log.ErrAndExit("Error parsing script: %v", err)
	}

	err = s.script.Execute()
	if err != nil {
		s.log.ErrAndExit("Error executing script: %v", err)
	}
	s.cleanExit()
}

func (s *engineState) cleanExit() {
	s.fs.Drain()    // wait for the save queue to finish
	s.table.flush() // flush table to disk
	if s.flags.Sync {
		timer.Print(os.Stdout)
	}
	s.metadata.Add("steps", s.solver.NSteps)
	s.metadata.End()
	s.log.Info("**************** Simulation Ended ****************** //")
	s.log.FlushToFile()
}

func (s *engineState) initializeMeshIfReady() {
	if s.mesh.ReadyToCreate() {
		s.mesh.Create()
		s.magnetization.initializeBuffer()
		s.regions.initializeBuffer()
		s.metadata.AddMesh(s.mesh)
	}
}
