package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/MathieuMoalic/amumax/src/engine_old"
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

func initSolverAPI(e *echo.Group, ws *WebSocketManager) *SolverState {
	solverState := SolverState{
		ws:         ws,
		Type:       getSolverName(engine_old.Solvertype),
		Steps:      &engine_old.NSteps,
		Time:       &engine_old.Time,
		Dt:         &engine_old.Dt_si,
		ErrPerStep: &engine_old.LastErr,
		MaxTorque:  &engine_old.LastTorque,
		Fixdt:      &engine_old.FixDt,
		Mindt:      &engine_old.MinDt,
		Maxdt:      &engine_old.MaxDt,
		Maxerr:     &engine_old.MaxErr,
	}

	e.POST("/api/solver/type", solverState.postSolverType)
	e.POST("/api/solver/run", solverState.postSolverRun)
	e.POST("/api/solver/steps", solverState.postSolverSteps)
	e.POST("/api/solver/minimize", solverState.postSolverMinimize)
	e.POST("/api/solver/relax", solverState.postSolverRelax)
	e.POST("/api/solver/break", solverState.postSolverBreak)
	e.POST("/api/solver/fixdt", solverState.postSolverFixDt)
	e.POST("/api/solver/mindt", solverState.postSolverMinDt)
	e.POST("/api/solver/maxdt", solverState.postSolverMaxDt)
	e.POST("/api/solver/maxerr", solverState.postSolverMaxErr)
	return &solverState
}

func (s *SolverState) Update() {
	s.Type = getSolverName(engine_old.Solvertype)
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
	engine_old.InjectAndWait(func() {
		solver := getSolverType(res.Type)

		// euler must have fixed time step
		if solver == engine_old.EULER && engine_old.FixDt == 0 {
			engine_old.EvalTryRecover("FixDt = 1e-15")
		}
		if solver == engine_old.BACKWARD_EULER && engine_old.FixDt == 0 {
			engine_old.EvalTryRecover("FixDt = 1e-13")
		}

		engine_old.EvalTryRecover(fmt.Sprint("SetSolver(", solver, ")"))
	})

	s.ws.broadcastEngineState()
	return c.JSON(http.StatusOK, engine_old.Solvertype)
}

func (s SolverState) postSolverRun(c echo.Context) error {
	type Request struct {
		Runtime string `msgpack:"runtime"`
	}

	req := new(Request)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "Invalid request payload: " + err.Error(),
		})
	}
	engine_old.Break()
	engine_old.InjectAndWait(func() { engine_old.EvalTryRecover("Run(" + req.Runtime + ")") })
	s.ws.broadcastEngineState()
	return c.JSON(http.StatusOK, "")
}

func (s SolverState) postSolverSteps(c echo.Context) error {
	type Request struct {
		Steps string `msgpack:"steps"`
	}

	req := new(Request)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request payload"})
	}

	engine_old.Break()
	engine_old.InjectAndWait(func() { engine_old.EvalTryRecover("Steps(" + req.Steps + ")") })
	s.ws.broadcastEngineState()
	return c.JSON(http.StatusOK, "")
}

func (s SolverState) postSolverRelax(c echo.Context) error {
	engine_old.Break()
	engine_old.InjectAndWait(func() { engine_old.EvalTryRecover("Relax()") })
	s.ws.broadcastEngineState()
	return c.JSON(http.StatusOK, "")
}

func (s SolverState) postSolverMinimize(c echo.Context) error {
	engine_old.Break()
	engine_old.InjectAndWait(func() { engine_old.EvalTryRecover("Minimize()") })
	s.ws.broadcastEngineState()
	return c.JSON(http.StatusOK, "")
}

func (s SolverState) postSolverBreak(c echo.Context) error {
	engine_old.Break()
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

	engine_old.InjectAndWait(func() { engine_old.EvalTryRecover("FixDt = " + strconv.FormatFloat(req.Fixdt, 'f', -1, 64)) })
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

	engine_old.InjectAndWait(func() { engine_old.EvalTryRecover("MinDt = " + strconv.FormatFloat(req.Mindt, 'f', -1, 64)) })
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

	engine_old.InjectAndWait(func() { engine_old.EvalTryRecover("MaxDt = " + strconv.FormatFloat(req.Maxdt, 'f', -1, 64)) })
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

	engine_old.InjectAndWait(func() { engine_old.EvalTryRecover("MaxErr = " + strconv.FormatFloat(req.Maxerr, 'f', -1, 64)) })
	s.ws.broadcastEngineState()
	return c.JSON(http.StatusOK, "")
}
