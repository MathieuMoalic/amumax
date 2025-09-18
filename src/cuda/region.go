package cuda

import (
	"unsafe"

	"github.com/MathieuMoalic/amumax/src/data"
	"github.com/MathieuMoalic/amumax/src/log"
)

// RegionAddV dst += LUT[region], for vectors. Used to add terms to excitation.
func RegionAddV(dst *data.Slice, lut LUTPtrs, regions *Bytes) {
	log.AssertMsg(dst.NComp() == 3, "Component mismatch: dst must have 3 components in RegionAddV")
	N := dst.Len()
	cfg := make1DConf(N)
	kRegionaddvAsync(dst.DevPtr(X), dst.DevPtr(Y), dst.DevPtr(Z),
		lut[X], lut[Y], lut[Z], regions.Ptr, N, cfg)
}

// RegionAddS dst += LUT[region], for scalar. Used to add terms to scalar excitation.
func RegionAddS(dst *data.Slice, lut LUTPtr, regions *Bytes) {
	log.AssertMsg(dst.NComp() == 1, "Component mismatch: dst must have 1 component in RegionAddS")
	N := dst.Len()
	cfg := make1DConf(N)
	kRegionaddsAsync(dst.DevPtr(0), unsafe.Pointer(lut), regions.Ptr, N, cfg)
}

// RegionDecode decode the regions+LUT pair into an uncompressed array
func RegionDecode(dst *data.Slice, lut LUTPtr, regions *Bytes) {
	N := dst.Len()
	cfg := make1DConf(N)
	kRegiondecodeAsync(dst.DevPtr(0), unsafe.Pointer(lut), regions.Ptr, N, cfg)
}

// RegionSelect select the part of src within the specified region, set 0's everywhere else.
func RegionSelect(dst, src *data.Slice, regions *Bytes, region byte) {
	log.AssertMsg(dst.NComp() == src.NComp(), "Component mismatch: dst and src must have the same number of components in RegionSelect")
	N := dst.Len()
	cfg := make1DConf(N)

	for c := 0; c < dst.NComp(); c++ {
		kRegionselectAsync(dst.DevPtr(c), src.DevPtr(c), regions.Ptr, region, N, cfg)
	}
}
