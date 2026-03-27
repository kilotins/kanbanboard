---
name: debrief
description: Review what happened, what was learned, and what should change in the skills. Use after planning phases, after a version is complete, or when a mid-development discovery warrants a skill update.
---

# Debrief

Structured retrospective that turns experience into skill improvements. There are three situations that call for a debrief, each with different inputs and focus.

The debrief is not a formality. Its output — skill updates — is the mechanism that makes each project better than the last.

## When to debrief

### 1. Planning debrief (after planning, before implementation)

No code exists yet. The input is the planning documents and the experience of producing them.

**Focus:** Did the planning skills produce good output? Were the phases in the right order? Did each phase have the right inputs from the previous phase?

**Example discoveries:**
- The user-stories skill didn't prompt for personas, so stories lacked perspective
- The ux-layout skill didn't reference the story map, so screen design happened in isolation
- The planning-checklist gates didn't check for non-functional requirements

### 2. Version debrief (after a version is complete)

Code has been written, tested, and accepted. The inputs are the git history, LEARNINGS.md, and the experience of building the version.

**Focus:** What was delivered vs planned? What was learned during development? Did the planning hold up during implementation?

### 3. Mid-development correction

A discovery during implementation reveals that a planning skill is missing something. This is smaller than a full debrief — it's a targeted skill update triggered by a specific problem.

**Focus:** What went wrong, what skill would have prevented it, and what's the minimal update?

**Example discoveries:**
- Validation rules weren't defined during domain encoding, causing ad-hoc decisions during implementation
- Sub-phases were being combined, causing scope creep within individual plans

Mid-development corrections don't need the full debrief process. Jump straight to Step 5 (draft skill updates).

---

## Planning Debrief Process

Capture first, then reflect. The planning debrief has two jobs: make sure nothing from the planning sessions is lost, then evaluate whether the process itself needs improving.

### Step 1: Capture decisions

List all decisions made during planning — architectural choices, technology picks, scope decisions, agreements. For each decision, note the reasoning if it was discussed.

Only capture what was actually discussed — don't invent or assume. Use the user's terminology and framing.

### Step 2: Capture open questions

List anything that was raised but not resolved, needs further research, or was explicitly deferred. These become inputs to early implementation phases or get resolved before coding starts.

### Step 3: Capture constraints

Note any constraints, requirements, or non-negotiables that were identified: deadlines, tech limitations, stakeholder requirements, performance targets. Convert any relative dates to absolute dates.

### Step 4: Capture action items

List concrete next steps with enough detail that they can be picked up in a future session. These typically feed into the implementation plan.

### Step 5: Save to memory

Save each category above as project memories so future sessions have full context. Keep memory entries focused and atomic — one topic per memory file. If something is ambiguous, ask the user before saving.

### Step 6: Review the planning output

Now shift from capturing to reflecting. For each planning phase, ask:

- Did this phase produce a clear, usable document?
- Was the input from the previous phase sufficient?
- Were there questions that should have been asked but weren't?
- Did the phase feel too heavy or too light for this project?

### Step 7: Review the phase flow

Look at the sequence of phases:

- Did any phase feel out of order?
- Were there moments where you needed information from a later phase?
- Were any phases missing? (e.g. "we should have discussed X before moving on")
- Were any phases redundant or overlapping?

### Step 8: What worked well

- Which planning phases produced the most value?
- Were there moments where a good question changed the direction?
- Did the grill-me process surface important issues?

### Step 9: What didn't work

- Where did a skill fail to ask the right questions?
- Where were the gate criteria too loose? (Let something through that caused problems later)
- Were there assumptions that went unchallenged?

Be specific. "User stories could be better" is not actionable. "The user-stories skill didn't prompt for edge cases, so we missed error handling stories" leads to a concrete skill update.

### Step 10: Draft skill updates

For each issue identified, draft a specific skill update:

- Reference the planning experience that motivated the change
- Write the update as guidance for the next project, not as a fix for this one
- Present each update to the user for approval before applying

**Format for presenting updates:**

> **Skill:** `user-stories`
> **Change:** Add "what could go wrong?" pass
> **Motivation:** During planning, we didn't consider edge cases until the UX phase, when it was harder to revise the stories.
> **Update:** New Step 5c that systematically asks about unexpected input, system failures, and boundary conditions for each must-have story.

Apply approved updates to the skill files.

### Step 11: Summary

Present a concise summary of everything captured and all proposed skill updates to the user for review before finalizing.

### Exit criteria (planning debrief)

- [ ] Decisions captured with reasoning
- [ ] Open questions listed
- [ ] Constraints and action items recorded
- [ ] Key items saved to project memory
- [ ] Each planning phase reviewed for output quality
- [ ] Phase flow reviewed — ordering and gaps checked
- [ ] What worked well identified
- [ ] What didn't work identified with specific causes
- [ ] Skill updates drafted, approved, and applied

---

## Version Debrief Process

### Step 1: Review what was delivered

Compare what was planned to what was actually built:

```bash
git log v{previous}..v{current} --oneline
```

- What was delivered as planned?
- What was deferred or dropped? Why?
- What was added that wasn't in the original plan?
- Did the scope change during development, and if so, was that handled well?

Add any deferred items to `docs/plan/backlog.md` so they aren't lost.

### Step 2: Review LEARNINGS.md

Go through each entry added since the last debrief:

- **One-off or pattern?** A one-off gotcha stays in LEARNINGS.md. A pattern (something that would bite you on the next project too) should become a skill update.
- **Was this preventable?** If a skill had mentioned this, would you have avoided the problem? If yes, update the skill.
- **Is the learning still accurate?** Remove or update entries that have been superseded.

### Step 3: What worked well

- Which skills produced the most value?
- Which development practices reduced friction?
- Were there moments where a planning decision saved time during implementation?
- Did any non-obvious approach turn out to be the right call?

These are important to capture — skills tend to accumulate warnings and corrections but lose track of what's already working. If something worked well, note why so you can judge whether it applies in a different context.

### Step 4: What didn't work

- Where did planning fail to prevent a problem?
- Where was the process too heavy or too light?
- Were there repeated mistakes that a skill could have caught?
- Did any phase feel like it was missing or in the wrong order?
- Was testing adequate? Were there categories of bugs that slipped through?
- Has the codebase degraded? Optionally run the `code-health` skill to get a concrete assessment.

Be specific. "Testing could be better" is not actionable. "Backend tests were planned but never written because the dev-workflow didn't enforce them" leads to a concrete skill update.

### Step 5: Draft skill updates

For each issue identified in Steps 2-4, draft a specific skill update:

- Reference the experience that motivated the change
- Write the update as guidance for the next project, not as a fix for this one
- Present each update to the user for approval before applying

**Format for presenting updates:**

> **Skill:** `testing-strategy`
> **Change:** Add enforcement section
> **Motivation:** In v1.0.0, backend tests were planned but never written. The strategy existed but nothing in the workflow enforced it.
> **Update:** "A strategy that isn't enforced is useless. The dev-workflow must include: automated tests must pass before a sub-phase is presented for acceptance."

Apply approved updates to the skill files.

### Step 6: Update CHANGELOG.md

- Rename `[Unreleased]` to `[X.Y.Z] - YYYY-MM-DD`
- Add a new empty `[Unreleased]` section above it
- Commit the changelog update

### Step 7: Capture the debrief

Add a summary to LEARNINGS.md:

```markdown
## YYYY-MM-DD: v{version} debrief

Delivered: {one-line summary of what was built}
Key learning: {the most important thing learned}
Skills updated: {list of skills that were changed and why}
```

This creates a trail that future debriefs can reference.

### Exit criteria (version debrief)

- [ ] Delivered scope compared to plan — deferrals added to backlog
- [ ] LEARNINGS.md entries reviewed — patterns promoted to skill updates
- [ ] What worked well identified and noted
- [ ] What didn't work identified with specific causes
- [ ] Skill updates drafted, approved, and applied
- [ ] CHANGELOG.md updated
- [ ] Debrief summary captured in LEARNINGS.md
