package engine_old

import (
	"reflect"

	"github.com/MathieuMoalic/amumax/src/engine_old/cuda_old"
	"github.com/MathieuMoalic/amumax/src/engine_old/data_old"
	"github.com/MathieuMoalic/amumax/src/engine_old/mesh_old"
)

var Quantities = make(map[string]Quantity)

// Arbitrary physical quantity.
type Quantity interface {
	NComp() int
	EvalTo(dst *data_old.Slice)
}

func addQuantity(name string, value interface{}, doc string) {
	_ = doc
	if v, ok := value.(Quantity); ok {
		Quantities[name] = v
	}
}

func meshSize() [3]int {
	return GetMesh().Size()
}

func sizeOf(q Quantity) [3]int {
	// quantity defines its own, custom, implementation:
	if s, ok := q.(interface {
		Mesh() *mesh_old.Mesh
	}); ok {
		return s.Mesh().Size()
	}
	// otherwise: default mesh
	return meshSize()
}

func AverageOf(q Quantity) []float64 {
	// quantity defines its own, custom, implementation:
	if s, ok := q.(interface {
		average() []float64
	}); ok {
		return s.average()
	}
	// otherwise: default mesh
	buf := ValueOf(q)
	defer cuda_old.Recycle(buf)
	return sAverageMagnet(buf)
}

func nameOf(q Quantity) string {
	// quantity defines its own, custom, implementation:
	if s, ok := q.(interface {
		Name() string
	}); ok {
		return s.Name()
	}
	return "unnamed." + reflect.TypeOf(q).String()
}

func unitOf(q Quantity) string {
	// quantity defines its own, custom, implementation:
	if s, ok := q.(interface {
		Unit() string
	}); ok {
		return s.Unit()
	}
	return "?"
}

func MeshOf(q Quantity) *mesh_old.Mesh {
	// quantity defines its own, custom, implementation:
	if s, ok := q.(interface {
		Mesh() *mesh_old.Mesh
	}); ok {
		return s.Mesh()
	}
	return GetMesh()
}

func ValueOf(q Quantity) *data_old.Slice {
	// TODO: check for Buffered() implementation
	buf := cuda_old.Buffer(q.NComp(), sizeOf(q))
	q.EvalTo(buf)
	return buf
}

// Temporary shim to fit Slice into evalTo
func evalTo(q interface {
	Slice() (*data_old.Slice, bool)
}, dst *data_old.Slice) {
	v, r := q.Slice()
	if r {
		defer cuda_old.Recycle(v)
	}
	data_old.Copy(dst, v)
}
