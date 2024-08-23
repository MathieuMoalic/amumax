package engine

// declare functionality for interpreted input scripts

import (
	"reflect"

	"github.com/MathieuMoalic/amumax/httpfs"
	"github.com/MathieuMoalic/amumax/script"
	"github.com/MathieuMoalic/amumax/util"
)

func CompileFile(fname string) (*script.BlockStmt, error) {
	bytes, err := httpfs.Read(fname)
	if err != nil {
		return nil, err
	}
	return World.Compile(string(bytes))
}

func EvalTryRecover(code string) {
	defer func() {
		if err := recover(); err != nil {
			if userErr, ok := err.(UserErr); ok {
				util.Log.Err("%v", userErr)
			} else {
				panic(err)
			}
		}
	}()
	Eval(code)
}

func Eval(code string) {
	tree, err := World.Compile(code)
	if err != nil {
		util.Log.Command(code)
		util.Log.Err("%v", err.Error())
		return
	}
	util.Log.Command(rmln(tree.Format()))
	tree.Eval()
}

func Eval1Line(code string) interface{} {
	tree, err := World.Compile(code)
	if err != nil {
		util.Log.Err("%v", err.Error())
		return nil
	}
	if len(tree.Children) != 1 {
		util.Log.Err("expected single statement:%v", code)
		return nil
	}
	return tree.Children[0].Eval()
}

// holds the script state (variables etc)
var World = script.NewWorld()

// Add a function to the script world
func DeclFunc(name string, f interface{}, doc string) {
	World.Func(name, f, doc)
}

// Add a constant to the script world
func DeclConst(name string, value float64, doc string) {
	World.Const(name, value, doc)
}

// Add a read-only variable to the script world.
// It can be changed, but not by the user.
func DeclROnly(name string, value interface{}, doc string) {
	World.ROnly(name, value, doc)
	AddQuantity(name, value, doc)
}

func Export(q interface {
	Name() string
	Unit() string
}, doc string) {
	DeclROnly(q.Name(), q, cat(doc, q.Unit()))
}

// Add a (pointer to) variable to the script world
func DeclVar(name string, value interface{}, doc string) {
	World.Var(name, value, doc)
	AddQuantity(name, value, doc)
}

// Hack for fixing the closure caveat:
// Defines "t", the time variable, handled specially by Fix()
func DeclTVar(name string, value interface{}, doc string) {
	World.TVar(name, value, doc)
	AddQuantity(name, value, doc)
}

// Add an LValue to the script world.
// Assign to LValue invokes SetValue()
func DeclLValue(name string, value LValue, doc string) {
	AddParameter(name, value, doc)
	World.LValue(name, newLValueWrapper(value), doc)
	AddQuantity(name, value, doc)
}

// LValue is settable
type LValue interface {
	SetValue(interface{}) // assigns a new value
	Eval() interface{}    // evaluate and return result (nil for void)
	Type() reflect.Type   // type that can be assigned and will be returned by Eval
}

// evaluate code, exit on error (behavior for input files)
func EvalFile(code *script.BlockStmt) {
	for i := range code.Children {
		formatted := rmln(script.Format(code.Node[i]))
		util.Log.Command(formatted)
		code.Children[i].Eval()
	}
}

// wraps LValue and provides empty Child()
type lValueWrapper struct {
	LValue
}

func newLValueWrapper(lv LValue) script.LValue {
	return &lValueWrapper{lv}
}

func (w *lValueWrapper) Child() []script.Expr { return nil }
func (w *lValueWrapper) Fix() script.Expr     { return script.NewConst(w) }

func (w *lValueWrapper) InputType() reflect.Type {
	if i, ok := w.LValue.(interface {
		InputType() reflect.Type
	}); ok {
		return i.InputType()
	} else {
		return w.Type()
	}
}
