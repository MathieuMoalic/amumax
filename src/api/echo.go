package api

import (
	"fmt"
	"io"
	"net"
	"strconv"
	"time"

	"github.com/MathieuMoalic/amumax/src/engine"
	"github.com/MathieuMoalic/amumax/src/engine/log"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Start(host string, port int, basePath string, tunnel string, debug bool) {
	defer func() {
		if r := recover(); r != nil {
			log.Log.Warn("WebUI crashed: %v", r)
		}
	}()
	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{"*"},
	}))

	e.HideBanner = true
	if debug {
		e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
			Format: "method=${method}, uri=${uri}, status=${status}\n",
		}))
	} else {
		e.Logger.SetOutput(io.Discard)
	}

	api := e.Group(basePath)

	// redirect "" to "/"
	api.GET("", func(c echo.Context) error {
		return c.Redirect(301, basePath+"/")
	})

	// Serve the `index.html` file at the root URL
	api.GET("/", indexFileHandler())

	// Serve the other embedded static files
	api.GET("/*", echo.WrapHandler(staticFileHandler(basePath)))

	wsManager := newWebSocketManager()
	api.GET("/ws", wsManager.websocketEntrypoint)
	wsManager.startBroadcastLoop()
	engineState := initEngineStateAPI(api, wsManager)
	wsManager.engineState = engineState

	startGuiServer(e, host, basePath, port, tunnel)
}

func startGuiServer(e *echo.Echo, host string, basePath string, port int, tunnel string) {
	const maxRetries = 5

	for i := 0; i < maxRetries; i++ {
		// Find an available port
		addr, port, err := FindAvailablePort(host, port)
		if err != nil {
			log.Log.ErrAndExit("Failed to find available port: %v", err)
		}
		log.Log.Info("Serving the web UI at http://%s%s", addr, basePath)

		if tunnel != "" {
			go startTunnel(tunnel)
		}

		engine.EngineState.Metadata.Add("webui", addr)
		engine.EngineState.Metadata.Add("port", port)

		// Attempt to start the server
		err = e.Start(addr)
		if err != nil {
			if opErr, ok := err.(*net.OpError); ok && opErr.Op == "listen" {
				// Port is already in use, retrying
				time.Sleep(1 * time.Second) // Wait before retrying
				continue
			}
			// If the error is not related to the port being busy, exit
			log.Log.Err("Failed to start server:  %v", err)
			break
		}

		// If the server started successfully, break out of the loop
		log.Log.Info("Successfully started server at http://%s", addr)
		return
	}

	// If the loop completes without successfully starting the server
	log.Log.Err("Failed to start server after multiple attempts")
}

func FindAvailablePort(host string, startPort int) (string, string, error) {
	// Loop to find the first available port
	for port := startPort; port <= 65535; port++ {
		address := net.JoinHostPort(host, strconv.Itoa(port))
		listener, err := net.Listen("tcp", address)
		if err == nil {
			// Close the listener immediately, we just wanted to check availability
			listener.Close()
			return address, strconv.Itoa(port), nil
		}
	}
	return "", "", fmt.Errorf("no available ports found")
}
