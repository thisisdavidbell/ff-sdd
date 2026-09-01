# Feature Specification: Report Datestamp and Dark Mode

**Created**: 2026-09-01
**Status**: Draft
**Input**: User description: "Add a datestamp and dark mode to the HTML report."

## User Scenarios & Testing

### User Story 1 - Verify Report Freshness (Priority: P1)

As a report viewer, I want to see when the report was generated so that I can judge whether its ownership and transfer information is current.

**Why this priority**: The report represents time-sensitive fantasy-football data; a visible generation time is necessary to interpret it correctly.

**Independent Test**: Generate a report and confirm that a human-readable UTC date and time appears directly beneath the report title.

**Acceptance Scenarios**:

1. **Given** a report has been generated, **When** a viewer opens it, **Then** the viewer sees its generation date and time in UTC directly beneath the report title.
2. **Given** the report is generated again at a later time, **When** the viewer opens the new report, **Then** the displayed generation timestamp reflects the later generation.

---

### User Story 2 - Choose a Comfortable Theme (Priority: P2)

As a report viewer, I want to switch between light and dark display modes so that I can comfortably read the report in my environment.

**Why this priority**: The report is reviewed in different lighting conditions, and all report content should remain comfortable to scan in either mode.

**Independent Test**: Open a generated report, change the display mode, and verify that all visible content adopts the chosen mode while remaining readable.

**Acceptance Scenarios**:

1. **Given** a viewer opens the report for the first time, **When** no theme preference has been saved, **Then** the report uses the viewer's device theme preference.
2. **Given** the report is displayed, **When** the viewer selects the alternative display mode, **Then** the page, tables, chart, legend, and expanded manager-change details immediately use that mode.
3. **Given** a viewer selects a display mode, **When** the viewer reloads the report in the same browser, **Then** the selected mode is retained.

### Edge Cases

- A viewer whose device has no declared theme preference sees the light display mode by default.
- The theme control remains visible and operable at narrow viewport widths.
- Text, data visualizations, borders, and interactive hover or expanded states remain distinguishable in both display modes.
- A report generation timestamp must not be confused with capture or player-transfer timestamps already shown in report data.

## Requirements

### Functional Requirements

- **FR-001**: The generated report MUST display one human-readable generation timestamp in UTC directly beneath the "Guy Sports Team Report" title.
- **FR-002**: The generation timestamp MUST identify itself as the report's generation time and MUST include a calendar date, time of day, and UTC designation.
- **FR-003**: Each newly generated report MUST display the time at which that report was generated, rather than reusing a timestamp from a previous report.
- **FR-004**: The report MUST provide a visible control that lets viewers select light or dark display mode.
- **FR-005**: On a viewer's first visit without a saved choice, the report MUST use the viewer's device theme preference; when no device preference is available, it MUST use light mode.
- **FR-006**: The report MUST retain a viewer-selected display mode across reloads in the same browser until the viewer changes it.
- **FR-007**: In both display modes, all report content and interactive states, including tables, the trends chart, legend, manager-change details, and theme control, MUST remain readable and visually distinguishable.
- **FR-008**: This feature MUST NOT alter the report's underlying player ownership, manager change, or historical trend data.

## Success Criteria

### Measurable Outcomes

- **SC-001**: Every newly generated report displays exactly one labeled UTC generation timestamp directly below its title.
- **SC-002**: A viewer can change the report between light and dark modes with one interaction.
- **SC-003**: In both modes, 100% of report sections and interactive states remain available and readable at desktop and narrow mobile viewport widths.
- **SC-004**: A viewer's selected mode is restored after a page reload in the same browser.

## Assumptions

- "Datestamp" means the time the static HTML report is generated, not the time of the latest data capture or manager change.
- The timestamp is presented in UTC to align with the report's existing captured-time values and to avoid ambiguity for viewers in different time zones.
- The report's existing light presentation remains the light mode; dark mode is an alternative presentation rather than a redesign of report content or layout.
- Theme selection is stored only in the viewer's browser and does not require accounts, server-side storage, or changes to captured data.