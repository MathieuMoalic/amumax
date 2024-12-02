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

// CUDA handle for settopologicalchargelattice kernel
var settopologicalchargelattice_code cu.Function

// Stores the arguments for settopologicalchargelattice kernel invocation
type settopologicalchargelattice_args_t struct {
	arg_s     unsafe.Pointer
	arg_mx    unsafe.Pointer
	arg_my    unsafe.Pointer
	arg_mz    unsafe.Pointer
	arg_icxcy float32
	arg_Nx    int
	arg_Ny    int
	arg_Nz    int
	arg_PBC   byte
	argptr    [9]unsafe.Pointer
	sync.Mutex
}

// Stores the arguments for settopologicalchargelattice kernel invocation
var settopologicalchargelattice_args settopologicalchargelattice_args_t

func init() {
	// CUDA driver kernel call wants pointers to arguments, set them up once.
	settopologicalchargelattice_args.argptr[0] = unsafe.Pointer(&settopologicalchargelattice_args.arg_s)
	settopologicalchargelattice_args.argptr[1] = unsafe.Pointer(&settopologicalchargelattice_args.arg_mx)
	settopologicalchargelattice_args.argptr[2] = unsafe.Pointer(&settopologicalchargelattice_args.arg_my)
	settopologicalchargelattice_args.argptr[3] = unsafe.Pointer(&settopologicalchargelattice_args.arg_mz)
	settopologicalchargelattice_args.argptr[4] = unsafe.Pointer(&settopologicalchargelattice_args.arg_icxcy)
	settopologicalchargelattice_args.argptr[5] = unsafe.Pointer(&settopologicalchargelattice_args.arg_Nx)
	settopologicalchargelattice_args.argptr[6] = unsafe.Pointer(&settopologicalchargelattice_args.arg_Ny)
	settopologicalchargelattice_args.argptr[7] = unsafe.Pointer(&settopologicalchargelattice_args.arg_Nz)
	settopologicalchargelattice_args.argptr[8] = unsafe.Pointer(&settopologicalchargelattice_args.arg_PBC)
}

// Wrapper for settopologicalchargelattice CUDA kernel, asynchronous.
func k_settopologicalchargelattice_async(s unsafe.Pointer, mx unsafe.Pointer, my unsafe.Pointer, mz unsafe.Pointer, icxcy float32, Nx int, Ny int, Nz int, PBC byte, cfg *config) {
	if Synchronous { // debug
		Sync()
		timer_old.Start("settopologicalchargelattice")
	}

	settopologicalchargelattice_args.Lock()
	defer settopologicalchargelattice_args.Unlock()

	if settopologicalchargelattice_code == 0 {
		settopologicalchargelattice_code = fatbinLoad(settopologicalchargelattice_map, "settopologicalchargelattice")
	}

	settopologicalchargelattice_args.arg_s = s
	settopologicalchargelattice_args.arg_mx = mx
	settopologicalchargelattice_args.arg_my = my
	settopologicalchargelattice_args.arg_mz = mz
	settopologicalchargelattice_args.arg_icxcy = icxcy
	settopologicalchargelattice_args.arg_Nx = Nx
	settopologicalchargelattice_args.arg_Ny = Ny
	settopologicalchargelattice_args.arg_Nz = Nz
	settopologicalchargelattice_args.arg_PBC = PBC

	args := settopologicalchargelattice_args.argptr[:]
	cu.LaunchKernel(settopologicalchargelattice_code, cfg.Grid.X, cfg.Grid.Y, cfg.Grid.Z, cfg.Block.X, cfg.Block.Y, cfg.Block.Z, 0, stream0, args)

	if Synchronous { // debug
		Sync()
		timer_old.Stop("settopologicalchargelattice")
	}
}

// maps compute capability on PTX code for settopologicalchargelattice kernel.
var settopologicalchargelattice_map = map[int]string{0: "",
	52: settopologicalchargelattice_ptx_52}

// settopologicalchargelattice PTX code for various compute capabilities.
const (
	settopologicalchargelattice_ptx_52 = `
.version 7.0
.target sm_52
.address_size 64

	// .globl	settopologicalchargelattice

.visible .entry settopologicalchargelattice(
	.param .u64 settopologicalchargelattice_param_0,
	.param .u64 settopologicalchargelattice_param_1,
	.param .u64 settopologicalchargelattice_param_2,
	.param .u64 settopologicalchargelattice_param_3,
	.param .f32 settopologicalchargelattice_param_4,
	.param .u32 settopologicalchargelattice_param_5,
	.param .u32 settopologicalchargelattice_param_6,
	.param .u32 settopologicalchargelattice_param_7,
	.param .u8 settopologicalchargelattice_param_8
)
{
	.reg .pred 	%p<83>;
	.reg .b16 	%rs<13>;
	.reg .f32 	%f<295>;
	.reg .b32 	%r<171>;
	.reg .b64 	%rd<46>;


	ld.param.u64 	%rd5, [settopologicalchargelattice_param_0];
	ld.param.u64 	%rd6, [settopologicalchargelattice_param_1];
	ld.param.u64 	%rd7, [settopologicalchargelattice_param_2];
	ld.param.u64 	%rd8, [settopologicalchargelattice_param_3];
	ld.param.f32 	%f52, [settopologicalchargelattice_param_4];
	ld.param.u32 	%r58, [settopologicalchargelattice_param_5];
	ld.param.u32 	%r59, [settopologicalchargelattice_param_6];
	ld.param.u32 	%r60, [settopologicalchargelattice_param_7];
	ld.param.u8 	%rs3, [settopologicalchargelattice_param_8];
	cvta.to.global.u64 	%rd1, %rd8;
	cvta.to.global.u64 	%rd2, %rd7;
	cvta.to.global.u64 	%rd3, %rd6;
	mov.u32 	%r61, %ntid.x;
	mov.u32 	%r62, %ctaid.x;
	mov.u32 	%r63, %tid.x;
	mad.lo.s32 	%r1, %r61, %r62, %r63;
	mov.u32 	%r64, %ntid.y;
	mov.u32 	%r65, %ctaid.y;
	mov.u32 	%r66, %tid.y;
	mad.lo.s32 	%r2, %r64, %r65, %r66;
	mov.u32 	%r67, %ntid.z;
	mov.u32 	%r68, %ctaid.z;
	mov.u32 	%r69, %tid.z;
	mad.lo.s32 	%r3, %r67, %r68, %r69;
	setp.ge.s32	%p3, %r2, %r59;
	setp.ge.s32	%p4, %r1, %r58;
	or.pred  	%p5, %p3, %p4;
	setp.ge.s32	%p6, %r3, %r60;
	or.pred  	%p7, %p5, %p6;
	@%p7 bra 	BB0_72;

	cvta.to.global.u64 	%rd9, %rd5;
	mul.lo.s32 	%r4, %r3, %r59;
	add.s32 	%r70, %r4, %r2;
	mul.lo.s32 	%r5, %r70, %r58;
	add.s32 	%r71, %r5, %r1;
	mul.wide.s32 	%rd10, %r71, 4;
	add.s64 	%rd11, %rd3, %rd10;
	add.s64 	%rd12, %rd2, %rd10;
	add.s64 	%rd13, %rd1, %rd10;
	ld.global.nc.f32 	%f1, [%rd11];
	ld.global.nc.f32 	%f2, [%rd12];
	mul.f32 	%f53, %f2, %f2;
	fma.rn.f32 	%f54, %f1, %f1, %f53;
	ld.global.nc.f32 	%f3, [%rd13];
	fma.rn.f32 	%f55, %f3, %f3, %f54;
	setp.eq.f32	%p8, %f55, 0f00000000;
	add.s64 	%rd4, %rd9, %rd10;
	@%p8 bra 	BB0_71;
	bra.uni 	BB0_2;

BB0_71:
	mov.u32 	%r158, 0;
	st.global.u32 	[%rd4], %r158;
	bra.uni 	BB0_72;

BB0_2:
	and.b16  	%rs1, %rs3, 1;
	setp.eq.s16	%p9, %rs1, 0;
	add.s32 	%r6, %r1, 1;
	@%p9 bra 	BB0_4;

	rem.s32 	%r72, %r6, %r58;
	add.s32 	%r73, %r72, %r58;
	rem.s32 	%r159, %r73, %r58;
	bra.uni 	BB0_5;

BB0_4:
	add.s32 	%r74, %r58, -1;
	min.s32 	%r159, %r6, %r74;

BB0_5:
	and.b16  	%rs2, %rs3, 2;
	setp.eq.s16	%p10, %rs2, 0;
	add.s32 	%r10, %r2, 1;
	@%p10 bra 	BB0_7;

	rem.s32 	%r75, %r10, %r59;
	add.s32 	%r76, %r75, %r59;
	rem.s32 	%r160, %r76, %r59;
	bra.uni 	BB0_8;

BB0_7:
	add.s32 	%r77, %r59, -1;
	min.s32 	%r160, %r10, %r77;

BB0_8:
	add.s32 	%r14, %r1, -1;
	@%p9 bra 	BB0_10;

	rem.s32 	%r78, %r14, %r58;
	add.s32 	%r79, %r78, %r58;
	rem.s32 	%r161, %r79, %r58;
	bra.uni 	BB0_11;

BB0_10:
	mov.u32 	%r80, 0;
	max.s32 	%r161, %r14, %r80;

BB0_11:
	add.s32 	%r18, %r159, %r5;
	add.s32 	%r81, %r160, %r4;
	mad.lo.s32 	%r19, %r81, %r58, %r1;
	add.s32 	%r20, %r161, %r5;
	add.s32 	%r21, %r2, -1;
	@%p10 bra 	BB0_13;

	rem.s32 	%r82, %r21, %r59;
	add.s32 	%r83, %r82, %r59;
	rem.s32 	%r162, %r83, %r59;
	bra.uni 	BB0_14;

BB0_13:
	mov.u32 	%r84, 0;
	max.s32 	%r162, %r21, %r84;

BB0_14:
	add.s32 	%r85, %r162, %r4;
	mad.lo.s32 	%r86, %r85, %r58, %r1;
	mul.wide.s32 	%rd14, %r18, 4;
	add.s64 	%rd15, %rd3, %rd14;
	ld.global.nc.f32 	%f4, [%rd15];
	add.s64 	%rd16, %rd2, %rd14;
	ld.global.nc.f32 	%f5, [%rd16];
	add.s64 	%rd17, %rd1, %rd14;
	ld.global.nc.f32 	%f6, [%rd17];
	mul.wide.s32 	%rd18, %r19, 4;
	add.s64 	%rd19, %rd3, %rd18;
	ld.global.nc.f32 	%f7, [%rd19];
	add.s64 	%rd20, %rd2, %rd18;
	ld.global.nc.f32 	%f8, [%rd20];
	add.s64 	%rd21, %rd1, %rd18;
	ld.global.nc.f32 	%f9, [%rd21];
	mul.wide.s32 	%rd22, %r20, 4;
	add.s64 	%rd23, %rd3, %rd22;
	ld.global.nc.f32 	%f10, [%rd23];
	add.s64 	%rd24, %rd2, %rd22;
	ld.global.nc.f32 	%f11, [%rd24];
	add.s64 	%rd25, %rd1, %rd22;
	ld.global.nc.f32 	%f12, [%rd25];
	mul.wide.s32 	%rd26, %r86, 4;
	add.s64 	%rd27, %rd3, %rd26;
	ld.global.nc.f32 	%f13, [%rd27];
	add.s64 	%rd28, %rd2, %rd26;
	ld.global.nc.f32 	%f14, [%rd28];
	add.s64 	%rd29, %rd1, %rd26;
	ld.global.nc.f32 	%f15, [%rd29];
	setp.ne.s16	%p14, %rs1, 0;
	setp.ge.s32	%p15, %r6, %r58;
	setp.lt.s32	%p16, %r6, %r58;
	or.pred  	%p1, %p16, %p14;
	mov.f32 	%f290, 0f00000000;
	and.pred  	%p17, %p15, %p9;
	@%p17 bra 	BB0_28;

	setp.ge.s32	%p18, %r10, %r59;
	and.pred  	%p20, %p18, %p10;
	@%p20 bra 	BB0_28;

	@%p10 bra 	BB0_18;

	rem.s32 	%r87, %r10, %r59;
	add.s32 	%r88, %r87, %r59;
	rem.s32 	%r163, %r88, %r59;
	bra.uni 	BB0_19;

BB0_18:
	add.s32 	%r89, %r59, -1;
	min.s32 	%r163, %r10, %r89;

BB0_19:
	@%p9 bra 	BB0_21;

	rem.s32 	%r90, %r6, %r58;
	add.s32 	%r91, %r90, %r58;
	rem.s32 	%r164, %r91, %r58;
	bra.uni 	BB0_22;

BB0_21:
	add.s32 	%r92, %r58, -1;
	min.s32 	%r164, %r6, %r92;

BB0_22:
	add.s32 	%r93, %r163, %r4;
	mad.lo.s32 	%r94, %r93, %r58, %r164;
	mul.wide.s32 	%rd30, %r94, 4;
	add.s64 	%rd31, %rd3, %rd30;
	add.s64 	%rd32, %rd2, %rd30;
	add.s64 	%rd33, %rd1, %rd30;
	ld.global.nc.f32 	%f58, [%rd31];
	ld.global.nc.f32 	%f59, [%rd32];
	mul.f32 	%f60, %f59, %f59;
	fma.rn.f32 	%f61, %f58, %f58, %f60;
	ld.global.nc.f32 	%f62, [%rd33];
	fma.rn.f32 	%f16, %f62, %f62, %f61;
	mul.f32 	%f63, %f6, %f8;
	mul.f32 	%f64, %f5, %f9;
	sub.f32 	%f65, %f64, %f63;
	mul.f32 	%f66, %f4, %f9;
	mul.f32 	%f67, %f6, %f7;
	sub.f32 	%f68, %f67, %f66;
	mul.f32 	%f69, %f5, %f7;
	mul.f32 	%f70, %f4, %f8;
	sub.f32 	%f71, %f70, %f69;
	mul.f32 	%f72, %f2, %f68;
	fma.rn.f32 	%f73, %f1, %f65, %f72;
	fma.rn.f32 	%f74, %f3, %f71, %f73;
	mul.f32 	%f75, %f2, %f5;
	fma.rn.f32 	%f76, %f1, %f4, %f75;
	fma.rn.f32 	%f77, %f3, %f6, %f76;
	add.f32 	%f78, %f77, 0f3F800000;
	mul.f32 	%f79, %f2, %f8;
	fma.rn.f32 	%f80, %f1, %f7, %f79;
	fma.rn.f32 	%f81, %f3, %f9, %f80;
	add.f32 	%f82, %f78, %f81;
	mul.f32 	%f83, %f5, %f8;
	fma.rn.f32 	%f84, %f4, %f7, %f83;
	fma.rn.f32 	%f85, %f6, %f9, %f84;
	add.f32 	%f86, %f85, %f82;
	abs.f32 	%f17, %f86;
	abs.f32 	%f18, %f74;
	setp.eq.f32	%p23, %f17, 0f00000000;
	setp.eq.f32	%p24, %f18, 0f00000000;
	and.pred  	%p25, %p23, %p24;
	mov.b32 	 %r31, %f86;
	mov.b32 	 %r95, %f74;
	and.b32  	%r32, %r95, -2147483648;
	@%p25 bra 	BB0_26;
	bra.uni 	BB0_23;

BB0_26:
	shr.s32 	%r102, %r31, 31;
	and.b32  	%r103, %r102, 1078530011;
	or.b32  	%r104, %r103, %r32;
	mov.b32 	 %f287, %r104;
	bra.uni 	BB0_27;

BB0_23:
	setp.eq.f32	%p26, %f17, 0f7F800000;
	setp.eq.f32	%p27, %f18, 0f7F800000;
	and.pred  	%p28, %p26, %p27;
	@%p28 bra 	BB0_25;
	bra.uni 	BB0_24;

BB0_25:
	shr.s32 	%r98, %r31, 31;
	and.b32  	%r99, %r98, 13483017;
	add.s32 	%r100, %r99, 1061752795;
	or.b32  	%r101, %r100, %r32;
	mov.b32 	 %f287, %r101;
	bra.uni 	BB0_27;

BB0_24:
	max.f32 	%f87, %f18, %f17;
	min.f32 	%f88, %f18, %f17;
	div.rn.f32 	%f89, %f88, %f87;
	mul.rn.f32 	%f90, %f89, %f89;
	mov.f32 	%f91, 0fC0B59883;
	mov.f32 	%f92, 0fBF52C7EA;
	fma.rn.f32 	%f93, %f90, %f92, %f91;
	mov.f32 	%f94, 0fC0D21907;
	fma.rn.f32 	%f95, %f93, %f90, %f94;
	mul.f32 	%f96, %f90, %f95;
	mul.f32 	%f97, %f89, %f96;
	add.f32 	%f98, %f90, 0f41355DC0;
	mov.f32 	%f99, 0f41E6BD60;
	fma.rn.f32 	%f100, %f98, %f90, %f99;
	mov.f32 	%f101, 0f419D92C8;
	fma.rn.f32 	%f102, %f100, %f90, %f101;
	rcp.rn.f32 	%f103, %f102;
	fma.rn.f32 	%f104, %f97, %f103, %f89;
	mov.f32 	%f105, 0f3FC90FDB;
	sub.f32 	%f106, %f105, %f104;
	setp.gt.f32	%p29, %f18, %f17;
	selp.f32	%f107, %f106, %f104, %p29;
	mov.f32 	%f108, 0f40490FDB;
	sub.f32 	%f109, %f108, %f107;
	setp.lt.s32	%p30, %r31, 0;
	selp.f32	%f110, %f109, %f107, %p30;
	mov.b32 	 %r96, %f110;
	or.b32  	%r97, %r96, %r32;
	mov.b32 	 %f111, %r97;
	add.f32 	%f112, %f17, %f18;
	setp.gtu.f32	%p31, %f112, 0f7F800000;
	selp.f32	%f287, %f112, %f111, %p31;

BB0_27:
	add.f32 	%f113, %f287, %f287;
	setp.eq.f32	%p32, %f16, 0f00000000;
	selp.f32	%f114, 0f3F800000, 0f3F000000, %p32;
	fma.rn.f32 	%f290, %f114, %f113, 0f00000000;

BB0_28:
	setp.lt.s32	%p33, %r14, 0;
	setp.gt.s32	%p34, %r14, -1;
	or.pred  	%p2, %p34, %p14;
	and.pred  	%p37, %p33, %p9;
	@%p37 bra 	BB0_42;

	setp.ge.s32	%p38, %r10, %r59;
	and.pred  	%p40, %p38, %p10;
	@%p40 bra 	BB0_42;

	@%p10 bra 	BB0_32;

	rem.s32 	%r105, %r10, %r59;
	add.s32 	%r106, %r105, %r59;
	rem.s32 	%r165, %r106, %r59;
	bra.uni 	BB0_33;

BB0_32:
	add.s32 	%r107, %r59, -1;
	min.s32 	%r165, %r10, %r107;

BB0_33:
	@%p9 bra 	BB0_35;

	rem.s32 	%r108, %r14, %r58;
	add.s32 	%r109, %r108, %r58;
	rem.s32 	%r166, %r109, %r58;
	bra.uni 	BB0_36;

BB0_35:
	mov.u32 	%r110, 0;
	max.s32 	%r166, %r14, %r110;

BB0_36:
	add.s32 	%r111, %r165, %r4;
	mad.lo.s32 	%r112, %r111, %r58, %r166;
	mul.wide.s32 	%rd34, %r112, 4;
	add.s64 	%rd35, %rd3, %rd34;
	add.s64 	%rd36, %rd2, %rd34;
	add.s64 	%rd37, %rd1, %rd34;
	ld.global.nc.f32 	%f115, [%rd35];
	ld.global.nc.f32 	%f116, [%rd36];
	mul.f32 	%f117, %f116, %f116;
	fma.rn.f32 	%f118, %f115, %f115, %f117;
	ld.global.nc.f32 	%f119, [%rd37];
	fma.rn.f32 	%f25, %f119, %f119, %f118;
	mul.f32 	%f120, %f9, %f11;
	mul.f32 	%f121, %f8, %f12;
	sub.f32 	%f122, %f121, %f120;
	mul.f32 	%f123, %f7, %f12;
	mul.f32 	%f124, %f9, %f10;
	sub.f32 	%f125, %f124, %f123;
	mul.f32 	%f126, %f8, %f10;
	mul.f32 	%f127, %f7, %f11;
	sub.f32 	%f128, %f127, %f126;
	mul.f32 	%f129, %f2, %f125;
	fma.rn.f32 	%f130, %f1, %f122, %f129;
	fma.rn.f32 	%f131, %f3, %f128, %f130;
	mul.f32 	%f132, %f2, %f8;
	fma.rn.f32 	%f133, %f1, %f7, %f132;
	fma.rn.f32 	%f134, %f3, %f9, %f133;
	add.f32 	%f135, %f134, 0f3F800000;
	mul.f32 	%f136, %f2, %f11;
	fma.rn.f32 	%f137, %f1, %f10, %f136;
	fma.rn.f32 	%f138, %f3, %f12, %f137;
	add.f32 	%f139, %f135, %f138;
	mul.f32 	%f140, %f8, %f11;
	fma.rn.f32 	%f141, %f7, %f10, %f140;
	fma.rn.f32 	%f142, %f9, %f12, %f141;
	add.f32 	%f143, %f142, %f139;
	abs.f32 	%f26, %f143;
	abs.f32 	%f27, %f131;
	setp.eq.f32	%p43, %f26, 0f00000000;
	setp.eq.f32	%p44, %f27, 0f00000000;
	and.pred  	%p45, %p43, %p44;
	mov.b32 	 %r39, %f143;
	mov.b32 	 %r113, %f131;
	and.b32  	%r40, %r113, -2147483648;
	@%p45 bra 	BB0_40;
	bra.uni 	BB0_37;

BB0_40:
	shr.s32 	%r120, %r39, 31;
	and.b32  	%r121, %r120, 1078530011;
	or.b32  	%r122, %r121, %r40;
	mov.b32 	 %f289, %r122;
	bra.uni 	BB0_41;

BB0_37:
	setp.eq.f32	%p46, %f26, 0f7F800000;
	setp.eq.f32	%p47, %f27, 0f7F800000;
	and.pred  	%p48, %p46, %p47;
	@%p48 bra 	BB0_39;
	bra.uni 	BB0_38;

BB0_39:
	shr.s32 	%r116, %r39, 31;
	and.b32  	%r117, %r116, 13483017;
	add.s32 	%r118, %r117, 1061752795;
	or.b32  	%r119, %r118, %r40;
	mov.b32 	 %f289, %r119;
	bra.uni 	BB0_41;

BB0_38:
	max.f32 	%f144, %f27, %f26;
	min.f32 	%f145, %f27, %f26;
	div.rn.f32 	%f146, %f145, %f144;
	mul.rn.f32 	%f147, %f146, %f146;
	mov.f32 	%f148, 0fC0B59883;
	mov.f32 	%f149, 0fBF52C7EA;
	fma.rn.f32 	%f150, %f147, %f149, %f148;
	mov.f32 	%f151, 0fC0D21907;
	fma.rn.f32 	%f152, %f150, %f147, %f151;
	mul.f32 	%f153, %f147, %f152;
	mul.f32 	%f154, %f146, %f153;
	add.f32 	%f155, %f147, 0f41355DC0;
	mov.f32 	%f156, 0f41E6BD60;
	fma.rn.f32 	%f157, %f155, %f147, %f156;
	mov.f32 	%f158, 0f419D92C8;
	fma.rn.f32 	%f159, %f157, %f147, %f158;
	rcp.rn.f32 	%f160, %f159;
	fma.rn.f32 	%f161, %f154, %f160, %f146;
	mov.f32 	%f162, 0f3FC90FDB;
	sub.f32 	%f163, %f162, %f161;
	setp.gt.f32	%p49, %f27, %f26;
	selp.f32	%f164, %f163, %f161, %p49;
	mov.f32 	%f165, 0f40490FDB;
	sub.f32 	%f166, %f165, %f164;
	setp.lt.s32	%p50, %r39, 0;
	selp.f32	%f167, %f166, %f164, %p50;
	mov.b32 	 %r114, %f167;
	or.b32  	%r115, %r114, %r40;
	mov.b32 	 %f168, %r115;
	add.f32 	%f169, %f26, %f27;
	setp.gtu.f32	%p51, %f169, 0f7F800000;
	selp.f32	%f289, %f169, %f168, %p51;

BB0_41:
	add.f32 	%f170, %f289, %f289;
	setp.eq.f32	%p52, %f25, 0f00000000;
	selp.f32	%f171, 0f3F800000, 0f3F000000, %p52;
	fma.rn.f32 	%f290, %f171, %f170, %f290;

BB0_42:
	@!%p2 bra 	BB0_56;
	bra.uni 	BB0_43;

BB0_43:
	setp.lt.s32	%p53, %r21, 0;
	and.pred  	%p55, %p53, %p10;
	@%p55 bra 	BB0_56;

	@%p10 bra 	BB0_46;

	rem.s32 	%r123, %r21, %r59;
	add.s32 	%r124, %r123, %r59;
	rem.s32 	%r167, %r124, %r59;
	bra.uni 	BB0_47;

BB0_46:
	mov.u32 	%r125, 0;
	max.s32 	%r167, %r21, %r125;

BB0_47:
	@%p9 bra 	BB0_49;

	rem.s32 	%r126, %r14, %r58;
	add.s32 	%r127, %r126, %r58;
	rem.s32 	%r168, %r127, %r58;
	bra.uni 	BB0_50;

BB0_49:
	mov.u32 	%r128, 0;
	max.s32 	%r168, %r14, %r128;

BB0_50:
	add.s32 	%r129, %r167, %r4;
	mad.lo.s32 	%r130, %r129, %r58, %r168;
	mul.wide.s32 	%rd38, %r130, 4;
	add.s64 	%rd39, %rd3, %rd38;
	add.s64 	%rd40, %rd2, %rd38;
	add.s64 	%rd41, %rd1, %rd38;
	ld.global.nc.f32 	%f172, [%rd39];
	ld.global.nc.f32 	%f173, [%rd40];
	mul.f32 	%f174, %f173, %f173;
	fma.rn.f32 	%f175, %f172, %f172, %f174;
	ld.global.nc.f32 	%f176, [%rd41];
	fma.rn.f32 	%f34, %f176, %f176, %f175;
	mul.f32 	%f177, %f12, %f14;
	mul.f32 	%f178, %f11, %f15;
	sub.f32 	%f179, %f178, %f177;
	mul.f32 	%f180, %f10, %f15;
	mul.f32 	%f181, %f12, %f13;
	sub.f32 	%f182, %f181, %f180;
	mul.f32 	%f183, %f11, %f13;
	mul.f32 	%f184, %f10, %f14;
	sub.f32 	%f185, %f184, %f183;
	mul.f32 	%f186, %f2, %f182;
	fma.rn.f32 	%f187, %f1, %f179, %f186;
	fma.rn.f32 	%f188, %f3, %f185, %f187;
	mul.f32 	%f189, %f2, %f11;
	fma.rn.f32 	%f190, %f1, %f10, %f189;
	fma.rn.f32 	%f191, %f3, %f12, %f190;
	add.f32 	%f192, %f191, 0f3F800000;
	mul.f32 	%f193, %f2, %f14;
	fma.rn.f32 	%f194, %f1, %f13, %f193;
	fma.rn.f32 	%f195, %f3, %f15, %f194;
	add.f32 	%f196, %f192, %f195;
	mul.f32 	%f197, %f11, %f14;
	fma.rn.f32 	%f198, %f10, %f13, %f197;
	fma.rn.f32 	%f199, %f12, %f15, %f198;
	add.f32 	%f200, %f199, %f196;
	abs.f32 	%f35, %f200;
	abs.f32 	%f36, %f188;
	setp.eq.f32	%p58, %f35, 0f00000000;
	setp.eq.f32	%p59, %f36, 0f00000000;
	and.pred  	%p60, %p58, %p59;
	mov.b32 	 %r47, %f200;
	mov.b32 	 %r131, %f188;
	and.b32  	%r48, %r131, -2147483648;
	@%p60 bra 	BB0_54;
	bra.uni 	BB0_51;

BB0_54:
	shr.s32 	%r138, %r47, 31;
	and.b32  	%r139, %r138, 1078530011;
	or.b32  	%r140, %r139, %r48;
	mov.b32 	 %f291, %r140;
	bra.uni 	BB0_55;

BB0_51:
	setp.eq.f32	%p61, %f35, 0f7F800000;
	setp.eq.f32	%p62, %f36, 0f7F800000;
	and.pred  	%p63, %p61, %p62;
	@%p63 bra 	BB0_53;
	bra.uni 	BB0_52;

BB0_53:
	shr.s32 	%r134, %r47, 31;
	and.b32  	%r135, %r134, 13483017;
	add.s32 	%r136, %r135, 1061752795;
	or.b32  	%r137, %r136, %r48;
	mov.b32 	 %f291, %r137;
	bra.uni 	BB0_55;

BB0_52:
	max.f32 	%f201, %f36, %f35;
	min.f32 	%f202, %f36, %f35;
	div.rn.f32 	%f203, %f202, %f201;
	mul.rn.f32 	%f204, %f203, %f203;
	mov.f32 	%f205, 0fC0B59883;
	mov.f32 	%f206, 0fBF52C7EA;
	fma.rn.f32 	%f207, %f204, %f206, %f205;
	mov.f32 	%f208, 0fC0D21907;
	fma.rn.f32 	%f209, %f207, %f204, %f208;
	mul.f32 	%f210, %f204, %f209;
	mul.f32 	%f211, %f203, %f210;
	add.f32 	%f212, %f204, 0f41355DC0;
	mov.f32 	%f213, 0f41E6BD60;
	fma.rn.f32 	%f214, %f212, %f204, %f213;
	mov.f32 	%f215, 0f419D92C8;
	fma.rn.f32 	%f216, %f214, %f204, %f215;
	rcp.rn.f32 	%f217, %f216;
	fma.rn.f32 	%f218, %f211, %f217, %f203;
	mov.f32 	%f219, 0f3FC90FDB;
	sub.f32 	%f220, %f219, %f218;
	setp.gt.f32	%p64, %f36, %f35;
	selp.f32	%f221, %f220, %f218, %p64;
	mov.f32 	%f222, 0f40490FDB;
	sub.f32 	%f223, %f222, %f221;
	setp.lt.s32	%p65, %r47, 0;
	selp.f32	%f224, %f223, %f221, %p65;
	mov.b32 	 %r132, %f224;
	or.b32  	%r133, %r132, %r48;
	mov.b32 	 %f225, %r133;
	add.f32 	%f226, %f35, %f36;
	setp.gtu.f32	%p66, %f226, 0f7F800000;
	selp.f32	%f291, %f226, %f225, %p66;

BB0_55:
	add.f32 	%f227, %f291, %f291;
	setp.eq.f32	%p67, %f34, 0f00000000;
	selp.f32	%f228, 0f3F800000, 0f3F000000, %p67;
	fma.rn.f32 	%f290, %f228, %f227, %f290;

BB0_56:
	@!%p1 bra 	BB0_70;
	bra.uni 	BB0_57;

BB0_57:
	setp.lt.s32	%p68, %r21, 0;
	and.pred  	%p70, %p68, %p10;
	@%p70 bra 	BB0_70;

	@%p10 bra 	BB0_60;

	rem.s32 	%r141, %r21, %r59;
	add.s32 	%r142, %r141, %r59;
	rem.s32 	%r169, %r142, %r59;
	bra.uni 	BB0_61;

BB0_60:
	mov.u32 	%r143, 0;
	max.s32 	%r169, %r21, %r143;

BB0_61:
	add.s32 	%r52, %r169, %r4;
	@%p9 bra 	BB0_63;

	rem.s32 	%r144, %r6, %r58;
	add.s32 	%r145, %r144, %r58;
	rem.s32 	%r170, %r145, %r58;
	bra.uni 	BB0_64;

BB0_63:
	add.s32 	%r146, %r58, -1;
	min.s32 	%r170, %r6, %r146;

BB0_64:
	mad.lo.s32 	%r147, %r52, %r58, %r170;
	mul.wide.s32 	%rd42, %r147, 4;
	add.s64 	%rd43, %rd3, %rd42;
	add.s64 	%rd44, %rd2, %rd42;
	add.s64 	%rd45, %rd1, %rd42;
	ld.global.nc.f32 	%f229, [%rd43];
	ld.global.nc.f32 	%f230, [%rd44];
	mul.f32 	%f231, %f230, %f230;
	fma.rn.f32 	%f232, %f229, %f229, %f231;
	ld.global.nc.f32 	%f233, [%rd45];
	fma.rn.f32 	%f43, %f233, %f233, %f232;
	mul.f32 	%f234, %f5, %f15;
	mul.f32 	%f235, %f6, %f14;
	sub.f32 	%f236, %f235, %f234;
	mul.f32 	%f237, %f6, %f13;
	mul.f32 	%f238, %f4, %f15;
	sub.f32 	%f239, %f238, %f237;
	mul.f32 	%f240, %f4, %f14;
	mul.f32 	%f241, %f5, %f13;
	sub.f32 	%f242, %f241, %f240;
	mul.f32 	%f243, %f2, %f239;
	fma.rn.f32 	%f244, %f1, %f236, %f243;
	fma.rn.f32 	%f245, %f3, %f242, %f244;
	mul.f32 	%f246, %f2, %f14;
	fma.rn.f32 	%f247, %f1, %f13, %f246;
	fma.rn.f32 	%f248, %f3, %f15, %f247;
	add.f32 	%f249, %f248, 0f3F800000;
	mul.f32 	%f250, %f2, %f5;
	fma.rn.f32 	%f251, %f1, %f4, %f250;
	fma.rn.f32 	%f252, %f3, %f6, %f251;
	add.f32 	%f253, %f252, %f249;
	mul.f32 	%f254, %f5, %f14;
	fma.rn.f32 	%f255, %f4, %f13, %f254;
	fma.rn.f32 	%f256, %f6, %f15, %f255;
	add.f32 	%f257, %f256, %f253;
	abs.f32 	%f44, %f257;
	abs.f32 	%f45, %f245;
	setp.eq.f32	%p73, %f44, 0f00000000;
	setp.eq.f32	%p74, %f45, 0f00000000;
	and.pred  	%p75, %p73, %p74;
	mov.b32 	 %r56, %f257;
	mov.b32 	 %r148, %f245;
	and.b32  	%r57, %r148, -2147483648;
	@%p75 bra 	BB0_68;
	bra.uni 	BB0_65;

BB0_68:
	shr.s32 	%r155, %r56, 31;
	and.b32  	%r156, %r155, 1078530011;
	or.b32  	%r157, %r156, %r57;
	mov.b32 	 %f293, %r157;
	bra.uni 	BB0_69;

BB0_65:
	setp.eq.f32	%p76, %f44, 0f7F800000;
	setp.eq.f32	%p77, %f45, 0f7F800000;
	and.pred  	%p78, %p76, %p77;
	@%p78 bra 	BB0_67;
	bra.uni 	BB0_66;

BB0_67:
	shr.s32 	%r151, %r56, 31;
	and.b32  	%r152, %r151, 13483017;
	add.s32 	%r153, %r152, 1061752795;
	or.b32  	%r154, %r153, %r57;
	mov.b32 	 %f293, %r154;
	bra.uni 	BB0_69;

BB0_66:
	max.f32 	%f258, %f45, %f44;
	min.f32 	%f259, %f45, %f44;
	div.rn.f32 	%f260, %f259, %f258;
	mul.rn.f32 	%f261, %f260, %f260;
	mov.f32 	%f262, 0fC0B59883;
	mov.f32 	%f263, 0fBF52C7EA;
	fma.rn.f32 	%f264, %f261, %f263, %f262;
	mov.f32 	%f265, 0fC0D21907;
	fma.rn.f32 	%f266, %f264, %f261, %f265;
	mul.f32 	%f267, %f261, %f266;
	mul.f32 	%f268, %f260, %f267;
	add.f32 	%f269, %f261, 0f41355DC0;
	mov.f32 	%f270, 0f41E6BD60;
	fma.rn.f32 	%f271, %f269, %f261, %f270;
	mov.f32 	%f272, 0f419D92C8;
	fma.rn.f32 	%f273, %f271, %f261, %f272;
	rcp.rn.f32 	%f274, %f273;
	fma.rn.f32 	%f275, %f268, %f274, %f260;
	mov.f32 	%f276, 0f3FC90FDB;
	sub.f32 	%f277, %f276, %f275;
	setp.gt.f32	%p79, %f45, %f44;
	selp.f32	%f278, %f277, %f275, %p79;
	mov.f32 	%f279, 0f40490FDB;
	sub.f32 	%f280, %f279, %f278;
	setp.lt.s32	%p80, %r56, 0;
	selp.f32	%f281, %f280, %f278, %p80;
	mov.b32 	 %r149, %f281;
	or.b32  	%r150, %r149, %r57;
	mov.b32 	 %f282, %r150;
	add.f32 	%f283, %f44, %f45;
	setp.gtu.f32	%p81, %f283, 0f7F800000;
	selp.f32	%f293, %f283, %f282, %p81;

BB0_69:
	add.f32 	%f284, %f293, %f293;
	setp.eq.f32	%p82, %f43, 0f00000000;
	selp.f32	%f285, 0f3F800000, 0f3F000000, %p82;
	fma.rn.f32 	%f290, %f285, %f284, %f290;

BB0_70:
	mul.f32 	%f286, %f290, %f52;
	st.global.f32 	[%rd4], %f286;

BB0_72:
	ret;
}


`
)
