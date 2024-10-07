package engine

import (
	"math"

	"github.com/MathieuMoalic/amumax/src/cuda"
	"github.com/MathieuMoalic/amumax/src/data"
)

var (
	TopologicalCharge        = newScalarValue("ext_topologicalcharge", "", "2D topological charge", getTopologicalCharge)
	TopologicalChargeDensity = newScalarField("ext_topologicalchargedensity", "1/m2",
		"2D topological charge density m·(∂m/∂x ✕ ∂m/∂y)", setTopologicalChargeDensity)
)

func setTopologicalChargeDensity(dst *data.Slice) {
	cuda.SetTopologicalCharge(dst, normMag.Buffer(), normMag.Mesh())
}

func getTopologicalCharge() float64 {
	s := ValueOf(TopologicalChargeDensity)
	defer cuda.Recycle(s)
	c := getMesh().CellSize()
	N := getMesh().Size()
	return (0.25 * c[X] * c[Y] / math.Pi / float64(N[Z])) * float64(cuda.Sum(s))
}
