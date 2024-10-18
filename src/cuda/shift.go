package cuda

import (
	"github.com/MathieuMoalic/amumax/src/data"
	"github.com/MathieuMoalic/amumax/src/log"
)

// shift dst by shx cells (positive or negative) along X-axis.
// new edge value is clampL at left edge or clampR at right edge.
func ShiftX(dst, src *data.Slice, shiftX int, clampL, clampR float32) {
	log.AssertMsg(dst.NComp() == 1 && src.NComp() == 1, "Component mismatch: dst and src must both have 1 component in ShiftX")
	log.AssertMsg(dst.Len() == src.Len(), "Length mismatch: dst and src must have the same length in ShiftX")
	N := dst.Size()
	cfg := make3DConf(N)
	k_shiftx_async(dst.DevPtr(0), src.DevPtr(0), N[X], N[Y], N[Z], shiftX, clampL, clampR, cfg)
}

// Shifts a component `src` of a vector field by `shiftX` cells along the X-axis.
// Unlike the normal `shift()`, the new edge value is the current edge value.
//
// To avoid the situation where the magnetization could be set to (0,0,0) within the geometry, it is
// also required to pass the two other vector components `othercomp` and `anothercomp` to this function.
// In cells where the vector (`src`, `othercomp`, `anothercomp`) is the zero-vector,
// `clampL` or `clampR` is used for the component `src` instead.
func ShiftEdgeCarryX(dst, src, othercomp, anothercomp *data.Slice, shiftX int, clampL, clampR float32) {
	log.AssertMsg(dst.NComp() == 1 && src.NComp() == 1 && othercomp.NComp() == 1 && anothercomp.NComp() == 1, "Component mismatch: dst, src, othercomp and anothercomp must all have 1 component in ShiftEdgeCarryX")
	log.AssertMsg(dst.Len() == src.Len(), "Length mismatch: dst and src must have the same length in ShiftEdgeCarryX")
	N := dst.Size()
	cfg := make3DConf(N)
	k_shiftedgecarryX_async(dst.DevPtr(0), src.DevPtr(0), othercomp.DevPtr(0), anothercomp.DevPtr(0), N[X], N[Y], N[Z], shiftX, clampL, clampR, cfg)
}

func ShiftY(dst, src *data.Slice, shiftY int, clampL, clampR float32) {
	log.AssertMsg(dst.NComp() == 1 && src.NComp() == 1, "Component mismatch: dst and src must both have 1 component in ShiftY")
	log.AssertMsg(dst.Len() == src.Len(), "Length mismatch: dst and src must have the same length in ShiftY")
	N := dst.Size()
	cfg := make3DConf(N)
	k_shifty_async(dst.DevPtr(0), src.DevPtr(0), N[X], N[Y], N[Z], shiftY, clampL, clampR, cfg)
}

// Shifts a component `src` of a vector field by `shiftY` cells along the Y-axis.
// Unlike the normal `shift()`, the new edge value is the current edge value.
//
// To avoid the situation where the magnetization could be set to (0,0,0) within the geometry, it is
// also required to pass the two other vector components `othercomp` and `anothercomp` to this function.
// In cells where the vector (`src`, `othercomp`, `anothercomp`) is the zero-vector,
// `clampD` or `clampU` is used for the component `src` instead.
func ShiftEdgeCarry(dst, src, othercomp, anothercomp *data.Slice, shiftY int, clampL, clampR float32) {
	log.AssertMsg(dst.NComp() == 1 && src.NComp() == 1 && othercomp.NComp() == 1 && anothercomp.NComp() == 1, "Component mismatch: dst, src, othercomp and anothercomp must all have 1 component in ShiftEdgeCarry")
	log.AssertMsg(dst.Len() == src.Len(), "Length mismatch: dst and src must have the same length in ShiftEdgeCarry")
	N := dst.Size()
	cfg := make3DConf(N)
	k_shiftedgecarryY_async(dst.DevPtr(0), src.DevPtr(0), othercomp.DevPtr(0), anothercomp.DevPtr(0), N[X], N[Y], N[Z], shiftY, clampL, clampR, cfg)
}

func ShiftZ(dst, src *data.Slice, shiftZ int, clampL, clampR float32) {
	log.AssertMsg(dst.NComp() == 1 && src.NComp() == 1, "Component mismatch: dst and src must both have 1 component in ShiftZ")
	log.AssertMsg(dst.Len() == src.Len(), "Length mismatch: dst and src must have the same length in ShiftZ")
	N := dst.Size()
	cfg := make3DConf(N)
	k_shiftz_async(dst.DevPtr(0), src.DevPtr(0), N[X], N[Y], N[Z], shiftZ, clampL, clampR, cfg)
}

// Like Shift, but for bytes
func ShiftBytes(dst, src *Bytes, m *data.Mesh, shiftX int, clamp byte) {
	N := m.Size()
	cfg := make3DConf(N)
	k_shiftbytes_async(dst.Ptr, src.Ptr, N[X], N[Y], N[Z], shiftX, clamp, cfg)
}

func ShiftBytesY(dst, src *Bytes, m *data.Mesh, shiftY int, clamp byte) {
	N := m.Size()
	cfg := make3DConf(N)
	k_shiftbytesy_async(dst.Ptr, src.Ptr, N[X], N[Y], N[Z], shiftY, clamp, cfg)
}
