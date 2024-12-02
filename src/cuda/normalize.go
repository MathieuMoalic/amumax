package cuda

import (
	"github.com/MathieuMoalic/amumax/src/log"
	"github.com/MathieuMoalic/amumax/src/slice"
)

// Normalize vec to unit length, unless length or vol are zero.
func Normalize(vec, vol *slice.Slice) {
	log.AssertMsg(vol == nil || vol.NComp() == 1, "Invalid volume component: vol must have 1 component or be nil in Normalize")

	N := vec.Len()
	cfg := make1DConf(N)
	k_normalize_async(vec.DevPtr(X), vec.DevPtr(Y), vec.DevPtr(Z), vol.DevPtr(0), N, cfg)
}
