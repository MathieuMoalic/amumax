package script

import (
	"go/ast"
	"go/token"
	"reflect"
)

func (w *World) compileIncDecStmt(n *ast.IncDecStmt) Expr {
	l := w.compileLvalue(n.X)
	switch n.Tok {
	case token.INC:
		rhsPlus1 := &addone{incdec{typeConv(n.Pos(), l, float64t)}}
		return &assignStmt{lhs: l, rhs: typeConv(n.Pos(), rhsPlus1, l.Type())}
	case token.DEC:
		rhsMinus1 := &subone{incdec{typeConv(n.Pos(), l, float64t)}}
		return &assignStmt{lhs: l, rhs: typeConv(n.Pos(), rhsMinus1, l.Type())}
	default:
		panic(err(n.Pos(), "not allowed:", n.Tok))
	}
}

type incdec struct{ x Expr }

func (e *incdec) Type() reflect.Type { return float64t }
func (e *incdec) Child() []Expr      { return []Expr{e.x} }
func (e *incdec) Fix() Expr          { panic(invalidClosure) }

type (
	addone struct{ incdec }
	subone struct{ incdec }
)

func (s *addone) Eval() any { return s.x.Eval().(float64) + 1 }
func (s *subone) Eval() any { return s.x.Eval().(float64) - 1 }
