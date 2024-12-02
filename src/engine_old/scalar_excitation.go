package engine_old

import (
	"fmt"
	"reflect"

	"github.com/MathieuMoalic/amumax/src/cuda_old"
	"github.com/MathieuMoalic/amumax/src/engine_old/data_old"
	"github.com/MathieuMoalic/amumax/src/engine_old/log_old"
	"github.com/MathieuMoalic/amumax/src/engine_old/mesh_old"
	"github.com/MathieuMoalic/amumax/src/engine_old/script_old"
)

// An excitation, typically field or current,
// can be defined region-wise plus extra mask*multiplier terms.
type scalarExcitation struct {
	name       string
	perRegion  regionwiseScalar // Region-based excitation
	extraTerms []mulmask        // add extra mask*multiplier terms
}

func newScalarExcitation(name, unit, desc string) *scalarExcitation {
	e := new(scalarExcitation)
	e.name = name
	e.perRegion.init("_"+name+"_perRegion", unit, desc, nil) // name starts with underscore: unexported
	declLValue(name, e, cat(desc, unit))
	return e
}
func (e *scalarExcitation) GetRegionToString(region int) string {
	return fmt.Sprintf("%g", e.perRegion.GetRegion(region))
}

func (p *scalarExcitation) MSlice() cuda_old.MSlice {
	buf, r := p.Slice()
	log_old.AssertMsg(r, "Failed to retrieve slice: invalid state in scalarExcitation.MSlice")
	return cuda_old.ToMSlice(buf)
}

func (e *scalarExcitation) AddTo(dst *data_old.Slice) {
	if !e.perRegion.isZero() {
		cuda_old.RegionAddS(dst, e.perRegion.gpuLUT1(), Regions.Gpu())
	}

	for _, t := range e.extraTerms {
		var mul float32 = 1
		if t.mul != nil {
			mul = float32(t.mul())
		}
		cuda_old.Madd2(dst, dst, t.mask, 1, mul)
	}
}

func (e *scalarExcitation) Slice() (*data_old.Slice, bool) {
	buf := cuda_old.Buffer(e.NComp(), e.Mesh().Size())
	cuda_old.Zero(buf)
	e.AddTo(buf)
	return buf, true
}

// After resizing the mesh, the extra terms don't fit the grid anymore
// and there is no reasonable way to resize them. So remove them and have
// the user re-add them.
func (e *scalarExcitation) RemoveExtraTerms() {
	if len(e.extraTerms) == 0 {
		return
	}

	// log.Log.Comment("REMOVING EXTRA TERMS FROM", e.Name())
	for _, m := range e.extraTerms {
		m.mask.Free()
	}
	e.extraTerms = nil
}

// Add an extra mask*multiplier term to the excitation.
func (e *scalarExcitation) Add(mask *data_old.Slice, f script_old.ScalarFunction) {
	var mul func() float64
	if f != nil {
		if isConst(f) {
			val := f.Float()
			mul = func() float64 {
				return val
			}
		} else {
			mul = func() float64 {
				return f.Float()
			}
		}
	}
	e.AddGo(mask, mul)
}

// An Add(mask, f) equivalent for Go use
func (e *scalarExcitation) AddGo(mask *data_old.Slice, mul func() float64) {
	if mask != nil {
		checkNaN(mask, e.Name()+".add()") // TODO: in more places
		mask = data_old.Resample(mask, e.Mesh().Size())
		mask = assureGPU(mask)
	}
	e.extraTerms = append(e.extraTerms, mulmask{mul, mask})
}

func (e *scalarExcitation) SetRegion(region int, f script_old.ScalarFunction) {
	e.perRegion.SetRegion(region, f)
}
func (e *scalarExcitation) SetValue(v interface{}) { e.perRegion.SetValue(v) }
func (e *scalarExcitation) Set(v float64)          { e.perRegion.setRegions(0, NREGION, []float64{v}) }

func (e *scalarExcitation) SetRegionFn(region int, f func() [3]float64) {
	e.perRegion.setFunc(region, region+1, func() []float64 {
		return slice(f())
	})
}

func (e *scalarExcitation) average() float64           { return qAverageUniverse(e)[0] }
func (e *scalarExcitation) Average() float64           { return e.average() }
func (e *scalarExcitation) IsUniform() bool            { return e.perRegion.IsUniform() }
func (e *scalarExcitation) Name() string               { return e.name }
func (e *scalarExcitation) Unit() string               { return e.perRegion.Unit() }
func (e *scalarExcitation) NComp() int                 { return e.perRegion.NComp() }
func (e *scalarExcitation) Mesh() *mesh_old.Mesh       { return GetMesh() }
func (e *scalarExcitation) Region(r int) *vOneReg      { return vOneRegion(e, r) }
func (e *scalarExcitation) Comp(c int) ScalarField     { return comp(e, c) }
func (e *scalarExcitation) Eval() interface{}          { return e }
func (e *scalarExcitation) Type() reflect.Type         { return reflect.TypeOf(new(scalarExcitation)) }
func (e *scalarExcitation) InputType() reflect.Type    { return script_old.ScalarFunction_t }
func (e *scalarExcitation) EvalTo(dst *data_old.Slice) { evalTo(e, dst) }
