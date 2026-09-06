---
type: Feature Specification
title: "Development Process Guidance"
description: Historical spec artifact for 003-readme-development-guidance.
tags: [sdd, feature-specification, 003-readme-development-guidance]
status: stable
feature: 003-readme-development-guidance
sdd_approach: streamlined
input_summary: Historical spec artifact.
generated: { by: process:okf-migration, at: 2026-09-06T00:00:00Z }
---

# Feature Specification: Development Process Guidance


## Requirements

- Move the existing Development Process section in `AGENTS.md` verbatim to the repository-root `DEVELOPMENT-PROCESS.md` document.
- Replace the moved section in `AGENTS.md` with a link to the canonical document.
- Replace the README Development prose with a link to the canonical document.
- Preserve the development-process requirements: full SpecKit SDD for major changes, streamlined SDD for smaller work, and documented requirements before implementation.

## Validation

- Confirm `DEVELOPMENT-PROCESS.md` contains the original Development Process section verbatim.
- Confirm `AGENTS.md` and `README.md` both link to `DEVELOPMENT-PROCESS.md`.
- Confirm the old duplicated Development Process headings and prose no longer appear in either entry-point document.