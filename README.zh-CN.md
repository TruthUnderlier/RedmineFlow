<p align="center">
  <img src="docs/public/favicon.svg" alt="rdflow logo" width="120" />
</p>

<h1 align="center">rdflow</h1>

<p align="center">
  <a href="https://www.redmine.org/">Redmine</a> 项目管理的命令行工具。
</p>

<p align="center">
  <a href="README.md">English</a> ·
  <b>简体中文</b> ·
  <a href="README.ja.md">日本語</a>
</p>

<p align="center">
  <a href="https://github.com/TruthUnderlier/RedmineFlow/releases/latest"><img src="https://img.shields.io/github/v/release/TruthUnderlier/RedmineFlow?style=for-the-badge&logo=github&logoColor=white" alt="Release"></a>
  <a href="https://opensource.org/licenses/MIT"><img src="https://img.shields.io/badge/License-MIT-yellow?style=for-the-badge&logo=opensourceinitiative&logoColor=white" alt="License"></a>
</p>

<p align="center">
  <a href="#安装">安装</a> ·
  <a href="#快速开始">快速开始</a> ·
  <a href="#自定义-ca-证书">自定义 CA 证书</a> ·
  <a href="#ai-agents">AI 代理</a>
</p>

rdflow 是 [aarondpn/redmine-cli](https://github.com/aarondpn/redmine-cli) 的独立分支，新增了：

- **自定义 CA 证书**：每个 profile 可通过 `ca_cert` 指定证书，用于连接使用私有或企业 CA 的 Redmine 服务器。
- **VPN 环境下的域名解析修复**：修复深信服等 VPN 的 DNS 代理拒绝 EDNS0 查询导致的 `no such host`。

命令名为 `rdflow`，配置文件为 `~/.redmine-flow.yaml`，可以与官方 `redmine` 同时安装，互不影响。

## 安装

**安装脚本**（macOS、Linux）：

```bash
curl -fsSL https://raw.githubusercontent.com/TruthUnderlier/RedmineFlow/main/install.sh | REDMINE_CLI_SKIP_ATTESTATION=1 bash
```

**PowerShell**（Windows）：

```powershell
$env:REDMINE_CLI_SKIP_ATTESTATION='1'; irm https://raw.githubusercontent.com/TruthUnderlier/RedmineFlow/main/install.ps1 | iex
```

也可以从 [GitHub Releases](https://github.com/TruthUnderlier/RedmineFlow/releases/latest) 下载压缩包，解压后把 `rdflow`（Windows 为 `rdflow.exe`）放到 `PATH` 中的目录：

| 系统 | 文件 |
|---|---|
| Linux x86_64 | `rdflow-linux-amd64.tar.gz` |
| Linux ARM64 | `rdflow-linux-arm64.tar.gz` |
| macOS Intel | `rdflow-darwin-amd64.tar.gz` |
| macOS Apple 芯片（M 系列） | `rdflow-darwin-arm64.tar.gz` |
| Windows x86_64 | `rdflow-windows-amd64.zip` |

升级：`rdflow update --skip-attestation`。发布版本没有 GitHub 构建证明，因此安装和升级会跳过证明检查，但仍会校验 SHA-256。

## 快速开始

```bash
# 配置 Redmine 服务器和 API 密钥
rdflow auth login

# 列出 issue
rdflow issues list

# 查看指定 issue
rdflow issues view 123

# 记录工时
rdflow time log

# 脚本化：JSON 输出，并用内置的 --jq 过滤
rdflow issues list --project myproject --jq '.[].id'
rdflow issues get 123 --jq '.status.name'
```

运行 `rdflow --help` 查看所有可用命令。

`-o json` 输出机器可读的 JSON，`-q/--jq <表达式>` 对其进行过滤（字符串原样输出，其他值输出为 JSON）。错误始终输出到 stderr；JSON 模式下，错误是 stderr 最后一行的单行 `{"error": {...}}`，stdout 为空，退出码非零。

## 自定义 CA 证书

如果服务器证书由私有或企业 CA 签发，请求会报错 `x509: certificate signed by unknown authority`。在 `~/.redmine-flow.yaml` 的 profile 中指定 CA 证书：

```yaml
profiles:
  work:
    server: https://redmine.internal
    auth_method: apikey
    api_key: your-api-key
    ca_cert: /etc/ssl/certs/company-ca.crt      # Windows：'C:\certs\company-ca.crt' 或 C:/certs/company-ca.crt
```

- 文件必须是 PEM 格式（以 `-----BEGIN CERTIFICATE-----` 开头），与扩展名无关。二进制 DER 证书可用 `openssl x509 -inform der -in ca.cer -out ca.crt` 转换，Windows 上可用 `certutil -encode ca.cer ca.crt`。
- 这些证书在系统根证书之外额外被信任，且只用于 rdflow 到 Redmine 服务器的连接，不会修改系统信任库，也不影响其他程序。
- 首次执行 `rdflow auth login` 前，请先在 profile 中写好 `ca_cert`；之后重新登录会保留它。

## 文档

上游文档 [redmine-cli.dev](https://redmine-cli.dev) 同样适用于 rdflow，只需把其中命令里的 `redmine` 换成 `rdflow`，把 `~/.redmine-cli.yaml` 换成 `~/.redmine-flow.yaml`。

<a name="ai-agents"></a>

## 与 AI 代理配合使用

rdflow 提供 [Agent Skill](https://agentskills.io) 和内置的 MCP 服务器：

| 工具 | 安装方式 |
|---|---|
| **任意代理（skill）** | `npx skills add TruthUnderlier/RedmineFlow` |
| **任意 MCP 宿主** | `rdflow mcp serve`（然后添加到宿主配置） |
| **Claude Code** | `/plugin marketplace add TruthUnderlier/RedmineFlow`，然后 `/plugin install redmine` |
| **Codex CLI** | `codex plugin marketplace add TruthUnderlier/RedmineFlow`，然后 `/plugins` |
| **Gemini CLI** | `gemini extensions install https://github.com/TruthUnderlier/RedmineFlow` |

MCP 服务器默认只读，并支持按分组或按工具设置允许和禁止列表。各宿主的具体配置请参阅上游的 [AI 代理集成指南](https://redmine-cli.dev/zh-cn/guides/ai-agents/)。

## 开发

```bash
make build   # 生成 bin/rdflow
make test
```

针对真实 Redmine 实例的本地 E2E 测试（Docker）：参见 [e2e/README.md](e2e/README.md)。

## 许可证

MIT。rdflow 基于 Aaron Döppner 的 [redmine-cli](https://github.com/aarondpn/redmine-cli)，详见 [LICENSE](LICENSE)。
