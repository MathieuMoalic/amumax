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
	e.POST("/solver", postSolverType)
	e.POST("/run", postRun)
	e.POST("/console", postConsole)
	e.POST("/steps", postSteps)
	e.POST("/relax", postRelax)
	e.POST("/break", postBreak)
	e.POST("/table", postTable)
	e.POST("/frontendstate", postPreviewState)

	// Start server
	e.Logger.Fatal(e.Start(":5001"))
}
