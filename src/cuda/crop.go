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
	log.AssertArgument(dst.NComp() == src.NComp())
	log.AssertArgument(D[X]+offX <= S[X] && D[Y]+offY <= S[Y] && D[Z]+offZ <= S[Z])

	cfg := make3DConf(D)

	for c := 0; c < dst.NComp(); c++ {
		k_crop_async(dst.DevPtr(c), D[X], D[Y], D[Z],
			src.DevPtr(c), S[X], S[Y], S[Z],
			offX, offY, offZ, cfg)
	}
}
