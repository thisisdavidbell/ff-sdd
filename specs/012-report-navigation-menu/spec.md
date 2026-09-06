# Feature Specification: Report Navigation Menu

**Feature Branch**: `012-report-navigation-menu`
**Created**: 2026-09-06
**Status**: Implemented
**Input**: User description: "Improve the desktop and mobile report menu, move the dark mode toggle into settings, rename the report to Guy Sports Data, and remove the Top menu item."

## Requested Behavior

The generated Guy Sports report MUST use `Guy Sports Data` as both its browser tab title and reader-facing report title. Its navigation MUST group report destinations by purpose, and its theme control MUST be available from a Settings group rather than the report header.

On desktop, the fixed left navigation MUST contain Statistics links for Player ownership and Team changes, a Trends link for Historical trends, and a Settings section containing the dark-mode toggle. The sidebar MUST not duplicate the report title. The main report header remains the only title and timestamp display.

On mobile, a familiar hamburger icon MUST open and close the navigation menu. It must be closed by default. The compact mobile header MUST show the report title and timestamp beside that control. The duplicate main report header MUST be hidden on mobile. The navigation menu MUST retain the Statistics, Trends, and Settings grouping.

The current `Top` navigation link MUST be removed. No replacement in-page navigation control is required because the persistent desktop title and mobile header establish report context.

## Scope Boundaries

- Preserve the report's three navigable content sections, their IDs, active-section indication, and existing browser-side table, chart, and manager-detail controls.
- Preserve the current persisted light/dark theme behavior and system-theme default; only relocate its control.
- Do not change capture, processing, report data, timestamp formatting, or chart contents.
- Do not add a server, external UI dependency, or additional pages.

## Affected Artifacts

- `internal/render/html.go`
- `tests/unit/test_render_output_test.go`
- `docs/index.html`
- `ARCHITECTURE.md`
- `specs/012-report-navigation-menu/spec.md`

## Acceptance Checks

- Rendering produces `Guy Sports Data` as both the HTML document title and visible report title.
- At desktop widths, the sidebar has Statistics, Trends, and Settings groups; it contains no report title; and the dark-mode toggle occurs only in Settings.
- Statistics contains Player ownership and Team changes, while Trends contains Historical trends.
- At mobile widths, the navigation menu is closed by default. A visible hamburger icon opens and closes it, exposes the same groups and destinations, and updates its accessible expanded state.
- At mobile widths, the header contains one visible report title and one generated timestamp; the main report header is not also visible.
- At desktop widths, the mobile menu is not visible or allocated layout space, including after resizing from an open mobile menu.
- Navigation does not include a `Top` link, and all three report content sections remain navigable.
- The dark-mode control continues to toggle themes and retain the selected theme in local storage.
- The checked-in `docs/index.html` matches the renderer's output.
- `ARCHITECTURE.md` is updated to describes the delivered navigation behavior and title change.
- After the implementation has been validated, the temporary `docs/menu-mockup.html` file is removed before this feature is marked complete.

## Assumptions And Decisions

- The accepted menu mockup is a temporary visual reference and is not part of the delivered report.
- The mobile menu is a full-width panel beneath the fixed compact header.
- `Guy Sports Data` is the intended current product name; no historical data identifiers or directory names change.
- This is a focused display and interaction change with one generated-report validation path, so Streamlined SDD applies.