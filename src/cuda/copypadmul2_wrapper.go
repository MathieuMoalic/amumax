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

// CUDA handle for copypadmul2 kernel
var copypadmul2_code cu.Function

// Stores the arguments for copypadmul2 kernel invocation
type copypadmul2_args_t struct {
	arg_dst    unsafe.Pointer
	arg_Dx     int
	arg_Dy     int
	arg_Dz     int
	arg_src    unsafe.Pointer
	arg_Sx     int
	arg_Sy     int
	arg_Sz     int
	arg_Ms_    unsafe.Pointer
	arg_Ms_mul float32
	arg_vol    unsafe.Pointer
	argptr     [11]unsafe.Pointer
	sync.Mutex
}

// Stores the arguments for copypadmul2 kernel invocation
var copypadmul2_args copypadmul2_args_t

func init() {
	// CUDA driver kernel call wants pointers to arguments, set them up once.
	copypadmul2_args.argptr[0] = unsafe.Pointer(&copypadmul2_args.arg_dst)
	copypadmul2_args.argptr[1] = unsafe.Pointer(&copypadmul2_args.arg_Dx)
	copypadmul2_args.argptr[2] = unsafe.Pointer(&copypadmul2_args.arg_Dy)
	copypadmul2_args.argptr[3] = unsafe.Pointer(&copypadmul2_args.arg_Dz)
	copypadmul2_args.argptr[4] = unsafe.Pointer(&copypadmul2_args.arg_src)
	copypadmul2_args.argptr[5] = unsafe.Pointer(&copypadmul2_args.arg_Sx)
	copypadmul2_args.argptr[6] = unsafe.Pointer(&copypadmul2_args.arg_Sy)
	copypadmul2_args.argptr[7] = unsafe.Pointer(&copypadmul2_args.arg_Sz)
	copypadmul2_args.argptr[8] = unsafe.Pointer(&copypadmul2_args.arg_Ms_)
	copypadmul2_args.argptr[9] = unsafe.Pointer(&copypadmul2_args.arg_Ms_mul)
	copypadmul2_args.argptr[10] = unsafe.Pointer(&copypadmul2_args.arg_vol)
}

// Wrapper for copypadmul2 CUDA kernel, asynchronous.
func k_copypadmul2_async(dst unsafe.Pointer, Dx int, Dy int, Dz int, src unsafe.Pointer, Sx int, Sy int, Sz int, Ms_ unsafe.Pointer, Ms_mul float32, vol unsafe.Pointer, cfg *config) {
	if Synchronous { // debug
		Sync()
		timer.Start("copypadmul2")
	}

	copypadmul2_args.Lock()
	defer copypadmul2_args.Unlock()

	if copypadmul2_code == 0 {
		copypadmul2_code = fatbinLoad(copypadmul2_map, "copypadmul2")
	}

	copypadmul2_args.arg_dst = dst
	copypadmul2_args.arg_Dx = Dx
	copypadmul2_args.arg_Dy = Dy
	copypadmul2_args.arg_Dz = Dz
	copypadmul2_args.arg_src = src
	copypadmul2_args.arg_Sx = Sx
	copypadmul2_args.arg_Sy = Sy
	copypadmul2_args.arg_Sz = Sz
	copypadmul2_args.arg_Ms_ = Ms_
	copypadmul2_args.arg_Ms_mul = Ms_mul
	copypadmul2_args.arg_vol = vol

	args := copypadmul2_args.argptr[:]
	cu.LaunchKernel(copypadmul2_code, cfg.Grid.X, cfg.Grid.Y, cfg.Grid.Z, cfg.Block.X, cfg.Block.Y, cfg.Block.Z, 0, stream0, args)

	if Synchronous { // debug
		Sync()
		timer.Stop("copypadmul2")
	}
}

// maps compute capability on PTX code for copypadmul2 kernel.
var copypadmul2_map = map[int]string{0: "",
	52: copypadmul2_ptx_52}

// copypadmul2 PTX code for various compute capabilities.
const (
	copypadmul2_ptx_52 = `
.version 7.0
.target sm_52
.address_size 64

	// .globl	copypadmul2

.visible .entry copypadmul2(
	.param .u64 copypadmul2_param_0,
	.param .u32 copypadmul2_param_1,
	.param .u32 copypadmul2_param_2,
	.param .u32 copypadmul2_param_3,
	.param .u64 copypadmul2_param_4,
	.param .u32 copypadmul2_param_5,
	.param .u32 copypadmul2_param_6,
	.param .u32 copypadmul2_param_7,
	.param .u64 copypadmul2_param_8,
	.param .f32 copypadmul2_param_9,
	.param .u64 copypadmul2_param_10
)
{
	.reg .pred 	%p<8>;
	.reg .f32 	%f<14>;
	.reg .b32 	%r<22>;
	.reg .f64 	%fd<3>;
	.reg .b64 	%rd<17>;


	ld.param.u64 	%rd1, [copypadmul2_param_0];
	ld.param.u32 	%r5, [copypadmul2_param_1];
	ld.param.u32 	%r6, [copypadmul2_param_2];
	ld.param.u64 	%rd2, [copypadmul2_param_4];
	ld.param.u32 	%r7, [copypadmul2_param_5];
	ld.param.u32 	%r8, [copypadmul2_param_6];
	ld.param.u32 	%r9, [copypadmul2_param_7];
	ld.param.u64 	%rd3, [copypadmul2_param_8];
	ld.param.f32 	%f12, [copypadmul2_param_9];
	ld.param.u64 	%rd4, [copypadmul2_param_10];
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
	setp.lt.s32	%p1, %r1, %r7;
	setp.lt.s32	%p2, %r2, %r8;
	and.pred  	%p3, %p1, %p2;
	setp.lt.s32	%p4, %r3, %r9;
	and.pred  	%p5, %p3, %p4;
	@!%p5 bra 	BB0_6;
	bra.uni 	BB0_1;

BB0_1:
	mad.lo.s32 	%r19, %r3, %r8, %r2;
	mad.lo.s32 	%r4, %r19, %r7, %r1;
	setp.eq.s64	%p6, %rd3, 0;
	@%p6 bra 	BB0_3;

	cvta.to.global.u64 	%rd5, %rd3;
	mul.wide.s32 	%rd6, %r4, 4;
	add.s64 	%rd7, %rd5, %rd6;
	ld.global.nc.f32 	%f6, [%rd7];
	mul.f32 	%f12, %f6, %f12;

BB0_3:
	setp.eq.s64	%p7, %rd4, 0;
	mov.f32 	%f13, 0f3F800000;
	@%p7 bra 	BB0_5;

	cvta.to.global.u64 	%rd8, %rd4;
	mul.wide.s32 	%rd9, %r4, 4;
	add.s64 	%rd10, %rd8, %rd9;
	ld.global.nc.f32 	%f13, [%rd10];

BB0_5:
	cvta.to.global.u64 	%rd11, %rd1;
	cvta.to.global.u64 	%rd12, %rd2;
	mul.wide.s32 	%rd13, %r4, 4;
	add.s64 	%rd14, %rd12, %rd13;
	ld.global.nc.f32 	%f8, [%rd14];
	cvt.f64.f32	%fd1, %f12;
	mul.f64 	%fd2, %fd1, 0d3EB515370F99F6CB;
	cvt.rn.f32.f64	%f9, %fd2;
	mul.f32 	%f10, %f9, %f13;
	mul.f32 	%f11, %f10, %f8;
	mad.lo.s32 	%r20, %r3, %r6, %r2;
	mad.lo.s32 	%r21, %r20, %r5, %r1;
	mul.wide.s32 	%rd15, %r21, 4;
	add.s64 	%rd16, %rd11, %rd15;
	st.global.f32 	[%rd16], %f11;

BB0_6:
	ret;
}


`
)
