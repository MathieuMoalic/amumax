package cuda

// Kernel multiplication for purely real kernel, symmetric around Y axis (apart from first row).
// Launch configs range over all complex elements of fft input. This could be optimized: range only over kernel.

import (
	"github.com/MathieuMoalic/amumax/src/log"
	"github.com/MathieuMoalic/amumax/src/slice"
)

// kernel multiplication for 3D demag convolution, exploiting full kernel symmetry.
func kernMulRSymm3D_async(fftM [3]*slice.Slice, Kxx, Kyy, Kzz, Kyz, Kxz, Kxy *slice.Slice, Nx, Ny, Nz int) {
	log.AssertMsg(fftM[X].NComp() == 1 && Kxx.NComp() == 1, "fftM[X] or Kxx has incorrect number of components in kernMulRSymm3D_async: expected NComp() == 1")

	cfg := make3DConf([3]int{Nx, Ny, Nz})
	k_kernmulRSymm3D_async(fftM[X].DevPtr(0), fftM[Y].DevPtr(0), fftM[Z].DevPtr(0),
		Kxx.DevPtr(0), Kyy.DevPtr(0), Kzz.DevPtr(0), Kyz.DevPtr(0), Kxz.DevPtr(0), Kxy.DevPtr(0),
		Nx, Ny, Nz, cfg)
}

// kernel multiplication for 2D demag convolution on X and Y, exploiting full kernel symmetry.
func kernMulRSymm2Dxy_async(fftMx, fftMy, Kxx, Kyy, Kxy *slice.Slice, Nx, Ny int) {
	log.AssertMsg(fftMy.NComp() == 1 && Kxx.NComp() == 1, "fftMy or Kxx has incorrect number of components in kernMulRSymm2Dxy_async: expected NComp() == 1")

	cfg := make3DConf([3]int{Nx, Ny, 1})
	k_kernmulRSymm2Dxy_async(fftMx.DevPtr(0), fftMy.DevPtr(0),
		Kxx.DevPtr(0), Kyy.DevPtr(0), Kxy.DevPtr(0),
		Nx, Ny, cfg)
}

// kernel multiplication for 2D demag convolution on Z, exploiting full kernel symmetry.
func kernMulRSymm2Dz_async(fftMz, Kzz *slice.Slice, Nx, Ny int) {
	log.AssertMsg(fftMz.NComp() == 1 && Kzz.NComp() == 1, "fftMz or Kzz has incorrect number of components in kernMulRSymm2Dz_async: expected NComp() == 1")

	cfg := make3DConf([3]int{Nx, Ny, 1})
	k_kernmulRSymm2Dz_async(fftMz.DevPtr(0), Kzz.DevPtr(0), Nx, Ny, cfg)
}

// kernel multiplication for general 1D convolution. Does not assume any symmetry.
// Used for MFM images.
func kernMulC_async(fftM, K *slice.Slice, Nx, Ny int) {
	log.AssertMsg(fftM.NComp() == 1 && K.NComp() == 1, "fftM or K has incorrect number of components in kernMulC_async: expected NComp() == 1")

	cfg := make3DConf([3]int{Nx, Ny, 1})
	k_kernmulC_async(fftM.DevPtr(0), K.DevPtr(0), Nx, Ny, cfg)
}
