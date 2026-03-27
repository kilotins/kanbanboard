# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/).

## [Unreleased]

### Added
- Version-planning skill (`docs/skills/version-planning/`) — orchestrates planning for next versions
- Code-health skill (`docs/skills/code-health/`) — codebase assessment for maintainability

### Changed
- Debrief skill updated with code-health references

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
