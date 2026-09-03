# 🌱 aw.tui

GitHub repositories as organizational data. `aw.tui` observes the **30 most recently updated bonsai repositories**, builds a similarity graph, and discovers communities without hard-coded departments.

## Model

```text
GitHub repos (latest 30)
        ↓
   repo features
 name / description / language / topics
        ↓
     TF-style vectors
        ↓
 cosine similarity + topic affinity
        ↓
   similarity graph
        ↓
 label propagation
        ↓
 self-organized communities
```

The `repos/` package is the repository-data submodule. The graph is deliberately generated from observed repository relationships; department names are not predefined.

## Run

```bash
go run ./cmd/aw-tui -owner bonsai -out graph.json
```

For private repositories, provide a GitHub token through `GITHUB_TOKEN`.

## Output

`graph.json` contains:

- `nodes`: repositories and discovered community IDs
- `edges`: weighted repository similarity links

This is the first layer of the AW orchestra:

```text
.company = organization / score
aw.tui   = conductor / observation
repos/   = organizational data
AW       = execution
```

Next layers can add embeddings, BQML clustering, GitHub Actions/AW state, and a terminal UI without changing the repository graph contract.
