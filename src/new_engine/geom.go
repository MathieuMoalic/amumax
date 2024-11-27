package new_engine

import (
	"math/rand"

	"github.com/MathieuMoalic/amumax/src/cuda"
	"github.com/MathieuMoalic/amumax/src/data"
	"github.com/MathieuMoalic/amumax/src/log"
)

type Geometry struct {
	EngineState *EngineStateStruct
	edgeSmooth  int
	GpuSlice    *data.Slice
	shape       shape
}

func NewGeom(engineState *EngineStateStruct) *Geometry {
	g := &Geometry{EngineState: engineState}
	g.EngineState.world.RegisterFunction("SetGeom", g.setGeom)
	g.EngineState.world.RegisterVariable("geom", g)
	return g
}

func (g *Geometry) getOrCreateGpuSlice() *data.Slice {
	if g.GpuSlice == nil {
		g.GpuSlice = data.NilSlice(1, g.EngineState.mesh.Size())
	}
	return g.GpuSlice
}

func (g *Geometry) Slice() (*data.Slice, bool) {
	s := g.getOrCreateGpuSlice()
	if s.IsNil() {
		buffer := cuda.Buffer(1, g.EngineState.mesh.Size())
		cuda.Memset(buffer, 1)
		return buffer, true
	} else {
		return s, false
	}
}
func (g *Geometry) NComp() int             { return 1 }
func (g *Geometry) Size() [3]int           { return g.GpuSlice.Size() }
func (g *Geometry) Name() string           { return "geom" }
func (g *Geometry) Unit() string           { return "" }
func (g *Geometry) EvalTo(dst *data.Slice) { data.Copy(dst, g.GpuSlice) }
func (g *Geometry) Value() *data.Slice     { return g.GpuSlice }

func (g *Geometry) Average() []float64 {
	s, r := g.Slice()
	if r {
		defer cuda.Recycle(s)
	}
	return sAverageUniverse(s)
}

func (g *Geometry) setGeom(s shape) {

	if s == nil {
		// TODO: would be nice not to save volume if entirely filled
		s = g.EngineState.shape.universeInner
	}

	g.shape = s
	if g.getOrCreateGpuSlice().IsNil() {
		g.GpuSlice = cuda.NewSlice(1, g.EngineState.mesh.Size())
	}

	CpuSlice := data.NewSlice(1, g.getOrCreateGpuSlice().Size())
	array := CpuSlice.Scalars()
	mesh := g.EngineState.mesh

	log.Log.Info("Initializing geometry")
	empty := true
	for iz := 0; iz < mesh.Nz; iz++ {
		for iy := 0; iy < mesh.Ny; iy++ {
			for ix := 0; ix < mesh.Nx; ix++ {

				r := g.EngineState.utils.Index2Coord(ix, iy, iz)
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
		log.Log.ErrAndExit("SetGeom: geometry completely empty")
	}

	data.Copy(g.GpuSlice, CpuSlice)
	// M inside geom but previously outside needs to be re-inited
	needupload := false
	geomlist := CpuSlice.Host()[0]
	mhost := g.EngineState.magnetization.slice.HostCopy()
	m := mhost.Host()
	rng := rand.New(rand.NewSource(0))
	for i := range m[0] {
		if geomlist[i] != 0 {
			mx, my, mz := m[X][i], m[Y][i], m[Z][i]
			if mx == 0 && my == 0 && mz == 0 {
				needupload = true
				rnd := randomDir(rng)
				m[X][i], m[Y][i], m[Z][i] = float32(rnd[X]), float32(rnd[Y]), float32(rnd[Z])
			}
		}
	}
	if needupload {
		data.Copy(g.EngineState.magnetization.slice, mhost)
	}

	g.EngineState.magnetization.normalize() // removes m outside vol
}

// Sample edgeSmooth^3 points inside the cell to estimate its volume.
func (g *Geometry) cellVolume(ix, iy, iz int) float32 {
	r := g.EngineState.utils.Index2Coord(ix, iy, iz)
	x0, y0, z0 := r[X], r[Y], r[Z]

	c := g.EngineState.mesh.CellSize()
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

func (g *Geometry) GetCell(ix, iy, iz int) float64 {
	return float64(cuda.GetCell(g.getOrCreateGpuSlice(), 0, ix, iy, iz))
}

func (g *Geometry) shift(dx int) {
	// empty mask, nothing to do
	if g == nil || g.GpuSlice.IsNil() {
		return
	}

	// allocated mask: shift
	s := g.GpuSlice
	s2 := cuda.Buffer(1, g.EngineState.mesh.Size())
	defer cuda.Recycle(s2)
	newv := float32(1) // initially fill edges with 1's
	cuda.ShiftX(s2, s, dx, newv, newv)
	data.Copy(s, s2)

	n := g.EngineState.mesh.Size()
	x1, x2 := g.shiftDirtyRange(dx)

	for iz := 0; iz < n[Z]; iz++ {
		for iy := 0; iy < n[Y]; iy++ {
			for ix := x1; ix < x2; ix++ {
				r := g.EngineState.utils.Index2Coord(ix, iy, iz) // includes shift
				if !g.shape(r[X], r[Y], r[Z]) {
					cuda.SetCell(g.GpuSlice, 0, ix, iy, iz, 0) // a bit slowish, but hardly reached
				}
			}
		}
	}

}

func (g *Geometry) shiftY(dy int) {
	// empty mask, nothing to do
	if g == nil || g.GpuSlice.IsNil() {
		return
	}

	// allocated mask: shift
	s := g.GpuSlice
	s2 := cuda.Buffer(1, g.EngineState.mesh.Size())
	defer cuda.Recycle(s2)
	newv := float32(1) // initially fill edges with 1's
	cuda.ShiftY(s2, s, dy, newv, newv)
	data.Copy(s, s2)

	n := g.EngineState.mesh.Size()
	y1, y2 := g.shiftDirtyRange(dy)

	for iz := 0; iz < n[Z]; iz++ {
		for ix := 0; ix < n[X]; ix++ {
			for iy := y1; iy < y2; iy++ {
				r := g.EngineState.utils.Index2Coord(ix, iy, iz) // includes shift
				if !g.shape(r[X], r[Y], r[Z]) {
					cuda.SetCell(g.GpuSlice, 0, ix, iy, iz, 0) // a bit slowish, but hardly reached
				}
			}
		}
	}

}

// x range that needs to be refreshed after shift over dx
func (g *Geometry) shiftDirtyRange(dx int) (x1, x2 int) {
	nx := g.EngineState.mesh.Size()[X]
	log.AssertMsg(dx != 0, "Invalid shift: dx must not be zero in shiftDirtyRange")

	if dx < 0 {
		x1 = nx + dx
		x2 = nx
	} else {
		x1 = 0
		x2 = dx
	}
	return
}
