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

// CUDA handle for reducesum kernel
var reducesum_code cu.Function

// Stores the arguments for reducesum kernel invocation
type reducesum_args_t struct {
	arg_src     unsafe.Pointer
	arg_dst     unsafe.Pointer
	arg_initVal float32
	arg_n       int
	argptr      [4]unsafe.Pointer
	sync.Mutex
}

// Stores the arguments for reducesum kernel invocation
var reducesum_args reducesum_args_t

func init() {
	// CUDA driver kernel call wants pointers to arguments, set them up once.
	reducesum_args.argptr[0] = unsafe.Pointer(&reducesum_args.arg_src)
	reducesum_args.argptr[1] = unsafe.Pointer(&reducesum_args.arg_dst)
	reducesum_args.argptr[2] = unsafe.Pointer(&reducesum_args.arg_initVal)
	reducesum_args.argptr[3] = unsafe.Pointer(&reducesum_args.arg_n)
}

// Wrapper for reducesum CUDA kernel, asynchronous.
func k_reducesum_async(src unsafe.Pointer, dst unsafe.Pointer, initVal float32, n int, cfg *config) {
	if Synchronous { // debug
		Sync()
		timer.Start("reducesum")
	}

	reducesum_args.Lock()
	defer reducesum_args.Unlock()

	if reducesum_code == 0 {
		reducesum_code = fatbinLoad(reducesum_map, "reducesum")
	}

	reducesum_args.arg_src = src
	reducesum_args.arg_dst = dst
	reducesum_args.arg_initVal = initVal
	reducesum_args.arg_n = n

	args := reducesum_args.argptr[:]
	cu.LaunchKernel(reducesum_code, cfg.Grid.X, cfg.Grid.Y, cfg.Grid.Z, cfg.Block.X, cfg.Block.Y, cfg.Block.Z, 0, stream0, args)

	if Synchronous { // debug
		Sync()
		timer.Stop("reducesum")
	}
}

// maps compute capability on PTX code for reducesum kernel.
var reducesum_map = map[int]string{0: "",
	52: reducesum_ptx_52}

// reducesum PTX code for various compute capabilities.
const (
	reducesum_ptx_52 = `
.version 7.0
.target sm_52
.address_size 64

	// .globl	reducesum

.visible .entry reducesum(
	.param .u64 reducesum_param_0,
	.param .u64 reducesum_param_1,
	.param .f32 reducesum_param_2,
	.param .u32 reducesum_param_3
)
{
	.reg .pred 	%p<8>;
	.reg .f32 	%f<31>;
	.reg .b32 	%r<21>;
	.reg .b64 	%rd<7>;
	// demoted variable
	.shared .align 4 .b8 _ZZ9reducesumE5sdata[2048];

	ld.param.u64 	%rd3, [reducesum_param_0];
	ld.param.u64 	%rd2, [reducesum_param_1];
	ld.param.f32 	%f30, [reducesum_param_2];
	ld.param.u32 	%r10, [reducesum_param_3];
	cvta.to.global.u64 	%rd1, %rd3;
	mov.u32 	%r20, %ntid.x;
	mov.u32 	%r11, %ctaid.x;
	mov.u32 	%r2, %tid.x;
	mad.lo.s32 	%r19, %r20, %r11, %r2;
	mov.u32 	%r12, %nctaid.x;
	mul.lo.s32 	%r4, %r12, %r20;
	setp.ge.s32	%p1, %r19, %r10;
	@%p1 bra 	BB0_2;

BB0_1:
	mul.wide.s32 	%rd4, %r19, 4;
	add.s64 	%rd5, %rd1, %rd4;
	ld.global.nc.f32 	%f5, [%rd5];
	add.f32 	%f30, %f30, %f5;
	add.s32 	%r19, %r19, %r4;
	setp.lt.s32	%p2, %r19, %r10;
	@%p2 bra 	BB0_1;

BB0_2:
	shl.b32 	%r13, %r2, 2;
	mov.u32 	%r14, _ZZ9reducesumE5sdata;
	add.s32 	%r7, %r14, %r13;
	st.shared.f32 	[%r7], %f30;
	bar.sync 	0;
	setp.lt.u32	%p3, %r20, 66;
	@%p3 bra 	BB0_6;

BB0_3:
	shr.u32 	%r9, %r20, 1;
	setp.ge.u32	%p4, %r2, %r9;
	@%p4 bra 	BB0_5;

	ld.shared.f32 	%f6, [%r7];
	add.s32 	%r15, %r9, %r2;
	shl.b32 	%r16, %r15, 2;
	add.s32 	%r18, %r14, %r16;
	ld.shared.f32 	%f7, [%r18];
	add.f32 	%f8, %f6, %f7;
	st.shared.f32 	[%r7], %f8;

BB0_5:
	bar.sync 	0;
	setp.gt.u32	%p5, %r20, 131;
	mov.u32 	%r20, %r9;
	@%p5 bra 	BB0_3;

BB0_6:
	setp.gt.s32	%p6, %r2, 31;
	@%p6 bra 	BB0_8;

	ld.volatile.shared.f32 	%f9, [%r7];
	ld.volatile.shared.f32 	%f10, [%r7+128];
	add.f32 	%f11, %f9, %f10;
	st.volatile.shared.f32 	[%r7], %f11;
	ld.volatile.shared.f32 	%f12, [%r7+64];
	ld.volatile.shared.f32 	%f13, [%r7];
	add.f32 	%f14, %f13, %f12;
	st.volatile.shared.f32 	[%r7], %f14;
	ld.volatile.shared.f32 	%f15, [%r7+32];
	ld.volatile.shared.f32 	%f16, [%r7];
	add.f32 	%f17, %f16, %f15;
	st.volatile.shared.f32 	[%r7], %f17;
	ld.volatile.shared.f32 	%f18, [%r7+16];
	ld.volatile.shared.f32 	%f19, [%r7];
	add.f32 	%f20, %f19, %f18;
	st.volatile.shared.f32 	[%r7], %f20;
	ld.volatile.shared.f32 	%f21, [%r7+8];
	ld.volatile.shared.f32 	%f22, [%r7];
	add.f32 	%f23, %f22, %f21;
	st.volatile.shared.f32 	[%r7], %f23;
	ld.volatile.shared.f32 	%f24, [%r7+4];
	ld.volatile.shared.f32 	%f25, [%r7];
	add.f32 	%f26, %f25, %f24;
	st.volatile.shared.f32 	[%r7], %f26;

BB0_8:
	setp.ne.s32	%p7, %r2, 0;
	@%p7 bra 	BB0_10;

	ld.shared.f32 	%f27, [_ZZ9reducesumE5sdata];
	cvta.to.global.u64 	%rd6, %rd2;
	atom.global.add.f32 	%f28, [%rd6], %f27;

BB0_10:
	ret;
}


`
)
