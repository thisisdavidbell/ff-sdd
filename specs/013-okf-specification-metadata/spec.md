---
type: Feature Specification
title: OKF Specification Metadata
description: Adopt Open Knowledge Format frontmatter and navigation for repository SDD artifacts.
tags: [documentation, sdd, okf, governance]
status: stable
feature: 013-okf-specification-metadata
sdd_approach: streamlined
input_summary: "Adopt OKF frontmatter for all Markdown files in specs, with index and log files, and document mandatory usage."
generated: { by: process:okf-migration, at: 2026-09-06T00:00:00Z }
sources:
  - id: okf-specification
    resource: https://github.com/GoogleCloudPlatform/knowledge-catalog/blob/main/okf/SPEC.md
    title: Open Knowledge Format Specification v0.2
---

# Requested Behavior

The `specs/` directory MUST be maintained as an OKF v0.2 knowledge bundle. Every non-reserved Markdown artifact in that directory tree MUST have parseable YAML frontmatter with a non-empty `type` field. Existing SDD body content remains authoritative and is preserved, while existing presentation-style headers are replaced by structured frontmatter.

The bundle root MUST include an `index.md` that declares `okf_version: "0.2"` and lists feature directories for progressive discovery, plus a `log.md` that records this migration. Documentation governing repository development MUST make the OKF requirement explicit.

The SDD process MUST have an explicit approval gate. An agent MAY create or amend the applicable SDD artifact to capture a request, but MUST NOT make implementation, configuration, automation, architecture-document, generated-output, or unrelated SDD-artifact changes until the user explicitly approves that artifact. The user may explicitly authorize a stated exception for the current request.

The metadata profile MUST use `feature` as the single feature identifier. It MUST NOT require or include `feature_branch`, because current branch names duplicate the feature identifier and provide no independent discovery or traceability value.

## Metadata Profile

Every concept in `specs/` MUST include the following frontmatter fields:

- `type`: OKF concept type.
- `title`: reader-facing artifact name.
- `description`: concise artifact summary.
- `tags`: searchable categories.
- `status`: `draft`, `stable`, or `deprecated`.
- `feature`: owning feature identifier.
- `sdd_approach`: `streamlined` or `full-speckit`.
- `input_summary`: concise original request or artifact purpose.
- `generated`: an OKF actor and timestamp for the artifact's latest meaningful content update.

`verified` remains optional. It MUST only be recorded when a named human or automated process has actually verified the artifact; migration alone does not constitute verification. Optional `sources` entries MAY record external standards or primary sources relevant to an artifact.

## Scope Boundaries

- Apply OKF only to the `specs/` bundle; do not require frontmatter in all repository Markdown files.
- Preserve the existing SDD selection gate and its streamlined and Full SpecKit workflows.
- Do not alter capture, processing, rendering, data files, or report behavior.
- Do not claim a human or process verification that did not occur.

## Affected Artifacts

- All existing Markdown artifacts beneath `specs/`.
- `specs/index.md` and `specs/log.md`.
- `DEVELOPMENT-PROCESS.md`.
- `ARCHITECTURE.md`.
- `README.md`.
- `AGENTS.md` and `.github/copilot-instructions.md`.
- A focused OKF bundle conformance test under `tests/unit/`.
- Approval-gate wording in the governing development-process and agent instructions.

## Acceptance Checks

- Every non-reserved Markdown file beneath `specs/` has parseable YAML frontmatter at the beginning of the file and a non-empty `type`.
- Every concept has the required repository metadata profile, and no concept carries `verified` unless its value is structurally valid.
- No concept frontmatter contains the redundant `feature_branch` field, and the conformance test does not require it.
- `specs/index.md` declares OKF version `0.2`, contains no concept metadata other than that declaration, and lists each feature directory.
- `specs/log.md` records the migration using the OKF date-grouped log format.
- Existing single-file streamlined specifications no longer use presentation-style `**Feature Branch**`, `**Created**`, `**Status**`, or `**Input**` header lines.
- `DEVELOPMENT-PROCESS.md`, `ARCHITECTURE.md`, `README.md`, `AGENTS.md`, and `.github/copilot-instructions.md` describe the mandatory `specs/` OKF bundle convention without extending it to unrelated repository Markdown.
- The development process and agent instructions define explicit user approval as the boundary between creating an SDD artifact and implementing its requested changes, including the sole pre-approval exception for that artifact itself.
- The focused conformance test passes with `go test ./tests/unit`.

## Assumptions And Decisions

- The accepted user reference to "OKP" means Open Knowledge Format (OKF) v0.2.
- `generated` records the date of this metadata migration where prior artifact-generation metadata is unavailable; it does not assert original authorship.
- Existing features numbered `001` through `012` are historical completed or draft artifacts. Their body content, rather than an inferred review assertion, establishes their historical context.
- `feature` is the sole feature identifier. It is stable within the bundle and coincides with historical branch names, so `feature_branch` is intentionally omitted.
- This governance and documentation migration has one validation path and uses Streamlined SDD.
- The user explicitly authorized this approval-gate correction while reviewing the current draft, so it is the stated exception for this request.