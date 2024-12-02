package cuda_old

import (
	"github.com/MathieuMoalic/amumax/src/engine_old/data_old"
	"github.com/MathieuMoalic/amumax/src/engine_old/log_old"
	"github.com/MathieuMoalic/amumax/src/mesh"
)

// Add magneto-elastic coupling field to the effective field.
// see magnetoelasticfield.cu
func AddMagnetoelasticField(Beff, m *data_old.Slice, exx, eyy, ezz, exy, exz, eyz, B1, B2, Msat MSlice) {
	log_old.AssertMsg(Beff.Size() == m.Size(), "Size mismatch: Beff and m must have the same dimensions in AddMagnetoelasticField")
	log_old.AssertMsg(Beff.Size() == exx.Size(), "Size mismatch: Beff and exx must have the same dimensions in AddMagnetoelasticField")
	log_old.AssertMsg(Beff.Size() == eyy.Size(), "Size mismatch: Beff and eyy must have the same dimensions in AddMagnetoelasticField")
	log_old.AssertMsg(Beff.Size() == ezz.Size(), "Size mismatch: Beff and ezz must have the same dimensions in AddMagnetoelasticField")
	log_old.AssertMsg(Beff.Size() == exy.Size(), "Size mismatch: Beff and exy must have the same dimensions in AddMagnetoelasticField")
	log_old.AssertMsg(Beff.Size() == exz.Size(), "Size mismatch: Beff and exz must have the same dimensions in AddMagnetoelasticField")
	log_old.AssertMsg(Beff.Size() == eyz.Size(), "Size mismatch: Beff and eyz must have the same dimensions in AddMagnetoelasticField")

	N := Beff.Len()
	cfg := make1DConf(N)

	k_addmagnetoelasticfield_async(Beff.DevPtr(X), Beff.DevPtr(Y), Beff.DevPtr(Z),
		m.DevPtr(X), m.DevPtr(Y), m.DevPtr(Z),
		exx.DevPtr(0), exx.Mul(0), eyy.DevPtr(0), eyy.Mul(0), ezz.DevPtr(0), ezz.Mul(0),
		exy.DevPtr(0), exy.Mul(0), exz.DevPtr(0), exz.Mul(0), eyz.DevPtr(0), eyz.Mul(0),
		B1.DevPtr(0), B1.Mul(0), B2.DevPtr(0), B2.Mul(0),
		Msat.DevPtr(0), Msat.Mul(0),
		N, cfg)
}

// Calculate magneto-elastic force density
// see magnetoelasticforce.cu
func GetMagnetoelasticForceDensity(out, m *data_old.Slice, B1, B2 MSlice, mesh mesh.MeshLike) {
	log_old.AssertMsg(out.Size() == m.Size(), "Size mismatch: out and m must have the same dimensions in GetMagnetoelasticForceDensity")

	cellsize := mesh.CellSize()
	N := mesh.Size()
	cfg := make3DConf(N)

	rcsx := float32(1.0 / cellsize[X])
	rcsy := float32(1.0 / cellsize[Y])
	rcsz := float32(1.0 / cellsize[Z])

	k_getmagnetoelasticforce_async(out.DevPtr(X), out.DevPtr(Y), out.DevPtr(Z),
		m.DevPtr(X), m.DevPtr(Y), m.DevPtr(Z),
		B1.DevPtr(0), B1.Mul(0), B2.DevPtr(0), B2.Mul(0),
		rcsx, rcsy, rcsz,
		N[X], N[Y], N[Z],
		mesh.PBC_code(), cfg)
}
