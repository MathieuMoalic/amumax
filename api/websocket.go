package api

import (
	"time"

	"github.com/MathieuMoalic/amumax/util"
	"github.com/labstack/echo/v4"
	"golang.org/x/net/websocket"
)

// WebSocket handler for engine state updates
func websocketEntrypoint(c echo.Context) error {
	websocket.Handler(func(ws *websocket.Conn) {
		defer ws.Close()
		util.Log("WebSocket client connected")
		for {
			// Create new engine state
			engineState := NewEngineState()

			// Send engine state to the WebSocket client
			err := websocket.JSON.Send(ws, engineState)
			if err != nil {
				util.Log("Error sending engine state via WebSocket:", err)
				break
			}

			// Wait before sending the next state update
			time.Sleep(1 * time.Second) // Adjust the interval as needed
		}
	}).ServeHTTP(c.Response(), c.Request())
	return nil
}
