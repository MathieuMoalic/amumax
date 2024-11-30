package script

import (
	"fmt"
	"strings"
)

func (p *ScriptParser) Execute() error {
	scriptLines := strings.Split(*p.script, "\n")
	var executeStatements func([]Statement, int) error
	executeStatements = func(statements []Statement, indentLevel int) error {
		indent := strings.Repeat("    ", indentLevel) // Indentation for nested blocks
		for _, stmt := range statements {
			// Log the line or block of code being executed
			if stmt.LineNum >= 0 && stmt.LineNum < len(scriptLines) {
				line := scriptLines[stmt.LineNum+p.lineOffset]
				p.log.Command(fmt.Sprintf("%s%s", indent, line))
			}

			// Execute the statement
			switch stmt.Type {
			case "assignment", "declaration":
				value, err := p.evaluateExpression(stmt.Expr)
				if err != nil {
					return fmt.Errorf("error evaluating expression: %v", err)
				}
				p.registerUserVariable(strings.ToLower(stmt.Name), value)
				p.metadata.Add(stmt.Name, value)
			case "function_call":
				p.log.Command(scriptLines[stmt.LineNum+p.lineOffset]) // Log the function call
				fn, ok := p.getFunction(stmt.Name)
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
				p.log.Command(indent)                              // Open the loop block
				err := executeStatements(stmt.Body, indentLevel+1) // Process loop body with increased indent
				if err != nil {
					return err
				}
				p.log.Command(fmt.Sprintf("%s}", indent)) // Close the loop block
			}
		}
		return nil
	}

	return executeStatements(p.statements, 0)
}
