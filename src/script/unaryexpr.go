package script

import (
	"go/ast"
	"go/token"
	"reflect"
)

func (w *World) compileUnaryExpr(n *ast.UnaryExpr) Expr {
	x := w.compileExpr(n.X)
	switch n.Op {
	default:
		panic(err(n.Pos(), "not allowed:", n.Op))
	case token.SUB:
		return &minus{typeConv(n.X.Pos(), x, float64t)}
	case token.NOT:
		return &not{typeConv(n.X.Pos(), x, boolt)}
	}
}

type minus struct{ x Expr }

func (m *minus) Type() reflect.Type { return float64t }
func (m *minus) Eval() any          { return -m.x.Eval().(float64) }
func (m *minus) Child() []Expr      { return []Expr{m.x} }
func (m *minus) Fix() Expr          { return &minus{m.x.Fix()} }

type not struct{ x Expr }

func (m *not) Type() reflect.Type { return boolt }
func (m *not) Eval() any          { return !m.x.Eval().(bool) }
func (m *not) Child() []Expr      { return []Expr{m.x} }
func (m *not) Fix() Expr          { return &not{m.x.Fix()} }
