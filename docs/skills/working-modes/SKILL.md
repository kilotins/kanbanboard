---
name: working-modes
description: When to use plan mode vs conversation mode. Plan mode is for implementation scoping. Conversation mode is for collaborative design, exploration, and decision-making.
---

# Working Modes

Claude Code has two working modes. Using the right mode for the right task prevents friction and produces better results.

## Conversation Mode (default)

Use for **collaborative design and decision-making**:

- Designing skills and workflows
- Eliciting user stories
- Domain modeling and encoding
- Architecture discussions
- UX layout design
- Grilling and challenging ideas
- Debriefs and retrospectives
- Any open-ended exploration where the outcome is decisions or documents

**Why:** These activities are iterative and conversational. The user and Claude go back and forth, challenge assumptions, change direction. Plan mode's structured workflow (explore → plan → approve → execute) would be restrictive here.

**Output:** Documents, decisions, skills — not code.

## Plan Mode

Use for **implementation scoping before writing code**:

- Planning a sub-phase implementation
- Identifying files to create/modify
- Defining verification steps
- Getting approval before coding

**Why:** Code changes need scoping. Plan mode enforces read-only exploration first, prevents premature coding, and ensures the user approves the approach before any files are changed.

**Output:** A concrete plan that leads to code changes.

## The Rule

- If the outcome is **decisions or documents** → conversation mode
- If the outcome is **code changes** → plan mode

## Common Mistake

Trying to use plan mode for design conversations. This leads to:
- Overly structured exploration when free-form discussion would be better
- Pressure to produce a "plan" when the goal is understanding
- Missing important design insights because the workflow is too linear

The planning phases (user stories, domain, architecture, UX, dev workflow) happen in **conversation mode**. The implementation of each sub-phase happens in **plan mode**.
