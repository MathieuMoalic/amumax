package api

import (
	"net/http"
	"sync"
	"time"

	"github.com/MathieuMoalic/amumax/engine"
	"github.com/MathieuMoalic/amumax/util"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/vmihailenco/msgpack/v5"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	connections = &connectionManager{
		conns: make(map[*websocket.Conn]struct{}),
		mu:    sync.Mutex{},
	}
	lastStep       int
	broadcastStop  chan struct{}
	broadcastStart sync.Once
)

type connectionManager struct {
	conns map[*websocket.Conn]struct{}
	mu    sync.Mutex
}

func (cm *connectionManager) add(ws *websocket.Conn) {
	util.Log.Debug("Websocket connection added")
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.conns[ws] = struct{}{}
	preview.Refresh = true
}

func (cm *connectionManager) remove(ws *websocket.Conn) {
	util.Log.Debug("Websocket connection removed")
	cm.mu.Lock()
	defer cm.mu.Unlock()
	delete(cm.conns, ws)
}

func (cm *connectionManager) broadcast(msg []byte) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	for ws := range cm.conns {
		err := ws.WriteMessage(websocket.BinaryMessage, msg)
		if err != nil {
			util.Log.Err("Error sending message via WebSocket: %v", err)
			ws.Close()
			delete(cm.conns, ws)
		}
	}
}

func websocketEntrypoint(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		util.Log.Err("Error upgrading connection to WebSocket: %v", err)
		return err
	}
	defer ws.Close()

	connections.add(ws)
	defer connections.remove(ws)
	broadcastEngineState()

	// Channel to signal when to stop the goroutine
	done := make(chan struct{})

	go func() {
		for {
			_, _, err := ws.ReadMessage()
			if err != nil {
				close(done)
				return
			}
		}
	}()

	select {
	case <-done:
		util.Log.Debug("Connection closed by client")
		return nil
	case <-broadcastStop:
		return nil
	}
}

func broadcastEngineState() {
	msg, err := msgpack.Marshal(NewEngineState())
	if err != nil {
		util.Log.Err("Error marshaling combined message: %v", err)
		return
	}

	connections.broadcast(msg)
	// Reset the refresh flag
	preview.Refresh = false
}

func startBroadcastLoop() {
	broadcastStop = make(chan struct{})
	broadcastStart.Do(func() {
		go func() {
			for {
				select {
				case <-broadcastStop:
					return
				default:
					if engine.NSteps != lastStep {
						if len(connections.conns) > 0 {
							broadcastEngineState()
							lastStep = engine.NSteps
						}
					}
					time.Sleep(1 * time.Second)
				}
			}
		}()
	})
}

func init() {
	startBroadcastLoop()
}
