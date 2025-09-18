package cuda

import (
	"unsafe"

	"github.com/MathieuMoalic/amumax/src/data"
)

// ZeroMask Sets vector dst to zero where mask != 0.
func ZeroMask(dst *data.Slice, mask LUTPtr, regions *Bytes) {
	N := dst.Len()
	cfg := make1DConf(N)

	for c := 0; c < dst.NComp(); c++ {
		kZeromaskAsync(dst.DevPtr(c), unsafe.Pointer(mask), regions.Ptr, N, cfg)
	}
}
