package cuda

import (
	"github.com/MathieuMoalic/amumax/data"
	"github.com/MathieuMoalic/amumax/util"
)

func SetPhi(s *data.Slice, m *data.Slice) {
	N := s.Size()
	util.Argument(m.Size() == N)
	cfg := make3DConf(N)
	k_setPhi_async(s.DevPtr(X), m.DevPtr(X), m.DevPtr(Y), N[X], N[Y], N[Z], cfg)
}

func SetTheta(s *data.Slice, m *data.Slice) {
	N := s.Size()
	util.Argument(m.Size() == N)
	cfg := make3DConf(N)
	k_setTheta_async(s.DevPtr(X), m.DevPtr(Z), N[X], N[Y], N[Z], cfg)
}
