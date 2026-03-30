# Domain Model

## Entities (7)

| Entity | Fields |
|---|---|
| **User** | name, email, credentials, roles (admin, team manager), state (active/inactive/deleted) |
| **Team** | name, owner (user with team manager role), members (users) |
| **Project** | name, owner (user or team), visibility (public/private), tag (unique, 2-4 uppercase letters), next task number (counter) |
| **Column** | name, position (within project) |
| **Task** | title, description, column, label (single), assignee (user), creator (user), parent task (optional), target version, priority, due date, task number (sequential within project) |
| **Label** | name, color (within project) |
| **Comment** | text, author (user), timestamp (on task) |

## Defaults on project creation

- **Columns:** Inbox, Todo, In Progress, Blocked, Done
- **Labels:** task (cyan), bug (red), feature (green), chore (grey)

## Validation rules

### Password policy
- Minimum 8 characters
- At least one letter (uppercase or lowercase)
- At least one number
- Special characters allowed

### Project tag
- 2-4 uppercase letters only (A-Z)
- Unique across all projects
- Required at project creation
- Immutable once the project has tasks

### Other input validation
- Email: valid email format, unique per user
- User name: required, non-empty
- Project name: required, non-empty
- Project tag: required, 2-4 uppercase letters, unique
- Task title: required, non-empty
- Column name: required, non-empty
- Label name: required, non-empty
- Priority: one of 'none', 'low', 'medium', 'high'
- Visibility: one of 'public', 'private'

## Key design decisions

- All work items are Tasks - no separate Bug/Feature/Subtask classes
- Subtasks are Tasks with a parent reference, move independently in columns
- Single label per task (not multiple)
- Labels are project-scoped - same text in different projects are independent
- Priority is a field on Task, not a label
- Columns must be defined before tasks are added
- Task assignee defaults to owner for personal projects, unassigned for team projects
- Creator and assignee are separate fields
- Task numbers are sequential per project, assigned atomically, never reused
- Project tag is editable only while the project has zero tasks
- User deletion is soft — record preserved with name for historical references
- Three user states: active (can log in), inactive (reversible, cannot log in), deleted (permanent, cannot log in)
- Deleting a user cascades: owned projects deleted, teams transferred, tasks unassigned

## Napkin diagram

```
User ──belongs to──▶ Team
 │                    │
 owns                 owns
 ▼                    ▼
Project ──has──▶ Column ──has──▶ Task ──parent──▶ Task
 │                                │
 has                              has
 ▼                                ▼
Label ◀──tagged on──────────── Comment
```
