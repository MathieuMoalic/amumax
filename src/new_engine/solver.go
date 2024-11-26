package new_engine

type Solver struct {
	EngineState *EngineStateStruct
	Time        float64
}

func NewSolver(EngineState *EngineStateStruct) *Solver {
	s := &Solver{EngineState: EngineState, Time: 0}
	return s
}
