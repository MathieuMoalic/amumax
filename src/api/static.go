package api

import (
	"embed"
	"io/fs"
	"net/http"

	"github.com/labstack/echo/v4"
)

//go:embed static/*
var staticFiles embed.FS

// Create a sub FS excluding the `index.html` to serve static files.
func staticFileHandler(basePath string) http.Handler {
	fsys, err := fs.Sub(staticFiles, "static")
	if err != nil {
		panic(err)
	}
	fileServer := http.FileServer(http.FS(fsys))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Strip the basePath from the request URL
		http.StripPrefix(basePath+"/", fileServer).ServeHTTP(w, r)
	})
}

// Serve the `index.html` file.
func indexFileHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		indexFile, err := staticFiles.ReadFile("static/index.html")
		if err != nil {
			return c.String(http.StatusInternalServerError, "Index file not found")
		}
		return c.Blob(http.StatusOK, "text/html", indexFile)
	}
}
