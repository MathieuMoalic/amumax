package cuda

import (
	"github.com/MathieuMoalic/amumax/src/log"
	"github.com/MathieuMoalic/amumax/src/mesh"
	"github.com/MathieuMoalic/amumax/src/slice"
)

// Topological charge according to Berg and Lüscher
func SetTopologicalChargeLattice(s *slice.Slice, m *slice.Slice, mesh mesh.Mesh) {
	cellsize := mesh.CellSize()
	N := s.Size()
	log.AssertMsg(m.Size() == N, "Size mismatch: m and s must have the same dimensions in SetTopologicalChargeLattice")
	cfg := make3DConf(N)
	icxcy := float32(1.0 / (cellsize[X] * cellsize[Y]))

	k_settopologicalchargelattice_async(s.DevPtr(X),
		m.DevPtr(X), m.DevPtr(Y), m.DevPtr(Z),
		icxcy, N[X], N[Y], N[Z], mesh.PBC_code(), cfg)
}
