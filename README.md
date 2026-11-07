<p align="center">
  <img src="docs/public/favicon.svg" alt="rdflow logo" width="120" />
</p>

<h1 align="center">rdflow</h1>

<p align="center">
  A command-line interface for <a href="https://www.redmine.org/">Redmine</a> project management.
</p>

<p align="center">
  <b>English</b> ·
  <a href="README.zh-CN.md">简体中文</a> ·
  <a href="README.ja.md">日本語</a>
</p>

<p align="center">
  <a href="https://github.com/TruthUnderlier/RedmineFlow/releases/latest"><img src="https://img.shields.io/github/v/release/TruthUnderlier/RedmineFlow?style=for-the-badge&logo=github&logoColor=white" alt="Release"></a>
  <a href="https://opensource.org/licenses/MIT"><img src="https://img.shields.io/badge/License-MIT-yellow?style=for-the-badge&logo=opensourceinitiative&logoColor=white" alt="License"></a>
</p>

<p align="center">
  <a href="#installation">Installation</a> ·
  <a href="#getting-started">Getting Started</a> ·
  <a href="#custom-ca-certificate">Custom CA certificate</a> ·
  <a href="#ai-agents">AI Agents</a>
</p>

rdflow is an independent fork of [aarondpn/redmine-cli](https://github.com/aarondpn/redmine-cli). It adds:

- **Custom CA certificates** per profile (`ca_cert`), for Redmine servers behind a private or corporate CA.
- **Working DNS behind VPN DNS proxies** such as Sangfor, which reject EDNS0 queries (`no such host`).

The command is `rdflow` and its config file is `~/.redmine-flow.yaml`, so it can be installed next to the upstream `redmine` without either affecting the other.

## Installation

**Install script** (macOS, Linux):

```bash
curl -fsSL https://raw.githubusercontent.com/TruthUnderlier/RedmineFlow/main/install.sh | REDMINE_CLI_SKIP_ATTESTATION=1 bash
```

**PowerShell** (Windows):

```powershell
$env:REDMINE_CLI_SKIP_ATTESTATION='1'; irm https://raw.githubusercontent.com/TruthUnderlier/RedmineFlow/main/install.ps1 | iex
```

Or download an archive from [GitHub Releases](https://github.com/TruthUnderlier/RedmineFlow/releases/latest) and put `rdflow` (`rdflow.exe` on Windows) on your `PATH`:

| System | Archive |
|---|---|
| Linux x86_64 | `rdflow-linux-amd64.tar.gz` |
| Linux ARM64 | `rdflow-linux-arm64.tar.gz` |
| macOS Intel | `rdflow-darwin-amd64.tar.gz` |
| macOS Apple silicon | `rdflow-darwin-arm64.tar.gz` |
| Windows x86_64 | `rdflow-windows-amd64.zip` |

Update with `rdflow update --skip-attestation`. Releases are published without GitHub build attestations, so installing and updating skip that check; the SHA-256 checksums are still verified.

## Getting Started

```bash
# Configure your Redmine server and API key
rdflow auth login

# List issues
rdflow issues list

# View a specific issue
rdflow issues view 123

# Log time
rdflow time log

# Script it: JSON output filtered with the built-in --jq
rdflow issues list --project myproject --jq '.[].id'
rdflow issues get 123 --jq '.status.name'
```

Run `rdflow --help` to see all available commands.

`-o json` prints machine-readable JSON, and `-q/--jq <expression>` filters it (strings are printed raw, everything else as JSON). Errors always go to stderr; in JSON mode they are a single-line `{"error": {...}}` envelope on the last line of stderr, stdout stays empty and the exit code is non-zero.

## Custom CA certificate

If your server's certificate is signed by a private or corporate CA, requests fail with `x509: certificate signed by unknown authority`. Point the profile at the CA certificate in `~/.redmine-flow.yaml`:

```yaml
profiles:
  work:
    server: https://redmine.internal
    auth_method: apikey
    api_key: your-api-key
    ca_cert: /etc/ssl/certs/company-ca.crt      # Windows: 'C:\certs\company-ca.crt' or C:/certs/company-ca.crt
```

- The file must be PEM encoded (starts with `-----BEGIN CERTIFICATE-----`); the extension does not matter. Convert a binary DER certificate with `openssl x509 -inform der -in ca.cer -out ca.crt`, or on Windows `certutil -encode ca.cer ca.crt`.
- The certificates are trusted in addition to the system roots and only for rdflow's connections to the Redmine server; the system trust store and other programs are not changed.
- Add `ca_cert` to the profile before running `rdflow auth login` for the first time; logging in again keeps it.

## Documentation

The upstream documentation at [redmine-cli.dev](https://redmine-cli.dev) applies to rdflow as well. Replace `redmine` with `rdflow` in its commands, and `~/.redmine-cli.yaml` with `~/.redmine-flow.yaml`.

<a name="ai-agents"></a>

## Using with AI Agents

rdflow ships an [Agent Skill](https://agentskills.io) and a built-in MCP server:

| Tool | Install |
|---|---|
| **Any agent (skill)** | `npx skills add TruthUnderlier/RedmineFlow` |
| **Any MCP host** | `rdflow mcp serve` (then add to host config) |
| **Claude Code** | `/plugin marketplace add TruthUnderlier/RedmineFlow` then `/plugin install redmine` |
| **Codex CLI** | `codex plugin marketplace add TruthUnderlier/RedmineFlow` then `/plugins` |
| **Gemini CLI** | `gemini extensions install https://github.com/TruthUnderlier/RedmineFlow` |

The MCP server is read-only by default, with per-group / per-tool allow- and deny-lists. See the upstream [AI Agent Integration guide](https://redmine-cli.dev/guides/ai-agents/) for host-specific configuration.

## Development

```bash
make build   # builds bin/rdflow
make test
```

Local E2E testing against a real Redmine instance (Docker): see [e2e/README.md](e2e/README.md).

## License

MIT. rdflow is based on [redmine-cli](https://github.com/aarondpn/redmine-cli) by Aaron Döppner; see [LICENSE](LICENSE).
