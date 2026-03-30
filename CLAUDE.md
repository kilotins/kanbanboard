# Kanban Board - Claude Code Navigation

## Project Overview

Lightweight kanban board for individuals and small teams. Go backend + Svelte frontend + PostgreSQL.

## Build & Run

```bash
# Docker (full stack)
docker compose up --build

# Dev mode - backend
cd backend && go run ./cmd/server

# Dev mode - frontend
cd frontend && npm run dev

# Run all backend tests (requires test database)
cd backend && go test -p 1 ./...

# Run unit tests only (no database needed)
cd backend && go test -short ./...
```

## Project Structure

- `backend/cmd/server/` - Go entry point
- `backend/internal/handler/` - REST API handlers
- `backend/internal/middleware/` - Auth middleware (RequireAuth, RequireAdmin)
- `backend/internal/model/` - Domain entities
- `backend/internal/store/` - Database access (PostgreSQL, plain SQL)
- `backend/internal/validate/` - Input validation (password policy)
- `backend/migrations/` - SQL migration files
- `frontend/src/lib/` - Svelte components

## API

REST API documented in `docs/api.md`.

## Key Decisions

- No ORM - standard `database/sql` with pgx driver
- REST API at `/api/v1/`
- JSON responses, camelCase naming
- Session-based auth with cookies
- Hand-rolled migration runner

## Planning & Skills

- `docs/plan/` - Planning documents (user stories, domain model, architecture, UX, dev workflow, testing, implementation plan)
- `docs/skills/` - Planning skills (reusable across projects)

## Versioning

`v0.{phase}.{subphase}` during development. `v1.0.0` at completion. Post-release: `v{major}.{minor}-snapshot-{N}` for phased work, `v{major}.{minor}.0` for release. Changes tracked in CHANGELOG.md.
