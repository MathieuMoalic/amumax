package cuda

import (
	"github.com/MathieuMoalic/amumax/src/log"
	"github.com/MathieuMoalic/amumax/src/slice"
)

// dst += prefactor * dot(a, b), as used for energy density
func AddDotProduct(dst *slice.Slice, prefactor float32, a, b *slice.Slice) {
	log.AssertMsg(dst.NComp() == 1 && a.NComp() == 3 && b.NComp() == 3,
		"Component mismatch: dst must have 1 component, and a and b must each have 3 components in AddDotProduct")

	log.AssertMsg(dst.Len() == a.Len() && dst.Len() == b.Len(),
		"Length mismatch: dst, a, and b must have the same length in AddDotProduct")

	N := dst.Len()
	cfg := make1DConf(N)
	k_dotproduct_async(dst.DevPtr(0), prefactor,
		a.DevPtr(X), a.DevPtr(Y), a.DevPtr(Z),
		b.DevPtr(X), b.DevPtr(Y), b.DevPtr(Z),
		N, cfg)
}
