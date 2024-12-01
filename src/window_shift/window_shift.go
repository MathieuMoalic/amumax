package window_shift

import (
	"github.com/MathieuMoalic/amumax/src/cuda"
	"github.com/MathieuMoalic/amumax/src/data"
	"github.com/MathieuMoalic/amumax/src/geometry"
	"github.com/MathieuMoalic/amumax/src/magnetization"
	"github.com/MathieuMoalic/amumax/src/mesh"
	"github.com/MathieuMoalic/amumax/src/regions"
)

type WindowShift struct {
	totalXShift, totalYShift float64
	shiftMagL                data.Vector
	// shiftMagR, shiftMagU, shiftMagD data.Vector // unused for now
	shiftM, shiftGeom, shiftRegions bool
	edgeCarryShift                  bool
	mesh                            *mesh.Mesh
	magnetization                   *magnetization.Magnetization
	regions                         *regions.Regions
	geometry                        *geometry.Geometry
}

func (w *WindowShift) Init() {
	// es.script.RegisterFunction("shift", w.shiftX)
	// es.script.RegisterFunction("yshift", w.shiftY)
}

func (w *WindowShift) AddToScope() []interface{} {
	return []interface{}{w.shiftX, w.shiftY}
}

// position of the window lab frame
// func (w *WindowShift) getShiftXPos() float64 { return -w.TotalXShift }
// func (w *WindowShift) getShiftYPos() float64 { return -w.TotalYShift }
func (w *WindowShift) shiftX(dx int) {
	w.totalXShift += float64(dx) * w.mesh.Dx
	if w.shiftM {
		w.shiftMagX(w.magnetization.Slice, dx)
	}
	if w.shiftRegions {
		w.regions.Shift(dx)
	}
	if w.shiftGeom {
		w.geometry.Shift(dx)
	}
	w.magnetization.Normalize()
}

func (w *WindowShift) shiftMagX(m *data.Slice, dx int) {
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
func (w *WindowShift) shiftY(dy int) {
	w.totalYShift += float64(dy) * w.mesh.Dy // needed to re-init geom, regions
	if w.shiftM {
		w.shiftMagY(w.magnetization.Slice, dy)
	}
	if w.shiftRegions {
		w.regions.ShiftY(dy)
	}
	if w.shiftGeom {
		w.geometry.ShiftY(dy)
	}
	w.magnetization.Normalize()
}

func (w *WindowShift) shiftMagY(m *data.Slice, dy int) {
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
