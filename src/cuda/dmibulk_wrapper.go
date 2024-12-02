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

// CUDA handle for adddmibulk kernel
var adddmibulk_code cu.Function

// Stores the arguments for adddmibulk kernel invocation
type adddmibulk_args_t struct {
	arg_Hx      unsafe.Pointer
	arg_Hy      unsafe.Pointer
	arg_Hz      unsafe.Pointer
	arg_mx      unsafe.Pointer
	arg_my      unsafe.Pointer
	arg_mz      unsafe.Pointer
	arg_Ms_     unsafe.Pointer
	arg_Ms_mul  float32
	arg_aLUT2d  unsafe.Pointer
	arg_DLUT2d  unsafe.Pointer
	arg_regions unsafe.Pointer
	arg_cx      float32
	arg_cy      float32
	arg_cz      float32
	arg_Nx      int
	arg_Ny      int
	arg_Nz      int
	arg_PBC     byte
	arg_OpenBC  byte
	argptr      [19]unsafe.Pointer
	sync.Mutex
}

// Stores the arguments for adddmibulk kernel invocation
var adddmibulk_args adddmibulk_args_t

func init() {
	// CUDA driver kernel call wants pointers to arguments, set them up once.
	adddmibulk_args.argptr[0] = unsafe.Pointer(&adddmibulk_args.arg_Hx)
	adddmibulk_args.argptr[1] = unsafe.Pointer(&adddmibulk_args.arg_Hy)
	adddmibulk_args.argptr[2] = unsafe.Pointer(&adddmibulk_args.arg_Hz)
	adddmibulk_args.argptr[3] = unsafe.Pointer(&adddmibulk_args.arg_mx)
	adddmibulk_args.argptr[4] = unsafe.Pointer(&adddmibulk_args.arg_my)
	adddmibulk_args.argptr[5] = unsafe.Pointer(&adddmibulk_args.arg_mz)
	adddmibulk_args.argptr[6] = unsafe.Pointer(&adddmibulk_args.arg_Ms_)
	adddmibulk_args.argptr[7] = unsafe.Pointer(&adddmibulk_args.arg_Ms_mul)
	adddmibulk_args.argptr[8] = unsafe.Pointer(&adddmibulk_args.arg_aLUT2d)
	adddmibulk_args.argptr[9] = unsafe.Pointer(&adddmibulk_args.arg_DLUT2d)
	adddmibulk_args.argptr[10] = unsafe.Pointer(&adddmibulk_args.arg_regions)
	adddmibulk_args.argptr[11] = unsafe.Pointer(&adddmibulk_args.arg_cx)
	adddmibulk_args.argptr[12] = unsafe.Pointer(&adddmibulk_args.arg_cy)
	adddmibulk_args.argptr[13] = unsafe.Pointer(&adddmibulk_args.arg_cz)
	adddmibulk_args.argptr[14] = unsafe.Pointer(&adddmibulk_args.arg_Nx)
	adddmibulk_args.argptr[15] = unsafe.Pointer(&adddmibulk_args.arg_Ny)
	adddmibulk_args.argptr[16] = unsafe.Pointer(&adddmibulk_args.arg_Nz)
	adddmibulk_args.argptr[17] = unsafe.Pointer(&adddmibulk_args.arg_PBC)
	adddmibulk_args.argptr[18] = unsafe.Pointer(&adddmibulk_args.arg_OpenBC)
}

// Wrapper for adddmibulk CUDA kernel, asynchronous.
func k_adddmibulk_async(Hx unsafe.Pointer, Hy unsafe.Pointer, Hz unsafe.Pointer, mx unsafe.Pointer, my unsafe.Pointer, mz unsafe.Pointer, Ms_ unsafe.Pointer, Ms_mul float32, aLUT2d unsafe.Pointer, DLUT2d unsafe.Pointer, regions unsafe.Pointer, cx float32, cy float32, cz float32, Nx int, Ny int, Nz int, PBC byte, OpenBC byte, cfg *config) {
	if Synchronous { // debug
		Sync()
		timer.Start("adddmibulk")
	}

	adddmibulk_args.Lock()
	defer adddmibulk_args.Unlock()

	if adddmibulk_code == 0 {
		adddmibulk_code = fatbinLoad(adddmibulk_map, "adddmibulk")
	}

	adddmibulk_args.arg_Hx = Hx
	adddmibulk_args.arg_Hy = Hy
	adddmibulk_args.arg_Hz = Hz
	adddmibulk_args.arg_mx = mx
	adddmibulk_args.arg_my = my
	adddmibulk_args.arg_mz = mz
	adddmibulk_args.arg_Ms_ = Ms_
	adddmibulk_args.arg_Ms_mul = Ms_mul
	adddmibulk_args.arg_aLUT2d = aLUT2d
	adddmibulk_args.arg_DLUT2d = DLUT2d
	adddmibulk_args.arg_regions = regions
	adddmibulk_args.arg_cx = cx
	adddmibulk_args.arg_cy = cy
	adddmibulk_args.arg_cz = cz
	adddmibulk_args.arg_Nx = Nx
	adddmibulk_args.arg_Ny = Ny
	adddmibulk_args.arg_Nz = Nz
	adddmibulk_args.arg_PBC = PBC
	adddmibulk_args.arg_OpenBC = OpenBC

	args := adddmibulk_args.argptr[:]
	cu.LaunchKernel(adddmibulk_code, cfg.Grid.X, cfg.Grid.Y, cfg.Grid.Z, cfg.Block.X, cfg.Block.Y, cfg.Block.Z, 0, stream0, args)

	if Synchronous { // debug
		Sync()
		timer.Stop("adddmibulk")
	}
}

// maps compute capability on PTX code for adddmibulk kernel.
var adddmibulk_map = map[int]string{0: "",
	52: adddmibulk_ptx_52}

// adddmibulk PTX code for various compute capabilities.
const (
	adddmibulk_ptx_52 = `
.version 7.0
.target sm_52
.address_size 64

	// .globl	adddmibulk

.visible .entry adddmibulk(
	.param .u64 adddmibulk_param_0,
	.param .u64 adddmibulk_param_1,
	.param .u64 adddmibulk_param_2,
	.param .u64 adddmibulk_param_3,
	.param .u64 adddmibulk_param_4,
	.param .u64 adddmibulk_param_5,
	.param .u64 adddmibulk_param_6,
	.param .f32 adddmibulk_param_7,
	.param .u64 adddmibulk_param_8,
	.param .u64 adddmibulk_param_9,
	.param .u64 adddmibulk_param_10,
	.param .f32 adddmibulk_param_11,
	.param .f32 adddmibulk_param_12,
	.param .f32 adddmibulk_param_13,
	.param .u32 adddmibulk_param_14,
	.param .u32 adddmibulk_param_15,
	.param .u32 adddmibulk_param_16,
	.param .u8 adddmibulk_param_17,
	.param .u8 adddmibulk_param_18
)
{
	.reg .pred 	%p<70>;
	.reg .b16 	%rs<43>;
	.reg .f32 	%f<292>;
	.reg .b32 	%r<128>;
	.reg .b64 	%rd<87>;


	ld.param.u64 	%rd9, [adddmibulk_param_0];
	ld.param.u64 	%rd10, [adddmibulk_param_1];
	ld.param.u64 	%rd11, [adddmibulk_param_2];
	ld.param.u64 	%rd13, [adddmibulk_param_3];
	ld.param.u64 	%rd14, [adddmibulk_param_4];
	ld.param.u64 	%rd15, [adddmibulk_param_5];
	ld.param.u64 	%rd12, [adddmibulk_param_6];
	ld.param.f32 	%f290, [adddmibulk_param_7];
	ld.param.u64 	%rd16, [adddmibulk_param_8];
	ld.param.u64 	%rd17, [adddmibulk_param_9];
	ld.param.u64 	%rd18, [adddmibulk_param_10];
	ld.param.f32 	%f87, [adddmibulk_param_11];
	ld.param.f32 	%f88, [adddmibulk_param_12];
	ld.param.f32 	%f89, [adddmibulk_param_13];
	ld.param.u32 	%r43, [adddmibulk_param_14];
	ld.param.u32 	%r44, [adddmibulk_param_15];
	ld.param.u32 	%r45, [adddmibulk_param_16];
	ld.param.u8 	%rs18, [adddmibulk_param_18];
	ld.param.u8 	%rs17, [adddmibulk_param_17];
	cvta.to.global.u64 	%rd1, %rd17;
	cvta.to.global.u64 	%rd2, %rd16;
	cvta.to.global.u64 	%rd3, %rd18;
	cvta.to.global.u64 	%rd4, %rd15;
	cvta.to.global.u64 	%rd5, %rd14;
	cvta.to.global.u64 	%rd6, %rd13;
	mov.u32 	%r46, %ntid.x;
	mov.u32 	%r47, %ctaid.x;
	mov.u32 	%r48, %tid.x;
	mad.lo.s32 	%r1, %r46, %r47, %r48;
	mov.u32 	%r49, %ntid.y;
	mov.u32 	%r50, %ctaid.y;
	mov.u32 	%r51, %tid.y;
	mad.lo.s32 	%r2, %r49, %r50, %r51;
	mov.u32 	%r52, %ntid.z;
	mov.u32 	%r53, %ctaid.z;
	mov.u32 	%r54, %tid.z;
	mad.lo.s32 	%r3, %r52, %r53, %r54;
	setp.ge.s32	%p1, %r2, %r44;
	setp.ge.s32	%p2, %r1, %r43;
	or.pred  	%p3, %p1, %p2;
	setp.ge.s32	%p4, %r3, %r45;
	or.pred  	%p5, %p3, %p4;
	@%p5 bra 	BB0_62;

	mul.lo.s32 	%r4, %r3, %r44;
	add.s32 	%r55, %r4, %r2;
	mul.lo.s32 	%r5, %r55, %r43;
	add.s32 	%r6, %r5, %r1;
	mul.wide.s32 	%rd19, %r6, 4;
	add.s64 	%rd20, %rd6, %rd19;
	cvt.s64.s32	%rd21, %r6;
	add.s64 	%rd22, %rd5, %rd19;
	add.s64 	%rd23, %rd4, %rd19;
	add.s64 	%rd24, %rd3, %rd21;
	ld.global.nc.u8 	%rs1, [%rd24];
	cvt.u32.u16	%r56, %rs1;
	and.b32  	%r7, %r56, 255;
	ld.global.nc.f32 	%f1, [%rd20];
	ld.global.nc.f32 	%f2, [%rd22];
	mul.f32 	%f90, %f2, %f2;
	fma.rn.f32 	%f91, %f1, %f1, %f90;
	ld.global.nc.f32 	%f3, [%rd23];
	fma.rn.f32 	%f92, %f3, %f3, %f91;
	setp.eq.f32	%p6, %f92, 0f00000000;
	@%p6 bra 	BB0_62;

	and.b16  	%rs2, %rs17, 1;
	setp.eq.s16	%p7, %rs2, 0;
	add.s32 	%r8, %r1, -1;
	@%p7 bra 	BB0_4;

	rem.s32 	%r57, %r8, %r43;
	add.s32 	%r58, %r57, %r43;
	rem.s32 	%r122, %r58, %r43;
	bra.uni 	BB0_5;

BB0_4:
	mov.u32 	%r59, 0;
	max.s32 	%r122, %r8, %r59;

BB0_5:
	add.s32 	%r12, %r122, %r5;
	setp.lt.s32	%p9, %r8, 0;
	mov.f32 	%f254, 0f00000000;
	and.pred  	%p10, %p9, %p7;
	mov.f32 	%f255, %f254;
	mov.f32 	%f256, %f254;
	@%p10 bra 	BB0_7;

	mul.wide.s32 	%rd25, %r12, 4;
	add.s64 	%rd26, %rd6, %rd25;
	ld.global.nc.f32 	%f254, [%rd26];
	add.s64 	%rd27, %rd5, %rd25;
	ld.global.nc.f32 	%f255, [%rd27];
	add.s64 	%rd28, %rd4, %rd25;
	ld.global.nc.f32 	%f256, [%rd28];

BB0_7:
	mul.f32 	%f96, %f255, %f255;
	fma.rn.f32 	%f97, %f254, %f254, %f96;
	fma.rn.f32 	%f10, %f256, %f256, %f97;
	setp.eq.f32	%p11, %f10, 0f00000000;
	mov.u16 	%rs37, %rs1;
	@%p11 bra 	BB0_9;

	cvt.s64.s32	%rd29, %r12;
	add.s64 	%rd30, %rd3, %rd29;
	ld.global.nc.u8 	%rs37, [%rd30];

BB0_9:
	setp.gt.u16	%p12, %rs37, %rs1;
	cvt.u32.u16	%r60, %rs37;
	and.b32  	%r61, %r60, 255;
	selp.b32	%r62, %r7, %r61, %p12;
	selp.b32	%r63, %r61, %r7, %p12;
	add.s32 	%r64, %r63, 1;
	mul.lo.s32 	%r65, %r64, %r63;
	shr.u32 	%r66, %r65, 1;
	add.s32 	%r13, %r66, %r62;
	setp.ne.s16	%p13, %rs18, 0;
	mov.f32 	%f263, 0f00000000;
	and.pred  	%p15, %p11, %p13;
	mov.f32 	%f264, %f263;
	mov.f32 	%f265, %f263;
	@%p15 bra 	BB0_11;

	mul.wide.s32 	%rd31, %r13, 4;
	add.s64 	%rd32, %rd2, %rd31;
	ld.global.nc.f32 	%f101, [%rd32];
	add.f32 	%f102, %f101, %f101;
	add.s64 	%rd33, %rd1, %rd31;
	ld.global.nc.f32 	%f103, [%rd33];
	div.rn.f32 	%f104, %f103, %f102;
	mul.f32 	%f105, %f104, %f87;
	fma.rn.f32 	%f106, %f3, %f105, %f2;
	mul.f32 	%f107, %f2, %f105;
	sub.f32 	%f108, %f3, %f107;
	selp.f32	%f109, %f1, %f254, %p11;
	selp.f32	%f110, %f106, %f255, %p11;
	selp.f32	%f111, %f108, %f256, %p11;
	mul.f32 	%f112, %f87, %f87;
	div.rn.f32 	%f113, %f102, %f112;
	sub.f32 	%f114, %f109, %f1;
	sub.f32 	%f115, %f110, %f2;
	sub.f32 	%f116, %f111, %f3;
	fma.rn.f32 	%f263, %f114, %f113, 0f00000000;
	fma.rn.f32 	%f117, %f115, %f113, 0f00000000;
	fma.rn.f32 	%f118, %f116, %f113, 0f00000000;
	div.rn.f32 	%f119, %f103, %f87;
	mul.f32 	%f120, %f111, %f119;
	sub.f32 	%f264, %f117, %f120;
	fma.rn.f32 	%f265, %f110, %f119, %f118;

BB0_11:
	add.s32 	%r14, %r1, 1;
	@%p7 bra 	BB0_13;

	rem.s32 	%r67, %r14, %r43;
	add.s32 	%r68, %r67, %r43;
	rem.s32 	%r123, %r68, %r43;
	bra.uni 	BB0_14;

BB0_13:
	add.s32 	%r69, %r43, -1;
	min.s32 	%r123, %r14, %r69;

BB0_14:
	add.s32 	%r18, %r123, %r5;
	setp.ge.s32	%p18, %r14, %r43;
	mov.f32 	%f260, 0f00000000;
	and.pred  	%p20, %p18, %p7;
	mov.f32 	%f261, %f260;
	mov.f32 	%f262, %f260;
	@%p20 bra 	BB0_16;

	mul.wide.s32 	%rd34, %r18, 4;
	add.s64 	%rd35, %rd6, %rd34;
	ld.global.nc.f32 	%f260, [%rd35];
	add.s64 	%rd36, %rd5, %rd34;
	ld.global.nc.f32 	%f261, [%rd36];
	add.s64 	%rd37, %rd4, %rd34;
	ld.global.nc.f32 	%f262, [%rd37];

BB0_16:
	mul.f32 	%f124, %f261, %f261;
	fma.rn.f32 	%f125, %f260, %f260, %f124;
	fma.rn.f32 	%f23, %f262, %f262, %f125;
	setp.eq.f32	%p21, %f23, 0f00000000;
	mov.u16 	%rs38, %rs1;
	@%p21 bra 	BB0_18;

	cvt.s64.s32	%rd38, %r18;
	add.s64 	%rd39, %rd3, %rd38;
	ld.global.nc.u8 	%rs38, [%rd39];

BB0_18:
	setp.gt.u16	%p22, %rs38, %rs1;
	cvt.u32.u16	%r70, %rs38;
	and.b32  	%r71, %r70, 255;
	selp.b32	%r72, %r7, %r71, %p22;
	selp.b32	%r73, %r71, %r7, %p22;
	add.s32 	%r74, %r73, 1;
	mul.lo.s32 	%r75, %r74, %r73;
	shr.u32 	%r76, %r75, 1;
	add.s32 	%r19, %r76, %r72;
	and.pred  	%p25, %p21, %p13;
	@%p25 bra 	BB0_20;

	mul.wide.s32 	%rd40, %r19, 4;
	add.s64 	%rd41, %rd2, %rd40;
	ld.global.nc.f32 	%f126, [%rd41];
	add.f32 	%f127, %f126, %f126;
	add.s64 	%rd42, %rd1, %rd40;
	ld.global.nc.f32 	%f128, [%rd42];
	div.rn.f32 	%f129, %f128, %f127;
	mul.f32 	%f130, %f129, %f87;
	mul.f32 	%f131, %f3, %f130;
	sub.f32 	%f132, %f2, %f131;
	fma.rn.f32 	%f133, %f2, %f130, %f3;
	selp.f32	%f134, %f1, %f260, %p21;
	selp.f32	%f135, %f132, %f261, %p21;
	selp.f32	%f136, %f133, %f262, %p21;
	mul.f32 	%f137, %f87, %f87;
	div.rn.f32 	%f138, %f127, %f137;
	sub.f32 	%f139, %f134, %f1;
	sub.f32 	%f140, %f135, %f2;
	sub.f32 	%f141, %f136, %f3;
	fma.rn.f32 	%f263, %f139, %f138, %f263;
	fma.rn.f32 	%f142, %f140, %f138, %f264;
	fma.rn.f32 	%f143, %f141, %f138, %f265;
	div.rn.f32 	%f144, %f128, %f87;
	fma.rn.f32 	%f264, %f136, %f144, %f142;
	mul.f32 	%f145, %f135, %f144;
	sub.f32 	%f265, %f143, %f145;

BB0_20:
	and.b16  	%rs7, %rs17, 2;
	setp.eq.s16	%p27, %rs7, 0;
	add.s32 	%r20, %r2, -1;
	@%p27 bra 	BB0_22;

	rem.s32 	%r77, %r20, %r44;
	add.s32 	%r78, %r77, %r44;
	rem.s32 	%r124, %r78, %r44;
	bra.uni 	BB0_23;

BB0_22:
	mov.u32 	%r79, 0;
	max.s32 	%r124, %r20, %r79;

BB0_23:
	add.s32 	%r80, %r124, %r4;
	mad.lo.s32 	%r24, %r80, %r43, %r1;
	setp.lt.s32	%p29, %r20, 0;
	mov.f32 	%f266, 0f00000000;
	and.pred  	%p30, %p29, %p27;
	mov.f32 	%f267, %f266;
	mov.f32 	%f268, %f266;
	@%p30 bra 	BB0_25;

	mul.wide.s32 	%rd43, %r24, 4;
	add.s64 	%rd44, %rd6, %rd43;
	ld.global.nc.f32 	%f266, [%rd44];
	add.s64 	%rd45, %rd5, %rd43;
	ld.global.nc.f32 	%f267, [%rd45];
	add.s64 	%rd46, %rd4, %rd43;
	ld.global.nc.f32 	%f268, [%rd46];

BB0_25:
	mul.f32 	%f149, %f267, %f267;
	fma.rn.f32 	%f150, %f266, %f266, %f149;
	fma.rn.f32 	%f36, %f268, %f268, %f150;
	setp.eq.f32	%p31, %f36, 0f00000000;
	mov.u16 	%rs39, %rs1;
	@%p31 bra 	BB0_27;

	cvt.s64.s32	%rd47, %r24;
	add.s64 	%rd48, %rd3, %rd47;
	ld.global.nc.u8 	%rs39, [%rd48];

BB0_27:
	setp.gt.u16	%p32, %rs39, %rs1;
	cvt.u32.u16	%r81, %rs39;
	and.b32  	%r82, %r81, 255;
	selp.b32	%r83, %r7, %r82, %p32;
	selp.b32	%r84, %r82, %r7, %p32;
	add.s32 	%r85, %r84, 1;
	mul.lo.s32 	%r86, %r85, %r84;
	shr.u32 	%r87, %r86, 1;
	add.s32 	%r25, %r87, %r83;
	and.pred  	%p35, %p31, %p13;
	@%p35 bra 	BB0_29;

	mul.wide.s32 	%rd49, %r25, 4;
	add.s64 	%rd50, %rd2, %rd49;
	ld.global.nc.f32 	%f151, [%rd50];
	add.f32 	%f152, %f151, %f151;
	add.s64 	%rd51, %rd1, %rd49;
	ld.global.nc.f32 	%f153, [%rd51];
	div.rn.f32 	%f154, %f153, %f152;
	mul.f32 	%f155, %f154, %f88;
	mul.f32 	%f156, %f3, %f155;
	sub.f32 	%f157, %f1, %f156;
	fma.rn.f32 	%f158, %f1, %f155, %f3;
	selp.f32	%f159, %f157, %f266, %p31;
	selp.f32	%f160, %f2, %f267, %p31;
	selp.f32	%f161, %f158, %f268, %p31;
	mul.f32 	%f162, %f88, %f88;
	div.rn.f32 	%f163, %f152, %f162;
	sub.f32 	%f164, %f159, %f1;
	sub.f32 	%f165, %f160, %f2;
	sub.f32 	%f166, %f161, %f3;
	fma.rn.f32 	%f167, %f164, %f163, %f263;
	fma.rn.f32 	%f264, %f165, %f163, %f264;
	fma.rn.f32 	%f168, %f166, %f163, %f265;
	div.rn.f32 	%f169, %f153, %f88;
	fma.rn.f32 	%f263, %f161, %f169, %f167;
	mul.f32 	%f170, %f159, %f169;
	sub.f32 	%f265, %f168, %f170;

BB0_29:
	add.s32 	%r26, %r2, 1;
	@%p27 bra 	BB0_31;

	rem.s32 	%r88, %r26, %r44;
	add.s32 	%r89, %r88, %r44;
	rem.s32 	%r125, %r89, %r44;
	bra.uni 	BB0_32;

BB0_31:
	add.s32 	%r90, %r44, -1;
	min.s32 	%r125, %r26, %r90;

BB0_32:
	add.s32 	%r91, %r125, %r4;
	mad.lo.s32 	%r30, %r91, %r43, %r1;
	setp.ge.s32	%p38, %r26, %r44;
	mov.f32 	%f272, 0f00000000;
	and.pred  	%p40, %p38, %p27;
	mov.f32 	%f273, %f272;
	mov.f32 	%f274, %f272;
	@%p40 bra 	BB0_34;

	mul.wide.s32 	%rd52, %r30, 4;
	add.s64 	%rd53, %rd6, %rd52;
	ld.global.nc.f32 	%f272, [%rd53];
	add.s64 	%rd54, %rd5, %rd52;
	ld.global.nc.f32 	%f273, [%rd54];
	add.s64 	%rd55, %rd4, %rd52;
	ld.global.nc.f32 	%f274, [%rd55];

BB0_34:
	mul.f32 	%f174, %f273, %f273;
	fma.rn.f32 	%f175, %f272, %f272, %f174;
	fma.rn.f32 	%f49, %f274, %f274, %f175;
	setp.eq.f32	%p41, %f49, 0f00000000;
	mov.u16 	%rs40, %rs1;
	@%p41 bra 	BB0_36;

	cvt.s64.s32	%rd56, %r30;
	add.s64 	%rd57, %rd3, %rd56;
	ld.global.nc.u8 	%rs40, [%rd57];

BB0_36:
	setp.gt.u16	%p42, %rs40, %rs1;
	cvt.u32.u16	%r92, %rs40;
	and.b32  	%r93, %r92, 255;
	selp.b32	%r94, %r7, %r93, %p42;
	selp.b32	%r95, %r93, %r7, %p42;
	add.s32 	%r96, %r95, 1;
	mul.lo.s32 	%r97, %r96, %r95;
	shr.u32 	%r98, %r97, 1;
	add.s32 	%r31, %r98, %r94;
	and.pred  	%p45, %p41, %p13;
	@%p45 bra 	BB0_38;

	mul.wide.s32 	%rd58, %r31, 4;
	add.s64 	%rd59, %rd2, %rd58;
	ld.global.nc.f32 	%f176, [%rd59];
	add.f32 	%f177, %f176, %f176;
	add.s64 	%rd60, %rd1, %rd58;
	ld.global.nc.f32 	%f178, [%rd60];
	div.rn.f32 	%f179, %f178, %f177;
	mul.f32 	%f180, %f179, %f88;
	fma.rn.f32 	%f181, %f3, %f180, %f1;
	mul.f32 	%f182, %f1, %f180;
	sub.f32 	%f183, %f3, %f182;
	selp.f32	%f184, %f181, %f272, %p41;
	selp.f32	%f185, %f2, %f273, %p41;
	selp.f32	%f186, %f183, %f274, %p41;
	mul.f32 	%f187, %f88, %f88;
	div.rn.f32 	%f188, %f177, %f187;
	sub.f32 	%f189, %f184, %f1;
	sub.f32 	%f190, %f185, %f2;
	sub.f32 	%f191, %f186, %f3;
	fma.rn.f32 	%f192, %f189, %f188, %f263;
	fma.rn.f32 	%f264, %f190, %f188, %f264;
	fma.rn.f32 	%f193, %f191, %f188, %f265;
	div.rn.f32 	%f194, %f178, %f88;
	mul.f32 	%f195, %f186, %f194;
	sub.f32 	%f263, %f192, %f195;
	fma.rn.f32 	%f265, %f184, %f194, %f193;

BB0_38:
	setp.eq.s32	%p47, %r45, 1;
	@%p47 bra 	BB0_57;

	and.b16  	%rs12, %rs17, 4;
	setp.eq.s16	%p48, %rs12, 0;
	add.s32 	%r32, %r3, -1;
	@%p48 bra 	BB0_41;

	rem.s32 	%r99, %r32, %r45;
	add.s32 	%r100, %r99, %r45;
	rem.s32 	%r126, %r100, %r45;
	bra.uni 	BB0_42;

BB0_41:
	mov.u32 	%r101, 0;
	max.s32 	%r126, %r32, %r101;

BB0_42:
	mad.lo.s32 	%r102, %r126, %r44, %r2;
	mad.lo.s32 	%r36, %r102, %r43, %r1;
	setp.lt.s32	%p50, %r32, 0;
	mov.f32 	%f278, 0f00000000;
	and.pred  	%p51, %p50, %p48;
	mov.f32 	%f279, %f278;
	mov.f32 	%f280, %f278;
	@%p51 bra 	BB0_44;

	mul.wide.s32 	%rd61, %r36, 4;
	add.s64 	%rd62, %rd6, %rd61;
	ld.global.nc.f32 	%f278, [%rd62];
	add.s64 	%rd63, %rd5, %rd61;
	ld.global.nc.f32 	%f279, [%rd63];
	add.s64 	%rd64, %rd4, %rd61;
	ld.global.nc.f32 	%f280, [%rd64];

BB0_44:
	mul.f32 	%f199, %f279, %f279;
	fma.rn.f32 	%f200, %f278, %f278, %f199;
	fma.rn.f32 	%f62, %f280, %f280, %f200;
	setp.eq.f32	%p52, %f62, 0f00000000;
	mov.u16 	%rs41, %rs1;
	@%p52 bra 	BB0_46;

	cvt.s64.s32	%rd65, %r36;
	add.s64 	%rd66, %rd3, %rd65;
	ld.global.nc.u8 	%rs41, [%rd66];

BB0_46:
	setp.gt.u16	%p53, %rs41, %rs1;
	cvt.u32.u16	%r103, %rs41;
	and.b32  	%r104, %r103, 255;
	selp.b32	%r105, %r7, %r104, %p53;
	selp.b32	%r106, %r104, %r7, %p53;
	add.s32 	%r107, %r106, 1;
	mul.lo.s32 	%r108, %r107, %r106;
	shr.u32 	%r109, %r108, 1;
	add.s32 	%r37, %r109, %r105;
	and.pred  	%p56, %p52, %p13;
	@%p56 bra 	BB0_48;

	mul.wide.s32 	%rd67, %r37, 4;
	add.s64 	%rd68, %rd2, %rd67;
	ld.global.nc.f32 	%f201, [%rd68];
	add.f32 	%f202, %f201, %f201;
	add.s64 	%rd69, %rd1, %rd67;
	ld.global.nc.f32 	%f203, [%rd69];
	div.rn.f32 	%f204, %f203, %f202;
	mul.f32 	%f205, %f204, %f89;
	fma.rn.f32 	%f206, %f2, %f205, %f1;
	mul.f32 	%f207, %f1, %f205;
	sub.f32 	%f208, %f2, %f207;
	selp.f32	%f209, %f206, %f278, %p52;
	selp.f32	%f210, %f208, %f279, %p52;
	selp.f32	%f211, %f3, %f280, %p52;
	mul.f32 	%f212, %f89, %f89;
	div.rn.f32 	%f213, %f202, %f212;
	sub.f32 	%f214, %f209, %f1;
	sub.f32 	%f215, %f210, %f2;
	sub.f32 	%f216, %f211, %f3;
	fma.rn.f32 	%f217, %f214, %f213, %f263;
	fma.rn.f32 	%f218, %f215, %f213, %f264;
	fma.rn.f32 	%f265, %f216, %f213, %f265;
	div.rn.f32 	%f219, %f203, %f89;
	mul.f32 	%f220, %f210, %f219;
	sub.f32 	%f263, %f217, %f220;
	fma.rn.f32 	%f264, %f209, %f219, %f218;

BB0_48:
	add.s32 	%r38, %r3, 1;
	@%p48 bra 	BB0_50;

	rem.s32 	%r110, %r38, %r45;
	add.s32 	%r111, %r110, %r45;
	rem.s32 	%r127, %r111, %r45;
	bra.uni 	BB0_51;

BB0_50:
	add.s32 	%r112, %r45, -1;
	min.s32 	%r127, %r38, %r112;

BB0_51:
	mad.lo.s32 	%r113, %r127, %r44, %r2;
	mad.lo.s32 	%r42, %r113, %r43, %r1;
	setp.ge.s32	%p59, %r38, %r45;
	mov.f32 	%f284, 0f00000000;
	and.pred  	%p61, %p59, %p48;
	mov.f32 	%f285, %f284;
	mov.f32 	%f286, %f284;
	@%p61 bra 	BB0_53;

	mul.wide.s32 	%rd70, %r42, 4;
	add.s64 	%rd71, %rd6, %rd70;
	ld.global.nc.f32 	%f286, [%rd71];
	add.s64 	%rd72, %rd5, %rd70;
	ld.global.nc.f32 	%f285, [%rd72];
	add.s64 	%rd73, %rd4, %rd70;
	ld.global.nc.f32 	%f284, [%rd73];

BB0_53:
	mul.f32 	%f224, %f286, %f286;
	fma.rn.f32 	%f225, %f285, %f285, %f224;
	fma.rn.f32 	%f75, %f284, %f284, %f225;
	setp.eq.f32	%p62, %f75, 0f00000000;
	mov.u16 	%rs42, %rs1;
	@%p62 bra 	BB0_55;

	cvt.s64.s32	%rd74, %r42;
	add.s64 	%rd75, %rd3, %rd74;
	ld.global.nc.u8 	%rs42, [%rd75];

BB0_55:
	setp.gt.u16	%p63, %rs42, %rs1;
	cvt.u32.u16	%r114, %rs42;
	and.b32  	%r115, %r114, 255;
	selp.b32	%r116, %r7, %r115, %p63;
	selp.b32	%r117, %r115, %r7, %p63;
	add.s32 	%r118, %r117, 1;
	mul.lo.s32 	%r119, %r118, %r117;
	shr.u32 	%r120, %r119, 1;
	add.s32 	%r121, %r120, %r116;
	mul.wide.s32 	%rd76, %r121, 4;
	add.s64 	%rd7, %rd2, %rd76;
	add.s64 	%rd8, %rd1, %rd76;
	and.pred  	%p66, %p62, %p13;
	@%p66 bra 	BB0_57;

	ld.global.nc.f32 	%f226, [%rd7];
	add.f32 	%f227, %f226, %f226;
	ld.global.nc.f32 	%f228, [%rd8];
	div.rn.f32 	%f229, %f228, %f227;
	mul.f32 	%f230, %f229, %f89;
	mul.f32 	%f231, %f2, %f230;
	sub.f32 	%f232, %f1, %f231;
	fma.rn.f32 	%f233, %f1, %f230, %f2;
	selp.f32	%f234, %f3, %f284, %p62;
	selp.f32	%f235, %f233, %f285, %p62;
	selp.f32	%f236, %f232, %f286, %p62;
	mul.f32 	%f237, %f89, %f89;
	div.rn.f32 	%f238, %f227, %f237;
	sub.f32 	%f239, %f236, %f1;
	sub.f32 	%f240, %f235, %f2;
	sub.f32 	%f241, %f234, %f3;
	fma.rn.f32 	%f242, %f239, %f238, %f263;
	fma.rn.f32 	%f243, %f240, %f238, %f264;
	fma.rn.f32 	%f265, %f241, %f238, %f265;
	div.rn.f32 	%f244, %f228, %f89;
	fma.rn.f32 	%f263, %f235, %f244, %f242;
	mul.f32 	%f245, %f236, %f244;
	sub.f32 	%f264, %f243, %f245;

BB0_57:
	setp.eq.s64	%p68, %rd12, 0;
	@%p68 bra 	BB0_59;

	cvta.to.global.u64 	%rd77, %rd12;
	add.s64 	%rd79, %rd77, %rd19;
	ld.global.nc.f32 	%f246, [%rd79];
	mul.f32 	%f290, %f246, %f290;

BB0_59:
	setp.eq.f32	%p69, %f290, 0f00000000;
	mov.f32 	%f291, 0f00000000;
	@%p69 bra 	BB0_61;

	rcp.rn.f32 	%f291, %f290;

BB0_61:
	cvta.to.global.u64 	%rd80, %rd11;
	cvta.to.global.u64 	%rd81, %rd10;
	cvta.to.global.u64 	%rd82, %rd9;
	add.s64 	%rd84, %rd82, %rd19;
	ld.global.f32 	%f248, [%rd84];
	fma.rn.f32 	%f249, %f263, %f291, %f248;
	st.global.f32 	[%rd84], %f249;
	add.s64 	%rd85, %rd81, %rd19;
	ld.global.f32 	%f250, [%rd85];
	fma.rn.f32 	%f251, %f264, %f291, %f250;
	st.global.f32 	[%rd85], %f251;
	add.s64 	%rd86, %rd80, %rd19;
	ld.global.f32 	%f252, [%rd86];
	fma.rn.f32 	%f253, %f265, %f291, %f252;
	st.global.f32 	[%rd86], %f253;

BB0_62:
	ret;
}


`
)
