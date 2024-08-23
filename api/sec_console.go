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

func postConsoleCommand(c echo.Context) error {
	type Request struct {
		Command string `msgpack:"command"`
	}

	req := new(Request)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request payload"})
	}
	engine.Inject <- func() { engine.EvalTryRecover(req.Command) }
	broadcastEngineState()
	return c.JSON(http.StatusOK, "")
}
