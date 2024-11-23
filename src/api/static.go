package api

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"strings"

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
func indexFileHandler(path string) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Read the index.html file
		indexFile, err := staticFiles.ReadFile("static/index.html")
		if err != nil {
			return c.String(http.StatusInternalServerError, "Index file not found")
		}

		// Convert to string for manipulation
		indexContent := string(indexFile)

		// Inject the `BASE_PATH` variable and `console.log("test")` before the closing </body> tag
		injectedContent := strings.Replace(indexContent, "</body>",
			fmt.Sprintf(`<script>const BASE_PATH = "%s";</script></body>`, path), 1)

		// Serve the modified content
		return c.Blob(http.StatusOK, "text/html", []byte(injectedContent))
	}
}
