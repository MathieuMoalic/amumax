package new_engine

import (
	"reflect"

	"github.com/MathieuMoalic/amumax/src/cuda"
	"github.com/MathieuMoalic/amumax/src/data"
)

// Special buffered quantity to store Magnetization
// makes sure it's normalized etc.
type Magnetization struct {
	EngineState *EngineStateStruct
	buffer_     *data.Slice
}

func NewMagnetization(es *EngineStateStruct) *Magnetization {
	return &Magnetization{
		EngineState: es,
		buffer_:     nil,
	}
}

// func (m *Magnetization) GetRegionToString(region int) string {
// 	v := unslice(AverageOf(m))
// 	return fmt.Sprintf("(%g,%g,%g)", v[0], v[1], v[2])
// }

func (m *Magnetization) NComp() int          { return 3 }
func (m *Magnetization) Name() string        { return "m" }
func (m *Magnetization) Unit() string        { return "" }
func (m *Magnetization) Buffer() *data.Slice { return m.buffer_ } // todo: rename Gpu()?

// func (m *Magnetization) Comp(c int) ScalarField  { return comp(m, c) }
func (m *Magnetization) SetValue(v interface{})  { m.SetInShape(nil, v.(config)) }
func (m *Magnetization) InputType() reflect.Type { return reflect.TypeOf(config(nil)) }
func (m *Magnetization) Type() reflect.Type      { return reflect.TypeOf(new(Magnetization)) }
func (m *Magnetization) Eval() interface{}       { return m }

// func (m *Magnetization) average() []float64      { return sAverageMagnet(NormMag.Buffer()) }
// func (m *Magnetization) Average() data.Vector    { return unslice(m.average()) }
func (m *Magnetization) normalize() { cuda.Normalize(m.Buffer(), m.EngineState.Geometry.Gpu()) }

// allocate storage (not done by init, as mesh size may not yet be known then)
func (m *Magnetization) alloc() {
	m.buffer_ = cuda.NewSlice(3, m.EngineState.Mesh.Size())
	m.Set(randomMag()) // sane starting config
}

func (m *Magnetization) SetArray(src *data.Slice) {
	if src.Size() != m.EngineState.Mesh.Size() {
		src = data.Resample(src, m.EngineState.Mesh.Size())
	}
	data.Copy(m.Buffer(), src)
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
	return m.Buffer(), false
}

func (m *Magnetization) EvalTo(dst *data.Slice) {
	data.Copy(dst, m.buffer_)
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
// 		cuda.SetCell(m.Buffer(), c, ix, iy, iz, float32(v[c]/vNorm))
// 	}
// }

// // Get the value of one cell.
// func (m *Magnetization) GetCell(ix, iy, iz int) data.Vector {
// 	mx := float64(cuda.GetCell(m.Buffer(), X, ix, iy, iz))
// 	my := float64(cuda.GetCell(m.Buffer(), Y, ix, iy, iz))
// 	mz := float64(cuda.GetCell(m.Buffer(), Z, ix, iy, iz))
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
		region = universeInner
	}
	host := m.Buffer().HostCopy()
	h := host.Vectors()
	n := m.EngineState.Mesh.Size()

	for iz := 0; iz < n[Z]; iz++ {
		for iy := 0; iy < n[Y]; iy++ {
			for ix := 0; ix < n[X]; ix++ {
				r := m.EngineState.Utils.Index2Coord(ix, iy, iz)
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
// 	host := m.Buffer().HostCopy()
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
