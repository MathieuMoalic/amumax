package script

import (
	"fmt"
	"go/token"
	"reflect"

	"github.com/MathieuMoalic/amumax/src/data"
)

// converts in to an expression of type OutT.
// also serves as type check (not convertible == type error)
// pos is used for error message on impossible conversion.
func typeConv(pos token.Pos, in Expr, outT reflect.Type) Expr {
	inT := in.Type()
	switch {
	default:
		panic(err(pos, "type mismatch: can not use type", inT, "as", outT))

	// treat 'void' (type nil) separately:
	case inT == nil && outT != nil:
		panic(err(pos, "void used as value"))
	case inT != nil && outT == nil:
		panic("script internal bug: void input type")

	// strict go conversions:
	case inT == outT:
		return in
	case inT.AssignableTo(outT):
		return in

	// extra conversions for ease-of-use:
	// int -> float64
	case outT == float64t && inT == intt:
		return &intToFloat64{in}

	// float64 -> int
	case outT == intt && inT == float64t:
		return &float64ToInt{in}

	case outT == float64t && inT.AssignableTo(ScalarIft):
		return &getScalar{in.Eval().(ScalarIf)}
	case outT == float64t && inT.AssignableTo(VectorIft):
		return &getVector{in.Eval().(VectorIf)}

	// magical expression -> function conversions
	case inT == float64t && outT.AssignableTo(ScalarFunctiont):
		return &scalFn{in}
	case inT == intt && outT.AssignableTo(ScalarFunctiont):
		return &scalFn{&intToFloat64{in}}
	case inT == vectort && outT.AssignableTo(VectorFunctiont):
		return &vecFn{in}
	case inT == boolt && outT == funcboolt:
		return &boolToFunc{in}
	}
}

// returns input type for expression. Usually this is the same as the return type,
// unless the expression has a method InputType()reflect.Type.
func inputType(e Expr) reflect.Type {
	if in, ok := e.(interface {
		InputType() reflect.Type
	}); ok {
		return in.InputType()
	}
	return e.Type()
}

// common type definitions
var (
	float64t        = reflect.TypeOf(float64(0))
	boolt           = reflect.TypeOf(false)
	funcboolt       = reflect.TypeOf(func() bool { panic(0) })
	intt            = reflect.TypeOf(int(0))
	stringt         = reflect.TypeOf("")
	vectort         = reflect.TypeOf(data.Vector{})
	ScalarFunctiont = reflect.TypeOf(dummyf).In(0)
	VectorFunctiont = reflect.TypeOf(dummyf3).In(0)
	ScalarIft       = reflect.TypeOf(dummyscalarif).In(0)
	VectorIft       = reflect.TypeOf(dummyvectorif).In(0)
)

// maneuvers to get interface type of Func (simpler way?)
func dummyf(ScalarFunction)  {}
func dummyf3(VectorFunction) {}
func dummyscalarif(ScalarIf) {}
func dummyvectorif(VectorIf) {}

// converts int to float64
type intToFloat64 struct{ in Expr }

func (c *intToFloat64) Eval() any          { return float64(c.in.Eval().(int)) }
func (c *intToFloat64) Type() reflect.Type { return float64t }
func (c *intToFloat64) Child() []Expr      { return []Expr{c.in} }
func (c *intToFloat64) Fix() Expr          { return &intToFloat64{in: c.in.Fix()} }

// converts float64 to int
type float64ToInt struct{ in Expr }

func (c *float64ToInt) Eval() any          { return safeInt(c.in.Eval().(float64)) }
func (c *float64ToInt) Type() reflect.Type { return intt }
func (c *float64ToInt) Child() []Expr      { return []Expr{c.in} }
func (c *float64ToInt) Fix() Expr          { return &float64ToInt{in: c.in.Fix()} }

type boolToFunc struct{ in Expr }

func (c *boolToFunc) Eval() any          { return func() bool { return c.in.Eval().(bool) } }
func (c *boolToFunc) Type() reflect.Type { return funcboolt }
func (c *boolToFunc) Child() []Expr      { return []Expr{c.in} }
func (c *boolToFunc) Fix() Expr          { return &boolToFunc{in: c.in.Fix()} }

type (
	getScalar struct{ in ScalarIf }
	getVector struct{ in VectorIf }
)

func (c *getScalar) Eval() any          { return c.in.Get() }
func (c *getScalar) Type() reflect.Type { return float64t }
func (c *getScalar) Child() []Expr      { return nil }
func (c *getScalar) Fix() Expr          { return NewConst(c) }

func (c *getVector) Eval() any          { return c.in.Get() }
func (c *getVector) Type() reflect.Type { return vectort }
func (c *getVector) Child() []Expr      { return nil }
func (c *getVector) Fix() Expr          { return NewConst(c) }

func safeInt(x float64) int {
	i := int(x)
	if float64(i) != x {
		panic(fmt.Errorf("can not use %v as int", x))
	}
	return i
}

type ScalarIf interface {
	Get() float64
} // TODO: Scalar

type VectorIf interface {
	Get() data.Vector
} // TODO: Vector
