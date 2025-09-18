package script

import (
	"go/ast"
)

// for statement
type forStmt struct {
	init, cond, post, body Expr
	void
}

var loopNestingCount int16

func init() {
	loopNestingCount = 0
}

func (p *forStmt) Eval() any {
	loopNestingCount++
	defer func() { loopNestingCount-- }() // Reset the flag after the loop

	for p.init.Eval(); p.cond.Eval().(bool); p.post.Eval() {
		p.body.Eval()
	}
	return nil // void
}

func (w *World) compileForStmt(n *ast.ForStmt) *forStmt {
	w.EnterScope()
	defer w.ExitScope()

	stmt := &forStmt{init: &nop{}, cond: &nop{}, post: &nop{}, body: &nop{}}
	if n.Init != nil {
		stmt.init = w.compileStmt(n.Init)
	}
	if n.Cond != nil {
		stmt.cond = typeConv(n.Cond.Pos(), w.compileExpr(n.Cond), boolt)
	} else {
		stmt.cond = boolLit(true)
	}
	if n.Post != nil {
		stmt.post = w.compileStmt(n.Post)
	}
	if n.Body != nil {
		stmt.body = w.compileBlockStmtNoScopeST(n.Body)
	}
	return stmt
}

type nop struct{ void }

func (e *nop) Child() []Expr { return nil }
func (e *nop) Eval() any     { return nil }
func (e *nop) Fix() Expr     { return e }

func (p *forStmt) Child() []Expr {
	return []Expr{p.init, p.cond, p.post, p.body}
}
