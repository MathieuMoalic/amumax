package quantity

import (
	"github.com/MathieuMoalic/amumax/src/engine_old/data_old"
)

type Quantity interface {
	EvalTo(*data_old.Slice)
	NComp() int
	Size() [3]int
	Average() []float64
	Name() string
	Unit() string
	Value() *data_old.Slice
}
