<p align="center">
  <img src="docs/public/favicon.svg" alt="rdflow logo" width="120" />
</p>

<h1 align="center">rdflow</h1>

<p align="center">
  <a href="https://www.redmine.org/">Redmine</a> プロジェクト管理のためのコマンドラインインターフェース。
</p>

<p align="center">
  <a href="README.md">English</a> ·
  <a href="README.zh-CN.md">简体中文</a> ·
  <b>日本語</b>
</p>

<p align="center">
  <a href="https://github.com/TruthUnderlier/RedmineFlow/releases/latest"><img src="https://img.shields.io/github/v/release/TruthUnderlier/RedmineFlow?style=for-the-badge&logo=github&logoColor=white" alt="Release"></a>
  <a href="https://opensource.org/licenses/MIT"><img src="https://img.shields.io/badge/License-MIT-yellow?style=for-the-badge&logo=opensourceinitiative&logoColor=white" alt="License"></a>
</p>

rdflow は [aarondpn/redmine-cli](https://github.com/aarondpn/redmine-cli) の独立したフォークで、次の機能を追加しています。

- **カスタム CA 証明書**：プロファイルごとに `ca_cert` を指定でき、プライベート CA や企業 CA の Redmine サーバーに接続できます。
- **VPN 環境での名前解決の修正**：Sangfor などの VPN の DNS プロキシが EDNS0 クエリを拒否して起きる `no such host` を修正します。

コマンド名は `rdflow`、設定ファイルは `~/.redmine-flow.yaml` です。公式の `redmine` と同時にインストールでき、互いに影響しません。

## インストール

**インストールスクリプト**（macOS、Linux）：

```bash
curl -fsSL https://raw.githubusercontent.com/TruthUnderlier/RedmineFlow/main/install.sh | REDMINE_CLI_SKIP_ATTESTATION=1 bash
```

**PowerShell**（Windows）：

```powershell
$env:REDMINE_CLI_SKIP_ATTESTATION='1'; irm https://raw.githubusercontent.com/TruthUnderlier/RedmineFlow/main/install.ps1 | iex
```

[GitHub Releases](https://github.com/TruthUnderlier/RedmineFlow/releases/latest) からアーカイブをダウンロードし、`rdflow`（Windows では `rdflow.exe`）を `PATH` に置くこともできます。アップデートは `rdflow update --skip-attestation`。リリースには GitHub のビルド証明（attestation）がないため、その確認はスキップされますが、SHA-256 チェックサムは検証されます。

## はじめに

```bash
rdflow auth login
rdflow issues list
rdflow issues view 123
rdflow time log
rdflow issues list --project myproject --jq '.[].id'
```

`rdflow --help` を実行すると、利用可能なすべてのコマンドが表示されます。

## カスタム CA 証明書

`~/.redmine-flow.yaml` のプロファイルに CA 証明書（PEM 形式）のパスを指定します。

```yaml
profiles:
  work:
    server: https://redmine.internal
    auth_method: apikey
    api_key: your-api-key
    ca_cert: /etc/ssl/certs/company-ca.crt      # Windows: 'C:\certs\company-ca.crt'
```

証明書はシステムのルート証明書に追加される形で、rdflow から Redmine サーバーへの接続にのみ使われます。システムの信頼ストアや他のプログラムには影響しません。

## ドキュメント

上流のドキュメント [redmine-cli.dev](https://redmine-cli.dev) は rdflow にもそのまま当てはまります。コマンドの `redmine` を `rdflow` に、`~/.redmine-cli.yaml` を `~/.redmine-flow.yaml` に読み替えてください。

## AI エージェント

| ツール | インストール |
|---|---|
| **任意のエージェント（skill）** | `npx skills add TruthUnderlier/RedmineFlow` |
| **任意の MCP ホスト** | `rdflow mcp serve` |
| **Claude Code** | `/plugin marketplace add TruthUnderlier/RedmineFlow`、続いて `/plugin install redmine` |

## ライセンス

MIT。rdflow は Aaron Döppner による [redmine-cli](https://github.com/aarondpn/redmine-cli) をベースにしています。[LICENSE](LICENSE) を参照してください。
