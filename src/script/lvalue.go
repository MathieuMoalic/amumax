package script

import (
	"go/ast"
	"reflect"
)

// LValue left-hand value in (single) assign statement
type LValue interface {
	Expr
	SetValue(any) // assigns a new value
}

func (w *World) compileLvalue(lhs ast.Node) LValue {
	switch lhs := lhs.(type) {
	default:
		panic(err(lhs.Pos(), "cannot assign to", typ(lhs)))
	case *ast.Ident:
		if l, ok := w.resolve(lhs.Pos(), lhs.Name).(LValue); ok {
			return l
		}
		panic(err(lhs.Pos(), "cannot assign to", lhs.Name))
	}
}

type reflectLvalue struct {
	elem reflect.Value
}

// newReflectLvalue general lvalue implementation using reflect.
// lhs must be settable, e.g. address of something:
//
//	var x float64
//	newReflectLValue(&x)
func newReflectLvalue(addr any) LValue {
	elem := reflect.ValueOf(addr).Elem()
	if elem.Kind() == 0 {
		panic("variable/constant needs to be passed as pointer to addressable value")
	}
	return &reflectLvalue{elem}
}

func (l *reflectLvalue) Eval() any {
	return l.elem.Interface()
}

func (l *reflectLvalue) Type() reflect.Type {
	return l.elem.Type()
}

func (l *reflectLvalue) SetValue(rvalue any) {
	l.elem.Set(reflect.ValueOf(rvalue))
}

func (l *reflectLvalue) Child() []Expr {
	return nil
}

func (l *reflectLvalue) Fix() Expr {
	return NewConst(l)
}

type TVar struct {
	LValue
}

func (t *TVar) Fix() Expr {
	return t // only variable that's not fixed
}
