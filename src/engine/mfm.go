package engine

import (
	"github.com/MathieuMoalic/amumax/src/cuda"
	"github.com/MathieuMoalic/amumax/src/data"
)

var (
	Mfm        = newScalarField("MFM", "arb.", "MFM image", setMFM)
	mfmLift    inputValue
	mfmTipSize inputValue
	mfmconv    *cuda.MFMConvolution
)

func init() {
	mfmLift = numParam(50e-9, "MFMLift", "m", reinitmfmconv)
	mfmTipSize = numParam(1e-3, "MFMDipole", "m", reinitmfmconv)
}

func setMFM(dst *data.Slice) {
	buf := cuda.Buffer(3, GetMesh().Size())
	defer cuda.Recycle(buf)
	if mfmconv == nil {
		reinitmfmconv()
	}

	msat := Msat.MSlice()
	defer msat.Recycle()

	mfmconv.Exec(buf, NormMag.Buffer(), Geometry.Gpu(), msat)
	cuda.Madd3(dst, buf.Comp(0), buf.Comp(1), buf.Comp(2), 1, 1, 1)
}

func reinitmfmconv() {
	setBusy(true)
	defer setBusy(false)
	if mfmconv == nil {
		mfmconv = cuda.NewMFM(GetMesh(), mfmLift.v, mfmTipSize.v, CacheDir)
	} else {
		mfmconv.Reinit(mfmLift.v, mfmTipSize.v, CacheDir)
	}
}
