package engine

// declare functionality for interpreted input scripts

import (
	"reflect"
	"strings"

	"github.com/MathieuMoalic/amumax/src/cuda"
	"github.com/MathieuMoalic/amumax/src/fsutil"
	"github.com/MathieuMoalic/amumax/src/log"
	"github.com/MathieuMoalic/amumax/src/mesh"
	"github.com/MathieuMoalic/amumax/src/script"
)

var Mesh mesh.Mesh

func init() {
	// Mesh = data.MeshType{Nx: 0, Ny: 0, Nz: 0, Dx: 0, Dy: 0, Dz: 0, Tx: 0, Ty: 0, Tz: 0, PBCx: 0, PBCy: 0, PBCz: 0, AutoMeshx: false, AutoMeshy: false, AutoMeshz: false}
}
func GetMesh() *mesh.Mesh {
	return &Mesh
}

func CompileFile(fname string) (*script.BlockStmt, error) {
	bytes, err := fsutil.Read(fname)
	if err != nil {
		return nil, err
	}
	script.AddMetadata = EngineState.Metadata.Add
	return World.Compile(string(bytes))
}

func EvalTryRecover(code string) {
	defer func() {
		if err := recover(); err != nil {
			log.Log.Err("%v", err)
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
	SetValue(any)       // assigns a new value
	Eval() any          // evaluate and return result (nil for void)
	Type() reflect.Type // type that can be assigned and will be returned by Eval
}

// evaluate code, exit on error (behavior for input files)
func EvalFile(code *script.BlockStmt) {
	for i := range code.Children {
		formatted := rmln(script.Format(code.Node[i]))
		log.Log.Command(formatted)
		exp := code.Children[i]
		exp.Eval()
		if isMeshExpression(exp) {
			if Mesh.ReadyToCreate() {
				setBusy(true)
				defer setBusy(false)
				CreateMesh()
				NormMag.Alloc()
				Regions.Alloc()
				EngineState.Metadata.Init(OD(), StartTime, cuda.GPUInfo_old)
				EngineState.Metadata.AddMesh(&Mesh)
			}
		}
	}
}

func isMeshExpression(exp script.Expr) bool {
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

func newLValueWrapper(name string, lv lValue) script.LValue {
	return &lValueWrapper{name: name, lValue: lv}
}
func (w *lValueWrapper) SetValue(val any) {
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
