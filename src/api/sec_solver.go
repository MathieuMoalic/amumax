package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/MathieuMoalic/amumax/src/engine"
	"github.com/labstack/echo/v4"
)

type SolverState struct {
	ws         *WebSocketManager
	Type       string   `msgpack:"type"`
	Steps      *int     `msgpack:"steps"`
	Time       *float64 `msgpack:"time"`
	Dt         *float64 `msgpack:"dt"`
	ErrPerStep *float64 `msgpack:"errPerStep"`
	MaxTorque  *float64 `msgpack:"maxTorque"`
	Fixdt      *float64 `msgpack:"fixdt"`
	Mindt      *float64 `msgpack:"mindt"`
	Maxdt      *float64 `msgpack:"maxdt"`
	Maxerr     *float64 `msgpack:"maxerr"`
}

func initSolverAPI(e *echo.Echo, ws *WebSocketManager) *SolverState {
	solverState := SolverState{
		ws:         ws,
		Type:       getSolverName(engine.Solvertype),
		Steps:      &engine.NSteps,
		Time:       &engine.Time,
		Dt:         &engine.Dt_si,
		ErrPerStep: &engine.LastErr,
		MaxTorque:  &engine.LastTorque,
		Fixdt:      &engine.FixDt,
		Mindt:      &engine.MinDt,
		Maxdt:      &engine.MaxDt,
		Maxerr:     &engine.MaxErr,
	}

	e.POST("/api/solver/type", solverState.postSolverType)
	e.POST("/api/solver/run", solverState.postSolverRun)
	e.POST("/api/solver/steps", solverState.postSolverSteps)
	e.POST("/api/solver/relax", solverState.postSolverRelax)
	e.POST("/api/solver/break", solverState.postSolverBreak)
	e.POST("/api/solver/fixdt", solverState.postSolverFixDt)
	e.POST("/api/solver/mindt", solverState.postSolverMinDt)
	e.POST("/api/solver/maxdt", solverState.postSolverMaxDt)
	e.POST("/api/solver/maxerr", solverState.postSolverMaxErr)
	return &solverState
}

func (s *SolverState) Update() {
	s.Type = getSolverName(engine.Solvertype)
}

func getSolverType(typeStr string) int {
	solvertypes := map[string]int{"bw_euler": -1, "euler": 1, "heun": 2, "rk23": 3, "rk4": 4, "rk45": 5, "rkf56": 6}
	if solver, ok := solvertypes[typeStr]; ok {
		return solver
	}
	return 0
}

func getSolverName(solver int) string {
	solvernames := map[int]string{-1: "bw_euler", 1: "euler", 2: "heun", 3: "rk23", 4: "rk4", 5: "rk45", 6: "rkf56"}
	if name, ok := solvernames[solver]; ok {
		return name
	}
	return ""
}

func (s SolverState) postSolverType(c echo.Context) error {
	type Response struct {
		Type string `msgpack:"type"`
	}

	res := new(Response)
	if err := c.Bind(res); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request payload"})
	}
	engine.InjectAndWait(func() {
		solver := getSolverType(res.Type)

		// euler must have fixed time step
		if solver == engine.EULER && engine.FixDt == 0 {
			engine.EvalTryRecover("FixDt = 1e-15")
		}
		if solver == engine.BACKWARD_EULER && engine.FixDt == 0 {
			engine.EvalTryRecover("FixDt = 1e-13")
		}

		engine.EvalTryRecover(fmt.Sprint("SetSolver(", solver, ")"))
	})

	s.ws.broadcastEngineState()
	return c.JSON(http.StatusOK, engine.Solvertype)
}

func (s SolverState) postSolverRun(c echo.Context) error {
	type Request struct {
		Runtime float64 `msgpack:"runtime"`
	}

	req := new(Request)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request payload"})
	}
	engine.Break()
	engine.InjectAndWait(func() { engine.EvalTryRecover("Run(" + strconv.FormatFloat(req.Runtime, 'f', -1, 64) + ")") })
	s.ws.broadcastEngineState()
	return c.JSON(http.StatusOK, "")
}

func (s SolverState) postSolverSteps(c echo.Context) error {
	type Request struct {
		Steps int `msgpack:"steps"`
	}

	req := new(Request)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request payload"})
	}

	engine.Break()
	engine.InjectAndWait(func() { engine.EvalTryRecover("Steps(" + strconv.Itoa(req.Steps) + ")") })
	s.ws.broadcastEngineState()
	return c.JSON(http.StatusOK, "")
}

func (s SolverState) postSolverRelax(c echo.Context) error {
	engine.Break()
	engine.InjectAndWait(func() { engine.EvalTryRecover("Relax()") })
	s.ws.broadcastEngineState()
	return c.JSON(http.StatusOK, "")
}

func (s SolverState) postSolverBreak(c echo.Context) error {
	engine.Break()
	s.ws.broadcastEngineState()
	return c.JSON(http.StatusOK, "")
}

func (s SolverState) postSolverFixDt(c echo.Context) error {
	type Request struct {
		Fixdt float64 `msgpack:"fixdt"`
	}

	req := new(Request)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request payload"})
	}

	engine.InjectAndWait(func() { engine.EvalTryRecover("FixDt = " + strconv.FormatFloat(req.Fixdt, 'f', -1, 64)) })
	s.ws.broadcastEngineState()
	return c.JSON(http.StatusOK, "")
}

func (s SolverState) postSolverMinDt(c echo.Context) error {
	type Request struct {
		Mindt float64 `msgpack:"mindt"`
	}

	req := new(Request)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request payload"})
	}

	engine.InjectAndWait(func() { engine.EvalTryRecover("MinDt = " + strconv.FormatFloat(req.Mindt, 'f', -1, 64)) })
	s.ws.broadcastEngineState()
	return c.JSON(http.StatusOK, "")
}

func (s SolverState) postSolverMaxDt(c echo.Context) error {
	type Request struct {
		Maxdt float64 `msgpack:"maxdt"`
	}

	req := new(Request)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request payload"})
	}

	engine.InjectAndWait(func() { engine.EvalTryRecover("MaxDt = " + strconv.FormatFloat(req.Maxdt, 'f', -1, 64)) })
	s.ws.broadcastEngineState()
	return c.JSON(http.StatusOK, "")
}

func (s SolverState) postSolverMaxErr(c echo.Context) error {
	type Request struct {
		Maxerr float64 `msgpack:"maxerr"`
	}

	req := new(Request)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request payload"})
	}

	engine.InjectAndWait(func() { engine.EvalTryRecover("MaxErr = " + strconv.FormatFloat(req.Maxerr, 'f', -1, 64)) })
	s.ws.broadcastEngineState()
	return c.JSON(http.StatusOK, "")
}
