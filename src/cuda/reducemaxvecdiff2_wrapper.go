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

// CUDA handle for reducemaxvecdiff2 kernel
var reducemaxvecdiff2_code cu.Function

// Stores the arguments for reducemaxvecdiff2 kernel invocation
type reducemaxvecdiff2_args_t struct{
	 arg_x1 unsafe.Pointer
	 arg_y1 unsafe.Pointer
	 arg_z1 unsafe.Pointer
	 arg_x2 unsafe.Pointer
	 arg_y2 unsafe.Pointer
	 arg_z2 unsafe.Pointer
	 arg_dst unsafe.Pointer
	 arg_initVal float32
	 arg_n int
	 argptr [9]unsafe.Pointer
	sync.Mutex
}

// Stores the arguments for reducemaxvecdiff2 kernel invocation
var reducemaxvecdiff2_args reducemaxvecdiff2_args_t

func init(){
	// CUDA driver kernel call wants pointers to arguments, set them up once.
	 reducemaxvecdiff2_args.argptr[0] = unsafe.Pointer(&reducemaxvecdiff2_args.arg_x1)
	 reducemaxvecdiff2_args.argptr[1] = unsafe.Pointer(&reducemaxvecdiff2_args.arg_y1)
	 reducemaxvecdiff2_args.argptr[2] = unsafe.Pointer(&reducemaxvecdiff2_args.arg_z1)
	 reducemaxvecdiff2_args.argptr[3] = unsafe.Pointer(&reducemaxvecdiff2_args.arg_x2)
	 reducemaxvecdiff2_args.argptr[4] = unsafe.Pointer(&reducemaxvecdiff2_args.arg_y2)
	 reducemaxvecdiff2_args.argptr[5] = unsafe.Pointer(&reducemaxvecdiff2_args.arg_z2)
	 reducemaxvecdiff2_args.argptr[6] = unsafe.Pointer(&reducemaxvecdiff2_args.arg_dst)
	 reducemaxvecdiff2_args.argptr[7] = unsafe.Pointer(&reducemaxvecdiff2_args.arg_initVal)
	 reducemaxvecdiff2_args.argptr[8] = unsafe.Pointer(&reducemaxvecdiff2_args.arg_n)
	 }

// Wrapper for reducemaxvecdiff2 CUDA kernel, asynchronous.
func k_reducemaxvecdiff2_async ( x1 unsafe.Pointer, y1 unsafe.Pointer, z1 unsafe.Pointer, x2 unsafe.Pointer, y2 unsafe.Pointer, z2 unsafe.Pointer, dst unsafe.Pointer, initVal float32, n int,  cfg *config) {
	if Synchronous{ // debug
		Sync()
		timer_old.Start("reducemaxvecdiff2")
	}

	reducemaxvecdiff2_args.Lock()
	defer reducemaxvecdiff2_args.Unlock()

	if reducemaxvecdiff2_code == 0{
		reducemaxvecdiff2_code = fatbinLoad(reducemaxvecdiff2_map, "reducemaxvecdiff2")
	}

	 reducemaxvecdiff2_args.arg_x1 = x1
	 reducemaxvecdiff2_args.arg_y1 = y1
	 reducemaxvecdiff2_args.arg_z1 = z1
	 reducemaxvecdiff2_args.arg_x2 = x2
	 reducemaxvecdiff2_args.arg_y2 = y2
	 reducemaxvecdiff2_args.arg_z2 = z2
	 reducemaxvecdiff2_args.arg_dst = dst
	 reducemaxvecdiff2_args.arg_initVal = initVal
	 reducemaxvecdiff2_args.arg_n = n
	

	args := reducemaxvecdiff2_args.argptr[:]
	cu.LaunchKernel(reducemaxvecdiff2_code, cfg.Grid.X, cfg.Grid.Y, cfg.Grid.Z, cfg.Block.X, cfg.Block.Y, cfg.Block.Z, 0, stream0, args)

	if Synchronous{ // debug
		Sync()
		timer_old.Stop("reducemaxvecdiff2")
	}
}

// maps compute capability on PTX code for reducemaxvecdiff2 kernel.
var reducemaxvecdiff2_map = map[int]string{ 0: "" ,
52: reducemaxvecdiff2_ptx_52  }

// reducemaxvecdiff2 PTX code for various compute capabilities.
const(
  reducemaxvecdiff2_ptx_52 = `
.version 7.0
.target sm_52
.address_size 64

	// .globl	reducemaxvecdiff2

.visible .entry reducemaxvecdiff2(
	.param .u64 reducemaxvecdiff2_param_0,
	.param .u64 reducemaxvecdiff2_param_1,
	.param .u64 reducemaxvecdiff2_param_2,
	.param .u64 reducemaxvecdiff2_param_3,
	.param .u64 reducemaxvecdiff2_param_4,
	.param .u64 reducemaxvecdiff2_param_5,
	.param .u64 reducemaxvecdiff2_param_6,
	.param .f32 reducemaxvecdiff2_param_7,
	.param .u32 reducemaxvecdiff2_param_8
)
{
	.reg .pred 	%p<8>;
	.reg .f32 	%f<42>;
	.reg .b32 	%r<23>;
	.reg .b64 	%rd<22>;
	// demoted variable
	.shared .align 4 .b8 _ZZ17reducemaxvecdiff2E5sdata[2048];

	ld.param.u64 	%rd8, [reducemaxvecdiff2_param_0];
	ld.param.u64 	%rd9, [reducemaxvecdiff2_param_1];
	ld.param.u64 	%rd10, [reducemaxvecdiff2_param_2];
	ld.param.u64 	%rd6, [reducemaxvecdiff2_param_3];
	ld.param.u64 	%rd11, [reducemaxvecdiff2_param_4];
	ld.param.u64 	%rd12, [reducemaxvecdiff2_param_5];
	ld.param.u64 	%rd7, [reducemaxvecdiff2_param_6];
	ld.param.f32 	%f41, [reducemaxvecdiff2_param_7];
	ld.param.u32 	%r10, [reducemaxvecdiff2_param_8];
	cvta.to.global.u64 	%rd1, %rd12;
	cvta.to.global.u64 	%rd2, %rd10;
	cvta.to.global.u64 	%rd3, %rd11;
	cvta.to.global.u64 	%rd4, %rd9;
	cvta.to.global.u64 	%rd5, %rd8;
	mov.u32 	%r22, %ntid.x;
	mov.u32 	%r11, %ctaid.x;
	mov.u32 	%r2, %tid.x;
	mad.lo.s32 	%r21, %r22, %r11, %r2;
	mov.u32 	%r12, %nctaid.x;
	mul.lo.s32 	%r4, %r12, %r22;
	setp.ge.s32	%p1, %r21, %r10;
	@%p1 bra 	BB0_3;

	cvta.to.global.u64 	%rd15, %rd6;

BB0_2:
	mul.wide.s32 	%rd13, %r21, 4;
	add.s64 	%rd14, %rd5, %rd13;
	add.s64 	%rd16, %rd15, %rd13;
	ld.global.nc.f32 	%f5, [%rd16];
	ld.global.nc.f32 	%f6, [%rd14];
	sub.f32 	%f7, %f6, %f5;
	add.s64 	%rd17, %rd4, %rd13;
	add.s64 	%rd18, %rd3, %rd13;
	ld.global.nc.f32 	%f8, [%rd18];
	ld.global.nc.f32 	%f9, [%rd17];
	sub.f32 	%f10, %f9, %f8;
	mul.f32 	%f11, %f10, %f10;
	fma.rn.f32 	%f12, %f7, %f7, %f11;
	add.s64 	%rd19, %rd2, %rd13;
	add.s64 	%rd20, %rd1, %rd13;
	ld.global.nc.f32 	%f13, [%rd20];
	ld.global.nc.f32 	%f14, [%rd19];
	sub.f32 	%f15, %f14, %f13;
	fma.rn.f32 	%f16, %f15, %f15, %f12;
	max.f32 	%f41, %f41, %f16;
	add.s32 	%r21, %r21, %r4;
	setp.lt.s32	%p2, %r21, %r10;
	@%p2 bra 	BB0_2;

BB0_3:
	shl.b32 	%r13, %r2, 2;
	mov.u32 	%r14, _ZZ17reducemaxvecdiff2E5sdata;
	add.s32 	%r7, %r14, %r13;
	st.shared.f32 	[%r7], %f41;
	bar.sync 	0;
	setp.lt.u32	%p3, %r22, 66;
	@%p3 bra 	BB0_7;

BB0_4:
	shr.u32 	%r9, %r22, 1;
	setp.ge.u32	%p4, %r2, %r9;
	@%p4 bra 	BB0_6;

	ld.shared.f32 	%f17, [%r7];
	add.s32 	%r15, %r9, %r2;
	shl.b32 	%r16, %r15, 2;
	add.s32 	%r18, %r14, %r16;
	ld.shared.f32 	%f18, [%r18];
	max.f32 	%f19, %f17, %f18;
	st.shared.f32 	[%r7], %f19;

BB0_6:
	bar.sync 	0;
	setp.gt.u32	%p5, %r22, 131;
	mov.u32 	%r22, %r9;
	@%p5 bra 	BB0_4;

BB0_7:
	setp.gt.s32	%p6, %r2, 31;
	@%p6 bra 	BB0_9;

	ld.volatile.shared.f32 	%f20, [%r7];
	ld.volatile.shared.f32 	%f21, [%r7+128];
	max.f32 	%f22, %f20, %f21;
	st.volatile.shared.f32 	[%r7], %f22;
	ld.volatile.shared.f32 	%f23, [%r7+64];
	ld.volatile.shared.f32 	%f24, [%r7];
	max.f32 	%f25, %f24, %f23;
	st.volatile.shared.f32 	[%r7], %f25;
	ld.volatile.shared.f32 	%f26, [%r7+32];
	ld.volatile.shared.f32 	%f27, [%r7];
	max.f32 	%f28, %f27, %f26;
	st.volatile.shared.f32 	[%r7], %f28;
	ld.volatile.shared.f32 	%f29, [%r7+16];
	ld.volatile.shared.f32 	%f30, [%r7];
	max.f32 	%f31, %f30, %f29;
	st.volatile.shared.f32 	[%r7], %f31;
	ld.volatile.shared.f32 	%f32, [%r7+8];
	ld.volatile.shared.f32 	%f33, [%r7];
	max.f32 	%f34, %f33, %f32;
	st.volatile.shared.f32 	[%r7], %f34;
	ld.volatile.shared.f32 	%f35, [%r7+4];
	ld.volatile.shared.f32 	%f36, [%r7];
	max.f32 	%f37, %f36, %f35;
	st.volatile.shared.f32 	[%r7], %f37;

BB0_9:
	setp.ne.s32	%p7, %r2, 0;
	@%p7 bra 	BB0_11;

	ld.shared.f32 	%f38, [_ZZ17reducemaxvecdiff2E5sdata];
	abs.f32 	%f39, %f38;
	mov.b32 	 %r19, %f39;
	cvta.to.global.u64 	%rd21, %rd7;
	atom.global.max.s32 	%r20, [%rd21], %r19;

BB0_11:
	ret;
}


`
 )
