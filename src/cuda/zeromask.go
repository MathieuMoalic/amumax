package cuda

import (
	"unsafe"

	"github.com/MathieuMoalic/amumax/src/engine_old/data_old"
)

// Sets vector dst to zero where mask != 0.
func ZeroMask(dst *data_old.Slice, mask LUTPtr, regions *Bytes) {
	N := dst.Len()
	cfg := make1DConf(N)

	for c := 0; c < dst.NComp(); c++ {
		k_zeromask_async(dst.DevPtr(c), unsafe.Pointer(mask), regions.Ptr, N, cfg)
	}
}
