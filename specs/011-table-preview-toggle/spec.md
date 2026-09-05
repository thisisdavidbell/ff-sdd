# Feature Specification: Table Preview Toggle

**Feature Branch**: `011-table-preview-toggle`
**Created**: 2026-09-05
**Status**: Implemented
**Input**: User description: "For both tables, display the initial 10 lines, and allow the user to click to expand and see all other lines with 1 click. There should also be a button to hide most of the teams again so its reversible."

## Requested Behavior

The Player ownership and Team changes tables in the generated report MUST initially display their first ten data rows. When a table contains additional rows, a reader can use one button to reveal all remaining rows and use that same button again to return to the ten-row preview.

## Scope Boundaries

- Apply the preview and toggle only to the Player ownership and Team changes tables.
- Keep all report data in the static HTML; do not change the data pipeline or add a server-side pagination mechanism.
- Preserve the existing theme, report navigation, and expandable manager-change details.

## Affected Artifacts

- `internal/render/html.go`
- `tests/unit/test_render_output_test.go`
- `docs/index.html`
- `ARCHITECTURE.md`
- `specs/011-table-preview-toggle/spec.md`

## Acceptance Checks

- A generated report with more than ten player rows initially displays ten Player ownership rows and offers a control labeled with the total number of players.
- A generated report with more than ten manager rows initially displays ten Team changes rows and offers a control labeled with the total number of teams.
- Activating either control reveals all rows and changes its label to `Show fewer players` or `Show fewer teams`.
- Activating the control again restores the ten-row preview and changes its label back to the applicable `Show all` label.
- A table with ten or fewer data rows does not show an unnecessary control.
- Existing manager-change details continue to be excluded from the ten-row count and remain independently expandable.

## Assumptions And Decisions

- The first ten rows use the existing sort order: ownership count descending for players and newest change first for managers.
- The controls are browser-side progressive enhancement because the report is a self-contained static page; without JavaScript, all rows remain readable.
- The existing generated `docs/index.html` is refreshed through the renderer so it matches the implementation.