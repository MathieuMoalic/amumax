package engine_old

// Cropped quantity refers to a cut-out piece of a large quantity

import (
	"fmt"

	"github.com/MathieuMoalic/amumax/src/cuda"
	"github.com/MathieuMoalic/amumax/src/engine_old/data_old"
	"github.com/MathieuMoalic/amumax/src/engine_old/log_old"
	"github.com/MathieuMoalic/amumax/src/engine_old/mesh_old"
)

type cropped struct {
	parent                 Quantity
	name                   string
	x1, x2, y1, y2, z1, z2 int
}

// Crop quantity to a box enclosing the given region.
// Used to output a region of interest, even if the region is non-rectangular.
func cropRegion(parent Quantity, region int) *cropped {
	n := MeshOf(parent).Size()
	// use -1 for unset values
	x1, y1, z1 := -1, -1, -1
	x2, y2, z2 := -1, -1, -1
	r := Regions.HostArray()
	for iz := 0; iz < n[Z]; iz++ {
		for iy := 0; iy < n[Y]; iy++ {
			for ix := 0; ix < n[X]; ix++ {
				if r[iz][iy][ix] == byte(region) {
					// initialize all indices if unset
					if x1 == -1 {
						x1, y1, z1 = ix, iy, iz
						x2, y2, z2 = ix, iy, iz
					}
					if ix < x1 {
						x1 = ix
					}
					if iy < y1 {
						y1 = iy
					}
					if iz < z1 {
						z1 = iz
					}
					if ix > x2 {
						x2 = ix
					}
					if iy > y2 {
						y2 = iy
					}
					if iz > z2 {
						z2 = iz
					}
				}
			}
		}
	}
	return crop(parent, x1, x2+1, y1, y2+1, z1, z2+1)
}

func cropLayer(parent Quantity, layer int) *cropped {
	n := MeshOf(parent).Size()
	return crop(parent, 0, n[X], 0, n[Y], layer, layer+1)
}

func cropX(parent Quantity, x1, x2 int) *cropped {
	n := MeshOf(parent).Size()
	return crop(parent, x1, x2, 0, n[Y], 0, n[Z])
}

func cropY(parent Quantity, y1, y2 int) *cropped {
	n := MeshOf(parent).Size()
	return crop(parent, 0, n[X], y1, y2, 0, n[Z])
}

func cropZ(parent Quantity, z1, z2 int) *cropped {
	n := MeshOf(parent).Size()
	return crop(parent, 0, n[X], 0, n[Y], z1, z2)
}

func crop(parent Quantity, x1, x2, y1, y2, z1, z2 int) *cropped {
	n := MeshOf(parent).Size()
	log_old.AssertMsg(x1 < x2 && y1 < y2 && z1 < z2,
		"Invalid crop range: x1 must be less than x2, y1 less than y2, and z1 less than z2 in crop")

	log_old.AssertMsg(x1 >= 0 && y1 >= 0 && z1 >= 0,
		"Invalid crop range: x1, y1, and z1 must be non-negative in crop")

	log_old.AssertMsg(x2 <= n[X] && y2 <= n[Y] && z2 <= n[Z],
		"Invalid crop range: x2, y2, and z2 must be within mesh dimensions in crop")

	name := nameOf(parent)
	if x1 != 0 || x2 != n[X] {
		name += "_x" + rangeStr(x1, x2)
	}
	if y1 != 0 || y2 != n[Y] {
		name += "_y" + rangeStr(y1, y2)
	}
	if z1 != 0 || z2 != n[Z] {
		name += "_z" + rangeStr(z1, z2)
	}

	return &cropped{parent, name, x1, x2, y1, y2, z1, z2}
}

func rangeStr(a, b int) string {
	if a+1 == b {
		return fmt.Sprint(a)
	} else {
		return fmt.Sprint(a, "-", b)
	}
	// (trailing underscore to separate from subsequent autosave number)
}

func (q *cropped) NComp() int                 { return q.parent.NComp() }
func (q *cropped) Name() string               { return q.name }
func (q *cropped) Unit() string               { return unitOf(q.parent) }
func (q *cropped) EvalTo(dst *data_old.Slice) { evalTo(q, dst) }

func (q *cropped) Mesh() *mesh_old.Mesh {
	c := MeshOf(q.parent) // currentMesh
	return mesh_old.NewMesh(q.x2-q.x1, q.y2-q.y1, q.z2-q.z1, c.Dx, c.Dy, c.Dz, c.PBCx, c.PBCy, c.PBCz)
}

func (q *cropped) average() []float64 { return qAverageUniverse(q) } // needed for table
func (q *cropped) Average() []float64 { return q.average() }         // handy for script

func (q *cropped) Slice() (*data_old.Slice, bool) {
	src := ValueOf(q.parent)
	defer cuda.Recycle(src)
	dst := cuda.Buffer(q.NComp(), q.Mesh().Size())
	cuda.Crop(dst, src, q.x1, q.y1, q.z1)
	return dst, true
}
