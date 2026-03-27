# Testing

## Test types

- **Unit tests** (`internal/handler/`, `internal/validate/`) — no database required
- **Integration tests** (`internal/store/`) — require a running PostgreSQL instance

## Setup

Integration tests need a test database. With Docker Compose running:

```bash
docker compose exec db psql -U kanban -d kanbanboard -c "CREATE DATABASE kanbanboard_test;"
```

## Running tests

```bash
# All tests (unit + integration)
cd backend && go test ./...

# Unit tests only (no database needed)
cd backend && go test -short ./...

# Verbose output
cd backend && go test -v ./...
```

## Environment variables

| Variable | Default | Purpose |
|----------|---------|---------|
| `TEST_DB_NAME` | `kanbanboard_test` | Test database name |
| `DB_HOST` | `localhost` | Database host |
| `DB_PORT` | `5432` | Database port |
| `DB_USER` | `kanban` | Database user |
| `DB_PASSWORD` | `kanban` | Database password |
