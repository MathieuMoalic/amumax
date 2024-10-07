package engine

// declare functionality for interpreted input scripts

import (
	"reflect"

	"github.com/MathieuMoalic/amumax/src/httpfs"
	"github.com/MathieuMoalic/amumax/src/log"
	"github.com/MathieuMoalic/amumax/src/script"
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
				log.Log.Err("%v", userErr)
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
		log.Log.Command(code)
		log.Log.Err("%v", err.Error())
		return
	}
	log.Log.Command(rmln(tree.Format()))
	tree.Eval()
}

// holds the script state (variables etc)
var World = script.NewWorld()

func export(q interface {
	Name() string
	Unit() string
}, doc string) {
	declROnly(q.Name(), q, cat(doc, q.Unit()))
}

// lValue is settable
type lValue interface {
	SetValue(interface{}) // assigns a new value
	Eval() interface{}    // evaluate and return result (nil for void)
	Type() reflect.Type   // type that can be assigned and will be returned by Eval
}

// evaluate code, exit on error (behavior for input files)
func EvalFile(code *script.BlockStmt) {
	for i := range code.Children {
		formatted := rmln(script.Format(code.Node[i]))
		log.Log.Command(formatted)
		code.Children[i].Eval()
	}
}

var QuantityChanged = make(map[string]bool)

// wraps LValue and provides empty Child()
type lValueWrapper struct {
	lValue
	name string
}

func newLValueWrapper(name string, lv lValue) script.LValue {
	return &lValueWrapper{name: name, lValue: lv}
}
func (w *lValueWrapper) SetValue(val interface{}) {
	w.lValue.SetValue(val)
	QuantityChanged[w.name] = true
}

func (w *lValueWrapper) Child() []script.Expr { return nil }
func (w *lValueWrapper) Fix() script.Expr     { return script.NewConst(w) }

func (w *lValueWrapper) InputType() reflect.Type {
	if i, ok := w.lValue.(interface {
		InputType() reflect.Type
	}); ok {
		return i.InputType()
	} else {
		return w.Type()
	}
}
