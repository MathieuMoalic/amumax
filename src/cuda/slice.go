package cuda

import (
	"math"
	"unsafe"

	"github.com/MathieuMoalic/amumax/src/cuda/cu"
	"github.com/MathieuMoalic/amumax/src/engine_old/timer_old"
	"github.com/MathieuMoalic/amumax/src/log"
	"github.com/MathieuMoalic/amumax/src/slice"
)

// Make a GPU Slice with nComp components each of size length.
func NewSlice(nComp int, size [3]int) *slice.Slice {
	return newSlice(nComp, size, MemAlloc, slice.GPUMemory)
}

// Make a GPU Slice with nComp components each of size length.
//func NewUnifiedSlice(nComp int, m *mesh.Mesh) *data.Slice {
//	return newSlice(nComp, m, cu.MemAllocHost, data.UnifiedMemory)
//}

func newSlice(nComp int, size [3]int, alloc func(int64) unsafe.Pointer, memType int8) *slice.Slice {
	slice.EnableGPU(memFree, cu.MemFreeHost, MemCpy, MemCpyDtoH, MemCpyHtoD)
	length := prod(size)
	bytes := int64(length) * cu.SIZEOF_FLOAT32
	ptrs := make([]unsafe.Pointer, nComp)
	for c := range ptrs {
		ptrs[c] = unsafe.Pointer(alloc(bytes))
		cu.MemsetD32(cu.DevicePtr(uintptr(ptrs[c])), 0, int64(length))
	}
	return slice.SliceFromPtrs(size, memType, ptrs)
}

// wrappers for data.EnableGPU arguments

func memFree(ptr unsafe.Pointer) { cu.MemFree(cu.DevicePtr(uintptr(ptr))) }

func MemCpyDtoH(dst, src unsafe.Pointer, bytes int64) {
	Sync() // sync previous kernels
	timer_old.Start("memcpyDtoH")
	cu.MemcpyDtoH(dst, cu.DevicePtr(uintptr(src)), bytes)
	Sync() // sync copy
	timer_old.Stop("memcpyDtoH")
}

func MemCpyHtoD(dst, src unsafe.Pointer, bytes int64) {
	Sync() // sync previous kernels
	timer_old.Start("memcpyHtoD")
	cu.MemcpyHtoD(cu.DevicePtr(uintptr(dst)), src, bytes)
	Sync() // sync copy
	timer_old.Stop("memcpyHtoD")
}

func MemCpy(dst, src unsafe.Pointer, bytes int64) {
	Sync()
	timer_old.Start("memcpy")
	cu.MemcpyAsync(cu.DevicePtr(uintptr(dst)), cu.DevicePtr(uintptr(src)), bytes, stream0)
	Sync()
	timer_old.Stop("memcpy")
}

// Memset sets the Slice's components to the specified values.
// To be carefully used on unified slice (need sync)
func Memset(s *slice.Slice, val ...float32) {
	if Synchronous { // debug
		Sync()
		timer_old.Start("memset")
	}
	log.AssertMsg(len(val) == s.NComp(), "Memset: wrong number of values")
	for c, v := range val {
		cu.MemsetD32Async(cu.DevicePtr(uintptr(s.DevPtr(c))), math.Float32bits(v), int64(s.Len()), stream0)
	}
	if Synchronous { //debug
		Sync()
		timer_old.Stop("memset")
	}
}

// Set all elements of all components to zero.
func Zero(s *slice.Slice) {
	Memset(s, make([]float32, s.NComp())...)
}

func SetCell(s *slice.Slice, comp int, ix, iy, iz int, value float32) {
	SetElem(s, comp, s.Index(ix, iy, iz), value)
}

func SetElem(s *slice.Slice, comp int, index int, value float32) {
	f := value
	dst := unsafe.Pointer(uintptr(s.DevPtr(comp)) + uintptr(index)*cu.SIZEOF_FLOAT32)
	MemCpyHtoD(dst, unsafe.Pointer(&f), cu.SIZEOF_FLOAT32)
}

func GetElem(s *slice.Slice, comp int, index int) float32 {
	var f float32
	src := unsafe.Pointer(uintptr(s.DevPtr(comp)) + uintptr(index)*cu.SIZEOF_FLOAT32)
	MemCpyDtoH(unsafe.Pointer(&f), src, cu.SIZEOF_FLOAT32)
	return f
}

func GetCell(s *slice.Slice, comp, ix, iy, iz int) float32 {
	return GetElem(s, comp, s.Index(ix, iy, iz))
}
