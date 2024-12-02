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

// CUDA handle for normalize kernel
var normalize_code cu.Function

// Stores the arguments for normalize kernel invocation
type normalize_args_t struct {
	arg_vx  unsafe.Pointer
	arg_vy  unsafe.Pointer
	arg_vz  unsafe.Pointer
	arg_vol unsafe.Pointer
	arg_N   int
	argptr  [5]unsafe.Pointer
	sync.Mutex
}

// Stores the arguments for normalize kernel invocation
var normalize_args normalize_args_t

func init() {
	// CUDA driver kernel call wants pointers to arguments, set them up once.
	normalize_args.argptr[0] = unsafe.Pointer(&normalize_args.arg_vx)
	normalize_args.argptr[1] = unsafe.Pointer(&normalize_args.arg_vy)
	normalize_args.argptr[2] = unsafe.Pointer(&normalize_args.arg_vz)
	normalize_args.argptr[3] = unsafe.Pointer(&normalize_args.arg_vol)
	normalize_args.argptr[4] = unsafe.Pointer(&normalize_args.arg_N)
}

// Wrapper for normalize CUDA kernel, asynchronous.
func k_normalize_async(vx unsafe.Pointer, vy unsafe.Pointer, vz unsafe.Pointer, vol unsafe.Pointer, N int, cfg *config) {
	if Synchronous { // debug
		Sync()
		timer_old.Start("normalize")
	}

	normalize_args.Lock()
	defer normalize_args.Unlock()

	if normalize_code == 0 {
		normalize_code = fatbinLoad(normalize_map, "normalize")
	}

	normalize_args.arg_vx = vx
	normalize_args.arg_vy = vy
	normalize_args.arg_vz = vz
	normalize_args.arg_vol = vol
	normalize_args.arg_N = N

	args := normalize_args.argptr[:]
	cu.LaunchKernel(normalize_code, cfg.Grid.X, cfg.Grid.Y, cfg.Grid.Z, cfg.Block.X, cfg.Block.Y, cfg.Block.Z, 0, stream0, args)

	if Synchronous { // debug
		Sync()
		timer_old.Stop("normalize")
	}
}

// maps compute capability on PTX code for normalize kernel.
var normalize_map = map[int]string{0: "",
	52: normalize_ptx_52}

// normalize PTX code for various compute capabilities.
const (
	normalize_ptx_52 = `
.version 7.0
.target sm_52
.address_size 64

	// .globl	normalize

.visible .entry normalize(
	.param .u64 normalize_param_0,
	.param .u64 normalize_param_1,
	.param .u64 normalize_param_2,
	.param .u64 normalize_param_3,
	.param .u32 normalize_param_4
)
{
	.reg .pred 	%p<4>;
	.reg .f32 	%f<22>;
	.reg .b32 	%r<9>;
	.reg .b64 	%rd<15>;


	ld.param.u64 	%rd4, [normalize_param_0];
	ld.param.u64 	%rd5, [normalize_param_1];
	ld.param.u64 	%rd6, [normalize_param_2];
	ld.param.u64 	%rd7, [normalize_param_3];
	ld.param.u32 	%r2, [normalize_param_4];
	mov.u32 	%r3, %nctaid.x;
	mov.u32 	%r4, %ctaid.y;
	mov.u32 	%r5, %ctaid.x;
	mad.lo.s32 	%r6, %r3, %r4, %r5;
	mov.u32 	%r7, %ntid.x;
	mov.u32 	%r8, %tid.x;
	mad.lo.s32 	%r1, %r6, %r7, %r8;
	setp.ge.s32	%p1, %r1, %r2;
	@%p1 bra 	BB0_6;

	setp.eq.s64	%p2, %rd7, 0;
	mov.f32 	%f20, 0f3F800000;
	@%p2 bra 	BB0_3;

	cvta.to.global.u64 	%rd8, %rd7;
	mul.wide.s32 	%rd9, %r1, 4;
	add.s64 	%rd10, %rd8, %rd9;
	ld.global.nc.f32 	%f20, [%rd10];

BB0_3:
	cvta.to.global.u64 	%rd11, %rd6;
	cvta.to.global.u64 	%rd12, %rd5;
	cvta.to.global.u64 	%rd13, %rd4;
	mul.wide.s32 	%rd14, %r1, 4;
	add.s64 	%rd1, %rd13, %rd14;
	ld.global.f32 	%f11, [%rd1];
	mul.f32 	%f3, %f20, %f11;
	add.s64 	%rd2, %rd12, %rd14;
	ld.global.f32 	%f12, [%rd2];
	mul.f32 	%f4, %f20, %f12;
	add.s64 	%rd3, %rd11, %rd14;
	ld.global.f32 	%f13, [%rd3];
	mul.f32 	%f5, %f20, %f13;
	mul.f32 	%f14, %f4, %f4;
	fma.rn.f32 	%f15, %f3, %f3, %f14;
	fma.rn.f32 	%f16, %f5, %f5, %f15;
	sqrt.rn.f32 	%f6, %f16;
	mov.f32 	%f21, 0f00000000;
	setp.eq.f32	%p3, %f6, 0f00000000;
	@%p3 bra 	BB0_5;

	rcp.rn.f32 	%f21, %f6;

BB0_5:
	mul.f32 	%f17, %f3, %f21;
	st.global.f32 	[%rd1], %f17;
	mul.f32 	%f18, %f4, %f21;
	st.global.f32 	[%rd2], %f18;
	mul.f32 	%f19, %f5, %f21;
	st.global.f32 	[%rd3], %f19;

BB0_6:
	ret;
}


`
)
