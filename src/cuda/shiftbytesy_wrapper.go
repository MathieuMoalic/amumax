package cuda

/*
 THIS FILE IS AUTO-GENERATED BY CUDA2GO.
 EDITING IS FUTILE.
*/

import (
	"sync"
	"unsafe"

	"github.com/MathieuMoalic/amumax/src/cuda/cu"
	"github.com/MathieuMoalic/amumax/src/engine_old/timer_old"
)

// CUDA handle for shiftbytesy kernel
var shiftbytesy_code cu.Function

// Stores the arguments for shiftbytesy kernel invocation
type shiftbytesy_args_t struct {
	arg_dst   unsafe.Pointer
	arg_src   unsafe.Pointer
	arg_Nx    int
	arg_Ny    int
	arg_Nz    int
	arg_shy   int
	arg_clamp byte
	argptr    [7]unsafe.Pointer
	sync.Mutex
}

// Stores the arguments for shiftbytesy kernel invocation
var shiftbytesy_args shiftbytesy_args_t

func init() {
	// CUDA driver kernel call wants pointers to arguments, set them up once.
	shiftbytesy_args.argptr[0] = unsafe.Pointer(&shiftbytesy_args.arg_dst)
	shiftbytesy_args.argptr[1] = unsafe.Pointer(&shiftbytesy_args.arg_src)
	shiftbytesy_args.argptr[2] = unsafe.Pointer(&shiftbytesy_args.arg_Nx)
	shiftbytesy_args.argptr[3] = unsafe.Pointer(&shiftbytesy_args.arg_Ny)
	shiftbytesy_args.argptr[4] = unsafe.Pointer(&shiftbytesy_args.arg_Nz)
	shiftbytesy_args.argptr[5] = unsafe.Pointer(&shiftbytesy_args.arg_shy)
	shiftbytesy_args.argptr[6] = unsafe.Pointer(&shiftbytesy_args.arg_clamp)
}

// Wrapper for shiftbytesy CUDA kernel, asynchronous.
func k_shiftbytesy_async(dst unsafe.Pointer, src unsafe.Pointer, Nx int, Ny int, Nz int, shy int, clamp byte, cfg *config) {
	if Synchronous { // debug
		Sync()
		timer_old.Start("shiftbytesy")
	}

	shiftbytesy_args.Lock()
	defer shiftbytesy_args.Unlock()

	if shiftbytesy_code == 0 {
		shiftbytesy_code = fatbinLoad(shiftbytesy_map, "shiftbytesy")
	}

	shiftbytesy_args.arg_dst = dst
	shiftbytesy_args.arg_src = src
	shiftbytesy_args.arg_Nx = Nx
	shiftbytesy_args.arg_Ny = Ny
	shiftbytesy_args.arg_Nz = Nz
	shiftbytesy_args.arg_shy = shy
	shiftbytesy_args.arg_clamp = clamp

	args := shiftbytesy_args.argptr[:]
	cu.LaunchKernel(shiftbytesy_code, cfg.Grid.X, cfg.Grid.Y, cfg.Grid.Z, cfg.Block.X, cfg.Block.Y, cfg.Block.Z, 0, stream0, args)

	if Synchronous { // debug
		Sync()
		timer_old.Stop("shiftbytesy")
	}
}

// maps compute capability on PTX code for shiftbytesy kernel.
var shiftbytesy_map = map[int]string{0: "",
	52: shiftbytesy_ptx_52}

// shiftbytesy PTX code for various compute capabilities.
const (
	shiftbytesy_ptx_52 = `
.version 7.0
.target sm_52
.address_size 64

	// .globl	shiftbytesy

.visible .entry shiftbytesy(
	.param .u64 shiftbytesy_param_0,
	.param .u64 shiftbytesy_param_1,
	.param .u32 shiftbytesy_param_2,
	.param .u32 shiftbytesy_param_3,
	.param .u32 shiftbytesy_param_4,
	.param .u32 shiftbytesy_param_5,
	.param .u8 shiftbytesy_param_6
)
{
	.reg .pred 	%p<9>;
	.reg .b16 	%rs<5>;
	.reg .b32 	%r<23>;
	.reg .b64 	%rd<9>;


	ld.param.u64 	%rd1, [shiftbytesy_param_0];
	ld.param.u64 	%rd2, [shiftbytesy_param_1];
	ld.param.u32 	%r6, [shiftbytesy_param_2];
	ld.param.u32 	%r7, [shiftbytesy_param_3];
	ld.param.u32 	%r9, [shiftbytesy_param_4];
	ld.param.u32 	%r8, [shiftbytesy_param_5];
	ld.param.u8 	%rs4, [shiftbytesy_param_6];
	mov.u32 	%r10, %ntid.x;
	mov.u32 	%r11, %ctaid.x;
	mov.u32 	%r12, %tid.x;
	mad.lo.s32 	%r1, %r10, %r11, %r12;
	mov.u32 	%r13, %ntid.y;
	mov.u32 	%r14, %ctaid.y;
	mov.u32 	%r15, %tid.y;
	mad.lo.s32 	%r2, %r13, %r14, %r15;
	mov.u32 	%r16, %ntid.z;
	mov.u32 	%r17, %ctaid.z;
	mov.u32 	%r18, %tid.z;
	mad.lo.s32 	%r3, %r16, %r17, %r18;
	setp.lt.s32	%p1, %r1, %r6;
	setp.lt.s32	%p2, %r2, %r7;
	and.pred  	%p3, %p1, %p2;
	setp.lt.s32	%p4, %r3, %r9;
	and.pred  	%p5, %p3, %p4;
	@!%p5 bra 	BB0_4;
	bra.uni 	BB0_1;

BB0_1:
	sub.s32 	%r4, %r2, %r8;
	setp.lt.s32	%p6, %r4, 0;
	setp.ge.s32	%p7, %r4, %r7;
	or.pred  	%p8, %p6, %p7;
	mul.lo.s32 	%r5, %r3, %r7;
	@%p8 bra 	BB0_3;

	cvta.to.global.u64 	%rd3, %rd2;
	add.s32 	%r19, %r5, %r4;
	mad.lo.s32 	%r20, %r19, %r6, %r1;
	cvt.s64.s32	%rd4, %r20;
	add.s64 	%rd5, %rd3, %rd4;
	ld.global.nc.u8 	%rs4, [%rd5];

BB0_3:
	cvta.to.global.u64 	%rd6, %rd1;
	add.s32 	%r21, %r5, %r2;
	mad.lo.s32 	%r22, %r21, %r6, %r1;
	cvt.s64.s32	%rd7, %r22;
	add.s64 	%rd8, %rd6, %rd7;
	st.global.u8 	[%rd8], %rs4;

BB0_4:
	ret;
}


`
)
