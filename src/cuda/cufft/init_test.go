package cufft

import (
	"fmt"

	"github.com/MathieuMoalic/amumax/src/cuda/cu"
)

// needed for all other tests.
func init() {
	cu.Init(0)
	ctx := cu.CtxCreate(cu.CTX_SCHED_AUTO, 0)
	cu.CtxSetCurrent(ctx)
	fmt.Println("Created CUDA context")
}
