package main

import (
	"embed"
	"io/fs"
	"net/http"
	"path"
	"strings"
)

//go:embed all:frontend/dist
var frontendDist embed.FS

// placeholderIndex is served when frontend/dist only contains .gitkeep. It
// keeps `make dev` and `go run .` usable before the first frontend build; the
// Vite dev server is still the normal UI during development.
const placeholderIndex = `<!doctype html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>gollem</title>
    <style>
      body { font-family: system-ui, sans-serif; margin: 0; background: #f8fafc; color: #0f172a; }
      main { max-width: 42rem; margin: 12vh auto; padding: 2rem; }
      code { background: #e2e8f0; padding: .15rem .35rem; border-radius: .25rem; }
    </style>
  </head>
  <body>
    <main>
      <h1>Backend is running</h1>
      <p>The Vue app has not been built into <code>frontend/dist</code> yet.</p>
      <p>During development, open the Vite server (usually <code>http://localhost:5173</code>) or run <code>make dev</code>. For a production build, run <code>make build</code>.</p>
      <p>The API is available at <code>/api/health</code>.</p>
    </main>
  </body>
</html>`

// newFrontendHandler serves the built Vue single-page app from the embedded
// frontend/dist tree. Paths that resolve to a built file are served directly;
// other extension-less paths fall back to index.html so history-mode routes
// survive a hard refresh. Asset misses and paths with a file extension return
// 404 rather than the SPA shell, and /api/* is never handled here (the API mux
// takes precedence).
func newFrontendHandler() http.Handler {
	sub, err := fs.Sub(frontendDist, "frontend/dist")
	if err != nil {
		return placeholderHandler()
	}
	index, err := fs.ReadFile(sub, "index.html")
	if err != nil {
		return placeholderHandler()
	}
	files := http.FileServer(http.FS(sub))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet && r.Method != http.MethodHead {
			w.Header().Set("Allow", "GET, HEAD")
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		upath := path.Clean("/" + r.URL.Path)
		name := strings.TrimPrefix(upath, "/")

		if name != "" {
			if f, err := sub.Open(name); err == nil {
				info, statErr := f.Stat()
				_ = f.Close()
				if statErr == nil && !info.IsDir() {
					setAssetCache(w, upath)
					files.ServeHTTP(w, r)
					return
				}
			}
		}

		if strings.HasPrefix(upath, "/assets/") || path.Ext(upath) != "" {
			http.NotFound(w, r)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Cache-Control", "no-cache")
		_, _ = w.Write(index)
	})
}

func placeholderHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Cache-Control", "no-cache")
		_, _ = w.Write([]byte(placeholderIndex))
	})
}

// setAssetCache marks Vite's content-hashed assets as immutable so browsers
// never revalidate them. Everything else keeps the default revalidation.
func setAssetCache(w http.ResponseWriter, upath string) {
	if strings.HasPrefix(upath, "/assets/") {
		w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
	}
}
