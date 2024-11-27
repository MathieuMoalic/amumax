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
	return w
}
func (w *World) register() {
	w.registerQuantities()
	w.registerTableFunctions()
	w.registerSaveFunctions()
	w.registerMeshVariables()

}

func (w *World) registerQuantities() {
	w.RegisterVariable("geom", w.EngineState.geometry)
}

func (w *World) registerTableFunctions() {
	w.RegisterFunction("TableAutoSave", w.EngineState.table.TableAutoSave)
	w.RegisterFunction("TableAdd", w.EngineState.table.tableAdd)
	w.RegisterFunction("TableAddAs", w.EngineState.table.tableAddAs)
	w.RegisterFunction("TableSave", w.EngineState.table.tableSave)
}

func (w *World) registerSaveFunctions() {
	w.RegisterFunction("save", w.EngineState.savedQuantities.save)
	w.RegisterFunction("saveAs", w.EngineState.savedQuantities.saveAs)
	w.RegisterFunction("SaveAsChunks", w.EngineState.savedQuantities.saveAsChunk)
	w.RegisterFunction("AutoSave", w.EngineState.savedQuantities.autoSave)
	w.RegisterFunction("AutoSaveAs", w.EngineState.savedQuantities.autoSaveAs)
	w.RegisterFunction("AutoSaveAsChunk", w.EngineState.savedQuantities.autoSaveAsChunk)
	w.RegisterFunction("Chunks", mx3chunks)
}

func (w *World) registerMeshVariables() {
	w.RegisterVariable("Nx", &w.EngineState.mesh.Nx)
	w.RegisterVariable("Ny", &w.EngineState.mesh.Ny)
	w.RegisterVariable("Nz", &w.EngineState.mesh.Nz)
	w.RegisterVariable("dx", &w.EngineState.mesh.Dx)
	w.RegisterVariable("dy", &w.EngineState.mesh.Dy)
	w.RegisterVariable("dz", &w.EngineState.mesh.Dz)
	w.RegisterVariable("Tx", &w.EngineState.mesh.Tx)
	w.RegisterVariable("Ty", &w.EngineState.mesh.Ty)
	w.RegisterVariable("Tz", &w.EngineState.mesh.Tz)
}

// RegisterFunction registers a pre-defined function in the world.
func (w *World) RegisterFunction(name string, function interface{}) {
	w.Functions[name] = w.WrapFunction(function, name)
}

// RegisterVariable registers a pre-defined variable in the world.
func (w *World) RegisterVariable(name string, value interface{}) {
	if value == nil {
		w.EngineState.log.ErrAndExit("Value is nil for variable: %s", name)
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
			w.EngineState.log.Warn("Unsupported type: %T", ptr)
		}
	} else {
		w.Variables[name] = value
	}
	if w.isMeshExpression(name) {
		if w.EngineState.mesh.ReadyToCreate() {
			w.EngineState.mesh.Create()
			w.EngineState.magnetization.InitializeBuffer()
			w.EngineState.regions.InitializeBuffer()
			w.EngineState.metadata.AddMesh(w.EngineState.mesh)
		}
	}
	w.EngineState.metadata.Add(name, value)
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

		numIn := fnType.NumIn()
		isVariadic := fnType.IsVariadic()
		numFixedArgs := numIn
		if isVariadic {
			numFixedArgs-- // The last parameter is variadic
		}

		// Check if the number of arguments is sufficient
		if (!isVariadic && len(args) != numIn) || (isVariadic && len(args) < numFixedArgs) {
			expectedArgs := numIn
			if isVariadic {
				expectedArgs = numFixedArgs
				return nil, fmt.Errorf(
					"%s expects at least %d arguments (%s), got %d",
					name,
					expectedArgs,
					w.formatFunctionSignature(fnType, name),
					len(args),
				)
			} else {
				return nil, fmt.Errorf(
					"%s expects %d arguments (%s), got %d",
					name,
					expectedArgs,
					w.formatFunctionSignature(fnType, name),
					len(args),
				)
			}
		}

		// Prepare arguments for the function call
		in := make([]reflect.Value, numFixedArgs)

		// Handle fixed arguments
		for i := 0; i < numFixedArgs; i++ {
			expectedType := fnType.In(i)
			if len(args) <= i {
				return nil, fmt.Errorf(
					"%s: missing argument for parameter %d\nExpected function signature: %s",
					name,
					i+1,
					w.formatFunctionSignature(fnType, name),
				)
			}
			arg := args[i]
			argVal := reflect.ValueOf(arg)

			// Check if the argument is assignable to the expected type
			if !argVal.Type().AssignableTo(expectedType) {
				if expectedType.Kind() == reflect.Interface && argVal.Type().Implements(expectedType) {
					// The argument implements the expected interface; proceed without conversion
				} else if argVal.Type().ConvertibleTo(expectedType) {
					argVal = argVal.Convert(expectedType)
				} else {
					return nil, fmt.Errorf(
						"%s: argument %d (%v) is not assignable to %s\nExpected function signature: %s",
						name,
						i+1,
						argVal.Type(),
						expectedType,
						w.formatFunctionSignature(fnType, name),
					)
				}
			}

			in[i] = argVal
		}

		// Handle variadic arguments
		if isVariadic {
			variadicType := fnType.In(numIn - 1).Elem() // Element type of variadic parameter
			numVariadicArgs := len(args) - numFixedArgs
			variadicSlice := reflect.MakeSlice(reflect.SliceOf(variadicType), numVariadicArgs, numVariadicArgs)
			for i := 0; i < numVariadicArgs; i++ {
				arg := args[numFixedArgs+i]
				argVal := reflect.ValueOf(arg)

				// Check if the argument is assignable to the variadic type
				if !argVal.Type().AssignableTo(variadicType) {
					if variadicType.Kind() == reflect.Interface && argVal.Type().Implements(variadicType) {
						// The argument implements the expected interface; proceed without conversion
					} else if argVal.Type().ConvertibleTo(variadicType) {
						argVal = argVal.Convert(variadicType)
					} else {
						return nil, fmt.Errorf(
							"%s: argument %d (%v) is not assignable to %s\nExpected function signature: %s",
							name,
							numFixedArgs+i+1,
							argVal.Type(),
							variadicType,
							w.formatFunctionSignature(fnType, name),
						)
					}
				}

				variadicSlice.Index(i).Set(argVal)
			}
			// Append the variadic slice to the arguments
			in = append(in, variadicSlice)
		}

		var out []reflect.Value
		// Call the function using reflection
		if isVariadic {
			out = fnValue.CallSlice(in)
		} else {
			out = fnValue.Call(in)
		}

		// Handle the function's return values
		switch len(out) {
		case 0:
			return nil, nil
		case 1:
			if fnType.Out(0).Name() == "error" {
				if !out[0].IsNil() {
					return nil, out[0].Interface().(error)
				}
				return nil, nil
			}
			return out[0].Interface(), nil
		case 2:
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
	numIn := fnType.NumIn()
	isVariadic := fnType.IsVariadic()
	for i := 0; i < numIn; i++ {
		if i > 0 {
			sb.WriteString(", ")
		}
		inType := fnType.In(i)
		if isVariadic && i == numIn-1 {
			sb.WriteString("...")
			sb.WriteString(inType.Elem().String())
		} else {
			sb.WriteString(inType.String())
		}
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
