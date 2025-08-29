package engine

import (
	"math"

	"github.com/MathieuMoalic/amumax/src/engine/cuda"
	"github.com/MathieuMoalic/amumax/src/engine/data"
)

var (
	TopologicalChargeLattice        = newScalarValue("ext_topologicalchargelattice", "", "2D topological charge according to Berg and Lüscher", getTopologicalChargeLattice)
	TopologicalChargeDensityLattice = newScalarField("ext_topologicalchargedensitylattice", "1/m2",
		"2D topological charge density according to Berg and Lüscher", setTopologicalChargeDensityLattice)
)

func setTopologicalChargeDensityLattice(dst *data.Slice) {
	cuda.SetTopologicalChargeLattice(dst, NormMag.Buffer(), NormMag.Mesh())
}

func getTopologicalChargeLattice() float64 {
	s := ValueOf(TopologicalChargeDensityLattice)
	defer cuda.Recycle(s)
	c := GetMesh().CellSize()
	N := GetMesh().Size()
	return (0.25 * c[X] * c[Y] / math.Pi / float64(N[Z])) * float64(cuda.Sum(s))
}
