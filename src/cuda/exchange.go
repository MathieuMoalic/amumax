package cuda

import (
	"unsafe"

	"github.com/MathieuMoalic/amumax/src/data"
	"github.com/MathieuMoalic/amumax/src/mesh"
)

// Add exchange field to Beff.
//
//	m: normalized magnetization
//	B: effective field in Tesla
//	Aex_red: Aex / (Msat * 1e18 m2)
//
// see exchange.cu
func AddExchange(B, m *data.Slice, AexRed SymmLUT, Msat MSlice, regions *Bytes, mesh mesh.MeshLike) {
	c := mesh.CellSize()
	wx := float32(2 / (c[X] * c[X]))
	wy := float32(2 / (c[Y] * c[Y]))
	wz := float32(2 / (c[Z] * c[Z]))
	N := mesh.Size()
	pbc := mesh.PBCCode()
	cfg := make3DConf(N)
	kAddexchangeAsync(B.DevPtr(X), B.DevPtr(Y), B.DevPtr(Z),
		m.DevPtr(X), m.DevPtr(Y), m.DevPtr(Z),
		Msat.DevPtr(0), Msat.Mul(0),
		unsafe.Pointer(AexRed), regions.Ptr,
		wx, wy, wz, N[X], N[Y], N[Z], pbc, cfg)
}

// Finds the average exchange strength around each cell, for debugging.
func ExchangeDecode(dst *data.Slice, AexRed SymmLUT, regions *Bytes, mesh mesh.MeshLike) {
	c := mesh.CellSize()
	wx := float32(2 / (c[X] * c[X]))
	wy := float32(2 / (c[Y] * c[Y]))
	wz := float32(2 / (c[Z] * c[Z]))
	N := mesh.Size()
	pbc := mesh.PBCCode()
	cfg := make3DConf(N)
	kExchangedecodeAsync(dst.DevPtr(0), unsafe.Pointer(AexRed), regions.Ptr, wx, wy, wz, N[X], N[Y], N[Z], pbc, cfg)
}
