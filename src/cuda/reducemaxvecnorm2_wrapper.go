package cuda

/*
 THIS FILE IS AUTO-GENERATED BY CUDA2GO.
 EDITING IS FUTILE.
*/

import(
	"unsafe"
	"github.com/MathieuMoalic/amumax/src/cuda/cu"
	"github.com/MathieuMoalic/amumax/src/engine_old/timer_old"
	"sync"
)

// CUDA handle for reducemaxvecnorm2 kernel
var reducemaxvecnorm2_code cu.Function

// Stores the arguments for reducemaxvecnorm2 kernel invocation
type reducemaxvecnorm2_args_t struct{
	 arg_x unsafe.Pointer
	 arg_y unsafe.Pointer
	 arg_z unsafe.Pointer
	 arg_dst unsafe.Pointer
	 arg_initVal float32
	 arg_n int
	 argptr [6]unsafe.Pointer
	sync.Mutex
}

// Stores the arguments for reducemaxvecnorm2 kernel invocation
var reducemaxvecnorm2_args reducemaxvecnorm2_args_t

func init(){
	// CUDA driver kernel call wants pointers to arguments, set them up once.
	 reducemaxvecnorm2_args.argptr[0] = unsafe.Pointer(&reducemaxvecnorm2_args.arg_x)
	 reducemaxvecnorm2_args.argptr[1] = unsafe.Pointer(&reducemaxvecnorm2_args.arg_y)
	 reducemaxvecnorm2_args.argptr[2] = unsafe.Pointer(&reducemaxvecnorm2_args.arg_z)
	 reducemaxvecnorm2_args.argptr[3] = unsafe.Pointer(&reducemaxvecnorm2_args.arg_dst)
	 reducemaxvecnorm2_args.argptr[4] = unsafe.Pointer(&reducemaxvecnorm2_args.arg_initVal)
	 reducemaxvecnorm2_args.argptr[5] = unsafe.Pointer(&reducemaxvecnorm2_args.arg_n)
	 }

// Wrapper for reducemaxvecnorm2 CUDA kernel, asynchronous.
func k_reducemaxvecnorm2_async ( x unsafe.Pointer, y unsafe.Pointer, z unsafe.Pointer, dst unsafe.Pointer, initVal float32, n int,  cfg *config) {
	if Synchronous{ // debug
		Sync()
		timer_old.Start("reducemaxvecnorm2")
	}

	reducemaxvecnorm2_args.Lock()
	defer reducemaxvecnorm2_args.Unlock()

	if reducemaxvecnorm2_code == 0{
		reducemaxvecnorm2_code = fatbinLoad(reducemaxvecnorm2_map, "reducemaxvecnorm2")
	}

	 reducemaxvecnorm2_args.arg_x = x
	 reducemaxvecnorm2_args.arg_y = y
	 reducemaxvecnorm2_args.arg_z = z
	 reducemaxvecnorm2_args.arg_dst = dst
	 reducemaxvecnorm2_args.arg_initVal = initVal
	 reducemaxvecnorm2_args.arg_n = n
	

	args := reducemaxvecnorm2_args.argptr[:]
	cu.LaunchKernel(reducemaxvecnorm2_code, cfg.Grid.X, cfg.Grid.Y, cfg.Grid.Z, cfg.Block.X, cfg.Block.Y, cfg.Block.Z, 0, stream0, args)

	if Synchronous{ // debug
		Sync()
		timer_old.Stop("reducemaxvecnorm2")
	}
}

// maps compute capability on PTX code for reducemaxvecnorm2 kernel.
var reducemaxvecnorm2_map = map[int]string{ 0: "" ,
52: reducemaxvecnorm2_ptx_52  }

// reducemaxvecnorm2 PTX code for various compute capabilities.
const(
  reducemaxvecnorm2_ptx_52 = `
.version 7.0
.target sm_52
.address_size 64

	// .globl	reducemaxvecnorm2

.visible .entry reducemaxvecnorm2(
	.param .u64 reducemaxvecnorm2_param_0,
	.param .u64 reducemaxvecnorm2_param_1,
	.param .u64 reducemaxvecnorm2_param_2,
	.param .u64 reducemaxvecnorm2_param_3,
	.param .f32 reducemaxvecnorm2_param_4,
	.param .u32 reducemaxvecnorm2_param_5
)
{
	.reg .pred 	%p<8>;
	.reg .f32 	%f<36>;
	.reg .b32 	%r<23>;
	.reg .b64 	%rd<13>;
	// demoted variable
	.shared .align 4 .b8 _ZZ17reducemaxvecnorm2E5sdata[2048];

	ld.param.u64 	%rd5, [reducemaxvecnorm2_param_0];
	ld.param.u64 	%rd6, [reducemaxvecnorm2_param_1];
	ld.param.u64 	%rd7, [reducemaxvecnorm2_param_2];
	ld.param.u64 	%rd4, [reducemaxvecnorm2_param_3];
	ld.param.f32 	%f35, [reducemaxvecnorm2_param_4];
	ld.param.u32 	%r10, [reducemaxvecnorm2_param_5];
	cvta.to.global.u64 	%rd1, %rd7;
	cvta.to.global.u64 	%rd2, %rd6;
	cvta.to.global.u64 	%rd3, %rd5;
	mov.u32 	%r22, %ntid.x;
	mov.u32 	%r11, %ctaid.x;
	mov.u32 	%r2, %tid.x;
	mad.lo.s32 	%r21, %r22, %r11, %r2;
	mov.u32 	%r12, %nctaid.x;
	mul.lo.s32 	%r4, %r12, %r22;
	setp.ge.s32	%p1, %r21, %r10;
	@%p1 bra 	BB0_2;

BB0_1:
	mul.wide.s32 	%rd8, %r21, 4;
	add.s64 	%rd9, %rd3, %rd8;
	ld.global.nc.f32 	%f5, [%rd9];
	add.s64 	%rd10, %rd2, %rd8;
	ld.global.nc.f32 	%f6, [%rd10];
	mul.f32 	%f7, %f6, %f6;
	fma.rn.f32 	%f8, %f5, %f5, %f7;
	add.s64 	%rd11, %rd1, %rd8;
	ld.global.nc.f32 	%f9, [%rd11];
	fma.rn.f32 	%f10, %f9, %f9, %f8;
	max.f32 	%f35, %f35, %f10;
	add.s32 	%r21, %r21, %r4;
	setp.lt.s32	%p2, %r21, %r10;
	@%p2 bra 	BB0_1;

BB0_2:
	shl.b32 	%r13, %r2, 2;
	mov.u32 	%r14, _ZZ17reducemaxvecnorm2E5sdata;
	add.s32 	%r7, %r14, %r13;
	st.shared.f32 	[%r7], %f35;
	bar.sync 	0;
	setp.lt.u32	%p3, %r22, 66;
	@%p3 bra 	BB0_6;

BB0_3:
	shr.u32 	%r9, %r22, 1;
	setp.ge.u32	%p4, %r2, %r9;
	@%p4 bra 	BB0_5;

	ld.shared.f32 	%f11, [%r7];
	add.s32 	%r15, %r9, %r2;
	shl.b32 	%r16, %r15, 2;
	add.s32 	%r18, %r14, %r16;
	ld.shared.f32 	%f12, [%r18];
	max.f32 	%f13, %f11, %f12;
	st.shared.f32 	[%r7], %f13;

BB0_5:
	bar.sync 	0;
	setp.gt.u32	%p5, %r22, 131;
	mov.u32 	%r22, %r9;
	@%p5 bra 	BB0_3;

BB0_6:
	setp.gt.s32	%p6, %r2, 31;
	@%p6 bra 	BB0_8;

	ld.volatile.shared.f32 	%f14, [%r7];
	ld.volatile.shared.f32 	%f15, [%r7+128];
	max.f32 	%f16, %f14, %f15;
	st.volatile.shared.f32 	[%r7], %f16;
	ld.volatile.shared.f32 	%f17, [%r7+64];
	ld.volatile.shared.f32 	%f18, [%r7];
	max.f32 	%f19, %f18, %f17;
	st.volatile.shared.f32 	[%r7], %f19;
	ld.volatile.shared.f32 	%f20, [%r7+32];
	ld.volatile.shared.f32 	%f21, [%r7];
	max.f32 	%f22, %f21, %f20;
	st.volatile.shared.f32 	[%r7], %f22;
	ld.volatile.shared.f32 	%f23, [%r7+16];
	ld.volatile.shared.f32 	%f24, [%r7];
	max.f32 	%f25, %f24, %f23;
	st.volatile.shared.f32 	[%r7], %f25;
	ld.volatile.shared.f32 	%f26, [%r7+8];
	ld.volatile.shared.f32 	%f27, [%r7];
	max.f32 	%f28, %f27, %f26;
	st.volatile.shared.f32 	[%r7], %f28;
	ld.volatile.shared.f32 	%f29, [%r7+4];
	ld.volatile.shared.f32 	%f30, [%r7];
	max.f32 	%f31, %f30, %f29;
	st.volatile.shared.f32 	[%r7], %f31;

BB0_8:
	setp.ne.s32	%p7, %r2, 0;
	@%p7 bra 	BB0_10;

	ld.shared.f32 	%f32, [_ZZ17reducemaxvecnorm2E5sdata];
	abs.f32 	%f33, %f32;
	mov.b32 	 %r19, %f33;
	cvta.to.global.u64 	%rd12, %rd4;
	atom.global.max.s32 	%r20, [%rd12], %r19;

BB0_10:
	ret;
}


`
 )
