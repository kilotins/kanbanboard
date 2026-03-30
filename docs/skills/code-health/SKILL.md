---
name: code-health
description: Scan the codebase for maintainability issues, architectural drift, dependency risks, test gaps, and documentation staleness. Use before planning a new version or during a debrief to assess whether the codebase has degraded.
---

# Code Health

Assess the current state of the codebase. The output is a list of findings with concrete recommendations — not a style report or a lecture on principles.

## Philosophy

Use clean code and SOLID principles as **diagnostic tools**, not rules. The question is always "is this causing a real problem?" not "does this violate a principle?"

Good finding: "This handler is 80 lines doing input parsing, authorization, business logic, and response formatting — it's hard to follow and hard to test."

Bad finding: "This violates the Single Responsibility Principle."

Every finding should describe the **problem it causes** (hard to understand, hard to change, hard to test, likely to break) and a **concrete recommendation**. If you can't articulate the problem, it's not a finding.

## When to run

New code is just as likely to have issues as old code. Don't assume recently written code is clean — the pre-release gate exists precisely because fresh code can introduce the same problems that previous versions fixed.

- **Before version planning** — discover technical debt as input to scope decisions. Referenced by `version-planning` Step 1.
- **During a version debrief** — check whether the codebase degraded during development. Referenced by `debrief` Step 4.
- **On demand** — when something feels wrong but you can't pinpoint it.

## Process

### Step 1: Code complexity

Scan for functions and files that have grown unwieldy:

- Functions longer than ~40 lines or doing more than one distinct job
- Files that have become catch-alls (too many unrelated functions)
- Deeply nested logic (multiple levels of if/else, loops within loops)
- Duplicated logic that should be consolidated

**What to look for, not what to count.** A 60-line function that does one clear thing is fine. A 30-line function that mixes three concerns is not.

### Step 2: Architecture adherence

Check whether the code still follows the original architecture decisions:

- Is the separation of concerns holding? (Business logic in handlers only, SQL in store only, no cross-layer leakage)
- Are the package boundaries clean? (No unexpected imports between packages)
- Are naming conventions consistent? (API naming, file organization, function signatures)
- Has any pattern drifted from what was agreed? (e.g. some handlers use a different error handling approach than others)

Compare against the architecture document (`docs/plan/architecture.md`) and CLAUDE.md. Flag divergences.

### Step 3: Error handling and validation

Look for gaps in defensive coding:

- Ignored errors (especially in Go: `result, _ := someFunc()` where the error matters)
- Missing input validation at system boundaries (API handlers, user input)
- Inconsistent error responses (different error formats or HTTP status codes for similar situations)
- Silent failures that should be logged or reported

### Step 4: Dead code

Identify code that is no longer used:

- Unused functions, types, or constants
- Unreachable code paths
- Commented-out code that was never cleaned up
- Exported functions that are never called from outside the package

Dead code is noise — it makes the codebase harder to understand and maintain.

### Step 5: Dependencies

Review external dependencies:

- Are dependencies up to date? Check for known vulnerabilities (`npm audit`, Go vulnerability database)
- Are there unused dependencies in go.mod or package.json?
- Have any dependencies been deprecated or abandoned?
- Is the dependency tree reasonable, or has it grown unexpectedly?

### Step 6: Test coverage

Assess whether the test suite matches the testing strategy:

- Which code paths have no tests?
- Are tests testing behavior (what the code does) or implementation (how it does it)?
- Are there areas of high-risk logic (authorization, data integrity) without test coverage?
- Do existing tests still match the current code? (Tests that pass but test outdated behavior)

Compare against the testing strategy document (`docs/plan/testing-strategy.md`).

### Step 7: Documentation staleness

Check whether documentation reflects the current state:

- Does CLAUDE.md accurately describe the project structure and build commands?
- Do planning documents match what was actually built?
- Are there undocumented API endpoints, features, or configuration options?
- Has the domain model changed without updating the domain model document?

## Output

Present findings grouped by severity:

**High** — actively causing problems or likely to cause bugs. Should be addressed before or during the next version.

**Medium** — making the codebase harder to work with. Should be considered for the next version's scope.

**Low** — minor issues worth noting. Add to backlog for future cleanup.

For each finding:

> **Area:** Code complexity / Architecture / Error handling / Dead code / Dependencies / Tests / Docs
> **Finding:** What the problem is, with specific file and function references
> **Impact:** What problem this causes (hard to test, likely to break, confusing to new contributors)
> **Recommendation:** What to do about it

Findings feed into `docs/plan/backlog.md` or directly into the current version's scope depending on severity and timing.

## What this skill does NOT cover

- **Style and formatting** — use a linter and formatter, enforce in CI
- **Performance** — requires profiling with real data, triggered by user-reported symptoms, not speculative scanning
- **Accessibility** — a feature decision addressed through user stories, not code quality

## Exit criteria

- [ ] Code complexity reviewed — long functions and duplicated logic flagged
- [ ] Architecture adherence checked against planning documents
- [ ] Error handling and validation gaps identified
- [ ] Dead code identified
- [ ] Dependencies reviewed for vulnerabilities and staleness
- [ ] Test coverage compared to testing strategy
- [ ] Documentation checked for staleness
- [ ] Findings presented with severity, impact, and recommendations
