package engine

import (
	"github.com/MathieuMoalic/amumax/cuda"
	"github.com/MathieuMoalic/amumax/data"
)

func SetPhi(dst *data.Slice) {
	cuda.SetPhi(dst, M.Buffer())
}

func SetTheta(dst *data.Slice) {
	cuda.SetTheta(dst, M.Buffer())
}
