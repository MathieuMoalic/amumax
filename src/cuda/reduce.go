package cuda

import (
	"math"
	"unsafe"

	"github.com/MathieuMoalic/amumax/src/cuda/cu"
	"github.com/MathieuMoalic/amumax/src/data"
	"github.com/MathieuMoalic/amumax/src/log"
)

//#include "reduce.h"
import "C"

// Block size for reduce kernels.
const REDUCE_BLOCKSIZE = C.REDUCE_BLOCKSIZE

// Sum of all elements.
func Sum(in *data.Slice) float32 {
	log.AssertMsg(in.NComp() == 1, "Component mismatch: input must have 1 component in Sum")
	out := reduceBuf(0)
	kReducesumAsync(in.DevPtr(0), out, 0, in.Len(), reducecfg)
	return copyback(out)
}

// Dot product.
func Dot(a, b *data.Slice) float32 {
	nComp := a.NComp()
	log.AssertMsg(nComp == b.NComp(), "Component mismatch: a and b must have the same number of components in Dot")
	out := reduceBuf(0)
	// not async over components
	for c := 0; c < nComp; c++ {
		kReducedotAsync(a.DevPtr(c), b.DevPtr(c), out, 0, a.Len(), reducecfg) // all components add to out
	}
	return copyback(out)
}

// Maximum of absolute values of all elements.
func MaxAbs(in *data.Slice) float32 {
	log.AssertMsg(in.NComp() == 1, "Component mismatch: input must have 1 component in MaxAbs")
	out := reduceBuf(0)
	kReducemaxabsAsync(in.DevPtr(0), out, 0, in.Len(), reducecfg)
	return copyback(out)
}

// Maximum of the norms of all vectors (x[i], y[i], z[i]).
//
//	max_i sqrt( x[i]*x[i] + y[i]*y[i] + z[i]*z[i] )
func MaxVecNorm(v *data.Slice) float64 {
	out := reduceBuf(0)
	kReducemaxvecnorm2Async(v.DevPtr(0), v.DevPtr(1), v.DevPtr(2), out, 0, v.Len(), reducecfg)
	return math.Sqrt(float64(copyback(out)))
}

// Maximum of the norms of the difference between all vectors (x1,y1,z1) and (x2,y2,z2)
//
//	(dx, dy, dz) = (x1, y1, z1) - (x2, y2, z2)
//	max_i sqrt( dx[i]*dx[i] + dy[i]*dy[i] + dz[i]*dz[i] )
func MaxVecDiff(x, y *data.Slice) float64 {
	log.AssertMsg(x.Len() == y.Len(), "Length mismatch: x and y must have the same length in MaxVecDiff")
	out := reduceBuf(0)
	kReducemaxvecdiff2Async(x.DevPtr(0), x.DevPtr(1), x.DevPtr(2),
		y.DevPtr(0), y.DevPtr(1), y.DevPtr(2),
		out, 0, x.Len(), reducecfg)
	return math.Sqrt(float64(copyback(out)))
}

var reduceBuffers chan unsafe.Pointer // pool of 1-float CUDA buffers for reduce

// return a 1-float CUDA reduction buffer from a pool
// initialized to initVal
func reduceBuf(initVal float32) unsafe.Pointer {
	if reduceBuffers == nil {
		initReduceBuf()
	}
	buf := <-reduceBuffers
	cu.MemsetD32Async(cu.DevicePtr(uintptr(buf)), math.Float32bits(initVal), 1, stream0)
	return buf
}

// copy back single float result from GPU and recycle buffer
func copyback(buf unsafe.Pointer) float32 {
	var result float32
	MemCpyDtoH(unsafe.Pointer(&result), buf, cu.SIZEOF_FLOAT32)
	reduceBuffers <- buf
	return result
}

// initialize pool of 1-float CUDA reduction buffers
func initReduceBuf() {
	const N = 128
	reduceBuffers = make(chan unsafe.Pointer, N)
	for i := 0; i < N; i++ {
		reduceBuffers <- MemAlloc(1 * cu.SIZEOF_FLOAT32)
	}
}

// launch configuration for reduce kernels
// 8 is typ. number of multiprocessors.
// could be improved but takes hardly ~1% of execution time
var reducecfg = &config{Grid: cu.Dim3{X: 8, Y: 1, Z: 1}, Block: cu.Dim3{X: REDUCE_BLOCKSIZE, Y: 1, Z: 1}}
