package api

import (
	"fmt"
	"io"
	"net/http"

	"github.com/MathieuMoalic/amumax/engine"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var Timeout = 42

func getDir(c echo.Context) error {
	res := fmt.Sprintf(`{"res":"%s"}`, engine.OD())
	return c.String(http.StatusOK, res)
}

func getTables(c echo.Context) error {
	return c.JSON(http.StatusOK, engine.ZTables)
}

func getImage(c echo.Context) error {
	img := engine.GUI.GetRenderedImg()
	return c.Stream(http.StatusOK, "image/png", img)
}
func getTimeout(c echo.Context) error {
	return c.JSON(http.StatusOK, engine.Timeout)
}
func updateTimeout(c echo.Context) error {
	return c.JSON(http.StatusOK, engine.Timeout)
}

func Start() {
	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{"*"},
	}))

	e.Logger.SetOutput(io.Discard)
	e.Static("/", "static")
	e.GET("/dir", getDir)
	e.GET("/tables", getTables)
	e.GET("/image", getImage)
	e.GET("/timeout", getTimeout)
	e.PUT("/timeout", updateTimeout)

	// Start server
	e.Logger.Fatal(e.Start(":5001"))
}
