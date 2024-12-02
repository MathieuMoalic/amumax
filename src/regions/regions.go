package regions

import (
	"sort"

	"github.com/MathieuMoalic/amumax/src/engine_old/cuda_old"
	"github.com/MathieuMoalic/amumax/src/engine_old/data_old"
	"github.com/MathieuMoalic/amumax/src/log"
	"github.com/MathieuMoalic/amumax/src/mesh"
	"github.com/MathieuMoalic/amumax/src/shape"
	"github.com/MathieuMoalic/amumax/src/utils"
)

// stores the region index for each cell
type Regions struct {
	mesh       *mesh.Mesh
	log        *log.Logs
	maxRegions int
	gpuBuffer  *cuda_old.Bytes             // region data on GPU
	hist       []func(x, y, z float64) int // history of region set operations
	indices    map[int]bool
}

func (r *Regions) Init(mesh *mesh.Mesh, log *log.Logs) {
	r.mesh = mesh
	r.log = log
	r.maxRegions = 256
	r.hist = make([]func(x, y, z float64) int, 0)
	r.indices = make(map[int]bool)
}

func (r *Regions) addIndex(i int) {
	r.indices[i] = true
}

func (r *Regions) Voronoi(minRegion, maxRegion int, getRegion func(float64, float64, float64) int) {
	r.hist = append(r.hist, getRegion)
	for i := minRegion; i < maxRegion; i++ {
		r.addIndex(i)
	}
	r.render(getRegion)
}

func (r *Regions) GetExistingIndices() []int {
	indices := make([]int, 0, len(r.indices))
	for i := range r.indices {
		indices = append(indices, i)
	}
	sort.Ints(indices)
	return indices
}

// func (r *regions) redefine(startId, endId int) {
// 	// Loop through all cells, if their region ID matches startId, change it to endId
// 	n := r.mesh.Size()
// 	l := r.regionListCPU() // need to start from previous state
// 	arr := reshapeBytes(l, r.mesh.Size())

// 	for iz := 0; iz < n[2]; iz++ {
// 		for iy := 0; iy < n[1]; iy++ {
// 			for ix := 0; ix < n[0]; ix++ {
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

func (r *Regions) InitializeBuffer() {
	r.gpuBuffer = cuda_old.NewBytes(r.mesh.NCell())
	r.defRegion(0, shape.Universe)
}

// Define a region with id (0-255) to be inside the Shape.
func (r *Regions) defRegion(id int, s shape.Shape) {
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
func (r *Regions) render(f func(x, y, z float64) int) {
	mesh := r.mesh
	regionArray1D := r.regionListCPU() // need to start from previous state
	regionArray3D := r.reshapeBytes(regionArray1D, mesh.Size())

	for iz := 0; iz < mesh.Nz; iz++ {
		for iy := 0; iy < mesh.Ny; iy++ {
			for ix := 0; ix < mesh.Nx; ix++ {
				r := r.mesh.Index2Coord(ix, iy, iz)
				region := f(r[0], r[1], r[2])
				if region >= 0 {
					regionArray3D[iz][iy][ix] = byte(region)
				}
			}
		}
	}
	r.gpuBuffer.Upload(regionArray1D)
}

// get the region for position R based on the history
func (r *Regions) get(R data_old.Vector) int {
	// reverse order, last one set wins.
	for i := len(r.hist) - 1; i >= 0; i-- {
		f := r.hist[i]
		region := f(R[0], R[1], R[2])
		if region >= 0 {
			return region
		}
	}
	return 0
}

// func (r *regions) hostArray() [][][]byte {
// 	return reshapeBytes(r.regionListCPU(), r.mesh.Size())
// }

func (r *Regions) regionListCPU() []byte {
	regionsList := make([]byte, r.mesh.NCell())
	r.gpuBuffer.Download(regionsList)
	return regionsList
}

// func (r *regions) defRegionCell(id int, x, y, z int) {
// 	r.defRegionId(id)
// 	index := data.Index(r.mesh.Size(), x, y, z)
// 	r.gpuBuffer.Set(index, byte(id))
// }

func (r *Regions) Average() float64 {
	s, recycle := r.slice()
	if recycle {
		defer cuda_old.Recycle(s)
	}
	return utils.AverageSlice(s)[0]
}

// Set the region of one cell
func (r *Regions) setCell(ix, iy, iz int, region int) {
	size := r.mesh.Size()
	i := data_old.Index(size, ix, iy, iz)
	r.gpuBuffer.Set(i, byte(region))
	r.addIndex(region)
}

// func (r *regions) getCell(ix, iy, iz int) int {
// 	size := r.mesh.Size()
// 	i := data.Index(size, ix, iy, iz)
// 	return int(r.gpuBuffer.Get(i))
// }

func (r *Regions) defRegionId(id int) {
	if id < 0 || id > r.maxRegions {
		r.log.ErrAndExit("region id should be 0 -%d, have: %d", r.maxRegions, id)
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
// 	V := float64(vol) / float64(r.mesh.NCell())
// 	return V
// }

// Get the region data on GPU
func (r *Regions) gpu() *cuda_old.Bytes {
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
func (r *Regions) slice() (*data_old.Slice, bool) {
	buf := cuda_old.Buffer(1, r.mesh.Size())
	// cuda.RegionDecode(buf, unitMap.gpuLUT1(), Regions.Gpu())
	return buf, true
}

// func (r *RegionsState) EvalTo(dst *data.Slice) { evalTo(r, dst) }

// var _ Quantity = &Regions

// Re-interpret a contiguous array as a multi-dimensional array of given size.
func (r *Regions) reshapeBytes(array []byte, size [3]int) [][][]byte {
	Nxx, Nyy, Nzz := size[0], size[1], size[2]
	if Nxx*Nyy*Nzz != len(array) {
		r.log.ErrAndExit("reshapeBytes: size mismatch")
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

func (r *Regions) Shift(dx int) {
	// TODO: return if no regions defined
	r1 := r.gpu()
	r2 := cuda_old.NewBytes(r.mesh.NCell()) // TODO: somehow recycle
	defer r2.Free()
	newreg := byte(0) // new region at edge
	cuda_old.ShiftBytes(r2, r1, r.mesh, dx, newreg)
	r1.Copy(r2)

	n := r.mesh.Size()
	x1, x2 := r.shiftDirtyRange(dx)

	for iz := 0; iz < n[2]; iz++ {
		for iy := 0; iy < n[1]; iy++ {
			for ix := x1; ix < x2; ix++ {
				i := r.mesh.Index2Coord(ix, iy, iz) // includes shift
				reg := r.get(i)
				if reg != 0 {
					r.setCell(ix, iy, iz, reg) // a bit slowish, but hardly reached
				}
			}
		}
	}
}

func (r *Regions) ShiftY(dy int) {
	// TODO: return if no regions defined
	r1 := r.gpu()
	r2 := cuda_old.NewBytes(r.mesh.NCell()) // TODO: somehow recycle
	defer r2.Free()
	newreg := byte(0) // new region at edge
	cuda_old.ShiftBytesY(r2, r1, r.mesh, dy, newreg)
	r1.Copy(r2)

	n := r.mesh.Size()
	y1, y2 := r.shiftDirtyRange(dy)

	for iz := 0; iz < n[2]; iz++ {
		for ix := 0; ix < n[0]; ix++ {
			for iy := y1; iy < y2; iy++ {
				i := r.mesh.Index2Coord(ix, iy, iz) // includes shift
				reg := r.get(i)
				if reg != 0 {
					r.setCell(ix, iy, iz, reg) // a bit slowish, but hardly reached
				}
			}
		}
	}
}

// x range that needs to be refreshed after shift over dx
func (r *Regions) shiftDirtyRange(dx int) (x1, x2 int) {
	Nx := r.mesh.Nx
	r.log.AssertMsg(dx != 0, "Invalid shift: dx must not be zero in shiftDirtyRange")
	if dx < 0 {
		x1 = Nx + dx
		x2 = Nx
	} else {
		x1 = 0
		x2 = dx
	}
	return
}
