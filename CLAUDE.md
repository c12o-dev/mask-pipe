# CLAUDE.md

このリポジトリは `mask-pipe` — ターミナル出力のシークレットをパイプ経由でマスクするCLIツール — の実装リポジトリです。

## ワークフロー: Spec-Driven × Issue-Driven ハイブリッド

このリポジトリでは**仕様駆動開発（SDD）とIssue駆動開発（IDD）のハイブリッド**を採用しています。

### 基本原則

1. **Specが真実の源（source of truth）**。振る舞い・契約・パターンは `docs/specs/` に記述
2. **Issueは作業の単位**。何を・なぜ・いつやるかを記録
3. **Specの変更もIssueになる**。仕様変更は `spec_change` Issue→議論→PR（Specとコードを同時に更新）
4. **PRは必ずIssueにリンク**する（`Fixes #N` または `Refs #N`）

### 作業の流れ

```
 [新機能・変更]               [バグ・質問]
     ↓                           ↓
 Issue（feature or              Issue（bug）
   spec_change）                  ↓
     ↓                        既存Specに照らして
 Spec更新が必要か？              再現・診断
     ↓ Yes                       ↓
 docs/specs/ に                PR（Fixes #N）
 差分を提案するPR                 ↓
     ↓                        テスト＋コード修正
 レビュー・マージ
     ↓
 実装PR（Fixes #N）
     ↓
 テスト＋コード＋必要ならSpec追記
```

### Claude Code セッションでのガイドライン

Claude Code がこのリポジトリで作業するときは以下に従ってください:

1. **コード変更前にSpec確認**: 該当する `docs/specs/NNN-*.md` を Read してから実装する
2. **Issue番号を要求**: 明示的な Issue 番号がない変更は、まず Issue を立てるか `--no-verify` 的なショートカットを取らない
3. **SpecとコードをPRで同時更新**: 振る舞いが変わるならSpecも必ず更新
4. **ADRを書く**: 後戻りしにくい技術選定（言語・依存関係・アーキテクチャ）は `docs/adr/NNNN-*.md` に記録
5. **テスト優先**: パターン追加や挙動変更は失敗するテストを先に書いてから実装する

### よくある作業パターン

**新しい組み込みパターンを追加したい:**
1. `pattern_proposal` Issueを立てる（5個以上のmatch例・no-match例を含める）
2. `docs/specs/002-pattern-library.md` を更新するPR
3. 実装PR（`patterns/` ディレクトリ＋テスト）

**CLIフラグを追加したい:**
1. `feature` または `spec_change` Issue
2. `docs/specs/001-cli-interface.md` を更新するPR（議論が分かれる場合はまずSpecだけ）
3. 実装PR

**バグ修正:**
1. `bug` Issue（再現手順・期待値・実際の値）
2. PR（Fixes #N）。既存Specに沿った振る舞いへの修正か、Specの曖昧さが原因かを判断

## ディレクトリ構造

| パス | 用途 |
|---|---|
| `docs/specs/` | 製品仕様。振る舞い・契約・パターン。**実装の前にここを更新する** |
| `docs/adr/` | アーキテクチャ決定記録。後戻りしにくい選択の理由 |
| `.github/ISSUE_TEMPLATE/` | Issueフォーム |
| `.github/PULL_REQUEST_TEMPLATE.md` | PRテンプレート |
| `.claude/` | Claude Code用のコマンド・エージェント・権限設定 |
| `CONTRIBUTING.md` | コントリビュータ向け（英語） |
| `README.md` | ユーザー向け（英語） |

## Claude Code 用ツール (`.claude/`)

このリポジトリには、頻出作業を支援するスラッシュコマンドとエージェントが用意されています:

| 種類 | 名前 | 用途 |
|---|---|---|
| コマンド | `/spec-new <NNN> <title>` | `docs/specs/NNN-*.md` を雛形から作成 + 索引更新 |
| コマンド | `/adr-new <NNNN> <title>` | `docs/adr/NNNN-*.md` をテンプレから作成 + 索引更新 |
| コマンド | `/pattern-new <pattern_id>` | Spec更新 + コード雛形 + テスト雛形を同時生成 |
| エージェント | `pattern-reviewer` | 新パターンの誤検知リスクを実コーパスで検証 |

`.claude/settings.json` は `go test`・`go build`・`gh issue view` 等を allow、`.env`・`~/.ssh/**` を deny、`git push` や `gh pr create` は ask。セキュリティツールとして**自身のシークレット取り扱いに厳格**であることを設定レベルで表明しています。

## 言語・技術スタック

- **実装言語**: Go（`docs/adr/0001-language-go.md` 参照）
- **ビルド**: `go build` / 将来は `goreleaser` でクロスコンパイル配布
- **テスト**: Go標準 `testing` パッケージ。パターンは `patterns/patterns_test.go` に網羅

## ルール

- **コメントは非自明な「なぜ」のみ**。命名で「何」が分かるならコメント不要
- **パターンの精度は再現率より優先**。誤検知1件は検出漏れ10件より嫌われる
- **依存関係を増やすときはADRで正当化**。単機能の軽量バイナリであることが価値
- **破壊的変更はmajor version bump必須**。CLIフラグの意味変更は破壊的変更
- **セキュリティツールとしての責任**: `--verify-patterns` のような自己診断機能を落とさない

## 参考

このプロジェクトの設計検討は別リポジトリ `~/idea/cli-drafts/` にあります（非公開）。そこにある深掘り・競合調査・批評は参考資料として残っていますが、このリポジトリ内では参照しないでください（パブリックリポジトリなので）。
