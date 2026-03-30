# API Reference

REST API at `/api/v1/`. All responses are JSON with camelCase field names. Authentication via session cookie (`session_token` cookie set on login).

## Authentication

### POST /api/v1/auth/login

Login with email and password. Sets a session cookie.

**Request:** `{ "email": "...", "password": "..." }`

**Response:** User object (see below).

### POST /api/v1/auth/logout

Clears the session cookie. No request body.

### GET /api/v1/auth/me

Returns the currently authenticated user, or 401.

**Response:** User object.

## Setup

### GET /api/v1/setup/status

Check if initial setup is required (no users exist).

**Response:** `{ "setupRequired": true|false }`

### POST /api/v1/setup

Create the initial admin account and set the application title. Only works when no users exist.

**Request:** `{ "name": "...", "email": "...", "password": "...", "appTitle": "..." }`

**Response:** User object (auto-logged in).

### GET /api/v1/app/title

Returns the application title.

**Response:** `{ "title": "..." }`

## Health

### GET /api/v1/health

**Response:** `{ "status": "ok", "database": "ok"|"error" }`

## Users (auth required)

### GET /api/v1/users

List all active users (basic info only: id and name).

**Response:** `[{ "id": "...", "name": "..." }, ...]`

### PUT /api/v1/users/me

Update the current user's profile.

**Request:** `{ "name": "...", "email": "..." }`

**Response:** User object.

### PUT /api/v1/users/me/password

Change the current user's password.

**Request:** `{ "currentPassword": "...", "newPassword": "..." }`

## Admin (admin role required)

### GET /api/v1/admin/users

List all users with full details.

**Response:** Array of User objects.

### POST /api/v1/admin/users

Create a new user account.

**Request:** `{ "name": "...", "email": "...", "password": "...", "isAdmin": false, "isTeamManager": false, "isActive": true }`

**Response:** User object.

### PUT /api/v1/admin/users/{userId}

Update a user's name, email, roles, and active status.

**Request:** `{ "name": "...", "email": "...", "isAdmin": false, "isTeamManager": false, "isActive": true }`

**Response:** User object.

### PUT /api/v1/admin/users/{userId}/password

Reset a user's password (admin operation, no current password required).

**Request:** `{ "password": "..." }`

## Teams (auth required)

### GET /api/v1/teams

List teams owned by the current user.

**Response:** Array of Team objects.

### POST /api/v1/teams

Create a new team. Requires team manager role.

**Request:** `{ "name": "..." }`

**Response:** Team object.

### PUT /api/v1/teams/{teamId}

Rename a team. Owner only.

**Request:** `{ "name": "..." }`

**Response:** Team object.

### DELETE /api/v1/teams/{teamId}

Delete a team. Owner only. Fails if the team owns any projects.

### GET /api/v1/teams/{teamId}/members

List team members. Owner only.

**Response:** Array of User objects.

### POST /api/v1/teams/{teamId}/members

Add a user to a team. Owner only. Idempotent.

**Request:** `{ "userId": "..." }`

### DELETE /api/v1/teams/{teamId}/members/{userId}

Remove a user from a team. Owner only.

## Projects (auth required)

### POST /api/v1/projects

Create a new project with default columns and labels.

**Request:** `{ "name": "...", "tag": "KB", "teamId": "..." }` (teamId optional — omit for personal project)

Tag must be 2-4 uppercase letters, unique across all projects.

**Response:** Project object with columns, labels, tasks, and canEdit flag.

### GET /api/v1/projects

List all projects visible to the current user (owned, team member, or public).

**Response:** Array of Project objects.

### GET /api/v1/projects/{id}

Get a project with its columns, labels, and tasks.

**Response:** `{ ...Project, "columns": [...], "labels": [...], "tasks": [...], "canEdit": true|false }`

### GET /api/v1/projects/{id}/members

List users who can work on a project. For user-owned: just the owner. For team-owned: team owner + members.

**Response:** `[{ "id": "...", "name": "..." }, ...]`

### PUT /api/v1/projects/{id}

Update project name, visibility, or tag. Owner only.

**Request:** `{ "name": "...", "visibility": "public"|"private", "tag": "..." }` (all fields optional)

Tag can only be changed when the project has zero tasks.

**Response:** Project object.

### POST /api/v1/projects/{id}/columns

Add a column at the end of the board.

**Request:** `{ "name": "..." }`

**Response:** Column object.

### PUT /api/v1/projects/{id}/columns/reorder

Reorder columns by providing the full ordered list of column IDs.

**Request:** `{ "columnIds": ["id1", "id2", ...] }`

### PUT /api/v1/projects/{id}/columns/{colId}

Rename a column.

**Request:** `{ "name": "..." }`

**Response:** Column object.

### DELETE /api/v1/projects/{id}/columns/{colId}

Delete a column. Fails if the column contains tasks.

### POST /api/v1/projects/{id}/labels

Add a label to a project.

**Request:** `{ "name": "...", "color": "#808080" }` (color defaults to #808080 if omitted)

**Response:** Label object.

### PUT /api/v1/projects/{id}/labels/{labelId}

Update a label's name and colour.

**Request:** `{ "name": "...", "color": "#..." }`

**Response:** Label object.

### DELETE /api/v1/projects/{id}/labels/{labelId}

Delete a label. Fails if any tasks use the label.

## Tasks (auth required)

### POST /api/v1/projects/{projectId}/tasks

Create a task. Task number is assigned automatically from the project's counter.

**Request:** `{ "title": "...", "columnId": "...", "parentTaskId": "..." }` (parentTaskId optional)

**Response:** Task object (201 Created).

### GET /api/v1/projects/{projectId}/tasks

List all tasks for a project, ordered by column position then task position.

**Response:** Array of Task objects.

### PUT /api/v1/projects/{projectId}/tasks/{taskId}

Update task fields. All fields optional. If columnId changes, the task is moved.

**Request:** `{ "title": "...", "description": "...", "columnId": "...", "labelId": "...", "assigneeId": "...", "priority": "none"|"low"|"medium"|"high", "targetVersion": "...", "dueDate": "YYYY-MM-DD" }`

Pass empty string for labelId, assigneeId, targetVersion, or dueDate to clear the field.

**Response:** Task object.

### PUT /api/v1/projects/{projectId}/tasks/{taskId}/move

Move a task to a specific column and position. Reorders both source and target columns.

**Request:** `{ "columnId": "...", "position": 0 }`

**Response:** Task object.

### DELETE /api/v1/projects/{projectId}/tasks/{taskId}

Delete a task. The task number is not reused.

## Comments (auth required)

### GET /api/v1/projects/{projectId}/tasks/{taskId}/comments

List comments on a task, ordered by creation time.

**Response:** `[{ ...Comment, "authorName": "..." }, ...]`

### POST /api/v1/projects/{projectId}/tasks/{taskId}/comments

Add a comment. Requires edit permission on the project.

**Request:** `{ "text": "..." }`

**Response:** Comment with authorName (201 Created).

### PUT /api/v1/projects/{projectId}/tasks/{taskId}/comments/{commentId}

Update a comment's text. Author only.

**Request:** `{ "text": "..." }`

**Response:** Comment with authorName.

### DELETE /api/v1/projects/{projectId}/tasks/{taskId}/comments/{commentId}

Delete a comment. Author only.

## Data Types

### User

```json
{
  "id": "uuid",
  "name": "string",
  "email": "string",
  "isAdmin": false,
  "isTeamManager": false,
  "isActive": true,
  "createdAt": "2026-03-27T12:00:00Z",
  "updatedAt": "2026-03-27T12:00:00Z"
}
```

### Team

```json
{
  "id": "uuid",
  "name": "string",
  "ownerId": "uuid",
  "createdAt": "2026-03-27T12:00:00Z",
  "updatedAt": "2026-03-27T12:00:00Z"
}
```

### Project

```json
{
  "id": "uuid",
  "name": "string",
  "visibility": "public|private",
  "tag": "KB",
  "nextTaskNumber": 8,
  "ownerUserId": "uuid (or omitted)",
  "ownerTeamId": "uuid (or omitted)",
  "createdAt": "2026-03-27T12:00:00Z",
  "updatedAt": "2026-03-27T12:00:00Z"
}
```

### Column

```json
{
  "id": "uuid",
  "projectId": "uuid",
  "name": "string",
  "position": 0
}
```

### Label

```json
{
  "id": "uuid",
  "projectId": "uuid",
  "name": "string",
  "color": "#hex"
}
```

### Task

```json
{
  "id": "uuid",
  "projectId": "uuid",
  "columnId": "uuid",
  "labelId": "uuid (or omitted)",
  "assigneeId": "uuid (or omitted)",
  "creatorId": "uuid",
  "parentTaskId": "uuid (or omitted)",
  "title": "string",
  "description": "string",
  "priority": "none|low|medium|high",
  "targetVersion": "string (or omitted)",
  "dueDate": "2026-03-27T00:00:00Z (or omitted)",
  "position": 0,
  "taskNumber": 7,
  "createdAt": "2026-03-27T12:00:00Z",
  "updatedAt": "2026-03-27T12:00:00Z"
}
```

### Comment

```json
{
  "id": "uuid",
  "taskId": "uuid",
  "authorId": "uuid",
  "text": "string",
  "createdAt": "2026-03-27T12:00:00Z",
  "updatedAt": "2026-03-27T12:00:00Z"
}
```

## Validation Rules

- **Password:** min 8 characters, at least one letter and one number
- **Project tag:** 2-4 uppercase letters (A-Z), unique, immutable after first task
- **Visibility:** must be "public" or "private"
- **Priority:** must be "none", "low", "medium", or "high"
- **Due date:** YYYY-MM-DD format

## Error Responses

All errors return JSON: `{ "error": "message" }`

Common status codes:
- 400 — invalid input
- 401 — not authenticated
- 403 — not authorized
- 404 — not found
- 409 — conflict (duplicate tag, can't delete column with tasks, etc.)
- 500 — server error
