package cuda

import (
	"unsafe"

	"github.com/MathieuMoalic/amumax/src/log"
	"github.com/MathieuMoalic/amumax/src/slice"
)

// dst += LUT[region], for vectors. Used to add terms to excitation.
func RegionAddV(dst *slice.Slice, lut LUTPtrs, regions *Bytes) {
	log.AssertMsg(dst.NComp() == 3, "Component mismatch: dst must have 3 components in RegionAddV")
	N := dst.Len()
	cfg := make1DConf(N)
	k_regionaddv_async(dst.DevPtr(X), dst.DevPtr(Y), dst.DevPtr(Z),
		lut[X], lut[Y], lut[Z], regions.Ptr, N, cfg)
}

// dst += LUT[region], for scalar. Used to add terms to scalar excitation.
func RegionAddS(dst *slice.Slice, lut LUTPtr, regions *Bytes) {
	log.AssertMsg(dst.NComp() == 1, "Component mismatch: dst must have 1 component in RegionAddS")
	N := dst.Len()
	cfg := make1DConf(N)
	k_regionadds_async(dst.DevPtr(0), unsafe.Pointer(lut), regions.Ptr, N, cfg)
}

// decode the regions+LUT pair into an uncompressed array
func RegionDecode(dst *slice.Slice, lut LUTPtr, regions *Bytes) {
	N := dst.Len()
	cfg := make1DConf(N)
	k_regiondecode_async(dst.DevPtr(0), unsafe.Pointer(lut), regions.Ptr, N, cfg)
}

// select the part of src within the specified region, set 0's everywhere else.
func RegionSelect(dst, src *slice.Slice, regions *Bytes, region byte) {
	log.AssertMsg(dst.NComp() == src.NComp(), "Component mismatch: dst and src must have the same number of components in RegionSelect")
	N := dst.Len()
	cfg := make1DConf(N)

	for c := 0; c < dst.NComp(); c++ {
		k_regionselect_async(dst.DevPtr(c), src.DevPtr(c), regions.Ptr, region, N, cfg)
	}
}
