package api

import (
	"net/http"
	"strings"

	"github.com/MathieuMoalic/amumax/engine"
	"github.com/MathieuMoalic/amumax/util"
	"github.com/labstack/echo/v4"
)

type Console struct {
	Hist string `msgpack:"hist"`
}

func newConsole() *Console {
	return &Console{
		Hist: util.Log.Hist,
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

	// Check if "RunShell" is in the command, case-insensitive
	if strings.Contains(strings.ToLower(req.Command), "runshell(") {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "RunShell command not allowed through the WebUI"})
	}

	engine.InjectAndWait(func() { engine.EvalTryRecover(req.Command) })
	broadcastEngineState()
	return c.JSON(http.StatusOK, "")
}
