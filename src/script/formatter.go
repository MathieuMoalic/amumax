package script

import (
	"fmt"
	"go/ast"
	"strings"
)

// formatExpr formats an expression into a string representation.
func (p *ScriptParser) formatExpr(expr ast.Expr) string {
	switch e := expr.(type) {
	case *ast.BasicLit:
		return e.Value
	case *ast.Ident:
		return e.Name
	case *ast.BinaryExpr:
		return fmt.Sprintf("%s %s %s", p.formatExpr(e.X), e.Op.String(), p.formatExpr(e.Y))
	case *ast.CallExpr:
		return p.formatCallExpr(e)
	case *ast.ParenExpr:
		return fmt.Sprintf("(%s)", p.formatExpr(e.X))
	case *ast.SelectorExpr:
		return fmt.Sprintf("%s.%s", p.formatExpr(e.X), e.Sel.Name)
	case *ast.IndexExpr:
		return fmt.Sprintf("%s[%s]", p.formatExpr(e.X), p.formatExpr(e.Index))
	case *ast.SliceExpr:
		return fmt.Sprintf("%s[%s:%s]", p.formatExpr(e.X), p.formatExpr(e.Low), p.formatExpr(e.High))
	case *ast.StarExpr:
		return fmt.Sprintf("*%s", p.formatExpr(e.X))
	case *ast.UnaryExpr:
		return fmt.Sprintf("%s%s", e.Op.String(), p.formatExpr(e.X))
	case *ast.TypeAssertExpr:
		return fmt.Sprintf("%s.(%s)", p.formatExpr(e.X), p.formatExpr(e.Type))
	case *ast.CompositeLit:
		elements := []string{}
		for _, elt := range e.Elts {
			elements = append(elements, p.formatExpr(elt))
		}
		return fmt.Sprintf("%s{%s}", p.formatExpr(e.Type), strings.Join(elements, ", "))
	case *ast.ArrayType:
		return fmt.Sprintf("[]%s", p.formatExpr(e.Elt))
	default:
		return fmt.Sprintf("%T", e)
	}
}

// Utility functions for formatting assignments and other statements

// formatCallExpr formats a function call expression with arguments.
func (p *ScriptParser) formatCallExpr(call *ast.CallExpr) string {
	funcName := p.formatExpr(call.Fun)
	args := make([]string, len(call.Args))
	for i, arg := range call.Args {
		args[i] = p.formatExpr(arg)
	}
	return fmt.Sprintf("%s(%s)", funcName, strings.Join(args, ", "))
}

// formatStmt formats a statement into a string.
func (p *ScriptParser) formatStmt(stmt ast.Stmt) string {
	switch s := stmt.(type) {
	case *ast.AssignStmt:
		return p.formatAssign(s)
	case *ast.IncDecStmt:
		return p.formatIncDec(s)
	case *ast.ExprStmt:
		return p.formatExpr(s.X)
	default:
		return ""
	}
}

// formatAssign formats an assignment statement into a string.
func (p *ScriptParser) formatAssign(assign *ast.AssignStmt) string {
	lhs := []string{}
	for _, l := range assign.Lhs {
		lhs = append(lhs, p.formatExpr(l))
	}
	rhs := []string{}
	for _, r := range assign.Rhs {
		rhs = append(rhs, p.formatExpr(r))
	}
	return fmt.Sprintf("%s %s %s", strings.Join(lhs, ", "), assign.Tok.String(), strings.Join(rhs, ", "))
}

// formatIncDec formats an increment/decrement statement.
func (p *ScriptParser) formatIncDec(stmt *ast.IncDecStmt) string {
	return fmt.Sprintf("%s%s", p.formatExpr(stmt.X), stmt.Tok.String())
}
