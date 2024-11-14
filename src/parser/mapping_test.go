package parser

import (
	"testing"
)

// Test setting parameters directly
func TestSetParameter(t *testing.T) {
	backend := NewSimulationBackend()
	backend.SetParameter("dx", 4e-9)

	if val, ok := backend.Parameters["dx"]; !ok || val != 4e-9 {
		t.Errorf("Expected dx to be 4e-9, got %v", val)
	}
}

// Test running the simulation
func TestRunSimulation(t *testing.T) {
	backend := NewSimulationBackend()

	// We don't expect any return value, but we could capture stdout to verify output
	backend.RunSimulation(1e-11)
}

// Test setting geometry
func TestSetGeometry(t *testing.T) {
	backend := NewSimulationBackend()
	args := []string{"Layer(1)"}
	backend.SetGeometry(args)

	// Capture output if needed to verify that SetGeometry was called correctly
}

// Test table autosave
func TestTableAutoSave(t *testing.T) {
	backend := NewSimulationBackend()
	backend.TableAutoSave(1e-11)

	// Since TableAutoSave doesn't return, use output capture if applicable
}

// Test the Execute function to interpret script commands and run backend operations
func TestExecute(t *testing.T) {
	parser := NewScriptParser()
	backend := NewSimulationBackend()

	script := `
		dx = 4e-9
		Nx = 32
		TableAutoSave(1e-11)
		Run(10e-11)
	`

	err := parser.Parse(script)
	if err != nil {
		t.Fatalf("Error parsing script: %v", err)
	}

	// Execute parsed statements
	if err := parser.Execute(backend); err != nil {
		t.Fatalf("Error executing script: %v", err)
	}

	// Verify that parameters were set correctly
	if val, ok := backend.Parameters["dx"]; !ok || val != 4e-9 {
		t.Errorf("Expected dx to be 4e-9, got %v", val)
	}
	if val, ok := backend.Parameters["Nx"]; !ok || val != 32 {
		t.Errorf("Expected Nx to be 32, got %v", val)
	}
}

func TestSetParameterWithDifferentTypes(t *testing.T) {
	backend := NewSimulationBackend()

	// Float parameter
	backend.SetParameter("dx", 4e-9)
	if val, ok := backend.Parameters["dx"].(float64); !ok || val != 4e-9 {
		t.Errorf("Expected dx to be 4e-9, got %v", val)
	}

	// Integer parameter
	backend.SetParameter("Nx", 32)
	if val, ok := backend.Parameters["Nx"].(int); !ok || val != 32 {
		t.Errorf("Expected Nx to be 32, got %v", val)
	}

	// String parameter
	backend.SetParameter("name", "Simulation")
	if val, ok := backend.Parameters["name"].(string); !ok || val != "Simulation" {
		t.Errorf("Expected name to be 'Simulation', got %v", val)
	}
}

func TestExecuteWithDifferentTypes(t *testing.T) {
	parser := NewScriptParser()
	backend := NewSimulationBackend()

	script := `
		dx = 4e-9
		Nx = 32
		name = "Simulation Test"
	`

	err := parser.Parse(script)
	if err != nil {
		t.Fatalf("Error parsing script: %v", err)
	}

	// Execute parsed statements
	if err := parser.Execute(backend); err != nil {
		t.Fatalf("Error executing script: %v", err)
	}

	// Verify that parameters were set correctly
	if val, ok := backend.Parameters["dx"].(float64); !ok || val != 4e-9 {
		t.Errorf("Expected dx to be 4e-9, got %v", val)
	}
	if val, ok := backend.Parameters["Nx"].(int); !ok || val != 32 {
		t.Errorf("Expected Nx to be 32, got %v", val)
	}
	if val, ok := backend.Parameters["name"].(string); !ok || val != "Simulation Test" {
		t.Errorf("Expected name to be 'Simulation Test', got %v", val)
	}
}

func TestExecute_WithEdgeCaseNumbers(t *testing.T) {
	parser := NewScriptParser()
	backend := NewSimulationBackend()

	// Test with extremely small, large, and negative values
	script := `
		dx = 1e-30
		dy = -1.5
		dz = 1e30
	`

	err := parser.Parse(script)
	if err != nil {
		t.Fatalf("Error parsing script: %v", err)
	}

	if err := parser.Execute(backend); err != nil {
		t.Fatalf("Error executing script: %v", err)
	}

	// Check values are correctly parsed and stored
	if val, ok := backend.Parameters["dx"].(float64); !ok || val != 1e-30 {
		t.Errorf("Expected dx to be 1e-30, got %v", val)
	}
	if val, ok := backend.Parameters["dy"].(float64); !ok || val != -1.5 {
		t.Errorf("Expected dy to be -1.5, got %v", val)
	}
	if val, ok := backend.Parameters["dz"].(float64); !ok || val != 1e30 {
		t.Errorf("Expected dz to be 1e30, got %v", val)
	}
}

func TestExecute_InvalidStatements(t *testing.T) {
	parser := NewScriptParser()
	backend := NewSimulationBackend()

	// Test with unsupported command and invalid value
	script := `
		unsupportedFunc(10)
		Nx = "not-a-number"
	`

	err := parser.Parse(script)
	if err != nil {
		t.Fatalf("Error parsing script: %v", err)
	}

	if err := parser.Execute(backend); err == nil {
		t.Errorf("Expected error for unsupported function or invalid value")
	}
}

func TestExecute_MultipleAssignmentsAndReassignments(t *testing.T) {
	parser := NewScriptParser()
	backend := NewSimulationBackend()

	// Test multiple assignments and reassignments
	script := `
		dx = 4e-9
		dy = 1e-9
		dx = 8e-9
		Nx = 32
		Nx = 64
	`

	err := parser.Parse(script)
	if err != nil {
		t.Fatalf("Error parsing script: %v", err)
	}

	if err := parser.Execute(backend); err != nil {
		t.Fatalf("Error executing script: %v", err)
	}

	// Verify reassignments
	if val, ok := backend.Parameters["dx"].(float64); !ok || val != 8e-9 {
		t.Errorf("Expected dx to be updated to 8e-9, got %v", val)
	}
	if val, ok := backend.Parameters["Nx"].(int); !ok || val != 64 {
		t.Errorf("Expected Nx to be updated to 64, got %v", val)
	}
}

func TestExecute_AssignmentToVariable(t *testing.T) {
	parser := NewScriptParser()
	backend := NewSimulationBackend()

	script := `
        a := 1
        b := a
    `

	err := parser.Parse(script)
	if err != nil {
		t.Fatalf("Error parsing script: %v", err)
	}

	if err := parser.Execute(backend); err != nil {
		t.Fatalf("Error executing script: %v", err)
	}

	// Verify assignments
	if val, ok := backend.Parameters["b"].(int); !ok || val != 1 {
		t.Errorf("Expected b to be 1, got %v", val)
	}
}

func TestExecute_WithDifferentFunctions(t *testing.T) {
	parser := NewScriptParser()
	backend := NewSimulationBackend()

	// Test multiple function calls with varying arguments
	script := `
		TableAutoSave(1e-11)
		TableAutoSave(5e-11)
		SetGeom(Layer(1))
		Run(1e-10)
	`

	err := parser.Parse(script)
	if err != nil {
		t.Fatalf("Error parsing script: %v", err)
	}

	// Execute parsed statements
	if err := parser.Execute(backend); err != nil {
		t.Fatalf("Error executing script: %v", err)
	}

	// This test assumes these functions would affect backend state or produce output
	// If functions do not return state changes, consider capturing stdout for validation.
}

func TestSetParameterWithMultipleTypes(t *testing.T) {
	backend := NewSimulationBackend()

	// Test setting a range of types to ensure proper storage and retrieval
	tests := []struct {
		name     string
		value    interface{}
		expected interface{}
	}{
		{"Float", 4e-9, 4e-9},
		{"Int", 32, 32},
		{"String", "TestSimulation", "TestSimulation"},
		{"NegativeInt", -42, -42},
		{"Zero", 0, 0},
	}

	for _, tt := range tests {
		backend.SetParameter(tt.name, tt.value)
		if val, ok := backend.Parameters[tt.name]; !ok || val != tt.expected {
			t.Errorf("Expected %s to be %v, got %v", tt.name, tt.expected, val)
		}
	}
}

func TestExecute_WithStringParameter(t *testing.T) {
	parser := NewScriptParser()
	backend := NewSimulationBackend()

	// Test string parameter assignment
	script := `description = "Simulation Run #1"`

	err := parser.Parse(script)
	if err != nil {
		t.Fatalf("Error parsing script: %v", err)
	}

	if err := parser.Execute(backend); err != nil {
		t.Fatalf("Error executing script: %v", err)
	}

	// Check if the string was stored correctly
	if val, ok := backend.Parameters["description"].(string); !ok || val != "Simulation Run #1" {
		t.Errorf("Expected description to be 'Simulation Run #1', got %v", val)
	}
}

func TestFloatArray(t *testing.T) {
	parser := NewScriptParser()
	backend := NewSimulationBackend()

	// Original array-like syntax
	script := `values = []float64{1.1, 2.2, 3.3}`

	err := parser.Parse(script)
	if err != nil {
		t.Fatalf("Error parsing script: %v", err)
	}

	if err := parser.Execute(backend); err != nil {
		t.Fatalf("Error executing script: %v", err)
	}

	// Validate storage of arrays
	val, ok := backend.Parameters["values"].([]float64)
	if !ok || len(val) != 3 || val[0] != 1.1 || val[1] != 2.2 || val[2] != 3.3 {
		t.Errorf("Expected values to be []float64{1.1, 2.2, 3.3}, got %v", val)
	}
}

func TestIntArray(t *testing.T) {
	parser := NewScriptParser()
	backend := NewSimulationBackend()

	// Original array-like syntax
	script := `values = []int{1, 2, 3}`

	err := parser.Parse(script)
	if err != nil {
		t.Fatalf("Error parsing script: %v", err)
	}

	if err := parser.Execute(backend); err != nil {
		t.Fatalf("Error executing script: %v", err)
	}

	// Validate storage of arrays
	val, ok := backend.Parameters["values"].([]int)
	if !ok || len(val) != 3 || val[0] != 1 || val[1] != 2 || val[2] != 3 {
		t.Errorf("Expected values to be []init{1, 2, 3}, got %v", val)
	}
}

func TestStringArray(t *testing.T) {
	parser := NewScriptParser()
	backend := NewSimulationBackend()

	// Original array-like syntax
	script := `values = []string{"A", "B", "C"}`

	err := parser.Parse(script)
	if err != nil {
		t.Fatalf("Error parsing script: %v", err)
	}

	if err := parser.Execute(backend); err != nil {
		t.Fatalf("Error executing script: %v", err)
	}

	// Validate storage of arrays
	val, ok := backend.Parameters["values"].([]string)
	if !ok || len(val) != 3 || val[0] != "A" || val[1] != "B" || val[2] != "C" {
		t.Errorf("Expected values to be []string{'A', 'B', 'C'}, got %v", val)
	}
}

func TestExecute_WithArraySyntax(t *testing.T) {
	parser := NewScriptParser()
	backend := NewSimulationBackend()

	// Original array-like syntax
	script := `
		values = []float64{1.1, 2.2, 3.3}
		indices = []int{1, 2, 3}
		names = []string{"A", "B", "C"}
	`

	err := parser.Parse(script)
	if err != nil {
		t.Fatalf("Error parsing script: %v", err)
	}

	if err := parser.Execute(backend); err != nil {
		t.Fatalf("Error executing script: %v", err)
	}

	// Validate storage of arrays
	if val, ok := backend.Parameters["values"].([]float64); !ok || len(val) != 3 || val[0] != 1.1 || val[1] != 2.2 || val[2] != 3.3 {
		t.Errorf("Expected values to be []float64{1.1, 2.2, 3.3}, got %v", val)
	}
	if val, ok := backend.Parameters["indices"].([]int); !ok || len(val) != 3 || val[0] != 1 || val[1] != 2 || val[2] != 3 {
		t.Errorf("Expected indices to be []int{1, 2, 3}, got %v", val)
	}
	if val, ok := backend.Parameters["names"].([]string); !ok || len(val) != 3 || val[0] != "A" || val[1] != "B" || val[2] != "C" {
		t.Errorf("Expected names to be []string{'A', 'B', 'C'}, got %v", val)
	}
}

func TestExecute_WithErrorHandling(t *testing.T) {
	parser := NewScriptParser()
	backend := NewSimulationBackend()

	// Test script with invalid statement to trigger error handling
	script := `
		dx = 4e-9
		invalidFunc(1e-9)
	`

	err := parser.Parse(script)
	if err != nil {
		t.Fatalf("Error parsing script: %v", err)
	}

	// Expect error during execution due to unsupported function call
	if err := parser.Execute(backend); err == nil {
		t.Error("Expected error due to unsupported function call, got nil")
	}
}
