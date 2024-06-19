package api

import (
	"time"

	"github.com/MathieuMoalic/amumax/engine"
	"github.com/MathieuMoalic/amumax/util"
	"github.com/labstack/echo/v4"
	"github.com/vmihailenco/msgpack/v5"
	"golang.org/x/net/websocket"
)

var (
	webSocketState WebSocketState
)

type WebSocketState struct {
	LastStep int
}

type WebSocketMessage struct {
	EngineState   *EngineState `msgpack:"engine_state"`
	PreviewBuffer *[]byte      `msgpack:"preview_buffer"`
}

// WebSocket handler for engine state updates
func websocketEntrypoint(c echo.Context) error {
	websocket.Handler(func(ws *websocket.Conn) {
		util.Log("WebSocket client connected: ", ws.RemoteAddr().String(), "->", ws.LocalAddr().String())
		defer ws.Close()

		// Send initial state when the client connects to the WebSocket
		sendMessage(ws)

		for {
			if engine.NSteps != webSocketState.LastStep {
				sendMessage(ws)
				webSocketState.LastStep = engine.NSteps
			}
			time.Sleep(1 * time.Second)
		}
	}).ServeHTTP(c.Response(), c.Request())
	return nil
}

func sendMessage(ws *websocket.Conn) {
	engine.InjectAndWait(PreparePreviewBuffer)
	rawMessage := WebSocketMessage{
		EngineState:   NewEngineState(),
		PreviewBuffer: &PreviewBuffer,
	}

	msg, err := msgpack.Marshal(rawMessage)
	if err != nil {
		util.LogErr("Error marshaling combined message:", err)
		return
	}

	err = websocket.Message.Send(ws, msg)
	if err != nil {
		util.LogErr("Error sending combined message via WebSocket:", err)
	}
}
