package cuda

import (
	"github.com/MathieuMoalic/amumax/src/data"
	"github.com/MathieuMoalic/amumax/src/log"
)

// Add uniaxial magnetocrystalline anisotropy field to Beff.
// see uniaxialanisotropy.cu
func AddCubicAnisotropy2(Beff, m *data.Slice, Msat, k1, k2, k3, c1, c2 MSlice) {
	log.AssertMsg(Beff.Size() == m.Size(), "AddCubicAnisotropy2: Size mismatch between Beff and m slices")

	N := Beff.Len()
	cfg := make1DConf(N)
	k_addcubicanisotropy2_async(
		Beff.DevPtr(X), Beff.DevPtr(Y), Beff.DevPtr(Z),
		m.DevPtr(X), m.DevPtr(Y), m.DevPtr(Z),
		Msat.DevPtr(0), Msat.Mul(0),
		k1.DevPtr(0), k1.Mul(0),
		k2.DevPtr(0), k2.Mul(0),
		k3.DevPtr(0), k3.Mul(0),
		c1.DevPtr(X), c1.Mul(X),
		c1.DevPtr(Y), c1.Mul(Y),
		c1.DevPtr(Z), c1.Mul(Z),
		c2.DevPtr(X), c2.Mul(X),
		c2.DevPtr(Y), c2.Mul(Y),
		c2.DevPtr(Z), c2.Mul(Z),
		N, cfg)
}

// Add uniaxial magnetocrystalline anisotropy field to Beff.
// see uniaxialanisotropy.cu
func AddUniaxialAnisotropy2(Beff, m *data.Slice, Msat, k1, k2, u MSlice) {
	log.AssertMsg(Beff.Size() == m.Size(), "AddUniaxialAnisotropy2: Size mismatch between Beff and m slices")

	checkSize(Beff, m, k1, k2, u, Msat)

	N := Beff.Len()
	cfg := make1DConf(N)

	k_adduniaxialanisotropy2_async(
		Beff.DevPtr(X), Beff.DevPtr(Y), Beff.DevPtr(Z),
		m.DevPtr(X), m.DevPtr(Y), m.DevPtr(Z),
		Msat.DevPtr(0), Msat.Mul(0),
		k1.DevPtr(0), k1.Mul(0),
		k2.DevPtr(0), k2.Mul(0),
		u.DevPtr(X), u.Mul(X),
		u.DevPtr(Y), u.Mul(Y),
		u.DevPtr(Z), u.Mul(Z),
		N, cfg)
}
