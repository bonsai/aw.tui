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

### Asset の 7 ジャンル

```text
bonsai.asset
│
├── repository    ← コード資産（旧 code）
├── dataset       ← データ・知識（旧 data）
├── model         ← 推論モデル
├── skill         ← 能力・関数
├── mcp           ← 外部接続
├── agent         ← 実行主体
└── workflow      ← 手順・AW
```

`soubi` はこの 7 ジャンルを閲覧する View。

## 3. エコシステム定義と責務分担

### エコシステム図

```text
                    BONSAI
                      │
          ┌───────────┴───────────┐
          │                       │
        MODEL                   VIEW
          │                       │
   ┌──────┼──────┐          ┌─────┴─────┐
   │      │      │          │           │
 repos   BQML  ontology   soubi.tui  aw.tui
   │      │                  │           │
   │    state                │ events   │ events
   │      │                  └─────┬─────┘
   └──────┴────────────────────────┘
                                  │
                              CONTROLLER
                                  │
                              AW / Agent
                                  │
                                  ▼
                                  pi
```

### 責務分担

| コンポーネント | 責務 | しないこと |
|----------------|------|------------|
| `ontology/` | 組織・資産・Agent・Workflow・State・Action の定義と契約 | データの取得・計算・表示 |
| `repos/` | 組織が所有するリポジトリの生データ保持 | 分析・状態計算 |
| `github-observatory/` | GitHub → BigQuery → BQML パイプライン。Feature 生成と State 計算 | UI・実行判断 |
| `yaml-as-agent/` | Agent の装備・能力・役割を YAML で定義 | 実際の実行 |
| `soubi.tui` | 資産の閲覧。What do we have? の表示 | 計算・判断・実行 |
| `aw.tui` | 実行・Workflow・承認の操作UI。What are we doing? の表示 | 計算・判断 |
| `AW / Agent` | イベント解釈・判断・Coordinator 制御 | 直接の GitHub 操作以外の低レベル作業は pi に委譲 |
| `pi` | 指示に基づく実作業・コード編集・コマンド実行 | 組織状態の計算・戦略判断 |

### データフロー

```text
GitHub API
    ↓
github-observatory (fetch + store)
    ↓
BigQuery (raw + feature tables)
    ↓
BQML (Health / Activity / Score / Cluster / Relation)
    ↓
Organization State (bonsai.state.*)
    ↓
soubi.tui （閲覧） / aw.tui （操作）
    ↓
event emit
    ↓
AW / Agent
    ↓
pi
```

### 人間とAIの役割

| 段階 | 人間 | AI（pi / Agent） |
|------|------|------------------|
| 設計・承認 | オントロジー・責務の最終判断 | 草案作成・影響分析 |
| Model 実装 | レビュー・承認 | github-observatory / ontology 実装 |
| View 実装 | 操作性の確認 | soubi / aw.tui 実装 |
| Controller 実装 | 重要イベントの承認設定 | AW / Agent 実装 |
| 運用 | Go/No-Go 判断 | 継続的な観測・実行・記録 |

## 4. MVC 分離

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

`soubi` が抱える資産ジャンルは **7 つ**に統一する。

```text
organization
├── repositories    ← コード資産
├── agents          ← 実行主体
├── skills          ← 能力・関数
├── mcps            ← 外部接続
├── models          ← 推論モデル
├── workflows       ← 手順・AW
└── datasets        ← データ・知識
```

統合ルール:

- `code` → `repositories`
- `data` → `datasets`
- `model` / `skill` / `mcp` / `agent` / `workflow` はそのまま

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

## 5. 装備のオントロジー

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

## 6. State は別 namespace

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

## 7. Action（ドラクエ風操作）

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

## 8. 現在の aw.tui からの移行

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

## 9. 実装順序

1. `ontology/` YAML スキーマ作成（`schema.json` 参照）
2. `github-observatory/` へ `repos` + `graph` を移動
3. `aw.tui` を View + events のみに整理
4. `soubi` 新規作成（資産閲覧 View）
5. `yaml-as-agent` 新規作成（Agent 装備定義）
6. `aw.tui` と `soubi` からのイベントを AW / Agent に接続

## 10. 非機能

- `aw.tui` は依然として Bubble Tea + lipgloss
- `soubi` も同じスタックで統一
- Model は Go の独立モジュール
- ontology は YAML + JSON Schema で検証
- BQML は BigQuery 上で実行

## 11. 完了基準

- [ ] `ontology/` に 7 つの YAML スキーマが存在
- [ ] `aw.tui` が View + events のみ
- [ ] `github-observatory` が Model として独立
- [ ] `soubi` が資産閲覧を表示
- [ ] `yaml-as-agent` で Agent 装備を定義
- [ ] `aw.tui` / `soubi` からのイベントが Agent / pi に到達
