package engine_old

// declare functionality for interpreted input scripts

import (
	"reflect"
	"strings"

	"github.com/MathieuMoalic/amumax/src/engine_old/cuda_old"
	"github.com/MathieuMoalic/amumax/src/engine_old/fsutil_old"
	"github.com/MathieuMoalic/amumax/src/engine_old/log_old"
	"github.com/MathieuMoalic/amumax/src/engine_old/mesh_old"
	"github.com/MathieuMoalic/amumax/src/engine_old/script_old"
)

var Mesh mesh_old.Mesh

func init() {
	// Mesh = data.MeshType{Nx: 0, Ny: 0, Nz: 0, Dx: 0, Dy: 0, Dz: 0, Tx: 0, Ty: 0, Tz: 0, PBCx: 0, PBCy: 0, PBCz: 0, AutoMeshx: false, AutoMeshy: false, AutoMeshz: false}
}
func GetMesh() *mesh_old.Mesh {
	return &Mesh
}

func CompileFile(fname string) (*script_old.BlockStmt, error) {
	bytes, err := fsutil_old.Read(fname)
	if err != nil {
		return nil, err
	}
	script_old.AddMetadata = EngineState.Metadata.Add
	return World.Compile(string(bytes))
}

func EvalTryRecover(code string) {
	defer func() {
		if err := recover(); err != nil {
			if userErr, ok := err.(UserErr); ok {
				log_old.Log.Err("%v", userErr)
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
		log_old.Log.Command(code)
		log_old.Log.Err("%v", err.Error())
		return
	}
	log_old.Log.Command(rmln(tree.Format()))
	tree.Eval()
}

// holds the script state (variables etc)
var World = script_old.NewWorld()

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
func EvalFile(code *script_old.BlockStmt) {
	for i := range code.Children {
		formatted := rmln(script_old.Format(code.Node[i]))
		log_old.Log.Command(formatted)
		exp := code.Children[i]
		exp.Eval()
		if isMeshExpression(exp) {
			if Mesh.ReadyToCreate() {
				setBusy(true)
				defer setBusy(false)
				CreateMesh()
				NormMag.Alloc()
				Regions.Alloc()
				EngineState.Metadata.Init(OD(), StartTime, cuda_old.GPUInfo_old)
				EngineState.Metadata.AddMesh(&Mesh)
			}
		}
	}
}

func isMeshExpression(exp script_old.Expr) bool {
	namesToCheck := []string{"Nx", "Ny", "Nz", "Dx", "Dy", "Dz", "Tx", "Ty", "Tz"}
	val := reflect.ValueOf(exp)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if nameField := val.FieldByName("name"); nameField.IsValid() {
		name := strings.ToLower(nameField.String())
		return containsIgnoreCase(namesToCheck, name)
	}
	return false
}

func containsIgnoreCase(s []string, target string) bool {
	target = strings.ToLower(target)
	for _, v := range s {
		if strings.ToLower(v) == target {
			return true
		}
	}
	return false
}

var QuantityChanged = make(map[string]bool)

// wraps LValue and provides empty Child()
type lValueWrapper struct {
	lValue
	name string
}

func newLValueWrapper(name string, lv lValue) script_old.LValue {
	return &lValueWrapper{name: name, lValue: lv}
}
func (w *lValueWrapper) SetValue(val interface{}) {
	w.lValue.SetValue(val)
	QuantityChanged[w.name] = true
}

func (w *lValueWrapper) Child() []script_old.Expr { return nil }
func (w *lValueWrapper) Fix() script_old.Expr     { return script_old.NewConst(w) }

func (w *lValueWrapper) InputType() reflect.Type {
	if i, ok := w.lValue.(interface {
		InputType() reflect.Type
	}); ok {
		return i.InputType()
	} else {
		return w.Type()
	}
}
