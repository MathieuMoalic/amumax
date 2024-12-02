package engine_old

import (
	"unsafe"

	"github.com/MathieuMoalic/amumax/src/cuda_old"
	"github.com/MathieuMoalic/amumax/src/cuda_old/cu"
	"github.com/MathieuMoalic/amumax/src/engine_old/data_old"
	"github.com/MathieuMoalic/amumax/src/engine_old/log_old"
)

// look-up table for region based parameters
type lut struct {
	gpu_buf cuda_old.LUTPtrs   // gpu copy of cpu buffer, only transferred when needed
	gpu_ok  bool               // gpu cache up-to date with cpu source?
	cpu_buf [][NREGION]float32 // table data on cpu
	source  updater            // updates cpu data
}

type updater interface {
	update() // updates cpu lookup table
}

func (p *lut) init(nComp int, source updater) {
	p.gpu_buf = make(cuda_old.LUTPtrs, nComp)
	p.cpu_buf = make([][NREGION]float32, nComp)
	p.source = source
}

// get an up-to-date version of the lookup-table on CPU
func (p *lut) cpuLUT() [][NREGION]float32 {
	p.source.update()
	return p.cpu_buf
}

// get an up-to-date version of the lookup-table on GPU
func (p *lut) gpuLUT() cuda_old.LUTPtrs {
	p.source.update()
	if !p.gpu_ok {
		// upload to GPU
		p.assureAlloc()
		cuda_old.Sync() // sync previous kernels, may still be using gpu lut
		for c := range p.gpu_buf {
			cuda_old.MemCpyHtoD(p.gpu_buf[c], unsafe.Pointer(&p.cpu_buf[c][0]), cu.SIZEOF_FLOAT32*NREGION)
		}
		p.gpu_ok = true
		cuda_old.Sync() //sync upload
	}
	return p.gpu_buf
}

// utility for LUT of single-component data
func (p *lut) gpuLUT1() cuda_old.LUTPtr {
	log_old.AssertMsg(len(p.gpu_buf) == 1, "Component mismatch: gpu_buf must have exactly 1 component in gpuLUT1")
	return cuda_old.LUTPtr(p.gpuLUT()[0])
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
	if p.gpu_buf[0] == nil {
		for i := range p.gpu_buf {
			p.gpu_buf[i] = cuda_old.MemAlloc(NREGION * cu.SIZEOF_FLOAT32)
		}
	}
}

func (b *lut) NComp() int { return len(b.cpu_buf) }

// uncompress the table to a full array with parameter values per cell.
func (p *lut) Slice() (*data_old.Slice, bool) {
	b := cuda_old.Buffer(p.NComp(), GetMesh().Size())
	p.EvalTo(b)
	return b, true
}

// uncompress the table to a full array in the dst Slice with parameter values per cell.
func (p *lut) EvalTo(dst *data_old.Slice) {
	gpu := p.gpuLUT()
	for c := 0; c < p.NComp(); c++ {
		cuda_old.RegionDecode(dst.Comp(c), cuda_old.LUTPtr(gpu[c]), Regions.Gpu())
	}
}
