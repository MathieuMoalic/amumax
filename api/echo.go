package api

import (
	"io"

	"github.com/MathieuMoalic/amumax/engine"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Start(addr string) {
	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{"*"},
	}))

	e.HideBanner = true
	if engine.VERSION == "dev" {
		e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
			Format: "method=${method}, uri=${uri}, status=${status}\n",
		}))
	} else {
		e.Logger.SetOutput(io.Discard)
	}

	// Serve the `index.html` file at the root URL
	e.GET("/", indexFileHandler())

	// Serve the other embedded static files
	e.GET("/*", echo.WrapHandler(staticFileHandler()))

	e.GET("/ws", websocketEntrypoint)

	e.POST("/api/preview/component", postPreviewComponent)
	e.POST("/api/preview/quantity", postPreviewQuantity)
	e.POST("/api/preview/layer", postPreviewLayer)
	e.POST("/api/preview/maxpoints", postPreviewMaxPoints)
	e.POST("/api/preview/refresh", postPreviewRefresh)

	e.POST("/api/tableplot/autosave-interval", postTablePlotAutoSaveInterval)
	e.POST("/api/tableplot/xcolumn", postTablePlotXColumn)
	e.POST("/api/tableplot/ycolumn", postTablePlotYColumn)

	e.POST("/api/console/command", postConsoleCommand)

	e.POST("/api/solver/type", postSolverType)
	e.POST("/api/solver/run", postSolverRun)
	e.POST("/api/solver/steps", postSolverSteps)
	e.POST("/api/solver/relax", postSolverRelax)
	e.POST("/api/solver/break", postSolverBreak)
	e.POST("/api/solver/fixdt", postSolverFixDt)
	e.POST("/api/solver/mindt", postSolverMinDt)
	e.POST("/api/solver/maxdt", postSolverMaxDt)
	e.POST("/api/solver/maxerr", postSolverMaxErr)

	e.POST("/mesh", postMesh)

	// Start server
	e.Logger.Fatal(e.Start(addr))
}
