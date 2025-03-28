package cuda_old

import (
	"unsafe"

	"github.com/MathieuMoalic/amumax/src/engine_old/data_old"
	"github.com/MathieuMoalic/amumax/src/engine_old/log_old"
	"github.com/MathieuMoalic/amumax/src/mesh"
)

// Add effective field due to bulk Dzyaloshinskii-Moriya interaction to Beff.
// See dmibulk.cu
func AddDMIBulk(Beff *data_old.Slice, m *data_old.Slice, Aex_red, D_red SymmLUT, Msat MSlice, regions *Bytes, mesh mesh.MeshLike, OpenBC bool) {
	cellsize := mesh.CellSize()
	N := Beff.Size()
	log_old.AssertMsg(m.Size() == N, "Size mismatch: m and Beff must have the same dimensions in AddDMIBulk")

	cfg := make3DConf(N)
	var openBC byte
	if OpenBC {
		openBC = 1
	}

	k_adddmibulk_async(Beff.DevPtr(X), Beff.DevPtr(Y), Beff.DevPtr(Z),
		m.DevPtr(X), m.DevPtr(Y), m.DevPtr(Z),
		Msat.DevPtr(0), Msat.Mul(0),
		unsafe.Pointer(Aex_red), unsafe.Pointer(D_red), regions.Ptr,
		float32(cellsize[X]), float32(cellsize[Y]), float32(cellsize[Z]), N[X], N[Y], N[Z], mesh.PBC_code(), openBC, cfg)
}
