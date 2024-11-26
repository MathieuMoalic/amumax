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
	w.RegisterFunction("TableAutoSave", w.EngineState.Table.TableAutoSave)

	w.RegisterVariable("Nx", &w.EngineState.Mesh.Nx)
	w.RegisterVariable("Ny", &w.EngineState.Mesh.Ny)
	w.RegisterVariable("Nz", &w.EngineState.Mesh.Nz)
	w.RegisterVariable("dx", &w.EngineState.Mesh.Dx)
	w.RegisterVariable("dy", &w.EngineState.Mesh.Dy)
	w.RegisterVariable("dz", &w.EngineState.Mesh.Dz)
	w.RegisterVariable("Tx", &w.EngineState.Mesh.Tx)
	w.RegisterVariable("Ty", &w.EngineState.Mesh.Ty)
	w.RegisterVariable("Tz", &w.EngineState.Mesh.Tz)
	return w
}

// RegisterFunction registers a pre-defined function in the world.
func (w *World) RegisterFunction(name string, function interface{}) {
	w.Functions[name] = w.WrapFunction(function, name)
}

// RegisterVariable registers a pre-defined variable in the world.
func (w *World) RegisterVariable(name string, value interface{}) {
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
			w.EngineState.Mesh.CreateMesh()
			w.EngineState.NormMag.alloc()
			// engine.Regions.Alloc()
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

			// Attempt to convert the argument to the expected type
			if !argValue.Type().AssignableTo(expectedType) {
				if argValue.Type().ConvertibleTo(expectedType) {
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
