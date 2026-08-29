# Command Contracts

## Capture command

Command: ./bin/capture --source guysports

Behavior:
- Fetches the current page set from the configured Guy Sports season table pages.
- Stores timestamped manager snapshots in YAML.
- Appends new records instead of replacing existing raw snapshots.
- Returns a summary of managers captured and files written.

Failure modes:
- Upstream page unavailable: log the error and exit non-zero.
- Partial page failure: retain successful snapshots and record the failed source page.

## Process command

Command: ./bin/process

Behavior:
- Reads the raw YAML snapshots.
- Produces ownership aggregates and manager change summaries.
- Writes derived YAML outputs under `data/<season>/processed/`.
- Does not mutate raw capture history.

Failure modes:
- Empty or malformed input: warn, skip invalid records, and continue where possible.
- Missing capture data: exit non-zero with clear diagnostics.

## Render command

Command: ./bin/render

Behavior:
- Reads the latest processed aggregate data currently present in the local repository state, whether from the latest committed data or the latest local processing run.
- Produces static HTML pages for local viewing and GitHub Pages hosting.
- Writes to docs/.
- Uses embedded templates and static assets with no live API dependency.

Failure modes:
- Missing processed data: fail clearly and instruct the user to run capture/process first.
- Template generation error: exit non-zero and preserve the previous good render.

## Reset/local cleanup policy

Local test runs may discard uncommitted generated working files, but must not modify the raw capture archive unless explicitly approved as a deliberate cleanup action.
