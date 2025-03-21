package cuda

import (
	"github.com/MathieuMoalic/amumax/src/log"
	"github.com/MathieuMoalic/amumax/src/mesh"
	"github.com/MathieuMoalic/amumax/src/slice"
)

// Set s to the topological charge density s = m · (∂m/∂x ❌ ∂m/∂y)
// See topologicalcharge.cu
func SetTopologicalCharge(s *slice.Slice, m *slice.Slice, mesh mesh.Mesh) {
	cellsize := mesh.CellSize()
	N := s.Size()
	log.AssertMsg(m.Size() == N, "Size mismatch: m and s must have the same dimensions in SetTopologicalCharge")

	cfg := make3DConf(N)
	icxcy := float32(1.0 / (cellsize[X] * cellsize[Y]))

	k_settopologicalcharge_async(s.DevPtr(X),
		m.DevPtr(X), m.DevPtr(Y), m.DevPtr(Z),
		icxcy, N[X], N[Y], N[Z], mesh.PBC_code(), cfg)
}
