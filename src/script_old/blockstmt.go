package script_old

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/token"
	"reflect"
	"strings"

	"github.com/MathieuMoalic/amumax/src/log_old"
)

// block statement is a list of statements.
type BlockStmt struct {
	Children []Expr
	Node     []ast.Node
}

// does not enter scope because it does not necessarily needs to (e.g. for, if).
func (w *World) compileBlockStmt_noScope(n *ast.BlockStmt) *BlockStmt {
	b := &BlockStmt{}
	for _, s := range n.List {
		b.append(w.compileStmt(s), s)
	}
	return b
}

func (b *BlockStmt) append(s Expr, n ast.Node) {
	b.Children = append(b.Children, s)
	b.Node = append(b.Node, n)
}

func (b *BlockStmt) Eval() interface{} {
	for _, s := range b.Children {
		s.Eval()
	}
	return nil
}

func (b *BlockStmt) Type() reflect.Type {
	return nil
}

func (b *BlockStmt) Child() []Expr {
	return b.Children
}

func Format(n ast.Node) string {
	var buf bytes.Buffer
	fset := token.NewFileSet()
	log_old.Log.PanicIfError(format.Node(&buf, fset, n))
	str := buf.String()
	str = strings.TrimSuffix(str, "\n")
	return str
}

func (b *BlockStmt) Format() string {
	var buf bytes.Buffer
	fset := token.NewFileSet()
	for i := range b.Children {
		log_old.Log.PanicIfError(format.Node(&buf, fset, b.Node[i]))
		fmt.Fprintln(&buf)
	}
	return buf.String()
}

func (b *BlockStmt) Fix() Expr {
	return &BlockStmt{Children: fixExprs(b.Children), Node: b.Node}
}
