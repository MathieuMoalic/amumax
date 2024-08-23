package engine

import (
	"sort"

	"github.com/MathieuMoalic/amumax/cuda"
	"github.com/MathieuMoalic/amumax/data"
	"github.com/MathieuMoalic/amumax/util"
)

var Regions = RegionsState{info: info{1, "regions", ""}, Indices: make(map[int]bool)} // global regions map

const NREGION = 256 // maximum number of regions, limited by size of byte.

func init() {
	DeclFunc("DefRegion", DefRegion, "Define a material region with given index (0-255) and shape")
	DeclFunc("ShapeFromRegion", ShapeFromRegion, "")
	DeclROnly("regions", &Regions, "Outputs the region index for each cell")
	DeclFunc("DefRegionCell", DefRegionCell, "Set a material region (first argument) in one cell "+
		"by the index of the cell (last three arguments)")
}

// stores the region index for each cell
type RegionsState struct {
	gpuBuffer *cuda.Bytes                 // region data on GPU
	hist      []func(x, y, z float64) int // history of region set operations
	info
	Indices map[int]bool
}

func (r *RegionsState) AddIndex(i int) {
	r.Indices[i] = true
}
func (r *RegionsState) GetExistingIndices() []int {
	indices := make([]int, 0, len(r.Indices))
	for i := range r.Indices {
		indices = append(indices, i)
	}
	sort.Ints(indices)
	return indices
}

func ShapeFromRegion(id int) Shape {
	return func(x, y, z float64) bool {
		return Regions.get(data.Vector{x, y, z}) == id
	}
}

func (r *RegionsState) alloc() {
	r.gpuBuffer = cuda.NewBytes(r.Mesh().NCell())
	DefRegion(0, universe)
}

// Define a region with id (0-255) to be inside the Shape.
func DefRegion(id int, s Shape) {
	defRegionId(id)
	f := func(x, y, z float64) int {
		if s(x, y, z) {
			return id
		} else {
			return -1
		}
	}
	Regions.render(f)
	Regions.AddIndex(id)
	// regions.hist = append(regions.hist, f)
}

// renders (rasterizes) shape, filling it with region number #id, between x1 and x2
// TODO: a tidbit expensive
func (r *RegionsState) render(f func(x, y, z float64) int) {
	mesh := GetMesh()
	regionArray1D := r.RegionListCPU() // need to start from previous state
	regionArray3D := reshapeBytes(regionArray1D, mesh.Size())

	for iz := 0; iz < mesh.Nz; iz++ {
		for iy := 0; iy < mesh.Ny; iy++ {
			for ix := 0; ix < mesh.Nx; ix++ {
				r := Index2Coord(ix, iy, iz)
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
func (r *RegionsState) get(R data.Vector) int {
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

func (r *RegionsState) HostArray() [][][]byte {
	return reshapeBytes(r.RegionListCPU(), r.Mesh().Size())
}

func (r *RegionsState) RegionListCPU() []byte {
	regionsList := make([]byte, r.Mesh().NCell())
	Regions.gpuBuffer.Download(regionsList)
	return regionsList
}

func DefRegionCell(id int, x, y, z int) {
	defRegionId(id)
	index := data.Index(GetMesh().Size(), x, y, z)
	Regions.gpuBuffer.Set(index, byte(id))
}

func (r *RegionsState) average() []float64 {
	s, recycle := r.Slice()
	if recycle {
		defer cuda.Recycle(s)
	}
	return sAverageUniverse(s)
}

func (r *RegionsState) Average() float64 { return r.average()[0] }

// Set the region of one cell
func (r *RegionsState) SetCell(ix, iy, iz int, region int) {
	size := GetMesh().Size()
	i := data.Index(size, ix, iy, iz)
	r.gpuBuffer.Set(i, byte(region))
	r.AddIndex(region)
}

func (r *RegionsState) GetCell(ix, iy, iz int) int {
	size := GetMesh().Size()
	i := data.Index(size, ix, iy, iz)
	return int(r.gpuBuffer.Get(i))
}

func defRegionId(id int) {
	if id < 0 || id > NREGION {
		util.Log.ErrAndExit("region id should be 0 -%d, have: %d", NREGION, id)
	}
	CreateMesh()
	Regions.AddIndex(id)
}

// normalized volume (0..1) of region.
// TODO: a tidbit too expensive
func (r *RegionsState) volume(region_ int) float64 {
	region := byte(region_)
	vol := 0
	list := r.RegionListCPU()
	for _, reg := range list {
		if reg == region {
			vol++
		}
	}
	V := float64(vol) / float64(r.Mesh().NCell())
	return V
}

// Get the region data on GPU
func (r *RegionsState) Gpu() *cuda.Bytes {
	return r.gpuBuffer
}

var unitMap regionwise // unit map used to output regions quantity

func init() {
	unitMap.init(1, "unit", "", nil)
	for r := 0; r < NREGION; r++ {
		unitMap.setRegion(r, []float64{float64(r)})
	}
}

// Get returns the regions as a slice of floats, so it can be output.
func (r *RegionsState) Slice() (*data.Slice, bool) {
	buf := cuda.Buffer(1, r.Mesh().Size())
	cuda.RegionDecode(buf, unitMap.gpuLUT1(), Regions.Gpu())
	return buf, true
}

func (r *RegionsState) EvalTo(dst *data.Slice) { EvalTo(r, dst) }

var _ Quantity = &Regions

// Re-interpret a contiguous array as a multi-dimensional array of given size.
func reshapeBytes(array []byte, size [3]int) [][][]byte {
	Nx, Ny, Nz := size[X], size[Y], size[Z]
	util.Argument(Nx*Ny*Nz == len(array))
	sliced := make([][][]byte, Nz)
	for i := range sliced {
		sliced[i] = make([][]byte, Ny)
	}
	for i := range sliced {
		for j := range sliced[i] {
			sliced[i][j] = array[(i*Ny+j)*Nx : (i*Ny+j)*Nx+Nx]
		}
	}
	return sliced
}

func (b *RegionsState) shift(dx int) {
	// TODO: return if no regions defined
	r1 := b.Gpu()
	r2 := cuda.NewBytes(b.Mesh().NCell()) // TODO: somehow recycle
	defer r2.Free()
	newreg := byte(0) // new region at edge
	cuda.ShiftBytes(r2, r1, b.Mesh(), dx, newreg)
	r1.Copy(r2)

	n := GetMesh().Size()
	x1, x2 := shiftDirtyRange(dx)

	for iz := 0; iz < n[Z]; iz++ {
		for iy := 0; iy < n[Y]; iy++ {
			for ix := x1; ix < x2; ix++ {
				r := Index2Coord(ix, iy, iz) // includes shift
				reg := b.get(r)
				if reg != 0 {
					b.SetCell(ix, iy, iz, reg) // a bit slowish, but hardly reached
				}
			}
		}
	}
}

func (b *RegionsState) shiftY(dy int) {
	// TODO: return if no regions defined
	r1 := b.Gpu()
	r2 := cuda.NewBytes(b.Mesh().NCell()) // TODO: somehow recycle
	defer r2.Free()
	newreg := byte(0) // new region at edge
	cuda.ShiftBytesY(r2, r1, b.Mesh(), dy, newreg)
	r1.Copy(r2)

	n := GetMesh().Size()
	y1, y2 := shiftDirtyRange(dy)

	for iz := 0; iz < n[Z]; iz++ {
		for ix := 0; ix < n[X]; ix++ {
			for iy := y1; iy < y2; iy++ {
				r := Index2Coord(ix, iy, iz) // includes shift
				reg := b.get(r)
				if reg != 0 {
					b.SetCell(ix, iy, iz, reg) // a bit slowish, but hardly reached
				}
			}
		}
	}
}

func (r *RegionsState) Mesh() *data.Mesh { return GetMesh() }

func prod(s [3]int) int {
	return s[0] * s[1] * s[2]
}
