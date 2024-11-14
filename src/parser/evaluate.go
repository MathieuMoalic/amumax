package parser

import (
	"errors"
	"fmt"
	"go/ast"
	"go/token"
	"strconv"
	"strings"
)

func (p *ScriptParser) Execute(backend *SimulationBackend) error {
	variables := make(map[string]interface{})
	for _, stmt := range p.statements {
		switch stmt.Type {
		case "assignment", "declaration":
			// Evaluate the RHS expression
			value, err := p.evaluateExpression(stmt.Expr, backend, variables)
			if err != nil {
				return fmt.Errorf("error evaluating expression: %v", err)
			}
			// Store the value in the local variables map
			variables[stmt.Name] = value
			// Also set the parameter in the backend
			backend.SetParameter(stmt.Name, value)
		case "function_call":
			fn, ok := functionRegistry[stmt.Name]
			if !ok {
				return fmt.Errorf("unsupported function: %s", stmt.Name)
			}
			args := []interface{}{}
			for _, argExpr := range stmt.ArgExprs {
				argValue, err := p.evaluateExpression(argExpr, backend, variables)
				if err != nil {
					return fmt.Errorf("error evaluating argument in function '%s': %v", stmt.Name, err)
				}
				args = append(args, argValue)
			}
			_, err := fn(backend, args)
			if err != nil {
				return fmt.Errorf("error executing function %s: %v", stmt.Name, err)
			}
		}
	}
	return nil
}

func (p *ScriptParser) evaluateExpression(expr ast.Expr, backend *SimulationBackend, variables map[string]interface{}) (interface{}, error) {
	switch e := expr.(type) {
	case *ast.BasicLit:
		return parseValue(e.Value)
	case *ast.Ident:
		if val, ok := variables[e.Name]; ok {
			return val, nil
		}
		return nil, fmt.Errorf("undefined variable: %s", e.Name)
	case *ast.CallExpr:
		funcName := p.formatExpr(e.Fun)
		args := []interface{}{}
		for _, argExpr := range e.Args {
			argValue, err := p.evaluateExpression(argExpr, backend, variables)
			if err != nil {
				return nil, err
			}
			args = append(args, argValue)
		}
		fn, ok := functionRegistry[funcName]
		if !ok {
			return nil, fmt.Errorf("unsupported function: %s", funcName)
		}
		return fn(backend, args)
	case *ast.UnaryExpr: // Handle unary expressions, e.g., -1.5
		val, err := p.evaluateExpression(e.X, backend, variables)
		if err != nil {
			return nil, err
		}
		switch e.Op {
		case token.SUB: // For negative numbers
			if floatVal, ok := val.(float64); ok {
				return -floatVal, nil
			}
			if intVal, ok := val.(int); ok {
				return -intVal, nil
			}
		}
		return nil, fmt.Errorf("unsupported unary operation: %s", e.Op)
	case *ast.CompositeLit: // Handle array literals
		return p.parseArrayLiteral(e)
	default:
		return nil, fmt.Errorf("unsupported expression type: %T", e)
	}
}

func (p *ScriptParser) parseArrayLiteral(lit *ast.CompositeLit) (interface{}, error) {
	elementType := p.formatExpr(lit.Type)
	var elements []interface{}

	for _, elt := range lit.Elts {
		val, err := p.evaluateExpression(elt, nil, nil)
		if err != nil {
			return nil, fmt.Errorf("error parsing array element: %v", err)
		}
		elements = append(elements, val)
	}

	// Convert elements to appropriate typed slices
	switch elementType {
	case "[]float64":
		floatArray := make([]float64, len(elements))
		for i, elem := range elements {
			if floatVal, ok := elem.(float64); ok {
				floatArray[i] = floatVal
			} else {
				return nil, fmt.Errorf("element %v is not of type float64", elem)
			}
		}
		return floatArray, nil
	case "[]int":
		intArray := make([]int, len(elements))
		for i, elem := range elements {
			if intVal, ok := elem.(int); ok {
				intArray[i] = intVal
			} else {
				return nil, fmt.Errorf("element %v is not of type int", elem)
			}
		}
		return intArray, nil
	case "[]string":
		strArray := make([]string, len(elements))
		for i, elem := range elements {
			if strVal, ok := elem.(string); ok {
				strArray[i] = strVal
			} else {
				return nil, fmt.Errorf("element %v is not of type string", elem)
			}
		}
		return strArray, nil
	default:
		return nil, fmt.Errorf("unsupported array element type: %s", elementType)
	}
}

// parseValue parses a string into the appropriate type: int, float64, or string.
func parseValue(value string) (interface{}, error) {
	// Attempt to parse as int
	if intVal, err := strconv.Atoi(value); err == nil {
		return intVal, nil
	}

	// Attempt to parse as float64
	if floatVal, err := strconv.ParseFloat(value, 64); err == nil {
		return floatVal, nil
	}

	// Check for string literal
	if strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"") {
		return value[1 : len(value)-1], nil // Remove quotes
	}

	return nil, errors.New("unsupported data type for value")
}
