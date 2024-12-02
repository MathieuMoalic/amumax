package cuda

// INTERNAL
// Base implementation for all FFT plans.

import (
	"github.com/MathieuMoalic/amumax/src/engine_old/cuda_old/cufft"
)

// Base implementation for all FFT plans.
type fftplan struct {
	handle cufft.Handle
}

func prod3(x, y, z int) int {
	return x * y * z
}

// Releases all resources associated with the FFT plan.
func (p *fftplan) Free() {
	if p.handle != 0 {
		p.handle.Destroy()
		p.handle = 0
	}
}
