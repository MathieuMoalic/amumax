package cuda

/*
 THIS FILE IS AUTO-GENERATED BY CUDA2GO.
 EDITING IS FUTILE.
*/

import(
	"unsafe"
	"github.com/MathieuMoalic/amumax/src/cuda/cu"
	"github.com/MathieuMoalic/amumax/src/timer"
	"sync"
)

// CUDA handle for shiftedgecarryY kernel
var shiftedgecarryY_code cu.Function

// Stores the arguments for shiftedgecarryY kernel invocation
type shiftedgecarryY_args_t struct{
	 arg_dst unsafe.Pointer
	 arg_src unsafe.Pointer
	 arg_othercomp unsafe.Pointer
	 arg_anothercomp unsafe.Pointer
	 arg_Nx int
	 arg_Ny int
	 arg_Nz int
	 arg_shy int
	 arg_clampD float32
	 arg_clampU float32
	 argptr [10]unsafe.Pointer
	sync.Mutex
}

// Stores the arguments for shiftedgecarryY kernel invocation
var shiftedgecarryY_args shiftedgecarryY_args_t

func init(){
	// CUDA driver kernel call wants pointers to arguments, set them up once.
	 shiftedgecarryY_args.argptr[0] = unsafe.Pointer(&shiftedgecarryY_args.arg_dst)
	 shiftedgecarryY_args.argptr[1] = unsafe.Pointer(&shiftedgecarryY_args.arg_src)
	 shiftedgecarryY_args.argptr[2] = unsafe.Pointer(&shiftedgecarryY_args.arg_othercomp)
	 shiftedgecarryY_args.argptr[3] = unsafe.Pointer(&shiftedgecarryY_args.arg_anothercomp)
	 shiftedgecarryY_args.argptr[4] = unsafe.Pointer(&shiftedgecarryY_args.arg_Nx)
	 shiftedgecarryY_args.argptr[5] = unsafe.Pointer(&shiftedgecarryY_args.arg_Ny)
	 shiftedgecarryY_args.argptr[6] = unsafe.Pointer(&shiftedgecarryY_args.arg_Nz)
	 shiftedgecarryY_args.argptr[7] = unsafe.Pointer(&shiftedgecarryY_args.arg_shy)
	 shiftedgecarryY_args.argptr[8] = unsafe.Pointer(&shiftedgecarryY_args.arg_clampD)
	 shiftedgecarryY_args.argptr[9] = unsafe.Pointer(&shiftedgecarryY_args.arg_clampU)
	 }

// Wrapper for shiftedgecarryY CUDA kernel, asynchronous.
func k_shiftedgecarryY_async ( dst unsafe.Pointer, src unsafe.Pointer, othercomp unsafe.Pointer, anothercomp unsafe.Pointer, Nx int, Ny int, Nz int, shy int, clampD float32, clampU float32,  cfg *config) {
	if Synchronous{ // debug
		Sync()
		timer.Start("shiftedgecarryY")
	}

	shiftedgecarryY_args.Lock()
	defer shiftedgecarryY_args.Unlock()

	if shiftedgecarryY_code == 0{
		shiftedgecarryY_code = fatbinLoad(shiftedgecarryY_map, "shiftedgecarryY")
	}

	 shiftedgecarryY_args.arg_dst = dst
	 shiftedgecarryY_args.arg_src = src
	 shiftedgecarryY_args.arg_othercomp = othercomp
	 shiftedgecarryY_args.arg_anothercomp = anothercomp
	 shiftedgecarryY_args.arg_Nx = Nx
	 shiftedgecarryY_args.arg_Ny = Ny
	 shiftedgecarryY_args.arg_Nz = Nz
	 shiftedgecarryY_args.arg_shy = shy
	 shiftedgecarryY_args.arg_clampD = clampD
	 shiftedgecarryY_args.arg_clampU = clampU
	

	args := shiftedgecarryY_args.argptr[:]
	cu.LaunchKernel(shiftedgecarryY_code, cfg.Grid.X, cfg.Grid.Y, cfg.Grid.Z, cfg.Block.X, cfg.Block.Y, cfg.Block.Z, 0, stream0, args)

	if Synchronous{ // debug
		Sync()
		timer.Stop("shiftedgecarryY")
	}
}

// maps compute capability on PTX code for shiftedgecarryY kernel.
var shiftedgecarryY_map = map[int]string{ 0: "" ,
52: shiftedgecarryY_ptx_52  }

// shiftedgecarryY PTX code for various compute capabilities.
const(
  shiftedgecarryY_ptx_52 = `
.version 7.0
.target sm_52
.address_size 64

	// .globl	shiftedgecarryY

.visible .entry shiftedgecarryY(
	.param .u64 shiftedgecarryY_param_0,
	.param .u64 shiftedgecarryY_param_1,
	.param .u64 shiftedgecarryY_param_2,
	.param .u64 shiftedgecarryY_param_3,
	.param .u32 shiftedgecarryY_param_4,
	.param .u32 shiftedgecarryY_param_5,
	.param .u32 shiftedgecarryY_param_6,
	.param .u32 shiftedgecarryY_param_7,
	.param .f32 shiftedgecarryY_param_8,
	.param .f32 shiftedgecarryY_param_9
)
{
	.reg .pred 	%p<14>;
	.reg .f32 	%f<14>;
	.reg .b32 	%r<27>;
	.reg .b64 	%rd<25>;


	ld.param.u64 	%rd4, [shiftedgecarryY_param_0];
	ld.param.u64 	%rd5, [shiftedgecarryY_param_1];
	ld.param.u64 	%rd6, [shiftedgecarryY_param_2];
	ld.param.u64 	%rd7, [shiftedgecarryY_param_3];
	ld.param.u32 	%r8, [shiftedgecarryY_param_4];
	ld.param.u32 	%r9, [shiftedgecarryY_param_5];
	ld.param.u32 	%r11, [shiftedgecarryY_param_6];
	ld.param.u32 	%r10, [shiftedgecarryY_param_7];
	ld.param.f32 	%f7, [shiftedgecarryY_param_8];
	ld.param.f32 	%f8, [shiftedgecarryY_param_9];
	cvta.to.global.u64 	%rd1, %rd7;
	cvta.to.global.u64 	%rd2, %rd6;
	cvta.to.global.u64 	%rd3, %rd5;
	mov.u32 	%r12, %ntid.x;
	mov.u32 	%r13, %ctaid.x;
	mov.u32 	%r14, %tid.x;
	mad.lo.s32 	%r1, %r12, %r13, %r14;
	mov.u32 	%r15, %ntid.y;
	mov.u32 	%r16, %ctaid.y;
	mov.u32 	%r17, %tid.y;
	mad.lo.s32 	%r2, %r15, %r16, %r17;
	mov.u32 	%r18, %ntid.z;
	mov.u32 	%r19, %ctaid.z;
	mov.u32 	%r20, %tid.z;
	mad.lo.s32 	%r3, %r18, %r19, %r20;
	setp.lt.s32	%p1, %r1, %r8;
	setp.lt.s32	%p2, %r2, %r9;
	and.pred  	%p3, %p1, %p2;
	setp.lt.s32	%p4, %r3, %r11;
	and.pred  	%p5, %p3, %p4;
	@!%p5 bra 	BB0_11;
	bra.uni 	BB0_1;

BB0_1:
	sub.s32 	%r4, %r2, %r10;
	setp.lt.s32	%p6, %r4, 0;
	mul.lo.s32 	%r5, %r3, %r9;
	@%p6 bra 	BB0_7;

	setp.lt.s32	%p7, %r4, %r9;
	@%p7 bra 	BB0_6;
	bra.uni 	BB0_3;

BB0_6:
	add.s32 	%r23, %r5, %r4;
	mad.lo.s32 	%r24, %r23, %r8, %r1;
	mul.wide.s32 	%rd14, %r24, 4;
	add.s64 	%rd15, %rd3, %rd14;
	ld.global.nc.f32 	%f13, [%rd15];
	bra.uni 	BB0_10;

BB0_7:
	mad.lo.s32 	%r7, %r5, %r8, %r1;
	mul.wide.s32 	%rd16, %r7, 4;
	add.s64 	%rd17, %rd3, %rd16;
	ld.global.nc.f32 	%f13, [%rd17];
	setp.neu.f32	%p11, %f13, 0f00000000;
	@%p11 bra 	BB0_10;

	add.s64 	%rd19, %rd2, %rd16;
	ld.global.nc.f32 	%f11, [%rd19];
	setp.neu.f32	%p12, %f11, 0f00000000;
	@%p12 bra 	BB0_10;

	add.s64 	%rd21, %rd1, %rd16;
	ld.global.nc.f32 	%f12, [%rd21];
	setp.eq.f32	%p13, %f12, 0f00000000;
	selp.f32	%f13, %f7, %f13, %p13;
	bra.uni 	BB0_10;

BB0_3:
	add.s32 	%r21, %r9, %r5;
	add.s32 	%r22, %r21, -1;
	mad.lo.s32 	%r6, %r22, %r8, %r1;
	mul.wide.s32 	%rd8, %r6, 4;
	add.s64 	%rd9, %rd3, %rd8;
	ld.global.nc.f32 	%f13, [%rd9];
	setp.neu.f32	%p8, %f13, 0f00000000;
	@%p8 bra 	BB0_10;

	add.s64 	%rd11, %rd2, %rd8;
	ld.global.nc.f32 	%f9, [%rd11];
	setp.neu.f32	%p9, %f9, 0f00000000;
	@%p9 bra 	BB0_10;

	add.s64 	%rd13, %rd1, %rd8;
	ld.global.nc.f32 	%f10, [%rd13];
	setp.eq.f32	%p10, %f10, 0f00000000;
	selp.f32	%f13, %f8, %f13, %p10;

BB0_10:
	cvta.to.global.u64 	%rd22, %rd4;
	add.s32 	%r25, %r5, %r2;
	mad.lo.s32 	%r26, %r25, %r8, %r1;
	mul.wide.s32 	%rd23, %r26, 4;
	add.s64 	%rd24, %rd22, %rd23;
	st.global.f32 	[%rd24], %f13;

BB0_11:
	ret;
}


`
 )
