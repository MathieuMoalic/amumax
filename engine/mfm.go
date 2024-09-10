package engine

import (
	"github.com/MathieuMoalic/amumax/cuda"
	"github.com/MathieuMoalic/amumax/data"
)

var (
	MFM        = NewScalarField("MFM", "arb.", "MFM image", SetMFM)
	MFMLift    inputValue
	MFMTipSize inputValue
	mfmconv_   *cuda.MFMConvolution
)

func init() {
	MFMLift = numParam(50e-9, "MFMLift", "m", reinitmfmconv)
	MFMTipSize = numParam(1e-3, "MFMDipole", "m", reinitmfmconv)
	DeclLValue("MFMLift", &MFMLift, "MFM lift height")
	DeclLValue("MFMDipole", &MFMTipSize, "Height of vertically magnetized part of MFM tip")
}

func SetMFM(dst *data.Slice) {
	buf := cuda.Buffer(3, GetMesh().Size())
	defer cuda.Recycle(buf)
	if mfmconv_ == nil {
		reinitmfmconv()
	}

	msat := Msat.MSlice()
	defer msat.Recycle()

	mfmconv_.Exec(buf, M.Buffer(), Geometry.Gpu(), msat)
	cuda.Madd3(dst, buf.Comp(0), buf.Comp(1), buf.Comp(2), 1, 1, 1)
}

func reinitmfmconv() {
	SetBusy(true)
	defer SetBusy(false)
	if mfmconv_ == nil {
		mfmconv_ = cuda.NewMFM(GetMesh(), MFMLift.v, MFMTipSize.v, CacheDir)
	} else {
		mfmconv_.Reinit(MFMLift.v, MFMTipSize.v, CacheDir)
	}
}
