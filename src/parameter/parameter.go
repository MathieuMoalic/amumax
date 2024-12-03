package parameter

/*
parameters are region- and time dependent input values,
like material parameters.
*/

import (
	"math"
	"reflect"

	"github.com/MathieuMoalic/amumax/src/cuda"
	"github.com/MathieuMoalic/amumax/src/log"
	"github.com/MathieuMoalic/amumax/src/mesh"
	"github.com/MathieuMoalic/amumax/src/slice"
)

// TODO BELOW
var Regions RegionsInterface

type RegionsInterface interface{ Gpu() *cuda.Bytes }
type quantity interface {
	EvalTo(*slice.Slice)
	NComp() int
	Name() string
	Unit() string
	Average() []float64
}

// lValue is settable
type lValue interface {
	SetValue(interface{}) // assigns a new value
	Eval() interface{}    // evaluate and return result (nil for void)
	Type() reflect.Type   // type that can be assigned and will be returned by Eval
}

func declLValue(name string, value lValue, doc string) {}
func qAverageUniverse(q quantity) []float64            { return q.Average() }

var Time = 0.0

//	func init() {
//		addParameter("B_ext", B_ext, "External magnetic field (T)")
//	}
const NREGION = 256

// type inputValue interface{ GetRegionToString(region int) string }

// checks if a script expression contains t (time)
//
//	func isConst(e script_old.Expr) bool {
//		t := World.Resolve("t")
//		return !script_old.Contains(e, t)
//	}

func cat(desc, unit string) string {
	if unit == "" {
		return desc
	} else {
		return desc + " (" + unit + ")"
	}
}

var Params map[string]field

type field struct {
	Name        string           `json:"name"`
	Value       func(int) string `json:"value"`
	Description string           `json:"description"`
}

// func (p *regionwiseVector) Comp(c int) ScalarField  { return comp(p, c) }
// END OF TODO

// input Parameter, settable by user
type Parameter struct {
	lut
	upd_reg    [NREGION]func() []float64 // time-dependent values
	timestamp  float64                   // used not to double-evaluate f(t)
	children   []derived                 // derived parameters
	name, unit string
	mesh       *mesh.Mesh
}

// func addParameter(name string, value interface{}, doc string) {
// 	if Params == nil {
// 		Params = make(map[string]field)
// 	}
// 	if v, ok := value.(*regionwiseScalar); ok {
// 		Params[name] = field{
// 			name,
// 			v.GetRegionToString,
// 			doc,
// 		}
// 	} else if v, ok := value.(*regionwiseVector); ok {
// 		Params[name] = field{
// 			name,
// 			v.GetRegionToString,
// 			doc,
// 		}
// 	} else if v, ok := value.(*inputValue); ok {
// 		Params[name] = field{
// 			name,
// 			v.GetRegionToString,
// 			doc,
// 		}
// 	} else if v, ok := value.(*excitation); ok {
// 		Params[name] = field{
// 			name,
// 			v.GetRegionToString,
// 			doc,
// 		}
// 	} else if v, ok := value.(*scalarExcitation); ok {
// 		Params[name] = field{
// 			name,
// 			v.GetRegionToString,
// 			doc,
// 		}
// 	}
// }

func (p *Parameter) init(nComp int, name, unit string, children []derived) {
	p.lut.init(nComp, p)
	p.name = name
	p.unit = unit
	p.children = children
	p.timestamp = math.Inf(-1)
}

func (p *Parameter) MSlice() cuda.MSlice {
	p.log.Debug("mesh: %v", p.mesh)
	if p.IsUniform() {
		return cuda.MakeMSlice(slice.NilSlice(p.NComp(), p.mesh.Size()), p.getRegion(0))
	} else {
		buf, r := p.Slice()
		log.AssertMsg(r, "Failed to retrieve slice: invalid state in regionwise.MSlice")
		return cuda.ToMSlice(buf)
	}
}

func (p *Parameter) Name() string     { return p.name }
func (p *Parameter) Unit() string     { return p.unit }
func (p *Parameter) Mesh() *mesh.Mesh { return p.mesh }

func (p *Parameter) update() {
	if p.timestamp != Time {
		changed := false
		// update functions of time
		for r := 0; r < NREGION; r++ {
			updFunc := p.upd_reg[r]
			if updFunc != nil {
				p.bufset_(r, updFunc())
				changed = true
			}
		}
		p.timestamp = Time
		if changed {
			p.invalidate()
		}
	}
}

// set in one region
func (p *Parameter) SetRegion(region int, v []float64) {
	if region == -1 {
		p.SetUniform(v)
	} else {
		p.setRegions(region, region+1, v)
	}
}

// set in all regions
func (p *Parameter) SetUniform(v []float64) {
	p.setRegions(0, NREGION, v)
}

// set in regions r1..r2(excl)
func (p *Parameter) setRegions(r1, r2 int, v []float64) {
	log.AssertMsg(len(v) == len(p.cpu_buf), "Size mismatch: the length of v must match the length of p.cpu_buf in setRegions")
	log.AssertMsg(r1 < r2, "Invalid region range: r1 must be less than r2 (exclusive upper bound) in setRegions")

	for r := r1; r < r2; r++ {
		p.upd_reg[r] = nil
		p.bufset_(r, v)
	}
	p.invalidate()
}

func (p *Parameter) bufset_(region int, v []float64) {
	for c := range p.cpu_buf {
		p.cpu_buf[c][region] = float32(v[c])
	}
}

func (p *Parameter) setFunc(r1, r2 int, f func() []float64) {
	log.AssertMsg(r1 < r2, "Invalid region range: r1 must be less than r2 (exclusive upper bound) in setFunc")

	for r := r1; r < r2; r++ {
		p.upd_reg[r] = f
	}
	p.invalidate()
}

// mark my GPU copy and my children as invalid (need update)
func (p *Parameter) invalidate() {
	p.gpu_ok = false
	for _, c := range p.children {
		c.invalidate()
	}
}

func (p *Parameter) getRegion(region int) []float64 {
	cpu := p.cpuLUT()
	v := make([]float64, p.NComp())
	for i := range v {
		v[i] = float64(cpu[i][region])
	}
	return v
}

func (p *Parameter) IsUniform() bool {
	cpu := p.cpuLUT()
	v1 := p.getRegion(0)
	for r := 1; r < NREGION; r++ {
		for c := range v1 {
			if cpu[c][r] != float32(v1[c]) {
				return false
			}
		}
	}
	return true
}

func (p *Parameter) Average() []float64 { return qAverageUniverse(p) }
