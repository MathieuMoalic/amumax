package parameter

import (
	"unsafe"

	"github.com/MathieuMoalic/amumax/src/cuda"
	"github.com/MathieuMoalic/amumax/src/cuda/cu"
	"github.com/MathieuMoalic/amumax/src/log"
	"github.com/MathieuMoalic/amumax/src/mesh"
	"github.com/MathieuMoalic/amumax/src/slice"
)

// look-up table for region based parameters
type lut struct {
	log     *log.Logs
	mesh    *mesh.Mesh         // mesh for region based parameters
	gpu_buf cuda.LUTPtrs       // gpu copy of cpu buffer, only transferred when needed
	gpu_ok  bool               // gpu cache up-to date with cpu source?
	cpu_buf [][NREGION]float32 // table data on cpu
	source  updater            // updates cpu slice
}

type updater interface {
	update() // updates cpu lookup table
}

func (p *lut) init(nComp int, source updater) {
	p.gpu_buf = make(cuda.LUTPtrs, nComp)
	p.cpu_buf = make([][NREGION]float32, nComp)
	p.source = source
}

// get an up-to-date version of the lookup-table on CPU
func (p *lut) cpuLUT() [][NREGION]float32 {
	p.source.update()
	return p.cpu_buf
}

// get an up-to-date version of the lookup-table on GPU
func (p *lut) gpuLUT() cuda.LUTPtrs {
	p.source.update()
	if !p.gpu_ok {
		// upload to GPU
		p.assureAlloc()
		cuda.Sync() // sync previous kernels, may still be using gpu lut
		for c := range p.gpu_buf {
			cuda.MemCpyHtoD(p.gpu_buf[c], unsafe.Pointer(&p.cpu_buf[c][0]), cu.SIZEOF_FLOAT32*NREGION)
		}
		p.gpu_ok = true
		cuda.Sync() //sync upload
	}
	return p.gpu_buf
}

// utility for LUT of single-component slice
func (p *lut) GpuLUT1() cuda.LUTPtr {
	p.log.AssertMsg(len(p.gpu_buf) == 1, "Component mismatch: gpu_buf must have exactly 1 component in gpuLUT1")
	return cuda.LUTPtr(p.gpuLUT()[0])
}

// all slice is 0?
func (p *lut) IsZero() bool {
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

func (p *lut) NonZero() bool { return !p.IsZero() }

func (p *lut) assureAlloc() {
	if p.gpu_buf[0] == nil {
		for i := range p.gpu_buf {
			p.gpu_buf[i] = cuda.MemAlloc(NREGION * cu.SIZEOF_FLOAT32)
		}
	}
}

func (b *lut) NComp() int { return len(b.cpu_buf) }

// uncompress the table to a full array with parameter values per cell.
func (p *lut) Slice() (*slice.Slice, bool) {
	b := cuda.Buffer(p.NComp(), p.mesh.Size())
	p.EvalTo(b)
	return b, true
}

// uncompress the table to a full array in the dst Slice with parameter values per cell.
func (p *lut) EvalTo(dst *slice.Slice) {
	gpu := p.gpuLUT()
	for c := 0; c < p.NComp(); c++ {
		cuda.RegionDecode(dst.Comp(c), cuda.LUTPtr(gpu[c]), Regions.Gpu())
	}
}
