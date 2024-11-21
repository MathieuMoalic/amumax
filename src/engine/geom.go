package engine

import (
	"math/rand"

	"github.com/MathieuMoalic/amumax/src/cuda"
	"github.com/MathieuMoalic/amumax/src/data"
	"github.com/MathieuMoalic/amumax/src/log"
)

func init() {
	Geometry.init()
}

var (
	Geometry   GeometryType
	edgeSmooth int = 0 // disabled by default
)

type GeometryType struct {
	info
	Buffer *data.Slice
	shape  shape
}

func (g *GeometryType) init() {
	g.Buffer = nil
	g.info = info{1, "geom", ""}
	declROnly("geom", g, "Cell fill fraction (0..1)")
}

func (g *GeometryType) Gpu() *data.Slice {
	if g.Buffer == nil {
		g.Buffer = data.NilSlice(1, Mesh.Size())
	}
	return g.Buffer
}

func (g *GeometryType) Slice() (*data.Slice, bool) {
	s := g.Gpu()
	if s.IsNil() {
		buffer := cuda.Buffer(g.NComp(), Mesh.Size())
		cuda.Memset(buffer, 1)
		return buffer, true
	} else {
		return s, false
	}
}

func (q *GeometryType) EvalTo(dst *data.Slice) { evalTo(q, dst) }

var _ Quantity = &Geometry

func (g *GeometryType) average() []float64 {
	s, r := g.Slice()
	if r {
		defer cuda.Recycle(s)
	}
	return sAverageUniverse(s)
}

func (g *GeometryType) Average() float64 { return g.average()[0] }

func (g *GeometryType) setGeom(s shape) {
	setBusy(true)
	defer setBusy(false)
	CreateMesh()

	if s == nil {
		// TODO: would be nice not to save volume if entirely filled
		s = universeInner
	}

	g.shape = s
	if g.Gpu().IsNil() {
		g.Buffer = cuda.NewSlice(1, Mesh.Size())
	}

	host := data.NewSlice(1, g.Gpu().Size())
	array := host.Scalars()
	V := host
	v := array

	empty := true
	for iz := 0; iz < Mesh.Nz; iz++ {
		for iy := 0; iy < Mesh.Ny; iy++ {
			for ix := 0; ix < Mesh.Nx; ix++ {
				r := index2Coord(ix, iy, iz)
				x0, y0, z0 := r[X], r[Y], r[Z]

				// check if center and all vertices lie inside or all outside
				allIn, allOut := true, true
				if s(x0, y0, z0) {
					allOut = false
				} else {
					allIn = false
				}

				if edgeSmooth != 0 { // center is sufficient if we're not really smoothing
					for _, Δx := range []float64{-Mesh.Dx / 2, Mesh.Dx / 2} {
						for _, Δy := range []float64{-Mesh.Dy / 2, Mesh.Dy / 2} {
							for _, Δz := range []float64{-Mesh.Dz / 2, Mesh.Dz / 2} {
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
					v[iz][iy][ix] = 1
					empty = false
				case allOut:
					v[iz][iy][ix] = 0
				default:
					v[iz][iy][ix] = g.cellVolume(ix, iy, iz)
					empty = empty && (v[iz][iy][ix] == 0)
				}
			}
		}
	}
	if empty {
		log.Log.ErrAndExit("SetGeom: geometry completely empty")
	}

	data.Copy(g.Buffer, V)

	// M inside geom but previously outside needs to be re-inited
	needupload := false
	geomlist := host.Host()[0]
	mhost := normMag.Buffer().HostCopy()
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
		data.Copy(normMag.Buffer(), mhost)
	}

	normMag.normalize() // removes m outside vol
}

// Sample edgeSmooth^3 points inside the cell to estimate its volume.
func (g *GeometryType) cellVolume(ix, iy, iz int) float32 {
	r := index2Coord(ix, iy, iz)
	x0, y0, z0 := r[X], r[Y], r[Z]

	s := Geometry.shape
	var vol float32

	N := edgeSmooth
	S := float64(edgeSmooth)

	for dx := 0; dx < N; dx++ {
		Δx := -Mesh.Dx/2 + (Mesh.Dx / (2 * S)) + (Mesh.Dx/S)*float64(dx)
		for dy := 0; dy < N; dy++ {
			Δy := -Mesh.Dy/2 + (Mesh.Dy / (2 * S)) + (Mesh.Dy/S)*float64(dy)
			for dz := 0; dz < N; dz++ {
				Δz := -Mesh.Dz/2 + (Mesh.Dz / (2 * S)) + (Mesh.Dz/S)*float64(dz)

				if s(x0+Δx, y0+Δy, z0+Δz) { // inside
					vol++
				}
			}
		}
	}
	return vol / float32(N*N*N)
}

func (g *GeometryType) GetCell(ix, iy, iz int) float64 {
	return float64(cuda.GetCell(g.Gpu(), 0, ix, iy, iz))
}

func (g *GeometryType) shift(dx int) {
	// empty mask, nothing to do
	if g == nil || g.Buffer.IsNil() {
		return
	}

	// allocated mask: shift
	s := g.Buffer
	s2 := cuda.Buffer(1, Mesh.Size())
	defer cuda.Recycle(s2)
	newv := float32(1) // initially fill edges with 1's
	cuda.ShiftX(s2, s, dx, newv, newv)
	data.Copy(s, s2)

	x1, x2 := shiftDirtyRange(dx)

	for iz := 0; iz < Mesh.Nz; iz++ {
		for iy := 0; iy < Mesh.Ny; iy++ {
			for ix := x1; ix < x2; ix++ {
				r := index2Coord(ix, iy, iz) // includes shift
				if !g.shape(r[X], r[Y], r[Z]) {
					cuda.SetCell(g.Buffer, 0, ix, iy, iz, 0) // a bit slowish, but hardly reached
				}
			}
		}
	}

}

func (g *GeometryType) shiftY(dy int) {
	// empty mask, nothing to do
	if g == nil || g.Buffer.IsNil() {
		return
	}

	// allocated mask: shift
	s := g.Buffer
	s2 := cuda.Buffer(1, Mesh.Size())
	defer cuda.Recycle(s2)
	newv := float32(1) // initially fill edges with 1's
	cuda.ShiftY(s2, s, dy, newv, newv)
	data.Copy(s, s2)

	y1, y2 := shiftDirtyRange(dy)

	for iz := 0; iz < Mesh.Nz; iz++ {
		for ix := 0; ix < Mesh.Nx; ix++ {
			for iy := y1; iy < y2; iy++ {
				r := index2Coord(ix, iy, iz) // includes shift
				if !g.shape(r[X], r[Y], r[Z]) {
					cuda.SetCell(g.Buffer, 0, ix, iy, iz, 0) // a bit slowish, but hardly reached
				}
			}
		}
	}

}

// x range that needs to be refreshed after shift over dx
func shiftDirtyRange(dx int) (x1, x2 int) {
	nx := Mesh.Size()[X]
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
