package api

import (
	"net/http"
	"strings"

	"github.com/MathieuMoalic/amumax/src/engine"
	"github.com/MathieuMoalic/amumax/src/log"
	"github.com/labstack/echo/v4"
)

type ConsoleState struct {
	ws   *WebSocketManager
	Hist string `msgpack:"hist"`
}

func (s ConsoleState) Update() {
	// s.Hist = log.Log.Hist
}

func initConsoleAPI(e *echo.Echo, ws *WebSocketManager) *ConsoleState {
	state := &ConsoleState{
		ws:   ws,
		Hist: log.Log.Hist,
	}
	e.POST("/api/console/command", state.postConsoleCommand)
	return state

}

func (s ConsoleState) postConsoleCommand(c echo.Context) error {
	// TODO: return error if the command wrong
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
	s.ws.broadcastEngineState() // Use the instance to call the method
	return c.JSON(http.StatusOK, "")
}
