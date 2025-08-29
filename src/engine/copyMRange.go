// engine/ext_copymrange.go
package engine

import "github.com/MathieuMoalic/amumax/src/engine/cuda"

func init() {
	DeclFunc("CopyMRange", CopyMRange,
		"Copy a block of m (GPU only): CopyMRange(dx,dy,dz, sx,sy,sz, W,H,D, wrap)")
}

// Copies a W×H×D block of m from (sx,sy,sz) to (dx,dy,dz).
func CopyMRange(dx, dy, dz, sx, sy, sz, W, H, D int, wrap bool) {
	m := NormMag.Buffer()
	dst0 := [3]int{dx, dy, dz}
	src0 := [3]int{sx, sy, sz}
	box := [3]int{W, H, D}
	for c := range m.NComp() {
		comp := m.Comp(c)
		cuda.CopyMRange(comp, comp, dst0, src0, box, wrap)
	}
}
