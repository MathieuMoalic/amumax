package script

import (
	"fmt"
	"go/ast"
	"go/parser"
)

// CompileExpr Compiles an expression, which can then be evaluated. E.g.:
//
//	expr, err := world.CompileExpr("1+1")
//	expr.Eval()   // returns 2
func (w *World) CompileExpr(src string) (code Expr, e error) {
	// parse
	tree, err := parser.ParseExpr(src)
	if err != nil {
		return nil, fmt.Errorf(`parse "%s": %v`, src, err)
	}
	if Debug {
		err = ast.Print(nil, tree)
		if err != nil {
			return nil, fmt.Errorf(`print "%s": %v`, src, err)
		}
	}

	// catch compile errors
	if !Debug {
		defer func() {
			err := recover()
			if err == nil {
				return
			}
			if er, ok := err.(*compileErr); ok {
				code = nil
				e = fmt.Errorf(`parse "%s": %v`, src, er)
			} else {
				panic(err)
			}
		}()
	}
	return w.compile(tree), nil
}

// MustCompileExpr CompileExpr with panic on error.
func (w *World) MustCompileExpr(src string) Expr {
	code, err := w.CompileExpr(src)
	if err != nil {
		panic(err)
	}
	return code
}

// Compile compiles source consisting of a number of statements. E.g.:
//
//	src = "a = 1; b = sin(x)"
//	code, err := world.Compile(src)
//	code.Eval()
func (w *World) Compile(src string) (code *BlockStmt, e error) {
	// parse
	exprSrc := "func(){\n" + src + "\n}" // wrap in func to turn into expression
	tree, err := parser.ParseExpr(exprSrc)
	if err != nil {
		return nil, fmt.Errorf("script line %v: ", err)
	}

	// catch compile errors and decode line number
	if !Debug {
		defer func() {
			err := recover()
			if err == nil {
				return
			}
			if compErr, ok := err.(*compileErr); ok {
				code = nil
				e = fmt.Errorf("script %v: %v", pos2line(compErr.pos, exprSrc), compErr.msg)
			} else {
				panic(err)
			}
		}()
	}

	// compile
	stmts := tree.(*ast.FuncLit).Body.List // strip func again
	if Debug {
		err = ast.Print(nil, stmts)
		if err != nil {
			return nil, fmt.Errorf(`print "%s": %v`, src, err)
		}
	}
	block := new(BlockStmt)
	for _, s := range stmts {
		block.append(w.compile(s), s)
	}
	return block, nil
}

// MustCompile Like Compile but panics on error
func (w *World) MustCompile(src string) Expr {
	code, err := w.Compile(src)
	if err != nil {
		panic(err)
	}
	return code
}
