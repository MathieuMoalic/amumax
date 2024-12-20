package engine_old

// Comp is a Derived Quantity pointing to a single component of vector Quantity

import (
	"fmt"

	"github.com/MathieuMoalic/amumax/src/engine_old/cuda_old"
	"github.com/MathieuMoalic/amumax/src/engine_old/data_old"
	"github.com/MathieuMoalic/amumax/src/engine_old/log_old"
	"github.com/MathieuMoalic/amumax/src/engine_old/mesh_old"
)

type component struct {
	parent Quantity
	comp   int
}

// comp returns vector component c of the parent Quantity
func comp(parent Quantity, c int) ScalarField {
	log_old.AssertMsg(c >= 0 && c < parent.NComp(),
		"Invalid component: component index c must be between 0 and NComp - 1 in comp")
	return AsScalarField(&component{parent, c})
}

func (q *component) NComp() int           { return 1 }
func (q *component) Name() string         { return fmt.Sprint(nameOf(q.parent), "_", compname[q.comp]) }
func (q *component) Unit() string         { return unitOf(q.parent) }
func (q *component) Mesh() *mesh_old.Mesh { return MeshOf(q.parent) }

func (q *component) Slice() (*data_old.Slice, bool) {
	p := q.parent
	src := ValueOf(p)
	defer cuda_old.Recycle(src)
	c := cuda_old.Buffer(1, src.Size())
	return c, true
}

func (q *component) EvalTo(dst *data_old.Slice) {
	src := ValueOf(q.parent)
	defer cuda_old.Recycle(src)
	data_old.Copy(dst, src.Comp(q.comp))
}

var compname = map[int]string{0: "x", 1: "y", 2: "z"}
