package cuda

// Generation of Magnetic Force Microscopy images.

import (
	"github.com/MathieuMoalic/amumax/src/kernels"
	"github.com/MathieuMoalic/amumax/src/mesh"
	"github.com/MathieuMoalic/amumax/src/slice"
)

// Stores the necessary state to perform FFT-accelerated convolution
type MFMConvolution struct {
	size        [3]int          // 3D size of the input/output data
	kernSize    [3]int          // Size of kernel and logical FFT size.
	fftKernSize [3]int          //
	fftRBuf     *slice.Slice    // FFT input buf for FFT, shares storage with fftCBuf.
	fftCBuf     *slice.Slice    // FFT output buf, shares storage with fftRBuf
	gpuFFTKern  [3]*slice.Slice // FFT kernel on device
	fwPlan      fft3DR2CPlan    // Forward FFT (1 component)
	bwPlan      fft3DC2RPlan    // Backward FFT (1 component)
	kern        [3]*slice.Slice // Real-space kernel (host)
	mesh        mesh.Mesh
}

func (c *MFMConvolution) Free() {
	if c == nil {
		return
	}
	c.size = [3]int{}
	c.kernSize = [3]int{}
	c.fftCBuf.Free() // shared with fftRbuf
	c.fftCBuf = nil
	c.fftRBuf = nil

	for j := 0; j < 3; j++ {
		c.gpuFFTKern[j].Free()
		c.gpuFFTKern[j] = nil
		c.kern[j] = nil
	}
	c.fwPlan.Free()
	c.bwPlan.Free()

	cudaCtx.SetCurrent()
}

func (c *MFMConvolution) init() {
	// init FFT plans
	padded := c.kernSize
	c.fwPlan = newFFT3DR2C(padded[X], padded[Y], padded[Z])
	c.bwPlan = newFFT3DC2R(padded[X], padded[Y], padded[Z])

	// init device buffers
	nc := fftR2COutputSizeFloats(c.kernSize)
	c.fftCBuf = NewSlice(1, nc)
	c.fftRBuf = NewSlice(1, c.kernSize)

	c.gpuFFTKern[X] = NewSlice(1, nc)
	c.gpuFFTKern[Y] = NewSlice(1, nc)
	c.gpuFFTKern[Z] = NewSlice(1, nc)

	c.initFFTKern3D()
}

func (c *MFMConvolution) initFFTKern3D() {
	c.fftKernSize = fftR2COutputSizeFloats(c.kernSize)

	for i := 0; i < 3; i++ {
		zero1_async(c.fftRBuf)
		slice.Copy(c.fftRBuf, c.kern[i])
		c.fwPlan.ExecAsync(c.fftRBuf, c.fftCBuf)
		scale := 2 / float32(c.fwPlan.InputLen()) // ??
		zero1_async(c.gpuFFTKern[i])
		Madd2(c.gpuFFTKern[i], c.gpuFFTKern[i], c.fftCBuf, 0, scale)
	}
}

// store MFM image in output, based on magnetization in inp.
func (c *MFMConvolution) Exec(outp, inp, vol *slice.Slice, Msat MSlice) {
	for i := 0; i < 3; i++ {
		zero1_async(c.fftRBuf)
		copyPadMul(c.fftRBuf, inp.Comp(i), vol, c.kernSize, c.size, Msat)
		c.fwPlan.ExecAsync(c.fftRBuf, c.fftCBuf)

		Nx, Ny := c.fftKernSize[X]/2, c.fftKernSize[Y] //   ??
		kernMulC_async(c.fftCBuf, c.gpuFFTKern[i], Nx, Ny)

		c.bwPlan.ExecAsync(c.fftCBuf, c.fftRBuf)
		copyUnPad(outp.Comp(i), c.fftRBuf, c.size, c.kernSize)
	}
}

func (c *MFMConvolution) Reinit(lift, tipsize float64, cachedir string) {
	c.kern = kernels.MFMKernel(c.mesh, lift, tipsize, cachedir)
	c.initFFTKern3D()
}

// Initializes a convolution to evaluate the demag field for the given mesh geometry.
func NewMFM(mesh mesh.Mesh, lift, tipsize float64, cachedir string) *MFMConvolution {
	k := kernels.MFMKernel(mesh, lift, tipsize, cachedir)
	size := mesh.Size()
	c := new(MFMConvolution)
	c.size = size
	c.kern = k
	c.kernSize = k[X].Size()
	c.init()
	c.mesh = mesh
	return c
}
