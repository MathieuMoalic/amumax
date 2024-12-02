package quantity

import (
	"github.com/MathieuMoalic/amumax/src/slice"
)

type Quantity interface {
	EvalTo(*slice.Slice)
	NComp() int
	Size() [3]int
	Average() []float64
	Name() string
	Unit() string
	Value() *slice.Slice
}
