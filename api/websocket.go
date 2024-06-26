package api

import (
	"sync"
	"time"

	"github.com/MathieuMoalic/amumax/engine"
	"github.com/MathieuMoalic/amumax/util"
	"github.com/labstack/echo/v4"
	"github.com/vmihailenco/msgpack/v5"
	"golang.org/x/net/websocket"
)

var (
	webSocketState WebSocketState
	connections    = &connectionManager{
		conns: make(map[*websocket.Conn]bool),
		mu:    sync.Mutex{},
	}
)

type WebSocketState struct {
	LastStep int
}

type connectionManager struct {
	conns map[*websocket.Conn]bool
	mu    sync.Mutex
}

func (cm *connectionManager) add(ws *websocket.Conn) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.conns[ws] = true
}

func (cm *connectionManager) remove(ws *websocket.Conn) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	delete(cm.conns, ws)
}

func (cm *connectionManager) broadcast(msg []byte) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	for ws := range cm.conns {
		err := websocket.Message.Send(ws, msg)
		if err != nil {
			util.LogErr("Error sending message via WebSocket:", err)
			ws.Close()
			delete(cm.conns, ws)
		}
	}
}

// WebSocket handler for engine state updates
func websocketEntrypoint(c echo.Context) error {
	websocket.Handler(func(ws *websocket.Conn) {
		util.Log("WebSocket client connected: ", ws.RemoteAddr().String(), "->", ws.LocalAddr().String())
		defer ws.Close()

		// Register the WebSocket connection
		connections.add(ws)
		defer connections.remove(ws)

		// Send initial state when the client connects to the WebSocket
		broadcastEngineState()

		for {
			if engine.NSteps != webSocketState.LastStep {
				broadcastEngineState()
				webSocketState.LastStep = engine.NSteps
			}
			time.Sleep(1 * time.Second)
		}
	}).ServeHTTP(c.Response(), c.Request())
	return nil
}

func broadcastEngineState() {
	engine.InjectAndWait(PreparePreviewBuffer)

	msg, err := msgpack.Marshal(NewEngineState())
	if err != nil {
		util.LogErr("Error marshaling combined message:", err)
		return
	}

	connections.broadcast(msg)
}
