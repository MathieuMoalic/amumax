package cuda

import (
	"fmt"
	"unsafe"

	"github.com/MathieuMoalic/amumax/src/cuda/cu"
	"github.com/MathieuMoalic/amumax/src/log"
	"github.com/MathieuMoalic/amumax/src/slice"
)

// Wrapper for cu.MemAlloc, fatal exit on out of memory.
func MemAlloc(bytes int64) unsafe.Pointer {
	defer func() {
		err := recover()
		if err == cu.ERROR_OUT_OF_MEMORY {
			log.PanicIfError(fmt.Errorf("out of memory"))
		}
		if err != nil {
			panic(err)
		}
	}()
	return unsafe.Pointer(uintptr(cu.MemAlloc(bytes)))
}

// Returns a copy of in, allocated on GPU.
func GPUCopy(in *slice.Slice) *slice.Slice {
	s := NewSlice(in.NComp(), in.Size())
	slice.Copy(s, in)
	return s
}
