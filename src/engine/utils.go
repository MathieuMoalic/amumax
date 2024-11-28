package engine

import (
	"encoding/binary"
	"fmt"
	"math"

	"github.com/MathieuMoalic/amumax/src/cuda"
	"github.com/MathieuMoalic/amumax/src/data"
)

type utils struct {
	e *engineState
}

func newUtils(engineState *engineState) *utils {
	u := &utils{e: engineState}
	u.e.script.RegisterFunction("Print", u.customPrint)
	return u
}

func (u *utils) Index2Coord(ix, iy, iz int) data.Vector {
	n := u.e.mesh.Size()
	c := u.e.mesh.CellSize()
	x := c[X]*(float64(ix)-0.5*float64(n[X]-1)) - u.e.windowShift.totalXShift
	y := c[Y]*(float64(iy)-0.5*float64(n[Y]-1)) - u.e.windowShift.totalYShift
	z := c[Z] * (float64(iz) - 0.5*float64(n[Z]-1))
	return data.Vector{x, y, z}
}

// x range that needs to be refreshed after shift over dx
func (u *utils) shiftDirtyRange(dx int) (x1, x2 int) {
	Nx := u.e.mesh.Nx
	u.e.log.AssertMsg(dx != 0, "Invalid shift: dx must not be zero in shiftDirtyRange")
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
func (u *utils) customPrint(msg ...interface{}) {
	u.e.log.Info("%v", u.customFmt(msg))
}

// mumax specific formatting (Slice -> average, etc).
func (u *utils) customFmt(msg []interface{}) (fmtMsg string) {
	for _, m := range msg {
		if e, ok := m.(quantity); ok {
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

func prod(s [3]int) int {
	return s[0] * s[1] * s[2]
}

// average of slice over universe
func averageSlice(s *data.Slice) []float64 {
	nCell := float64(prod(s.Size()))
	avg := make([]float64, s.NComp())
	for i := range avg {
		avg[i] = float64(cuda.Sum(s.Comp(i))) / nCell
		if math.IsNaN(avg[i]) {
			panic("NaN")
		}
	}
	return avg
}
func sign(x float64) float64 {
	switch {
	case x > 0:
		return 1
	case x < 0:
		return -1
	default:
		return 0
	}
}

func fmod(a, b float64) float64 {
	if b == 0 || math.IsInf(b, 1) {
		return a
	}
	if math.Abs(a) > b/2 {
		return sign(a) * (math.Mod(math.Abs(a+b/2), b) - b/2)
	} else {
		return a
	}
}

func sqr64(x float64) float64 { return x * x }

func float64ToBytes(f float64) []byte {
	var buf [8]byte
	binary.LittleEndian.PutUint64(buf[:], math.Float64bits(f))
	return buf[:]
}

// func bytesToFloat64(bytes []byte) float64 {
// 	bits := binary.LittleEndian.Uint64(bytes)
// 	float := math.Float64frombits(bits)
// 	return float
// }

// func bytesToFloat32(bytes []byte) float32 {
// 	bits := binary.LittleEndian.Uint32(bytes)
// 	float := math.Float32frombits(bits)
// 	return float
// }

// func float32ToBytes(f float32) []byte {
// 	var buf [4]byte
// 	binary.LittleEndian.PutUint32(buf[:], math.Float32bits(f))
// 	return buf[:]
// }
