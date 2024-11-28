package new_engine

import (
	"github.com/MathieuMoalic/amumax/src/cuda"
	"github.com/MathieuMoalic/amumax/src/data"
)

type windowShift struct {
	e                        *engineState
	totalXShift, totalYShift float64
	shiftMagL                data.Vector
	// shiftMagR, shiftMagU, shiftMagD data.Vector // unused for now
	shiftM, shiftGeom, shiftRegions bool
	edgeCarryShift                  bool
}

func newWindowShift(es *engineState) *windowShift {
	w := new(windowShift)
	w.e = es
	es.world.registerFunction("shift", w.shiftX)
	es.world.registerFunction("yshift", w.shiftY)
	return w
}

// position of the window lab frame
// func (w *WindowShift) getShiftXPos() float64 { return -w.TotalXShift }
// func (w *WindowShift) getShiftYPos() float64 { return -w.TotalYShift }
func (w *windowShift) shiftX(dx int) {
	w.totalXShift += float64(dx) * w.e.mesh.Dx
	if w.shiftM {
		w.shiftMagX(w.e.magnetization.slice, dx)
	}
	if w.shiftRegions {
		w.e.regions.shift(dx)
	}
	if w.shiftGeom {
		w.e.geometry.shift(dx)
	}
	w.e.magnetization.normalize()
}

func (w *windowShift) shiftMagX(m *data.Slice, dx int) {
	m2 := cuda.Buffer(1, m.Size())
	defer cuda.Recycle(m2)
	for c := 0; c < m.NComp(); c++ {
		comp := m.Comp(c)
		if w.edgeCarryShift {
			cuda.ShiftEdgeCarryX(m2, comp, m.Comp((c+1)%3), m.Comp((c+2)%3), dx, float32(w.shiftMagL[c]), float32(w.shiftMagL[c]))
		} else {
			cuda.ShiftX(m2, comp, dx, float32(w.shiftMagL[c]), float32(w.shiftMagL[c]))
		}
		data.Copy(comp, m2) // str0 ?
	}
}

// shift the simulation window over dy cells in Y direction
func (w *windowShift) shiftY(dy int) {
	w.totalYShift += float64(dy) * w.e.mesh.Dy // needed to re-init geom, regions
	if w.shiftM {
		w.shiftMagY(w.e.magnetization.slice, dy)
	}
	if w.shiftRegions {
		w.e.regions.shiftY(dy)
	}
	if w.shiftGeom {
		w.e.geometry.shiftY(dy)
	}
	w.e.magnetization.normalize()
}

func (w *windowShift) shiftMagY(m *data.Slice, dy int) {
	m2 := cuda.Buffer(1, m.Size())
	defer cuda.Recycle(m2)
	for c := 0; c < m.NComp(); c++ {
		comp := m.Comp(c)
		if w.edgeCarryShift {
			cuda.ShiftEdgeCarryX(m2, comp, m.Comp((c+1)%3), m.Comp((c+2)%3), dy, float32(w.shiftMagL[c]), float32(w.shiftMagL[c]))
		} else {
			cuda.ShiftX(m2, comp, dy, float32(w.shiftMagL[c]), float32(w.shiftMagL[c]))
		}
		data.Copy(comp, m2) // str0 ?
	}
}
