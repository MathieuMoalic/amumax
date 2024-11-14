package parser

import "fmt"

type ScriptFunction func(backend *SimulationBackend, args []interface{}) (interface{}, error)

var functionRegistry = make(map[string]ScriptFunction)

func RegisterFunction(name string, function ScriptFunction) {
	functionRegistry[name] = function
}
func init() {
	RegisterFunction("Vortex", vortexFunction)
	RegisterFunction("Layer", layerFunction)
	RegisterFunction("SetGeom", setGeomFunction)
	RegisterFunction("TableAutoSave", tableAutoSaveFunction)
	RegisterFunction("Run", runSimulationFunction)
}

type Config struct {
	// Define the fields that describe a vortex configuration
	Param1 int
	Param2 int
}

func vortexFunction(backend *SimulationBackend, args []interface{}) (interface{}, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("vortex expects 2 arguments, got %d", len(args))
	}
	param1, ok := args[0].(int)
	if !ok {
		return nil, fmt.Errorf("vortex first argument must be int")
	}
	param2, ok := args[1].(int)
	if !ok {
		return nil, fmt.Errorf("vortex second argument must be int")
	}
	config := Config{Param1: param1, Param2: param2}
	return config, nil
}
func layerFunction(backend *SimulationBackend, args []interface{}) (interface{}, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("layerFunction expects 1 arguments, got %d", len(args))
	}
	param1, ok := args[0].(int)
	if !ok {
		return nil, fmt.Errorf("layerFunction first argument must be int")
	}
	config := Config{Param1: param1, Param2: 1}
	return config, nil
}

func setGeomFunction(backend *SimulationBackend, args []interface{}) (interface{}, error) {
	// Implement SetGeom function
	// Convert args as needed
	return nil, nil
}

func tableAutoSaveFunction(backend *SimulationBackend, args []interface{}) (interface{}, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("TableAutoSave expects 1 argument, got %d", len(args))
	}
	interval, ok := args[0].(float64)
	if !ok {
		return nil, fmt.Errorf("TableAutoSave argument must be float64")
	}
	backend.TableAutoSave(interval)
	return nil, nil
}

func runSimulationFunction(backend *SimulationBackend, args []interface{}) (interface{}, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("run expects 1 argument, got %d", len(args))
	}
	timestep, ok := args[0].(float64)
	if !ok {
		return nil, fmt.Errorf("run argument must be float64")
	}
	backend.RunSimulation(timestep)
	return nil, nil
}
