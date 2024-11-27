package new_engine

import (
	"github.com/MathieuMoalic/amumax/src/data"
)

type Utils struct {
	EngineState *EngineStateStruct
}

func NewUtils(engineState *EngineStateStruct) *Utils {
	u := &Utils{EngineState: engineState}
	return u
}

func (u *Utils) Index2Coord(ix, iy, iz int) data.Vector {
	n := u.EngineState.Mesh.Size()
	c := u.EngineState.Mesh.CellSize()
	x := c[X]*(float64(ix)-0.5*float64(n[X]-1)) - u.EngineState.WindowShift.TotalXShift
	y := c[Y]*(float64(iy)-0.5*float64(n[Y]-1)) - u.EngineState.WindowShift.TotalYShift
	z := c[Z] * (float64(iz) - 0.5*float64(n[Z]-1))
	return data.Vector{x, y, z}
}

// x range that needs to be refreshed after shift over dx
func (u *Utils) shiftDirtyRange(dx int) (x1, x2 int) {
	Nx := u.EngineState.Mesh.Nx
	u.EngineState.Log.AssertMsg(dx != 0, "Invalid shift: dx must not be zero in shiftDirtyRange")
	if dx < 0 {
		x1 = Nx + dx
		x2 = Nx
	} else {
		x1 = 0
		x2 = dx
	}
	return
}
