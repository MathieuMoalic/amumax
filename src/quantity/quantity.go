package quantity

import (
	"github.com/MathieuMoalic/amumax/src/data"
)

type Quantity interface {
	EvalTo(*data.Slice)
	NComp() int
	Size() [3]int
	Average() []float64
	Name() string
	Unit() string
	Value() *data.Slice
}
