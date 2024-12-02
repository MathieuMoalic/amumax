package window_shift

import (
	"github.com/MathieuMoalic/amumax/src/cuda"
	"github.com/MathieuMoalic/amumax/src/geometry"
	"github.com/MathieuMoalic/amumax/src/magnetization"
	"github.com/MathieuMoalic/amumax/src/mesh"
	"github.com/MathieuMoalic/amumax/src/regions"
	"github.com/MathieuMoalic/amumax/src/slice"
	"github.com/MathieuMoalic/amumax/src/vector"
)

type WindowShift struct {
	totalXShift, totalYShift                   float64
	ShiftMagL, ShiftMagR, ShiftMagU, ShiftMagD vector.Vector // unused for now
	shiftM, shiftGeom, shiftRegions            bool
	edgeCarryShift                             bool
	mesh                                       *mesh.Mesh
	magnetization                              *magnetization.Magnetization
	regions                                    *regions.Regions
	geometry                                   *geometry.Geometry
}

func (w *WindowShift) Init() {}

// position of the window lab frame
// func (w *WindowShift) getShiftXPos() float64 { return -w.TotalXShift }
// func (w *WindowShift) getShiftYPos() float64 { return -w.TotalYShift }
func (w *WindowShift) ShiftX(dx int) {
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

func (w *WindowShift) shiftMagX(m *slice.Slice, dx int) {
	m2 := cuda.Buffer(1, m.Size())
	defer cuda.Recycle(m2)
	for c := 0; c < m.NComp(); c++ {
		comp := m.Comp(c)
		if w.edgeCarryShift {
			cuda.ShiftEdgeCarryX(m2, comp, m.Comp((c+1)%3), m.Comp((c+2)%3), dx, float32(w.ShiftMagL[c]), float32(w.ShiftMagL[c]))
		} else {
			cuda.ShiftX(m2, comp, dx, float32(w.ShiftMagL[c]), float32(w.ShiftMagL[c]))
		}
		slice.Copy(comp, m2) // str0 ?
	}
}

// shift the simulation window over dy cells in Y direction
func (w *WindowShift) ShiftY(dy int) {
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

func (w *WindowShift) shiftMagY(m *slice.Slice, dy int) {
	m2 := cuda.Buffer(1, m.Size())
	defer cuda.Recycle(m2)
	for c := 0; c < m.NComp(); c++ {
		comp := m.Comp(c)
		if w.edgeCarryShift {
			cuda.ShiftEdgeCarryX(m2, comp, m.Comp((c+1)%3), m.Comp((c+2)%3), dy, float32(w.ShiftMagL[c]), float32(w.ShiftMagL[c]))
		} else {
			cuda.ShiftX(m2, comp, dy, float32(w.ShiftMagL[c]), float32(w.ShiftMagL[c]))
		}
		slice.Copy(comp, m2) // str0 ?
	}
}
