# Data Model: Automated Data Pipeline & Capture Source Separation

**Feature**: [002-automated-data-pipeline](../spec.md)
**Date**: 2026-08-30

## Artifact Lineage & Persistence

This feature orchestrates and automates the flow of data through existing storage entities. No schema changes are made to existing YAML files; however, the pipeline enforces strict additive state transitions across the storage lineage.

```
+--------------------------+
| Upstream (guysports.co.uk)|
+--------------------------+
             |
             | capture-guysports
             v
+--------------------------+
|  Raw Data Snapshots      |  (Additive YAML in data/{season}/raw/)
+--------------------------+
             |
             | process
             v
+--------------------------+
|  Processed State         |  (Ownership & Manager Changes YAML in data/{season}/processed/)
+--------------------------+
             |
             | render
             v
+--------------------------+
|  Static HTML Report      |  (Rendered report in docs/index.html)
+--------------------------+
             |
             | git commit & push (CI)
             v
+--------------------------+
|  Version Control Branch  |  (Persisted repository state)
+--------------------------+
```

## Pipeline Execution Artifacts

### 1. Raw Snapshot Entity (Existing, Additive)
- **Path**: `data/{season}/raw/{TeamName}_{ManagerID}/{Timestamp}.yaml`
- **Mutability**: Append-only (new snapshots created per run if changes detected; unchanged snapshots skipped).
- **Source**: `capture-guysports` CLI target.

### 2. Processed Derived Entities (Existing, Updated)
- **Paths**:
  - `data/{season}/processed/player-ownership.yaml`
  - `data/{season}/processed/manager-changes/{TeamName}_{ManagerID}.yaml`
- **Mutability**: Overwritten/updated deterministically based on raw snapshot lineage.
- **Source**: `process` CLI target.

### 3. Rendered Presentation Entity (Existing, Updated)
- **Path**: `docs/index.html`
- **Mutability**: Overwritten deterministically based on processed state.
- **Source**: `render` CLI target.

### 4. Git Commit Artifact (New Automated Behavior)
- **Scope**: Includes modified/added files in `data/` and `docs/`.
- **Commit Message Format**: `chore(data): automated weekly pipeline run [skip ci]` or equivalent.
- **Trigger**: GitHub Actions scheduled or manual workflow completion.
