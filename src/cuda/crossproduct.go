package cuda

import (
	"github.com/MathieuMoalic/amumax/src/log"
	"github.com/MathieuMoalic/amumax/src/slice"
)

func CrossProduct(dst, a, b *slice.Slice) {
	log.AssertMsg(dst.NComp() == 3 && a.NComp() == 3 && b.NComp() == 3,
		"Invalid number of components: dst, a, and b must all have 3 components for CrossProduct")

	log.AssertMsg(dst.Len() == a.Len() && dst.Len() == b.Len(),
		"Length mismatch: dst, a, and b must have the same length for CrossProduct")

	N := dst.Len()
	cfg := make1DConf(N)
	k_crossproduct_async(dst.DevPtr(X), dst.DevPtr(Y), dst.DevPtr(Z),
		a.DevPtr(X), a.DevPtr(Y), a.DevPtr(Z),
		b.DevPtr(X), b.DevPtr(Y), b.DevPtr(Z),
		N, cfg)
}
