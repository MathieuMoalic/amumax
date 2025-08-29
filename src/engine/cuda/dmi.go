package cuda

import (
	"unsafe"

	"github.com/MathieuMoalic/amumax/src/engine/data"
	"github.com/MathieuMoalic/amumax/src/engine/log"
	"github.com/MathieuMoalic/amumax/src/mesh"
)

// Add effective field of Dzyaloshinskii-Moriya interaction to Beff (Tesla).
// According to Bagdanov and Röβler, PRL 87, 3, 2001. eq.8 (out-of-plane symmetry breaking).
// See dmi.cu
func AddDMI(Beff *data.Slice, m *data.Slice, Aex_red, Dex_red SymmLUT, Msat MSlice, regions *Bytes, mesh mesh.MeshLike, OpenBC bool) {
	cellsize := mesh.CellSize()
	N := Beff.Size()
	log.AssertMsg(m.Size() == N, "Size mismatch: m and Beff must have the same dimensions in AddDMI")

	cfg := make3DConf(N)
	var openBC byte
	if OpenBC {
		openBC = 1
	}

	k_adddmi_async(Beff.DevPtr(X), Beff.DevPtr(Y), Beff.DevPtr(Z),
		m.DevPtr(X), m.DevPtr(Y), m.DevPtr(Z),
		Msat.DevPtr(0), Msat.Mul(0),
		unsafe.Pointer(Aex_red), unsafe.Pointer(Dex_red), regions.Ptr,
		float32(cellsize[X]), float32(cellsize[Y]), float32(cellsize[Z]), N[X], N[Y], N[Z], mesh.PBC_code(), openBC, cfg)
}
