# Pull Request Review Instructions

For every pull request, read `DEVELOPMENT-PROCESS.md` and `ARCHITECTURE.md` before reviewing changed implementation, configuration, automation, SDD artifacts, or repository behavior.

Treat the following as blocking review findings. Use the highest available severity, explain the missing artifact or inconsistency, and cite the relevant changed files:

- The pull request changes implementation, configuration, automation, or repository behavior but does not include the applicable SDD artifact required by the SDD Selection Gate in `DEVELOPMENT-PROCESS.md`. Do not accept code-only pull requests.
- The pull request adds or changes an SDD artifact for a current-state behavior change but does not make the matching `ARCHITECTURE.md` update. Do not accept SDD-only pull requests that leave the architecture document out of sync.
- `ARCHITECTURE.md` conflicts with the behavior introduced or modified by the pull request.
- A changed or added non-reserved Markdown artifact beneath `specs/` lacks the mandatory OKF frontmatter profile required by `DEVELOPMENT-PROCESS.md`, or a changed `specs/index.md` or `specs/log.md` violates its reserved OKF role.
- The pull request implements behavior while its applicable SDD artifact remains `status: draft`, unless that artifact records an explicit user-authorized exception for the current request.

Do not report a missing `ARCHITECTURE.md` update for editorial-only documentation corrections that leave current repository behavior unchanged. Do not request a particular SDD approach without applying the SDD Selection Gate.