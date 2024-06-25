package api

import (
	"net/http"

	"github.com/MathieuMoalic/amumax/engine"
	"github.com/labstack/echo/v4"
)

type Console struct {
	Hist string `msgpack:"hist"`
}

func newConsole() *Console {
	return &Console{
		Hist: engine.Hist,
	}
}

func postConsole(c echo.Context) error {
	type Request struct {
		Command string `msgpack:"command"`
	}

	req := new(Request)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request payload"})
	}
	engine.Inject <- func() { engine.GUI.EvalGUI(req.Command) }
	return c.JSON(http.StatusOK, "")
}
