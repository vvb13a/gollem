# gollem

A Go + SQLite + Vue application.

## Stack

- Go (net/http, no framework)
- SQLite via glebarez/sqlite + GORM (pure Go, no CGO)
- Goose migrations embedded in the binary
- Vue 3 + TypeScript
- Vite
- PrimeVue 5
- Tailwind CSS 4

## Development

```sh
make install
make dev
```

- Vite UI: http://localhost:5173
- Go API: http://localhost:8080
- Health check: http://localhost:8080/api/health

`make dev` runs both processes and stops them together on Ctrl-C. It installs
frontend dependencies automatically if `frontend/node_modules` is missing.
Development data defaults to `./.dev/data`; override it with `APP_DATA_DIR`.

## Project structure

See [docs/architecture.md](docs/architecture.md).

## AI / contributor guidance

See [AGENTS.md](AGENTS.md).

## Module path

This project uses `github.com/vvb13a/gollem`. Update `go.mod` and Go imports if you rename it.

## Make targets

| Target | Description |
| --- | --- |
| `make` | List targets |
| `make install` | Install frontend dependencies |
| `make dev` | Run the Go API and Vite together |
| `make api` | Run the Go server only |
| `make web` | Run Vite only |
| `make build-web` | Build the Vue app into `frontend/dist` |
| `make build` | Build frontend + `bin/app` |
| `make release` | Cross-compile a static Linux binary |
| `make test` | Run Go tests |
| `make tidy` | Tidy Go modules |
| `make clean` | Remove build artifacts |

## Docker

Docker is optional. The app uses a pure-Go SQLite driver, so a single-stage
runtime image can run the compiled binary without CGO or system SQLite.
