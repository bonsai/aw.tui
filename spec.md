# BONSAI MVC + Ontology 再設計

## 1. 目的

`soubi` と `aw.tui` の責務を **MVC + Ontology** で再定義し、競合を解消する。

| 問い | 答え |
|------|------|
| What do we have? | `soubi.tui` |
| What are we doing? | `aw.tui` |
| What is it? | `ontology/` |
| What is its state? | `BQML` |
| What should we do? | `AW / Agent` |
| Do it. | `pi` |

## 2. オントロジー（命名空間）

```text
bonsai
│
├── ontology/        ← 何があるかの定義（What is it?）
│   ├── organization.yaml
│   ├── repository.yaml
│   ├── asset.yaml
│   ├── agent.yaml
│   ├── workflow.yaml
│   ├── state.yaml
│   └── action.yaml
│
├── repos/           ← 組織の資産データ
├── github-observatory/  ← GitHub → BigQuery パイプライン
├── yaml-as-agent/   ← YAML で記述された Agent 定義
├── soubi/           ← 資産閲覧 TUI（What do we have?）
└── aw.tui/          ← 実行操作 TUI（What are we doing?）
```

### Namespace

```text
bonsai.organization
bonsai.repository
bonsai.asset
bonsai.agent
bonsai.workflow
bonsai.state
bonsai.action
```

## 3. MVC 分離

### Model（真実のデータ・契約・計算）

```text
bonsai/
├── repos/              ← 生のリポジトリデータ
├── github-observatory/ ← GitHub → BigQuery → BQML
├── yaml-as-agent/      ← Agent 定義
└── repo.yaml           ← 組織全体の状態・契約
```

計算フロー:

```text
GitHub
  ↓
Raw
  ↓
BigQuery
  ↓
BQML
  ├── Health
  ├── Activity
  ├── Score
  ├── Cluster
  └── Relationship
          ↓
    Organization State
          ↓
    soubi.tui / aw.tui
```

BQML は **State Generator / Analyzer** として扱う。

### View（表示のみ）

```text
bonsai/soubi          = 資産を見る
bonsai/aw.tui         = 行動を見る・操作する
```

#### soubi の View

```text
organization
├── repos
├── agents
├── skills
├── mcp
├── models
└── workflows
```

#### aw.tui の View

```text
execution
├── plans
├── workflows
├── runs
├── approvals
└── results
```

### Controller（判断はここで行う）

TUI には Controller ロジックを埋めない。

```text
soubi.tui ──────┐
                │ event
aw.tui ─────────┤
                ▼
               AW
                │
                ▼
          Agent / Coordinator
                │
                ▼
               pi
```

TUI はイベントを発行するだけ:

- `github-observatory を見る`
- `この repo を装備する`
- `この AW を実行する`
- `approve`

判断は AW / Agent 側。

## 4. 装備のオントロジー

```text
Repository
     │
     │ owns
     ▼
Asset
     │
     │ capability
     ▼
Agent
     │
     │ equipped-with
     ├──── Skill
     ├──── MCP
     ├──── Model
     └──── Repository
```

例:

```yaml
agent: repo-coordinator

equipment:
  repositories:
    - github-observatory
    - repos
  skills:
    - repo-discovery
    - bqml-analysis
  models:
    - gemma
  mcp:
    - github
```

## 5. State は別 namespace

```text
bonsai.state
│
├── activity
├── health
├── score
├── cluster
├── relation
├── execution
└── assignment
```

例:

```yaml
repository: github-observatory

state:
  activity: 94.8
  health: 91.2
  score: 87.4
  cluster: data
```

TUI はこれを表示するだけ。

## 6. Action（ドラクエ風操作）

```text
みる      → view
せってい  → configure
そうび    → equip
かくにん  → inspect
けいさん  → calculate
クエスト  → run / workflow
承認      → approve
```

これらを `bonsai.action` として統一する。

## 7. 現在の aw.tui からの移行

### 現状

```text
aw.tui/
├── cmd/aw-tui/main.go   ← View + Controller が混在
├── graph/graph.go       ← Model（似ている）
├── repos/repos.go       ← Model（GitHub クライアント）
└── repos/               ← データサブモジュール（git submodule?）
```

### 移行後

```text
bonsai/aw.tui/
├── cmd/aw-tui/main.go   ← View のみ
├── internal/
│   ├── tui/             ← Bubble Tea View
│   └── events/          ← イベント発行のみ
└── go.mod

bonsai/github-observatory/
├── cmd/observatory/
├── internal/
│   ├── github/          ← repos.FetchRecent 相当
│   ├── bigquery/
│   └── bqml/
└── go.mod

bonsai/ontology/
├── organization.yaml
├── repository.yaml
├── asset.yaml
├── agent.yaml
├── workflow.yaml
├── state.yaml
└── action.yaml

bonsai/yaml-as-agent/
├── agents/
│   └── repo-coordinator.yaml
└── schemas/
    └── equipment.json
```

### 責務移動

| 現コード | 移行先 |
|----------|--------|
| `repos.Repo` | `bonsai.github-observatory` / `ontology/repository.yaml` |
| `repos.FetchRecent` | `github-observatory/internal/github` |
| `graph.Build` | `github-observatory/internal/bqml` または `ontology` 経由の計算 |
| `cmd/aw-tui` の TUI | `aw.tui/internal/tui`（View のみ） |
| `cmd/aw-tui` のキー処理 | `aw.tui/internal/events`（イベント発行） |

## 8. 実装順序

1. `ontology/` YAML スキーマ作成（`schema.json` 参照）
2. `github-observatory/` へ `repos` + `graph` を移動
3. `aw.tui` を View + events のみに整理
4. `soubi` 新規作成（資産閲覧 View）
5. `yaml-as-agent` 新規作成（Agent 装備定義）
6. `aw.tui` と `soubi` からのイベントを AW / Agent に接続

## 9. 非機能

- `aw.tui` は依然として Bubble Tea + lipgloss
- `soubi` も同じスタックで統一
- Model は Go の独立モジュール
- ontology は YAML + JSON Schema で検証
- BQML は BigQuery 上で実行

## 10. 完了基準

- [ ] `ontology/` に 7 つの YAML スキーマが存在
- [ ] `aw.tui` が View + events のみ
- [ ] `github-observatory` が Model として独立
- [ ] `soubi` が資産閲覧を表示
- [ ] `yaml-as-agent` で Agent 装備を定義
- [ ] `aw.tui` / `soubi` からのイベントが Agent / pi に到達
