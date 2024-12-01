package script

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"

	"github.com/MathieuMoalic/amumax/src/log"
	"github.com/MathieuMoalic/amumax/src/metadata"
)

// ScriptParser parses and stores statements from a script.
type ScriptParser struct {
	statements []Statement
	fset       *token.FileSet
	lineOffset int // Adjusts line numbers for printing
	log        *log.Logs
	scriptStr  *string
	metadata   *metadata.Metadata
	// scope holds the state of the script execution environment.
	functionsScope        map[string]interface{}
	variablesScope        map[string]interface{}
	initializeMeshIfReady func()
}

// Warning: It is called more than once, for example for each for loop and if statements
func (p *ScriptParser) Init(script *string, log *log.Logs, metadata *metadata.Metadata, f func()) {
	p.scriptStr = script
	p.log = log
	p.metadata = metadata
	p.initializeMeshIfReady = f
	p.fset = token.NewFileSet()
	p.lineOffset = -3
	p.functionsScope = make(map[string]interface{})
	p.variablesScope = make(map[string]interface{})
}

// Parse parses a script, wrapping it in a main function to process each statement.
func (p *ScriptParser) Parse() error {
	wrappedScript := "package main\nfunc main() {\n" + *p.scriptStr + "\n}"
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
		bodyParser := ScriptParser{}
		bodyParser.Init(p.scriptStr, p.log, p.metadata, p.initializeMeshIfReady)
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
		bodyParser := ScriptParser{}
		bodyParser.Init(p.scriptStr, p.log, p.metadata, p.initializeMeshIfReady)
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
		bodyParser := ScriptParser{}
		bodyParser.Init(p.scriptStr, p.log, p.metadata, p.initializeMeshIfReady)
		bodyParser.fset = p.fset
		for _, node := range loop.Body.List {
			bodyParser.processNode(node)
		}
		stmt.Body = bodyParser.statements
	}
	p.statements = append(p.statements, stmt)
}
