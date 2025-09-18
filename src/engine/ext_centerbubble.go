package engine

import (
	"math"

	"github.com/MathieuMoalic/amumax/src/data"
)

func centerBubbleInner() {
	c := GetMesh().CellSize()

	position := bubblePos()
	var centerIdx [2]int
	centerIdx[X] = int(math.Floor((position[X] - getShiftPos()) / c[X]))
	centerIdx[Y] = int(math.Floor((position[Y] - getShiftYPos()) / c[Y]))

	zero := data.Vector{0, 0, 0}
	if shiftMagL == zero || shiftMagR == zero || shiftMagD == zero || shiftMagU == zero {
		shiftMagL[Z] = -BubbleMz
		shiftMagR[Z] = -BubbleMz
		shiftMagD[Z] = -BubbleMz
		shiftMagU[Z] = -BubbleMz
	}

	// put bubble to center
	if centerIdx[X] != 0 {
		shift(-centerIdx[X])
	}
	if centerIdx[Y] != 0 {
		yShift(-centerIdx[Y])
	}
}

// This post-step function centers the simulation window on a bubble
func centerBubble() {
	PostStep(func() { centerBubbleInner() })
}
