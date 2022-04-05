package cuda

import (
	"github.com/MathieuMoalic/amumax/data"
	"github.com/MathieuMoalic/amumax/util"
)

// Topological charge according to Berg and Lüscher
func SetTopologicalChargeLattice(s *data.Slice, m *data.Slice, mesh *data.Mesh) {
	cellsize := mesh.CellSize()
	N := s.Size()
	util.Argument(m.Size() == N)
	cfg := make3DConf(N)
	icxcy := float32(1.0 / (cellsize[X] * cellsize[Y]))

	k_settopologicalchargelattice_async(s.DevPtr(X),
		m.DevPtr(X), m.DevPtr(Y), m.DevPtr(Z),
		icxcy, N[X], N[Y], N[Z], mesh.PBC_code(), cfg)
}
