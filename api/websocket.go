package api

import (
	"time"

	"github.com/MathieuMoalic/amumax/engine"
	"github.com/MathieuMoalic/amumax/util"
	"github.com/labstack/echo/v4"
	"golang.org/x/net/websocket"
)

var (
	webSocketState WebSocketState
)

type WebSocketState struct {
	LastStep int
}

// WebSocket handler for engine state updates
func websocketEntrypoint(c echo.Context) error {
	websocket.Handler(func(ws *websocket.Conn) {
		util.Log("WebSocket client connected: ", ws.RemoteAddr().String(), "->", ws.LocalAddr().String())
		defer ws.Close()

		// Send initial state when the client connects to the WebSocket
		sendEngineState(ws)
		sendDisplayVectorField(ws)

		for {
			if engine.NSteps != webSocketState.LastStep {
				sendEngineState(ws)
				sendDisplayVectorField(ws)
				webSocketState.LastStep = engine.NSteps
			}
			time.Sleep(1 * time.Second)
		}
	}).ServeHTTP(c.Response(), c.Request())
	return nil
}

func sendEngineState(ws *websocket.Conn) {
	engineState := NewEngineState()
	err := websocket.JSON.Send(ws, engineState)
	if err != nil {
		util.LogErr("Error sending engine state via WebSocket:", err)
	}
}

func sendDisplayVectorField(ws *websocket.Conn) {
	engine.InjectAndWait(engine.GetVectorField)
	err := websocket.Message.Send(ws, engine.DisplayVectorField)
	if err != nil {
		util.LogErr("Error sending binary data via WebSocket:", err)
	}
}
