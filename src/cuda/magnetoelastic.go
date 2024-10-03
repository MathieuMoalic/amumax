package cuda

import (
	"github.com/MathieuMoalic/amumax/src/data"
	"github.com/MathieuMoalic/amumax/src/log"
)

// Add magneto-elasticit coupling field to the effective field.
// see magnetoelasticfield.cu
func AddMagnetoelasticField(Beff, m *data.Slice, exx, eyy, ezz, exy, exz, eyz, B1, B2, Msat MSlice) {
	log.AssertArgument(Beff.Size() == m.Size())
	log.AssertArgument(Beff.Size() == exx.Size())
	log.AssertArgument(Beff.Size() == eyy.Size())
	log.AssertArgument(Beff.Size() == ezz.Size())
	log.AssertArgument(Beff.Size() == exy.Size())
	log.AssertArgument(Beff.Size() == exz.Size())
	log.AssertArgument(Beff.Size() == eyz.Size())

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

// Calculate magneto-elasticit force density
// see magnetoelasticforce.cu
func GetMagnetoelasticForceDensity(out, m *data.Slice, B1, B2 MSlice, mesh *data.Mesh) {
	log.AssertArgument(out.Size() == m.Size())

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
