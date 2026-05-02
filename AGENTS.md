# Go5 — Agent Context

## Stack
- Go 1.26.1, Fiber v3, GORM/PostgreSQL, JWT (golang-jwt/jwt/v5)
- Module: `github.com/karelmolina/go5`

## Dev Server
```bash
air              # uses .air.toml (builds ./tmp/main, then runs with APP_ENV=dev)
```
- **Do not** use `go run ./cmd/api` — it skips the env wrapper and build pipeline.

## Required Env Vars
```
JWT_SECRET=<min 32 chars>    # app panics at startup if shorter
PORT=3000
DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME
```
Copy from `.env.example` → `.env`.

## Migrations
No migration CLI. `database.ConnectDB()` runs `AutoMigrate(&model.User{}, &model.Event{}, &model.EventResponse{})` on every startup. Drop/recreate for schema changes.

## Test Commands
```bash
go test ./...                          # all tests
go test -run TestHashAndCheck ./...    # single test
go vet ./...                           # type check (no golangci-lint installed)
```

## API Base
```
http://localhost:3000/api/v1
```

## Routes
| Route | Auth | Notes |
|-------|------|-------|
| `POST /register`, `POST /login` | — | Public |
| `GET /me`, `PATCH /me` | JWT | Any approved user |
| `POST /me/photo` | JWT | Photo upload |
| `GET /users`, `PATCH /users/:id/approve` | JWT + Admin | |
| `POST /events`, `GET /events`, `PATCH /events/:id`, `DELETE /events/:id` | JWT (+Admin for mutating) | |
| `POST /events/:id/responses`, `PATCH /events/:id/responses/me` | JWT | |
| `GET /events/:id/responses` | JWT + Admin | |

## i18n
Language via `Accept-Language` header or `?lang=en|es`.

## README is Stale
README.md references Fiber v2 and older structure. Code uses **Fiber v3** and routes are under `/api/v1`. Trust code over README.

## Key Files
- `cmd/api/main.go` — entry point
- `database/connect.go` — DB connect + AutoMigrate
- `router/router.go` — route definitions
- `internal/utils/jwt.go` — JWT secret management
- `internal/utils/password.go` — bcrypt
- `model/user.go`, `model/event.go` — GORM models
