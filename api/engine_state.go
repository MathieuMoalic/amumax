package api

type EngineState struct {
	Header    Header     `json:"header"`
	Solver    Solver     `json:"solver"`
	Console   Console    `json:"console"`
	Mesh      Mesh       `json:"mesh"`
	Params    Parameters `json:"parameters"`
	TablePlot TablePlot  `json:"tablePlot"`
}

func NewEngineState() *EngineState {
	return &EngineState{
		Header:    *newHeader(),
		Solver:    *newSolver(),
		Console:   *newConsole(),
		Mesh:      *newMesh(),
		Params:    *newParameters(),
		TablePlot: *newTablePlot(),
	}
}
