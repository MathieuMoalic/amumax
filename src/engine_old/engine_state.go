package engine_old

import (
	"github.com/MathieuMoalic/amumax/src/engine_old/zarr_old"
)

type EngineStateStruct struct {
	Metadata zarr_old.Metadata
}

var EngineState EngineStateStruct

func init() {
	EngineState = EngineStateStruct{}
}
