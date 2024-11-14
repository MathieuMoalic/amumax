package main

import (
	"reflect"
	"testing"
)

func statementsEqual(a, b []Statement) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i].Type != b[i].Type ||
			a[i].Name != b[i].Name ||
			a[i].Value != b[i].Value ||
			!reflect.DeepEqual(a[i].Args, b[i].Args) ||
			a[i].LineNum != b[i].LineNum ||
			a[i].Init != b[i].Init ||
			a[i].Cond != b[i].Cond ||
			a[i].Post != b[i].Post {
			return false
		}
		// Recursively check the Body field without comparing Original
		if !statementsEqual(a[i].Body, b[i].Body) {
			return false
		}
	}
	return true
}

func TestParseAssignment(t *testing.T) {
	parser := NewScriptParser()
	script := `x := 10`

	statements, err := parser.Parse(script)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	expected := []Statement{
		{
			Type:    "declaration",
			Name:    "x",
			Value:   "10",
			LineNum: 1,
		},
	}

	if !statementsEqual(statements, expected) {
		t.Errorf("Expected %v, got %v", expected, statements)
	}
}

func TestParseFunctionCall(t *testing.T) {
	parser := NewScriptParser()
	script := `Println("Hello, World")`

	statements, err := parser.Parse(script)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	expected := []Statement{
		{
			Type:    "function_call",
			Name:    "Println",
			Args:    []string{`"Hello, World"`},
			LineNum: 1,
		},
	}

	if !statementsEqual(statements, expected) {
		t.Errorf("Expected %v, got %v", expected, statements)
	}
}

func TestParseForLoop(t *testing.T) {
	parser := NewScriptParser()
	script := `for i := 0; i < 5; i++ { Println(i) }`

	statements, err := parser.Parse(script)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	expected := []Statement{
		{
			Type:    "for_loop",
			Init:    "i := 0",
			Cond:    "i < 5",
			Post:    "i++",
			LineNum: 1,
			Body: []Statement{
				{
					Type:    "function_call",
					Name:    "Println",
					Args:    []string{"i"},
					LineNum: 1,
				},
			},
		},
	}

	if !statementsEqual(statements, expected) {
		t.Errorf("Expected %v, got %v", expected, statements)
	}
}

func TestParseRangeLoop(t *testing.T) {
	parser := NewScriptParser()
	script := `values := []int{1, 2, 3}
for _, v := range values { Println(v) }`

	statements, err := parser.Parse(script)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	expected := []Statement{
		{
			Type:    "declaration",
			Name:    "values",
			Value:   "[]int{1, 2, 3}",
			LineNum: 1,
		},
		{
			Type:    "range_loop",
			Name:    "_",
			Value:   "v",
			Args:    []string{"values"},
			LineNum: 2,
			Body: []Statement{
				{
					Type:    "function_call",
					Name:    "Println",
					Args:    []string{"v"},
					LineNum: 2,
				},
			},
		},
	}

	if !statementsEqual(statements, expected) {
		t.Errorf("Expected %v, got %v", expected, statements)
	}
}

func TestParseMultipleAssignments(t *testing.T) {
	parser := NewScriptParser()
	script := `x := 10
y := x + 5`

	statements, err := parser.Parse(script)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	expected := []Statement{
		{
			Type:    "declaration",
			Name:    "x",
			Value:   "10",
			LineNum: 1,
		},
		{
			Type:    "declaration",
			Name:    "y",
			Value:   "x + 5",
			LineNum: 2,
		},
	}

	if !statementsEqual(statements, expected) {
		t.Errorf("Expected %v, got %v", expected, statements)
	}
}

func TestParseIncDec(t *testing.T) {
	parser := NewScriptParser()
	script := `x := 10
x++`

	statements, err := parser.Parse(script)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	expected := []Statement{
		{
			Type:    "declaration",
			Name:    "x",
			Value:   "10",
			LineNum: 1,
		},
		{
			Type:    "assignment",
			Name:    "x",
			Value:   "x++",
			LineNum: 2,
		},
	}

	if !statementsEqual(statements, expected) {
		t.Errorf("Expected %v, got %v", expected, statements)
	}
}

func TestParseBinaryOperation(t *testing.T) {
	parser := NewScriptParser()
	script := `x := 5 + 3`

	statements, err := parser.Parse(script)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	expected := []Statement{
		{
			Type:    "declaration",
			Name:    "x",
			Value:   "5 + 3",
			LineNum: 1,
		},
	}

	if !statementsEqual(statements, expected) {
		t.Errorf("Expected %v, got %v", expected, statements)
	}
}

func TestParseArrayIndexing(t *testing.T) {
	parser := NewScriptParser()
	script := `arr := []int{1, 2, 3}
x := arr[1]`

	statements, err := parser.Parse(script)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	expected := []Statement{
		{
			Type:    "declaration",
			Name:    "arr",
			Value:   "[]int{1, 2, 3}",
			LineNum: 1,
		},
		{
			Type:    "declaration",
			Name:    "x",
			Value:   "arr[1]",
			LineNum: 2,
		},
	}

	if !statementsEqual(statements, expected) {
		t.Errorf("Expected %v, got %v", expected, statements)
	}
}

func TestParseNestedFunctionCall(t *testing.T) {
	parser := NewScriptParser()
	script := `x := max(5, min(2, 3))`

	statements, err := parser.Parse(script)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	expected := []Statement{
		{
			Type:    "declaration",
			Name:    "x",
			Value:   "max(5, min(2, 3))",
			LineNum: 1,
		},
	}

	if !statementsEqual(statements, expected) {
		t.Errorf("Expected %v, got %v", expected, statements)
	}
}

func TestParseIfStatement(t *testing.T) {
	parser := NewScriptParser()
	script := `if x > 5 { Println("Greater") }`

	statements, err := parser.Parse(script)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	expected := []Statement{
		{
			Type:    "if_statement",
			Cond:    "x > 5",
			LineNum: 1,
			Body: []Statement{
				{
					Type:    "function_call",
					Name:    "Println",
					Args:    []string{`"Greater"`},
					LineNum: 1,
				},
			},
		},
	}

	if !statementsEqual(statements, expected) {
		t.Errorf("Expected %v, got %v", expected, statements)
	}
}

func TestParseAssignmentWithBinaryExpr(t *testing.T) {
	parser := NewScriptParser()
	script := `y := x * (z + 1)`

	statements, err := parser.Parse(script)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	expected := []Statement{
		{
			Type:    "declaration",
			Name:    "y",
			Value:   "x * (z + 1)",
			LineNum: 1,
		},
	}

	if !statementsEqual(statements, expected) {
		t.Errorf("Expected %v, got %v", expected, statements)
	}
}
