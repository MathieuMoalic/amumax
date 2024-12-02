package engine_old

import (
	"github.com/MathieuMoalic/amumax/src/cuda_old"
	"github.com/MathieuMoalic/amumax/src/engine_old/data_old"
)

var (
	MaxAngle  = newScalarValue("MaxAngle", "rad", "maximum angle between neighboring spins", getMaxAngle)
	SpinAngle = newScalarField("spinAngle", "rad", "Angle between neighboring spins", setSpinAngle)
)

func setSpinAngle(dst *data_old.Slice) {
	cuda_old.SetMaxAngle(dst, NormMag.Buffer(), lex2.Gpu(), Regions.Gpu(), NormMag.Mesh())
}

func getMaxAngle() float64 {
	s := ValueOf(SpinAngle)
	defer cuda_old.Recycle(s)
	return float64(cuda_old.MaxAbs(s)) // just a max would be fine, but not currently implemented
}
