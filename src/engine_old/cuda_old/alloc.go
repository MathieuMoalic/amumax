package cuda_old

import (
	"fmt"
	"unsafe"

	"github.com/MathieuMoalic/amumax/src/engine_old/cuda_old/cu"
	"github.com/MathieuMoalic/amumax/src/engine_old/data_old"
	"github.com/MathieuMoalic/amumax/src/engine_old/log_old"
)

// Wrapper for cu.MemAlloc, fatal exit on out of memory.
func MemAlloc(bytes int64) unsafe.Pointer {
	defer func() {
		err := recover()
		if err == cu.ERROR_OUT_OF_MEMORY {
			log_old.Log.PanicIfError(fmt.Errorf("out of memory"))
		}
		if err != nil {
			panic(err)
		}
	}()
	return unsafe.Pointer(uintptr(cu.MemAlloc(bytes)))
}

// Returns a copy of in, allocated on GPU.
func GPUCopy(in *data_old.Slice) *data_old.Slice {
	s := data_old.NewSlice(in.NComp(), in.Size())
	data_old.Copy(s, in)
	return s
}
