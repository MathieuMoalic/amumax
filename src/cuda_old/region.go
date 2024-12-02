package cuda_old

import (
	"unsafe"

	"github.com/MathieuMoalic/amumax/src/engine_old/data_old"
	"github.com/MathieuMoalic/amumax/src/engine_old/log_old"
)

// dst += LUT[region], for vectors. Used to add terms to excitation.
func RegionAddV(dst *data_old.Slice, lut LUTPtrs, regions *Bytes) {
	log_old.AssertMsg(dst.NComp() == 3, "Component mismatch: dst must have 3 components in RegionAddV")
	N := dst.Len()
	cfg := make1DConf(N)
	k_regionaddv_async(dst.DevPtr(X), dst.DevPtr(Y), dst.DevPtr(Z),
		lut[X], lut[Y], lut[Z], regions.Ptr, N, cfg)
}

// dst += LUT[region], for scalar. Used to add terms to scalar excitation.
func RegionAddS(dst *data_old.Slice, lut LUTPtr, regions *Bytes) {
	log_old.AssertMsg(dst.NComp() == 1, "Component mismatch: dst must have 1 component in RegionAddS")
	N := dst.Len()
	cfg := make1DConf(N)
	k_regionadds_async(dst.DevPtr(0), unsafe.Pointer(lut), regions.Ptr, N, cfg)
}

// decode the regions+LUT pair into an uncompressed array
func RegionDecode(dst *data_old.Slice, lut LUTPtr, regions *Bytes) {
	N := dst.Len()
	cfg := make1DConf(N)
	k_regiondecode_async(dst.DevPtr(0), unsafe.Pointer(lut), regions.Ptr, N, cfg)
}

// select the part of src within the specified region, set 0's everywhere else.
func RegionSelect(dst, src *data_old.Slice, regions *Bytes, region byte) {
	log_old.AssertMsg(dst.NComp() == src.NComp(), "Component mismatch: dst and src must have the same number of components in RegionSelect")
	N := dst.Len()
	cfg := make1DConf(N)

	for c := 0; c < dst.NComp(); c++ {
		k_regionselect_async(dst.DevPtr(c), src.DevPtr(c), regions.Ptr, region, N, cfg)
	}
}
