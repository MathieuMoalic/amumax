package magnetization

import (
	"math"

	"github.com/MathieuMoalic/amumax/src/engine_old/cuda_old"
	"github.com/MathieuMoalic/amumax/src/engine_old/data_old"
	"github.com/MathieuMoalic/amumax/src/geometry"
	"github.com/MathieuMoalic/amumax/src/mag_config"
	"github.com/MathieuMoalic/amumax/src/mesh"
	"github.com/MathieuMoalic/amumax/src/shape"
)

// Special buffered quantity to store Magnetization
// makes sure it's normalized etc.
type Magnetization struct {
	mesh     *mesh.Mesh
	config   *mag_config.ConfigList
	geometry *geometry.Geometry
	Slice    *data_old.Slice
}

func (m *Magnetization) Init(mesh *mesh.Mesh, config *mag_config.ConfigList, geometry *geometry.Geometry) {
	m.mesh = mesh
	m.config = config
	m.geometry = geometry
}

// These methods are defined for the Quantity interface

func (m *Magnetization) Size() [3]int               { return m.Slice.Size() }
func (m *Magnetization) EvalTo(dst *data_old.Slice) { data_old.Copy(dst, m.Slice) }
func (m *Magnetization) NComp() int                 { return 3 }
func (m *Magnetization) Name() string               { return "m" }
func (m *Magnetization) Unit() string               { return "" }
func (m *Magnetization) Value() *data_old.Slice     { return m.Slice }

func (m *Magnetization) Average() []float64 {
	s := m.Slice
	geom := m.geometry.GetOrCreateGpuSlice()
	avg := make([]float64, s.NComp())
	for i := range avg {
		avg[i] = float64(cuda_old.Dot(s.Comp(i), geom)) / float64(cuda_old.Sum(geom))
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
func (m *Magnetization) Normalize() {
	cuda_old.Normalize(m.Slice, m.geometry.GetOrCreateGpuSlice())
}

// allocate storage (not done by init, as mesh size may not yet be known then)
func (m *Magnetization) InitializeBuffer() {
	m.Slice = cuda_old.NewSlice(3, m.mesh.Size())
	m.set(m.config.RandomMag()) // sane starting config
}

func (m *Magnetization) setArray(src *data_old.Slice) {
	if src.Size() != m.mesh.Size() {
		src = data_old.Resample(src, m.mesh.Size())
	}
	data_old.Copy(m.Slice, src)
	m.Normalize()
}

func (m *Magnetization) set(c mag_config.Config) {
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
func (m *Magnetization) setInShape(region shape.Shape, conf mag_config.Config) {
	if region == nil {
		region = shape.Universe
	}
	cpuSlice := m.Slice.HostCopy()
	vectors := cpuSlice.Vectors()
	Nx, Ny, Nz := m.mesh.GetNi()

	for iz := 0; iz < Nz; iz++ {
		for iy := 0; iy < Ny; iy++ {
			for ix := 0; ix < Nx; ix++ {
				r := m.mesh.Index2Coord(ix, iy, iz)
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
