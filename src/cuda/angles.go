package cuda

import (
	"github.com/MathieuMoalic/amumax/src/log"
	"github.com/MathieuMoalic/amumax/src/slice"
)

func SetPhi(s *slice.Slice, m *slice.Slice) {
	N := s.Size()
	log.AssertMsg(m.Size() == N, "SetPhi: Size mismatch between slices s and m")
	cfg := make3DConf(N)
	k_setPhi_async(s.DevPtr(X), m.DevPtr(X), m.DevPtr(Y), N[X], N[Y], N[Z], cfg)
}

func SetTheta(s *slice.Slice, m *slice.Slice) {
	N := s.Size()
	log.AssertMsg(m.Size() == N, "SetTheta: Size mismatch between slices s and m")
	cfg := make3DConf(N)
	k_setTheta_async(s.DevPtr(X), m.DevPtr(Z), N[X], N[Y], N[Z], cfg)
}
