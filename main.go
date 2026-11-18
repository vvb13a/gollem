package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/vvb13a/gollem/api"
	"github.com/vvb13a/gollem/service"
	"github.com/vvb13a/gollem/store"
)

// defaultAddr is the HTTP listen address for the API and embedded frontend.
const defaultAddr = ":8080"

// version is overridden at build time with -ldflags "-X main.version=...".
var version = "dev"

func main() {
	addr := envOrDefault("APP_ADDR", defaultAddr)

	cfg := service.DefaultConfig()
	cfg.AppName = envOrDefault("APP_NAME", cfg.AppName)

	// 1. Open the database. DataDir can be overridden with APP_DATA_DIR so
	// local development does not write into the per-user OS data directory.
	dbPath, err := store.DatabasePath()
	if err != nil {
		log.Fatalf("failed to resolve database path: %v", err)
	}
	if _, err := store.Open(dbPath); err != nil {
		log.Fatalf("failed to open database: %v", err)
	}

	// 2. Wire the services. Handlers depend on services, never on the DB.

	log.Printf("%s database: %s", cfg.AppName, dbPath)
	runServer(addr, api.Deps{
		AppName: cfg.AppName,
		Version: version,
	})
}

// runServer starts the HTTP server that serves both the JSON API and the
// embedded Vue frontend.
func runServer(addr string, deps api.Deps) {
	frontend := newFrontendHandler()

	apiServer := api.NewServer(deps)
	mux := http.NewServeMux()
	mux.Handle("/api/", apiServer)
	mux.Handle("/", frontend)

	log.Printf("listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("server: %v", err)
	}
}

func envOrDefault(key, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return fallback
}
