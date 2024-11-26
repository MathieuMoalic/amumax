package new_engine

import (
	"fmt"
	"reflect"
	"strings"
)

type World struct {
	EngineState *EngineStateStruct
	Functions   map[string]interface{}
	Variables   map[string]interface{}
}

func NewWorld(engineState *EngineStateStruct) *World {
	w := &World{
		EngineState: engineState,
		Functions:   make(map[string]interface{}),
		Variables:   make(map[string]interface{}),
	}

	w.registerQuantities()
	w.registerTableFunctions()
	w.registerSaveFunctions()
	w.registerMeshVariables()
	w.registerShapeFunctions()
	return w
}

func (w *World) registerQuantities() {
	w.RegisterVariable("geom", w.EngineState.Geometry)
}

func (w *World) registerTableFunctions() {
	w.RegisterFunction("TableAutoSave", w.EngineState.Table.TableAutoSave)
	w.RegisterFunction("TableAdd", w.EngineState.Table.tableAdd)
	w.RegisterFunction("TableAddAs", w.EngineState.Table.tableAddAs)
	w.RegisterFunction("TableSave", w.EngineState.Table.tableSave)
}

func (w *World) registerSaveFunctions() {
	w.RegisterFunction("save", w.EngineState.SavedQuantities.save)
	w.RegisterFunction("saveAs", w.EngineState.SavedQuantities.saveAs)
	w.RegisterFunction("SaveAsChunks", w.EngineState.SavedQuantities.saveAsChunk)
	w.RegisterFunction("AutoSave", w.EngineState.SavedQuantities.autoSave)
	w.RegisterFunction("AutoSaveAs", w.EngineState.SavedQuantities.autoSaveAs)
	w.RegisterFunction("AutoSaveAsChunk", w.EngineState.SavedQuantities.autoSaveAsChunk)
	w.RegisterFunction("Chunks", mx3chunks)
}

func (w *World) registerShapeFunctions() {
	w.RegisterFunction("wave", wave)
	w.RegisterFunction("ellipsoid", ellipsoid)
	w.RegisterFunction("ellipse", ellipse)
	w.RegisterFunction("cone", cone)
	w.RegisterFunction("circle", circle)
	w.RegisterFunction("cylinder", cylinder)
	w.RegisterFunction("cuboid", cuboid)
	w.RegisterFunction("rect", rect)
	w.RegisterFunction("triangle", triangle)
	w.RegisterFunction("rTriangle", rTriangle)
	w.RegisterFunction("hexagon", hexagon)
	w.RegisterFunction("diamond", diamond)
	w.RegisterFunction("squircle", squircle)
	w.RegisterFunction("square", square)
	w.RegisterFunction("xRange", xRange)
	w.RegisterFunction("yRange", yRange)
	w.RegisterFunction("zRange", zRange)
	w.RegisterFunction("universe", universe)
}

func (w *World) registerMeshVariables() {
	w.RegisterVariable("Nx", &w.EngineState.Mesh.Nx)
	w.RegisterVariable("Ny", &w.EngineState.Mesh.Ny)
	w.RegisterVariable("Nz", &w.EngineState.Mesh.Nz)
	w.RegisterVariable("dx", &w.EngineState.Mesh.Dx)
	w.RegisterVariable("dy", &w.EngineState.Mesh.Dy)
	w.RegisterVariable("dz", &w.EngineState.Mesh.Dz)
	w.RegisterVariable("Tx", &w.EngineState.Mesh.Tx)
	w.RegisterVariable("Ty", &w.EngineState.Mesh.Ty)
	w.RegisterVariable("Tz", &w.EngineState.Mesh.Tz)
}

// RegisterFunction registers a pre-defined function in the world.
func (w *World) RegisterFunction(name string, function interface{}) {
	w.Functions[name] = w.WrapFunction(function, name)
}

// RegisterVariable registers a pre-defined variable in the world.
func (w *World) RegisterVariable(name string, value interface{}) {
	if value == nil {
		w.EngineState.Log.ErrAndExit("Value is nil for variable: %s", name)
	}
	w.Variables[name] = value
}

// RegisterUserVariable registers a user-defined variable in the world.
func (w *World) RegisterUserVariable(name string, value interface{}) {
	if existingValue, ok := w.Variables[name]; ok {
		switch ptr := existingValue.(type) {
		case *int:
			if v, ok := value.(int); ok {
				*ptr = v
			}
		case *float64:
			if v, ok := value.(float64); ok {
				*ptr = v
			}
		case *bool:
			if v, ok := value.(bool); ok {
				*ptr = v
			}
		default:
			w.EngineState.Log.Warn("Unsupported type: %T", ptr)
		}
	} else {
		w.Variables[name] = value
	}
	if w.isMeshExpression(name) {
		if w.EngineState.Mesh.ReadyToCreate() {
			w.EngineState.Mesh.Create()
			w.EngineState.Magnetization.InitializeBuffer()
			w.EngineState.Regions.InitializeBuffer()
			w.EngineState.Metadata.AddMesh(w.EngineState.Mesh)
		}
	}
	w.EngineState.Metadata.Add(name, value)
	w.Variables[name] = value
}

// WrapFunction creates a universal wrapper for any function.
func (w *World) WrapFunction(fn interface{}, name string) func([]interface{}) (interface{}, error) {
	return func(args []interface{}) (interface{}, error) {
		fnValue := reflect.ValueOf(fn)
		fnType := fnValue.Type()

		// Ensure the provided function is callable
		if fnType.Kind() != reflect.Func {
			return nil, fmt.Errorf("provided argument is not a function")
		}

		// Check if the number of arguments matches
		if len(args) != fnType.NumIn() {
			return nil, fmt.Errorf(
				"%s expects %d arguments (%s), got %d",
				name,
				fnType.NumIn(),
				w.formatFunctionSignature(fnType, name),
				len(args),
			)
		}

		// Prepare arguments for the function call
		in := make([]reflect.Value, len(args))
		for i, arg := range args {
			expectedType := fnType.In(i)

			argValue := reflect.ValueOf(arg)

			// Check if the argument is assignable to the expected type
			if !argValue.Type().AssignableTo(expectedType) {
				if expectedType.Kind() == reflect.Interface && argValue.Type().Implements(expectedType) {
					// The argument implements the expected interface; proceed without conversion
					// No action needed here
				} else if argValue.Type().ConvertibleTo(expectedType) {
					argValue = argValue.Convert(expectedType)
				} else {
					w.EngineState.Log.Err("%s is not assignable to %s", argValue.Type(), expectedType)
					w.EngineState.Log.Err("Expected function signature: %s", w.formatFunctionSignature(fnType, name))
					return nil, fmt.Errorf("wrong argument type")
				}
			}

			in[i] = argValue
		}

		// Call the function using reflection
		out := fnValue.Call(in)

		// Handle the function's return values
		switch len(out) {
		case 0:
			// Function does not return any values
			return nil, nil
		case 1:
			// Function returns a single value
			if fnType.Out(0).Name() == "error" {
				if !out[0].IsNil() {
					return nil, out[0].Interface().(error)
				}
				return nil, nil
			}
			return out[0].Interface(), nil
		case 2:
			// Function returns (interface{}, error)
			var err error
			if !out[1].IsNil() {
				err = out[1].Interface().(error)
			}
			return out[0].Interface(), err
		default:
			return nil, fmt.Errorf("%s has unsupported number of return values: %d", name, len(out))
		}
	}
}

// formatFunctionSignature returns a string representation of the function's signature.
func (w *World) formatFunctionSignature(fnType reflect.Type, name string) string {
	var sb strings.Builder
	sb.WriteString(name)
	sb.WriteString("(")
	for i := 0; i < fnType.NumIn(); i++ {
		if i > 0 {
			sb.WriteString(", ")
		}
		inType := fnType.In(i)
		sb.WriteString(inType.String())
	}
	sb.WriteString(")")
	if fnType.NumOut() > 0 {
		sb.WriteString(" (")
		for i := 0; i < fnType.NumOut(); i++ {
			if i > 0 {
				sb.WriteString(", ")
			}
			outType := fnType.Out(i)
			sb.WriteString(outType.String())
		}
		sb.WriteString(")")
	}
	return sb.String()
}

func (w *World) isMeshExpression(name string) bool {
	namesToCheck := []string{"Nx", "Ny", "Nz", "Dx", "Dy", "Dz", "Tx", "Ty", "Tz"}
	for _, v := range namesToCheck {
		if strings.EqualFold(v, name) {
			return true
		}
	}
	return false
}