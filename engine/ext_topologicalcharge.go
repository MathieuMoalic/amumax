package engine

import (
	"math"

	"github.com/MathieuMoalic/amumax/cuda"
	"github.com/MathieuMoalic/amumax/data"
)

var (
	Ext_TopologicalCharge        = NewScalarValue("ext_topologicalcharge", "", "2D topological charge", GetTopologicalCharge)
	Ext_TopologicalChargeDensity = NewScalarField("ext_topologicalchargedensity", "1/m2",
		"2D topological charge density m·(∂m/∂x ✕ ∂m/∂y)", SetTopologicalChargeDensity)
)

func SetTopologicalChargeDensity(dst *data.Slice) {
	cuda.SetTopologicalCharge(dst, M.Buffer(), M.Mesh())
}

func GetTopologicalCharge() float64 {
	s := ValueOf(Ext_TopologicalChargeDensity)
	defer cuda.Recycle(s)
	c := GetMesh().CellSize()
	N := GetMesh().Size()
	return (0.25 * c[X] * c[Y] / math.Pi / float64(N[Z])) * float64(cuda.Sum(s))
}
