package new_engine

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

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

// ScriptParser parses and stores statements from a script.
type ScriptParser struct {
	e          *engineState
	statements []Statement
	fset       *token.FileSet
	lineOffset int // this is used to adjust line numbers for printing because the wrappedScript adds 3 lines
}

// newScriptParser initializes and returns a new ScriptParser instance.
func newScriptParser(es *engineState) *ScriptParser {
	return &ScriptParser{
		e:          es,
		fset:       token.NewFileSet(),
		lineOffset: -3,
	}
}

// Parse parses a script, wrapping it in a main function to process each statement.
func (p *ScriptParser) Parse(script string) error {
	wrappedScript := "package main\nfunc main() {\n" + script + "\n}"
	file, err := parser.ParseFile(p.fset, "", wrappedScript, parser.AllErrors)
	if err != nil {
		return fmt.Errorf("parsing error: %v", err)
	}

	// Extract and process the main function statements
	var mainFunc *ast.FuncDecl
	for _, decl := range file.Decls {
		if fn, ok := decl.(*ast.FuncDecl); ok && fn.Name.Name == "main" {
			mainFunc = fn
			break
		}
	}
	if mainFunc != nil && mainFunc.Body != nil {
		for _, stmt := range mainFunc.Body.List {
			p.processNode(stmt)
		}
	}
	return nil
}

// processNode routes each AST node type to the appropriate processing function.
func (p *ScriptParser) processNode(node ast.Node) {
	switch n := node.(type) {
	case *ast.AssignStmt:
		p.processAssignment(n)
	case *ast.ExprStmt:
		if call, ok := n.X.(*ast.CallExpr); ok {
			p.processFunctionCall(call)
		}
	case *ast.ForStmt:
		p.processForLoop(n)
	case *ast.RangeStmt:
		p.processRangeLoop(n)
	case *ast.IncDecStmt:
		p.processIncDec(n)
	case *ast.IfStmt:
		p.processIfStmt(n)
	}
}

// processAssignment handles variable assignments and declarations.
func (p *ScriptParser) processAssignment(assign *ast.AssignStmt) {
	stmt := Statement{
		LineNum:  p.fset.Position(assign.Pos()).Line,
		Original: assign,
	}
	if assign.Tok == token.DEFINE {
		stmt.Type = "declaration"
	} else {
		stmt.Type = "assignment"
	}

	stmt.Name = p.formatExpr(assign.Lhs[0])
	stmt.Expr = assign.Rhs[0] // Store the RHS expression

	// Add this line to set the Value field
	stmt.Value = p.formatExpr(assign.Rhs[0])

	p.statements = append(p.statements, stmt)
}

// processFunctionCall handles function calls by storing the function name and arguments.
func (p *ScriptParser) processFunctionCall(call *ast.CallExpr) {
	stmt := Statement{
		Type:     "function_call",
		Name:     p.formatExpr(call.Fun),
		LineNum:  p.fset.Position(call.Pos()).Line,
		ArgExprs: call.Args, // Store argument expressions
	}
	// Keep Args field for backward compatibility if needed
	for _, arg := range call.Args {
		stmt.Args = append(stmt.Args, p.formatExpr(arg))
	}
	p.statements = append(p.statements, stmt)
}

// processIncDec handles increment and decrement statements (e.g., x++, y--).
func (p *ScriptParser) processIncDec(incDec *ast.IncDecStmt) {
	stmt := Statement{
		Type:     "assignment",
		Name:     p.formatExpr(incDec.X),
		Value:    fmt.Sprintf("%s%s", p.formatExpr(incDec.X), incDec.Tok.String()),
		LineNum:  p.fset.Position(incDec.Pos()).Line,
		Original: incDec,
	}
	p.statements = append(p.statements, stmt)
}

// processIfStmt processes an if-statement, capturing the condition and body.
func (p *ScriptParser) processIfStmt(ifStmt *ast.IfStmt) {
	stmt := Statement{
		Type:    "if_statement",
		Cond:    p.formatExpr(ifStmt.Cond),
		LineNum: p.fset.Position(ifStmt.Pos()).Line,
	}
	if ifStmt.Body != nil {
		bodyParser := newScriptParser(p.e)
		bodyParser.fset = p.fset
		for _, node := range ifStmt.Body.List {
			bodyParser.processNode(node)
		}
		stmt.Body = bodyParser.statements
	}
	p.statements = append(p.statements, stmt)
}

// processForLoop processes a for-loop, capturing init, condition, post, and body statements.
func (p *ScriptParser) processForLoop(loop *ast.ForStmt) {
	stmt := Statement{
		Type:    "for_loop",
		LineNum: p.fset.Position(loop.Pos()).Line,
	}
	if loop.Init != nil {
		stmt.Init = p.formatStmt(loop.Init)
	}
	if loop.Cond != nil {
		stmt.Cond = p.formatExpr(loop.Cond)
	}
	if loop.Post != nil {
		stmt.Post = p.formatStmt(loop.Post)
	}

	// Process the body
	if loop.Body != nil {
		bodyParser := newScriptParser(p.e)
		bodyParser.fset = p.fset
		for _, node := range loop.Body.List {
			bodyParser.processNode(node)
		}
		stmt.Body = bodyParser.statements
	}
	p.statements = append(p.statements, stmt)
}

// processRangeLoop processes a range-loop, capturing key, value, range expression, and body.
func (p *ScriptParser) processRangeLoop(loop *ast.RangeStmt) {
	stmt := Statement{
		Type:    "range_loop",
		LineNum: p.fset.Position(loop.Pos()).Line,
	}
	if loop.Key != nil {
		stmt.Name = p.formatExpr(loop.Key)
	}
	if loop.Value != nil {
		stmt.Value = p.formatExpr(loop.Value)
	}
	if loop.X != nil {
		stmt.Args = []string{p.formatExpr(loop.X)}
	}

	// Process body statements
	if loop.Body != nil {
		bodyParser := newScriptParser(p.e)
		bodyParser.fset = p.fset
		for _, node := range loop.Body.List {
			bodyParser.processNode(node)
		}
		stmt.Body = bodyParser.statements
	}
	p.statements = append(p.statements, stmt)
}

// formatExpr formats an expression into a string representation.
func (p *ScriptParser) formatExpr(expr ast.Expr) string {
	switch e := expr.(type) {
	case *ast.BasicLit:
		return e.Value
	case *ast.Ident:
		return e.Name
	case *ast.BinaryExpr:
		return fmt.Sprintf("%s %s %s", p.formatExpr(e.X), e.Op.String(), p.formatExpr(e.Y))
	case *ast.CallExpr:
		return p.formatCallExpr(e)
	case *ast.ParenExpr:
		return fmt.Sprintf("(%s)", p.formatExpr(e.X))
	case *ast.SelectorExpr:
		return fmt.Sprintf("%s.%s", p.formatExpr(e.X), e.Sel.Name)
	case *ast.IndexExpr:
		return fmt.Sprintf("%s[%s]", p.formatExpr(e.X), p.formatExpr(e.Index))
	case *ast.SliceExpr:
		return fmt.Sprintf("%s[%s:%s]", p.formatExpr(e.X), p.formatExpr(e.Low), p.formatExpr(e.High))
	case *ast.StarExpr:
		return fmt.Sprintf("*%s", p.formatExpr(e.X))
	case *ast.UnaryExpr:
		return fmt.Sprintf("%s%s", e.Op.String(), p.formatExpr(e.X))
	case *ast.TypeAssertExpr:
		return fmt.Sprintf("%s.(%s)", p.formatExpr(e.X), p.formatExpr(e.Type))
	case *ast.CompositeLit:
		elements := []string{}
		for _, elt := range e.Elts {
			elements = append(elements, p.formatExpr(elt))
		}
		return fmt.Sprintf("%s{%s}", p.formatExpr(e.Type), strings.Join(elements, ", "))
	case *ast.ArrayType:
		return fmt.Sprintf("[]%s", p.formatExpr(e.Elt))
	default:
		return fmt.Sprintf("%T", e)
	}
}

// Utility functions for formatting assignments and other statements

// formatCallExpr formats a function call expression with arguments.
func (p *ScriptParser) formatCallExpr(call *ast.CallExpr) string {
	funcName := p.formatExpr(call.Fun)
	args := make([]string, len(call.Args))
	for i, arg := range call.Args {
		args[i] = p.formatExpr(arg)
	}
	return fmt.Sprintf("%s(%s)", funcName, strings.Join(args, ", "))
}

// formatStmt formats a statement into a string.
func (p *ScriptParser) formatStmt(stmt ast.Stmt) string {
	switch s := stmt.(type) {
	case *ast.AssignStmt:
		return p.formatAssign(s)
	case *ast.IncDecStmt:
		return p.formatIncDec(s)
	case *ast.ExprStmt:
		return p.formatExpr(s.X)
	default:
		return ""
	}
}

// formatAssign formats an assignment statement into a string.
func (p *ScriptParser) formatAssign(assign *ast.AssignStmt) string {
	lhs := []string{}
	for _, l := range assign.Lhs {
		lhs = append(lhs, p.formatExpr(l))
	}
	rhs := []string{}
	for _, r := range assign.Rhs {
		rhs = append(rhs, p.formatExpr(r))
	}
	return fmt.Sprintf("%s %s %s", strings.Join(lhs, ", "), assign.Tok.String(), strings.Join(rhs, ", "))
}

// formatIncDec formats an increment/decrement statement.
func (p *ScriptParser) formatIncDec(stmt *ast.IncDecStmt) string {
	return fmt.Sprintf("%s%s", p.formatExpr(stmt.X), stmt.Tok.String())
}
