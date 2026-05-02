# aict - AI 命令行工具

用中文描述你想做什么，AI 生成对应的 shell 命令并解释，确认后直接执行。

## 功能特性

- **自然语言转命令** — 输入中文需求描述，AI 自动生成可执行命令并给出详细解释
- **多 AI 后端** — 支持 Claude、OpenAI、DeepSeek、Ollama 及任意 OpenAI 兼容 API
- **流式输出** — 实时显示 AI 回复，体验流畅
- **安全黑名单** — 内置高危命令拦截（`rm -rf /`、`mkfs` 等），支持用户自定义规则
- **交互式确认** — 执行前可选择：执行 / 取消 / 编辑 / 重新生成 / 查看详情 / 加入黑名单
- **跨平台** — 支持 Windows（PowerShell/CMD）、macOS、Linux（Bash/Zsh/Fish）
- **便携/安装双模式** — 可执行文件同目录放 `.aict/` 即为便携版，否则使用 `~/.aict/`
- **命令日志** — 自动记录执行过的命令，按天归档到 `log/` 目录
- **多轮对话** — 支持上下文连续对话，AI 能理解前后文关系

## 快速开始

### 从源码构建

```bash
git clone https://github.com/aict-tool/aict.git
cd aict
go build -o aict ./cmd/aict/
```

Windows 用户也可使用 PowerShell 构建脚本：

```powershell
.\build.ps1 -Target windows
```

### 首次运行

```bash
./aict
```

首次启动会进入配置向导，引导你选择 AI provider、填写 API Key 和模型名称。也可以手动运行向导：

```bash
./aict init
```

### 使用示例

```text
> 列出当前目录下最大的10个文件

## 解释
使用 du 命令统计文件大小，sort 排序后取前 10 条...

## 命令
$ du -ah . | sort -rh | head -10

[y] 执行  [N] 取消  [e] 编辑  [r] 重新生成  [d] 更详细  [b] 加黑名单
> y
```

## CLI 命令

| 命令              | 说明                         |
| ----------------- | ---------------------------- |
| `aict`            | 进入交互式 REPL              |
| `aict init`       | 重新运行配置向导             |
| `aict config show`| 显示当前配置（API Key 脱敏） |
| `aict version`    | 显示版本号                   |
| `aict help`       | 显示帮助信息                 |

## REPL 元命令

在交互界面中输入以下命令：

| 命令               | 说明                                                  |
| ------------------ | ----------------------------------------------------- |
| `:exit` / `:quit`  | 退出                                                  |
| `:reset`           | 清空当前对话历史                                      |
| `:provider <name>` | 切换 AI provider（如 `claude`、`openai`、`ollama`）   |
| `:model <name>`    | 切换当前 provider 的模型                              |
| `:blacklist`       | 列出所有黑名单规则                                    |
| `:config dir`      | 显示配置目录及运行模式                                |
| `:help`            | 显示帮助                                              |

## 配置

配置文件为 TOML 格式，位于配置目录下的 `config.toml`：

```toml
default_provider = "claude"

[providers.claude]
base_url = "https://api.anthropic.com"
api_key  = "sk-ant-xxx"
model    = "claude-sonnet-4-6"

[providers.ollama]
base_url = "http://localhost:11434/v1"
model    = "qwen2.5-coder:7b"

[ui]
stream = true
color  = "auto"
```

### 环境变量覆盖

| 环境变量                     | 说明                              |
| ---------------------------- | --------------------------------- |
| `AICT_PROVIDER`              | 覆盖默认 provider                 |
| `AICT_CLAUDE_API_KEY`        | 覆盖 Claude API Key               |
| `AICT_OPENAI_API_KEY`        | 覆盖 OpenAI API Key               |
| `AICT_<PROVIDER>_MODEL`      | 覆盖指定 provider 的模型          |
| `AICT_<PROVIDER>_BASE_URL`   | 覆盖指定 provider 的 API 地址     |

## 便携版

将 `aict` 可执行文件和 `.aict/` 空目录放在同一文件夹中，程序会自动检测为便携模式，配置和日志都保存在该目录下。

```powershell
# Windows 打包便携版
.\build.ps1 -Target portable-windows
```

## 项目结构

```text
cmd/aict/main.go          入口、子命令分发
internal/config/           配置加载、保存、向导
internal/provider/         AI Provider（Claude/OpenAI/兼容）
internal/parser/           流式 Markdown 解析，提取命令和解释
internal/shell/            Shell 检测与命令执行
internal/prompt/           System prompt 构建
internal/safety/           安全黑名单
internal/cmdlog/           命令执行日志
internal/repl/             REPL 主循环、渲染、确认菜单、元命令
```

## License

MIT
