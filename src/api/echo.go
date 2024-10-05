package api

import (
	"fmt"
	"io"
	"net"
	"strconv"
	"time"

	"github.com/MathieuMoalic/amumax/src/log"
	"github.com/MathieuMoalic/amumax/src/script"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Start(host string, port int, tunnel string, debug bool) {
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

	// Serve the `index.html` file at the root URL
	e.GET("/", indexFileHandler())

	// Serve the other embedded static files
	e.GET("/*", echo.WrapHandler(staticFileHandler()))

	wsManager := NewWebSocketManager()
	e.GET("/ws", wsManager.websocketEntrypoint)
	wsManager.startBroadcastLoop()

	// e.POST("/api/preview/component", postPreviewComponent)
	// e.POST("/api/preview/quantity", postPreviewQuantity)
	// e.POST("/api/preview/layer", postPreviewLayer)
	// e.POST("/api/preview/maxpoints", postPreviewMaxPoints)
	// e.POST("/api/preview/refresh", postPreviewRefresh)

	// e.POST("/api/tableplot/autosave-interval", postTablePlotAutoSaveInterval)
	// e.POST("/api/tableplot/xcolumn", postTablePlotXColumn)
	// e.POST("/api/tableplot/ycolumn", postTablePlotYColumn)
	// e.POST("/api/tableplot/maxpoints", postTablePlotMaxPoints)
	// e.POST("/api/tableplot/step", postTablePlotStep)

	initPreviewAPI(e, wsManager)
	// _ = initTablePlotAPI(e, wsManager)
	_ = initConsoleAPI(e, wsManager)
	// _ = initSolverAPI(e, wsManager)
	// _ = initParameterAPI(e, wsManager)
	// _ = initMeshAPI(e, wsManager)

	// e.POST("/api/solver/type", postSolverType)
	// e.POST("/api/solver/run", postSolverRun)
	// e.POST("/api/solver/steps", postSolverSteps)
	// e.POST("/api/solver/relax", postSolverRelax)
	// e.POST("/api/solver/break", postSolverBreak)
	// e.POST("/api/solver/fixdt", postSolverFixDt)
	// e.POST("/api/solver/mindt", postSolverMinDt)
	// e.POST("/api/solver/maxdt", postSolverMaxDt)
	// e.POST("/api/solver/maxerr", postSolverMaxErr)

	// e.POST("/api/parameter/selected-region", postSelectParameterRegion)

	// e.POST("/mesh", postMesh)

	startGuiServer(e, host, port, tunnel)
}

func startGuiServer(e *echo.Echo, host string, port int, tunnel string) {
	const maxRetries = 5

	for i := 0; i < maxRetries; i++ {
		// Find an available port
		addr, port, err := findAvailablePort(host, port)
		if err != nil {
			log.Log.ErrAndExit("Failed to find available port: %v", err)
		}
		log.Log.Info("Serving the web UI at http://%s", addr)

		if tunnel != "" {
			go startTunnel(tunnel)
		}

		script.MMetadata.Add("webui", addr)
		script.MMetadata.Add("port", port)

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

func findAvailablePort(host string, startPort int) (string, string, error) {
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
