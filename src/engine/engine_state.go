package engine

import (
	"os"

	"github.com/MathieuMoalic/amumax/src/log"
	"github.com/MathieuMoalic/amumax/src/parser"
	"github.com/MathieuMoalic/amumax/src/zarr"
)

type EngineStateStruct struct {
	Metadata zarr.Metadata
}

var EngineState EngineStateStruct

func init() {
	EngineState = EngineStateStruct{}
}

func (s *EngineStateStruct) Start(filename string) {
	scriptParser := parser.NewScriptParser()
	backend := parser.NewSimulationBackend()
	script, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	err = scriptParser.Parse(string(script))
	if err != nil {
		panic(err)
	}
	log.Log.Debug("Parsed script")

	// Execute parsed statements
	if err := scriptParser.Execute(backend); err != nil {
		panic(err)
	}
}
