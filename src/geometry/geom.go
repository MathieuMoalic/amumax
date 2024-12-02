package geometry

import (
	"math/rand"

	"github.com/MathieuMoalic/amumax/src/cuda"
	"github.com/MathieuMoalic/amumax/src/log"
	"github.com/MathieuMoalic/amumax/src/mag_config"
	"github.com/MathieuMoalic/amumax/src/mesh"
	"github.com/MathieuMoalic/amumax/src/shape"
	"github.com/MathieuMoalic/amumax/src/slice"
)

type Geometry struct {
	mesh          *mesh.Mesh
	log           *log.Logs
	config        *mag_config.ConfigList
	mag_slice     *slice.Slice
	mag_Normalize func()
	EdgeSmooth    int
	gpuSlice      *slice.Slice
	shapeFunc     shape.Shape
}

func (g *Geometry) Init(mesh *mesh.Mesh, log *log.Logs, config *mag_config.ConfigList, mag_Normalize func()) {
	g.mesh = mesh
	g.log = log
	g.config = config
	g.mag_Normalize = mag_Normalize
	g.EdgeSmooth = 0
}

func (g *Geometry) InitializeBuffer(mag_slice *slice.Slice) {
	g.mag_slice = mag_slice
}

func (g *Geometry) GetOrCreateGpuSlice() *slice.Slice {
	if g.gpuSlice == nil {
		g.gpuSlice = slice.NilSlice(1, g.mesh.Size())
	}
	return g.gpuSlice
}

func (g *Geometry) slice() (*slice.Slice, bool) {
	s := g.GetOrCreateGpuSlice()
	if s.IsNil() {
		buffer := cuda.Buffer(1, g.mesh.Size())
		cuda.Memset(buffer, 1)
		return buffer, true
	} else {
		return s, false
	}
}
func (g *Geometry) NComp() int              { return 1 }
func (g *Geometry) Size() [3]int            { return g.gpuSlice.Size() }
func (g *Geometry) Name() string            { return "geom" }
func (g *Geometry) Unit() string            { return "" }
func (g *Geometry) EvalTo(dst *slice.Slice) { slice.Copy(dst, g.gpuSlice) }
func (g *Geometry) Value() *slice.Slice     { return g.gpuSlice }

func (g *Geometry) Average() []float64 {
	s, r := g.slice()
	if r {
		defer cuda.Recycle(s)
	}
	return cuda.AverageSlice(s)
}

func (g *Geometry) SetGeom(s shape.Shape) {

	if s == nil {
		// TODO: would be nice not to save volume if entirely filled
		s = shape.Universe
	}

	g.shapeFunc = s
	if g.GetOrCreateGpuSlice().IsNil() {
		g.gpuSlice = cuda.NewSlice(1, g.mesh.Size())
	}

	CpuSlice := slice.NewSlice(1, g.GetOrCreateGpuSlice().Size())
	array := CpuSlice.Scalars()
	mesh := g.mesh

	g.log.Info("Initializing geometry")
	empty := true
	for iz := 0; iz < mesh.Nz; iz++ {
		for iy := 0; iy < mesh.Ny; iy++ {
			for ix := 0; ix < mesh.Nx; ix++ {

				r := g.mesh.Index2Coord(ix, iy, iz)
				x0, y0, z0 := r[0], r[1], r[2]

				// check if center and all vertices lie inside or all outside
				allIn, allOut := true, true
				if s(x0, y0, z0) {
					allOut = false
				} else {
					allIn = false
				}

				if g.EdgeSmooth != 0 { // center is sufficient if we're not really smoothing
					for _, Δx := range []float64{-mesh.Dx / 2, mesh.Dx / 2} {
						for _, Δy := range []float64{-mesh.Dy / 2, mesh.Dy / 2} {
							for _, Δz := range []float64{-mesh.Dz / 2, mesh.Dz / 2} {
								if s(x0+Δx, y0+Δy, z0+Δz) { // inside
									allOut = false
								} else {
									allIn = false
								}
							}
						}
					}
				}

				switch {
				case allIn:
					array[iz][iy][ix] = 1
					empty = false
				case allOut:
					array[iz][iy][ix] = 0
				default:
					array[iz][iy][ix] = g.cellVolume(ix, iy, iz)
					empty = empty && (array[iz][iy][ix] == 0)
				}
			}
		}
	}

	if empty {
		g.log.ErrAndExit("SetGeom: geometry completely empty")
	}

	slice.Copy(g.gpuSlice, CpuSlice)
	// M inside geom but previously outside needs to be re-inited
	needupload := false
	geomlist := CpuSlice.Host()[0]
	g.log.Debug("2")
	mhost := g.mag_slice.HostCopy()
	g.log.Debug("3")
	m := mhost.Host()
	rng := rand.New(rand.NewSource(0))
	for i := range m[0] {
		if geomlist[i] != 0 {
			mx, my, mz := m[0][i], m[1][i], m[2][i]
			if mx == 0 && my == 0 && mz == 0 {
				needupload = true
				rnd := g.config.RandomDir(rng)
				m[0][i], m[1][i], m[2][i] = float32(rnd[0]), float32(rnd[1]), float32(rnd[2])
			}
		}
	}
	if needupload {
		slice.Copy(g.mag_slice, mhost)
	}

	g.mag_Normalize() // removes m outside vol
}

// Sample edgeSmooth^3 points inside the cell to estimate its volume.
func (g *Geometry) cellVolume(ix, iy, iz int) float32 {
	r := g.mesh.Index2Coord(ix, iy, iz)
	x0, y0, z0 := r[0], r[1], r[2]

	c := g.mesh.CellSize()
	cx, cy, cz := c[0], c[1], c[2]
	s := g.shapeFunc
	var vol float32

	N := g.EdgeSmooth
	S := float64(g.EdgeSmooth)

	for dx := 0; dx < N; dx++ {
		Δx := -cx/2 + (cx / (2 * S)) + (cx/S)*float64(dx)
		for dy := 0; dy < N; dy++ {
			Δy := -cy/2 + (cy / (2 * S)) + (cy/S)*float64(dy)
			for dz := 0; dz < N; dz++ {
				Δz := -cz/2 + (cz / (2 * S)) + (cz/S)*float64(dz)

				if s(x0+Δx, y0+Δy, z0+Δz) { // inside
					vol++
				}
			}
		}
	}
	return vol / float32(N*N*N)
}

// func (g *geometry) getCell(ix, iy, iz int) float64 {
// 	return float64(cuda.GetCell(g.getOrCreateGpuSlice(), 0, ix, iy, iz))
// }

func (g *Geometry) Shift(dx int) {
	// empty mask, nothing to do
	if g == nil || g.gpuSlice.IsNil() {
		return
	}

	// allocated mask: shift
	s := g.gpuSlice
	s2 := cuda.Buffer(1, g.mesh.Size())
	defer cuda.Recycle(s2)
	newv := float32(1) // initially fill edges with 1's
	cuda.ShiftX(s2, s, dx, newv, newv)
	slice.Copy(s, s2)

	n := g.mesh.Size()
	x1, x2 := g.shiftDirtyRange(dx)

	for iz := 0; iz < n[2]; iz++ {
		for iy := 0; iy < n[1]; iy++ {
			for ix := x1; ix < x2; ix++ {
				r := g.mesh.Index2Coord(ix, iy, iz) // includes shift
				if !g.shapeFunc(r[0], r[1], r[2]) {
					cuda.SetCell(g.gpuSlice, 0, ix, iy, iz, 0) // a bit slowish, but hardly reached
				}
			}
		}
	}

}

func (g *Geometry) ShiftY(dy int) {
	// empty mask, nothing to do
	if g == nil || g.gpuSlice.IsNil() {
		return
	}

	// allocated mask: shift
	s := g.gpuSlice
	s2 := cuda.Buffer(1, g.mesh.Size())
	defer cuda.Recycle(s2)
	newv := float32(1) // initially fill edges with 1's
	cuda.ShiftY(s2, s, dy, newv, newv)
	slice.Copy(s, s2)

	n := g.mesh.Size()
	y1, y2 := g.shiftDirtyRange(dy)

	for iz := 0; iz < n[2]; iz++ {
		for ix := 0; ix < n[0]; ix++ {
			for iy := y1; iy < y2; iy++ {
				r := g.mesh.Index2Coord(ix, iy, iz) // includes shift
				if !g.shapeFunc(r[0], r[1], r[2]) {
					cuda.SetCell(g.gpuSlice, 0, ix, iy, iz, 0) // a bit slowish, but hardly reached
				}
			}
		}
	}

}

// x range that needs to be refreshed after shift over dx
func (g *Geometry) shiftDirtyRange(dx int) (x1, x2 int) {
	nx := g.mesh.Size()[0]
	g.log.AssertMsg(dx != 0, "Invalid shift: dx must not be zero in shiftDirtyRange")

	if dx < 0 {
		x1 = nx + dx
		x2 = nx
	} else {
		x1 = 0
		x2 = dx
	}
	return
}
