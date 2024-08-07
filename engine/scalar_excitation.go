package engine

import (
	"fmt"
	"reflect"

	"github.com/MathieuMoalic/amumax/cuda"
	"github.com/MathieuMoalic/amumax/data"
	"github.com/MathieuMoalic/amumax/script"
	"github.com/MathieuMoalic/amumax/util"
)

// An excitation, typically field or current,
// can be defined region-wise plus extra mask*multiplier terms.
type ScalarExcitation struct {
	name       string
	perRegion  RegionwiseScalar // Region-based excitation
	extraTerms []mulmask        // add extra mask*multiplier terms
}

func NewScalarExcitation(name, unit, desc string) *ScalarExcitation {
	e := new(ScalarExcitation)
	e.name = name
	e.perRegion.init("_"+name+"_perRegion", unit, desc, nil) // name starts with underscore: unexported
	DeclLValue(name, e, cat(desc, unit))
	return e
}
func (e *ScalarExcitation) GetRegionToString(region int) string {
	return fmt.Sprintf("%g", e.perRegion.GetRegion(region))
}

func (p *ScalarExcitation) MSlice() cuda.MSlice {
	buf, r := p.Slice()
	util.Assert(r)
	return cuda.ToMSlice(buf)
}

func (e *ScalarExcitation) AddTo(dst *data.Slice) {
	if !e.perRegion.isZero() {
		cuda.RegionAddS(dst, e.perRegion.gpuLUT1(), Regions.Gpu())
	}

	for _, t := range e.extraTerms {
		var mul float32 = 1
		if t.mul != nil {
			mul = float32(t.mul())
		}
		cuda.Madd2(dst, dst, t.mask, 1, mul)
	}
}

func (e *ScalarExcitation) Slice() (*data.Slice, bool) {
	buf := cuda.Buffer(e.NComp(), e.Mesh().Size())
	cuda.Zero(buf)
	e.AddTo(buf)
	return buf, true
}

// After resizing the mesh, the extra terms don't fit the grid anymore
// and there is no reasonable way to resize them. So remove them and have
// the user re-add them.
func (e *ScalarExcitation) RemoveExtraTerms() {
	if len(e.extraTerms) == 0 {
		return
	}

	// LogOut("REMOVING EXTRA TERMS FROM", e.Name())
	for _, m := range e.extraTerms {
		m.mask.Free()
	}
	e.extraTerms = nil
}

// Add an extra mask*multiplier term to the excitation.
func (e *ScalarExcitation) Add(mask *data.Slice, f script.ScalarFunction) {
	var mul func() float64
	if f != nil {
		if IsConst(f) {
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
func (e *ScalarExcitation) AddGo(mask *data.Slice, mul func() float64) {
	if mask != nil {
		checkNaN(mask, e.Name()+".add()") // TODO: in more places
		mask = data.Resample(mask, e.Mesh().Size())
		mask = assureGPU(mask)
	}
	e.extraTerms = append(e.extraTerms, mulmask{mul, mask})
}

func (e *ScalarExcitation) SetRegion(region int, f script.ScalarFunction) {
	e.perRegion.SetRegion(region, f)
}
func (e *ScalarExcitation) SetValue(v interface{}) { e.perRegion.SetValue(v) }
func (e *ScalarExcitation) Set(v float64)          { e.perRegion.setRegions(0, NREGION, []float64{v}) }

func (e *ScalarExcitation) SetRegionFn(region int, f func() [3]float64) {
	e.perRegion.setFunc(region, region+1, func() []float64 {
		return slice(f())
	})
}

func (e *ScalarExcitation) average() float64        { return qAverageUniverse(e)[0] }
func (e *ScalarExcitation) Average() float64        { return e.average() }
func (e *ScalarExcitation) IsUniform() bool         { return e.perRegion.IsUniform() }
func (e *ScalarExcitation) Name() string            { return e.name }
func (e *ScalarExcitation) Unit() string            { return e.perRegion.Unit() }
func (e *ScalarExcitation) NComp() int              { return e.perRegion.NComp() }
func (e *ScalarExcitation) Mesh() *data.Mesh        { return GetMesh() }
func (e *ScalarExcitation) Region(r int) *vOneReg   { return vOneRegion(e, r) }
func (e *ScalarExcitation) Comp(c int) ScalarField  { return Comp(e, c) }
func (e *ScalarExcitation) Eval() interface{}       { return e }
func (e *ScalarExcitation) Type() reflect.Type      { return reflect.TypeOf(new(ScalarExcitation)) }
func (e *ScalarExcitation) InputType() reflect.Type { return script.ScalarFunction_t }
func (e *ScalarExcitation) EvalTo(dst *data.Slice)  { EvalTo(e, dst) }
