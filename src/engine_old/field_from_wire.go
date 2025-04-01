package engine_old

import (
	"math"

	"github.com/MathieuMoalic/amumax/src/engine_old/data_old"
)

func init() {
	DeclFunc("fieldFromWireVectorMask", fieldFromWireVectorMask, "fieldFromWireVectorMask")
	DeclFunc("magModulatedByMask", magModulatedByMask, "magModulatedByMask")
	DeclFunc("cpwVectorMask", cpwVectorMask, "cpwVectorMask")
	DeclFunc("fieldFromWire", fieldFromWire, "fieldFromWire")
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
func fieldFromWireVectorMask(I, width, height, xcenter, zcenter float64) *data_old.Slice {
	Nx := GetMesh().Nx
	Nz := GetMesh().Nz
	maskSlice := newVectorMask(Nx, 1, Nz)
	B_max := 0.0
	for ix := range Nx {
		for iz := range Nz {
			r := index2Coord(ix, 0, iz)
			x := r[0] - xcenter
			z := r[2] - zcenter
			Bx, By, Bz := fieldFromWire(I, x, z, width/2, height/2)
			B_max = math.Max(B_max, math.Abs(Bx))
			B_max = math.Max(B_max, math.Abs(Bz))
			maskSlice.Set(0, ix, 0, iz, Bx)
			maskSlice.Set(1, ix, 0, iz, By)
			maskSlice.Set(2, ix, 0, iz, Bz)
		}
	}
	for ix := range Nx {
		for iz := range Nz {
			maskSlice.Set(0, ix, 0, iz, maskSlice.Get(0, ix, 0, iz)/B_max)
			maskSlice.Set(2, ix, 0, iz, maskSlice.Get(2, ix, 0, iz)/B_max)
		}
	}
	return maskSlice
}

// cpwVectorMask creates a vector mask of the magnetic field created by a coplanar waveguide (CPW).
// The CPW is centered on xoffset, zoffset and has a rectangular cross-section of the given width and height.
// The mask is normalized to the maximum value of the magnetic field.
func cpwVectorMask(I, width, height, distance, xoffset, zoffset float64) *data_old.Slice {
	Nx, _, Nz := Mesh.GetNi()
	maskSlice := newVectorMask(Nx, 1, Nz)
	B_max := 0.0
	for ix := range Nx {
		for iz := range Nz {
			r := index2Coord(ix, 0, iz)

			bx1, by1, bz1 := fieldFromWire(-I, r.X()-distance-xoffset, r.Z()-zoffset, width/2, height/2)
			bx2, by2, bz2 := fieldFromWire(I, r.X()-xoffset, r.Z()-zoffset, width/2, height/2)
			bx3, by3, bz3 := fieldFromWire(-I, r.X()+distance-xoffset, r.Z()-zoffset, width/2, height/2)

			Bx := bx1 + bx2 + bx3
			By := by1 + by2 + by3
			Bz := bz1 + bz2 + bz3

			B_max = math.Max(B_max, math.Abs(Bx))
			B_max = math.Max(B_max, math.Abs(Bz))
			maskSlice.Set(0, ix, 0, iz, Bx)
			maskSlice.Set(1, ix, 0, iz, By)
			maskSlice.Set(2, ix, 0, iz, Bz)
		}
	}
	if B_max > 0 {
		for ix := range Nx {
			for iz := range Nz {
				maskSlice.Set(0, ix, 0, iz, maskSlice.Get(0, ix, 0, iz)/B_max)
				maskSlice.Set(2, ix, 0, iz, maskSlice.Get(2, ix, 0, iz)/B_max)
			}
		}
	}
	return maskSlice
}

// magModulatedByMask computes the signal of the magnetic field modulated by a mask.
func magModulatedByMask(mask_slice *data_old.Slice) float64 {
	Nx, Ny, Nz := Mesh.GetNi()
	var signal float32
	mag_slice, _ := NormMag.Slice()
	mag := mag_slice.HostCopy().Tensors() // [c][z][y][x]
	mask := mask_slice.Tensors()          // [c][z][y][x]

	for iz := range Nz {
		for iy := range Ny {
			for ix := range Nx {
				for ic := range 3 {
					signal += mask[ic][iz][0][ix] * mag[ic][iz][iy][ix]
				}
			}
		}
	}
	return float64(signal)
}
