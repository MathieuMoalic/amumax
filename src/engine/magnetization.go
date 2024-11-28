package engine

import (
	"math"

	"github.com/MathieuMoalic/amumax/src/cuda"
	"github.com/MathieuMoalic/amumax/src/data"
)

// Special buffered quantity to store magnetization
// makes sure it's normalized etc.
type magnetization struct {
	e     *engineState
	slice *data.Slice
}

func newMagnetization(es *engineState) *magnetization {
	m := &magnetization{
		e: es,
	}
	m.e.script.RegisterVariable("m", m)
	return m
}

// These methods are defined for the Quantity interface

func (m *magnetization) Size() [3]int           { return m.slice.Size() }
func (m *magnetization) EvalTo(dst *data.Slice) { data.Copy(dst, m.slice) }
func (m *magnetization) NComp() int             { return 3 }
func (m *magnetization) Name() string           { return "m" }
func (m *magnetization) Unit() string           { return "" }
func (m *magnetization) Value() *data.Slice     { return m.slice }

func (m *magnetization) Average() []float64 {
	s := m.slice
	geom := m.e.geometry.getOrCreateGpuSlice()
	avg := make([]float64, s.NComp())
	for i := range avg {
		avg[i] = float64(cuda.Dot(s.Comp(i), geom)) / float64(cuda.Sum(geom))
		if math.IsNaN(avg[i]) {
			panic("NaN")
		}
	}
	return avg
}

// These are the other methods

// func (m *Magnetization) GetRegionToString(region int) string {
// 	v := unslice(AverageOf(m))
// 	return fmt.Sprintf("(%g,%g,%g)", v[0], v[1], v[2])
// }

// func (m *Magnetization) Comp(c int) ScalarField  { return comp(m, c) }
// func (m *magnetization) setValue(v interface{})  { m.SetInShape(nil, v.(config)) }
// func (m *magnetization) inputType() reflect.Type { return reflect.TypeOf(config(nil)) }
// func (m *magnetization) getType() reflect.Type   { return reflect.TypeOf(new(magnetization)) }
// func (m *magnetization) eval() interface{}       { return m }

// func (m *Magnetization) Average() data.Vector    { return unslice(m.average()) }
func (m *magnetization) normalize() {
	cuda.Normalize(m.slice, m.e.geometry.getOrCreateGpuSlice())
}

// allocate storage (not done by init, as mesh size may not yet be known then)
func (m *magnetization) initializeBuffer() {
	m.slice = cuda.NewSlice(3, m.e.mesh.Size())
	m.set(m.e.config.randomMag()) // sane starting config
}

func (m *magnetization) setArray(src *data.Slice) {
	if src.Size() != m.e.mesh.Size() {
		src = data.Resample(src, m.e.mesh.Size())
	}
	data.Copy(m.slice, src)
	m.normalize()
}

func (m *magnetization) set(c config) {
	m.setInShape(nil, c)
}

// func (m *Magnetization) LoadFile(fname string) {
// 	m.SetArray(loadFile(fname))
// }
// func (m *Magnetization) LoadOvfFile(fname string) {
// 	m.SetArray(loadOvfFile(fname))
// }

// func (m *Magnetization) Region(r int) *vOneReg { return vOneRegion(m, r) }

// // Set the value of one cell.
// func (m *Magnetization) SetCell(ix, iy, iz int, v data.Vector) {
// 	r := index2Coord(ix, iy, iz)
// 	if Geometry.shape != nil && !Geometry.shape(r[X], r[Y], r[Z]) {
// 		return
// 	}
// 	vNorm := v.Len()
// 	for c := 0; c < 3; c++ {
// 		cuda.SetCell(m.buffer, c, ix, iy, iz, float32(v[c]/vNorm))
// 	}
// }

// // Get the value of one cell.
// func (m *Magnetization) GetCell(ix, iy, iz int) data.Vector {
// 	mx := float64(cuda.GetCell(m.buffer, X, ix, iy, iz))
// 	my := float64(cuda.GetCell(m.buffer, Y, ix, iy, iz))
// 	mz := float64(cuda.GetCell(m.buffer, Z, ix, iy, iz))
// 	return vector(mx, my, mz)
// }

// func (m *Magnetization) Quantity() []float64 { return slice(m.Average()) }
const (
	X = 0
	Y = 1
	Z = 2
)

// Sets the magnetization inside the shape
func (m *magnetization) setInShape(region shape, conf config) {
	if region == nil {
		region = m.e.shape.universeInner
	}
	cpuSlice := m.slice.HostCopy()
	vectors := cpuSlice.Vectors()
	Nx, Ny, Nz := m.e.mesh.GetNi()

	for iz := 0; iz < Nz; iz++ {
		for iy := 0; iy < Ny; iy++ {
			for ix := 0; ix < Nx; ix++ {
				r := m.e.utils.Index2Coord(ix, iy, iz)
				x, y, z := r[X], r[Y], r[Z]
				if region(x, y, z) { // inside
					m := conf(x, y, z)
					vectors[X][iz][iy][ix] = float32(m[X])
					vectors[Y][iz][iy][ix] = float32(m[Y])
					vectors[Z][iz][iy][ix] = float32(m[Z])
				}
			}
		}
	}
	m.setArray(cpuSlice)
}

// // set m to config in region
// func (m *Magnetization) SetRegion(region int, conf config) {
// 	host := m.buffer.HostCopy()
// 	h := host.Vectors()
// 	n := m.EngineState.Mesh.Size()
// 	r := byte(region)

// 	regionsArr := Regions.HostArray()

// 	for iz := 0; iz < n[Z]; iz++ {
// 		for iy := 0; iy < n[Y]; iy++ {
// 			for ix := 0; ix < n[X]; ix++ {
// 				pos := index2Coord(ix, iy, iz)
// 				x, y, z := pos[X], pos[Y], pos[Z]
// 				if regionsArr[iz][iy][ix] == r {
// 					m := conf(x, y, z)
// 					h[X][iz][iy][ix] = float32(m[X])
// 					h[Y][iz][iy][ix] = float32(m[Y])
// 					h[Z][iz][iy][ix] = float32(m[Z])
// 				}
// 			}
// 		}
// 	}
// 	m.SetArray(host)
// }
