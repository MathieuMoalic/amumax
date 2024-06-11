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
	Type       string  `json:"type"`
	Steps      int     `json:"steps"`
	Time       float64 `json:"time"`
	Dt         float64 `json:"dt"`
	ErrPerStep float64 `json:"errPerStep"`
	MaxTorque  float64 `json:"maxTorque"`
	Fixdt      float64 `json:"fixdt"`
	Mindt      float64 `json:"mindt"`
	Maxdt      float64 `json:"maxdt"`
	Maxerr     float64 `json:"maxerr"`
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

func postSolverType(c echo.Context) error {
	type Response struct {
		Type string `json:"type"`
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
		Duration float64 `json:"duration"`
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
		Steps int `json:"steps"`
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
