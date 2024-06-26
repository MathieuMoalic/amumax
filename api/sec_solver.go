package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/MathieuMoalic/amumax/engine"
	"github.com/MathieuMoalic/amumax/util"
	"github.com/labstack/echo/v4"
)

type Solver struct {
	Type       string  `msgpack:"type"`
	Steps      int     `msgpack:"steps"`
	Time       float64 `msgpack:"time"`
	Dt         float64 `msgpack:"dt"`
	ErrPerStep float64 `msgpack:"errPerStep"`
	MaxTorque  float64 `msgpack:"maxTorque"`
	Fixdt      float64 `msgpack:"fixdt"`
	Mindt      float64 `msgpack:"mindt"`
	Maxdt      float64 `msgpack:"maxdt"`
	Maxerr     float64 `msgpack:"maxerr"`
}

func newSolver() *Solver {
	return &Solver{
		Type:       engine.Solvernames[engine.Solvertype],
		Steps:      engine.NSteps,
		Time:       engine.Time,
		Dt:         engine.Dt_si,
		ErrPerStep: engine.LastErr,
		MaxTorque:  engine.LastTorque,
		Fixdt:      engine.FixDt,
		Mindt:      engine.MinDt,
		Maxdt:      engine.MaxDt,
		Maxerr:     engine.MaxErr,
	}
}

func postSolverName(c echo.Context) error {
	type Response struct {
		Type string `msgpack:"type"`
	}

	res := new(Response)
	if err := c.Bind(res); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request payload"})
	}
	engine.Inject <- func() {
		solver := engine.Solvertypes[res.Type]

		// euler must have fixed time step
		if solver == engine.EULER && engine.FixDt == 0 {
			engine.GUI.EvalGUI("FixDt = 1e-15")
		}
		if solver == engine.BACKWARD_EULER && engine.FixDt == 0 {
			engine.GUI.EvalGUI("FixDt = 1e-13")
		}
		util.Log("SetSolver: %v", solver)

		engine.GUI.EvalGUI(fmt.Sprint("SetSolver(", solver, ")"))
	}
	// engine.Solvertype = engine.Solvertypes[res.Type]
	return c.JSON(http.StatusOK, engine.Solvertype)
}

func postRun(c echo.Context) error {
	type Request struct {
		Duration float64 `msgpack:"duration"`
	}

	req := new(Request)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request payload"})
	}
	engine.Break()
	engine.Inject <- func() { engine.GUI.EvalGUI("Run(" + strconv.FormatFloat(req.Duration, 'f', -1, 64) + ")") }
	return c.JSON(http.StatusOK, "")
}

func postSteps(c echo.Context) error {
	type Request struct {
		Steps int `msgpack:"steps"`
	}

	req := new(Request)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request payload"})
	}

	engine.Break()
	engine.Inject <- func() { engine.GUI.EvalGUI("Steps(" + strconv.Itoa(req.Steps) + ")") }
	return c.JSON(http.StatusOK, "")
}

func postRelax(c echo.Context) error {
	engine.Break()
	engine.Inject <- func() { engine.GUI.EvalGUI("Relax()") }
	return c.JSON(http.StatusOK, "")
}

func postBreak(c echo.Context) error {
	engine.Break()
	return c.JSON(http.StatusOK, "")
}
