package parameter

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/MathieuMoalic/amumax/src/engine_old/script_old"
	"github.com/MathieuMoalic/amumax/src/utils"
)

// vector input parameter, settable by user
type VectorParam struct {
	Parameter
}

func NewVectorParam(name, unit, desc string) *VectorParam {
	p := new(VectorParam)
	p.Parameter.init(3, name, unit, nil) // no vec param has children (yet)
	if !strings.HasPrefix(name, "_") {   // don't export names beginning with "_" (e.g. from exciation)
		declLValue(name, p, cat(desc, unit))
	}
	return p
}

func (p *VectorParam) SetRegion(region int, f script_old.VectorFunction) {
	if region == -1 {
		p.setRegionsFunc(0, NREGION, f) //uniform
	} else {
		p.setRegionsFunc(region, region+1, f)
	}
}

func (p *VectorParam) SetValue(v interface{}) {
	f := v.(script_old.VectorFunction)
	p.setRegionsFunc(0, NREGION, f)
}

func (p *VectorParam) setRegionsFunc(r1, r2 int, f script_old.VectorFunction) {
	p.setRegions(r1, r2, utils.Slice(f.Float3()))
	// if isConst(f) {
	// 	p.setRegions(r1, r2, utils.Slice(f.Float3()))
	// } else {
	// 	f := f.Fix() // fix values of all variables except t
	// 	p.setFunc(r1, r2, func() []float64 {
	// 		return utils.Slice(f.Eval().(script_old.VectorFunction).Float3())
	// 	})
	// }
}

func (p *VectorParam) SetRegionFn(region int, f func() [3]float64) {
	p.setFunc(region, region+1, func() []float64 {
		return utils.Slice(f())
	})
}

func (p *VectorParam) GetRegion(region int) [3]float64 {
	v := p.getRegion(region)
	return utils.Unslice(v)
}
func (p *VectorParam) GetRegionToString(region int) string {
	v := utils.Unslice(p.getRegion(region))
	return fmt.Sprintf("(%g,%g,%g)", v[0], v[1], v[2])
}
func (p *VectorParam) Eval() interface{}  { return p }
func (p *VectorParam) Type() reflect.Type { return reflect.TypeOf(new(VectorParam)) }

// func (p *regionwiseVector) InputType() reflect.Type { return script_old.VectorFunction_t }

// func (p *regionwiseVector) Region(r int) *vOneReg   { return vOneRegion(p, r) }
// func (p *regionwiseVector) Average() vector.Vector { return utils.Unslice(qAverageUniverse(p)) }
