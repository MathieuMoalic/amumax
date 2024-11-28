package new_engine

import (
	"fmt"
	"reflect"
	"strings"
)

type world struct {
	e         *engineState
	functions map[string]interface{}
	variables map[string]interface{}
}

func newWorld(engineState *engineState) *world {
	w := &world{
		e:         engineState,
		functions: make(map[string]interface{}),
		variables: make(map[string]interface{}),
	}
	return w
}

func (w *world) register() {
	w.registerQuantities()
	w.registerTableFunctions()
	w.registerMesh()
}

func (w *world) registerQuantities() {
	w.registerVariable("geom", w.e.geometry)
}

func (w *world) registerTableFunctions() {
	w.registerFunction("TableAutoSave", w.e.table.tableAutoSave)
	w.registerFunction("TableAdd", w.e.table.tableAdd)
	w.registerFunction("TableAddAs", w.e.table.tableAddAs)
	w.registerFunction("TableSave", w.e.table.tableSave)
}

func (w *world) registerMesh() {
	w.registerVariable("Nx", &w.e.mesh.Nx)
	w.registerVariable("Ny", &w.e.mesh.Ny)
	w.registerVariable("Nz", &w.e.mesh.Nz)
	w.registerVariable("dx", &w.e.mesh.Dx)
	w.registerVariable("dy", &w.e.mesh.Dy)
	w.registerVariable("dz", &w.e.mesh.Dz)
	w.registerVariable("Tx", &w.e.mesh.Tx)
	w.registerVariable("Ty", &w.e.mesh.Ty)
	w.registerVariable("Tz", &w.e.mesh.Tz)
	w.registerFunction("SetGridSize", w.e.mesh.SetGridSize)
	w.registerFunction("SetCellSize", w.e.mesh.SetCellSize)

}

// registerFunction registers a pre-defined function in the world.
func (w *world) registerFunction(name string, function interface{}) {
	w.functions[strings.ToLower(name)] = w.wrapFunction(function, strings.ToLower(name))
}

// registerVariable registers a pre-defined variable in the world.
func (w *world) registerVariable(name string, value interface{}) {
	if value == nil {
		w.e.log.ErrAndExit("Value is nil for variable: %s", name)
	}
	w.variables[strings.ToLower(name)] = value
}

// registerUserVariable registers a user-defined variable in the world.
func (w *world) registerUserVariable(name string, value interface{}) {
	if existingValue, ok := w.variables[name]; ok {
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
			w.e.log.Warn("Unsupported type: %T", ptr)
		}
	} else {
		w.variables[strings.ToLower(name)] = value
	}
	if w.isMeshExpression(name) {
		if w.e.mesh.ReadyToCreate() {
			w.e.mesh.Create()
			w.e.magnetization.initializeBuffer()
			w.e.regions.initializeBuffer()
			w.e.metadata.AddMesh(w.e.mesh)
		}
	}
	w.e.metadata.Add(name, value)
	w.variables[name] = value
}

// wrapFunction creates a universal wrapper for any function.
func (w *world) wrapFunction(fn interface{}, name string) func([]interface{}) (interface{}, error) {
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
func (w *world) formatFunctionSignature(fnType reflect.Type, name string) string {
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

func (w *world) isMeshExpression(name string) bool {
	namesToCheck := []string{"Nx", "Ny", "Nz", "Dx", "Dy", "Dz", "Tx", "Ty", "Tz"}
	for _, v := range namesToCheck {
		if strings.EqualFold(v, name) {
			return true
		}
	}
	return false
}

func (w *world) getVariable(name string) (interface{}, bool) {
	value, ok := w.variables[strings.ToLower(name)]
	return value, ok
}

func (w *world) getFunction(name string) (interface{}, bool) {
	value, ok := w.functions[strings.ToLower(name)]
	return value, ok
}
