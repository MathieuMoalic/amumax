package engine_old

import (
	"math"
)

var (
	BubblePos   = newVectorValue("ext_bubblepos", "m", "Bubble core position", bubblePos)
	BubbleDist  = newScalarValue("ext_bubbledist", "m", "Bubble traveled distance", bubbleDist)
	BubbleSpeed = newScalarValue("ext_bubblespeed", "m/s", "Bubble velocity", bubbleSpeed)
	BubbleMz    = 1.0
)

func bubblePos() []float64 {
	m := NormMag.Buffer()
	n := GetMesh().Size()
	c := GetMesh().CellSize()
	mz := m.Comp(Z).HostCopy().Scalars()[0]

	posx, posy := 0., 0.

	if BubbleMz != -1.0 && BubbleMz != 1.0 {
		panic("ext_BubbleMz should be 1.0 or -1.0")
	}

	{
		var magsum float32
		var weightedsum float32

		for iy := range mz {
			for ix := range mz[0] {
				magsum += ((mz[iy][ix]*float32(BubbleMz) + 1.) / 2.)
				weightedsum += ((mz[iy][ix]*float32(BubbleMz) + 1.) / 2.) * float32(iy)
			}
		}
		posy = float64(weightedsum / magsum)
	}

	{
		var magsum float32
		var weightedsum float32

		for ix := range mz[0] {
			for iy := range mz {
				magsum += ((mz[iy][ix]*float32(BubbleMz) + 1.) / 2.)
				weightedsum += ((mz[iy][ix]*float32(BubbleMz) + 1.) / 2.) * float32(ix)
			}
		}
		posx = float64(weightedsum / magsum)
	}

	return []float64{(posx-float64(n[X]/2))*c[X] + getShiftPos(), (posy-float64(n[Y]/2))*c[Y] + getShiftYPos(), 0}
}

var (
	prevBpos = [2]float64{-1e99, -1e99}
	bdist    = 0.0
)

func bubbleDist() float64 {
	pos := bubblePos()
	if prevBpos == [2]float64{-1e99, -1e99} {
		prevBpos = [2]float64{pos[X], pos[Y]}
		return 0
	}

	w := GetMesh().WorldSize()
	dx := pos[X] - prevBpos[X]
	dy := pos[Y] - prevBpos[Y]
	prevBpos = [2]float64{pos[X], pos[Y]}

	// PBC wrap
	if dx > w[X]/2 {
		dx -= w[X]
	}
	if dx < -w[X]/2 {
		dx += w[X]
	}
	if dy > w[Y]/2 {
		dy -= w[Y]
	}
	if dy < -w[Y]/2 {
		dy += w[Y]
	}

	bdist += math.Sqrt(dx*dx + dy*dy)
	return bdist
}

var (
	prevBdist = 0.0
	prevBt    = -999.0
)

func bubbleSpeed() float64 {
	dist := bubbleDist()

	if prevBt < 0 {
		prevBdist = dist
		prevBt = Time
		return 0
	}

	v := (dist - prevBdist) / (Time - prevBt)

	prevBt = Time
	prevBdist = dist

	return v
}
