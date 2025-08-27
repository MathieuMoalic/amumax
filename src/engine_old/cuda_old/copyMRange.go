package cuda_old

import (
	"github.com/MathieuMoalic/amumax/src/engine_old/data_old"
)

func CopyMRange(dst, src *data_old.Slice, dst0, src0, box [3]int, wrap bool) {
	N := dst.Size()
	W, H, D := box[X], box[Y], box[Z]
	cfg := make3DConf([3]int{W, H, D})
	w := 0
	if wrap {
		w = 1
	}
	k_copyMRange_async(
		dst.DevPtr(0), src.DevPtr(0),
		N[X], N[Y], N[Z],
		dst0[X], dst0[Y], dst0[Z],
		src0[X], src0[Y], src0[Z],
		W, H, D, w, cfg,
	)
}
