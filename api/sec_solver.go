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

func postSolverType(c echo.Context) error {
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

func postSolverRun(c echo.Context) error {
	type Request struct {
		Runtime float64 `msgpack:"runtime"`
	}

	req := new(Request)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request payload"})
	}
	engine.Break()
	engine.Inject <- func() { engine.GUI.EvalGUI("Run(" + strconv.FormatFloat(req.Runtime, 'f', -1, 64) + ")") }
	return c.JSON(http.StatusOK, "")
}

func postSolverSteps(c echo.Context) error {
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

func postSolverRelax(c echo.Context) error {
	engine.Break()
	engine.Inject <- func() { engine.GUI.EvalGUI("Relax()") }
	return c.JSON(http.StatusOK, "")
}

func postSolverBreak(c echo.Context) error {
	engine.Break()
	return c.JSON(http.StatusOK, "")
}

func postSolverFixDt(c echo.Context) error {
	type Request struct {
		Fixdt float64 `msgpack:"fixdt"`
	}

	req := new(Request)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request payload"})
	}

	engine.Inject <- func() { engine.GUI.EvalGUI("FixDt = " + strconv.FormatFloat(req.Fixdt, 'f', -1, 64)) }
	return c.JSON(http.StatusOK, "")
}

func postSolverMinDt(c echo.Context) error {
	type Request struct {
		Mindt float64 `msgpack:"mindt"`
	}

	req := new(Request)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request payload"})
	}

	engine.Inject <- func() { engine.GUI.EvalGUI("MinDt = " + strconv.FormatFloat(req.Mindt, 'f', -1, 64)) }
	return c.JSON(http.StatusOK, "")
}

func postSolverMaxDt(c echo.Context) error {
	type Request struct {
		Maxdt float64 `msgpack:"maxdt"`
	}

	req := new(Request)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request payload"})
	}

	engine.Inject <- func() { engine.GUI.EvalGUI("MaxDt = " + strconv.FormatFloat(req.Maxdt, 'f', -1, 64)) }
	return c.JSON(http.StatusOK, "")
}

func postSolverMaxErr(c echo.Context) error {
	type Request struct {
		Maxerr float64 `msgpack:"maxerr"`
	}

	req := new(Request)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request payload"})
	}

	engine.Inject <- func() { engine.GUI.EvalGUI("MaxErr = " + strconv.FormatFloat(req.Maxerr, 'f', -1, 64)) }
	return c.JSON(http.StatusOK, "")
}
