---
type: Feature Specification
title: "Copilot PR SDD Enforcement"
description: Historical spec artifact for 010-copilot-pr-sdd-enforcement.
tags: [sdd, feature-specification, 010-copilot-pr-sdd-enforcement]
status: stable
feature: 010-copilot-pr-sdd-enforcement
sdd_approach: streamlined
input_summary: Historical spec artifact.
generated: { by: process:okf-migration, at: 2026-09-06T00:00:00Z }
---

# Feature Specification: Copilot PR SDD Enforcement


## Requested Behavior

GitHub Copilot's pull-request review MUST treat the repository development process as a merge requirement. It must identify and raise a blocking review finding when a pull request changes repository behavior or implementation without the applicable SDD artifact, or when the resulting repository state is not reflected in `ARCHITECTURE.md`.

`DEVELOPMENT-PROCESS.md` remains the authoritative definition of the two SDD approaches and their selection gate. Copilot review instructions must refer reviewers to that policy rather than duplicate it.

## Scope Boundaries

- Add repository-level GitHub Copilot instructions for pull-request review.
- Clarify the existing development-process policy so its pull-request requirements are explicit.
- Record the review-enforcement capability in `ARCHITECTURE.md`.
- Do not alter application code, pipeline behavior, GitHub Actions behavior, or existing feature specifications.
- Do not create Full SpecKit planning artifacts; this is a focused configuration and documentation change with one validation path.

## Affected Artifacts

- `.github/copilot-instructions.md`
- `DEVELOPMENT-PROCESS.md`
- `ARCHITECTURE.md`
- `specs/010-copilot-pr-sdd-enforcement/spec.md`

## Acceptance Checks

- `.github/copilot-instructions.md` instructs GitHub Copilot reviewers to read `DEVELOPMENT-PROCESS.md` and `ARCHITECTURE.md` before reviewing relevant changes.
- A PR containing implementation or behavior changes with no applicable SDD artifact is reported as a blocking issue.
- A PR containing a new or changed SDD artifact without a corresponding `ARCHITECTURE.md` update for a current-state change is reported as a blocking issue.
- Documentation-only changes that do not alter current repository behavior are not incorrectly required to change `ARCHITECTURE.md`.
- `DEVELOPMENT-PROCESS.md` is the sole detailed policy definition; other guidance refers to it instead of duplicating it.
- `ARCHITECTURE.md` describes the repository's Copilot review governance capability.

## Assumptions And Decisions

- A review finding marked as blocking/high priority is the available mechanism to discourage merge; Copilot cannot technically prevent a merge by itself.
- The existing SDD Selection Gate determines whether a change requires Streamlined SDD or Full SpecKit SDD.
- `ARCHITECTURE.md` must be updated whenever a PR changes current application or repository behavior, including a newly introduced review control. It is not required for editorial-only documentation corrections that leave current behavior unchanged.