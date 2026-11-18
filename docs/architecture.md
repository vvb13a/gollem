# Architecture

This project is a small full-stack application with a deliberately strict
layer split. The same shape is used by every feature.

## Folder map

```text
.
├── api/                  HTTP layer
├── domain/               Business types and domain errors
├── service/              Use cases / business logic
├── store/                SQLite persistence, models, migrations
├── frontend/             Vue + Vite application
│   ├── src/api/          Typed API client
│   ├── src/components/   Vue components
│   ├── src/lib/          Small shared frontend modules
│   └── dist/             Built SPA embedded into the Go binary
├── main.go               Composition root (config, DB, services, server)
├── static.go             Serves the embedded SPA with history-mode fallback
├── Makefile              Developer tasks, including `make dev`
└── .env.example          Local configuration example
```

### `api/`

Owns HTTP only: routing, request decoding, status codes, and response
encoding. Handlers call services. Keep business rules and SQL out of this
package.

- `server.go` — `Deps`, route registration, JSON/error helpers
- Feature handler files, for example `widgets.go`

### `domain/`

Plain Go types and domain errors. No database or HTTP imports. This is what the
service and API layers exchange.

- Feature types
- `errors.go` — stable sentinel errors such as `ErrNotFound`

### `service/`

Application behavior: validation, orchestration, transactions, and the rules
that should hold regardless of transport. Services are the only layer that
GORM models should cross into.

- `config.go` — deployment/app configuration
- `errors.go` — service-level errors such as `ErrInvalidInput`
- Feature service files, for example `widget_service.go`

### `store/`

Persistence details: GORM models, table names, migrations, and database
opening/path resolution. Only this package should import GORM or the SQLite
driver.

- `open.go` — opens SQLite, applies migrations, serializes writers
- `migrations.go` — embeds and runs Goose migrations
- `migrations/00001_init.sql` — ground-zero schema
- `path.go` — per-user data directory, overridable with `APP_DATA_DIR`
- Feature persistence models, for example `widget.go`

### `frontend/`

Vue 3 + TypeScript app built by Vite.

- `src/api/client.ts` — the only place that knows API URLs
- `src/components/` — feature UI
- `src/lib/theme.ts` — light/dark theme helper
- `src/router.ts` — Vue Router setup
- `src/style.css` — Tailwind import and class-based dark variant
- `dist/` — build output. The committed `.gitkeep` lets `go:embed` compile on
  a clean checkout; `make build-web` replaces the contents.

### Root files

- `main.go` wires configuration, the database, services, and the HTTP server.
- `static.go` embeds `frontend/dist` and serves the SPA. If `index.html` has
  not been built yet, it serves a small placeholder so `go run .` still starts.
- `Makefile` is the entry point for development and builds.
- `AGENTS.md` contains the rules for contributors and AI assistants.

## Request flow

```text
Vue component
  -> frontend/src/api/client.ts
  -> /api/... (Vite proxy in dev)
  -> api.Server handler
  -> service.FeatureService
  -> store.Feature (GORM)
  -> SQLite
```

## Development topology

`make dev` starts two processes:

1. Go on `:8080` serving `/api` and, when built, the embedded SPA.
2. Vite on `:5173` serving the frontend and proxying `/api` to Go.

The database lives in `./.dev/data` during development by default. Outside
development it uses the platform data directory, unless `APP_DATA_DIR` is set.

## Database and migrations

The database is SQLite through a pure-Go driver. WAL mode, foreign keys, and a
busy timeout are enabled, and the connection pool is intentionally limited to
one connection because SQLite serializes writers.

Add schema changes as new Goose migrations:

```sql
-- store/migrations/00002_add_widget_status.sql
-- +goose Up
ALTER TABLE widgets ADD COLUMN status TEXT NOT NULL DEFAULT 'open';

-- +goose Down
ALTER TABLE widgets DROP COLUMN status;
```

Migrations run automatically on startup. Never edit an applied migration.

## Layer rules

- `api` never touches the database.
- `service` never writes HTTP status codes.
- `domain` never imports `store`, `service`, or `api`.
- `store` owns migrations and GORM tags.
- `frontend/src/api/client.ts` owns API URLs.
