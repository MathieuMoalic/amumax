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
	e.POST("/solver", postSolverName)
	e.POST("/run", postRun)
	e.POST("/console", postConsole)
	e.POST("/steps", postSteps)
	e.POST("/relax", postRelax)
	e.POST("/break", postBreak)
	e.POST("/table", postTable)
	e.POST("/preview-component", postPreviewComponent)
	e.POST("/preview-quantity", postPreviewQuantity)
	e.POST("/preview-layer", postPreviewLayer)
	e.POST("/preview-maxpoints", postPreviewMaxPoints)
	e.POST("/preview-refresh", postPreviewRefresh)

	// Start server
	e.Logger.Fatal(e.Start(":5001"))
}
