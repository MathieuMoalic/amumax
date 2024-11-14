package main

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

	statements, err := parser.Parse(script)
	if err != nil {
		t.Fatalf("Error parsing script: %v", err)
	}
	parser.statements = statements

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

	statements, err := parser.Parse(script)
	if err != nil {
		t.Fatalf("Error parsing script: %v", err)
	}
	parser.statements = statements

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
