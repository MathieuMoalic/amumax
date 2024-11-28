package engine_old

import (
	"github.com/MathieuMoalic/amumax/src/cuda"
	"github.com/MathieuMoalic/amumax/src/data"
)

var (
	MaxAngle  = newScalarValue("MaxAngle", "rad", "maximum angle between neighboring spins", getMaxAngle)
	SpinAngle = newScalarField("spinAngle", "rad", "Angle between neighboring spins", setSpinAngle)
)

func setSpinAngle(dst *data.Slice) {
	cuda.SetMaxAngle(dst, NormMag.Buffer(), lex2.Gpu(), Regions.Gpu(), NormMag.Mesh())
}

func getMaxAngle() float64 {
	s := ValueOf(SpinAngle)
	defer cuda.Recycle(s)
	return float64(cuda.MaxAbs(s)) // just a max would be fine, but not currently implemented
}
