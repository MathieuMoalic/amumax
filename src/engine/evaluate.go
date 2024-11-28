package engine

import (
	"errors"
	"fmt"
	"go/ast"
	"go/token"
	"strconv"
	"strings"
)

func (p *ScriptParser) execute() error {
	scriptLines := strings.Split(p.e.script, "\n")
	var executeStatements func([]Statement, int) error
	executeStatements = func(statements []Statement, indentLevel int) error {
		indent := strings.Repeat("    ", indentLevel) // Indentation for nested blocks
		for _, stmt := range statements {
			// Log the line or block of code being executed
			if stmt.LineNum >= 0 && stmt.LineNum < len(scriptLines) {
				line := scriptLines[stmt.LineNum-3]
				p.e.log.Command(fmt.Sprintf("%s%s", indent, line))
			}

			// Execute the statement
			switch stmt.Type {
			case "assignment", "declaration":
				value, err := p.evaluateExpression(stmt.Expr)
				if err != nil {
					return fmt.Errorf("error evaluating expression: %v", err)
				}
				p.e.world.registerUserVariable(stmt.Name, value)
			case "function_call":
				fn, ok := p.e.world.getFunction(stmt.Name)
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
				p.e.log.Command(indent)                            // Open the loop block
				err := executeStatements(stmt.Body, indentLevel+1) // Process loop body with increased indent
				if err != nil {
					return err
				}
				p.e.log.Command(fmt.Sprintf("%s}", indent)) // Close the loop block
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
		val, ok := p.e.world.getVariable(e.Name)
		if ok {
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
		fn, ok := p.e.world.getFunction(funcName)
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
	case *ast.BinaryExpr:
		return p.evaluateBinaryExpr(e)
	case *ast.ParenExpr:
		return p.evaluateExpression(e.X)
	case *ast.IndexExpr:
		return nil, fmt.Errorf("unsupported index expression: %v", e)
	case *ast.SliceExpr:
		return nil, fmt.Errorf("unsupported slice expression: %v", e)
	case *ast.SelectorExpr:
		return nil, fmt.Errorf("unsupported selector expression: %v", e)
	case *ast.StarExpr:
		return nil, fmt.Errorf("unsupported star expression: %v", e)
	case *ast.TypeAssertExpr:
		return nil, fmt.Errorf("unsupported type assertion expression: %v", e)
	case *ast.ArrayType:
		return nil, fmt.Errorf("unsupported array type: %v", e)
	case *ast.MapType:
		return nil, fmt.Errorf("unsupported map type: %v", e)
	case *ast.StructType:
		return nil, fmt.Errorf("unsupported struct type: %v", e)
	case *ast.FuncLit:
		return nil, fmt.Errorf("unsupported function literal: %v", e)
	case *ast.FuncType:
		return nil, fmt.Errorf("unsupported function type: %v", e)
	case *ast.InterfaceType:
		return nil, fmt.Errorf("unsupported interface type: %v", e)
	case *ast.ChanType:
		return nil, fmt.Errorf("unsupported channel type: %v", e)
	case *ast.Ellipsis:
		return nil, fmt.Errorf("unsupported ellipsis: %v", e)
	case *ast.KeyValueExpr:
		return nil, fmt.Errorf("unsupported key-value expression: %v", e)
	case *ast.BadExpr:
		return nil, fmt.Errorf("unsupported bad expression: %v", e)
	default:
		return nil, fmt.Errorf("unsupported expression type: %T", e)
	}
}
func (p *ScriptParser) evaluateBinaryExpr(e *ast.BinaryExpr) (interface{}, error) {
	// Evaluate left operand
	leftVal, err := p.evaluateExpression(e.X)
	if err != nil {
		return nil, err
	}

	// Evaluate right operand
	rightVal, err := p.evaluateExpression(e.Y)
	if err != nil {
		return nil, err
	}

	// Perform the operation based on the operator
	switch e.Op {
	case token.ADD: // '+'
		return p.addValues(leftVal, rightVal)
	case token.SUB: // '-'
		return p.subtractValues(leftVal, rightVal)
	case token.MUL: // '*'
		return p.multiplyValues(leftVal, rightVal)
	case token.QUO: // '/'
		return p.divideValues(leftVal, rightVal)
	case token.REM: // '%'
		return p.remainderValues(leftVal, rightVal)
	case token.EQL, token.NEQ, token.LSS, token.LEQ, token.GTR, token.GEQ: // Comparison operators
		return p.compareValues(leftVal, rightVal, e.Op)
	case token.LAND: // '&&'
		return p.logicalAndValues(leftVal, rightVal)
	case token.LOR: // '||'
		return p.logicalOrValues(leftVal, rightVal)
	default:
		return nil, fmt.Errorf("unsupported binary operator: %s", e.Op)
	}
}
func (p *ScriptParser) addValues(left, right interface{}) (interface{}, error) {
	switch leftVal := left.(type) {
	case int:
		if rightVal, ok := right.(int); ok {
			return leftVal + rightVal, nil
		}
		if rightVal, ok := right.(float64); ok {
			return float64(leftVal) + rightVal, nil
		}
	case float64:
		if rightVal, ok := right.(int); ok {
			return leftVal + float64(rightVal), nil
		}
		if rightVal, ok := right.(float64); ok {
			return leftVal + rightVal, nil
		}
	case string:
		if rightVal, ok := right.(string); ok {
			return leftVal + rightVal, nil
		}
	}

	return nil, fmt.Errorf("unsupported addition between %T and %T", left, right)
}
func (p *ScriptParser) subtractValues(left, right interface{}) (interface{}, error) {
	switch leftVal := left.(type) {
	case int:
		if rightVal, ok := right.(int); ok {
			return leftVal - rightVal, nil
		}
		if rightVal, ok := right.(float64); ok {
			return float64(leftVal) - rightVal, nil
		}
	case float64:
		if rightVal, ok := right.(int); ok {
			return leftVal - float64(rightVal), nil
		}
		if rightVal, ok := right.(float64); ok {
			return leftVal - rightVal, nil
		}
	}

	return nil, fmt.Errorf("unsupported subtraction between %T and %T", left, right)
}
func (p *ScriptParser) multiplyValues(left, right interface{}) (interface{}, error) {
	switch leftVal := left.(type) {
	case int:
		if rightVal, ok := right.(int); ok {
			return leftVal * rightVal, nil
		}
		if rightVal, ok := right.(float64); ok {
			return float64(leftVal) * rightVal, nil
		}
	case float64:
		if rightVal, ok := right.(int); ok {
			return leftVal * float64(rightVal), nil
		}
		if rightVal, ok := right.(float64); ok {
			return leftVal * rightVal, nil
		}
	}

	return nil, fmt.Errorf("unsupported multiplication between %T and %T", left, right)
}
func (p *ScriptParser) divideValues(left, right interface{}) (interface{}, error) {
	switch leftVal := left.(type) {
	case int:
		if rightVal, ok := right.(int); ok {
			if rightVal == 0 {
				return nil, fmt.Errorf("division by zero")
			}
			return float64(leftVal) / float64(rightVal), nil // Return float64 for division
		}
		if rightVal, ok := right.(float64); ok {
			if rightVal == 0 {
				return nil, fmt.Errorf("division by zero")
			}
			return float64(leftVal) / rightVal, nil
		}
	case float64:
		if rightVal, ok := right.(int); ok {
			if rightVal == 0 {
				return nil, fmt.Errorf("division by zero")
			}
			return leftVal / float64(rightVal), nil
		}
		if rightVal, ok := right.(float64); ok {
			if rightVal == 0 {
				return nil, fmt.Errorf("division by zero")
			}
			return leftVal / rightVal, nil
		}
	}

	return nil, fmt.Errorf("unsupported division between %T and %T", left, right)
}
func (p *ScriptParser) remainderValues(left, right interface{}) (interface{}, error) {
	leftInt, ok1 := left.(int)
	rightInt, ok2 := right.(int)
	if ok1 && ok2 {
		if rightInt == 0 {
			return nil, fmt.Errorf("division by zero")
		}
		return leftInt % rightInt, nil
	}
	return nil, fmt.Errorf("unsupported remainder operation between %T and %T", left, right)
}
func (p *ScriptParser) compareValues(left, right interface{}, op token.Token) (interface{}, error) {
	switch leftVal := left.(type) {
	case int:
		var leftFloat float64 = float64(leftVal)
		var rightFloat float64
		switch rightVal := right.(type) {
		case int:
			rightFloat = float64(rightVal)
		case float64:
			rightFloat = rightVal
		default:
			return nil, fmt.Errorf("unsupported comparison between %T and %T", left, right)
		}
		return compareFloats(leftFloat, rightFloat, op), nil
	case float64:
		var rightFloat float64
		switch rightVal := right.(type) {
		case int:
			rightFloat = float64(rightVal)
		case float64:
			rightFloat = rightVal
		default:
			return nil, fmt.Errorf("unsupported comparison between %T and %T", left, right)
		}
		return compareFloats(leftVal, rightFloat, op), nil
	case string:
		rightVal, ok := right.(string)
		if !ok {
			return nil, fmt.Errorf("unsupported comparison between %T and %T", left, right)
		}
		return compareStrings(leftVal, rightVal, op), nil
	default:
		return nil, fmt.Errorf("unsupported comparison between %T and %T", left, right)
	}
}
func compareFloats(left, right float64, op token.Token) bool {
	switch op {
	case token.EQL:
		return left == right
	case token.NEQ:
		return left != right
	case token.LSS:
		return left < right
	case token.LEQ:
		return left <= right
	case token.GTR:
		return left > right
	case token.GEQ:
		return left >= right
	default:
		return false
	}
}
func compareStrings(left, right string, op token.Token) bool {
	switch op {
	case token.EQL:
		return left == right
	case token.NEQ:
		return left != right
	case token.LSS:
		return left < right
	case token.LEQ:
		return left <= right
	case token.GTR:
		return left > right
	case token.GEQ:
		return left >= right
	default:
		return false
	}
}
func (p *ScriptParser) logicalAndValues(left, right interface{}) (interface{}, error) {
	leftBool, ok1 := left.(bool)
	rightBool, ok2 := right.(bool)
	if ok1 && ok2 {
		return leftBool && rightBool, nil
	}
	return nil, fmt.Errorf("unsupported logical AND between %T and %T", left, right)
}

func (p *ScriptParser) logicalOrValues(left, right interface{}) (interface{}, error) {
	leftBool, ok1 := left.(bool)
	rightBool, ok2 := right.(bool)
	if ok1 && ok2 {
		return leftBool || rightBool, nil
	}
	return nil, fmt.Errorf("unsupported logical OR between %T and %T", left, right)
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
