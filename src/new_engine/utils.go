package new_engine

import "github.com/MathieuMoalic/amumax/src/data"

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
	x := c[X]*(float64(ix)-0.5*float64(n[X]-1)) - totalShift
	y := c[Y]*(float64(iy)-0.5*float64(n[Y]-1)) - totalYShift
	z := c[Z] * (float64(iz) - 0.5*float64(n[Z]-1))
	return data.Vector{x, y, z}
}
