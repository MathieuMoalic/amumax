package engine

// Mangeto-elastic coupling.

import (
	"github.com/MathieuMoalic/amumax/src/cuda"
	"github.com/MathieuMoalic/amumax/src/data"
	"github.com/MathieuMoalic/amumax/src/log"
)

var (
	B1        = newScalarParam("B1", "J/m3", "First magneto-elastic coupling constant")
	B2        = newScalarParam("B2", "J/m3", "Second magneto-elastic coupling constant")
	exx       = newScalarExcitation("exx", "", "exx component of the strain tensor")
	eyy       = newScalarExcitation("eyy", "", "eyy component of the strain tensor")
	ezz       = newScalarExcitation("ezz", "", "ezz component of the strain tensor")
	exy       = newScalarExcitation("exy", "", "exy component of the strain tensor")
	exz       = newScalarExcitation("exz", "", "exz component of the strain tensor")
	eyz       = newScalarExcitation("eyz", "", "eyz component of the strain tensor")
	B_mel     = newVectorField("B_mel", "T", "Magneto-elastic filed", addMagnetoelasticField)
	F_mel     = newVectorField("F_mel", "N/m3", "Magneto-elastic force density", getMagnetoelasticForceDensity)
	Edens_mel = newScalarField("Edens_mel", "J/m3", "Magneto-elastic energy density", addMagnetoelasticEnergyDensity)
	E_mel     = newScalarValue("E_mel", "J", "Magneto-elastic energy", getMagnetoelasticEnergy)
)

var (
	zeroMel = newScalarParam("_zeroMel", "", "utility zero parameter")
)

func init() {
	registerEnergy(getMagnetoelasticEnergy, addMagnetoelasticEnergyDensity)
}

func addMagnetoelasticField(dst *data.Slice) {
	haveMel := B1.nonZero() || B2.nonZero()
	if !haveMel {
		return
	}

	Exx := exx.MSlice()
	defer Exx.Recycle()

	Eyy := eyy.MSlice()
	defer Eyy.Recycle()

	Ezz := ezz.MSlice()
	defer Ezz.Recycle()

	Exy := exy.MSlice()
	defer Exy.Recycle()

	Exz := exz.MSlice()
	defer Exz.Recycle()

	Eyz := eyz.MSlice()
	defer Eyz.Recycle()

	b1 := B1.MSlice()
	defer b1.Recycle()

	b2 := B2.MSlice()
	defer b2.Recycle()

	ms := Msat.MSlice()
	defer ms.Recycle()

	cuda.AddMagnetoelasticField(dst, normMag.Buffer(),
		Exx, Eyy, Ezz,
		Exy, Exz, Eyz,
		b1, b2, ms)
}

func getMagnetoelasticForceDensity(dst *data.Slice) {
	haveMel := B1.nonZero() || B2.nonZero()
	if !haveMel {
		return
	}

	log.AssertMsg(B1.IsUniform() && B2.IsUniform(), "Magnetoelastic: B1, B2 must be uniform")

	b1 := B1.MSlice()
	defer b1.Recycle()

	b2 := B2.MSlice()
	defer b2.Recycle()

	cuda.GetMagnetoelasticForceDensity(dst, normMag.Buffer(),
		b1, b2, normMag.Mesh())
}

func addMagnetoelasticEnergyDensity(dst *data.Slice) {
	haveMel := B1.nonZero() || B2.nonZero()
	if !haveMel {
		return
	}

	buf := cuda.Buffer(B_mel.NComp(), B_mel.Mesh().Size())
	defer cuda.Recycle(buf)

	// unnormalized magnetization:
	Mf := ValueOf(M_full)
	defer cuda.Recycle(Mf)

	Exx := exx.MSlice()
	defer Exx.Recycle()

	Eyy := eyy.MSlice()
	defer Eyy.Recycle()

	Ezz := ezz.MSlice()
	defer Ezz.Recycle()

	Exy := exy.MSlice()
	defer Exy.Recycle()

	Exz := exz.MSlice()
	defer Exz.Recycle()

	Eyz := eyz.MSlice()
	defer Eyz.Recycle()

	b1 := B1.MSlice()
	defer b1.Recycle()

	b2 := B2.MSlice()
	defer b2.Recycle()

	ms := Msat.MSlice()
	defer ms.Recycle()

	zeromel := zeroMel.MSlice()
	defer zeromel.Recycle()

	// 1st
	cuda.Zero(buf)
	cuda.AddMagnetoelasticField(buf, normMag.Buffer(),
		Exx, Eyy, Ezz,
		Exy, Exz, Eyz,
		b1, zeromel, ms)
	cuda.AddDotProduct(dst, -1./2., buf, Mf)

	// 1nd
	cuda.Zero(buf)
	cuda.AddMagnetoelasticField(buf, normMag.Buffer(),
		Exx, Eyy, Ezz,
		Exy, Exz, Eyz,
		zeromel, b2, ms)
	cuda.AddDotProduct(dst, -1./1., buf, Mf)
}

// Returns magneto-ell energy in joules.
func getMagnetoelasticEnergy() float64 {
	haveMel := B1.nonZero() || B2.nonZero()
	if !haveMel {
		return float64(0.0)
	}

	buf := cuda.Buffer(1, getMesh().Size())
	defer cuda.Recycle(buf)

	cuda.Zero(buf)
	addMagnetoelasticEnergyDensity(buf)
	return cellVolume() * float64(cuda.Sum(buf))
}
