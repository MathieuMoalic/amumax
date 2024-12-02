package engine_old

import (
	"fmt"

	"github.com/MathieuMoalic/amumax/src/cuda_old"
	"github.com/MathieuMoalic/amumax/src/engine_old/data_old"
	"github.com/MathieuMoalic/amumax/src/engine_old/log_old"
	"github.com/MathieuMoalic/amumax/src/engine_old/mesh_old"
)

func sOneRegion(q Quantity, r int) *sOneReg {
	log_old.AssertMsg(q.NComp() == 1, "Component mismatch: q must have 1 component in sOneRegion")
	return &sOneReg{oneReg{q, r}}
}

func vOneRegion(q Quantity, r int) *vOneReg {
	log_old.AssertMsg(q.NComp() == 3, "Component mismatch: q must have 3 components in vOneRegion")
	return &vOneReg{oneReg{q, r}}
}

type sOneReg struct{ oneReg }

func (q *sOneReg) Average() float64 { return q.average()[0] }

type vOneReg struct{ oneReg }

func (q *vOneReg) Average() data_old.Vector { return unslice(q.average()) }

// represents a new quantity equal to q in the given region, 0 outside.
type oneReg struct {
	parent Quantity
	region int
}

func inRegion(q Quantity, region int) Quantity {
	return &oneReg{q, region}
}

func (q *oneReg) NComp() int                 { return q.parent.NComp() }
func (q *oneReg) Name() string               { return fmt.Sprint(nameOf(q.parent), ".region", q.region) }
func (q *oneReg) Unit() string               { return unitOf(q.parent) }
func (q *oneReg) Mesh() *mesh_old.Mesh       { return MeshOf(q.parent) }
func (q *oneReg) EvalTo(dst *data_old.Slice) { evalTo(q, dst) }

// returns a new slice equal to q in the given region, 0 outside.
func (q *oneReg) Slice() (*data_old.Slice, bool) {
	src := ValueOf(q.parent)
	defer cuda_old.Recycle(src)
	out := cuda_old.Buffer(q.NComp(), q.Mesh().Size())
	cuda_old.RegionSelect(out, src, Regions.Gpu(), byte(q.region))
	return out, true
}

func (q *oneReg) average() []float64 {
	slice, r := q.Slice()
	if r {
		defer cuda_old.Recycle(slice)
	}
	avg := sAverageUniverse(slice)
	sDiv(avg, Regions.volume(q.region))
	return avg
}

func (q *oneReg) Average() []float64 { return q.average() }

// slice division
func sDiv(v []float64, x float64) {
	for i := range v {
		v[i] /= x
	}
}
