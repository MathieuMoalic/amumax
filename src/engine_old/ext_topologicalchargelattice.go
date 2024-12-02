package engine_old

import (
	"math"

	"github.com/MathieuMoalic/amumax/src/cuda_old"
	"github.com/MathieuMoalic/amumax/src/engine_old/data_old"
)

var (
	TopologicalChargeLattice        = newScalarValue("ext_topologicalchargelattice", "", "2D topological charge according to Berg and Lüscher", getTopologicalChargeLattice)
	TopologicalChargeDensityLattice = newScalarField("ext_topologicalchargedensitylattice", "1/m2",
		"2D topological charge density according to Berg and Lüscher", setTopologicalChargeDensityLattice)
)

func setTopologicalChargeDensityLattice(dst *data_old.Slice) {
	cuda_old.SetTopologicalChargeLattice(dst, NormMag.Buffer(), NormMag.Mesh())
}

func getTopologicalChargeLattice() float64 {
	s := ValueOf(TopologicalChargeDensityLattice)
	defer cuda_old.Recycle(s)
	c := GetMesh().CellSize()
	N := GetMesh().Size()
	return (0.25 * c[X] * c[Y] / math.Pi / float64(N[Z])) * float64(cuda_old.Sum(s))
}
