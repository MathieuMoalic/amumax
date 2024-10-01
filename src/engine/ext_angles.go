package engine

import (
	"github.com/MathieuMoalic/amumax/src/cuda"
	"github.com/MathieuMoalic/amumax/src/data"
)

func SetPhi(dst *data.Slice) {
	cuda.SetPhi(dst, M.Buffer())
}

func SetTheta(dst *data.Slice) {
	cuda.SetTheta(dst, M.Buffer())
}
