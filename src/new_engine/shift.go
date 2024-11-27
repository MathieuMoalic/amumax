package new_engine

import (
	"github.com/MathieuMoalic/amumax/src/cuda"
	"github.com/MathieuMoalic/amumax/src/data"
)

type WindowShift struct {
	EngineState                                *EngineStateStruct
	TotalXShift, TotalYShift                   float64
	ShiftMagL, ShiftMagR, ShiftMagU, ShiftMagD data.Vector
	ShiftM, ShiftGeom, ShiftRegions            bool
	EdgeCarryShift                             bool
}

func NewWindowShift(es *EngineStateStruct) *WindowShift {
	w := new(WindowShift)
	w.EngineState = es
	es.world.RegisterFunction("shift", w.shiftX)
	es.world.RegisterFunction("yshift", w.ShiftY)
	return w
}

// position of the window lab frame
// func (w *WindowShift) getShiftXPos() float64 { return -w.TotalXShift }
// func (w *WindowShift) getShiftYPos() float64 { return -w.TotalYShift }
func (w *WindowShift) shiftX(dx int) {
	w.TotalXShift += float64(dx) * w.EngineState.mesh.Dx
	if w.ShiftM {
		w.shiftMagX(w.EngineState.magnetization.slice, dx)
	}
	if w.ShiftRegions {
		w.EngineState.regions.shift(dx)
	}
	if w.ShiftGeom {
		w.EngineState.geometry.shift(dx)
	}
	w.EngineState.magnetization.normalize()
}

func (w *WindowShift) shiftMagX(m *data.Slice, dx int) {
	m2 := cuda.Buffer(1, m.Size())
	defer cuda.Recycle(m2)
	for c := 0; c < m.NComp(); c++ {
		comp := m.Comp(c)
		if w.EdgeCarryShift {
			cuda.ShiftEdgeCarryX(m2, comp, m.Comp((c+1)%3), m.Comp((c+2)%3), dx, float32(w.ShiftMagL[c]), float32(w.ShiftMagL[c]))
		} else {
			cuda.ShiftX(m2, comp, dx, float32(w.ShiftMagL[c]), float32(w.ShiftMagL[c]))
		}
		data.Copy(comp, m2) // str0 ?
	}
}

// shift the simulation window over dy cells in Y direction
func (w *WindowShift) ShiftY(dy int) {
	w.TotalYShift += float64(dy) * w.EngineState.mesh.Dy // needed to re-init geom, regions
	if w.ShiftM {
		w.shiftMagY(w.EngineState.magnetization.slice, dy)
	}
	if w.ShiftRegions {
		w.EngineState.regions.shiftY(dy)
	}
	if w.ShiftGeom {
		w.EngineState.geometry.shiftY(dy)
	}
	w.EngineState.magnetization.normalize()
}

func (w *WindowShift) shiftMagY(m *data.Slice, dy int) {
	m2 := cuda.Buffer(1, m.Size())
	defer cuda.Recycle(m2)
	for c := 0; c < m.NComp(); c++ {
		comp := m.Comp(c)
		if w.EdgeCarryShift {
			cuda.ShiftEdgeCarryX(m2, comp, m.Comp((c+1)%3), m.Comp((c+2)%3), dy, float32(w.ShiftMagL[c]), float32(w.ShiftMagL[c]))
		} else {
			cuda.ShiftX(m2, comp, dy, float32(w.ShiftMagL[c]), float32(w.ShiftMagL[c]))
		}
		data.Copy(comp, m2) // str0 ?
	}
}
