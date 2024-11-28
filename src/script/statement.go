package script

import "go/ast"

// Statement represents a parsed code statement, with relevant metadata and attributes.
type Statement struct {
	Type     string      // Type of the statement, e.g., "declaration", "for_loop", etc.
	Name     string      // Variable or function name
	Value    string      // Assigned or returned value
	Args     []string    // Arguments for function calls
	LineNum  int         // Line number in the script
	Original ast.Node    // Original AST node (not used in comparisons)
	Init     string      // For-loop initialization
	Cond     string      // Condition expression (for-loops, if-statements)
	Post     string      // Post expression in for-loops
	Body     []Statement // Body statements for loops or if-statements
	Expr     ast.Expr    // RHS expression for assignments
	ArgExprs []ast.Expr  // Argument expressions for function calls
}
