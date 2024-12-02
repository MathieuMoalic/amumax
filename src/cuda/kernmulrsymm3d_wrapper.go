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

// CUDA handle for kernmulRSymm3D kernel
var kernmulRSymm3D_code cu.Function

// Stores the arguments for kernmulRSymm3D kernel invocation
type kernmulRSymm3D_args_t struct {
	arg_fftMx  unsafe.Pointer
	arg_fftMy  unsafe.Pointer
	arg_fftMz  unsafe.Pointer
	arg_fftKxx unsafe.Pointer
	arg_fftKyy unsafe.Pointer
	arg_fftKzz unsafe.Pointer
	arg_fftKyz unsafe.Pointer
	arg_fftKxz unsafe.Pointer
	arg_fftKxy unsafe.Pointer
	arg_Nx     int
	arg_Ny     int
	arg_Nz     int
	argptr     [12]unsafe.Pointer
	sync.Mutex
}

// Stores the arguments for kernmulRSymm3D kernel invocation
var kernmulRSymm3D_args kernmulRSymm3D_args_t

func init() {
	// CUDA driver kernel call wants pointers to arguments, set them up once.
	kernmulRSymm3D_args.argptr[0] = unsafe.Pointer(&kernmulRSymm3D_args.arg_fftMx)
	kernmulRSymm3D_args.argptr[1] = unsafe.Pointer(&kernmulRSymm3D_args.arg_fftMy)
	kernmulRSymm3D_args.argptr[2] = unsafe.Pointer(&kernmulRSymm3D_args.arg_fftMz)
	kernmulRSymm3D_args.argptr[3] = unsafe.Pointer(&kernmulRSymm3D_args.arg_fftKxx)
	kernmulRSymm3D_args.argptr[4] = unsafe.Pointer(&kernmulRSymm3D_args.arg_fftKyy)
	kernmulRSymm3D_args.argptr[5] = unsafe.Pointer(&kernmulRSymm3D_args.arg_fftKzz)
	kernmulRSymm3D_args.argptr[6] = unsafe.Pointer(&kernmulRSymm3D_args.arg_fftKyz)
	kernmulRSymm3D_args.argptr[7] = unsafe.Pointer(&kernmulRSymm3D_args.arg_fftKxz)
	kernmulRSymm3D_args.argptr[8] = unsafe.Pointer(&kernmulRSymm3D_args.arg_fftKxy)
	kernmulRSymm3D_args.argptr[9] = unsafe.Pointer(&kernmulRSymm3D_args.arg_Nx)
	kernmulRSymm3D_args.argptr[10] = unsafe.Pointer(&kernmulRSymm3D_args.arg_Ny)
	kernmulRSymm3D_args.argptr[11] = unsafe.Pointer(&kernmulRSymm3D_args.arg_Nz)
}

// Wrapper for kernmulRSymm3D CUDA kernel, asynchronous.
func k_kernmulRSymm3D_async(fftMx unsafe.Pointer, fftMy unsafe.Pointer, fftMz unsafe.Pointer, fftKxx unsafe.Pointer, fftKyy unsafe.Pointer, fftKzz unsafe.Pointer, fftKyz unsafe.Pointer, fftKxz unsafe.Pointer, fftKxy unsafe.Pointer, Nx int, Ny int, Nz int, cfg *config) {
	if Synchronous { // debug
		Sync()
		timer_old.Start("kernmulRSymm3D")
	}

	kernmulRSymm3D_args.Lock()
	defer kernmulRSymm3D_args.Unlock()

	if kernmulRSymm3D_code == 0 {
		kernmulRSymm3D_code = fatbinLoad(kernmulRSymm3D_map, "kernmulRSymm3D")
	}

	kernmulRSymm3D_args.arg_fftMx = fftMx
	kernmulRSymm3D_args.arg_fftMy = fftMy
	kernmulRSymm3D_args.arg_fftMz = fftMz
	kernmulRSymm3D_args.arg_fftKxx = fftKxx
	kernmulRSymm3D_args.arg_fftKyy = fftKyy
	kernmulRSymm3D_args.arg_fftKzz = fftKzz
	kernmulRSymm3D_args.arg_fftKyz = fftKyz
	kernmulRSymm3D_args.arg_fftKxz = fftKxz
	kernmulRSymm3D_args.arg_fftKxy = fftKxy
	kernmulRSymm3D_args.arg_Nx = Nx
	kernmulRSymm3D_args.arg_Ny = Ny
	kernmulRSymm3D_args.arg_Nz = Nz

	args := kernmulRSymm3D_args.argptr[:]
	cu.LaunchKernel(kernmulRSymm3D_code, cfg.Grid.X, cfg.Grid.Y, cfg.Grid.Z, cfg.Block.X, cfg.Block.Y, cfg.Block.Z, 0, stream0, args)

	if Synchronous { // debug
		Sync()
		timer_old.Stop("kernmulRSymm3D")
	}
}

// maps compute capability on PTX code for kernmulRSymm3D kernel.
var kernmulRSymm3D_map = map[int]string{0: "",
	52: kernmulRSymm3D_ptx_52}

// kernmulRSymm3D PTX code for various compute capabilities.
const (
	kernmulRSymm3D_ptx_52 = `
.version 7.0
.target sm_52
.address_size 64

	// .globl	kernmulRSymm3D

.visible .entry kernmulRSymm3D(
	.param .u64 kernmulRSymm3D_param_0,
	.param .u64 kernmulRSymm3D_param_1,
	.param .u64 kernmulRSymm3D_param_2,
	.param .u64 kernmulRSymm3D_param_3,
	.param .u64 kernmulRSymm3D_param_4,
	.param .u64 kernmulRSymm3D_param_5,
	.param .u64 kernmulRSymm3D_param_6,
	.param .u64 kernmulRSymm3D_param_7,
	.param .u64 kernmulRSymm3D_param_8,
	.param .u32 kernmulRSymm3D_param_9,
	.param .u32 kernmulRSymm3D_param_10,
	.param .u32 kernmulRSymm3D_param_11
)
{
	.reg .pred 	%p<8>;
	.reg .f32 	%f<38>;
	.reg .b32 	%r<32>;
	.reg .b64 	%rd<30>;


	ld.param.u64 	%rd1, [kernmulRSymm3D_param_0];
	ld.param.u64 	%rd2, [kernmulRSymm3D_param_1];
	ld.param.u64 	%rd3, [kernmulRSymm3D_param_2];
	ld.param.u64 	%rd4, [kernmulRSymm3D_param_3];
	ld.param.u64 	%rd5, [kernmulRSymm3D_param_4];
	ld.param.u64 	%rd6, [kernmulRSymm3D_param_5];
	ld.param.u64 	%rd7, [kernmulRSymm3D_param_6];
	ld.param.u64 	%rd8, [kernmulRSymm3D_param_7];
	ld.param.u64 	%rd9, [kernmulRSymm3D_param_8];
	ld.param.u32 	%r4, [kernmulRSymm3D_param_9];
	ld.param.u32 	%r5, [kernmulRSymm3D_param_10];
	ld.param.u32 	%r6, [kernmulRSymm3D_param_11];
	mov.u32 	%r7, %ntid.x;
	mov.u32 	%r8, %ctaid.x;
	mov.u32 	%r9, %tid.x;
	mad.lo.s32 	%r1, %r7, %r8, %r9;
	mov.u32 	%r10, %ntid.y;
	mov.u32 	%r11, %ctaid.y;
	mov.u32 	%r12, %tid.y;
	mad.lo.s32 	%r2, %r10, %r11, %r12;
	mov.u32 	%r13, %ntid.z;
	mov.u32 	%r14, %ctaid.z;
	mov.u32 	%r15, %tid.z;
	mad.lo.s32 	%r3, %r13, %r14, %r15;
	setp.ge.s32	%p1, %r2, %r5;
	setp.ge.s32	%p2, %r1, %r4;
	or.pred  	%p3, %p1, %p2;
	setp.ge.s32	%p4, %r3, %r6;
	or.pred  	%p5, %p3, %p4;
	@%p5 bra 	BB0_2;

	cvta.to.global.u64 	%rd10, %rd3;
	cvta.to.global.u64 	%rd11, %rd2;
	cvta.to.global.u64 	%rd12, %rd1;
	cvta.to.global.u64 	%rd13, %rd4;
	mad.lo.s32 	%r16, %r3, %r5, %r2;
	mad.lo.s32 	%r17, %r16, %r4, %r1;
	shl.b32 	%r18, %r17, 1;
	mul.wide.s32 	%rd14, %r18, 4;
	add.s64 	%rd15, %rd12, %rd14;
	ld.global.f32 	%f1, [%rd15+4];
	add.s64 	%rd16, %rd11, %rd14;
	ld.global.f32 	%f2, [%rd16+4];
	add.s64 	%rd17, %rd10, %rd14;
	ld.global.f32 	%f3, [%rd17+4];
	shr.u32 	%r19, %r5, 31;
	add.s32 	%r20, %r5, %r19;
	shr.s32 	%r21, %r20, 1;
	setp.gt.s32	%p6, %r2, %r21;
	sub.s32 	%r22, %r5, %r2;
	selp.b32	%r23, %r22, %r2, %p6;
	selp.f32	%f4, 0fBF800000, 0f3F800000, %p6;
	shr.u32 	%r24, %r6, 31;
	add.s32 	%r25, %r6, %r24;
	shr.s32 	%r26, %r25, 1;
	setp.gt.s32	%p7, %r3, %r26;
	neg.f32 	%f5, %f4;
	sub.s32 	%r27, %r6, %r3;
	selp.b32	%r28, %r27, %r3, %p7;
	selp.f32	%f6, %f5, %f4, %p7;
	selp.f32	%f7, 0fBF800000, 0f3F800000, %p7;
	add.s32 	%r29, %r21, 1;
	mad.lo.s32 	%r30, %r28, %r29, %r23;
	mad.lo.s32 	%r31, %r30, %r4, %r1;
	mul.wide.s32 	%rd18, %r31, 4;
	add.s64 	%rd19, %rd13, %rd18;
	cvta.to.global.u64 	%rd20, %rd5;
	add.s64 	%rd21, %rd20, %rd18;
	ld.global.nc.f32 	%f8, [%rd21];
	cvta.to.global.u64 	%rd22, %rd6;
	add.s64 	%rd23, %rd22, %rd18;
	ld.global.nc.f32 	%f9, [%rd23];
	cvta.to.global.u64 	%rd24, %rd7;
	add.s64 	%rd25, %rd24, %rd18;
	ld.global.nc.f32 	%f10, [%rd25];
	mul.f32 	%f11, %f6, %f10;
	cvta.to.global.u64 	%rd26, %rd8;
	add.s64 	%rd27, %rd26, %rd18;
	ld.global.nc.f32 	%f12, [%rd27];
	mul.f32 	%f13, %f7, %f12;
	cvta.to.global.u64 	%rd28, %rd9;
	add.s64 	%rd29, %rd28, %rd18;
	ld.global.nc.f32 	%f14, [%rd29];
	mul.f32 	%f15, %f4, %f14;
	ld.global.nc.f32 	%f16, [%rd19];
	ld.global.f32 	%f17, [%rd15];
	ld.global.f32 	%f18, [%rd16];
	mul.f32 	%f19, %f18, %f15;
	fma.rn.f32 	%f20, %f17, %f16, %f19;
	ld.global.f32 	%f21, [%rd17];
	fma.rn.f32 	%f22, %f21, %f13, %f20;
	st.global.f32 	[%rd15], %f22;
	mul.f32 	%f23, %f2, %f15;
	fma.rn.f32 	%f24, %f1, %f16, %f23;
	fma.rn.f32 	%f25, %f3, %f13, %f24;
	st.global.f32 	[%rd15+4], %f25;
	mul.f32 	%f26, %f17, %f15;
	fma.rn.f32 	%f27, %f18, %f8, %f26;
	fma.rn.f32 	%f28, %f21, %f11, %f27;
	st.global.f32 	[%rd16], %f28;
	mul.f32 	%f29, %f1, %f15;
	fma.rn.f32 	%f30, %f2, %f8, %f29;
	fma.rn.f32 	%f31, %f3, %f11, %f30;
	st.global.f32 	[%rd16+4], %f31;
	mul.f32 	%f32, %f17, %f13;
	fma.rn.f32 	%f33, %f18, %f11, %f32;
	fma.rn.f32 	%f34, %f21, %f9, %f33;
	st.global.f32 	[%rd17], %f34;
	mul.f32 	%f35, %f1, %f13;
	fma.rn.f32 	%f36, %f2, %f11, %f35;
	fma.rn.f32 	%f37, %f3, %f9, %f36;
	st.global.f32 	[%rd17+4], %f37;

BB0_2:
	ret;
}


`
)
