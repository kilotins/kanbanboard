# Testing

## Test types

- **Unit tests** (`internal/handler/authz_test.go`, `internal/validate/`) — no database required
- **Integration tests** (`internal/store/`, `internal/handler/handler_test.go`) — require a running PostgreSQL instance

## Setup

Integration tests need a test database. With Docker Compose running:

```bash
docker compose exec db psql -U kanban -d kanbanboard -c "CREATE DATABASE kanbanboard_test;"
```

## Running tests

**Important:** Use `-p 1` to run packages sequentially. Multiple test packages share the same test database and will interfere if run in parallel.

```bash
# All tests (unit + integration)
cd backend && go test -p 1 ./...

# Unit tests only (no database needed)
cd backend && go test -short ./...

# Verbose output
cd backend && go test -p 1 -v ./...
```

## Environment variables

| Variable | Default | Purpose |
|----------|---------|---------|
| `TEST_DB_NAME` | `kanbanboard_test` | Test database name |
| `DB_HOST` | `localhost` | Database host |
| `DB_PORT` | `5432` | Database port |
| `DB_USER` | `kanban` | Database user |
| `DB_PASSWORD` | `kanban` | Database password |
