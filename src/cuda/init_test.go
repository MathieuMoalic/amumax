package cuda

import (
	"github.com/MathieuMoalic/amumax/src/engine_old/cuda_old/cu"
)

// needed for all other tests.
func init() {
	cu.Init(0)
	ctx := cu.CtxCreate(cu.CTX_SCHED_AUTO, 0)
	cu.CtxSetCurrent(ctx)
}
