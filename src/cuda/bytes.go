package cuda

// This file provides GPU byte slices, used to store regions.

import (
	"fmt"
	"unsafe"

	"github.com/MathieuMoalic/amumax/src/cuda/cu"
	"github.com/MathieuMoalic/amumax/src/engine_old/log_old"
)

// 3D byte slice, used for region lookup.
type Bytes struct {
	Ptr unsafe.Pointer
	Len int
}

// Construct new byte slice with given length,
// initialised to zeros.
func NewBytes(Len int) *Bytes {
	ptr := cu.MemAlloc(int64(Len))
	cu.MemsetD8(cu.DevicePtr(ptr), 0, int64(Len))
	return &Bytes{unsafe.Pointer(uintptr(ptr)), Len}
}

// Upload src (host) to dst (gpu).
func (dst *Bytes) Upload(src []byte) {
	log_old.AssertMsg(dst.Len == len(src), "Upload: Length mismatch between destination (gpu) and source (host) data")
	MemCpyHtoD(dst.Ptr, unsafe.Pointer(&src[0]), int64(dst.Len))
}

// Copy on device: dst = src.
func (dst *Bytes) Copy(src *Bytes) {
	log_old.AssertMsg(dst.Len == src.Len, "Copy: Length mismatch between source and destination data on device")
	MemCpy(dst.Ptr, src.Ptr, int64(dst.Len))
}

// Copy to host: dst = src.
func (src *Bytes) Download(dst []byte) {
	log_old.AssertMsg(src.Len == len(dst), "Download: Length mismatch between source (gpu) and destination (host) data")
	MemCpyDtoH(unsafe.Pointer(&dst[0]), src.Ptr, int64(src.Len))
}

// Set one element to value.
// data.Index can be used to find the index for x,y,z.
func (dst *Bytes) Set(index int, value byte) {
	if index < 0 || index >= dst.Len {
		log_old.Log.PanicIfError(fmt.Errorf("Bytes.Set: index out of range: %d", index))
	}
	src := value
	MemCpyHtoD(unsafe.Pointer(uintptr(dst.Ptr)+uintptr(index)), unsafe.Pointer(&src), 1)
}

// Get one element.
// data.Index can be used to find the index for x,y,z.
func (src *Bytes) Get(index int) byte {
	if index < 0 || index >= src.Len {
		log_old.Log.PanicIfError(fmt.Errorf("Bytes.Set: index out of range: %v", index))
	}
	var dst byte
	MemCpyDtoH(unsafe.Pointer(&dst), unsafe.Pointer(uintptr(src.Ptr)+uintptr(index)), 1)
	return dst
}

// Frees the GPU memory and disables the slice.
func (b *Bytes) Free() {
	if b.Ptr != nil {
		cu.MemFree(cu.DevicePtr(uintptr(b.Ptr)))
	}
	b.Ptr = nil
	b.Len = 0
}
