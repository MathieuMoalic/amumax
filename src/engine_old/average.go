package engine_old

// Averaging of quantities over entire universe or just magnet.

import (
	"github.com/MathieuMoalic/amumax/src/engine_old/cuda_old"
	"github.com/MathieuMoalic/amumax/src/engine_old/data_old"
)

// average of quantity over universe
func qAverageUniverse(q Quantity) []float64 {
	s := ValueOf(q)
	defer cuda_old.Recycle(s)
	return sAverageUniverse(s)
}

// average of slice over universe
func sAverageUniverse(s *data_old.Slice) []float64 {
	nCell := float64(prod(s.Size()))
	avg := make([]float64, s.NComp())
	for i := range avg {
		avg[i] = float64(cuda_old.Sum(s.Comp(i))) / nCell
		checkNaN1(avg[i])
	}
	return avg
}

// average of slice over the magnet volume
func sAverageMagnet(s *data_old.Slice) []float64 {
	if Geometry.Gpu().IsNil() {
		return sAverageUniverse(s)
	} else {
		avg := make([]float64, s.NComp())
		for i := range avg {
			avg[i] = float64(cuda_old.Dot(s.Comp(i), Geometry.Gpu())) / magnetNCell()
			checkNaN1(avg[i])
		}
		return avg
	}
}

// number of cells in the magnet.
// not necessarily integer as cells can have fractional volume.
func magnetNCell() float64 {
	if Geometry.Gpu().IsNil() {
		return float64(GetMesh().NCell())
	} else {
		return float64(cuda_old.Sum(Geometry.Gpu()))
	}
}
