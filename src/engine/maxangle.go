package engine

import (
	"github.com/MathieuMoalic/amumax/src/cuda"
	"github.com/MathieuMoalic/amumax/src/data"
)

var (
	MaxAngle  = NewScalarValue("MaxAngle", "rad", "maximum angle between neighboring spins", GetMaxAngle)
	SpinAngle = NewScalarField("spinAngle", "rad", "Angle between neighboring spins", SetSpinAngle)
)

func SetSpinAngle(dst *data.Slice) {
	cuda.SetMaxAngle(dst, M.Buffer(), lex2.Gpu(), Regions.Gpu(), M.Mesh())
}

func GetMaxAngle() float64 {
	s := ValueOf(SpinAngle)
	defer cuda.Recycle(s)
	return float64(cuda.MaxAbs(s)) // just a max would be fine, but not currently implemented
}
