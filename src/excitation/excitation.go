package excitation

// import (
// 	"fmt"
// 	"math"
// 	"reflect"

// 	"github.com/MathieuMoalic/amumax/src/engine_old/data_old"
// 	"github.com/MathieuMoalic/amumax/src/engine_old/log_old"
// 	"github.com/MathieuMoalic/amumax/src/engine_old/script_old"
// 	"github.com/MathieuMoalic/amumax/src/parameter"
// )

// // An Excitation, typically field or current,
// // can be defined region-wise plus extra mask*multiplier terms.
// type Excitation struct {
// 	name       string
// 	perRegion  *parameter.VectorParam // Region-based excitation
// 	extraTerms []mulmask              // add extra mask*multiplier terms
// }

// // space-dependent mask plus time dependent multiplier
// type mulmask struct {
// 	mul  func() float64
// 	mask *data_old.Slice
// }

// func NewExcitation(name, unit, desc string) *Excitation {
// 	e := new(Excitation)
// 	e.name = name
// 	// e.perRegion.init(3, "_"+name+"_perRegion", unit, nil) // name starts with underscore: unexported
// 	e.perRegion = parameter.NewVectorParam(name+"_perRegion", unit, nil)
// 	// declLValue(name, e, cat(desc, unit))
// 	return e
// }

// // func newExcitation(name, unit, desc string) *excitation {
// // 	e := new(excitation)
// // 	e.name = name
// // 	e.perRegion.init(3, "_"+name+"_perRegion", unit, nil) // name starts with underscore: unexported
// // 	declLValue(name, e, cat(desc, unit))
// // 	return e
// // }

// // func (p *excitation) MSlice() cuda_old.MSlice {
// // 	buf, r := p.Slice()
// // 	log_old.AssertMsg(r, "Failed to retrieve slice: invalid state in excitation.MSlice")
// // 	return cuda_old.ToMSlice(buf)
// // }

// // func (e *excitation) AddTo(dst *data_old.Slice) {
// // 	if !e.perRegion.isZero() {
// // 		cuda_old.RegionAddV(dst, e.perRegion.gpuLUT(), Regions.Gpu())
// // 	}

// // 	for _, t := range e.extraTerms {
// // 		var mul float32 = 1
// // 		if t.mul != nil {
// // 			mul = float32(t.mul())
// // 		}
// // 		cuda_old.Madd2(dst, dst, t.mask, 1, mul)
// // 	}
// // }

// // func (e *excitation) isZero() bool {
// // 	return e.perRegion.isZero() && len(e.extraTerms) == 0
// // }

// // func (e *excitation) Slice() (*data_old.Slice, bool) {
// // 	buf := cuda_old.Buffer(e.NComp(), e.Mesh().Size())
// // 	cuda_old.Zero(buf)
// // 	e.AddTo(buf)
// // 	return buf, true
// // }

// // After resizing the mesh, the extra terms don't fit the grid anymore
// // and there is no reasonable way to resize them. So remove them and have
// // the user re-add them.
// func (e *Excitation) RemoveExtraTerms() {
// 	if len(e.extraTerms) == 0 {
// 		return
// 	}

// 	// log.Log.Comment("REMOVING EXTRA TERMS FROM", e.Name())
// 	for _, m := range e.extraTerms {
// 		m.mask.Free()
// 	}
// 	e.extraTerms = nil
// }

// // // Add an extra mask*multiplier term to the excitation.
// // func (e *excitation) Add(mask *data_old.Slice, f script_old.ScalarFunction) {
// // 	var mul func() float64
// // 	if f != nil {
// // 		if isConst(f) {
// // 			val := f.Float()
// // 			mul = func() float64 {
// // 				return val
// // 			}
// // 		} else {
// // 			mul = func() float64 {
// // 				return f.Float()
// // 			}
// // 		}
// // 	}
// // 	e.AddGo(mask, mul)
// // }

// // // An Add(mask, f) equivalent for Go use
// // func (e *excitation) AddGo(mask *data_old.Slice, mul func() float64) {
// // 	if mask != nil {
// // 		checkNaN(mask, e.Name()+".add()") // TODO: in more places
// // 		mask = data_old.Resample(mask, e.Mesh().Size())
// // 		mask = assureGPU(mask)
// // 	}
// // 	e.extraTerms = append(e.extraTerms, mulmask{mul, mask})
// // }

// func (e *Excitation) SetRegion(region int, f script_old.VectorFunction) {
// 	e.perRegion.SetRegion(region, f)
// }
// func (e *Excitation) SetValue(v interface{}) { e.perRegion.SetValue(v) }

// // func (e *excitation) Set(v data_old.Vector)  { e.perRegion.setRegions(0, NREGION, slice(v)) }

// // func (e *excitation) SetRegionFn(region int, f func() [3]float64) {
// // 	e.perRegion.setFunc(region, region+1, func() []float64 {
// // 		return slice(f())
// // 	})
// // }

// // func (e *excitation) average() []float64         { return qAverageUniverse(e) }
// // func (e *excitation) Average() data_old.Vector   { return unslice(qAverageUniverse(e)) }
// func (e *Excitation) IsUniform() bool { return e.perRegion.IsUniform() }
// func (e *Excitation) Name() string    { return e.name }
// func (e *Excitation) Unit() string    { return e.perRegion.Unit() }
// func (e *Excitation) NComp() int      { return e.perRegion.NComp() }

// // func (e *excitation) Mesh() *mesh_old.Mesh       { return GetMesh() }
// // func (e *excitation) Region(r int) *vOneReg      { return vOneRegion(e, r) }
// // func (e *excitation) Comp(c int) ScalarField     { return comp(e, c) }
// func (e *Excitation) Eval() interface{}       { return e }
// func (e *Excitation) Type() reflect.Type      { return reflect.TypeOf(new(Excitation)) }
// func (e *Excitation) InputType() reflect.Type { return script_old.VectorFunction_t }

// // func (e *excitation) EvalTo(dst *data_old.Slice) { evalTo(e, dst) }

// func (e *Excitation) GetRegionToString(region int) string {
// 	v := e.perRegion.GetRegion(region)
// 	return fmt.Sprintf("(%g,%g,%g)", v[0], v[1], v[2])
// }

// func checkNaN(s *data_old.Slice, name string) {
// 	h := s.Host()
// 	for _, h := range h {
// 		for _, v := range h {
// 			if math.IsNaN(float64(v)) || math.IsInf(float64(v), 0) {
// 				log_old.Log.ErrAndExit("NaN or Inf in %v", name)
// 			}
// 		}
// 	}
// }
