package new_engine

import (
	"fmt"

	"github.com/MathieuMoalic/amumax/src/data"
)

type Utils struct {
	EngineState *EngineStateStruct
}

func NewUtils(engineState *EngineStateStruct) *Utils {
	u := &Utils{EngineState: engineState}
	u.EngineState.world.RegisterFunction("Print", u.customPrint)
	return u
}

func (u *Utils) Index2Coord(ix, iy, iz int) data.Vector {
	n := u.EngineState.mesh.Size()
	c := u.EngineState.mesh.CellSize()
	x := c[X]*(float64(ix)-0.5*float64(n[X]-1)) - u.EngineState.windowShift.TotalXShift
	y := c[Y]*(float64(iy)-0.5*float64(n[Y]-1)) - u.EngineState.windowShift.TotalYShift
	z := c[Z] * (float64(iz) - 0.5*float64(n[Z]-1))
	return data.Vector{x, y, z}
}

// x range that needs to be refreshed after shift over dx
func (u *Utils) shiftDirtyRange(dx int) (x1, x2 int) {
	Nx := u.EngineState.mesh.Nx
	u.EngineState.log.AssertMsg(dx != 0, "Invalid shift: dx must not be zero in shiftDirtyRange")
	if dx < 0 {
		x1 = Nx + dx
		x2 = Nx
	} else {
		x1 = 0
		x2 = dx
	}
	return
}

// print with special formatting for some known types
func (u *Utils) customPrint(msg ...interface{}) {
	u.EngineState.log.Info("%v", u.customFmt(msg))
}

// mumax specific formatting (Slice -> average, etc).
func (u *Utils) customFmt(msg []interface{}) (fmtMsg string) {
	for _, m := range msg {
		if e, ok := m.(Quantity); ok {
			str := fmt.Sprint(e.Average())
			str = str[1 : len(str)-1] // remove [ ]
			fmtMsg += fmt.Sprintf("%v, ", str)
		} else {
			fmtMsg += fmt.Sprintf("%v, ", m)
		}
	}
	// remove trailing ", "
	if len(fmtMsg) > 2 {
		fmtMsg = fmtMsg[:len(fmtMsg)-2]
	}
	return
}
