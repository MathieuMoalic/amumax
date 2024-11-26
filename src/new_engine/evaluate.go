package new_engine

import (
	"errors"
	"fmt"
	"go/ast"
	"go/token"
	"strconv"
	"strings"
)

func (p *ScriptParser) Execute() error {
	scriptLines := strings.Split(p.EngineState.Script, "\n")
	var executeStatements func([]Statement, int) error
	executeStatements = func(statements []Statement, indentLevel int) error {
		indent := strings.Repeat("    ", indentLevel) // Indentation for nested blocks
		for _, stmt := range statements {
			// Log the line or block of code being executed
			if stmt.LineNum >= 0 && stmt.LineNum < len(scriptLines) {
				line := scriptLines[stmt.LineNum-3]
				p.EngineState.Log.Command(fmt.Sprintf("%s%s", indent, line))
			}

			// Execute the statement
			switch stmt.Type {
			case "assignment", "declaration":
				value, err := p.evaluateExpression(stmt.Expr)
				if err != nil {
					return fmt.Errorf("error evaluating expression: %v", err)
				}
				p.EngineState.World.RegisterUserVariable(stmt.Name, value)
			case "function_call":
				fn, ok := p.EngineState.World.Functions[stmt.Name]
				if !ok {
					return fmt.Errorf("unsupported function: %s", stmt.Name)
				}
				args := []interface{}{}
				for _, argExpr := range stmt.ArgExprs {
					argValue, err := p.evaluateExpression(argExpr)
					if err != nil {
						return fmt.Errorf("error evaluating argument in function '%s': %v", stmt.Name, err)
					}
					args = append(args, argValue)
				}
				fnTyped, ok := fn.(func([]interface{}) (interface{}, error))
				if !ok {
					return fmt.Errorf("invalid function type for: %s", stmt.Name)
				}
				_, err := fnTyped(args)
				if err != nil {
					return fmt.Errorf("error executing function %s: %v", stmt.Name, err)
				}
			case "for_loop":
				p.EngineState.Log.Command(indent)                  // Open the loop block
				err := executeStatements(stmt.Body, indentLevel+1) // Process loop body with increased indent
				if err != nil {
					return err
				}
				p.EngineState.Log.Command(fmt.Sprintf("%s}", indent)) // Close the loop block
			}
		}
		return nil
	}

	return executeStatements(p.statements, 0)
}

func (p *ScriptParser) evaluateExpression(expr ast.Expr) (interface{}, error) {
	switch e := expr.(type) {
	case *ast.BasicLit:
		return parseValue(e.Value)
	case *ast.Ident:
		if val, ok := p.EngineState.World.Variables[e.Name]; ok {
			return val, nil
		}
		return nil, fmt.Errorf("undefined variable: %s", e.Name)
	case *ast.CallExpr:
		funcName := p.formatExpr(e.Fun)
		args := []interface{}{}
		for _, argExpr := range e.Args {
			argValue, err := p.evaluateExpression(argExpr)
			if err != nil {
				return nil, err
			}
			args = append(args, argValue)
		}
		fn, ok := p.EngineState.World.Functions[funcName]
		if !ok {
			return nil, fmt.Errorf("unsupported function: %s", funcName)
		}
		fnTyped, ok := fn.(func([]interface{}) (interface{}, error))
		if !ok {
			return nil, fmt.Errorf("invalid function type for: %s", funcName)
		}
		return fnTyped(args)
	case *ast.UnaryExpr: // Handle unary expressions, e.g., -1.5
		val, err := p.evaluateExpression(e.X)
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
		val, err := p.evaluateExpression(elt)
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
