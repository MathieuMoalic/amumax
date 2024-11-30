// Package cuda provides GPU interaction
package cuda

import (
	"fmt"
	"runtime"

	"github.com/MathieuMoalic/amumax/src/cuda/cu"
	"github.com/MathieuMoalic/amumax/src/log_old"
)

var (
	GPUInfo_old string     // Human-readable GPU description
	Synchronous bool       // for debug: synchronize stream0 at every kernel launch
	cudaCtx     cu.Context // global CUDA context
)

// Locks to an OS thread and initializes CUDA for that thread.
func Init(gpu int) [6]string {
	if cudaCtx != 0 {
		return [6]string{"", "", "", "", "", ""} // needed for tests
	}

	runtime.LockOSThread()
	tryCuInit()
	dev := cu.Device(gpu)
	cudaCtx = cu.CtxCreate(cu.CTX_SCHED_YIELD, dev)
	cudaCtx.SetCurrent()

	M, m := dev.ComputeCapability()
	DriverVersion := cu.Version()
	DevName := dev.Name()
	TotalMem := dev.TotalMem()
	GPUInfo_old = fmt.Sprintf("%s(%dMB), CUDA Driver %d.%d, cc=%d.%d",
		DevName, (TotalMem)/(1024*1024), DriverVersion/1000, (DriverVersion%1000)/10, M, m)

	if M < 5 {
		log_old.Log.ErrAndExit("GPU has insufficient compute capability, need 5.0 or higher.")
	}
	if Synchronous {
		log_old.Log.Info("DEBUG: synchronized CUDA calls")
	}

	// test PTX load so that we can catch CUDA_ERROR_NO_BINARY_FOR_GPU early
	fatbinLoad(madd2_map, "madd2")
	GpuInfo1 := [6]string{
		fmt.Sprintf("%d.%d", cu.CUDA_VERSION/1000, (cu.CUDA_VERSION%1000)/10),
		fmt.Sprintf("%d", UseCC),
		DevName,
		fmt.Sprintf("%d", TotalMem/(1024*1024)),
		fmt.Sprintf("%d.%d", DriverVersion/1000, (DriverVersion%1000)/10),
		fmt.Sprintf("%d.%d", M, m),
	}
	return GpuInfo1
}

// cu.Init(), but error is fatal and does not dump stack.
func tryCuInit() {
	defer func() {
		err := recover()
		if err == cu.ERROR_UNKNOWN {
			log_old.Log.ErrAndExit("\n CUDA unknown error\n")
		}
		if err != nil {
			log_old.Log.PanicIfError(fmt.Errorf("%v", fmt.Sprint(err)))
		}
	}()
	cu.Init(0)
}

// Global stream used for everything
const stream0 = cu.Stream(0)

// Synchronize the global stream
// This is called before and after all memcopy operations between host and device.
func Sync() {
	stream0.Synchronize()
}
