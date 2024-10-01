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

// CUDA handle for setmaxangle kernel
var setmaxangle_code cu.Function

// Stores the arguments for setmaxangle kernel invocation
type setmaxangle_args_t struct{
	 arg_dst unsafe.Pointer
	 arg_mx unsafe.Pointer
	 arg_my unsafe.Pointer
	 arg_mz unsafe.Pointer
	 arg_aLUT2d unsafe.Pointer
	 arg_regions unsafe.Pointer
	 arg_Nx int
	 arg_Ny int
	 arg_Nz int
	 arg_PBC byte
	 argptr [10]unsafe.Pointer
	sync.Mutex
}

// Stores the arguments for setmaxangle kernel invocation
var setmaxangle_args setmaxangle_args_t

func init(){
	// CUDA driver kernel call wants pointers to arguments, set them up once.
	 setmaxangle_args.argptr[0] = unsafe.Pointer(&setmaxangle_args.arg_dst)
	 setmaxangle_args.argptr[1] = unsafe.Pointer(&setmaxangle_args.arg_mx)
	 setmaxangle_args.argptr[2] = unsafe.Pointer(&setmaxangle_args.arg_my)
	 setmaxangle_args.argptr[3] = unsafe.Pointer(&setmaxangle_args.arg_mz)
	 setmaxangle_args.argptr[4] = unsafe.Pointer(&setmaxangle_args.arg_aLUT2d)
	 setmaxangle_args.argptr[5] = unsafe.Pointer(&setmaxangle_args.arg_regions)
	 setmaxangle_args.argptr[6] = unsafe.Pointer(&setmaxangle_args.arg_Nx)
	 setmaxangle_args.argptr[7] = unsafe.Pointer(&setmaxangle_args.arg_Ny)
	 setmaxangle_args.argptr[8] = unsafe.Pointer(&setmaxangle_args.arg_Nz)
	 setmaxangle_args.argptr[9] = unsafe.Pointer(&setmaxangle_args.arg_PBC)
	 }

// Wrapper for setmaxangle CUDA kernel, asynchronous.
func k_setmaxangle_async ( dst unsafe.Pointer, mx unsafe.Pointer, my unsafe.Pointer, mz unsafe.Pointer, aLUT2d unsafe.Pointer, regions unsafe.Pointer, Nx int, Ny int, Nz int, PBC byte,  cfg *config) {
	if Synchronous{ // debug
		Sync()
		timer.Start("setmaxangle")
	}

	setmaxangle_args.Lock()
	defer setmaxangle_args.Unlock()

	if setmaxangle_code == 0{
		setmaxangle_code = fatbinLoad(setmaxangle_map, "setmaxangle")
	}

	 setmaxangle_args.arg_dst = dst
	 setmaxangle_args.arg_mx = mx
	 setmaxangle_args.arg_my = my
	 setmaxangle_args.arg_mz = mz
	 setmaxangle_args.arg_aLUT2d = aLUT2d
	 setmaxangle_args.arg_regions = regions
	 setmaxangle_args.arg_Nx = Nx
	 setmaxangle_args.arg_Ny = Ny
	 setmaxangle_args.arg_Nz = Nz
	 setmaxangle_args.arg_PBC = PBC
	

	args := setmaxangle_args.argptr[:]
	cu.LaunchKernel(setmaxangle_code, cfg.Grid.X, cfg.Grid.Y, cfg.Grid.Z, cfg.Block.X, cfg.Block.Y, cfg.Block.Z, 0, stream0, args)

	if Synchronous{ // debug
		Sync()
		timer.Stop("setmaxangle")
	}
}

// maps compute capability on PTX code for setmaxangle kernel.
var setmaxangle_map = map[int]string{ 0: "" ,
52: setmaxangle_ptx_52  }

// setmaxangle PTX code for various compute capabilities.
const(
  setmaxangle_ptx_52 = `
.version 7.0
.target sm_52
.address_size 64

	// .globl	setmaxangle

.visible .entry setmaxangle(
	.param .u64 setmaxangle_param_0,
	.param .u64 setmaxangle_param_1,
	.param .u64 setmaxangle_param_2,
	.param .u64 setmaxangle_param_3,
	.param .u64 setmaxangle_param_4,
	.param .u64 setmaxangle_param_5,
	.param .u32 setmaxangle_param_6,
	.param .u32 setmaxangle_param_7,
	.param .u32 setmaxangle_param_8,
	.param .u8 setmaxangle_param_9
)
{
	.reg .pred 	%p<44>;
	.reg .b16 	%rs<26>;
	.reg .f32 	%f<255>;
	.reg .b32 	%r<128>;
	.reg .b64 	%rd<69>;


	ld.param.u64 	%rd6, [setmaxangle_param_0];
	ld.param.u64 	%rd7, [setmaxangle_param_1];
	ld.param.u64 	%rd8, [setmaxangle_param_2];
	ld.param.u64 	%rd9, [setmaxangle_param_3];
	ld.param.u64 	%rd10, [setmaxangle_param_4];
	ld.param.u64 	%rd11, [setmaxangle_param_5];
	ld.param.u32 	%r32, [setmaxangle_param_6];
	ld.param.u32 	%r33, [setmaxangle_param_7];
	ld.param.u32 	%r34, [setmaxangle_param_8];
	ld.param.u8 	%rs5, [setmaxangle_param_9];
	cvta.to.global.u64 	%rd1, %rd10;
	cvta.to.global.u64 	%rd2, %rd11;
	cvta.to.global.u64 	%rd3, %rd9;
	cvta.to.global.u64 	%rd4, %rd8;
	cvta.to.global.u64 	%rd5, %rd7;
	mov.u32 	%r35, %ntid.x;
	mov.u32 	%r36, %ctaid.x;
	mov.u32 	%r37, %tid.x;
	mad.lo.s32 	%r1, %r35, %r36, %r37;
	mov.u32 	%r38, %ntid.y;
	mov.u32 	%r39, %ctaid.y;
	mov.u32 	%r40, %tid.y;
	mad.lo.s32 	%r2, %r38, %r39, %r40;
	mov.u32 	%r41, %ntid.z;
	mov.u32 	%r42, %ctaid.z;
	mov.u32 	%r43, %tid.z;
	mad.lo.s32 	%r3, %r41, %r42, %r43;
	setp.ge.s32	%p1, %r2, %r33;
	setp.ge.s32	%p2, %r1, %r32;
	or.pred  	%p3, %p1, %p2;
	setp.ge.s32	%p4, %r3, %r34;
	or.pred  	%p5, %p3, %p4;
	@%p5 bra 	BB0_34;

	mul.lo.s32 	%r4, %r3, %r33;
	add.s32 	%r44, %r4, %r2;
	mul.lo.s32 	%r5, %r44, %r32;
	add.s32 	%r6, %r5, %r1;
	mul.wide.s32 	%rd12, %r6, 4;
	add.s64 	%rd13, %rd5, %rd12;
	add.s64 	%rd14, %rd4, %rd12;
	add.s64 	%rd15, %rd3, %rd12;
	ld.global.nc.f32 	%f1, [%rd13];
	ld.global.nc.f32 	%f2, [%rd14];
	mul.f32 	%f38, %f2, %f2;
	fma.rn.f32 	%f39, %f1, %f1, %f38;
	ld.global.nc.f32 	%f3, [%rd15];
	fma.rn.f32 	%f40, %f3, %f3, %f39;
	setp.eq.f32	%p6, %f40, 0f00000000;
	@%p6 bra 	BB0_34;

	cvt.s64.s32	%rd16, %r6;
	add.s64 	%rd17, %rd2, %rd16;
	ld.global.nc.u8 	%rs1, [%rd17];
	cvt.u32.u16	%r45, %rs1;
	and.b32  	%r7, %r45, 255;
	and.b16  	%rs2, %rs5, 1;
	setp.eq.s16	%p7, %rs2, 0;
	add.s32 	%r8, %r1, -1;
	@%p7 bra 	BB0_4;

	rem.s32 	%r46, %r8, %r32;
	add.s32 	%r47, %r46, %r32;
	rem.s32 	%r122, %r47, %r32;
	bra.uni 	BB0_5;

BB0_4:
	mov.u32 	%r48, 0;
	max.s32 	%r122, %r8, %r48;

BB0_5:
	add.s32 	%r49, %r122, %r5;
	cvt.s64.s32	%rd18, %r49;
	mul.wide.s32 	%rd19, %r49, 4;
	add.s64 	%rd20, %rd5, %rd19;
	add.s64 	%rd21, %rd4, %rd19;
	add.s64 	%rd22, %rd3, %rd19;
	ld.global.nc.f32 	%f4, [%rd20];
	ld.global.nc.f32 	%f5, [%rd21];
	mul.f32 	%f42, %f5, %f5;
	fma.rn.f32 	%f43, %f4, %f4, %f42;
	ld.global.nc.f32 	%f6, [%rd22];
	fma.rn.f32 	%f7, %f6, %f6, %f43;
	add.s64 	%rd23, %rd2, %rd18;
	ld.global.nc.u8 	%rs6, [%rd23];
	setp.gt.u16	%p8, %rs6, %rs1;
	cvt.u32.u16	%r50, %rs6;
	and.b32  	%r51, %r50, 255;
	selp.b32	%r52, %r7, %r51, %p8;
	selp.b32	%r53, %r51, %r7, %p8;
	add.s32 	%r54, %r53, 1;
	mul.lo.s32 	%r55, %r54, %r53;
	shr.u32 	%r56, %r55, 1;
	add.s32 	%r57, %r56, %r52;
	mul.wide.s32 	%rd24, %r57, 4;
	add.s64 	%rd25, %rd1, %rd24;
	ld.global.nc.f32 	%f44, [%rd25];
	mov.f32 	%f250, 0f00000000;
	setp.eq.f32	%p9, %f44, 0f00000000;
	@%p9 bra 	BB0_7;

	setp.eq.f32	%p10, %f7, 0f00000000;
	selp.f32	%f45, %f1, %f4, %p10;
	selp.f32	%f46, %f2, %f5, %p10;
	mul.f32 	%f47, %f2, %f46;
	fma.rn.f32 	%f48, %f1, %f45, %f47;
	selp.f32	%f49, %f3, %f6, %p10;
	fma.rn.f32 	%f50, %f3, %f49, %f48;
	abs.f32 	%f51, %f50;
	mov.f32 	%f52, 0f3F800000;
	sub.f32 	%f53, %f52, %f51;
	mul.f32 	%f54, %f53, 0f3F000000;
	sqrt.rn.f32 	%f55, %f54;
	setp.gt.f32	%p11, %f51, 0f3F11EB85;
	selp.f32	%f56, %f55, %f51, %p11;
	mul.f32 	%f57, %f56, %f56;
	mov.f32 	%f58, 0f3C94D2E9;
	mov.f32 	%f59, 0f3D53F941;
	fma.rn.f32 	%f60, %f59, %f57, %f58;
	mov.f32 	%f61, 0f3D3F841F;
	fma.rn.f32 	%f62, %f60, %f57, %f61;
	mov.f32 	%f63, 0f3D994929;
	fma.rn.f32 	%f64, %f62, %f57, %f63;
	mov.f32 	%f65, 0f3E2AAB94;
	fma.rn.f32 	%f66, %f64, %f57, %f65;
	mul.f32 	%f67, %f57, %f66;
	fma.rn.f32 	%f68, %f67, %f56, %f56;
	add.f32 	%f69, %f68, %f68;
	mov.f32 	%f70, 0f3FC90FDB;
	sub.f32 	%f71, %f70, %f68;
	selp.f32	%f72, %f69, %f71, %p11;
	setp.lt.f32	%p12, %f50, 0f00000000;
	mov.f32 	%f73, 0f40490FDB;
	sub.f32 	%f74, %f73, %f72;
	selp.f32	%f75, %f74, %f72, %p12;
	mov.f32 	%f76, 0f00000000;
	max.f32 	%f250, %f76, %f75;

BB0_7:
	add.s32 	%r12, %r1, 1;
	@%p7 bra 	BB0_9;

	rem.s32 	%r58, %r12, %r32;
	add.s32 	%r59, %r58, %r32;
	rem.s32 	%r123, %r59, %r32;
	bra.uni 	BB0_10;

BB0_9:
	add.s32 	%r60, %r32, -1;
	min.s32 	%r123, %r12, %r60;

BB0_10:
	add.s32 	%r61, %r123, %r5;
	cvt.s64.s32	%rd26, %r61;
	mul.wide.s32 	%rd27, %r61, 4;
	add.s64 	%rd28, %rd5, %rd27;
	add.s64 	%rd29, %rd4, %rd27;
	add.s64 	%rd30, %rd3, %rd27;
	ld.global.nc.f32 	%f10, [%rd28];
	ld.global.nc.f32 	%f11, [%rd29];
	mul.f32 	%f77, %f11, %f11;
	fma.rn.f32 	%f78, %f10, %f10, %f77;
	ld.global.nc.f32 	%f12, [%rd30];
	fma.rn.f32 	%f13, %f12, %f12, %f78;
	add.s64 	%rd31, %rd2, %rd26;
	ld.global.nc.u8 	%rs9, [%rd31];
	setp.gt.u16	%p14, %rs9, %rs1;
	cvt.u32.u16	%r62, %rs9;
	and.b32  	%r63, %r62, 255;
	selp.b32	%r64, %r7, %r63, %p14;
	selp.b32	%r65, %r63, %r7, %p14;
	add.s32 	%r66, %r65, 1;
	mul.lo.s32 	%r67, %r66, %r65;
	shr.u32 	%r68, %r67, 1;
	add.s32 	%r69, %r68, %r64;
	mul.wide.s32 	%rd32, %r69, 4;
	add.s64 	%rd33, %rd1, %rd32;
	ld.global.nc.f32 	%f79, [%rd33];
	setp.eq.f32	%p15, %f79, 0f00000000;
	@%p15 bra 	BB0_12;

	setp.eq.f32	%p16, %f13, 0f00000000;
	selp.f32	%f80, %f1, %f10, %p16;
	selp.f32	%f81, %f2, %f11, %p16;
	mul.f32 	%f82, %f2, %f81;
	fma.rn.f32 	%f83, %f1, %f80, %f82;
	selp.f32	%f84, %f3, %f12, %p16;
	fma.rn.f32 	%f85, %f3, %f84, %f83;
	abs.f32 	%f86, %f85;
	mov.f32 	%f87, 0f3F800000;
	sub.f32 	%f88, %f87, %f86;
	mul.f32 	%f89, %f88, 0f3F000000;
	sqrt.rn.f32 	%f90, %f89;
	setp.gt.f32	%p17, %f86, 0f3F11EB85;
	selp.f32	%f91, %f90, %f86, %p17;
	mul.f32 	%f92, %f91, %f91;
	mov.f32 	%f93, 0f3C94D2E9;
	mov.f32 	%f94, 0f3D53F941;
	fma.rn.f32 	%f95, %f94, %f92, %f93;
	mov.f32 	%f96, 0f3D3F841F;
	fma.rn.f32 	%f97, %f95, %f92, %f96;
	mov.f32 	%f98, 0f3D994929;
	fma.rn.f32 	%f99, %f97, %f92, %f98;
	mov.f32 	%f100, 0f3E2AAB94;
	fma.rn.f32 	%f101, %f99, %f92, %f100;
	mul.f32 	%f102, %f92, %f101;
	fma.rn.f32 	%f103, %f102, %f91, %f91;
	add.f32 	%f104, %f103, %f103;
	mov.f32 	%f105, 0f3FC90FDB;
	sub.f32 	%f106, %f105, %f103;
	selp.f32	%f107, %f104, %f106, %p17;
	setp.lt.f32	%p18, %f85, 0f00000000;
	mov.f32 	%f108, 0f40490FDB;
	sub.f32 	%f109, %f108, %f107;
	selp.f32	%f110, %f109, %f107, %p18;
	max.f32 	%f250, %f250, %f110;

BB0_12:
	and.b16  	%rs3, %rs5, 2;
	setp.eq.s16	%p19, %rs3, 0;
	add.s32 	%r16, %r2, -1;
	@%p19 bra 	BB0_14;

	rem.s32 	%r70, %r16, %r33;
	add.s32 	%r71, %r70, %r33;
	rem.s32 	%r124, %r71, %r33;
	bra.uni 	BB0_15;

BB0_14:
	mov.u32 	%r72, 0;
	max.s32 	%r124, %r16, %r72;

BB0_15:
	add.s32 	%r73, %r124, %r4;
	mad.lo.s32 	%r74, %r73, %r32, %r1;
	cvt.s64.s32	%rd34, %r74;
	mul.wide.s32 	%rd35, %r74, 4;
	add.s64 	%rd36, %rd5, %rd35;
	add.s64 	%rd37, %rd4, %rd35;
	add.s64 	%rd38, %rd3, %rd35;
	ld.global.nc.f32 	%f16, [%rd36];
	ld.global.nc.f32 	%f17, [%rd37];
	mul.f32 	%f111, %f17, %f17;
	fma.rn.f32 	%f112, %f16, %f16, %f111;
	ld.global.nc.f32 	%f18, [%rd38];
	fma.rn.f32 	%f19, %f18, %f18, %f112;
	add.s64 	%rd39, %rd2, %rd34;
	ld.global.nc.u8 	%rs12, [%rd39];
	setp.gt.u16	%p20, %rs12, %rs1;
	cvt.u32.u16	%r75, %rs12;
	and.b32  	%r76, %r75, 255;
	selp.b32	%r77, %r7, %r76, %p20;
	selp.b32	%r78, %r76, %r7, %p20;
	add.s32 	%r79, %r78, 1;
	mul.lo.s32 	%r80, %r79, %r78;
	shr.u32 	%r81, %r80, 1;
	add.s32 	%r82, %r81, %r77;
	mul.wide.s32 	%rd40, %r82, 4;
	add.s64 	%rd41, %rd1, %rd40;
	ld.global.nc.f32 	%f113, [%rd41];
	setp.eq.f32	%p21, %f113, 0f00000000;
	@%p21 bra 	BB0_17;

	setp.eq.f32	%p22, %f19, 0f00000000;
	selp.f32	%f114, %f1, %f16, %p22;
	selp.f32	%f115, %f2, %f17, %p22;
	mul.f32 	%f116, %f2, %f115;
	fma.rn.f32 	%f117, %f1, %f114, %f116;
	selp.f32	%f118, %f3, %f18, %p22;
	fma.rn.f32 	%f119, %f3, %f118, %f117;
	abs.f32 	%f120, %f119;
	mov.f32 	%f121, 0f3F800000;
	sub.f32 	%f122, %f121, %f120;
	mul.f32 	%f123, %f122, 0f3F000000;
	sqrt.rn.f32 	%f124, %f123;
	setp.gt.f32	%p23, %f120, 0f3F11EB85;
	selp.f32	%f125, %f124, %f120, %p23;
	mul.f32 	%f126, %f125, %f125;
	mov.f32 	%f127, 0f3C94D2E9;
	mov.f32 	%f128, 0f3D53F941;
	fma.rn.f32 	%f129, %f128, %f126, %f127;
	mov.f32 	%f130, 0f3D3F841F;
	fma.rn.f32 	%f131, %f129, %f126, %f130;
	mov.f32 	%f132, 0f3D994929;
	fma.rn.f32 	%f133, %f131, %f126, %f132;
	mov.f32 	%f134, 0f3E2AAB94;
	fma.rn.f32 	%f135, %f133, %f126, %f134;
	mul.f32 	%f136, %f126, %f135;
	fma.rn.f32 	%f137, %f136, %f125, %f125;
	add.f32 	%f138, %f137, %f137;
	mov.f32 	%f139, 0f3FC90FDB;
	sub.f32 	%f140, %f139, %f137;
	selp.f32	%f141, %f138, %f140, %p23;
	setp.lt.f32	%p24, %f119, 0f00000000;
	mov.f32 	%f142, 0f40490FDB;
	sub.f32 	%f143, %f142, %f141;
	selp.f32	%f144, %f143, %f141, %p24;
	max.f32 	%f250, %f250, %f144;

BB0_17:
	add.s32 	%r20, %r2, 1;
	@%p19 bra 	BB0_19;

	rem.s32 	%r83, %r20, %r33;
	add.s32 	%r84, %r83, %r33;
	rem.s32 	%r125, %r84, %r33;
	bra.uni 	BB0_20;

BB0_19:
	add.s32 	%r85, %r33, -1;
	min.s32 	%r125, %r20, %r85;

BB0_20:
	add.s32 	%r86, %r125, %r4;
	mad.lo.s32 	%r87, %r86, %r32, %r1;
	cvt.s64.s32	%rd42, %r87;
	mul.wide.s32 	%rd43, %r87, 4;
	add.s64 	%rd44, %rd5, %rd43;
	add.s64 	%rd45, %rd4, %rd43;
	add.s64 	%rd46, %rd3, %rd43;
	ld.global.nc.f32 	%f22, [%rd44];
	ld.global.nc.f32 	%f23, [%rd45];
	mul.f32 	%f145, %f23, %f23;
	fma.rn.f32 	%f146, %f22, %f22, %f145;
	ld.global.nc.f32 	%f24, [%rd46];
	fma.rn.f32 	%f25, %f24, %f24, %f146;
	add.s64 	%rd47, %rd2, %rd42;
	ld.global.nc.u8 	%rs16, [%rd47];
	setp.gt.u16	%p26, %rs16, %rs1;
	cvt.u32.u16	%r88, %rs16;
	and.b32  	%r89, %r88, 255;
	selp.b32	%r90, %r7, %r89, %p26;
	selp.b32	%r91, %r89, %r7, %p26;
	add.s32 	%r92, %r91, 1;
	mul.lo.s32 	%r93, %r92, %r91;
	shr.u32 	%r94, %r93, 1;
	add.s32 	%r95, %r94, %r90;
	mul.wide.s32 	%rd48, %r95, 4;
	add.s64 	%rd49, %rd1, %rd48;
	ld.global.nc.f32 	%f147, [%rd49];
	setp.eq.f32	%p27, %f147, 0f00000000;
	@%p27 bra 	BB0_22;

	setp.eq.f32	%p28, %f25, 0f00000000;
	selp.f32	%f148, %f1, %f22, %p28;
	selp.f32	%f149, %f2, %f23, %p28;
	mul.f32 	%f150, %f2, %f149;
	fma.rn.f32 	%f151, %f1, %f148, %f150;
	selp.f32	%f152, %f3, %f24, %p28;
	fma.rn.f32 	%f153, %f3, %f152, %f151;
	abs.f32 	%f154, %f153;
	mov.f32 	%f155, 0f3F800000;
	sub.f32 	%f156, %f155, %f154;
	mul.f32 	%f157, %f156, 0f3F000000;
	sqrt.rn.f32 	%f158, %f157;
	setp.gt.f32	%p29, %f154, 0f3F11EB85;
	selp.f32	%f159, %f158, %f154, %p29;
	mul.f32 	%f160, %f159, %f159;
	mov.f32 	%f161, 0f3C94D2E9;
	mov.f32 	%f162, 0f3D53F941;
	fma.rn.f32 	%f163, %f162, %f160, %f161;
	mov.f32 	%f164, 0f3D3F841F;
	fma.rn.f32 	%f165, %f163, %f160, %f164;
	mov.f32 	%f166, 0f3D994929;
	fma.rn.f32 	%f167, %f165, %f160, %f166;
	mov.f32 	%f168, 0f3E2AAB94;
	fma.rn.f32 	%f169, %f167, %f160, %f168;
	mul.f32 	%f170, %f160, %f169;
	fma.rn.f32 	%f171, %f170, %f159, %f159;
	add.f32 	%f172, %f171, %f171;
	mov.f32 	%f173, 0f3FC90FDB;
	sub.f32 	%f174, %f173, %f171;
	selp.f32	%f175, %f172, %f174, %p29;
	setp.lt.f32	%p30, %f153, 0f00000000;
	mov.f32 	%f176, 0f40490FDB;
	sub.f32 	%f177, %f176, %f175;
	selp.f32	%f178, %f177, %f175, %p30;
	max.f32 	%f250, %f250, %f178;

BB0_22:
	setp.eq.s32	%p31, %r34, 1;
	@%p31 bra 	BB0_33;

	and.b16  	%rs4, %rs5, 4;
	setp.eq.s16	%p32, %rs4, 0;
	add.s32 	%r24, %r3, -1;
	@%p32 bra 	BB0_25;

	rem.s32 	%r96, %r24, %r34;
	add.s32 	%r97, %r96, %r34;
	rem.s32 	%r126, %r97, %r34;
	bra.uni 	BB0_26;

BB0_25:
	mov.u32 	%r98, 0;
	max.s32 	%r126, %r24, %r98;

BB0_26:
	mad.lo.s32 	%r99, %r126, %r33, %r2;
	mad.lo.s32 	%r100, %r99, %r32, %r1;
	cvt.s64.s32	%rd50, %r100;
	mul.wide.s32 	%rd51, %r100, 4;
	add.s64 	%rd52, %rd5, %rd51;
	add.s64 	%rd53, %rd4, %rd51;
	add.s64 	%rd54, %rd3, %rd51;
	ld.global.nc.f32 	%f179, [%rd52];
	ld.global.nc.f32 	%f180, [%rd53];
	mul.f32 	%f181, %f180, %f180;
	fma.rn.f32 	%f182, %f179, %f179, %f181;
	ld.global.nc.f32 	%f183, [%rd54];
	fma.rn.f32 	%f184, %f183, %f183, %f182;
	setp.eq.f32	%p33, %f184, 0f00000000;
	selp.f32	%f28, %f1, %f179, %p33;
	selp.f32	%f29, %f2, %f180, %p33;
	selp.f32	%f30, %f3, %f183, %p33;
	add.s64 	%rd55, %rd2, %rd50;
	ld.global.nc.u8 	%rs19, [%rd55];
	setp.gt.u16	%p34, %rs19, %rs1;
	cvt.u32.u16	%r101, %rs19;
	and.b32  	%r102, %r101, 255;
	selp.b32	%r103, %r7, %r102, %p34;
	selp.b32	%r104, %r102, %r7, %p34;
	add.s32 	%r105, %r104, 1;
	mul.lo.s32 	%r106, %r105, %r104;
	shr.u32 	%r107, %r106, 1;
	add.s32 	%r108, %r107, %r103;
	mul.wide.s32 	%rd56, %r108, 4;
	add.s64 	%rd57, %rd1, %rd56;
	ld.global.nc.f32 	%f185, [%rd57];
	setp.eq.f32	%p35, %f185, 0f00000000;
	@%p35 bra 	BB0_28;

	mul.f32 	%f186, %f2, %f29;
	fma.rn.f32 	%f187, %f1, %f28, %f186;
	fma.rn.f32 	%f188, %f3, %f30, %f187;
	abs.f32 	%f189, %f188;
	mov.f32 	%f190, 0f3F800000;
	sub.f32 	%f191, %f190, %f189;
	mul.f32 	%f192, %f191, 0f3F000000;
	sqrt.rn.f32 	%f193, %f192;
	setp.gt.f32	%p36, %f189, 0f3F11EB85;
	selp.f32	%f194, %f193, %f189, %p36;
	mul.f32 	%f195, %f194, %f194;
	mov.f32 	%f196, 0f3C94D2E9;
	mov.f32 	%f197, 0f3D53F941;
	fma.rn.f32 	%f198, %f197, %f195, %f196;
	mov.f32 	%f199, 0f3D3F841F;
	fma.rn.f32 	%f200, %f198, %f195, %f199;
	mov.f32 	%f201, 0f3D994929;
	fma.rn.f32 	%f202, %f200, %f195, %f201;
	mov.f32 	%f203, 0f3E2AAB94;
	fma.rn.f32 	%f204, %f202, %f195, %f203;
	mul.f32 	%f205, %f195, %f204;
	fma.rn.f32 	%f206, %f205, %f194, %f194;
	add.f32 	%f207, %f206, %f206;
	mov.f32 	%f208, 0f3FC90FDB;
	sub.f32 	%f209, %f208, %f206;
	selp.f32	%f210, %f207, %f209, %p36;
	setp.lt.f32	%p37, %f188, 0f00000000;
	mov.f32 	%f211, 0f40490FDB;
	sub.f32 	%f212, %f211, %f210;
	selp.f32	%f213, %f212, %f210, %p37;
	max.f32 	%f250, %f250, %f213;

BB0_28:
	add.s32 	%r28, %r3, 1;
	@%p32 bra 	BB0_30;

	rem.s32 	%r109, %r28, %r34;
	add.s32 	%r110, %r109, %r34;
	rem.s32 	%r127, %r110, %r34;
	bra.uni 	BB0_31;

BB0_30:
	add.s32 	%r111, %r34, -1;
	min.s32 	%r127, %r28, %r111;

BB0_31:
	mad.lo.s32 	%r112, %r127, %r33, %r2;
	mad.lo.s32 	%r113, %r112, %r32, %r1;
	cvt.s64.s32	%rd58, %r113;
	mul.wide.s32 	%rd59, %r113, 4;
	add.s64 	%rd60, %rd5, %rd59;
	add.s64 	%rd61, %rd4, %rd59;
	add.s64 	%rd62, %rd3, %rd59;
	ld.global.nc.f32 	%f214, [%rd60];
	ld.global.nc.f32 	%f215, [%rd61];
	mul.f32 	%f216, %f215, %f215;
	fma.rn.f32 	%f217, %f214, %f214, %f216;
	ld.global.nc.f32 	%f218, [%rd62];
	fma.rn.f32 	%f219, %f218, %f218, %f217;
	setp.eq.f32	%p39, %f219, 0f00000000;
	selp.f32	%f33, %f1, %f214, %p39;
	selp.f32	%f34, %f2, %f215, %p39;
	selp.f32	%f35, %f3, %f218, %p39;
	add.s64 	%rd63, %rd2, %rd58;
	ld.global.nc.u8 	%rs23, [%rd63];
	setp.gt.u16	%p40, %rs23, %rs1;
	cvt.u32.u16	%r114, %rs23;
	and.b32  	%r115, %r114, 255;
	selp.b32	%r116, %r7, %r115, %p40;
	selp.b32	%r117, %r115, %r7, %p40;
	add.s32 	%r118, %r117, 1;
	mul.lo.s32 	%r119, %r118, %r117;
	shr.u32 	%r120, %r119, 1;
	add.s32 	%r121, %r120, %r116;
	mul.wide.s32 	%rd64, %r121, 4;
	add.s64 	%rd65, %rd1, %rd64;
	ld.global.nc.f32 	%f220, [%rd65];
	setp.eq.f32	%p41, %f220, 0f00000000;
	@%p41 bra 	BB0_33;

	mul.f32 	%f221, %f2, %f34;
	fma.rn.f32 	%f222, %f1, %f33, %f221;
	fma.rn.f32 	%f223, %f3, %f35, %f222;
	abs.f32 	%f224, %f223;
	mov.f32 	%f225, 0f3F800000;
	sub.f32 	%f226, %f225, %f224;
	mul.f32 	%f227, %f226, 0f3F000000;
	sqrt.rn.f32 	%f228, %f227;
	setp.gt.f32	%p42, %f224, 0f3F11EB85;
	selp.f32	%f229, %f228, %f224, %p42;
	mul.f32 	%f230, %f229, %f229;
	mov.f32 	%f231, 0f3C94D2E9;
	mov.f32 	%f232, 0f3D53F941;
	fma.rn.f32 	%f233, %f232, %f230, %f231;
	mov.f32 	%f234, 0f3D3F841F;
	fma.rn.f32 	%f235, %f233, %f230, %f234;
	mov.f32 	%f236, 0f3D994929;
	fma.rn.f32 	%f237, %f235, %f230, %f236;
	mov.f32 	%f238, 0f3E2AAB94;
	fma.rn.f32 	%f239, %f237, %f230, %f238;
	mul.f32 	%f240, %f230, %f239;
	fma.rn.f32 	%f241, %f240, %f229, %f229;
	add.f32 	%f242, %f241, %f241;
	mov.f32 	%f243, 0f3FC90FDB;
	sub.f32 	%f244, %f243, %f241;
	selp.f32	%f245, %f242, %f244, %p42;
	setp.lt.f32	%p43, %f223, 0f00000000;
	mov.f32 	%f246, 0f40490FDB;
	sub.f32 	%f247, %f246, %f245;
	selp.f32	%f248, %f247, %f245, %p43;
	max.f32 	%f250, %f250, %f248;

BB0_33:
	cvta.to.global.u64 	%rd66, %rd6;
	add.s64 	%rd68, %rd66, %rd12;
	st.global.f32 	[%rd68], %f250;

BB0_34:
	ret;
}


`
 )