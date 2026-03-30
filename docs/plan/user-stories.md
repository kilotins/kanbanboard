# User Stories

## User types

One user type with roles:
- Every user can create projects and tasks
- **Team manager** role: can create teams and manage membership of teams you own
- **Administrator** role: can create/manage user accounts, system-level settings
- Roles can overlap (a user can be both admin and team manager)
- No self-registration - admin creates accounts

## Personas

- **Tom, Developer & Admin** — built the app because Jira is too complicated for small teams. Uses it for personal and team projects. Creates accounts and assigns roles.
- **Arne, Scrum Master / Team Owner** — manages a small dev team (Tom, Kåre, Siri). Wants at-a-glance status. Creates team projects.
- **Kåre, Developer** — on Arne's team, also has private projects. Uses private visibility.
- **Siri, Developer** — on Arne's team. Day-to-day user — picks up tasks, subtasks, comments.

## Must have (v1.0)

1. **As a user, I want to create a project with customizable columns** so that I can organize work the way I prefer. *(Default columns: Inbox, Todo, In Progress, Blocked, Done)*
2. **As a user, I want to create tasks and move them between columns** so that I can track progress.
3. **As a user, I want to add subtasks to a task** so that I can break work into smaller pieces. *(Subtasks appear and move independently in columns.)*
4. **As a user, I want to label tasks** (single label per task, from project-scoped labels) so that I can categorize and filter my work. *(Default labels: bug, feature, chore)*
5. **As a user, I want to log in and manage my profile** so that my work is secure and personal.
6. **As a team manager, I want to create teams and manage members** so that my team can collaborate on shared projects.
7. **As a user, I want to control project visibility** - public (everyone can view, only owner edits) or private (only owner views and edits). Owner is a user or team.
8. **As an administrator, I want to create and manage user accounts** so that I control who has access.

## Should have (v1.0)

9. **As an administrator, I want to assign roles** (team manager, administrator) so that I can delegate responsibilities.

## v1.1 Stories

10. **As a user, I want tasks to have a sequential number with a project tag** (e.g. KB-7) so that I can reference tasks in branch names and commits. *(Tag: 2-4 uppercase letters, unique, immutable after first task. Numbers never reused.)*
11. **As Arne, I want to see at a glance who is assigned to each task** so that I know what everyone is working on at standup. *(Assignee initials displayed on task cards.)*
12. **As a user, I want task cards to be visually tinted by their label colour** so that I can quickly scan the board and see the mix of work types.
13. **As a user, I want subtask cards to show their parent task name** so that I can see the relationship when subtasks are in different columns.
14. **As Kåre, I want to see at a glance whether I'm on a team, personal, or private board** so that I don't accidentally create tasks in the wrong project. *(Column backgrounds tinted by board type.)*

## Acceptance criteria (v1.0)

| # | Done when |
|---|-----------|
| 1 | Create project, add/remove/reorder columns. Default columns added on creation. Owner can edit column names. |
| 2 | Tasks on board with title and description. Move between columns (drag and drop). |
| 3 | Subtask linked to parent. Appears in column independently. Moves independently. |
| 4 | Task has single label from project's label set. Can filter board by label. Default labels on project creation. |
| 5 | Login, logout. Sessions persist. User can edit own profile. Password policy: min 8 chars, at least one letter and one number. |
| 6 | Create team, add/remove members. Only team manager who owns the team can manage it. |
| 7 | Default public. Toggle to private. Non-owners see public projects read-only. |
| 8 | Admin creates users with name/email/password. Admin can deactivate users. Password policy enforced on user creation. |
| 9 | Admin assigns/removes roles. Users can have multiple roles. |

## Acceptance criteria (v1.1)

| # | Done when |
|---|-----------|
| 10 | Project has a unique tag (2-4 uppercase letters). Tasks auto-numbered sequentially. Number displayed on cards, detail panel, and API. Tag editable only when project has zero tasks. Deleted task numbers not reused. |
| 11 | Assignee initials shown on task cards when assigned. No indicator when unassigned. |
| 12 | Task card background uses a light tint of the label colour. Task label default colour is cyan (#0891b2). |
| 13 | Subtask cards show parent name above title. ↳ prefix before subtask title. Parent cards show ▤ icon. |
| 14 | Column backgrounds tinted: light blue (personal), green (team), amber (private). |
