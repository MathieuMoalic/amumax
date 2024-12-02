package cuda

import (
	"github.com/MathieuMoalic/amumax/src/log"
	"github.com/MathieuMoalic/amumax/src/slice"
)

// Copies src (larger) into dst (smaller).
// Used to extract demag field after convolution on padded m.
func copyUnPad(dst, src *slice.Slice, dstsize, srcsize [3]int) {
	log.AssertMsg(dst.NComp() == 1 && src.NComp() == 1, "copyUnPad: Both destination and source must have a single component")
	log.AssertMsg(dst.Len() == prod(dstsize) && src.Len() == prod(srcsize), "copyUnPad: Length mismatch between destination and source sizes")

	cfg := make3DConf(dstsize)

	k_copyunpad_async(dst.DevPtr(0), dstsize[X], dstsize[Y], dstsize[Z],
		src.DevPtr(0), srcsize[X], srcsize[Y], srcsize[Z], cfg)
}

// Copies src into dst, which is larger, and multiplies by vol*Bsat.
// The remainder of dst is not filled with zeros.
// Used to zero-pad magnetization before convolution and in the meanwhile multiply m by its length.
func copyPadMul(dst, src, vol *slice.Slice, dstsize, srcsize [3]int, Msat MSlice) {
	log.AssertMsg(dst.NComp() == 1 && src.NComp() == 1, "copyPadMul: Both destination and source must have a single component")
	log.AssertMsg(dst.Len() == prod(dstsize) && src.Len() == prod(srcsize), "copyPadMul: Length mismatch between destination and source sizes")

	cfg := make3DConf(srcsize)

	k_copypadmul2_async(dst.DevPtr(0), dstsize[X], dstsize[Y], dstsize[Z],
		src.DevPtr(0), srcsize[X], srcsize[Y], srcsize[Z],
		Msat.DevPtr(0), Msat.Mul(0), vol.DevPtr(0), cfg)
}
