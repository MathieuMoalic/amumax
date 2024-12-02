package cuda_old

/*
 THIS FILE IS AUTO-GENERATED BY CUDA2GO.
 EDITING IS FUTILE.
*/

import (
	"sync"
	"unsafe"

	"github.com/MathieuMoalic/amumax/src/engine_old/cuda_old/cu"
	"github.com/MathieuMoalic/amumax/src/timer"
)

// CUDA handle for shiftz kernel
var shiftz_code cu.Function

// Stores the arguments for shiftz kernel invocation
type shiftz_args_t struct {
	arg_dst    unsafe.Pointer
	arg_src    unsafe.Pointer
	arg_Nx     int
	arg_Ny     int
	arg_Nz     int
	arg_shz    int
	arg_clampL float32
	arg_clampR float32
	argptr     [8]unsafe.Pointer
	sync.Mutex
}

// Stores the arguments for shiftz kernel invocation
var shiftz_args shiftz_args_t

func init() {
	// CUDA driver kernel call wants pointers to arguments, set them up once.
	shiftz_args.argptr[0] = unsafe.Pointer(&shiftz_args.arg_dst)
	shiftz_args.argptr[1] = unsafe.Pointer(&shiftz_args.arg_src)
	shiftz_args.argptr[2] = unsafe.Pointer(&shiftz_args.arg_Nx)
	shiftz_args.argptr[3] = unsafe.Pointer(&shiftz_args.arg_Ny)
	shiftz_args.argptr[4] = unsafe.Pointer(&shiftz_args.arg_Nz)
	shiftz_args.argptr[5] = unsafe.Pointer(&shiftz_args.arg_shz)
	shiftz_args.argptr[6] = unsafe.Pointer(&shiftz_args.arg_clampL)
	shiftz_args.argptr[7] = unsafe.Pointer(&shiftz_args.arg_clampR)
}

// Wrapper for shiftz CUDA kernel, asynchronous.
func k_shiftz_async(dst unsafe.Pointer, src unsafe.Pointer, Nx int, Ny int, Nz int, shz int, clampL float32, clampR float32, cfg *config) {
	if Synchronous { // debug
		Sync()
		timer.Start("shiftz")
	}

	shiftz_args.Lock()
	defer shiftz_args.Unlock()

	if shiftz_code == 0 {
		shiftz_code = fatbinLoad(shiftz_map, "shiftz")
	}

	shiftz_args.arg_dst = dst
	shiftz_args.arg_src = src
	shiftz_args.arg_Nx = Nx
	shiftz_args.arg_Ny = Ny
	shiftz_args.arg_Nz = Nz
	shiftz_args.arg_shz = shz
	shiftz_args.arg_clampL = clampL
	shiftz_args.arg_clampR = clampR

	args := shiftz_args.argptr[:]
	cu.LaunchKernel(shiftz_code, cfg.Grid.X, cfg.Grid.Y, cfg.Grid.Z, cfg.Block.X, cfg.Block.Y, cfg.Block.Z, 0, stream0, args)

	if Synchronous { // debug
		Sync()
		timer.Stop("shiftz")
	}
}

// maps compute capability on PTX code for shiftz kernel.
var shiftz_map = map[int]string{0: "",
	52: shiftz_ptx_52}

// shiftz PTX code for various compute capabilities.
const (
	shiftz_ptx_52 = `
.version 7.0
.target sm_52
.address_size 64

	// .globl	shiftz

.visible .entry shiftz(
	.param .u64 shiftz_param_0,
	.param .u64 shiftz_param_1,
	.param .u32 shiftz_param_2,
	.param .u32 shiftz_param_3,
	.param .u32 shiftz_param_4,
	.param .u32 shiftz_param_5,
	.param .f32 shiftz_param_6,
	.param .f32 shiftz_param_7
)
{
	.reg .pred 	%p<8>;
	.reg .f32 	%f<6>;
	.reg .b32 	%r<22>;
	.reg .b64 	%rd<9>;


	ld.param.u64 	%rd1, [shiftz_param_0];
	ld.param.u64 	%rd2, [shiftz_param_1];
	ld.param.u32 	%r5, [shiftz_param_2];
	ld.param.u32 	%r6, [shiftz_param_3];
	ld.param.u32 	%r7, [shiftz_param_4];
	ld.param.u32 	%r8, [shiftz_param_5];
	ld.param.f32 	%f5, [shiftz_param_6];
	ld.param.f32 	%f4, [shiftz_param_7];
	mov.u32 	%r9, %ntid.x;
	mov.u32 	%r10, %ctaid.x;
	mov.u32 	%r11, %tid.x;
	mad.lo.s32 	%r1, %r9, %r10, %r11;
	mov.u32 	%r12, %ntid.y;
	mov.u32 	%r13, %ctaid.y;
	mov.u32 	%r14, %tid.y;
	mad.lo.s32 	%r2, %r12, %r13, %r14;
	mov.u32 	%r15, %ntid.z;
	mov.u32 	%r16, %ctaid.z;
	mov.u32 	%r17, %tid.z;
	mad.lo.s32 	%r3, %r15, %r16, %r17;
	setp.lt.s32	%p1, %r1, %r5;
	setp.lt.s32	%p2, %r2, %r6;
	and.pred  	%p3, %p1, %p2;
	setp.lt.s32	%p4, %r3, %r7;
	and.pred  	%p5, %p3, %p4;
	@!%p5 bra 	BB0_5;
	bra.uni 	BB0_1;

BB0_1:
	sub.s32 	%r4, %r3, %r8;
	setp.lt.s32	%p6, %r4, 0;
	@%p6 bra 	BB0_4;

	setp.ge.s32	%p7, %r4, %r7;
	mov.f32 	%f5, %f4;
	@%p7 bra 	BB0_4;

	cvta.to.global.u64 	%rd3, %rd2;
	mad.lo.s32 	%r18, %r4, %r6, %r2;
	mad.lo.s32 	%r19, %r18, %r5, %r1;
	mul.wide.s32 	%rd4, %r19, 4;
	add.s64 	%rd5, %rd3, %rd4;
	ld.global.nc.f32 	%f5, [%rd5];

BB0_4:
	cvta.to.global.u64 	%rd6, %rd1;
	mad.lo.s32 	%r20, %r3, %r6, %r2;
	mad.lo.s32 	%r21, %r20, %r5, %r1;
	mul.wide.s32 	%rd7, %r21, 4;
	add.s64 	%rd8, %rd6, %rd7;
	st.global.f32 	[%rd8], %f5;

BB0_5:
	ret;
}


`
)