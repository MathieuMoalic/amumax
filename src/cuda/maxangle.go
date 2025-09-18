package cuda

import (
	"unsafe"

	"github.com/MathieuMoalic/amumax/src/data"
	"github.com/MathieuMoalic/amumax/src/mesh"
)

// SetMaxAngle sets dst to the maximum angle of each cells magnetization with all of its neighbors,
// provided the exchange stiffness with that neighbor is nonzero.
func SetMaxAngle(dst, m *data.Slice, AexRed SymmLUT, regions *Bytes, mesh mesh.MeshLike) {
	N := mesh.Size()
	pbc := mesh.PBCCode()
	cfg := make3DConf(N)
	kSetmaxangleAsync(dst.DevPtr(0),
		m.DevPtr(X), m.DevPtr(Y), m.DevPtr(Z),
		unsafe.Pointer(AexRed), regions.Ptr,
		N[X], N[Y], N[Z], pbc, cfg)
}
