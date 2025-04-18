package cuda_old

import (
	"unsafe"

	"github.com/MathieuMoalic/amumax/src/engine_old/data_old"
	"github.com/MathieuMoalic/amumax/src/mesh"
)

// SetMaxAngle sets dst to the maximum angle of each cells magnetization with all of its neighbors,
// provided the exchange stiffness with that neighbor is nonzero.
func SetMaxAngle(dst, m *data_old.Slice, Aex_red SymmLUT, regions *Bytes, mesh mesh.MeshLike) {
	N := mesh.Size()
	pbc := mesh.PBC_code()
	cfg := make3DConf(N)
	k_setmaxangle_async(dst.DevPtr(0),
		m.DevPtr(X), m.DevPtr(Y), m.DevPtr(Z),
		unsafe.Pointer(Aex_red), regions.Ptr,
		N[X], N[Y], N[Z], pbc, cfg)
}
