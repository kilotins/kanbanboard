# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/).

## [1.1.0] - 2026-03-30

### Added
- Task numbering: projects have a unique tag (2-4 uppercase letters), tasks are auto-numbered sequentially (e.g. KB-7). Numbers are never reused.
- Label-tinted task cards: card backgrounds use a light tint of the label colour
- Task number displayed on cards and in the task detail panel
- Parent task indicator (▤) on cards with subtasks
- Subtask cards show parent name above title with ↳ prefix
- Assignee initials displayed on task cards when assigned
- Board context icons in header: 👤 personal, 👥 team, 🔒 private
- Subtle board background tinting by project type
- Tag input with auto-suggest on project creation
- Tag display in project settings (locked after tasks exist)
- API documentation (`docs/api.md`)
- Version-planning skill (`docs/skills/version-planning/`)
- Code-health skill (`docs/skills/code-health/`)
- Personas added to user stories document

### Changed
- Default "task" label colour changed from blue (#4a90d9) to cyan (#0891b2)
- Debrief skill updated with code-health references
- Dev-workflow skill updated with snapshot versioning for post-release development
- Planning documents updated with v1.1 stories and implementation plan

## [1.0.1] - 2026-03-27

### Added
- Unit tests for handler authorization logic (visibility, edit permission, ownership)
- Integration tests for store layer against PostgreSQL (users, sessions, teams, projects, tasks, comments)
- Test infrastructure: test database helpers, seed functions, `-short` flag support
- `backend/TESTING.md` with setup and run instructions
- Debrief skill (`docs/skills/debrief/`)
- This changelog

### Changed
- Extracted authorization decisions into pure functions (`authz.go`) for testability
- Handler functions delegate to pure authorization functions via `resolveTeamContext`
- Updated dev-workflow and planning-checklist skills with changelog and debrief references
- Added validation rules to domain-encoding exit criteria

## [1.0.0] - 2026-03-27

Initial release. Kanban board for individuals and small teams with projects, tasks, subtasks, comments, teams, labels, drag-and-drop, and role-based access control.
