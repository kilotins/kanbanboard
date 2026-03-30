# Implementation Plan

## Phase 1: Foundation (v0.1.x)

| Sub-phase | Description | Version |
|---|---|---|
| 1.1 | Project setup - Go backend, Svelte frontend, PostgreSQL, Docker Compose, project structure | v0.1.1 |
| 1.2 | Database schema and migrations (all entities, hand-rolled runner) | v0.1.2 |
| 1.3 | Onboarding - first-time setup screen (admin account, application title) | v0.1.3 |
| 1.4 | Authentication - login, logout, sessions | v0.1.4 |

## Phase 2: Core Board (v0.2.x)

| Sub-phase | Description | Version |
|---|---|---|
| 2.1 | Project CRUD - create project with default columns and labels | v0.2.1 |
| 2.2 | Board view - render columns and task cards | v0.2.2 |
| 2.3 | Task CRUD - create, edit, move tasks. Side panel for detail. | v0.2.3 |
| 2.4 | Drag and drop between columns | v0.2.4 |
| 2.5 | Labels - assign to task, filter board by label | v0.2.5 |

## Phase 3: Customization (v0.3.x)

| Sub-phase | Description | Version |
|---|---|---|
| 3.1 | Profile editing | v0.3.1 |
| 3.2 | Project settings - edit columns, labels, visibility | v0.3.2 |

## Phase 4: Subtasks and Comments (v0.4.x)

| Sub-phase | Description | Version |
|---|---|---|
| 4.1 | Subtasks - create, show independently in columns with indicator | v0.4.1 |
| 4.2 | Subtask progress on parent card, warning on parent to Done | v0.4.2 |
| 4.3 | Comments - add, edit, delete own | v0.4.3 |

## Phase 5: Teams and Collaboration (v0.5.x)

| Sub-phase | Description | Version |
|---|---|---|
| 5.1 | Admin area - user management (create, edit, deactivate, assign roles) | v0.5.1 |
| 5.2 | Team management - create teams, add/remove members | v0.5.2 |
| 5.3 | Team project ownership - all members can edit | v0.5.3 |
| 5.4 | Project visibility - public/private, read-only for non-owners | v0.5.4 |

## Phase 6: Documentation and Release (v0.6.x)

| Sub-phase | Description | Version |
|---|---|---|
| 6.1 | User guide | v0.6.1 |
| 6.2 | README and LICENSE | v0.6.2 |
| 6.3 | Final review | v0.6.3 |
| - | Release | v1.0.0 |

## v1.0.1: Backend Tests and Skills

| Description | Version |
|---|---|
| Unit tests (handler authorization), integration tests (store layer), changelog, debrief skill | v1.0.1 |

## v1.1: Task Numbering and Visual Improvements

| Phase | Description | Version |
|---|---|---|
| 1 | Task numbering — backend, migration (add tag/number fields, backfill), API, tests | v1.1-snapshot-1 |
| 2 | Card visual improvements — label tinting, task number display, parent/subtask indicators, assignee initials | v1.1-snapshot-2 |
| 3 | Board context colours — column backgrounds by project type (blue/green/amber) | v1.1-snapshot-3 |
| 4 | API documentation + release | v1.1.0 |
