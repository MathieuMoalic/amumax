package api

import (
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/MathieuMoalic/amumax/engine"
	"github.com/MathieuMoalic/amumax/util"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func postSolver(c echo.Context) error {
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

func postConsole(c echo.Context) error {
	type Request struct {
		Command string `json:"command"`
	}

	req := new(Request)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request payload"})
	}
	engine.Inject <- func() { engine.GUI.EvalGUI(req.Command) }
	return c.JSON(http.StatusOK, "")
}

func postTable(c echo.Context) error {
	type Request struct {
		XColumn string `json:"XColumn"`
		YColumn string `json:"YColumn"`
	}

	req := new(Request)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request payload"})
	}
	engine.Tableplot.X = req.XColumn
	engine.Tableplot.Y = req.YColumn
	data := newTablePlot()
	return c.JSON(http.StatusOK, data)
}

func postFrontendState(c echo.Context) error {
	req := new(engine.FrontendState)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request payload"})
	}
	engine.Frontend = *req
	return c.JSON(http.StatusOK, "")
}

func Start() {
	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{"*"},
	}))

	e.Logger.SetOutput(io.Discard)
	e.Static("/", "static")
	e.GET("/ws", websocketEntrypoint)
	e.POST("/solver", postSolver)
	e.POST("/run", postRun)
	e.POST("/console", postConsole)
	e.POST("/steps", postSteps)
	e.POST("/relax", postRelax)
	e.POST("/break", postBreak)
	e.POST("/table", postTable)
	e.POST("/frontendstate", postFrontendState)

	// Start server
	e.Logger.Fatal(e.Start(":5001"))
}
