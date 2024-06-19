package engine

import (
	"reflect"

	"github.com/MathieuMoalic/amumax/cuda"
	"github.com/MathieuMoalic/amumax/data"
)

var Quantities = make(map[string]Quantity)

// Arbitrary physical quantity.
type Quantity interface {
	NComp() int
	EvalTo(dst *data.Slice)
}

func AddQuantity(name string, value interface{}, doc string) {
	if v, ok := value.(Quantity); ok {
		Quantities[name] = v
	}
}

func MeshSize() [3]int {
	return GetMesh().Size()
}

func SizeOf(q Quantity) [3]int {
	// quantity defines its own, custom, implementation:
	if s, ok := q.(interface {
		Mesh() *data.Mesh
	}); ok {
		return s.Mesh().Size()
	}
	// otherwise: default mesh
	return MeshSize()
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
	defer cuda.Recycle(buf)
	return sAverageMagnet(buf)
}

func NameOf(q Quantity) string {
	// quantity defines its own, custom, implementation:
	if s, ok := q.(interface {
		Name() string
	}); ok {
		return s.Name()
	}
	return "unnamed." + reflect.TypeOf(q).String()
}

func UnitOf(q Quantity) string {
	// quantity defines its own, custom, implementation:
	if s, ok := q.(interface {
		Unit() string
	}); ok {
		return s.Unit()
	}
	return "?"
}

func MeshOf(q Quantity) *data.Mesh {
	// quantity defines its own, custom, implementation:
	if s, ok := q.(interface {
		Mesh() *data.Mesh
	}); ok {
		return s.Mesh()
	}
	return GetMesh()
}

func ValueOf(q Quantity) *data.Slice {
	// TODO: check for Buffered() implementation
	buf := cuda.Buffer(q.NComp(), SizeOf(q))
	q.EvalTo(buf)
	return buf
}

// Temporary shim to fit Slice into EvalTo
func EvalTo(q interface {
	Slice() (*data.Slice, bool)
}, dst *data.Slice) {
	v, r := q.Slice()
	if r {
		defer cuda.Recycle(v)
	}
	data.Copy(dst, v)
}
