# Development Process

This project follows spec-driven development (SDD), though with 2 different approaches depending on the task. Both approaches store their requirements and supporting artifacts in `specs/`.

## Architecture Source of Truth

`ARCHITECTURE.md` is the source of truth for the current repository and application:

- Architecture.
- Functionality and capabilities.
- High-level approach.

All changes MUST be made through one of the SDD processes and MUST update `ARCHITECTURE.md` to keep it in sync with the resulting current state.

## SDD Selection Gate

Select the least complex SDD approach that fully captures the change before creating artifacts.

| Change characteristics | Required approach |
| --- | --- |
| Focused fix, maintenance, configuration, display change, or small behavior change with one validation path | Streamlined SDD |
| New major capability, architectural redesign, external integration, data-model or contract change, or several independently testable delivery slices | Full SpecKit SDD |
| User explicitly selects an approach | The user-selected approach |

Do not infer Full SpecKit SDD merely because a SpecKit skill is available. When there is any doubt or ambiguity about the appropriate approach, the AI Agent MUST ask the user to select Streamlined SDD or Full SpecKit SDD before creating or changing any SDD artifact. Do not default to either approach in that situation.

## Full SpecKit SDD

Use the full SpecKit workflow for major new features or substantial changes. Maintain `spec.md`, `plan.md`, `tasks.md`, and any necessary supporting artifacts, such as contracts, data models, and research. Use the SpecKit workflow to specify, plan, task, and implement the work.

## Streamlined SDD

Use streamlined SDD for smaller changes, fixes, or maintenance. Record all implementation-relevant requirements in a single spec file under `specs/` before changing code. The spec must be sufficient to define the desired behavior and validation, but separate plan and task artifacts are not required.

A streamlined spec MUST record the requested behavior, scope boundaries, affected artifacts, acceptance checks, and assumptions or decisions. It may use concise headings appropriate to the change. Do not create `plan.md`, `tasks.md`, research, data-model, contracts, or generated quality checklists unless the user requests Full SpecKit SDD.

Do not make implementation-only changes: every change must be represented in the applicable SDD artifact before implementation.

## Candidate Work
[TODO.md](TODO.md) contains lists of candidate work for consideration in the future.