package engine

import (
	"math/rand"

	"github.com/MathieuMoalic/amumax/src/cuda"
	"github.com/MathieuMoalic/amumax/src/data"
	"github.com/MathieuMoalic/amumax/src/utils"
)

type geometry struct {
	e          *engineState
	edgeSmooth int
	gpuSlice   *data.Slice
	shape      shape
}

func newGeom(engineState *engineState) *geometry {
	g := &geometry{e: engineState}
	g.e.script.RegisterFunction("SetGeom", g.setGeom)
	g.e.script.RegisterVariable("geom", g)
	return g
}

func (g *geometry) getOrCreateGpuSlice() *data.Slice {
	if g.gpuSlice == nil {
		g.gpuSlice = data.NilSlice(1, g.e.mesh.Size())
	}
	return g.gpuSlice
}

func (g *geometry) slice() (*data.Slice, bool) {
	s := g.getOrCreateGpuSlice()
	if s.IsNil() {
		buffer := cuda.Buffer(1, g.e.mesh.Size())
		cuda.Memset(buffer, 1)
		return buffer, true
	} else {
		return s, false
	}
}
func (g *geometry) NComp() int             { return 1 }
func (g *geometry) Size() [3]int           { return g.gpuSlice.Size() }
func (g *geometry) Name() string           { return "geom" }
func (g *geometry) Unit() string           { return "" }
func (g *geometry) EvalTo(dst *data.Slice) { data.Copy(dst, g.gpuSlice) }
func (g *geometry) Value() *data.Slice     { return g.gpuSlice }

func (g *geometry) Average() []float64 {
	s, r := g.slice()
	if r {
		defer cuda.Recycle(s)
	}
	return utils.AverageSlice(s)
}

func (g *geometry) setGeom(s shape) {

	if s == nil {
		// TODO: would be nice not to save volume if entirely filled
		s = g.e.shape.universeInner
	}

	g.shape = s
	if g.getOrCreateGpuSlice().IsNil() {
		g.gpuSlice = cuda.NewSlice(1, g.e.mesh.Size())
	}

	CpuSlice := data.NewSlice(1, g.getOrCreateGpuSlice().Size())
	array := CpuSlice.Scalars()
	mesh := g.e.mesh

	g.e.log.Info("Initializing geometry")
	empty := true
	for iz := 0; iz < mesh.Nz; iz++ {
		for iy := 0; iy < mesh.Ny; iy++ {
			for ix := 0; ix < mesh.Nx; ix++ {

				r := g.e.mesh.Index2Coord(ix, iy, iz)
				x0, y0, z0 := r[X], r[Y], r[Z]

				// check if center and all vertices lie inside or all outside
				allIn, allOut := true, true
				if s(x0, y0, z0) {
					allOut = false
				} else {
					allIn = false
				}

				if g.edgeSmooth != 0 { // center is sufficient if we're not really smoothing
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
		g.e.log.ErrAndExit("SetGeom: geometry completely empty")
	}

	data.Copy(g.gpuSlice, CpuSlice)
	// M inside geom but previously outside needs to be re-inited
	needupload := false
	geomlist := CpuSlice.Host()[0]
	mhost := g.e.magnetization.slice.HostCopy()
	m := mhost.Host()
	rng := rand.New(rand.NewSource(0))
	for i := range m[0] {
		if geomlist[i] != 0 {
			mx, my, mz := m[X][i], m[Y][i], m[Z][i]
			if mx == 0 && my == 0 && mz == 0 {
				needupload = true
				rnd := g.e.config.RandomDir(rng)
				m[X][i], m[Y][i], m[Z][i] = float32(rnd[X]), float32(rnd[Y]), float32(rnd[Z])
			}
		}
	}
	if needupload {
		data.Copy(g.e.magnetization.slice, mhost)
	}

	g.e.magnetization.normalize() // removes m outside vol
}

// Sample edgeSmooth^3 points inside the cell to estimate its volume.
func (g *geometry) cellVolume(ix, iy, iz int) float32 {
	r := g.e.mesh.Index2Coord(ix, iy, iz)
	x0, y0, z0 := r[X], r[Y], r[Z]

	c := g.e.mesh.CellSize()
	cx, cy, cz := c[X], c[Y], c[Z]
	s := g.shape
	var vol float32

	N := g.edgeSmooth
	S := float64(g.edgeSmooth)

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

func (g *geometry) shift(dx int) {
	// empty mask, nothing to do
	if g == nil || g.gpuSlice.IsNil() {
		return
	}

	// allocated mask: shift
	s := g.gpuSlice
	s2 := cuda.Buffer(1, g.e.mesh.Size())
	defer cuda.Recycle(s2)
	newv := float32(1) // initially fill edges with 1's
	cuda.ShiftX(s2, s, dx, newv, newv)
	data.Copy(s, s2)

	n := g.e.mesh.Size()
	x1, x2 := g.shiftDirtyRange(dx)

	for iz := 0; iz < n[Z]; iz++ {
		for iy := 0; iy < n[Y]; iy++ {
			for ix := x1; ix < x2; ix++ {
				r := g.e.mesh.Index2Coord(ix, iy, iz) // includes shift
				if !g.shape(r[X], r[Y], r[Z]) {
					cuda.SetCell(g.gpuSlice, 0, ix, iy, iz, 0) // a bit slowish, but hardly reached
				}
			}
		}
	}

}

func (g *geometry) shiftY(dy int) {
	// empty mask, nothing to do
	if g == nil || g.gpuSlice.IsNil() {
		return
	}

	// allocated mask: shift
	s := g.gpuSlice
	s2 := cuda.Buffer(1, g.e.mesh.Size())
	defer cuda.Recycle(s2)
	newv := float32(1) // initially fill edges with 1's
	cuda.ShiftY(s2, s, dy, newv, newv)
	data.Copy(s, s2)

	n := g.e.mesh.Size()
	y1, y2 := g.shiftDirtyRange(dy)

	for iz := 0; iz < n[Z]; iz++ {
		for ix := 0; ix < n[X]; ix++ {
			for iy := y1; iy < y2; iy++ {
				r := g.e.mesh.Index2Coord(ix, iy, iz) // includes shift
				if !g.shape(r[X], r[Y], r[Z]) {
					cuda.SetCell(g.gpuSlice, 0, ix, iy, iz, 0) // a bit slowish, but hardly reached
				}
			}
		}
	}

}

// x range that needs to be refreshed after shift over dx
func (g *geometry) shiftDirtyRange(dx int) (x1, x2 int) {
	nx := g.e.mesh.Size()[X]
	g.e.log.AssertMsg(dx != 0, "Invalid shift: dx must not be zero in shiftDirtyRange")

	if dx < 0 {
		x1 = nx + dx
		x2 = nx
	} else {
		x1 = 0
		x2 = dx
	}
	return
}
