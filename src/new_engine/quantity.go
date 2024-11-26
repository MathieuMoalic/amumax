package new_engine

import (
	"github.com/MathieuMoalic/amumax/src/data"
)

type Quantities struct {
	list map[string]Quantity
}

func NewQuantities() *Quantities {
	return &Quantities{list: make(map[string]Quantity)}
}

func (qs *Quantities) Add(q Quantity) {
	qs.list[q.Name()] = q
}

type Quantity interface {
	EvalTo(*data.Slice)
	NComp() int
	Size() [3]int
	Average() []float64
	Name() string
	Unit() string
	Value() *data.Slice
}

// type Quantity struct {
// 	engineState *EngineStateStruct
// 	size        [3]int
// }

// func (q *Quantity) EvalTo(dst *data.Slice) {
// }

// func (q *Quantity) NComp() int {
// 	return 0
// }

// //	func (q *Quantity) Mesh() *data.MeshType {
// //		return nil
// //	}

// func (q *Quantity) sizeOf() [3]int {
// 	return q.size
// }

// func (q *Quantity) AverageOf() []float64 {
// 	buf := q.ValueOf()
// 	defer cuda.Recycle(buf)
// 	return q.sAverageMagnet(buf)
// }

// average of slice over the magnet volume
// func (q *Quantity) sAverageMagnet(s *data.Slice) []float64 {
// 	if q.engineState.Geometry.Gpu().IsNil() {
// 		return sAverageUniverse(s)
// 	} else {
// 		avg := make([]float64, s.NComp())
// 		for i := range avg {
// 			avg[i] = float64(cuda.Dot(s.Comp(i), q.engineState.Geometry.Gpu())) / q.magnetNCell()
// 			if math.IsNaN(avg[i]) {
// 				panic("NaN")
// 			}
// 		}
// 		return avg
// 	}
// }

// func (q *Quantity) magnetNCell() float64 {
// 	if q.engineState.Geometry.Gpu().IsNil() {
// 		return float64(q.engineState.Mesh.NCell())
// 	} else {
// 		return float64(cuda.Sum(q.engineState.Geometry.Gpu()))
// 	}
// }

// func (q *Quantity) nameOf() string {
// 	return "unnamed." + reflect.TypeOf(q).String()
// }

// func (q *Quantity) unitOf() string {
// 	return "?"
// }

// // func (q *Quantity) MeshOf(q Quantity) *data.MeshType {
// // 	// quantity defines its own, custom, implementation:
// // 	if s, ok := q.(interface {
// // 		Mesh() *data.MeshType
// // 	}); ok {
// // 		return s.Mesh()
// // 	}
// // 	return GetMesh()
// // }

// func (q *Quantity) ValueOf() *data.Slice {
// 	// TODO: check for Buffered() implementation
// 	buf := cuda.Buffer(q.NComp(), q.sizeOf())
// 	q.EvalTo(buf)
// 	return buf
// }

// // // Temporary shim to fit Slice into evalTo
// // func evalTo(q interface {
// // 	Slice() (*data.Slice, bool)
// // }, dst *data.Slice) {
// // 	v, r := q.Slice()
// // 	if r {
// // 		defer cuda.Recycle(v)
// // 	}
// // 	data.Copy(dst, v)
// // }
