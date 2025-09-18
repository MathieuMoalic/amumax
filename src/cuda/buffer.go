package cuda

// Pool of re-usable GPU buffers.
// Synchronization subtlety:
// async kernel launches mean a buffer may already be recycled when still in use.
// That should be fine since the next launch run in the same stream (0), and will
// effectively wait for the previous operation on the buffer.

import (
	"fmt"
	"unsafe"

	"github.com/MathieuMoalic/amumax/src/cuda/cu"
	"github.com/MathieuMoalic/amumax/src/data"
	"github.com/MathieuMoalic/amumax/src/log"
)

var (
	bufPool  = make(map[int][]unsafe.Pointer)    // pool of GPU buffers indexed by size
	bufCheck = make(map[unsafe.Pointer]struct{}) // checks if pointer originates here to avoid unintended recycle
)

const bufMax = 100 // maximum number of buffers to allocate (detect memory leak early)

// Buffer Returns a GPU slice for temporary use. To be returned to the pool with Recycle
func Buffer(nComp int, size [3]int) *data.Slice {
	if Synchronous {
		Sync()
	}

	ptrs := make([]unsafe.Pointer, nComp)

	// re-use as many buffers as possible form our stack
	N := prod(size)
	pool := bufPool[N]
	nFromPool := iMin(nComp, len(pool))
	for i := 0; i < nFromPool; i++ {
		ptrs[i] = pool[len(pool)-i-1]
	}
	bufPool[N] = pool[:len(pool)-nFromPool]

	// allocate as much new memory as needed
	for i := nFromPool; i < nComp; i++ {
		if len(bufCheck) >= bufMax {
			log.Log.PanicIfError(fmt.Errorf("too many buffers in use, possible memory leak"))
		}
		ptrs[i] = MemAlloc(int64(cu.SIZEOF_FLOAT32 * N))
		bufCheck[ptrs[i]] = struct{}{} // mark this pointer as mine
	}

	return data.SliceFromPtrs(size, data.GPUMemory, ptrs)
}

// Recycle Returns a buffer obtained from GetBuffer to the pool.
func Recycle(s *data.Slice) {
	if Synchronous {
		Sync()
	}

	N := s.Len()
	pool := bufPool[N]
	// put each component buffer back on the stack
	for i := 0; i < s.NComp(); i++ {
		ptr := s.DevPtr(i)
		if ptr == unsafe.Pointer(uintptr(0)) {
			continue
		}
		if _, ok := bufCheck[ptr]; !ok {
			log.Log.PanicIfError(fmt.Errorf("recyle: was not obtained with getbuffer"))
		}
		pool = append(pool, ptr)
	}
	s.Disable() // make it unusable, protect against accidental use after recycle
	bufPool[N] = pool
}

// FreeBuffers Frees all buffers. Called after mesh resize.
func FreeBuffers() {
	Sync()
	for _, size := range bufPool {
		for i := range size {
			cu.DevicePtr(uintptr(size[i])).Free()
			size[i] = nil
		}
	}
	bufPool = make(map[int][]unsafe.Pointer)
	bufCheck = make(map[unsafe.Pointer]struct{})
}
