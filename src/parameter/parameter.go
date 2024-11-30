package parameter

// /*
// parameters are region- and time dependent input values,
// like material parameters.
// */

// import (
// 	"fmt"
// 	"math"
// 	"reflect"
// 	"strings"

// 	"github.com/MathieuMoalic/amumax/src/cuda"
// 	"github.com/MathieuMoalic/amumax/src/data"
// 	"github.com/MathieuMoalic/amumax/src/engine_old/log_old"
// 	"github.com/MathieuMoalic/amumax/src/mesh"
// 	"github.com/MathieuMoalic/amumax/src/engine_old/script_old"
// )

// // TODO BELOW
// var Regions RegionsInterface

// type RegionsInterface interface{ Gpu() *cuda.Bytes }
// type quantity interface {
// 	EvalTo(*data.Slice)
// 	NComp() int
// 	Name() string
// 	Unit() string
// }

// // lValue is settable
// type lValue interface {
// 	SetValue(interface{}) // assigns a new value
// 	Eval() interface{}    // evaluate and return result (nil for void)
// 	Type() reflect.Type   // type that can be assigned and will be returned by Eval
// }

// func declLValue(name string, value lValue, doc string) {}
// func qAverageUniverse(q quantity) []float64            { return nil }

// var Time = 0.0

// //	func init() {
// //		addParameter("B_ext", B_ext, "External magnetic field (T)")
// //	}
// const NREGION = 256

// // type inputValue interface{ GetRegionToString(region int) string }
// func unslice(v []float64) [3]float64 {
// 	log_old.AssertMsg(len(v) == 3, "Length mismatch: input slice must have exactly 3 elements in unslice")
// 	return [3]float64{v[0], v[1], v[2]}
// }

// // checks if a script expression contains t (time)
// //
// //	func isConst(e script_old.Expr) bool {
// //		t := World.Resolve("t")
// //		return !script_old.Contains(e, t)
// //	}
// func isConst(e script_old.Expr) bool {
// 	return true
// }

// func cat(desc, unit string) string {
// 	if unit == "" {
// 		return desc
// 	} else {
// 		return desc + " (" + unit + ")"
// 	}
// }

// var Params map[string]field

// // func (p *regionwiseVector) Comp(c int) ScalarField  { return comp(p, c) }
// // END OF TODO

// // input parameter, settable by user
// type parameter struct {
// 	lut
// 	upd_reg    [NREGION]func() []float64 // time-dependent values
// 	timestamp  float64                   // used not to double-evaluate f(t)
// 	children   []derived                 // derived parameters
// 	name, unit string
// 	mesh       *mesh.Mesh
// }

// type field struct {
// 	Name        string           `json:"name"`
// 	Value       func(int) string `json:"value"`
// 	Description string           `json:"description"`
// }

// // func addParameter(name string, value interface{}, doc string) {
// // 	if Params == nil {
// // 		Params = make(map[string]field)
// // 	}
// // 	if v, ok := value.(*regionwiseScalar); ok {
// // 		Params[name] = field{
// // 			name,
// // 			v.GetRegionToString,
// // 			doc,
// // 		}
// // 	} else if v, ok := value.(*regionwiseVector); ok {
// // 		Params[name] = field{
// // 			name,
// // 			v.GetRegionToString,
// // 			doc,
// // 		}
// // 	} else if v, ok := value.(*inputValue); ok {
// // 		Params[name] = field{
// // 			name,
// // 			v.GetRegionToString,
// // 			doc,
// // 		}
// // 	} else if v, ok := value.(*excitation); ok {
// // 		Params[name] = field{
// // 			name,
// // 			v.GetRegionToString,
// // 			doc,
// // 		}
// // 	} else if v, ok := value.(*scalarExcitation); ok {
// // 		Params[name] = field{
// // 			name,
// // 			v.GetRegionToString,
// // 			doc,
// // 		}
// // 	}
// // }

// func (p *parameter) init(nComp int, name, unit string, children []derived) {
// 	p.lut.init(nComp, p)
// 	p.name = name
// 	p.unit = unit
// 	p.children = children
// 	p.timestamp = math.Inf(-1)
// }

// func (p *parameter) MSlice() cuda.MSlice {
// 	if p.IsUniform() {
// 		return cuda.MakeMSlice(data.NilSlice(p.NComp(), p.mesh.Size()), p.getRegion(0))
// 	} else {
// 		buf, r := p.Slice()
// 		log_old.AssertMsg(r, "Failed to retrieve slice: invalid state in regionwise.MSlice")
// 		return cuda.ToMSlice(buf)
// 	}
// }

// func (p *parameter) Name() string     { return p.name }
// func (p *parameter) Unit() string     { return p.unit }
// func (p *parameter) Mesh() *mesh.Mesh { return p.mesh }

// func (p *parameter) update() {
// 	if p.timestamp != Time {
// 		changed := false
// 		// update functions of time
// 		for r := 0; r < NREGION; r++ {
// 			updFunc := p.upd_reg[r]
// 			if updFunc != nil {
// 				p.bufset_(r, updFunc())
// 				changed = true
// 			}
// 		}
// 		p.timestamp = Time
// 		if changed {
// 			p.invalidate()
// 		}
// 	}
// }

// // set in one region
// func (p *parameter) setRegion(region int, v []float64) {
// 	if region == -1 {
// 		p.setUniform(v)
// 	} else {
// 		p.setRegions(region, region+1, v)
// 	}
// }

// // set in all regions
// func (p *parameter) setUniform(v []float64) {
// 	p.setRegions(0, NREGION, v)
// }

// // set in regions r1..r2(excl)
// func (p *parameter) setRegions(r1, r2 int, v []float64) {
// 	log_old.AssertMsg(len(v) == len(p.cpu_buf), "Size mismatch: the length of v must match the length of p.cpu_buf in setRegions")
// 	log_old.AssertMsg(r1 < r2, "Invalid region range: r1 must be less than r2 (exclusive upper bound) in setRegions")

// 	for r := r1; r < r2; r++ {
// 		p.upd_reg[r] = nil
// 		p.bufset_(r, v)
// 	}
// 	p.invalidate()
// }

// func (p *parameter) bufset_(region int, v []float64) {
// 	for c := range p.cpu_buf {
// 		p.cpu_buf[c][region] = float32(v[c])
// 	}
// }

// func (p *parameter) setFunc(r1, r2 int, f func() []float64) {
// 	log_old.AssertMsg(r1 < r2, "Invalid region range: r1 must be less than r2 (exclusive upper bound) in setFunc")

// 	for r := r1; r < r2; r++ {
// 		p.upd_reg[r] = f
// 	}
// 	p.invalidate()
// }

// // mark my GPU copy and my children as invalid (need update)
// func (p *parameter) invalidate() {
// 	p.gpu_ok = false
// 	for _, c := range p.children {
// 		c.invalidate()
// 	}
// }

// func (p *parameter) getRegion(region int) []float64 {
// 	cpu := p.cpuLUT()
// 	v := make([]float64, p.NComp())
// 	for i := range v {
// 		v[i] = float64(cpu[i][region])
// 	}
// 	return v
// }

// func (p *parameter) IsUniform() bool {
// 	cpu := p.cpuLUT()
// 	v1 := p.getRegion(0)
// 	for r := 1; r < NREGION; r++ {
// 		for c := range v1 {
// 			if cpu[c][r] != float32(v1[c]) {
// 				return false
// 			}
// 		}
// 	}
// 	return true
// }

// func (p *parameter) average() []float64 { return qAverageUniverse(p) }

// // any parameter that depends on an inputParam
// type derived interface {
// 	invalidate()
// }

// // specialized param with 1 component
// type regionwiseScalar struct {
// 	parameter
// }

// func (p *regionwiseScalar) init(name, unit, desc string, children []derived) {
// 	p.parameter.init(1, name, unit, children)
// 	if !strings.HasPrefix(name, "_") { // don't export names beginning with "_" (e.g. from exciation)
// 		declLValue(name, p, cat(desc, unit))
// 	}
// }

// // TODO: auto derived
// func newScalarParam(name, unit, desc string, children ...derived) *regionwiseScalar {
// 	p := new(regionwiseScalar)
// 	p.parameter.init(1, name, unit, children)
// 	if !strings.HasPrefix(name, "_") { // don't export names beginning with "_" (e.g. from exciation)
// 		declLValue(name, p, cat(desc, unit))
// 	}
// 	return p
// }

// func (p *regionwiseScalar) SetRegion(region int, f script_old.ScalarFunction) {
// 	if region == -1 {
// 		p.setRegionsFunc(0, NREGION, f) // uniform
// 	} else {
// 		p.setRegionsFunc(region, region+1, f) // upper bound exclusive
// 	}
// }

// func (p *regionwiseScalar) SetValue(v interface{}) {
// 	f := v.(script_old.ScalarFunction)
// 	p.setRegionsFunc(0, NREGION, f)
// }

// func (p *regionwiseScalar) Set(v float64) {
// 	p.setRegions(0, NREGION, []float64{v})
// }

// func (p *regionwiseScalar) setRegionsFunc(r1, r2 int, f script_old.ScalarFunction) {
// 	if isConst(f) {
// 		p.setRegions(r1, r2, []float64{f.Float()})
// 	} else {
// 		f := f.Fix() // fix values of all variables except t
// 		p.setFunc(r1, r2, func() []float64 {
// 			return []float64{f.Eval().(script_old.ScalarFunction).Float()}
// 		})
// 	}
// }

// func (p *regionwiseScalar) GetRegion(region int) float64 {
// 	return float64(p.getRegion(region)[0])
// }
// func (p *regionwiseScalar) GetRegionToString(region int) string {
// 	v := float64(p.getRegion(region)[0])
// 	return fmt.Sprintf("%g", v)
// }

// func (p *regionwiseScalar) Eval() interface{}       { return p }
// func (p *regionwiseScalar) Type() reflect.Type      { return reflect.TypeOf(new(regionwiseScalar)) }
// func (p *regionwiseScalar) InputType() reflect.Type { return script_old.ScalarFunction_t }
// func (p *regionwiseScalar) Average() float64        { return qAverageUniverse(p)[0] }

// // func (p *regionwiseScalar) Region(r int) *sOneReg   { return sOneRegion(p, r) }

// // these methods should only be accesible from Go

// // vector input parameter, settable by user
// type regionwiseVector struct {
// 	parameter
// }

// func newVectorParam(name, unit, desc string) *regionwiseVector {
// 	p := new(regionwiseVector)
// 	p.parameter.init(3, name, unit, nil) // no vec param has children (yet)
// 	if !strings.HasPrefix(name, "_") {   // don't export names beginning with "_" (e.g. from exciation)
// 		declLValue(name, p, cat(desc, unit))
// 	}
// 	return p
// }

// func (p *regionwiseVector) SetRegion(region int, f script_old.VectorFunction) {
// 	if region == -1 {
// 		p.setRegionsFunc(0, NREGION, f) //uniform
// 	} else {
// 		p.setRegionsFunc(region, region+1, f)
// 	}
// }

// func (p *regionwiseVector) SetValue(v interface{}) {
// 	f := v.(script_old.VectorFunction)
// 	p.setRegionsFunc(0, NREGION, f)
// }

// func (p *regionwiseVector) setRegionsFunc(r1, r2 int, f script_old.VectorFunction) {
// 	if isConst(f) {
// 		p.setRegions(r1, r2, slice(f.Float3()))
// 	} else {
// 		f := f.Fix() // fix values of all variables except t
// 		p.setFunc(r1, r2, func() []float64 {
// 			return slice(f.Eval().(script_old.VectorFunction).Float3())
// 		})
// 	}
// }

// func (p *regionwiseVector) SetRegionFn(region int, f func() [3]float64) {
// 	p.setFunc(region, region+1, func() []float64 {
// 		return slice(f())
// 	})
// }
// func slice(v [3]float64) []float64 {
// 	return v[:]
// }

// func (p *regionwiseVector) GetRegion(region int) [3]float64 {
// 	v := p.getRegion(region)
// 	return unslice(v)
// }
// func (p *regionwiseVector) GetRegionToString(region int) string {
// 	v := unslice(p.getRegion(region))
// 	return fmt.Sprintf("(%g,%g,%g)", v[0], v[1], v[2])
// }
// func (p *regionwiseVector) Eval() interface{}       { return p }
// func (p *regionwiseVector) Type() reflect.Type      { return reflect.TypeOf(new(regionwiseVector)) }
// func (p *regionwiseVector) InputType() reflect.Type { return script_old.VectorFunction_t }

// // func (p *regionwiseVector) Region(r int) *vOneReg   { return vOneRegion(p, r) }
// func (p *regionwiseVector) Average() data.Vector { return unslice(qAverageUniverse(p)) }
