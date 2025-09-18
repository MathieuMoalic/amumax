package cuda

// This file provides GPU byte slices, used to store regions.

import (
	"fmt"
	"unsafe"

	"github.com/MathieuMoalic/amumax/src/cuda/cu"
	"github.com/MathieuMoalic/amumax/src/log"
)

// Bytes 3D byte slice, used for region lookup.
type Bytes struct {
	Ptr unsafe.Pointer
	Len int
}

// NewBytes Construct new byte slice with given length,
// initialised to zeros.
func NewBytes(Len int) *Bytes {
	ptr := cu.MemAlloc(int64(Len))
	cu.MemsetD8(cu.DevicePtr(ptr), 0, int64(Len))
	return &Bytes{unsafe.Pointer(uintptr(ptr)), Len}
}

// Upload src (host) to dst (gpu).
func (b *Bytes) Upload(src []byte) {
	log.AssertMsg(b.Len == len(src), "Upload: Length mismatch between destination (gpu) and source (host) data")
	MemCpyHtoD(b.Ptr, unsafe.Pointer(&src[0]), int64(b.Len))
}

// Copy on device: dst = src.
func (b *Bytes) Copy(src *Bytes) {
	log.AssertMsg(b.Len == src.Len, "Copy: Length mismatch between source and destination data on device")
	MemCpy(b.Ptr, src.Ptr, int64(b.Len))
}

// Download Copy to host: dst = src.
func (b *Bytes) Download(dst []byte) {
	log.AssertMsg(b.Len == len(dst), "Download: Length mismatch between source (gpu) and destination (host) data")
	MemCpyDtoH(unsafe.Pointer(&dst[0]), b.Ptr, int64(b.Len))
}

// Set one element to value.
// data.Index can be used to find the index for x,y,z.
func (b *Bytes) Set(index int, value byte) {
	if index < 0 || index >= b.Len {
		log.Log.PanicIfError(fmt.Errorf("Bytes.Set: index out of range: %d", index))
	}
	src := value
	MemCpyHtoD(unsafe.Pointer(uintptr(b.Ptr)+uintptr(index)), unsafe.Pointer(&src), 1)
}

// Get one element.
// data.Index can be used to find the index for x,y,z.
func (b *Bytes) Get(index int) byte {
	if index < 0 || index >= b.Len {
		log.Log.PanicIfError(fmt.Errorf("Bytes.Set: index out of range: %v", index))
	}
	var dst byte
	MemCpyDtoH(unsafe.Pointer(&dst), unsafe.Pointer(uintptr(b.Ptr)+uintptr(index)), 1)
	return dst
}

// Free Frees the GPU memory and disables the slice.
func (b *Bytes) Free() {
	if b.Ptr != nil {
		cu.MemFree(cu.DevicePtr(uintptr(b.Ptr)))
	}
	b.Ptr = nil
	b.Len = 0
}
