---
type: Feature Specification
title: "Report Navigation"
description: Historical spec artifact for 008-report-navigation.
tags: [sdd, feature-specification, 008-report-navigation]
status: stable
feature: 008-report-navigation
sdd_approach: streamlined
input_summary: Historical spec artifact.
generated: { by: process:okf-migration, at: 2026-09-06T00:00:00Z }
---

# Feature Specification: Report Navigation


## User Scenarios & Testing

### User Story 1 - Navigate the report (Priority: P1)

As a report reader, I want persistent navigation between the main report sections so that I can reach the data I need without manually scrolling through the full report.

**Acceptance Scenarios**:

1. **Given** a reader is viewing the report, **When** they choose Player ownership, Team changes, or Historical trends in navigation, **Then** the selected section is brought into view.
2. **Given** a reader changes report sections by scrolling or navigation, **When** a section is in view, **Then** navigation identifies that current section.

### User Story 2 - Use report navigation on any screen (Priority: P1)

As a report reader, I want report navigation that works comfortably on a narrow screen so that no report content or navigation control becomes inaccessible.

**Acceptance Scenarios**:

1. **Given** a reader views the report on a narrow screen, **When** they open the section navigation, **Then** they can select every report destination without horizontal overflow.
2. **Given** a keyboard-only reader uses navigation, **When** they tab through its controls and follow a link, **Then** focus indication and the destination remain clear.

### Edge Cases

- A section link remains useful when a section has no table rows or has no historical chart data.
- Navigation remains keyboard-operable and indicates the current location.
- Existing interactive manager-change rows and theme selection continue to work.
- The navigation structure can link to a separate report page in future without changing its visual or accessibility model.

## Requirements

### Functional Requirements

- **FR-001**: The report MUST provide navigation to the top of the report, Player ownership, Team changes, and Historical trends, in that order.
- **FR-002**: On desktop viewports, navigation MUST appear as a persistent full-height left sidebar that does not obscure report content.
- **FR-003**: On narrow viewports, navigation MUST become a compact control that reveals the same destinations.
- **FR-004**: Navigation links MUST be ordinary links to stable destinations, allowing a later change from in-page destinations to individual report pages.
- **FR-005**: Navigation MUST identify the section currently in view.
- **FR-006**: Navigation MUST retain the report's existing theme toggle and manager-change detail interaction.
- **FR-007**: Navigation MUST be keyboard-operable and expose meaningful labels to assistive technologies.

## Success Criteria

- **SC-001**: A reader can reach every primary report section from navigation with one action.
- **SC-002**: Navigation renders without horizontal overflow at viewport widths from 320 to 1440 CSS pixels.
- **SC-003**: Navigation identifies the current section after scrolling or following a section link.
- **SC-004**: The report remains usable with keyboard navigation and without JavaScript for following links or revealing narrow-screen navigation.

## Decision Record

- **Goal**: Make the long single-page report quick to navigate while preserving full access to its data.
- **Options considered**: A sticky top menu, a responsive side rail, and collapsible sections with jump links.
- **Selected approach**: A responsive side rail. On desktop it is a full-height left sidebar with visible destinations; on narrow screens it becomes a compact expandable navigation menu.
- **Rationale**: The side rail provides the clearest orientation for a long data report and naturally supports future links to individual report pages. Collapsible sections would hide data rather than improve navigation, and a top menu has less persistent context on desktop.

## Implementation Outcome

- Implemented in `internal/render.BuildHTML` and published in `docs/index.html`.
- The generated report provides a full-height desktop sidebar and a compact native disclosure menu below the desktop breakpoint.
- Navigation order is Top, Player ownership, Team changes, and Historical trends; active-section feedback follows the visible report section.
- Focused renderer tests and browser checks at desktop and narrow viewport widths passed on 2026-09-02.