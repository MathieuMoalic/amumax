package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/MathieuMoalic/amumax/engine"
	"github.com/MathieuMoalic/amumax/util"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func getImage(c echo.Context) error {
	quantity := c.QueryParam("quantity")
	component := c.QueryParam("component")
	zslice := c.QueryParam("zslice")
	zsliceInt, err := strconv.Atoi(zslice)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid query parameter: zlice is not an integer"})
	}
	if quantity == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Missing required query parameter: quantity"})
	}
	if component == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Missing required query parameter: component"})
	}
	img := engine.GUI.GetRenderedImg(quantity, component, zsliceInt)
	return c.Stream(http.StatusOK, "image/png", img)
}

func getTablePlot(c echo.Context) error {
	x := c.QueryParam("x")
	y := c.QueryParam("y")
	engine.Tableplot.SelectDataColumns(x, y)
	img, err := engine.Tableplot.Render()
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error rendering table plot")
	}
	return c.Stream(http.StatusOK, "image/png", img)
}

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

func Start() {
	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{"*"},
	}))

	// e.Logger.SetOutput(io.Discard)
	e.Static("/", "static")
	e.GET("/ws", websocketEntrypoint)
	e.GET("/tableplot", getTablePlot)
	e.GET("/image", getImage)
	e.POST("/solver", postSolver)
	e.POST("/run", postRun)
	e.POST("/console", postConsole)
	e.POST("/steps", postSteps)
	e.POST("/relax", postRelax)
	e.POST("/break", postBreak)

	// Start server
	e.Logger.Fatal(e.Start(":5001"))
}
