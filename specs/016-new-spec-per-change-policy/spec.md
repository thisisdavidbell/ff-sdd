---
type: Feature Specification
title: New Specification Per Change Policy
description: Require a new standalone SDD specification for each new body of work while allowing related refinements within active feature work.
tags: [sdd, process, governance, specifications]
status: stable
feature: 016-new-spec-per-change-policy
sdd_approach: streamlined
input_summary: Ensure new work is documented in a new specification instead of changing completed historical specifications, while allowing changes within active feature work.
generated: { by: GitHub Copilot, at: 2026-09-06T00:00:00Z }
---

# Requested Behavior

The development process MUST require a new standalone SDD specification for every new body of work. A completed historical specification MUST NOT be amended to describe a later fix, maintenance change, display change, or behavior change.

A related user-requested refinement MAY be added to an active feature specification before its implementation is complete. An approved `status: stable` specification remains active until its scoped acceptance checks are delivered; its status records approval, not completion.

A new specification MAY reference prior specifications for context, but it MUST contain the implementation-relevant requirements, scope boundaries, affected artifacts, acceptance checks, and decisions for its own work.

## Scope Boundaries

- Update only the documented streamlined SDD process rule.
- Preserve the existing SDD Selection Gate, approval requirement, OKF requirements, and Full SpecKit workflow.
- Do not amend completed historical specifications to describe later work.
- Do not alter application behavior, tooling, automation, or generated report output.

## Affected Artifacts

- `DEVELOPMENT-PROCESS.md`
- `specs/016-new-spec-per-change-policy/spec.md`
- `specs/index.md`
- `specs/log.md`

## Acceptance Checks

- The development process explicitly requires a new standalone specification for new work.
- The development process explicitly prohibits amending a completed historical specification to describe later work.
- The development process permits related user-requested refinements within active feature work.
- The process permits a new specification to reference prior work for context.
- Existing SDD selection, approval, and OKF rules remain present and unchanged in intent.
- The specification bundle index and log include the delivered specification.

## Assumptions And Decisions

- This is a focused documentation and governance change with one validation path, so Streamlined SDD applies.
- A "new body of work" is a distinct request, fix, maintenance task, configuration change, display change, or behavior change that is not a related refinement of active feature work.
- An active feature has approved scope that has not yet completed its acceptance checks; a completed feature is a historical record and may not be reused for later work.
- This clarification applies to future work and does not require historical specification restructuring.
