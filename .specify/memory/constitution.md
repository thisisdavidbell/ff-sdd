<!--
Sync Impact Report:
- Version change: 1.0.0 -> 1.1.0
- Modified principles: I. Spec-First Delivery; II. Data Integrity & Traceability; III. Test-First Validation; IV. Observable, Reproducible Execution; V. Simplicity and Change Control; VI. Historical Data Preservation (added)
- Added sections: Project Constraints, Development Workflow, Historical Data Preservation principle
- Removed sections: none
- Follow-up TODOs: none
-->

# ff-sdd Constitution

This repository exists to collect and present fantasy football data, using upstream
sources including www.guysports.co.uk and dreamteamfc.com. All governance,
implementation, and validation decisions must preserve the reliability and clarity of
that data work.

## Terminology
- **Manager**: a person entered into the Guy Sports competition who selects a team.
- **Player**: a footballer who a manager may select for their team.
- **Team**: a manager's selection of 11 players that follows the game rules and
  roster constraints for the competition.

These definitions apply throughout this project and must be used consistently in
specifications, implementation work, and project documentation.

## Core Principles

### I. Spec-First Delivery
Every feature, bug fix, and data change MUST start with a written specification,
acceptance criteria, and explicit constraints before implementation begins. Any work
without a clear scope and validation target is non-compliant and must be rejected
or re-scoped before proceeding.

This project is data-heavy and decisions are easy to misinterpret when implicit.
A written specification keeps the team aligned on expected behavior, edge cases,
and downstream impacts.

### II. Data Integrity & Traceability
The project MUST treat source data, transformations, and generated outputs as
versioned artifacts with traceable lineage. No silent schema drift, undocumented
mutations, or undocumented assumptions are allowed in data pipelines or reports.

Fantasy football data is noisy, evolving, and often sourced from multiple systems.
Traceability is the mechanism that keeps analyses trustworthy and makes regressions
reproducible.

### III. Test-First Validation
Any behavior change MUST begin with a failing test, reproducible check, or explicit
validation case and must be followed by the minimal implementation needed to make
that validation pass. Data contracts, API changes, and model behavior require
relevant unit and integration checks before merge.

This is non-negotiable because edge cases in data systems often fail silently unless
validation is required before and after a change.

### IV. Observable, Reproducible Execution
Repository workflows, scripts, and operational tasks MUST be runnable in a
reproducible environment with declared dependencies, clear commands, and observable
outputs. Logging, validation evidence, and failure signals are required for data
jobs, analysis runs, and automation.

Reliable execution matters because sports-data work depends on repeatability and
clear debugging paths when outputs diverge from expectations.

### V. Simplicity and Change Control
We prefer the smallest, clearest solution that satisfies the current requirement and
must document trade-offs when extra complexity is justified. Breaking changes must
include explicit review, migration notes, and version-aware handling for consumers.

Simple design reduces maintenance cost and makes it easier to reason about feature
and data changes without hidden complexity accumulating across the project.

### VI. Historical Data Preservation
The project MUST capture data from the original sources at scheduled intervals and
store each collection additively. It MUST preserve historical state, prior versions,
and change over time rather than replacing the existing data with only the latest
snapshot.

The project MUST NOT issue live queries to www.guysports.co.uk or dreamteamfc.com
during results presentation or while handling user queries. Presentation and user
interaction MUST rely on stored data that has already been captured and retained from
those upstream sources.

This requirement is essential because fantasy football data is time-sensitive and
comparative analysis depends on historical context, trend tracking, and replayable
data lineage.

## Project Constraints
The project MUST maintain a structured repository workflow with explicit ownership,
readable task tracking, and clear documentation for all significant changes. This
repository is dedicated to collecting and presenting fantasy football data sourced
from www.guysports.co.uk and dreamteamfc.com. New external dependencies, data
sources, or integration points require justification, review, and documentation
before adoption.

Periodic ingestion MUST fetch data from the original sources and store captured
records additively so that historical state remains available for reporting,
analysis, and reconciliation. The latest record alone is not sufficient; prior
state must remain accessible and queryable.

Secrets, credentials, and sensitive operational details MUST remain out of source
control. The team must preserve auditability and avoid undocumented assumptions in
production-facing workflows.

## Development Workflow
Features, fixes, and experiments MUST follow a defined sequence: specify the
requirement, plan the impact, implement in small increments, validate with relevant
checks, and document contract or behavior changes. Pull requests must be reviewed
for correctness, compatibility, and compliance with this Constitution before merge.

Data capture workflows MUST run on a scheduled, repeatable basis and persist each
snapshot as part of the historical record. Any user-facing query or presentation path
MUST use persisted data rather than live upstream requests.

Changes involving data contracts, schema expectations, or operational automation MUST
include validation evidence and, when necessary, a migration note describing the
transition path for affected consumers.

## Governance
This Constitution governs all repository work and supersedes informal practice when
policies conflict. Amendments require documented rationale, a version bump, and
review by the maintainers before they are enforced.

Compliance is checked during review by verifying that implemented work matches these
principles, that required validation is present, and that any deviation is
explicitly justified and approved. When a conflict arises, this Constitution takes
precedence and the exception must be recorded.

**Version**: 1.1.0 | **Ratified**: 2026-08-29 | **Last Amended**: 2026-08-29
