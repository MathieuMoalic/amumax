package mag

import (
	"fmt"
	"math"

	"github.com/MathieuMoalic/amumax/src/data"
	"github.com/MathieuMoalic/amumax/src/log_old"
	"github.com/MathieuMoalic/amumax/src/mesh"
	"github.com/MathieuMoalic/amumax/src/mesh_old"
)

func MFMKernel(mesh mesh.MeshLike, lift, tipsize float64, cacheDir string) (kernel [3]*data.Slice) {
	return CalcMFMKernel(mesh, lift, tipsize)

}

// Kernel for the vertical derivative of the force on an MFM tip due to mx, my, mz.
// This is the 2nd derivative of the energy w.r.t. z.
func CalcMFMKernel(kernelMesh mesh.MeshLike, lift, tipsize float64) (kernel [3]*data.Slice) {

	const TipCharge = 1 / Mu0 // tip charge
	const Δ = 1e-9            // tip oscillation, take 2nd derivative over this distance
	log_old.AssertMsg(lift > 0, "MFM tip crashed into sample, please lift the new one higher")

	{ // Kernel mesh is 2x larger than input, instead in case of PBC
		pbc := kernelMesh.PBC()
		sz := padSize(kernelMesh.Size(), pbc)
		cs := kernelMesh.CellSize()
		kernelMesh = mesh_old.NewMesh(sz[X], sz[Y], sz[Z], cs[X], cs[Y], cs[Z], pbc[0], pbc[1], pbc[2])
	}

	// Shorthand
	size := kernelMesh.Size()
	pbc := kernelMesh.PBC()
	cellsize := kernelMesh.CellSize()
	volume := cellsize[X] * cellsize[Y] * cellsize[Z]
	fmt.Println("calculating MFM kernel")

	// Sanity check
	{
		log_old.AssertMsg(size[Z] >= 1 && size[Y] >= 2 && size[X] >= 2, "Invalid MFM kernel size: size must be at least 1x2x2")
		log_old.AssertMsg(cellsize[X] > 0 && cellsize[Y] > 0 && cellsize[Z] > 0, "Invalid cell size: all dimensions must be positive in CalcMFMKernel")
		log_old.AssertMsg(size[X]%2 == 0 && size[Y]%2 == 0, "MFM must have even cellsize on the X and Y axis")
		if size[Z] > 1 {
			log_old.AssertMsg(size[Z]%2 == 0, "MFM only supports one cell thickness on the Z axis")
		}
	}

	// Allocate only upper diagonal part. The rest is symmetric due to reciprocity.
	var K [3][][][]float32
	for i := 0; i < 3; i++ {
		kernel[i] = data.NewSlice(1, kernelMesh.Size())
		K[i] = kernel[i].Scalars()
	}

	r1, r2 := kernelRanges(size, pbc)

	for iz := r1[Z]; iz <= r2[Z]; iz++ {
		zw := wrap(iz, size[Z])
		z := float64(iz) * cellsize[Z]

		for iy := r1[Y]; iy <= r2[Y]; iy++ {
			yw := wrap(iy, size[Y])
			y := float64(iy) * cellsize[Y]

			for ix := r1[X]; ix <= r2[X]; ix++ {
				x := float64(ix) * cellsize[X]
				xw := wrap(ix, size[X])

				for s := 0; s < 3; s++ { // source index Ksxyz
					m := data.Vector{0, 0, 0}
					m[s] = 1

					var E [3]float64 // 3 energies for 2nd derivative

					for i := -1; i <= 1; i++ {
						I := float64(i)
						R := data.Vector{-x, -y, z - (lift + (I * Δ))}
						r := R.Len()
						B := R.Mul(TipCharge / (4 * math.Pi * r * r * r))

						R = data.Vector{-x, -y, z - (lift + tipsize + (I * Δ))}
						r = R.Len()
						B = B.Add(R.Mul(-TipCharge / (4 * math.Pi * r * r * r)))

						E[i+1] = B.Dot(m) * volume // i=-1 stored in  E[0]
					}

					dFdz_tip := ((E[0] - E[1]) + (E[2] - E[1])) / (Δ * Δ) // dFz/dz = d2E/dz2

					K[s][zw][yw][xw] += float32(dFdz_tip) // += needed in case of PBC
				}
			}
		}
	}

	return kernel
}
