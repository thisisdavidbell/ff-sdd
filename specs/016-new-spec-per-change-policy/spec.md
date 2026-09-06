---
type: Feature Specification
title: New Specification Per Change Policy
description: Require a new standalone SDD specification for each new body of work rather than amending completed historical specifications.
tags: [sdd, process, governance, specifications]
status: stable
feature: 016-new-spec-per-change-policy
sdd_approach: streamlined
input_summary: Ensure new work is documented in a new specification instead of changing an existing specification.
generated: { by: GitHub Copilot, at: 2026-09-06T00:00:00Z }
---

# Requested Behavior

The development process MUST require a new standalone SDD specification for every new body of work. A completed historical specification MUST NOT be amended to describe a later fix, maintenance change, display change, or behavior change.

A new specification MAY reference prior specifications for context, but it MUST contain the implementation-relevant requirements, scope boundaries, affected artifacts, acceptance checks, and decisions for its own work.

## Scope Boundaries

- Update only the documented streamlined SDD process rule.
- Preserve the existing SDD Selection Gate, approval requirement, OKF requirements, and Full SpecKit workflow.
- Do not retroactively modify prior specifications.
- Do not alter application behavior, tooling, automation, or generated report output.

## Affected Artifacts

- `DEVELOPMENT-PROCESS.md`
- `specs/016-new-spec-per-change-policy/spec.md`
- `specs/index.md`
- `specs/log.md`

## Acceptance Checks

- The development process explicitly requires a new standalone specification for new work.
- The development process explicitly prohibits amending a completed historical specification to describe later work.
- The process permits a new specification to reference prior work for context.
- Existing SDD selection, approval, and OKF rules remain present and unchanged in intent.
- The specification bundle index and log include the delivered specification.

## Assumptions And Decisions

- This is a focused documentation and governance change with one validation path, so Streamlined SDD applies.
- A "new body of work" is a distinct request, fix, maintenance task, configuration change, display change, or behavior change that is not merely correcting the specification currently under review before its approval.
- Corrections to a draft specification remain permitted before approval because they refine the same work rather than document a later change.
- The policy applies to future work and does not require historical specification restructuring.
