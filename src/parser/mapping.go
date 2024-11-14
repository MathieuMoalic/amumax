package parser

import (
	"fmt"
)

// Backend structure to hold and control simulation settings
type SimulationBackend struct {
	Parameters map[string]interface{} // Use interface{} to allow any data type
}

// NewSimulationBackend initializes a new simulation backend
func NewSimulationBackend() *SimulationBackend {
	return &SimulationBackend{
		Parameters: make(map[string]interface{}),
	}
}

// SetParameter sets a simulation parameter (e.g., dx, dy, dz)
func (backend *SimulationBackend) SetParameter(name string, value interface{}) {
	backend.Parameters[name] = value
	fmt.Printf("Set %s to %v\n", name, value)
}

// RunSimulation runs the simulation with the specified timestep
func (backend *SimulationBackend) RunSimulation(timestep float64) {
	fmt.Printf("Running simulation with timestep %e\n", timestep)
	// Here, trigger actual simulation logic
}

// Backend function for other specific functions, e.g., SetGeom
func (backend *SimulationBackend) SetGeometry(params []string) {
	fmt.Printf("Setting geometry with params: %v\n", params)
}

// Example for TableAutoSave function
func (backend *SimulationBackend) TableAutoSave(interval float64) {
	fmt.Printf("Table autosave interval set to %e\n", interval)
}
