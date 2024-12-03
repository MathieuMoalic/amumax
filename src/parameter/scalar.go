package parameter

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/MathieuMoalic/amumax/src/engine_old/script_old"
)

// any parameter that depends on an inputParam
type derived interface {
	invalidate()
}

// specialized param with 1 component
type ScalarParam struct {
	Parameter
}

func (p *ScalarParam) init(name, unit, desc string, children []derived) {
	p.Parameter.init(1, name, unit, children)
	if !strings.HasPrefix(name, "_") { // don't export names beginning with "_" (e.g. from exciation)
		declLValue(name, p, cat(desc, unit))
	}
}

// TODO: auto derived
func NewScalarParam(name, unit, desc string, children ...derived) *ScalarParam {
	p := new(ScalarParam)
	p.Parameter.init(1, name, unit, children)
	if !strings.HasPrefix(name, "_") { // don't export names beginning with "_" (e.g. from exciation)
		declLValue(name, p, cat(desc, unit))
	}
	return p
}

func (p *ScalarParam) SetRegion(region int, f script_old.ScalarFunction) {
	if region == -1 {
		p.setRegionsFunc(0, NREGION, f) // uniform
	} else {
		p.setRegionsFunc(region, region+1, f) // upper bound exclusive
	}
}

func (p *ScalarParam) SetValue(v interface{}) {
	f := v.(script_old.ScalarFunction)
	p.setRegionsFunc(0, NREGION, f)
}

func (p *ScalarParam) Set(v float64) {
	p.setRegions(0, NREGION, []float64{v})
}

func (p *ScalarParam) setRegionsFunc(r1, r2 int, f script_old.ScalarFunction) {
	p.setRegions(r1, r2, []float64{f.Float()})
	// if isConst(f) {
	// 	p.setRegions(r1, r2, []float64{f.Float()})
	// } else {
	// 	f := f.Fix() // fix values of all variables except t
	// 	p.setFunc(r1, r2, func() []float64 {
	// 		return []float64{f.Eval().(script_old.ScalarFunction).Float()}
	// 	})
	// }
}

func (p *ScalarParam) GetRegion(region int) float64 {
	return float64(p.getRegion(region)[0])
}
func (p *ScalarParam) GetRegionToString(region int) string {
	v := float64(p.getRegion(region)[0])
	return fmt.Sprintf("%g", v)
}

func (p *ScalarParam) Eval() interface{}  { return p }
func (p *ScalarParam) Type() reflect.Type { return reflect.TypeOf(new(ScalarParam)) }

// func (p *regionwiseScalar) InputType() reflect.Type { return script_old.ScalarFunction_t }
// func (p *regionwiseScalar) Average() float64        { return qAverageUniverse(p)[0] }

// func (p *regionwiseScalar) Region(r int) *sOneReg   { return sOneRegion(p, r) }
