package api

// type EngineState struct {
// 	Header    Header          `msgpack:"header"`
// 	Solver    Solver          `msgpack:"solver"`
// 	Console   ConsoleResponse `msgpack:"console"`
// 	Mesh      Mesh            `msgpack:"mesh"`
// 	Params    Parameters      `msgpack:"parameters"`
// 	TablePlot TablePlot       `msgpack:"tablePlot"`
// 	Preview   Preview         `msgpack:"preview"`
// }

// func NewEngineState() *EngineState {
// 	return &EngineState{
// 		Header:    *newHeader(),
// 		Solver:    *newSolver(),
// 		Console:   *newConsole(),
// 		Mesh:      *newMesh(),
// 		Params:    *newParameters(),
// 		TablePlot: *newTablePlot(),
// 		Preview:   *newPreview(),
// 	}
// }

type EngineState struct {
}

func NewEngineState() *EngineState {
	return &EngineState{}
}
