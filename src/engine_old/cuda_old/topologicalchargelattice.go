package cuda_old

import (
	"github.com/MathieuMoalic/amumax/src/engine_old/data_old"
	"github.com/MathieuMoalic/amumax/src/engine_old/log_old"
	"github.com/MathieuMoalic/amumax/src/mesh"
)

// Topological charge according to Berg and Lüscher
func SetTopologicalChargeLattice(s *data_old.Slice, m *data_old.Slice, mesh mesh.MeshLike) {
	cellsize := mesh.CellSize()
	N := s.Size()
	log_old.AssertMsg(m.Size() == N, "Size mismatch: m and s must have the same dimensions in SetTopologicalChargeLattice")
	cfg := make3DConf(N)
	icxcy := float32(1.0 / (cellsize[X] * cellsize[Y]))

	k_settopologicalchargelattice_async(s.DevPtr(X),
		m.DevPtr(X), m.DevPtr(Y), m.DevPtr(Z),
		icxcy, N[X], N[Y], N[Z], mesh.PBC_code(), cfg)
}
