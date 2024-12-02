package slice

// Slice stores N-component GPU or host data.

import (
	"fmt"
	"unsafe"

	"github.com/MathieuMoalic/amumax/src/utils"
	"github.com/MathieuMoalic/amumax/src/vector"
)

// Slice is like a [][]float32, but may be stored in GPU or host memory.
type Slice struct {
	ptrs    []unsafe.Pointer
	size    [3]int
	memType int8
}

// this package must not depend on CUDA. If CUDA is
// loaded, these functions are set to cu.MemFree, ...
// NOTE: cpyDtoH and cpuHtoD are only needed to support 32-bit builds,
// otherwise, it could be removed in favor of memCpy only.
var (
	memFree                        func(unsafe.Pointer)
	memCpy, memCpyDtoH, memCpyHtoD func(dst, src unsafe.Pointer, bytes int64)
)

// value for Slice.memType
const (
	CPUMemory      = 1 << 0
	GPUMemory      = 1 << 1
	SIZEOF_FLOAT32 = 4
)

// Internal: enables slices on GPU. Called upon cuda init.
func EnableGPU(free, freeHost func(unsafe.Pointer),
	cpy, cpyDtoH, cpyHtoD func(dst, src unsafe.Pointer, bytes int64)) {
	memFree = free
	memCpy = cpy
	memCpyDtoH = cpyDtoH
	memCpyHtoD = cpyHtoD
}

// Make a CPU Slice with nComp components of size length.
func NewSlice(nComp int, size [3]int) (*Slice, error) {
	length := utils.Prod(size)
	if nComp <= 0 || length <= 0 {
		return nil, fmt.Errorf("invalid input: number of components and size must be greater than 0")
	}
	ptrs := make([]unsafe.Pointer, nComp)
	for i := range ptrs {
		ptrs[i] = unsafe.Pointer(&(make([]float32, length)[0]))
	}
	return SliceFromPtrs(size, CPUMemory, ptrs)
}

func SliceFromArray(data [][]float32, size [3]int) (*Slice, error) {
	nComp := len(data)
	length := utils.Prod(size)
	ptrs := make([]unsafe.Pointer, nComp)
	for i := range ptrs {
		if len(data[i]) != length {
			return nil, fmt.Errorf("invalid input: data[%v] has length %v, expected %v", i, len(data[i]), length)
		}
		ptrs[i] = unsafe.Pointer(&data[i][0])
	}
	return SliceFromPtrs(size, CPUMemory, ptrs)
}

// Return a slice without underlying storage. Used to represent a mask containing all 1's.
func NilSlice(nComp int, size [3]int) (*Slice, error) {
	return SliceFromPtrs(size, GPUMemory, make([]unsafe.Pointer, nComp))
}

// Internal: construct a Slice using bare memory pointers.
func SliceFromPtrs(size [3]int, memType int8, ptrs []unsafe.Pointer) (*Slice, error) {
	length := utils.Prod(size)
	nComp := len(ptrs)
	if nComp <= 0 || length <= 0 {
		return nil, fmt.Errorf("invalid input: number of components and size must be greater than 0 in SliceFromPtrs")
	}

	s := new(Slice)
	s.ptrs = make([]unsafe.Pointer, nComp)
	s.size = size
	copy(s.ptrs, ptrs)
	s.memType = memType
	return s, nil
}

// Frees the underlying storage and zeros the Slice header to avoid accidental use.
// Slices sharing storage will be invalid after Free. Double free is OK.
func (s *Slice) Free() error {
	if s == nil {
		return nil
	}
	// free storage
	switch s.memType {
	case 0:
		return nil // already freed
	case GPUMemory:
		for _, ptr := range s.ptrs {
			memFree(ptr)
		}
	//case UnifiedMemory:
	//	for _, ptr := range s.ptrs {
	//		memFreeHost(ptr)
	//	}
	case CPUMemory:
		// nothing to do
	default:
		return fmt.Errorf("invalid memory type")
	}
	s.Disable()
	return nil
}

// INTERNAL. Overwrite struct fields with zeros to avoid
// accidental use after Free.
func (s *Slice) Disable() {
	s.ptrs = s.ptrs[:0]
	s.size = [3]int{0, 0, 0}
	s.memType = 0
}

// MemType returns the memory type of the underlying storage:
// CPUMemory, GPUMemory or UnifiedMemory
func (s *Slice) MemType() int {
	return int(s.memType)
}

// GPUAccess returns whether the Slice is accessible by the GPU.
// true means it is either stored on GPU or in unified host memory.
func (s *Slice) GPUAccess() bool {
	return s.memType&GPUMemory != 0
}

// CPUAccess returns whether the Slice is accessible by the CPU.
// true means it is stored in host memory.
func (s *Slice) CPUAccess() bool {
	return s.memType&CPUMemory != 0
}

// NComp returns the number of components.
func (s *Slice) NComp() int {
	return len(s.ptrs)
}

// Len returns the number of elements per component.
func (s *Slice) Len() int {
	return utils.Prod(s.size)
}

func (s *Slice) Size() [3]int {
	if s == nil {
		return [3]int{0, 0, 0}
	}
	return s.size
}

// Comp returns a single component of the Slice.
func (s *Slice) Comp(i int) (*Slice, error) {
	if i < 0 || i >= len(s.ptrs) {
		return nil, fmt.Errorf("component index out of bounds")
	}
	sl := new(Slice)
	sl.ptrs = make([]unsafe.Pointer, 1)
	sl.ptrs[0] = s.ptrs[i]
	sl.size = s.size
	sl.memType = s.memType
	return sl, nil
}

// DevPtr returns a CUDA device pointer to a component.
// Slice must have GPUAccess.
// It is safe to call on a nil slice, returns NULL.
func (s *Slice) DevPtr(component int) (unsafe.Pointer, error) {
	if s == nil {
		return nil, nil
	}
	if !s.GPUAccess() {
		return nil, fmt.Errorf("slice not accessible by GPU")
	}
	if component < 0 || component >= len(s.ptrs) {
		return nil, fmt.Errorf("component index out of bounds")
	}
	return s.ptrs[component], nil
}

// Host returns the Slice as a [][]float32 indexed by component, cell number.
// It should have CPUAccess() == true.
func (s *Slice) Host() ([][]float32, error) {
	if !s.CPUAccess() {
		return nil, fmt.Errorf("slice not accessible by CPU")
	}
	list := make([][]float32, s.NComp())
	for c := range list {
		list[c] = unsafe.Slice((*float32)(unsafe.Pointer(s.ptrs[c])), s.Len())
	}
	return list, nil
}

// Returns a copy of the Slice, allocated on CPU.
func (s *Slice) HostCopy() (*Slice, error) {
	if s == nil {
		return nil, fmt.Errorf("nil slice")
	}
	cpy, err := NewSlice(s.NComp(), s.Size())
	if err != nil {
		return nil, fmt.Errorf("error creating slice copy: %v", err)
	}
	err = s.CopyTo(cpy)
	if err != nil {
		return nil, err
	}
	return cpy, nil
}

func (s *Slice) CopyTo(dst *Slice) error {
	if dst.NComp() != s.NComp() || dst.Len() != s.Len() {
		return fmt.Errorf("slice copy: illegal sizes: dst: %vx%v, src: %vx%v", dst.NComp(), dst.Len(), s.NComp(), s.Len())
	}
	dstIsGpu, srcIsGpu := dst.GPUAccess(), s.GPUAccess()
	bytes := SIZEOF_FLOAT32 * int64(dst.Len())
	switch {
	default:
		return fmt.Errorf("unexpected case in Copy()")
	case dstIsGpu && srcIsGpu:
		for c := 0; c < dst.NComp(); c++ {
			dstPtr, err := dst.DevPtr(c)
			if err != nil {
				return err
			}
			srcPtr, err := s.DevPtr(c)
			if err != nil {
				return err
			}
			memCpy(dstPtr, srcPtr, bytes)
		}
	case srcIsGpu && !dstIsGpu:
		for c := 0; c < dst.NComp(); c++ {
			srcPtr, err := s.DevPtr(c)
			if err != nil {
				return err
			}
			memCpyDtoH(dst.ptrs[c], srcPtr, bytes)
		}
	case !srcIsGpu && dstIsGpu:
		for c := 0; c < dst.NComp(); c++ {
			dstPtr, err := dst.DevPtr(c)
			if err != nil {
				return err
			}
			memCpyHtoD(dstPtr, s.ptrs[c], bytes)
		}
	case !dstIsGpu && !srcIsGpu:
		dstHost, err := dst.Host()
		if err != nil {
			return err
		}
		srcHost, err := s.Host()
		if err != nil {
			return err
		}
		for c := range dstHost {
			copy(dstHost[c], srcHost[c])
		}
	}
	return nil
}

// Floats returns the data as 3D array,
// indexed by cell position. Data should be
// scalar (1 component) and have CPUAccess() == true.
func (s *Slice) Scalars() ([][][]float32, error) {
	x, err := s.Tensors()
	if err != nil {
		return nil, err
	}
	if len(x) != 1 {
		return nil, fmt.Errorf("expecting 1 component, got %v", s.NComp())
	}
	return x[0], nil
}

// Vectors returns the data as 4D array,
// indexed by component, cell position. Data should have
// 3 components and have CPUAccess() == true.
func (s *Slice) Vectors() ([3][][][]float32, error) {
	x, err := s.Tensors()
	if err != nil {
		return [3][][][]float32{}, err
	}
	if len(x) != 3 {
		return [3][][][]float32{}, fmt.Errorf("expecting 3 components, got %v", s.NComp())
	}
	return [3][][][]float32{x[0], x[1], x[2]}, nil
}

// Tensors returns the data as 4D array,
// indexed by component, cell position.
// Requires CPUAccess() == true.
func (s *Slice) Tensors() ([][][][]float32, error) {
	tensors := make([][][][]float32, s.NComp())
	host, err := s.Host()
	if err != nil {
		return nil, err
	}
	for i := range tensors {
		tensors[i] = utils.Reshape(host[i], s.Size())
	}
	return tensors, nil
}

// IsNil returns true if either s is nil or s.pointer[0] == nil
func (s *Slice) IsNil() bool {
	if s == nil {
		return true
	}
	return s.ptrs[0] == nil
}

func (s *Slice) Set(comp, ix, iy, iz int, value float64) error {
	if err := s.checkComp(comp); err != nil {
		return err
	}
	host, err := s.Host()
	if err != nil {
		return err
	}
	index, err := s.Index(ix, iy, iz)
	if err != nil {
		return err
	}
	host[comp][index] = float32(value)
	return nil
}

func (s *Slice) SetVector(ix, iy, iz int, v vector.Vector) error {
	index, err := s.Index(ix, iy, iz)
	if err != nil {
		return err
	}
	host, err := s.Host()
	if err != nil {
		return err
	}
	for c := range v {
		host[c][index] = float32(v[c])
	}
	return nil
}

func (s *Slice) SetScalar(ix, iy, iz int, v float64) error {
	host, err := s.Host()
	if err != nil {
		return err
	}
	index, err := s.Index(ix, iy, iz)
	if err != nil {
		return err
	}
	host[0][index] = float32(v)
	return nil
}

func (s *Slice) Get(comp, ix, iy, iz int) (float64, error) {
	if err := s.checkComp(comp); err != nil {
		return 0, err
	}
	host, err := s.Host()
	if err != nil {
		return 0, err
	}
	index, err := s.Index(ix, iy, iz)
	if err != nil {
		return 0, err
	}
	return float64(host[comp][index]), nil
}

func (s *Slice) checkComp(comp int) error {
	if comp < 0 || comp >= s.NComp() {
		return fmt.Errorf("slice: invalid component index: %v (number of components=%v)", comp, s.NComp())
	}
	return nil
}

func (s *Slice) Index(ix, iy, iz int) (int, error) {
	return Index(s.Size(), ix, iy, iz)
}

func Index(size [3]int, ix, iy, iz int) (int, error) {
	if ix < 0 || ix >= size[0] || iy < 0 || iy >= size[1] || iz < 0 || iz >= size[2] {
		return 0, fmt.Errorf("Slice index out of bounds: %v,%v,%v (bounds=%v)", ix, iy, iz, size)
	}
	return (iz*size[1]+iy)*size[0] + ix, nil
}

// Resample returns a slice of new size N,
// using nearest neighbor interpolation over the input slice.
func (s *Slice) Resample(N [3]int) (*Slice, error) {
	if s.Size() == N {
		return s, nil // nothing to do
	}
	In, err := s.Tensors()
	if err != nil {
		return nil, err
	}
	out, err := NewSlice(s.NComp(), N)
	if err != nil {
		return nil, fmt.Errorf("error creating slice: %v", err)
	}
	Out, err := out.Tensors()
	if err != nil {
		return nil, err
	}
	size1 := utils.SizeOf(In[0])
	size2 := utils.SizeOf(Out[0])
	for c := range Out {
		for i := range Out[c] {
			i1 := (i * size1[2]) / size2[2]
			for j := range Out[c][i] {
				j1 := (j * size1[1]) / size2[1]
				for k := range Out[c][i][j] {
					k1 := (k * size1[0]) / size2[0]
					Out[c][i][j][k] = In[c][i1][j1][k1]
				}
			}
		}
	}
	return out, nil
}

// Downsample returns a slice of new size N, smaller than in.Size().
// Averaging interpolation over the input slice.
// in is returned untouched if the sizes are equal.
func Downsample(In [][][][]float32, N [3]int) ([][][][]float32, error) {
	if utils.SizeOf(In[0]) == N {
		return In, nil // nothing to do
	}

	nComp := len(In)
	out, err := NewSlice(nComp, N)
	if err != nil {
		return nil, fmt.Errorf("error creating slice: %v", err)
	}
	Out, err := out.Tensors()
	if err != nil {
		return nil, err
	}

	srcsize := utils.SizeOf(In[0])
	dstsize := utils.SizeOf(Out[0])

	Dx := dstsize[0]
	Dy := dstsize[1]
	Dz := dstsize[2]
	Sx := srcsize[0]
	Sy := srcsize[1]
	Sz := srcsize[2]
	scalex := Sx / Dx
	scaley := Sy / Dy
	scalez := Sz / Dz
	if scalex <= 0 || scaley <= 0 || scalez <= 0 {
		return nil, fmt.Errorf("scaling factors must be positive in Downsample")
	}

	for c := range Out {

		for iz := 0; iz < Dz; iz++ {
			for iy := 0; iy < Dy; iy++ {
				for ix := 0; ix < Dx; ix++ {
					sum, n := 0.0, 0.0

					for I := 0; I < scalez; I++ {
						i2 := iz*scalez + I
						for J := 0; J < scaley; J++ {
							j2 := iy*scaley + J
							for K := 0; K < scalex; K++ {
								k2 := ix*scalex + K

								if i2 < Sz && j2 < Sy && k2 < Sx {
									sum += float64(In[c][i2][j2][k2])
									n++
								}
							}
						}
					}
					Out[c][iz][iy][ix] = float32(sum / n)
				}
			}
		}
	}

	return Out, nil
}
