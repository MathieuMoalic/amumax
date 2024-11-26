package engine

import (
	"github.com/MathieuMoalic/amumax/src/zarr"
)

type EngineStateStruct struct {
	Metadata zarr.Metadata
}

var EngineState EngineStateStruct

func init() {
	EngineState = EngineStateStruct{}
}
