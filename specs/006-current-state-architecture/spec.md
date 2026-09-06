---
type: Feature Specification
title: "Current-State Architecture Document"
description: Historical spec artifact for 006-current-state-architecture.
tags: [sdd, feature-specification, 006-current-state-architecture]
status: stable
feature: 006-current-state-architecture
sdd_approach: streamlined
input_summary: Historical spec artifact.
generated: { by: process:okf-migration, at: 2026-09-06T00:00:00Z }
---

# Feature Specification: Current-State Architecture Document


## User Scenarios & Testing

### User Story 1 - Understand the application quickly (Priority: P1)

As a human contributor or AI agent, I want one concise document describing the repository's purpose, current capabilities, architecture, and operating approach so that I can orient myself before making changes.

**Acceptance Scenarios**:

1. **Given** a reader opens the architecture document, **When** they scan its opening sections, **Then** they can identify the application's purpose and the major stages of its data pipeline.
2. **Given** a reader needs to trace data, **When** they consult the document, **Then** they can follow data from capture through processing to the rendered report and locate the relevant repository areas.
3. **Given** a reader needs to run or validate the application, **When** they consult the document, **Then** they can find the main entry points, configuration/data conventions, and testing approach.

### User Story 2 - Find authoritative project context (Priority: P1)

As a reader of the README or agent instructions, I want a direct link to the architecture document so that the current-state description is easy to find from both human and AI entry points.

**Acceptance Scenarios**:

1. **Given** a reader opens `README.md`, **When** they look at the project or development guidance, **Then** they can follow a link to `ARCHITECTURE.md`.
2. **Given** an AI agent opens `AGENTS.md`, **When** it looks at repository guidance, **Then** it can follow a link to `ARCHITECTURE.md`.

## Functional Requirements

- **FR-001**: The repository MUST contain a root-level `ARCHITECTURE.md` document.
- **FR-002**: The document MUST describe the repository and application's purpose, current capabilities, high-level architecture, end-to-end data flow, major ownership boundaries, operating entry points, and relevant testing/validation approach.
- **FR-003**: The document MUST describe the current implementation accurately using concise, technology-aware language appropriate for both human contributors and AI agents.
- **FR-004**: The document MUST identify the principal input, intermediate, and output data locations and the role of the durable data format between pipeline stages.
- **FR-005**: The document MUST state meaningful current constraints or assumptions that affect how the application operates.
- **FR-006**: The document MUST describe `config.yaml`, including its role and the `season`, `pages`, and `source_url` settings, and MUST mention the supported configuration overrides.
- **FR-007**: The document MUST use the agreed project terminology consistently: **Manager** means a person entered into the Guy Sports competition who selects a team; **Player** means a footballer a manager may select; and **Team** means a manager's selection of 11 players that follows the competition rules and roster constraints.
- **FR-008**: Where the document presents several related items, it SHOULD use bullets or other scannable structure rather than embedding the items in long sentences.
- **FR-009**: The document MUST use **Team** as the canonical reader-facing term and MUST introduce it on its first relevant reference as **Team (roster)** or equivalent explanatory wording.
- **FR-010**: The terminology convention MUST distinguish reader-facing prose from implementation identifiers, data fields, and historical specifications, which MAY retain `roster` where changing it would add churn or break an existing contract.
- **FR-011**: `README.md` MUST link to `ARCHITECTURE.md`.
- **FR-012**: `AGENTS.md` MUST link to `ARCHITECTURE.md`.
- **FR-013**: The architecture document MUST remain a current-state reference and MUST NOT introduce a speculative roadmap or duplicate detailed feature specifications.
- **FR-014**: `DEVELOPMENT-PROCESS.md` MUST identify `ARCHITECTURE.md` as the source of truth for the repository and application's architecture, functionality, and approach.
- **FR-015**: The SDD processes MUST require updates to `ARCHITECTURE.md` whenever an implementation change alters the repository or application's current architecture, functionality, or approach.
- **FR-016**: The constitution MUST require `ARCHITECTURE.md` to be maintained as the authoritative current-state architecture reference and kept in sync through the SDD processes.

## Success Criteria

- **SC-001**: A new reader can identify the purpose, three pipeline stages, and generated report location from `ARCHITECTURE.md` without consulting source code.
- **SC-002**: A reader can locate the implementation area for capture, processing, rendering, configuration, models, storage, and tests from links or clearly named paths in the document.
- **SC-003**: Both `README.md` and `AGENTS.md` contain a working relative link to the root architecture document.
- **SC-004**: The architecture document is concise enough to serve as a quick orientation reference while covering every subject in FR-002 through FR-005.
- **SC-005**: `DEVELOPMENT-PROCESS.md` and the constitution explicitly identify `ARCHITECTURE.md` as the authoritative current-state reference and require relevant SDD changes to keep it in sync.

## Assumptions

- The existing Go pipeline and repository structure are the source of truth for the current-state description.
- Documentation links use repository-relative Markdown paths.
- No application code or generated report behavior needs to change for this feature.

## Out of Scope

- Changing the capture, processing, rendering, or test implementation.
- Creating a future-state architecture, roadmap, or detailed API contract.
- Replacing the existing development-process or feature-specific specification documents.
