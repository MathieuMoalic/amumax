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

func ShiftY(dst, src *data.Slice, shiftY int, clampL, clampR float32) {
	log.AssertMsg(dst.NComp() == 1 && src.NComp() == 1, "Component mismatch: dst and src must both have 1 component in ShiftY")
	log.AssertMsg(dst.Len() == src.Len(), "Length mismatch: dst and src must have the same length in ShiftY")
	N := dst.Size()
	cfg := make3DConf(N)
	k_shifty_async(dst.DevPtr(0), src.DevPtr(0), N[X], N[Y], N[Z], shiftY, clampL, clampR, cfg)
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
