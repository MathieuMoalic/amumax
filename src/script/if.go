package script

import (
	"go/ast"
)

// if statement
type ifStmt struct {
	cond, body, elseExpression Expr
	void
}

func (b *ifStmt) Eval() any {
	if b.cond.Eval().(bool) {
		b.body.Eval()
	} else {
		if b.elseExpression != nil {
			b.elseExpression.Eval()
		}
	}
	return nil // void
}

func (w *World) compileIfStmt(n *ast.IfStmt) *ifStmt {
	w.EnterScope()
	defer w.ExitScope()

	stmt := &ifStmt{
		cond: typeConv(n.Cond.Pos(), w.compileExpr(n.Cond), boolt),
		body: w.compileBlockStmtNoScopeST(n.Body),
	}
	if n.Else != nil {
		stmt.elseExpression = w.compileStmt(n.Else)
	}

	return stmt
}

func (b *ifStmt) Child() []Expr {
	child := []Expr{b.cond, b.body, b.elseExpression}
	if b.elseExpression == nil {
		child = child[:2]
	}
	return child
}
