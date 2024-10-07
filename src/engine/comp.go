package engine

// Comp is a Derived Quantity pointing to a single component of vector Quantity

import (
	"fmt"

	"github.com/MathieuMoalic/amumax/src/cuda"
	"github.com/MathieuMoalic/amumax/src/data"
	"github.com/MathieuMoalic/amumax/src/log"
)

type component struct {
	parent Quantity
	comp   int
}

// comp returns vector component c of the parent Quantity
func comp(parent Quantity, c int) ScalarField {
	log.AssertArgument(c >= 0 && c < parent.NComp())
	return AsScalarField(&component{parent, c})
}

func (q *component) NComp() int       { return 1 }
func (q *component) Name() string     { return fmt.Sprint(nameOf(q.parent), "_", compname[q.comp]) }
func (q *component) Unit() string     { return unitOf(q.parent) }
func (q *component) Mesh() *data.Mesh { return MeshOf(q.parent) }

func (q *component) Slice() (*data.Slice, bool) {
	p := q.parent
	src := ValueOf(p)
	defer cuda.Recycle(src)
	c := cuda.Buffer(1, src.Size())
	return c, true
}

func (q *component) EvalTo(dst *data.Slice) {
	src := ValueOf(q.parent)
	defer cuda.Recycle(src)
	data.Copy(dst, src.Comp(q.comp))
}

var compname = map[int]string{0: "x", 1: "y", 2: "z"}
