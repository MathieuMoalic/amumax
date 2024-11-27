package new_engine

import (
	"math"
	"reflect"

	"github.com/MathieuMoalic/amumax/src/cuda"
	"github.com/MathieuMoalic/amumax/src/data"
)

// Special buffered quantity to store Magnetization
// makes sure it's normalized etc.
type Magnetization struct {
	EngineState *EngineStateStruct
	slice       *data.Slice
}

func NewMagnetization(es *EngineStateStruct) *Magnetization {
	m := &Magnetization{
		EngineState: es,
	}
	m.EngineState.world.RegisterVariable("m", m)
	return m
}

// These methods are defined for the Quantity interface

func (m *Magnetization) Size() [3]int           { return m.slice.Size() }
func (m *Magnetization) EvalTo(dst *data.Slice) { data.Copy(dst, m.slice) }
func (m *Magnetization) NComp() int             { return 3 }
func (m *Magnetization) Name() string           { return "m" }
func (m *Magnetization) Unit() string           { return "" }
func (m *Magnetization) Value() *data.Slice     { return m.slice }

func (m *Magnetization) Average() []float64 {
	s := m.slice
	geom := m.EngineState.geometry.getOrCreateGpuSlice()
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
func (m *Magnetization) SetValue(v interface{})  { m.SetInShape(nil, v.(config)) }
func (m *Magnetization) InputType() reflect.Type { return reflect.TypeOf(config(nil)) }
func (m *Magnetization) Type() reflect.Type      { return reflect.TypeOf(new(Magnetization)) }
func (m *Magnetization) Eval() interface{}       { return m }

// func (m *Magnetization) Average() data.Vector    { return unslice(m.average()) }
func (m *Magnetization) normalize() {
	cuda.Normalize(m.slice, m.EngineState.geometry.getOrCreateGpuSlice())
}

// allocate storage (not done by init, as mesh size may not yet be known then)
func (m *Magnetization) InitializeBuffer() {
	m.slice = cuda.NewSlice(3, m.EngineState.mesh.Size())
	m.Set(randomMag()) // sane starting config
}

func (m *Magnetization) SetArray(src *data.Slice) {
	if src.Size() != m.EngineState.mesh.Size() {
		src = data.Resample(src, m.EngineState.mesh.Size())
	}
	data.Copy(m.slice, src)
	m.normalize()
}

func (m *Magnetization) Set(c config) {
	m.SetInShape(nil, c)
}

// func (m *Magnetization) LoadFile(fname string) {
// 	m.SetArray(loadFile(fname))
// }
// func (m *Magnetization) LoadOvfFile(fname string) {
// 	m.SetArray(loadOvfFile(fname))
// }

func (m *Magnetization) Slice() (s *data.Slice, recycle bool) {
	return m.slice, false
}

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
func (m *Magnetization) SetInShape(region shape, conf config) {
	if region == nil {
		region = m.EngineState.shape.universeInner
	}
	host := m.slice.HostCopy()
	h := host.Vectors()
	n := m.EngineState.mesh.Size()

	for iz := 0; iz < n[Z]; iz++ {
		for iy := 0; iy < n[Y]; iy++ {
			for ix := 0; ix < n[X]; ix++ {
				r := m.EngineState.utils.Index2Coord(ix, iy, iz)
				x, y, z := r[X], r[Y], r[Z]
				if region(x, y, z) { // inside
					m := conf(x, y, z)
					h[X][iz][iy][ix] = float32(m[X])
					h[Y][iz][iy][ix] = float32(m[Y])
					h[Z][iz][iy][ix] = float32(m[Z])
				}
			}
		}
	}
	m.SetArray(host)
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
