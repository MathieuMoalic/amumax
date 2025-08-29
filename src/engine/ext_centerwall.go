package engine

import (
	"fmt"

	"github.com/MathieuMoalic/amumax/src/engine/data"
)

var (
	DWPos   = newScalarValue("ext_dwpos", "m", "Position of the simulation window while following a domain wall", getShiftPos) // TODO: make more accurate
	DWxPos  = newScalarValue("ext_dwxpos", "m", "Position of the simulation window while following a domain wall", getDWxPos)
	DWSpeed = newScalarValue("ext_dwspeed", "m/s", "Speed of the simulation window while following a domain wall", getShiftSpeed)
)

func centerWallInner(c int) {
	M := &NormMag
	mc := sAverageUniverse(M.Buffer().Comp(c))[0]
	n := GetMesh().Size()
	tolerance := 4 / float64(n[X]) // x*2 * expected <m> change for 1 cell shift

	zero := data.Vector{0, 0, 0}
	if shiftMagL == zero || shiftMagR == zero {
		sign := magsign(M.GetCell(0, n[Y]/2, n[Z]/2)[c])
		shiftMagL[c] = float64(sign)
		shiftMagR[c] = -float64(sign)
	}

	sign := magsign(shiftMagL[c])

	if mc < -tolerance {
		shift(sign)
	} else if mc > tolerance {
		shift(-sign)
	}
}

// This post-step function centers the simulation window on a domain wall
// between up-down (or down-up) domains (like in perpendicular media). E.g.:
//
//	PostStep(CenterPMAWall)
func centerWall(magComp int) {
	PostStep(func() { centerWallInner(magComp) })
}

func magsign(x float64) int {
	if x > 0.1 {
		return 1
	}
	if x < -0.1 {
		return -1
	}
	panic(fmt.Errorf("center wall: unclear in which direction to shift: magnetization at border=%v. Set ShiftMagL, ShiftMagR", x))
}

// used for speed
var (
	lastShift float64 // shift the last time we queried speed
	lastT     float64 // time the last time we queried speed
	lastV     float64 // speed the last time we queried speed
)

func getShiftSpeed() float64 {
	if lastShift != getShiftPos() {
		lastV = (getShiftPos() - lastShift) / (Time - lastT)
		lastShift = getShiftPos()
		lastT = Time
	}
	return lastV
}

func getDWxPos() float64 {
	M := &NormMag
	mx := sAverageUniverse(M.Buffer().Comp(0))[0]
	c := GetMesh().CellSize()
	n := GetMesh().Size()
	position := mx * c[0] * float64(n[0]) / 2.
	return getShiftPos() + position
}
