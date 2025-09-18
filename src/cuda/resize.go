package cuda

import (
	"github.com/MathieuMoalic/amumax/src/data"
	"github.com/MathieuMoalic/amumax/src/log"
)

// Select and resize one layer for interactive output
func Resize(dst, src *data.Slice, layer int) {
	dstsize := dst.Size()
	srcsize := src.Size()
	log.AssertMsg(dstsize[Z] == 1, "Destination slice must have a single layer (Z=1) in Resize")
	log.AssertMsg(dst.NComp() == 1 && src.NComp() == 1, "Component mismatch: dst and src must both have 1 component in Resize")

	scalex := srcsize[X] / dstsize[X]
	scaley := srcsize[Y] / dstsize[Y]
	log.AssertMsg(scalex > 0 && scaley > 0, "Scaling factors must be positive in Resize")

	cfg := make3DConf(dstsize)

	kResizeAsync(dst.DevPtr(0), dstsize[X], dstsize[Y], dstsize[Z],
		src.DevPtr(0), srcsize[X], srcsize[Y], srcsize[Z], layer, scalex, scaley, cfg)
}
