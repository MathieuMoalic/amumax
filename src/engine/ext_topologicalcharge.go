package engine

import (
	"math"

	"github.com/MathieuMoalic/amumax/src/engine/cuda"
	"github.com/MathieuMoalic/amumax/src/engine/data"
)

var (
	TopologicalCharge        = newScalarValue("ext_topologicalcharge", "", "2D topological charge", getTopologicalCharge)
	TopologicalChargeDensity = newScalarField("ext_topologicalchargedensity", "1/m2",
		"2D topological charge density m·(∂m/∂x ✕ ∂m/∂y)", setTopologicalChargeDensity)
)

func setTopologicalChargeDensity(dst *data.Slice) {
	cuda.SetTopologicalCharge(dst, NormMag.Buffer(), NormMag.Mesh())
}

func getTopologicalCharge() float64 {
	s := ValueOf(TopologicalChargeDensity)
	defer cuda.Recycle(s)
	c := GetMesh().CellSize()
	N := GetMesh().Size()
	return (0.25 * c[X] * c[Y] / math.Pi / float64(N[Z])) * float64(cuda.Sum(s))
}
