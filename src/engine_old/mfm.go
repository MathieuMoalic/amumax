package engine_old

import (
	"github.com/MathieuMoalic/amumax/src/engine_old/cuda_old"
	"github.com/MathieuMoalic/amumax/src/engine_old/data_old"
)

var (
	Mfm        = newScalarField("MFM", "arb.", "MFM image", setMFM)
	mfmLift    inputValue
	mfmTipSize inputValue
	mfmconv_   *cuda_old.MFMConvolution
)

func init() {
	mfmLift = numParam(50e-9, "MFMLift", "m", reinitmfmconv)
	mfmTipSize = numParam(1e-3, "MFMDipole", "m", reinitmfmconv)
}

func setMFM(dst *data_old.Slice) {
	buf := cuda_old.Buffer(3, GetMesh().Size())
	defer cuda_old.Recycle(buf)
	if mfmconv_ == nil {
		reinitmfmconv()
	}

	msat := Msat.MSlice()
	defer msat.Recycle()

	mfmconv_.Exec(buf, NormMag.Buffer(), Geometry.Gpu(), msat)
	cuda_old.Madd3(dst, buf.Comp(0), buf.Comp(1), buf.Comp(2), 1, 1, 1)
}

func reinitmfmconv() {
	setBusy(true)
	defer setBusy(false)
	if mfmconv_ == nil {
		mfmconv_ = cuda_old.NewMFM(GetMesh(), mfmLift.v, mfmTipSize.v, CacheDir)
	} else {
		mfmconv_.Reinit(mfmLift.v, mfmTipSize.v, CacheDir)
	}
}
