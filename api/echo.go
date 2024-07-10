package api

import (
	"io"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Start() {
	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{"*"},
	}))

	e.Logger.SetOutput(io.Discard)
	e.Static("/", "static")
	e.GET("/ws", websocketEntrypoint)

	e.POST("/preview-component", postPreviewComponent)
	e.POST("/preview-quantity", postPreviewQuantity)
	e.POST("/preview-layer", postPreviewLayer)
	e.POST("/preview-maxpoints", postPreviewMaxPoints)
	e.POST("/preview-refresh", postPreviewRefresh)

	e.POST("/tableplot-autosave-interval", postTablePlotAutoSaveInterval)
	e.POST("/tableplot-xcolumn", postTablePlotXColumn)
	e.POST("/tableplot-ycolumn", postTablePlotYColumn)

	e.POST("/console-command", postConsoleCommand)

	e.POST("/solver-type", postSolverType)
	e.POST("/solver-run", postSolverRun)
	e.POST("/solver-steps", postSolverSteps)
	e.POST("/solver-relax", postSolverRelax)
	e.POST("/solver-break", postSolverBreak)
	e.POST("/solver-fixdt", postSolverFixDt)
	e.POST("/solver-mindt", postSolverMinDt)
	e.POST("/solver-maxdt", postSolverMaxDt)
	e.POST("/solver-maxerr", postSolverMaxErr)

	e.POST("/mesh", postMesh)

	// Start server
	e.Logger.Fatal(e.Start(":5001"))
}
