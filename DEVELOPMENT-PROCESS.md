# Development Process

This project follows spec-driven development (SDD), though with 2 different approaches depending on the task. Both approaches store their requirements and supporting artifacts in `specs/`.

## Architecture Source of Truth

`ARCHITECTURE.md` is the source of truth for the current repository and application:

- Architecture.
- Functionality and capabilities.
- High-level approach.

All changes MUST be made through one of the SDD processes. A pull request that changes implementation, configuration, automation, or repository behavior MUST include the applicable SDD artifact. It also MUST update `ARCHITECTURE.md` whenever the resulting code behaviour current state differs from that defined in `ARCHITECTURE.md. Neither code-only pull requests nor SDD-only pull requests that leave `ARCHITECTURE.md` out of sync are eligible to merge.

Editorial documentation corrections that do not change the repository's current behavior do not require an `ARCHITECTURE.md` update.

## OKF Specification Bundle

`specs/` is an Open Knowledge Format (OKF) v0.2 bundle. Every non-reserved Markdown file beneath it is an OKF concept and MUST begin with parseable YAML frontmatter containing the repository metadata profile: `type`, `title`, `description`, `tags`, `status`, `feature`, `sdd_approach`, `input_summary`, and `generated`.

`verified` is optional and MUST only record an actual human or automated verification. `sources` is optional and SHOULD identify a primary external standard or source when it materially informs an artifact. `index.md` and `log.md` are reserved OKF files: the root index declares the targeted OKF version and lists features, while the log records date-grouped bundle changes. This requirement is limited to `specs/`; it does not apply frontmatter to unrelated repository Markdown.

## SDD Selection Gate

Select the least complex SDD approach that fully captures the change before creating artifacts.

| Change characteristics | Required approach |
| --- | --- |
| Focused fix, maintenance, configuration, display change, or small behavior change with one validation path | Streamlined SDD |
| New major capability, architectural redesign, external integration, data-model or contract change, or several independently testable delivery slices | Full SpecKit SDD |
| User explicitly selects an approach | The user-selected approach |

Do not infer Full SpecKit SDD merely because a SpecKit skill is available. When there is any doubt or ambiguity about the appropriate approach, the AI Agent MUST ask the user to select Streamlined SDD or Full SpecKit SDD before creating or changing any SDD artifact. Do not default to either approach in that situation.

## SDD Approval Gate

After selecting an SDD approach, an agent MAY create or amend only the applicable SDD artifact so the user can review the proposed requirements, scope, decisions, and acceptance checks. The artifact MUST have `status: draft` while that review is pending.

The agent MUST receive explicit user approval of the applicable SDD artifact before changing implementation, configuration, automation, generated output, `ARCHITECTURE.md`, or any unrelated SDD artifact. On approval, update the artifact to `status: stable` before implementation begins. A user may explicitly authorize a stated exception for the current request; record that exception in the applicable SDD artifact.

Do not treat a request to create, update, inspect, or discuss an SDD artifact as approval to implement the behavior it describes. Do not infer approval from silence, continued conversation, or a request for status.

## Full SpecKit SDD

Use the full SpecKit workflow for major new features or substantial changes. Maintain `spec.md`, `plan.md`, `tasks.md`, and any necessary supporting artifacts, such as contracts, data models, and research. Use the SpecKit workflow to specify, plan, task, and implement the work.

## Streamlined SDD

Use streamlined SDD for smaller changes, fixes, or maintenance. Record all implementation-relevant requirements in a single spec file under `specs/` before changing code. The spec must be sufficient to define the desired behavior and validation, but separate plan and task artifacts are not required.

A streamlined spec MUST record the requested behavior, scope boundaries, affected artifacts, acceptance checks, and assumptions or decisions. It may use concise headings appropriate to the change. Do not create `plan.md`, `tasks.md`, research, data-model, contracts, or generated quality checklists unless the user requests Full SpecKit SDD.

The streamlined `spec.md` is an OKF concept with `type: Feature Specification` and `sdd_approach: streamlined`; record its feature status in frontmatter rather than a presentation-style metadata header.

Create or amend the single streamlined spec for review, then stop for the SDD Approval Gate before making any change outside that spec.

Do not make implementation-only changes: every change must be represented in the applicable SDD artifact before implementation.

## Candidate Work
[TODO.md](TODO.md) contains lists of candidate work for consideration in the future.