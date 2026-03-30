# Learnings

## 2026-03-30: Shared test database requires sequential package execution

When multiple Go test packages (handler and store) share the same PostgreSQL test database, running packages in parallel causes FK constraint violations and data corruption. Each package's `cleanTables` call conflicts with the other package's active tests.

**Solution:** Run tests with `go test -p 1 ./...` to execute packages sequentially. The `-p 1` flag limits parallelism to one package at a time. Within a package, tests still run sequentially (no `t.Parallel()` used).

**Alternative not taken:** Separate test databases per package. More isolated but adds setup complexity.

## 2026-03-30: v1.1.2 debrief

Delivered: Code health fixes — authorization on column/label handlers, priority validation, transactional task creation, writeJSON helper, duplicate email detection, handler tests, doc corrections.
Key learning: Security bugs (missing authorization) existed since v1.0.0 and weren't caught until the first code-health scan. Run code-health before every release, not just when you think of it.
Skills updated: planning-checklist (pre-release gate with code-health scan).

## 2026-03-30: v1.1.0 debrief

Delivered: Task numbering with project tags, label-tinted cards, assignee initials, parent/subtask indicators, board context icons, API documentation.
Key learning: Multiple visual changes to the same area can clash — column background tinting combined with label-tinted cards was too visually heavy. Test visual interactions before committing to multiple colour layers.
Skills updated: ux-layout (check for feature removal during redesign), version-planning (visual interaction check in impact analysis).

## 2026-03-27: v1.0.1 debrief

Delivered: Backend test suite (62 tests — unit + integration), changelog, debrief skill.
Key learning: Extracting pure functions from DB-coupled handlers is an effective way to get testable authorization logic without mocks or interfaces.
Skills updated: domain-encoding (validation exit criteria), dev-workflow (changelog), planning-checklist (changelog + debrief reference).

Discoveries, gotchas, and workarounds found during development.

## 2026-03-27: svelte-dnd-action + Svelte 5

`svelte-dnd-action` has known compatibility issues with Svelte 5's `$state` proxy objects. Drag-and-drop silently breaks — items become invisible or the library's internal state fights with Svelte's reactivity.

**Solution:** Switched to `@thisux/sveltednd` which is built specifically for Svelte 5 runes. Simpler API (draggable + droppable actions) and works correctly with `$state`.

**Lesson:** Always verify library compatibility with the specific framework version before committing. Check GitHub issues for the library + framework version combination.

## 2026-03-27: PostgreSQL 18 volume mount path

PostgreSQL 18 changed the Docker image to use version-specific subdirectories under `/var/lib/postgresql/`. The old volume mount at `/var/lib/postgresql/data` causes a startup error.

**Solution:** Mount the volume at `/var/lib/postgresql` instead of `/var/lib/postgresql/data`.

## 2026-03-27: Click vs drag detection

When using a drag-and-drop library, click events on draggable items are unreliable. The drag library captures pointer events, preventing normal click handling.

Tried: timing-based (200ms threshold) — unreliable, small mouse movements during click triggered drag state.
Tried: pointermove-based — any tiny movement set isDragging=true.

**Solution:** Distance-based detection. Track pointer position on `pointerdown`, compare with position on `click`. If delta < 5px, it's a click; otherwise it's a drag.

## 2026-03-27: UNIQUE constraint during column reorder

The `columns` table has a UNIQUE constraint on `(project_id, position)`. When swapping column positions in a loop, two columns temporarily have the same position, violating the constraint.

**Solution:** Two-pass approach within a transaction: first set all positions to negative values (-(i+1)), then set to final values. Negative positions never conflict with the UNIQUE constraint.

## 2026-03-27: PostgreSQL UUID columns and "not found" tests

When testing "not found" paths for tables with UUID primary keys, passing a plain string like `"nonexistent-id"` causes a PostgreSQL type error (`invalid input syntax for type uuid`) instead of returning `sql.ErrNoRows`. The store functions then return a wrapped error instead of the expected sentinel (e.g. `ErrTaskNotFound`).

**Solution:** Use a valid UUID format that doesn't exist: `"00000000-0000-0000-0000-000000000000"`. This produces the expected `ErrNoRows` → sentinel error path. This applies to any function that takes a UUID parameter — including "exclude" parameters where you'd pass empty string to mean "don't exclude anything".

## 2026-03-27: Svelte 5 $state with prop initial values

Svelte 5 warns when using `$state(prop.value)` because it only captures the initial value. The prop is reactive but the `$state` isn't synced automatically.

**Solution:** Use `$effect` to resync the local state when the prop changes. This is intentional in Svelte 5 — local mutable copies of props need explicit sync.
