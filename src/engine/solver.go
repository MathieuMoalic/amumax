package engine

type solver struct {
	e      *engineState
	time   float64
	NSteps int
}

func newSolver(EngineState *engineState) *solver {
	s := &solver{e: EngineState, time: 0}
	return s
}
