package engine

import (
	"unsafe"

	"github.com/MathieuMoalic/amumax/src/cuda"
	"github.com/MathieuMoalic/amumax/src/cuda/cu"
	"github.com/MathieuMoalic/amumax/src/data"
	"github.com/MathieuMoalic/amumax/src/log"
)

// look-up table for region based parameters
type lut struct {
	gpuBuf cuda.LUTPtrs       // gpu copy of cpu buffer, only transferred when needed
	gpuOk  bool               // gpu cache up-to date with cpu source?
	cpuBuf [][NREGION]float32 // table data on cpu
	source updater            // updates cpu data
}

type updater interface {
	update() // updates cpu lookup table
}

func (l *lut) init(nComp int, source updater) {
	l.gpuBuf = make(cuda.LUTPtrs, nComp)
	l.cpuBuf = make([][NREGION]float32, nComp)
	l.source = source
}

// get an up-to-date version of the lookup-table on CPU
func (l *lut) cpuLUT() [][NREGION]float32 {
	l.source.update()
	return l.cpuBuf
}

// get an up-to-date version of the lookup-table on GPU
func (l *lut) gpuLUT() cuda.LUTPtrs {
	l.source.update()
	if !l.gpuOk {
		// upload to GPU
		l.assureAlloc()
		cuda.Sync() // sync previous kernels, may still be using gpu lut
		for c := range l.gpuBuf {
			cuda.MemCpyHtoD(l.gpuBuf[c], unsafe.Pointer(&l.cpuBuf[c][0]), cu.SIZEOF_FLOAT32*NREGION)
		}
		l.gpuOk = true
		cuda.Sync() // sync upload
	}
	return l.gpuBuf
}

// utility for LUT of single-component data
func (l *lut) gpuLUT1() cuda.LUTPtr {
	log.AssertMsg(len(l.gpuBuf) == 1, "Component mismatch: gpu_buf must have exactly 1 component in gpuLUT1")
	return cuda.LUTPtr(l.gpuLUT()[0])
}

// all data is 0?
func (l *lut) isZero() bool {
	v := l.cpuLUT()
	for c := range v {
		for i := 0; i < NREGION; i++ {
			if v[c][i] != 0 {
				return false
			}
		}
	}
	return true
}

func (l *lut) nonZero() bool { return !l.isZero() }

func (l *lut) assureAlloc() {
	if l.gpuBuf[0] == nil {
		for i := range l.gpuBuf {
			l.gpuBuf[i] = cuda.MemAlloc(NREGION * cu.SIZEOF_FLOAT32)
		}
	}
}

func (l *lut) NComp() int { return len(l.cpuBuf) }

// Slice uncompress the table to a full array with parameter values per cell.
func (l *lut) Slice() (*data.Slice, bool) {
	b := cuda.Buffer(l.NComp(), GetMesh().Size())
	l.EvalTo(b)
	return b, true
}

// EvalTo uncompress the table to a full array in the dst Slice with parameter values per cell.
func (l *lut) EvalTo(dst *data.Slice) {
	gpu := l.gpuLUT()
	for c := 0; c < l.NComp(); c++ {
		cuda.RegionDecode(dst.Comp(c), cuda.LUTPtr(gpu[c]), Regions.Gpu())
	}
}
