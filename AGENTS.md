# Contributor and AI guidelines

This project is a Go + SQLite + Vue application. Keep the layers small and
predictable. The goal is that a new feature touches the same files in the same
order every time.

## Architecture rules

- `api/` is the HTTP layer only. It decodes requests, validates transport
  details, calls services, and writes JSON/status codes. It must not import
  GORM or run SQL.
- `domain/` contains plain Go types and domain sentinel errors. It must not
  import `api`, `service`, `store`, GORM, or HTTP packages.
- `service/` owns business rules, validation, orchestration, and transactions.
  It may import `domain` and `store`. It must not write HTTP responses.
- `store/` owns GORM models, table names, SQLite access, and migrations. Only
  this package imports GORM or the SQLite driver.
- `frontend/src/api/client.ts` is the only frontend module that knows API URLs.
- Keep `make dev` working at all times.

Dependency direction:

```text
api -> service -> store -> SQLite
        |
        v
      domain
```

## Adding a vertical slice

Use this order so the code stays compiling and the layering stays intact.

1. **Domain**: add plain types in `domain/`. Add sentinel errors there only if
   they are domain-wide.
2. **Store**: add the GORM model in `store/`, set `TableName()`, and add a new
   numbered migration in `store/migrations/`. Add conversion helpers between
   the store model and the domain type.
3. **Service**: add the use case in `service/`. Validate input here. Return
   domain or service errors, never HTTP errors.
4. **API**: add a handler file in `api/`, register routes in `api/server.go`,
   and map errors to status codes in `writeError`.
5. **Frontend**: add types in `src/types.ts` (or the relevant client module),
   API calls in `src/api/client.ts`, a component in `src/components/`, and wire
   it from `src/router.ts` or an existing view.
6. **Docs**: update `README.md` and `docs/architecture.md` if the feature or
   folder responsibilities changed.

## Migrations

- Add new migrations as `store/migrations/000NN_description.sql`.
- Include both `-- +goose Up` and `-- +goose Down`.
- Never edit a migration that has already been applied.
- SQLite is opened with WAL mode, foreign keys, and a busy timeout. Keep writes
  serialized; do not raise the connection pool without revisiting that choice.

## Verification

Run these before finishing a change:

```sh
make tidy
make test
make build-web
make build
```

Run `make dev` when the change affects startup, routing, the Vite proxy, or the
embedded frontend. The API listens on `:8080` and Vite on `:5173` by default.

## Conventions

- Use `snake_case` for JSON fields and types that cross the API boundary.
- Keep error messages lowercase and wrap errors with context.
- Do not add a business rule to `api/` just because it is convenient there.
- Do not put GORM tags on `domain/` types.
- Prefer small files named after the feature, for example `widget_service.go`,
  `widget.go`, and `widgets.go`.
- Update this file if the architecture changes.
