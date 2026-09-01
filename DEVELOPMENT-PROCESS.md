# Development Process

This project follows spec-driven development (SDD), though with 2 different approaches depending on the task. Both approaches store their requirements and supporting artifacts in `specs/`.

## Full SpecKit SDD

Use the full SpecKit workflow for major new features or substantial changes. Maintain `spec.md`, `plan.md`, `tasks.md`, and any necessary supporting artifacts, such as contracts, data models, and research. Use the SpecKit workflow to specify, plan, task, and implement the work.

## Streamlined SDD

Use streamlined SDD for smaller changes, fixes, or maintenance. Record all implementation-relevant requirements in a single spec file under `specs/` before changing code. The spec must be sufficient to define the desired behavior and validation, but separate plan and task artifacts are not required.

Do not make implementation-only changes: every change must be represented in the applicable SDD artifact before implementation.