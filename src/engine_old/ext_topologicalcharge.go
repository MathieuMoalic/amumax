package engine_old

import (
	"math"

	"github.com/MathieuMoalic/amumax/src/cuda_old"
	"github.com/MathieuMoalic/amumax/src/engine_old/data_old"
)

var (
	TopologicalCharge        = newScalarValue("ext_topologicalcharge", "", "2D topological charge", getTopologicalCharge)
	TopologicalChargeDensity = newScalarField("ext_topologicalchargedensity", "1/m2",
		"2D topological charge density m·(∂m/∂x ✕ ∂m/∂y)", setTopologicalChargeDensity)
)

func setTopologicalChargeDensity(dst *data_old.Slice) {
	cuda_old.SetTopologicalCharge(dst, NormMag.Buffer(), NormMag.Mesh())
}

func getTopologicalCharge() float64 {
	s := ValueOf(TopologicalChargeDensity)
	defer cuda_old.Recycle(s)
	c := GetMesh().CellSize()
	N := GetMesh().Size()
	return (0.25 * c[X] * c[Y] / math.Pi / float64(N[Z])) * float64(cuda_old.Sum(s))
}
