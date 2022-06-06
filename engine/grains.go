package engine

import "github.com/MathieuMoalic/amumax/data"

func init() {
	DeclFunc("GrainsInShape", GrainsInShape, "GrainsInShape")
}

func NewVectorMask2(Nx, Ny, Nz int) *data.Slice {
	return data.NewSlice(3, [3]int{Nx, Ny, Nz})
}

func GrainsInShape(grainsize float64, numRegions, seed int) {
	SetBusy(true)
	defer SetBusy(false)

	t := newTesselation(grainsize, numRegions, int64(seed))
	regions.hist = append(regions.hist, t.RegionOf)

	regions.render(t.RegionOf)
}
