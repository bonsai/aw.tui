# BONSAI 設計図

## 概要

`soubi` と `aw.tui` は同じ組織を見る 2 つの窓。

- `soubi` = 倉庫の棚を見る窓
- `aw.tui` = 作業現場を見る窓

## 1. 画面構成

### 1.1 soubi（資産閲覧）

```text
┌─────────────────────────────────────────────────────────────┐
│ 🌱 SOUBI        組織の装備を見る                            │
├─────────────────────────────────────────────────────────────┤
│ [R]epos  [A]gents  [S]kills  [M]CPs  [Mo]dels  [W]f  [D]ata │
├───────────────┬─────────────────────────────────────────────┤
│ repositories  │ github-observatory                          │
│ > agents      │   score: 94.8                               │
│   skills      │   health: 91.2                              │
│   mcps        │   cluster: data                             │
│   models      │                                             │
│   workflows   │ 装備: repo-discovery, bqml-analysis         │
│   datasets    │                                             │
├───────────────┴─────────────────────────────────────────────┤
│ [e]quip  [i]nspect  [c]alculate  [q]uit                     │
└─────────────────────────────────────────────────────────────┘
```

### 1.2 aw.tui（実行操作）

```text
┌─────────────────────────────────────────────────────────────┐
│ ⚡ AW.TUI       実行と承認の現場                             │
├─────────────────────────────────────────────────────────────┤
│ [P]lans  [W]orkflows  [R]uns  [A]pprovals  [Re]sults        │
├───────────────┬─────────────────────────────────────────────┤
│ active runs   │ workflow: repo-health-check                 │
│ > #42         │   agent: repo-coordinator                   │
│   #41         │   status: waiting approval                  │
│   #40         │                                             │
│               │ ステップ 3/5: BQML スコア更新               │
│               │                                             │
│               │ [a]pprove  [r]eject  [v]iew logs            │
├───────────────┴─────────────────────────────────────────────┤
│ [n]ew run  [s]top  [q]uit                                   │
└─────────────────────────────────────────────────────────────┘
```

## 2. ナビゲーション

### soubi

| キー | 動作 |
|------|------|
| `r` | repositories タブ |
| `a` | agents タブ |
| `s` | skills タブ |
| `m` | mcps タブ |
| `o` | models タブ |
| `w` | workflows タブ |
| `d` | datasets タブ |
| `↑↓` | 項目選択 |
| `e` | equip（選択中資産を Agent に装備） |
| `i` | inspect（詳細確認） |
| `c` | calculate（スコア再計算） |
| `q` | 終了 |

### aw.tui

| キー | 動作 |
|------|------|
| `p` | plans タブ |
| `w` | workflows タブ |
| `r` | runs タブ |
| `a` | approvals タブ |
| `e` | results タブ |
| `↑↓` | 項目選択 |
| `n` | new run（新規実行） |
| `s` | stop（実行停止） |
| `a` | approve（承認） |
| `x` | reject（却下） |
| `v` | view logs（ログ閲覧） |
| `q` | 終了 |

## 3. イベント発行

TUI はすべての操作をイベントとして発行する。

```text
soubi:
  VIEW_ASSET      {genre, id}
  EQUIP           {agent, asset}
  INSPECT         {genre, id}
  CALCULATE       {target}

aw.tui:
  VIEW_RUN        {run_id}
  NEW_RUN         {workflow, input}
  STOP_RUN        {run_id}
  APPROVE         {run_id, step}
  REJECT          {run_id, step}
```

これらのイベントは `AW` / `Agent` / `pi` で処理される。

## 4. 状態の反映

```text
github-observatory
  ↓
BQML
  ↓
Organization State
  ↓
soubi / aw.tui （表示）
  ↓
user event
  ↓
AW / Agent / pi
  ↓
状態変化 → github-observatory → BQML → ...
```

## 5. 統一アクション

ドラクエ風メニューを両 TUI で共通キーにする。

```text
みる      → m / view
せってい  → c / configure
そうび    → e / equip
かくにん  → i / inspect
けいさん  → = / calculate
クエスト  → q / run
承認      → a / approve
```
