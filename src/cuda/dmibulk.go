package cuda

import (
	"unsafe"

	"github.com/MathieuMoalic/amumax/src/data"
	"github.com/MathieuMoalic/amumax/src/log"
	"github.com/MathieuMoalic/amumax/src/mesh"
)

// Add effective field due to bulk Dzyaloshinskii-Moriya interaction to Beff.
// See dmibulk.cu
func AddDMIBulk(Beff *data.Slice, m *data.Slice, AexRed, DRed SymmLUT, Msat MSlice, regions *Bytes, mesh mesh.MeshLike, OpenBC bool) {
	cellsize := mesh.CellSize()
	N := Beff.Size()
	log.AssertMsg(m.Size() == N, "Size mismatch: m and Beff must have the same dimensions in AddDMIBulk")

	cfg := make3DConf(N)
	var openBC byte
	if OpenBC {
		openBC = 1
	}

	kAdddmibulkAsync(Beff.DevPtr(X), Beff.DevPtr(Y), Beff.DevPtr(Z),
		m.DevPtr(X), m.DevPtr(Y), m.DevPtr(Z),
		Msat.DevPtr(0), Msat.Mul(0),
		unsafe.Pointer(AexRed), unsafe.Pointer(DRed), regions.Ptr,
		float32(cellsize[X]), float32(cellsize[Y]), float32(cellsize[Z]), N[X], N[Y], N[Z], mesh.PBCCode(), openBC, cfg)
}
