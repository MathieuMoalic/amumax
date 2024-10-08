package cuda

import (
	"github.com/MathieuMoalic/amumax/src/data"
	"github.com/MathieuMoalic/amumax/src/log"
)

// Normalize vec to unit length, unless length or vol are zero.
func Normalize(vec, vol *data.Slice) {
	log.AssertArgument(vol == nil || vol.NComp() == 1)
	N := vec.Len()
	cfg := make1DConf(N)
	k_normalize_async(vec.DevPtr(X), vec.DevPtr(Y), vec.DevPtr(Z), vol.DevPtr(0), N, cfg)
}
