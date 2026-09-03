# Issue 2 — Make TUI a control surface for AW

## Decision
Keep `graph/` temporarily, but move the canonical organization graph to `organization.json` + `organization.mmd`.

## Goal
Make `aw.tui` / `soubi.tui` useful because it can trigger a real AW workflow, while keeping all orchestration logic outside the TUI.

## First executable loop

```text
TUI RUN
  ↓
AW
  ↓
latest 30 repos
  ↓
repo.yaml discovery
  ↓
organization.json + organization.mmd
  ↓
TUI refresh
```

## TUI events
- `refresh`
- `inspect`
- `plan`
- `run`
- `approve`

## Boundaries
- TUI: View/Input/Event only
- AW: workflow orchestration
- Coordinator Agent: decisions
- repo.yaml: repo capability plug
- pi agent: execution runtime
- JSON: machine state
- MMD: graph visualization

## Acceptance criteria
- A TUI action can trigger one real AW workflow.
- The workflow refreshes the organization state.
- The resulting JSON/MMD can be displayed by the TUI.
- No agent assignment/orchestration policy is embedded in Go UI code.
- The old Go graph algorithm is removable after the new loop is validated.
