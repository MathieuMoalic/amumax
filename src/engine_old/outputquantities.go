package engine_old

/*
The metadata layer wraps basic micromagnetic functions (e.g. func SetDemagField())
in objects that provide:

- additional information (Name, Unit, ...) used for saving output,
- additional methods (Comp, Region, ...) handy for input scripting.
*/

import (
	"fmt"

	"github.com/MathieuMoalic/amumax/src/engine_old/cuda_old"
	"github.com/MathieuMoalic/amumax/src/engine_old/data_old"
	"github.com/MathieuMoalic/amumax/src/engine_old/mesh_old"
)

// info provides an Info implementation intended for embedding in other types.
type info struct {
	nComp int
	name  string
	unit  string
}

func (i *info) Name() string { return i.name }
func (i *info) Unit() string { return i.unit }
func (i *info) NComp() int   { return i.nComp }

// outputFunc is an outputValue implementation where a function provides the output value.
// It can be scalar or vector.
// Used internally by NewScalarValue and NewVectorValue.
type valueFunc struct {
	info
	f func() []float64
}

func (g *valueFunc) get() []float64     { return g.f() }
func (g *valueFunc) average() []float64 { return g.get() }
func (g *valueFunc) EvalTo(dst *data_old.Slice) {
	v := g.get()
	for c, v := range v {
		cuda_old.Memset(dst.Comp(c), float32(v))
	}
}

// ScalarValue enhances an outputValue with methods specific to
// a space-independent scalar quantity (e.g. total energy).
type ScalarValue struct {
	*valueFunc
}

// newScalarValue constructs an outputable space-independent scalar quantity whose
// value is provided by function f.
func newScalarValue(name, unit, desc string, f func() float64) *ScalarValue {
	g := func() []float64 { return []float64{f()} }
	v := &ScalarValue{&valueFunc{info{1, name, unit}, g}}
	export(v, desc)
	return v
}

func (s ScalarValue) Get() float64     { return s.average()[0] }
func (s ScalarValue) Average() float64 { return s.Get() }

// VectorValue enhances an outputValue with methods specific to
// a space-independent vector quantity (e.g. averaged magnetization).
type VectorValue struct {
	*valueFunc
}

// newVectorValue constructs an outputable space-independent vector quantity whose
// value is provided by function f.
func newVectorValue(name, unit, desc string, f func() []float64) *VectorValue {
	v := &VectorValue{&valueFunc{info{3, name, unit}, f}}
	export(v, desc)
	return v
}

func (v *VectorValue) Get() data_old.Vector     { return unslice(v.average()) }
func (v *VectorValue) Average() data_old.Vector { return v.Get() }

// newVectorField constructs an outputable space-dependent vector quantity whose
// value is provided by function f.
func newVectorField(name, unit, desc string, f func(dst *data_old.Slice)) VectorField {
	v := AsVectorField(&fieldFunc{info{3, name, unit}, f})
	declROnly(name, v, cat(desc, unit))
	return v
}

// NewVectorField constructs an outputable space-dependent scalar quantity whose
// value is provided by function f.
func newScalarField(name, unit, desc string, f func(dst *data_old.Slice)) ScalarField {
	q := AsScalarField(&fieldFunc{info{1, name, unit}, f})
	declROnly(name, q, cat(desc, unit))
	return q
}

type fieldFunc struct {
	info
	f func(*data_old.Slice)
}

func (c *fieldFunc) Mesh() *mesh_old.Mesh       { return GetMesh() }
func (c *fieldFunc) average() []float64         { return qAverageUniverse(c) }
func (c *fieldFunc) EvalTo(dst *data_old.Slice) { evalTo(c, dst) }

// Calculates and returns the quantity.
// recycle is true: slice needs to be recycled.
func (q *fieldFunc) Slice() (s *data_old.Slice, recycle bool) {
	buf := cuda_old.Buffer(q.NComp(), q.Mesh().Size())
	cuda_old.Zero(buf)
	q.f(buf)
	return buf, true
}

// ScalarField enhances an outputField with methods specific to
// a space-dependent scalar quantity.
type ScalarField struct {
	Quantity
}

// AsScalarField promotes a quantity to a ScalarField,
// enabling convenience methods particular to scalars.
func AsScalarField(q Quantity) ScalarField {
	if q.NComp() != 1 {
		panic(fmt.Errorf("ScalarField(%v): need 1 component, have: %v", nameOf(q), q.NComp()))
	}
	return ScalarField{q}
}

func (s ScalarField) average() []float64       { return AverageOf(s.Quantity) }
func (s ScalarField) Average() float64         { return s.average()[0] }
func (s ScalarField) Region(r int) ScalarField { return AsScalarField(inRegion(s.Quantity, r)) }
func (s ScalarField) Name() string             { return nameOf(s.Quantity) }
func (s ScalarField) Unit() string             { return unitOf(s.Quantity) }

// VectorField enhances an outputField with methods specific to
// a space-dependent vector quantity.
type VectorField struct {
	Quantity
}

// AsVectorField promotes a quantity to a VectorField,
// enabling convenience methods particular to vectors.
func AsVectorField(q Quantity) VectorField {
	if q.NComp() != 3 {
		panic(fmt.Errorf("VectorField(%v): need 3 components, have: %v", nameOf(q), q.NComp()))
	}
	return VectorField{q}
}

func (v VectorField) average() []float64       { return AverageOf(v.Quantity) }
func (v VectorField) Average() data_old.Vector { return unslice(v.average()) }
func (v VectorField) Region(r int) VectorField { return AsVectorField(inRegion(v.Quantity, r)) }
func (v VectorField) Comp(c int) ScalarField   { return AsScalarField(comp(v.Quantity, c)) }
func (v VectorField) Mesh() *mesh_old.Mesh     { return MeshOf(v.Quantity) }
func (v VectorField) Name() string             { return nameOf(v.Quantity) }
func (v VectorField) Unit() string             { return unitOf(v.Quantity) }
func (v VectorField) HostCopy() *data_old.Slice {
	s := ValueOf(v.Quantity)
	defer cuda_old.Recycle(s)
	return s.HostCopy()
}
