package cuda_old

import (
	"github.com/MathieuMoalic/amumax/src/cuda_old/cu"
	"github.com/MathieuMoalic/amumax/src/cuda_old/cufft"
	"github.com/MathieuMoalic/amumax/src/engine_old/data_old"
	"github.com/MathieuMoalic/amumax/src/engine_old/log_old"
	"github.com/MathieuMoalic/amumax/src/timer"
)

// 3D single-precission real-to-complex FFT plan.
type fft3DR2CPlan struct {
	fftplan
	size [3]int
}

// 3D single-precission real-to-complex FFT plan.
func newFFT3DR2C(Nx, Ny, Nz int) fft3DR2CPlan {
	handle := cufft.Plan3d(Nz, Ny, Nx, cufft.R2C) // new xyz swap
	handle.SetStream(stream0)
	return fft3DR2CPlan{fftplan{handle}, [3]int{Nx, Ny, Nz}}
}

// Execute the FFT plan, asynchronous.
// src and dst are 3D arrays stored as 1D arrays.
func (p *fft3DR2CPlan) ExecAsync(src, dst *data_old.Slice) {
	if Synchronous {
		Sync()
		timer.Start("fft")
	}
	log_old.AssertMsg(src.NComp() == 1 && dst.NComp() == 1,
		"Component mismatch: both src and dst must have 1 component in fft3DR2CPlan.ExecAsync")

	oksrclen := p.InputLen()
	if src.Len() != oksrclen {
		log_old.Log.ErrAndExit("FFT size mismatch: expected src length %v, but got %v", oksrclen, src.Len())
	}

	okdstlen := p.OutputLen()
	if dst.Len() != okdstlen {
		log_old.Log.ErrAndExit("FFT size mismatch: expected dst length %v, but got %v", okdstlen, dst.Len())
	}

	p.handle.ExecR2C(cu.DevicePtr(uintptr(src.DevPtr(0))), cu.DevicePtr(uintptr(dst.DevPtr(0))))

	if Synchronous {
		Sync()
		timer.Stop("fft")
	}
}

// 3D size of the input array.
func (p *fft3DR2CPlan) InputSizeFloats() (Nx, Ny, Nz int) {
	return p.size[X], p.size[Y], p.size[Z]
}

// 3D size of the output array.
func (p *fft3DR2CPlan) OutputSizeFloats() (Nx, Ny, Nz int) {
	return 2 * (p.size[X]/2 + 1), p.size[Y], p.size[Z]
}

// Required length of the (1D) input array.
func (p *fft3DR2CPlan) InputLen() int {
	return prod3(p.InputSizeFloats())
}

// Required length of the (1D) output array.
func (p *fft3DR2CPlan) OutputLen() int {
	return prod3(p.OutputSizeFloats())
}
