package mag

import (
	"bufio"
	"fmt"
	"math"
	"os"

	"github.com/DataDog/zstd"
	"github.com/MathieuMoalic/amumax/data"
	"github.com/MathieuMoalic/amumax/httpfs"
	"github.com/MathieuMoalic/amumax/timer"
	"github.com/MathieuMoalic/amumax/util"
	"github.com/MathieuMoalic/amumax/zarr"
)

// Obtains the demag kernel either from cacheDir/ or by calculating (and then storing in cacheDir for next time).
// Empty cacheDir disables caching.
func DemagKernel(gridsize, pbc [3]int, cellsize [3]float64, accuracy float64, cacheDir string, showMagnets bool) (kernel [3][3]*data.Slice) {
	timer.Start("kernel_init")
	timer.Stop("kernel_init") // warm-up

	timer.Start("kernel_init")
	defer timer.Stop("kernel_init")

	sanityCheck(cellsize)
	// Cache disabled
	if cacheDir == "" {
		util.Log.Warn(`Kernel cache disabled.`)
		kernel = calcDemagKernel(gridsize, pbc, cellsize, accuracy, showMagnets)
		// return without saving
		return
	}
	// Make sure cache directory exists
	if !httpfs.Exists(cacheDir) {
		err := httpfs.Mkdir(cacheDir)
		if err != nil {
			util.Log.Warn("Unable to create kernel cache directory: %v", err)
		}
	}
	defer func() {
		if err := recover(); err != nil {
			util.Log.Warn("Unable to use kernel cache: %v", err)
			kernel = calcDemagKernel(gridsize, pbc, cellsize, accuracy, showMagnets)
		}
	}()

	basename := kernelName(gridsize, pbc, cellsize, accuracy, cacheDir)
	if httpfs.Exists(basename) {
		util.Log.Info("Loading kernel from cache")
		var err error
		kernel, err = loadKernel(basename, padSize(gridsize, pbc))
		if err != nil {
			util.Log.Warn("Couldn't load kernel from cache: %v", err)
			kernel = calcDemagKernel(gridsize, pbc, cellsize, accuracy, showMagnets)
		}
		return kernel
	} else {
		util.Log.Info("Calculating kernel and saving to cache")
		kernel = calcDemagKernel(gridsize, pbc, cellsize, accuracy, showMagnets)
		err := saveKernel(basename, kernel)
		if err != nil {
			util.Log.Warn("Couldn't save kernel to cache: %v \n %v", basename, err.Error())
		}
		return
	}
}

func bytesToKernel(kernelBytes []byte, size [3]int) (kernel [3][3]*data.Slice) {
	offset := 0
	sliceLength := size[X] * size[Y] * size[Z] * 4
	for i := 0; i < 3; i++ {
		for j := i; j < 3; j++ {
			// if Nz is 1, we have a 2D simulation and these elements are zero
			if size[Z] == 1 {
				if (i == 0 && j == 2) || (i == 1 && j == 2) {
					kernel[i][j] = nil
					continue
				}
			}
			end := offset + sliceLength
			kernel[i][j] = bytesToSlice(kernelBytes[offset:end], size)
			offset = end
		}
	}
	return
}

func kernelToBytes(kernel [3][3]*data.Slice) (bytes []byte) {
	for i := 0; i < 3; i++ {
		for j := i; j < 3; j++ {
			kernelBytes := sliceToBytes(kernel[i][j])
			for k := 0; k < len(kernelBytes); k++ {
				bytes = append(bytes, kernelBytes[k])
			}
		}
	}
	return bytes
}

func sliceToBytes(slice *data.Slice) (bytes []byte) {
	size := slice.Size()
	if size[X] == 0 && size[Y] == 0 && size[Z] == 0 {
		return
	}
	data := slice.Tensors()
	for iz := 0; iz < size[Z]; iz++ {
		for iy := 0; iy < size[Y]; iy++ {
			for ix := 0; ix < size[X]; ix++ {
				bytes = append(bytes, zarr.Float32ToBytes(data[0][iz][iy][ix])...)
			}
		}
	}
	return bytes
}

func bytesToSlice(kernelBytes []byte, size [3]int) (slice *data.Slice) {
	slice = data.NewSlice(1, size)
	tensors := slice.Tensors()
	count := 0
	for iz := 0; iz < size[Z]; iz++ {
		for iy := 0; iy < size[Y]; iy++ {
			for ix := 0; ix < size[X]; ix++ {
				tensors[0][iz][iy][ix] = zarr.BytesToFloat32(kernelBytes[count*4 : (count+1)*4])
				count++
			}
		}
	}
	return
}

func kernelName(gridsize, pbc [3]int, cellsize [3]float64, accuracy float64, cacheDir string) string {
	sSize := fmt.Sprintf("%d_%d_%d", gridsize[X], gridsize[Y], gridsize[Z])
	sPBC := fmt.Sprintf("%d_%d_%d", pbc[X], pbc[Y], pbc[Z])
	sCellsize := fmt.Sprintf("%e_%e_%e", cellsize[X], cellsize[Y], cellsize[Z])
	return fmt.Sprintf("%s/%s_%s_%s_%v.cache", cacheDir, sSize, sPBC, sCellsize, accuracy)
}

func loadKernel(fname string, size [3]int) ([3][3]*data.Slice, error) {
	compressedData, err := httpfs.Read(fname)
	if err != nil {
		return [3][3]*data.Slice{}, err
	}
	kernelBytes, err := zstd.Decompress(nil, compressedData)
	if err != nil {
		return [3][3]*data.Slice{}, err
	}
	return bytesToKernel(kernelBytes, size), nil
}

func saveKernel(fname string, kernel [3][3]*data.Slice) error {
	kernelBytes := kernelToBytes(kernel)
	compressedData, err := zstd.Compress(nil, kernelBytes)
	if err != nil {
		return err
	}
	f, err := os.OpenFile(fname, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	out := bufio.NewWriter(f)
	defer out.Flush()
	_, err = out.Write(compressedData)
	if err != nil {
		return err
	}
	return nil
}

// Calculates the magnetostatic kernel by brute-force integration
// of magnetic charges over the faces and averages over cell volumes.
func calcDemagKernel(gridsize, pbc [3]int, cellsize [3]float64, accuracy float64, showMagnets bool) (kernel [3][3]*data.Slice) {
	// Add zero-padding in non-PBC directions
	size := padSize(gridsize, pbc)

	// Sanity check
	{
		util.Assert(size[Z] > 0 && size[Y] > 0 && size[X] > 0)
		util.Assert(cellsize[X] > 0 && cellsize[Y] > 0 && cellsize[Z] > 0)
		util.Assert(pbc[X] >= 0 && pbc[Y] >= 0 && pbc[Z] >= 0)
		util.Assert(accuracy > 0)
	}

	// Allocate only upper diagonal part. The rest is symmetric due to reciprocity.
	var array [3][3][][][]float32
	for i := 0; i < 3; i++ {
		for j := i; j < 3; j++ {
			kernel[i][j] = data.NewSlice(1, size)
			array[i][j] = kernel[i][j].Scalars()
		}
	}
	// Field (destination) loop ranges
	r1, r2 := kernelRanges(size, pbc)

	// smallest cell dimension is our typical length scale
	L := cellsize[X]
	{
		if cellsize[Y] < L {
			L = cellsize[Y]
		}
		if cellsize[Z] < L {
			L = cellsize[Z]
		}
	}

	progress, progmax := 0, (1+(r2[Y]-r1[Y]))*(1+(r2[Z]-r1[Z])) // progress bar
	done := make(chan struct{}, 3)                              // parallel calculation of one component done?

	ProgressBar := zarr.ProgressBar{}
	ProgressBar.New(0, float64(progmax), "🔧", showMagnets)
	// Start brute integration
	// 9 nested loops, does that stress you out?
	// Fortunately, the 5 inner ones usually loop over just one element.

	for s := 0; s < 3; s++ { // source index Ksdxyz (parallelized over)
		go func(s int) {
			u, v, w := s, (s+1)%3, (s+2)%3 // u = direction of source (s), v & w are the orthogonal directions
			var (
				R, R2  [3]float64 // field and source cell center positions
				pole   [3]float64 // position of point charge on the surface
				points int        // counts used integration points
			)

			for z := r1[Z]; z <= r2[Z]; z++ {
				zw := wrap(z, size[Z])
				// skip one half, reconstruct from symmetry later
				// check on wrapped index instead of loop range so it also works for PBC
				if zw > size[Z]/2 {
					if s == 0 {
						progress += (1 + (r2[Y] - r1[Y]))
					}
					continue
				}
				R[Z] = float64(z) * cellsize[Z]

				for y := r1[Y]; y <= r2[Y]; y++ {

					if s == 0 { // show progress of only one component
						progress++
						ProgressBar.Update(float64(progress))
					}

					yw := wrap(y, size[Y])
					if yw > size[Y]/2 {
						continue
					}
					R[Y] = float64(y) * cellsize[Y]

					for x := r1[X]; x <= r2[X]; x++ {
						xw := wrap(x, size[X])
						if xw > size[X]/2 {
							continue
						}
						R[X] = float64(x) * cellsize[X]

						// choose number of integration points depending on how far we are from source.
						dx, dy, dz := delta(x)*cellsize[X], delta(y)*cellsize[Y], delta(z)*cellsize[Z]
						d := math.Sqrt(dx*dx + dy*dy + dz*dz)
						if d == 0 {
							d = L
						}
						maxSize := d / accuracy // maximum acceptable integration size

						nv := int(math.Max(cellsize[v]/maxSize, 1) + 0.5)
						nw := int(math.Max(cellsize[w]/maxSize, 1) + 0.5)
						nx := int(math.Max(cellsize[X]/maxSize, 1) + 0.5)
						ny := int(math.Max(cellsize[Y]/maxSize, 1) + 0.5)
						nz := int(math.Max(cellsize[Z]/maxSize, 1) + 0.5)
						// Stagger source and destination grids.
						// Massively improves accuracy, see note.
						nv *= 2
						nw *= 2

						util.Assert(nv > 0 && nw > 0 && nx > 0 && ny > 0 && nz > 0)

						scale := 1 / float64(nv*nw*nx*ny*nz)
						surface := cellsize[v] * cellsize[w] // the two directions perpendicular to direction s
						charge := surface * scale
						pu1 := cellsize[u] / 2. // positive pole center
						pu2 := -pu1             // negative pole center

						// Do surface integral over source cell, accumulate  in B
						var B [3]float64
						for i := 0; i < nv; i++ {
							pv := -(cellsize[v] / 2.) + cellsize[v]/float64(2*nv) + float64(i)*(cellsize[v]/float64(nv))
							pole[v] = pv
							for j := 0; j < nw; j++ {
								pw := -(cellsize[w] / 2.) + cellsize[w]/float64(2*nw) + float64(j)*(cellsize[w]/float64(nw))
								pole[w] = pw

								// Do volume integral over destination cell
								for α := 0; α < nx; α++ {
									rx := R[X] - cellsize[X]/2 + cellsize[X]/float64(2*nx) + (cellsize[X]/float64(nx))*float64(α)

									for β := 0; β < ny; β++ {
										ry := R[Y] - cellsize[Y]/2 + cellsize[Y]/float64(2*ny) + (cellsize[Y]/float64(ny))*float64(β)

										for γ := 0; γ < nz; γ++ {
											rz := R[Z] - cellsize[Z]/2 + cellsize[Z]/float64(2*nz) + (cellsize[Z]/float64(nz))*float64(γ)
											points++

											pole[u] = pu1
											R2[X], R2[Y], R2[Z] = rx-pole[X], ry-pole[Y], rz-pole[Z]
											r := math.Sqrt(R2[X]*R2[X] + R2[Y]*R2[Y] + R2[Z]*R2[Z])
											qr := charge / (4 * math.Pi * r * r * r)
											bx := R2[X] * qr
											by := R2[Y] * qr
											bz := R2[Z] * qr

											pole[u] = pu2
											R2[X], R2[Y], R2[Z] = rx-pole[X], ry-pole[Y], rz-pole[Z]
											r = math.Sqrt(R2[X]*R2[X] + R2[Y]*R2[Y] + R2[Z]*R2[Z])
											qr = -charge / (4 * math.Pi * r * r * r)
											B[X] += (bx + R2[X]*qr) // addition ordered for accuracy
											B[Y] += (by + R2[Y]*qr)
											B[Z] += (bz + R2[Z]*qr)

										}
									}
								}
							}
						}
						for d := s; d < 3; d++ { // destination index Ksdxyz
							array[s][d][zw][yw][xw] += float32(B[d]) // += needed in case of PBC
						}
					}
				}
			}
			done <- struct{}{} // notify parallel computation of this component is done
		}(s)
	}
	// wait for all 3 components to finish
	<-done
	<-done
	<-done
	// Reconstruct skipped parts from symmetry (X)
	for z := 0; z < size[Z]; z++ {
		for y := 0; y < size[Y]; y++ {
			for x := size[X]/2 + 1; x < size[X]; x++ {
				x2 := size[X] - x
				array[X][X][z][y][x] = array[X][X][z][y][x2]
				array[X][Y][z][y][x] = -array[X][Y][z][y][x2]
				array[X][Z][z][y][x] = -array[X][Z][z][y][x2]
				array[Y][Y][z][y][x] = array[Y][Y][z][y][x2]
				array[Y][Z][z][y][x] = array[Y][Z][z][y][x2]
				array[Z][Z][z][y][x] = array[Z][Z][z][y][x2]
			}
		}
	}

	// Reconstruct skipped parts from symmetry (Y)
	for z := 0; z < size[Z]; z++ {
		for y := size[Y]/2 + 1; y < size[Y]; y++ {
			y2 := size[Y] - y
			for x := 0; x < size[X]; x++ {
				array[X][X][z][y][x] = array[X][X][z][y2][x]
				array[X][Y][z][y][x] = -array[X][Y][z][y2][x]
				array[X][Z][z][y][x] = array[X][Z][z][y2][x]
				array[Y][Y][z][y][x] = array[Y][Y][z][y2][x]
				array[Y][Z][z][y][x] = -array[Y][Z][z][y2][x]
				array[Z][Z][z][y][x] = array[Z][Z][z][y2][x]

			}
		}
	}

	// Reconstruct skipped parts from symmetry (Z)
	for z := size[Z]/2 + 1; z < size[Z]; z++ {
		z2 := size[Z] - z
		for y := 0; y < size[Y]; y++ {
			for x := 0; x < size[X]; x++ {
				array[X][X][z][y][x] = array[X][X][z2][y][x]
				array[X][Y][z][y][x] = array[X][Y][z2][y][x]
				array[X][Z][z][y][x] = -array[X][Z][z2][y][x]
				array[Y][Y][z][y][x] = array[Y][Y][z2][y][x]
				array[Y][Z][z][y][x] = -array[Y][Z][z2][y][x]
				array[Z][Z][z][y][x] = array[Z][Z][z2][y][x]
			}
		}
	}

	// for 2D these elements are zero:
	if size[Z] == 1 {
		kernel[X][Z] = nil
		kernel[Y][Z] = nil
	}

	// make result symmetric for tools that expect it so.
	kernel[Y][X] = kernel[X][Y]
	kernel[Z][X] = kernel[X][Z]
	kernel[Z][Y] = kernel[Y][Z]

	ProgressBar.Finish()
	return kernel
}

// integration ranges for kernel. size=kernelsize, so padded for no PBC, not padded for PBC
func kernelRanges(size, pbc [3]int) (r1, r2 [3]int) {
	for c := 0; c < 3; c++ {
		if pbc[c] == 0 {
			r1[c], r2[c] = -(size[c]-1)/2, (size[c]-1)/2
		} else {
			r1[c], r2[c] = -(size[c]*pbc[c] - 1), (size[c]*pbc[c] - 1) // no /2 here, or we would take half right and half left image
		}
	}
	// support for 2D simulations (thickness 1)
	if size[Z] == 1 && pbc[Z] == 0 {
		r2[Z] = 0
	}
	return
}

const (
	X = 0
	Y = 1
	Z = 2
)

// closest distance between cells, given center distance d.
// if cells touch by just even a corner, the distance is zero.
func delta(d int) float64 {
	if d < 0 {
		d = -d
	}
	if d > 0 {
		d -= 1
	}
	return float64(d)
}

// Wraps an index to [0, max] by adding/subtracting a multiple of max.
func wrap(number, max int) int {
	for number < 0 {
		number += max
	}
	for number >= max {
		number -= max
	}
	return number
}

const maxAspect = 100.0 // maximum sane cell aspect ratio

func sanityCheck(cellsize [3]float64) {
	a3 := cellsize[X] / cellsize[Y]
	a2 := cellsize[Y] / cellsize[Z]
	a1 := cellsize[Z] / cellsize[X]

	aMax := math.Max(a1, math.Max(a2, a3))
	aMin := math.Min(a1, math.Min(a2, a3))

	if aMax > maxAspect || aMin < 1./maxAspect {
		util.Log.PanicIfError(fmt.Errorf("unrealistic cell aspect ratio %v", cellsize))
	}
}

// Returns the size after zero-padding, taking into account periodic boundary conditions.
// In a certain direction, there is no padding in case of PBC (it should wrap around).
// Without PBC there should be zero padding up to at least 2*N - 1. In that case there
// is a trade-off: for large N, padding up to 2*N can be much more efficient since
// power-of-two sized FFT's are ludicrously fast on CUDA. However for very small N,
// in particular N=1, we should not over-pad.
func padSize(gridsize, pbc [3]int) [3]int {
	var padded [3]int
	for i := range gridsize {
		if pbc[i] != 0 {
			padded[i] = gridsize[i]
			continue
		}
		if i != Z || gridsize[i] > SMALL_N { // for some reason it only works for Z, perhaps we assume even FFT size elsewhere?
			// large N: zero pad * 2 for FFT performance
			padded[i] = gridsize[i] * 2
		} else {
			// small N: minimal zero padding for memory/performance
			padded[i] = gridsize[i]*2 - 1
		}
	}
	return padded
}

// Use 2N-1 padding instead of 2N for sizes up to SMALL_N.
// 5 seems a good choice since for all n<=5, 2*n-1 only has
// prime factors 2,3,5,7 (good CUFFT performance).
// starting from 6 it becomes problematic so we use 2*n.
const SMALL_N = 5

// "If brute force doesn't solve your problem,
// you're not using enough of it."

/*
	Note: error for cubic self-kernel for different stagger decisions:

       1 ++--+----+-++---+----+--++---+----+-++---+----+--++---+----+-++--++
         +           +            +           +            +           +   +
         |                                                                 |
         +                 A                                               +
     0.1 ++                           A      A                            ++
         +                                         A    A   A   A          +
         +                       C                                 A  A  A +
    0.01 ++    B               D          E      C                        ++
 e       +                               B      D     E    C               +
 r       |                        F                  B    D   BE  C        |
 r       +                                   F                   D  BE DC B+
 o 0.001 ++                                                               ++
 r       +                                           F                     +
         +                                                 F               +
  0.0001 ++                                                             F +F
         +                                                          F      +
         |                                                      F          |
         +           +            +           +            +           +   +
   1e-05 ++--+----+-++---+----+--++---+----+-++---+----+--++---+----+-++--++
        100         1000        10000       100000       1e+06       1e+07
                                  evaluation points
A: no staggering

B:
	nv = ((nv + 1) / 2) * 2
	nw = ((nw + 1) / 2) * 2
	nx = ((nx+1)/2)*2 - 1
	ny = ((ny+1)/2)*2 - 1
	nz = ((nz+1)/2)*2 - 1

C:
	nv = ((nv + 1) / 2) * 2
	nw = ((nw + 1) / 2) * 2
	nx = ((nx+1)/2)*2 + 1
	ny = ((ny+1)/2)*2 + 1
	nz = ((nz+1)/2)*2 + 1

D:
	nv += 1
	nw += 1

E:
	nx += 1
	ny += 1
	nz += 1

F: best with accuracy 6
	nv *= 2
	nw *= 2
*/
