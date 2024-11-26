package engine

import (
	"fmt"
	"reflect"

	"github.com/MathieuMoalic/amumax/src/cuda"
	"github.com/MathieuMoalic/amumax/src/data"
)

var NormMag magnetization // reduced magnetization (unit length)

// Special buffered quantity to store magnetization
// makes sure it's normalized etc.
type magnetization struct {
	buffer_ *data.Slice
}

func (m *magnetization) GetRegionToString(region int) string {
	v := unslice(AverageOf(m))
	return fmt.Sprintf("(%g,%g,%g)", v[0], v[1], v[2])
}

func (m *magnetization) Mesh() *data.MeshType { return GetMesh() }
func (m *magnetization) NComp() int           { return 3 }
func (m *magnetization) Name() string         { return "m" }
func (m *magnetization) Unit() string         { return "" }
func (m *magnetization) Buffer() *data.Slice  { return m.buffer_ } // todo: rename Gpu()?

func (m *magnetization) Comp(c int) ScalarField  { return comp(m, c) }
func (m *magnetization) SetValue(v interface{})  { m.SetInShape(nil, v.(config)) }
func (m *magnetization) InputType() reflect.Type { return reflect.TypeOf(config(nil)) }
func (m *magnetization) Type() reflect.Type      { return reflect.TypeOf(new(magnetization)) }
func (m *magnetization) Eval() interface{}       { return m }
func (m *magnetization) average() []float64      { return sAverageMagnet(NormMag.Buffer()) }
func (m *magnetization) Average() data.Vector    { return unslice(m.average()) }
func (m *magnetization) normalize()              { cuda.Normalize(m.Buffer(), Geometry.Gpu()) }

// allocate storage (not done by init, as mesh size may not yet be known then)
func (m *magnetization) Alloc() {
	m.buffer_ = cuda.NewSlice(3, m.Mesh().Size())
	m.Set(randomMag()) // sane starting config
}

func (b *magnetization) SetArray(src *data.Slice) {
	if src.Size() != b.Mesh().Size() {
		src = data.Resample(src, b.Mesh().Size())
	}
	data.Copy(b.Buffer(), src)
	b.normalize()
}

func (m *magnetization) Set(c config) {
	m.SetInShape(nil, c)
}

func (m *magnetization) LoadFile(fname string) {
	m.SetArray(loadFile(fname))
}
func (m *magnetization) LoadOvfFile(fname string) {
	m.SetArray(loadOvfFile(fname))
}

func (m *magnetization) Slice() (s *data.Slice, recycle bool) {
	return m.Buffer(), false
}

func (m *magnetization) EvalTo(dst *data.Slice) {
	data.Copy(dst, m.buffer_)
}

func (m *magnetization) Region(r int) *vOneReg { return vOneRegion(m, r) }

// Set the value of one cell.
func (m *magnetization) SetCell(ix, iy, iz int, v data.Vector) {
	r := index2Coord(ix, iy, iz)
	if Geometry.shape != nil && !Geometry.shape(r[X], r[Y], r[Z]) {
		return
	}
	vNorm := v.Len()
	for c := 0; c < 3; c++ {
		cuda.SetCell(m.Buffer(), c, ix, iy, iz, float32(v[c]/vNorm))
	}
}

// Get the value of one cell.
func (m *magnetization) GetCell(ix, iy, iz int) data.Vector {
	mx := float64(cuda.GetCell(m.Buffer(), X, ix, iy, iz))
	my := float64(cuda.GetCell(m.Buffer(), Y, ix, iy, iz))
	mz := float64(cuda.GetCell(m.Buffer(), Z, ix, iy, iz))
	return vector(mx, my, mz)
}

func (m *magnetization) Quantity() []float64 { return slice(m.Average()) }

// Sets the magnetization inside the shape
func (m *magnetization) SetInShape(region shape, conf config) {
	if region == nil {
		region = universeInner
	}
	host := m.Buffer().HostCopy()
	h := host.Vectors()
	n := m.Mesh().Size()

	for iz := 0; iz < n[Z]; iz++ {
		for iy := 0; iy < n[Y]; iy++ {
			for ix := 0; ix < n[X]; ix++ {
				r := index2Coord(ix, iy, iz)
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

// set m to config in region
func (m *magnetization) SetRegion(region int, conf config) {
	host := m.Buffer().HostCopy()
	h := host.Vectors()
	n := m.Mesh().Size()
	r := byte(region)

	regionsArr := Regions.HostArray()

	for iz := 0; iz < n[Z]; iz++ {
		for iy := 0; iy < n[Y]; iy++ {
			for ix := 0; ix < n[X]; ix++ {
				pos := index2Coord(ix, iy, iz)
				x, y, z := pos[X], pos[Y], pos[Z]
				if regionsArr[iz][iy][ix] == r {
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
