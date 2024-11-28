package new_engine

import (
	"github.com/MathieuMoalic/amumax/src/cuda"
	"github.com/MathieuMoalic/amumax/src/data"
	"github.com/MathieuMoalic/amumax/src/log"
)

// stores the region index for each cell
type regions struct {
	e          *engineState
	maxRegions int
	gpuBuffer  *cuda.Bytes                 // region data on GPU
	hist       []func(x, y, z float64) int // history of region set operations
	indices    map[int]bool
}

func newRegions(engineState *engineState) *regions {
	return &regions{e: engineState, maxRegions: 256, indices: make(map[int]bool), hist: make([]func(x, y, z float64) int, 0)}
}

func (r *regions) addIndex(i int) {
	r.indices[i] = true
}

// func (r *regions) getExistingIndices() []int {
// 	indices := make([]int, 0, len(r.indices))
// 	for i := range r.indices {
// 		indices = append(indices, i)
// 	}
// 	sort.Ints(indices)
// 	return indices
// }

// func (r *regions) redefine(startId, endId int) {
// 	// Loop through all cells, if their region ID matches startId, change it to endId
// 	n := r.e.mesh.Size()
// 	l := r.regionListCPU() // need to start from previous state
// 	arr := reshapeBytes(l, r.e.mesh.Size())

// 	for iz := 0; iz < n[Z]; iz++ {
// 		for iy := 0; iy < n[Y]; iy++ {
// 			for ix := 0; ix < n[X]; ix++ {
// 				if arr[iz][iy][ix] == byte(startId) {
// 					arr[iz][iy][ix] = byte(endId)
// 				}
// 			}
// 		}
// 	}
// 	r.gpuBuffer.Upload(l)
// }

// func (r *regions) shapeFromRegion(id int) shape {
// 	return func(x, y, z float64) bool {
// 		return r.get(data.Vector{x, y, z}) == id
// 	}
// }

func (r *regions) initializeBuffer() {
	r.gpuBuffer = cuda.NewBytes(r.e.mesh.NCell())
	r.defRegion(0, r.e.shape.universeInner)
}

// Define a region with id (0-255) to be inside the Shape.
func (r *regions) defRegion(id int, s shape) {
	r.defRegionId(id)
	f := func(x, y, z float64) int {
		if s(x, y, z) {
			return id
		} else {
			return -1
		}
	}
	r.render(f)
	r.addIndex(id)
	// regions.hist = append(regions.hist, f)
}

// // Redefine a region with a given ID to a new ID
// func (r *regions) redefRegion(startId, endId int) {
// 	// Checks validity of input region IDs
// 	r.defRegionId(startId)
// 	r.defRegionId(endId)

// 	hist_len := len(r.hist) // Only consider hist before this Redef to avoid recursion
// 	f := func(x, y, z float64) int {
// 		value := -1
// 		for i := hist_len - 1; i >= 0; i-- {
// 			f_other := r.hist[i]
// 			region := f_other(x, y, z)
// 			if region >= 0 {
// 				value = region
// 				break
// 			}
// 		}
// 		if value == startId {
// 			return endId
// 		} else {
// 			return value
// 		}
// 	}
// 	r.redefine(startId, endId)
// 	r.hist = append(r.hist, f)
// }

// renders (rasterizes) shape, filling it with region number #id, between x1 and x2
// TODO: a tidbit expensive
func (r *regions) render(f func(x, y, z float64) int) {
	mesh := r.e.mesh
	regionArray1D := r.regionListCPU() // need to start from previous state
	regionArray3D := reshapeBytes(regionArray1D, mesh.Size())

	for iz := 0; iz < mesh.Nz; iz++ {
		for iy := 0; iy < mesh.Ny; iy++ {
			for ix := 0; ix < mesh.Nx; ix++ {
				r := r.e.utils.Index2Coord(ix, iy, iz)
				region := f(r[X], r[Y], r[Z])
				if region >= 0 {
					regionArray3D[iz][iy][ix] = byte(region)
				}
			}
		}
	}
	r.gpuBuffer.Upload(regionArray1D)
}

// get the region for position R based on the history
func (r *regions) get(R data.Vector) int {
	// reverse order, last one set wins.
	for i := len(r.hist) - 1; i >= 0; i-- {
		f := r.hist[i]
		region := f(R[X], R[Y], R[Z])
		if region >= 0 {
			return region
		}
	}
	return 0
}

// func (r *regions) hostArray() [][][]byte {
// 	return reshapeBytes(r.regionListCPU(), r.e.mesh.Size())
// }

func (r *regions) regionListCPU() []byte {
	regionsList := make([]byte, r.e.mesh.NCell())
	r.gpuBuffer.Download(regionsList)
	return regionsList
}

// func (r *regions) defRegionCell(id int, x, y, z int) {
// 	r.defRegionId(id)
// 	index := data.Index(r.e.mesh.Size(), x, y, z)
// 	r.gpuBuffer.Set(index, byte(id))
// }

func (r *regions) Average() float64 {
	s, recycle := r.slice()
	if recycle {
		defer cuda.Recycle(s)
	}
	return averageSlice(s)[0]
}

// Set the region of one cell
func (r *regions) setCell(ix, iy, iz int, region int) {
	size := r.e.mesh.Size()
	i := data.Index(size, ix, iy, iz)
	r.gpuBuffer.Set(i, byte(region))
	r.addIndex(region)
}

// func (r *regions) getCell(ix, iy, iz int) int {
// 	size := r.e.mesh.Size()
// 	i := data.Index(size, ix, iy, iz)
// 	return int(r.gpuBuffer.Get(i))
// }

func (r *regions) defRegionId(id int) {
	if id < 0 || id > r.maxRegions {
		log.Log.ErrAndExit("region id should be 0 -%d, have: %d", r.maxRegions, id)
	}
	r.addIndex(id)
}

// // normalized volume (0..1) of region.
// // TODO: a tidbit too expensive
// func (r *Regions) volume(region_ int) float64 {
// 	region := byte(region_)
// 	vol := 0
// 	list := r.RegionListCPU()
// 	for _, reg := range list {
// 		if reg == region {
// 			vol++
// 		}
// 	}
// 	V := float64(vol) / float64(r.e.mesh.NCell())
// 	return V
// }

// Get the region data on GPU
func (r *regions) gpu() *cuda.Bytes {
	return r.gpuBuffer
}

// var unitMap regionwise // unit map used to output regions quantity

// func init() {
// 	unitMap.init(1, "unit", "", nil)
// 	for r := 0; r < NREGION; r++ {
// 		unitMap.setRegion(r, []float64{float64(r)})
// 	}
// }

// Get returns the regions as a slice of floats, so it can be output.
func (r *regions) slice() (*data.Slice, bool) {
	buf := cuda.Buffer(1, r.e.mesh.Size())
	// cuda.RegionDecode(buf, unitMap.gpuLUT1(), Regions.Gpu())
	return buf, true
}

// func (r *RegionsState) EvalTo(dst *data.Slice) { evalTo(r, dst) }

// var _ Quantity = &Regions

// Re-interpret a contiguous array as a multi-dimensional array of given size.
func reshapeBytes(array []byte, size [3]int) [][][]byte {
	Nxx, Nyy, Nzz := size[X], size[Y], size[Z]
	if Nxx*Nyy*Nzz != len(array) {
		log.Log.ErrAndExit("reshapeBytes: size mismatch")
	}
	sliced := make([][][]byte, Nzz)
	for i := range sliced {
		sliced[i] = make([][]byte, Nyy)
	}
	for i := range sliced {
		for j := range sliced[i] {
			sliced[i][j] = array[(i*Nyy+j)*Nxx : (i*Nyy+j)*Nxx+Nxx]
		}
	}
	return sliced
}

func (r *regions) shift(dx int) {
	// TODO: return if no regions defined
	r1 := r.gpu()
	r2 := cuda.NewBytes(r.e.mesh.NCell()) // TODO: somehow recycle
	defer r2.Free()
	newreg := byte(0) // new region at edge
	cuda.ShiftBytes(r2, r1, r.e.mesh, dx, newreg)
	r1.Copy(r2)

	n := r.e.mesh.Size()
	x1, x2 := r.e.utils.shiftDirtyRange(dx)

	for iz := 0; iz < n[Z]; iz++ {
		for iy := 0; iy < n[Y]; iy++ {
			for ix := x1; ix < x2; ix++ {
				i := r.e.utils.Index2Coord(ix, iy, iz) // includes shift
				reg := r.get(i)
				if reg != 0 {
					r.setCell(ix, iy, iz, reg) // a bit slowish, but hardly reached
				}
			}
		}
	}
}

func (r *regions) shiftY(dy int) {
	// TODO: return if no regions defined
	r1 := r.gpu()
	r2 := cuda.NewBytes(r.e.mesh.NCell()) // TODO: somehow recycle
	defer r2.Free()
	newreg := byte(0) // new region at edge
	cuda.ShiftBytesY(r2, r1, r.e.mesh, dy, newreg)
	r1.Copy(r2)

	n := r.e.mesh.Size()
	y1, y2 := r.e.utils.shiftDirtyRange(dy)

	for iz := 0; iz < n[Z]; iz++ {
		for ix := 0; ix < n[X]; ix++ {
			for iy := y1; iy < y2; iy++ {
				i := r.e.utils.Index2Coord(ix, iy, iz) // includes shift
				reg := r.get(i)
				if reg != 0 {
					r.setCell(ix, iy, iz, reg) // a bit slowish, but hardly reached
				}
			}
		}
	}
}
