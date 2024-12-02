package cuda

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

// CUDA handle for adduniaxialanisotropy2 kernel
var adduniaxialanisotropy2_code cu.Function

// Stores the arguments for adduniaxialanisotropy2 kernel invocation
type adduniaxialanisotropy2_args_t struct {
	arg_Bx     unsafe.Pointer
	arg_By     unsafe.Pointer
	arg_Bz     unsafe.Pointer
	arg_mx     unsafe.Pointer
	arg_my     unsafe.Pointer
	arg_mz     unsafe.Pointer
	arg_Ms_    unsafe.Pointer
	arg_Ms_mul float32
	arg_K1_    unsafe.Pointer
	arg_K1_mul float32
	arg_K2_    unsafe.Pointer
	arg_K2_mul float32
	arg_ux_    unsafe.Pointer
	arg_ux_mul float32
	arg_uy_    unsafe.Pointer
	arg_uy_mul float32
	arg_uz_    unsafe.Pointer
	arg_uz_mul float32
	arg_N      int
	argptr     [19]unsafe.Pointer
	sync.Mutex
}

// Stores the arguments for adduniaxialanisotropy2 kernel invocation
var adduniaxialanisotropy2_args adduniaxialanisotropy2_args_t

func init() {
	// CUDA driver kernel call wants pointers to arguments, set them up once.
	adduniaxialanisotropy2_args.argptr[0] = unsafe.Pointer(&adduniaxialanisotropy2_args.arg_Bx)
	adduniaxialanisotropy2_args.argptr[1] = unsafe.Pointer(&adduniaxialanisotropy2_args.arg_By)
	adduniaxialanisotropy2_args.argptr[2] = unsafe.Pointer(&adduniaxialanisotropy2_args.arg_Bz)
	adduniaxialanisotropy2_args.argptr[3] = unsafe.Pointer(&adduniaxialanisotropy2_args.arg_mx)
	adduniaxialanisotropy2_args.argptr[4] = unsafe.Pointer(&adduniaxialanisotropy2_args.arg_my)
	adduniaxialanisotropy2_args.argptr[5] = unsafe.Pointer(&adduniaxialanisotropy2_args.arg_mz)
	adduniaxialanisotropy2_args.argptr[6] = unsafe.Pointer(&adduniaxialanisotropy2_args.arg_Ms_)
	adduniaxialanisotropy2_args.argptr[7] = unsafe.Pointer(&adduniaxialanisotropy2_args.arg_Ms_mul)
	adduniaxialanisotropy2_args.argptr[8] = unsafe.Pointer(&adduniaxialanisotropy2_args.arg_K1_)
	adduniaxialanisotropy2_args.argptr[9] = unsafe.Pointer(&adduniaxialanisotropy2_args.arg_K1_mul)
	adduniaxialanisotropy2_args.argptr[10] = unsafe.Pointer(&adduniaxialanisotropy2_args.arg_K2_)
	adduniaxialanisotropy2_args.argptr[11] = unsafe.Pointer(&adduniaxialanisotropy2_args.arg_K2_mul)
	adduniaxialanisotropy2_args.argptr[12] = unsafe.Pointer(&adduniaxialanisotropy2_args.arg_ux_)
	adduniaxialanisotropy2_args.argptr[13] = unsafe.Pointer(&adduniaxialanisotropy2_args.arg_ux_mul)
	adduniaxialanisotropy2_args.argptr[14] = unsafe.Pointer(&adduniaxialanisotropy2_args.arg_uy_)
	adduniaxialanisotropy2_args.argptr[15] = unsafe.Pointer(&adduniaxialanisotropy2_args.arg_uy_mul)
	adduniaxialanisotropy2_args.argptr[16] = unsafe.Pointer(&adduniaxialanisotropy2_args.arg_uz_)
	adduniaxialanisotropy2_args.argptr[17] = unsafe.Pointer(&adduniaxialanisotropy2_args.arg_uz_mul)
	adduniaxialanisotropy2_args.argptr[18] = unsafe.Pointer(&adduniaxialanisotropy2_args.arg_N)
}

// Wrapper for adduniaxialanisotropy2 CUDA kernel, asynchronous.
func k_adduniaxialanisotropy2_async(Bx unsafe.Pointer, By unsafe.Pointer, Bz unsafe.Pointer, mx unsafe.Pointer, my unsafe.Pointer, mz unsafe.Pointer, Ms_ unsafe.Pointer, Ms_mul float32, K1_ unsafe.Pointer, K1_mul float32, K2_ unsafe.Pointer, K2_mul float32, ux_ unsafe.Pointer, ux_mul float32, uy_ unsafe.Pointer, uy_mul float32, uz_ unsafe.Pointer, uz_mul float32, N int, cfg *config) {
	if Synchronous { // debug
		Sync()
		timer.Start("adduniaxialanisotropy2")
	}

	adduniaxialanisotropy2_args.Lock()
	defer adduniaxialanisotropy2_args.Unlock()

	if adduniaxialanisotropy2_code == 0 {
		adduniaxialanisotropy2_code = fatbinLoad(adduniaxialanisotropy2_map, "adduniaxialanisotropy2")
	}

	adduniaxialanisotropy2_args.arg_Bx = Bx
	adduniaxialanisotropy2_args.arg_By = By
	adduniaxialanisotropy2_args.arg_Bz = Bz
	adduniaxialanisotropy2_args.arg_mx = mx
	adduniaxialanisotropy2_args.arg_my = my
	adduniaxialanisotropy2_args.arg_mz = mz
	adduniaxialanisotropy2_args.arg_Ms_ = Ms_
	adduniaxialanisotropy2_args.arg_Ms_mul = Ms_mul
	adduniaxialanisotropy2_args.arg_K1_ = K1_
	adduniaxialanisotropy2_args.arg_K1_mul = K1_mul
	adduniaxialanisotropy2_args.arg_K2_ = K2_
	adduniaxialanisotropy2_args.arg_K2_mul = K2_mul
	adduniaxialanisotropy2_args.arg_ux_ = ux_
	adduniaxialanisotropy2_args.arg_ux_mul = ux_mul
	adduniaxialanisotropy2_args.arg_uy_ = uy_
	adduniaxialanisotropy2_args.arg_uy_mul = uy_mul
	adduniaxialanisotropy2_args.arg_uz_ = uz_
	adduniaxialanisotropy2_args.arg_uz_mul = uz_mul
	adduniaxialanisotropy2_args.arg_N = N

	args := adduniaxialanisotropy2_args.argptr[:]
	cu.LaunchKernel(adduniaxialanisotropy2_code, cfg.Grid.X, cfg.Grid.Y, cfg.Grid.Z, cfg.Block.X, cfg.Block.Y, cfg.Block.Z, 0, stream0, args)

	if Synchronous { // debug
		Sync()
		timer.Stop("adduniaxialanisotropy2")
	}
}

// maps compute capability on PTX code for adduniaxialanisotropy2 kernel.
var adduniaxialanisotropy2_map = map[int]string{0: "",
	52: adduniaxialanisotropy2_ptx_52}

// adduniaxialanisotropy2 PTX code for various compute capabilities.
const (
	adduniaxialanisotropy2_ptx_52 = `
.version 7.0
.target sm_52
.address_size 64

	// .globl	adduniaxialanisotropy2

.visible .entry adduniaxialanisotropy2(
	.param .u64 adduniaxialanisotropy2_param_0,
	.param .u64 adduniaxialanisotropy2_param_1,
	.param .u64 adduniaxialanisotropy2_param_2,
	.param .u64 adduniaxialanisotropy2_param_3,
	.param .u64 adduniaxialanisotropy2_param_4,
	.param .u64 adduniaxialanisotropy2_param_5,
	.param .u64 adduniaxialanisotropy2_param_6,
	.param .f32 adduniaxialanisotropy2_param_7,
	.param .u64 adduniaxialanisotropy2_param_8,
	.param .f32 adduniaxialanisotropy2_param_9,
	.param .u64 adduniaxialanisotropy2_param_10,
	.param .f32 adduniaxialanisotropy2_param_11,
	.param .u64 adduniaxialanisotropy2_param_12,
	.param .f32 adduniaxialanisotropy2_param_13,
	.param .u64 adduniaxialanisotropy2_param_14,
	.param .f32 adduniaxialanisotropy2_param_15,
	.param .u64 adduniaxialanisotropy2_param_16,
	.param .f32 adduniaxialanisotropy2_param_17,
	.param .u32 adduniaxialanisotropy2_param_18
)
{
	.reg .pred 	%p<10>;
	.reg .f32 	%f<72>;
	.reg .b32 	%r<9>;
	.reg .b64 	%rd<44>;


	ld.param.u64 	%rd1, [adduniaxialanisotropy2_param_0];
	ld.param.u64 	%rd2, [adduniaxialanisotropy2_param_1];
	ld.param.u64 	%rd3, [adduniaxialanisotropy2_param_2];
	ld.param.u64 	%rd4, [adduniaxialanisotropy2_param_3];
	ld.param.u64 	%rd5, [adduniaxialanisotropy2_param_4];
	ld.param.u64 	%rd6, [adduniaxialanisotropy2_param_5];
	ld.param.u64 	%rd7, [adduniaxialanisotropy2_param_6];
	ld.param.f32 	%f68, [adduniaxialanisotropy2_param_7];
	ld.param.u64 	%rd8, [adduniaxialanisotropy2_param_8];
	ld.param.f32 	%f70, [adduniaxialanisotropy2_param_9];
	ld.param.u64 	%rd9, [adduniaxialanisotropy2_param_10];
	ld.param.f32 	%f71, [adduniaxialanisotropy2_param_11];
	ld.param.u64 	%rd10, [adduniaxialanisotropy2_param_12];
	ld.param.f32 	%f64, [adduniaxialanisotropy2_param_13];
	ld.param.u64 	%rd11, [adduniaxialanisotropy2_param_14];
	ld.param.f32 	%f65, [adduniaxialanisotropy2_param_15];
	ld.param.u64 	%rd12, [adduniaxialanisotropy2_param_16];
	ld.param.f32 	%f66, [adduniaxialanisotropy2_param_17];
	ld.param.u32 	%r2, [adduniaxialanisotropy2_param_18];
	mov.u32 	%r3, %nctaid.x;
	mov.u32 	%r4, %ctaid.y;
	mov.u32 	%r5, %ctaid.x;
	mad.lo.s32 	%r6, %r3, %r4, %r5;
	mov.u32 	%r7, %ntid.x;
	mov.u32 	%r8, %tid.x;
	mad.lo.s32 	%r1, %r6, %r7, %r8;
	setp.ge.s32	%p1, %r1, %r2;
	@%p1 bra 	BB0_18;

	setp.eq.s64	%p2, %rd10, 0;
	@%p2 bra 	BB0_3;

	cvta.to.global.u64 	%rd13, %rd10;
	mul.wide.s32 	%rd14, %r1, 4;
	add.s64 	%rd15, %rd13, %rd14;
	ld.global.nc.f32 	%f27, [%rd15];
	mul.f32 	%f64, %f27, %f64;

BB0_3:
	setp.eq.s64	%p3, %rd11, 0;
	@%p3 bra 	BB0_5;

	cvta.to.global.u64 	%rd16, %rd11;
	mul.wide.s32 	%rd17, %r1, 4;
	add.s64 	%rd18, %rd16, %rd17;
	ld.global.nc.f32 	%f28, [%rd18];
	mul.f32 	%f65, %f28, %f65;

BB0_5:
	setp.eq.s64	%p4, %rd12, 0;
	@%p4 bra 	BB0_7;

	cvta.to.global.u64 	%rd19, %rd12;
	mul.wide.s32 	%rd20, %r1, 4;
	add.s64 	%rd21, %rd19, %rd20;
	ld.global.nc.f32 	%f29, [%rd21];
	mul.f32 	%f66, %f29, %f66;

BB0_7:
	mul.f32 	%f31, %f65, %f65;
	fma.rn.f32 	%f32, %f64, %f64, %f31;
	fma.rn.f32 	%f33, %f66, %f66, %f32;
	sqrt.rn.f32 	%f7, %f33;
	mov.f32 	%f67, 0f00000000;
	setp.eq.f32	%p5, %f7, 0f00000000;
	@%p5 bra 	BB0_9;

	rcp.rn.f32 	%f67, %f7;

BB0_9:
	mul.f32 	%f10, %f64, %f67;
	mul.f32 	%f11, %f65, %f67;
	mul.f32 	%f12, %f66, %f67;
	setp.eq.s64	%p6, %rd7, 0;
	@%p6 bra 	BB0_11;

	cvta.to.global.u64 	%rd22, %rd7;
	mul.wide.s32 	%rd23, %r1, 4;
	add.s64 	%rd24, %rd22, %rd23;
	ld.global.nc.f32 	%f34, [%rd24];
	mul.f32 	%f68, %f34, %f68;

BB0_11:
	setp.eq.f32	%p7, %f68, 0f00000000;
	mov.f32 	%f69, 0f00000000;
	@%p7 bra 	BB0_13;

	rcp.rn.f32 	%f69, %f68;

BB0_13:
	setp.eq.s64	%p8, %rd8, 0;
	@%p8 bra 	BB0_15;

	cvta.to.global.u64 	%rd25, %rd8;
	mul.wide.s32 	%rd26, %r1, 4;
	add.s64 	%rd27, %rd25, %rd26;
	ld.global.nc.f32 	%f36, [%rd27];
	mul.f32 	%f70, %f36, %f70;

BB0_15:
	setp.eq.s64	%p9, %rd9, 0;
	@%p9 bra 	BB0_17;

	cvta.to.global.u64 	%rd28, %rd9;
	mul.wide.s32 	%rd29, %r1, 4;
	add.s64 	%rd30, %rd28, %rd29;
	ld.global.nc.f32 	%f37, [%rd30];
	mul.f32 	%f71, %f37, %f71;

BB0_17:
	cvta.to.global.u64 	%rd31, %rd4;
	mul.wide.s32 	%rd32, %r1, 4;
	add.s64 	%rd33, %rd31, %rd32;
	cvta.to.global.u64 	%rd34, %rd5;
	add.s64 	%rd35, %rd34, %rd32;
	cvta.to.global.u64 	%rd36, %rd6;
	add.s64 	%rd37, %rd36, %rd32;
	ld.global.nc.f32 	%f38, [%rd33];
	ld.global.nc.f32 	%f39, [%rd35];
	mul.f32 	%f40, %f11, %f39;
	fma.rn.f32 	%f41, %f10, %f38, %f40;
	ld.global.nc.f32 	%f42, [%rd37];
	fma.rn.f32 	%f43, %f12, %f42, %f41;
	mul.f32 	%f44, %f69, %f70;
	fma.rn.f32 	%f45, %f69, %f70, %f44;
	mul.f32 	%f46, %f45, %f43;
	mul.f32 	%f47, %f69, %f71;
	mul.f32 	%f48, %f47, 0f40800000;
	mul.f32 	%f49, %f43, %f43;
	mul.f32 	%f50, %f43, %f49;
	mul.f32 	%f51, %f48, %f50;
	mul.f32 	%f52, %f10, %f51;
	mul.f32 	%f53, %f11, %f51;
	mul.f32 	%f54, %f12, %f51;
	fma.rn.f32 	%f55, %f10, %f46, %f52;
	fma.rn.f32 	%f56, %f11, %f46, %f53;
	fma.rn.f32 	%f57, %f12, %f46, %f54;
	cvta.to.global.u64 	%rd38, %rd1;
	add.s64 	%rd39, %rd38, %rd32;
	ld.global.f32 	%f58, [%rd39];
	add.f32 	%f59, %f58, %f55;
	st.global.f32 	[%rd39], %f59;
	cvta.to.global.u64 	%rd40, %rd2;
	add.s64 	%rd41, %rd40, %rd32;
	ld.global.f32 	%f60, [%rd41];
	add.f32 	%f61, %f60, %f56;
	st.global.f32 	[%rd41], %f61;
	cvta.to.global.u64 	%rd42, %rd3;
	add.s64 	%rd43, %rd42, %rd32;
	ld.global.f32 	%f62, [%rd43];
	add.f32 	%f63, %f62, %f57;
	st.global.f32 	[%rd43], %f63;

BB0_18:
	ret;
}


`
)
