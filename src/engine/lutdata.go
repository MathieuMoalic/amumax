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
	source  updater            // updates cpu data
}

type updater interface {
	update() // updates cpu lookup table
}

func (p *lut) init(nComp int, source updater) {
	p.gpuBuf = make(cuda.LUTPtrs, nComp)
	p.cpuBuf = make([][NREGION]float32, nComp)
	p.source = source
}

// get an up-to-date version of the lookup-table on CPU
func (p *lut) cpuLUT() [][NREGION]float32 {
	p.source.update()
	return p.cpuBuf
}

// get an up-to-date version of the lookup-table on GPU
func (p *lut) gpuLUT() cuda.LUTPtrs {
	p.source.update()
	if !p.gpuOk {
		// upload to GPU
		p.assureAlloc()
		cuda.Sync() // sync previous kernels, may still be using gpu lut
		for c := range p.gpuBuf {
			cuda.MemCpyHtoD(p.gpuBuf[c], unsafe.Pointer(&p.cpuBuf[c][0]), cu.SIZEOF_FLOAT32*NREGION)
		}
		p.gpuOk = true
		cuda.Sync() // sync upload
	}
	return p.gpuBuf
}

// utility for LUT of single-component data
func (p *lut) gpuLUT1() cuda.LUTPtr {
	log.AssertMsg(len(p.gpuBuf) == 1, "Component mismatch: gpu_buf must have exactly 1 component in gpuLUT1")
	return cuda.LUTPtr(p.gpuLUT()[0])
}

// all data is 0?
func (p *lut) isZero() bool {
	v := p.cpuLUT()
	for c := range v {
		for i := 0; i < NREGION; i++ {
			if v[c][i] != 0 {
				return false
			}
		}
	}
	return true
}

func (p *lut) nonZero() bool { return !p.isZero() }

func (p *lut) assureAlloc() {
	if p.gpuBuf[0] == nil {
		for i := range p.gpuBuf {
			p.gpuBuf[i] = cuda.MemAlloc(NREGION * cu.SIZEOF_FLOAT32)
		}
	}
}

func (b *lut) NComp() int { return len(b.cpuBuf) }

// Slice uncompress the table to a full array with parameter values per cell.
func (p *lut) Slice() (*data.Slice, bool) {
	b := cuda.Buffer(p.NComp(), GetMesh().Size())
	p.EvalTo(b)
	return b, true
}

// EvalTo uncompress the table to a full array in the dst Slice with parameter values per cell.
func (p *lut) EvalTo(dst *data.Slice) {
	gpu := p.gpuLUT()
	for c := 0; c < p.NComp(); c++ {
		cuda.RegionDecode(dst.Comp(c), cuda.LUTPtr(gpu[c]), Regions.Gpu())
	}
}
