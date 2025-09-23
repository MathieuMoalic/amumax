package engine

import (
	"math"

	"github.com/MathieuMoalic/amumax/src/data"
	"github.com/MathieuMoalic/amumax/src/log"
)

func init() {
	DeclFunc("fieldFromWireVectorMask", fieldFromWireVectorMask, "fieldFromWireVectorMask")
	DeclFunc("magModulatedByMask", magModulatedByMask, "magModulatedByMask")
	DeclFunc("cpwVectorMask", cpwVectorMask, "cpwVectorMask")
}

// fieldFromWire computes the magnetic field at a point (x, z) due to a wire carrying a current I.
// The wire is centered at the origin and has a rectangular cross-section of width 2a and height 2b.
// The wire is oriented along the y-axis.
func fieldFromWire(I, x, z, a, b float64) (float64, float64, float64) {
	const eps = 1e-12
	const mu0 = 4 * math.Pi * 1e-7

	ax := a - x
	bx := -a - x
	bz := b - z
	nbz := -b - z

	if ax == 0 {
		ax = eps
	}
	if bx == 0 {
		bx = eps
	}
	if bz == 0 {
		bz = eps
	}
	if nbz == 0 {
		nbz = eps
	}

	// Bx component
	t1 := (a - x) * (0.5*math.Log((bz*bz+ax*ax)/(nbz*nbz+ax*ax)) +
		(bz/ax)*math.Atan(ax/bz) -
		(nbz/ax)*math.Atan(ax/nbz))

	t2 := -(-a - x) * (0.5*math.Log((bz*bz+bx*bx)/(bx*bx+nbz*nbz)) +
		(bz/bx)*math.Atan(bx/bz) -
		(nbz/bx)*math.Atan(bx/nbz))

	Bx := (I * mu0 / (8 * math.Pi * a * b)) * (t1 + t2)

	// Bz component
	t3 := bz * (0.5*math.Log((bz*bz+ax*ax)/(bz*bz+bx*bx)) +
		(ax/bz)*math.Atan(bz/ax) -
		(bx/bz)*math.Atan(bz/bx))

	t4 := -nbz * (0.5*math.Log((nbz*nbz+ax*ax)/(nbz*nbz+bx*bx)) +
		(ax/nbz)*math.Atan(nbz/ax) -
		(bx/nbz)*math.Atan(nbz/bx))

	Bz := -(I * mu0 / (8 * math.Pi * a * b)) * (t3 + t4)

	return Bx, 0, Bz
}

// fieldFromWireVectorMask creates a vector mask of the magnetic field created by a wire carrying a current I.
// The wire is centered on xcenter, zcenter and has a rectangular cross-section of the given width and height.
// The mask is normalized to the maximum value of the magnetic field.
func fieldFromWireVectorMask(I, width, height, xcenter, zcenter float64) *data.Slice {
	Nx := GetMesh().Nx
	Nz := GetMesh().Nz
	maskSlice := newVectorMask(Nx, 1, Nz)
	for ix := range Nx {
		for iz := range Nz {
			r := index2Coord(ix, 0, iz)
			x := r[0] - xcenter
			z := r[2] - zcenter
			Bx, By, Bz := fieldFromWire(I, x, z, width/2, height/2)
			maskSlice.Set(0, ix, 0, iz, Bx)
			maskSlice.Set(1, ix, 0, iz, By)
			maskSlice.Set(2, ix, 0, iz, Bz)
		}
	}
	return maskSlice
}

// cpwVectorMask creates a vector mask of the magnetic field created by a coplanar waveguide (CPW).
// The CPW is centered on xoffset, zoffset and has a rectangular cross-section of the given width and height.
// The mask is normalized to the maximum value of the magnetic field.
func cpwVectorMask(I, width, height, distance, xoffset, zoffset float64) *data.Slice {
	Nx, _, Nz := Mesh.GetNi()
	maskSlice := newVectorMask(Nx, 1, Nz)
	for ix := range Nx {
		for iz := range Nz {
			r := index2Coord(ix, 0, iz)

			bx1, by1, bz1 := fieldFromWire(-I, r.X()-distance-xoffset, r.Z()-zoffset, width/2, height/2)
			bx2, by2, bz2 := fieldFromWire(I, r.X()-xoffset, r.Z()-zoffset, width/2, height/2)
			bx3, by3, bz3 := fieldFromWire(-I, r.X()+distance-xoffset, r.Z()-zoffset, width/2, height/2)

			Bx := bx1 + bx2 + bx3
			By := by1 + by2 + by3
			Bz := bz1 + bz2 + bz3

			maskSlice.Set(0, ix, 0, iz, Bx)
			maskSlice.Set(1, ix, 0, iz, By)
			maskSlice.Set(2, ix, 0, iz, Bz)
		}
	}
	return maskSlice
}

// magModulatedByMask computes ∑_{c,z,y,x} mask[c][z][y][x] * mag[c][z][y][x].
func magModulatedByMask(maskSlice *data.Slice) (sum float64) {
	magSlice, _ := NormMag.Slice()

	mag := magSlice.HostCopy().Tensors() // [c][z][y][x]float32
	mask := maskSlice.Tensors()          // [c][z][y][x]float32

	// Expect 3 components, but be defensive.
	if len(mag) != 3 || len(mask) != 3 {
		log.Log.Err("magModulatedByMask: expected 3 components, got %d and %d", len(mag), len(mask))
		return 0
	}

	for iz := range mag[0] { // z
		m0z, m1z, m2z := mag[0][iz], mag[1][iz], mag[2][iz]
		k0z, k1z, k2z := mask[0][iz], mask[1][iz], mask[2][iz]

		for iy := range m0z { // y
			m0x, m1x, m2x := m0z[iy], m1z[iy], m2z[iy]
			k0x, k1x, k2x := k0z[iy], k1z[iy], k2z[iy]

			// (Optional) hints for bounds-check elimination
			_ = m0x[len(m0x)-1]
			_ = m1x[len(m1x)-1]
			_ = m2x[len(m2x)-1]
			_ = k0x[len(k0x)-1]
			_ = k1x[len(k1x)-1]
			_ = k2x[len(k2x)-1]

			for ix := range m0x { // x
				// Unrolled channel sum
				sum += float64(
					m0x[ix]*k0x[ix] +
						m1x[ix]*k1x[ix] +
						m2x[ix]*k2x[ix],
				)
			}
		}
	}
	return sum
}
