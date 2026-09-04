# BONSAI UX 物語

## 朝、8:45。bonsai 本社、地下ラボ。

Sさんはターミナルを開いた。最初に立ち上がるのはいつもの `soubi`。

```text
🌱 SOUBI        組織の装備を見る
```

「今日の組織はどうかな」

Sさんは `r` を押す。Repositories のタブが開く。`github-observatory` がリストの先頭にあった。右パネルには score: 94.8、health: 91.2、cluster: data と表示されている。

「観測所は相変わらず元気だ」

`↓` で `repos` に移動。score は 72.4。health は 68.1。cluster は legacy。

「……これはちょっと心配だな」

Sさんは `i` を押した。inspect イベントが発行される。3秒後、右パネルに詳細が表示された。最後の更新は 3 週間前。最近のコミットはメンテナンスのみ。

「放置されてる。誰かに見てもらおう」

`a` で Agents タブ。`repo-coordinator` を選び、`e` を押す。equip イベント。

```text
equip repo-coordinator with repository: repos
```

Sさんは続けて `w` を押した。Workflows タブ。`repo-health-check` という Workflow があった。`i` で中身を確認する。

```yaml
workflow: repo-health-check
steps:
  - action: inspect
    target: repository
  - action: calculate
    model: bqml-health
  - action: run
    agent: repo-coordinator
  - action: approve
    if: score < 80
```

「これでいこう」

Sさんは `q` で `soubi` を閉じ、`aw.tui` を開いた。

```text
⚡ AW.TUI       実行と承認の現場
```

`w` で Workflows タブ。`repo-health-check` を選び、`n` を押す。new run イベント。

```text
new run: repo-health-check
input: {target: repos}
```

しばらくして `r` の Runs タブに `#47` が追加された。status: running。

Sさんはコーヒーを淹れに行った。

---

## 8:52。aw.tui の通知音。

Sさんは席に戻る。`a` の Approvals タブが点滅している。

```text
run #47  step 4/5  waiting approval
repo-coordinator: repos の health score が 68.1 です。
推奨アクション: リファクタリング計画を作成
```

Sさんは `v` を押してログを確認した。BQML の計算結果、最近の活動低下、依存関係の老朽化。すべて数字で示されている。

「うん、これは本当にやらないといけないな」

`a` を押す。approve イベント。

```text
approved: run #47 step 4/5
```

---

## 9:15。pi の作業終了通知。

Sさんは `e` の Results タブを開く。`#47` の結果が表示されている。

```text
result #47
├── created: repo-health-check-20250903.md
├── assigned: repo-coordinator
└── next: 週次レビューに追加
```

Sさんはもう一度 `soubi` を開き、`r` で repositories を見た。

`repos` の health はまだ 68.1 のままだ。でも score の横に小さく表示されていた。

```text
score: 72.4  → 計画中（run #47）
```

「動き始めた」

Sさんは満足してターミナルを閉じた。

---

## 同日、14:00。Mさんのターン。

Mさんは `soubi` を開き、`d` で Datasets タブに移動した。

「新しい訓練データが入ってるかな」

` Fine-tuning dataset v3` が追加されていた。score: 88.2。cluster: ml。

Mさんは `m` で Models タブに移動し、`gemma-finetuned` を選んだ。`e` を押す。

```text
equip agent: ml-engineer with model: gemma-finetuned, dataset: fine-tuning-v3
```

次に `aw.tui` で `p` の Plans タブを開く。先週の計画 `ml-model-refresh` があった。`n` で新しい run を作成。

```text
new run: ml-model-refresh
input:
  model: gemma-finetuned
  dataset: fine-tuning-v3
```

Mさんは承認待ちの間、他の仕事をした。

---

## 16:30。承認画面。

```text
run #48  waiting approval
ml-engineer: ファインチューニングを開始します。
推定コスト: $12.40
推定時間: 45分
```

Mさんは `a` で承認した。

```text
approved: run #48
pi: ファインチューニングを開始します
```

45分後、Results タブに新しいモデルが表示されていた。

```text
result #48
├── model: gemma-finetuned-v4
├── metrics: accuracy 0.91 (↑ +0.04)
└── deployed: false
```

Mさんは `soubi` に戻り、Models の `gemma-finetuned-v4` を確認した。score: 91.0。cluster: ml。

「よし、これで次の週次レビューに出せる」

---

## 夜、22:00。自動実行。

組織のスケジューラが動いた。`github-observatory` が全リポジトリを観測し、BQML が新しい State を計算する。

`soubi` を開かなくても、組織の状態は更新されていた。

明日の朝、誰かが `soubi` を開いたとき、新しい数字が待っている。
