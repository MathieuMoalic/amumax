package cuda

import (
	"github.com/MathieuMoalic/amumax/src/data"
	"github.com/MathieuMoalic/amumax/src/log"
)

// Crop stores in dst a rectangle cropped from src at given offset position.
// dst size may be smaller than src.
func Crop(dst, src *data.Slice, offX, offY, offZ int) {
	D := dst.Size()
	S := src.Size()
	log.AssertMsg(dst.NComp() == src.NComp(), "dst and src must have the same number of components in Crop function")

	log.AssertMsg(D[X]+offX <= S[X] && D[Y]+offY <= S[Y] && D[Z]+offZ <= S[Z],
		"Invalid crop parameters: destination size plus offset exceeds source dimensions in Crop function")

	cfg := make3DConf(D)

	for c := 0; c < dst.NComp(); c++ {
		kCropAsync(dst.DevPtr(c), D[X], D[Y], D[Z],
			src.DevPtr(c), S[X], S[Y], S[Z],
			offX, offY, offZ, cfg)
	}
}
