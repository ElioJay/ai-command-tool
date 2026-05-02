# aict 使用指南（Windows 环境）

## 1. 安装

### 1.1 前置条件

- Windows 10/11
- Go 1.22+（仅从源码构建时需要）
- 一个 AI 服务的 API Key（Claude / OpenAI / DeepSeek），或本地运行的 Ollama

### 1.2 从源码构建

打开 PowerShell，执行：

```powershell
git clone https://github.com/aict-tool/aict.git
cd aict
go build -o aict.exe ./cmd/aict/
```

或使用项目内置的构建脚本：

```powershell
.\build.ps1 -Target windows
```

构建产物位于 `dist\aict-windows-amd64.exe`。

### 1.3 加入 PATH（可选）

将 `aict.exe` 所在目录加入系统 PATH，即可在任意位置直接调用：

```powershell
# 临时生效（当前会话）
$env:PATH += ";D:\tools\aict"

# 永久生效（需以管理员运行）
[Environment]::SetEnvironmentVariable("PATH", $env:PATH + ";D:\tools\aict", "User")
```

---

## 2. 首次配置

### 2.1 配置向导

首次运行 `aict.exe` 时，程序会自动进入配置向导：

```text
欢迎使用 aict！我们来完成首次配置。

[1/5] 选择默认 AI provider：
  1) Claude (官方)
  2) OpenAI (官方)
  3) DeepSeek
  4) Ollama (本地)
  5) 自定义 OpenAI 兼容接口
> 1

[2/5] Base URL（回车使用默认 https://api.anthropic.com）：
>

[3/5] API Key：
> sk-ant-xxxxx

[4/5] 模型名称（回车使用默认 claude-sonnet-4-6）：
>

[5/5] 配置文件保存位置：
  1) 用户目录 ~/.aict（安装版，推荐）
  2) 可执行文件同目录 .aict/（便携版）
> 1

配置已保存到 C:\Users\你的用户名\.aict
```

也可以随时重新运行向导：

```powershell
.\aict.exe init
```

### 2.2 配置文件

配置向导会在配置目录下生成 `config.toml`，你也可以直接编辑它：

```toml
default_provider = "claude"

[providers.claude]
base_url = "https://api.anthropic.com"
api_key  = "sk-ant-xxxxx"
model    = "claude-sonnet-4-6"

[ui]
stream = true
color  = "auto"
```

#### 配置多个 Provider

在同一个文件中添加多个 provider 段落，运行时用 `:provider` 命令切换：

```toml
default_provider = "claude"

[providers.claude]
base_url = "https://api.anthropic.com"
api_key  = "sk-ant-xxxxx"
model    = "claude-sonnet-4-6"

[providers.openai]
base_url = "https://api.openai.com/v1"
api_key  = "sk-xxxxx"
model    = "gpt-4o"

[providers.deepseek]
base_url = "https://api.deepseek.com/v1"
api_key  = "sk-xxxxx"
model    = "deepseek-chat"

[providers.ollama]
base_url = "http://localhost:11434/v1"
model    = "qwen2.5-coder:7b"
```

### 2.3 环境变量覆盖

不想把 API Key 写进文件时，可以用环境变量：

```powershell
# 当前会话生效
$env:AICT_CLAUDE_API_KEY = "sk-ant-xxxxx"

# 切换默认 provider
$env:AICT_PROVIDER = "openai"
```

支持的环境变量：

- `AICT_PROVIDER` — 覆盖默认 provider
- `AICT_<大写PROVIDER名>_API_KEY` — 覆盖 API Key
- `AICT_<大写PROVIDER名>_MODEL` — 覆盖模型
- `AICT_<大写PROVIDER名>_BASE_URL` — 覆盖 API 地址

---

## 3. 基本使用

### 3.1 启动

```powershell
.\aict.exe
```

启动后进入交互界面：

```text
aict 已启动（provider: claude，shell: pwsh）
输入中文描述需求，或输入 :help 查看帮助

>
```

### 3.2 日常使用流程

输入中文描述你想完成的操作，AI 会生成命令并解释：

```text
> 查看当前目录下所有 .go 文件的总行数

## 解释
使用 Get-ChildItem 递归查找所有 .go 文件，通过管道传给 Get-Content 读取内容，
再用 Measure-Object -Line 统计总行数。

## 命令
$ Get-ChildItem -Recurse -Filter "*.go" | Get-Content | Measure-Object -Line

─────────────────────────────────

$ Get-ChildItem -Recurse -Filter "*.go" | Get-Content | Measure-Object -Line

[y] 执行  [N] 取消  [e] 编辑  [r] 重新生成  [d] 更详细  [b] 加黑名单
>
```

### 3.3 确认菜单

AI 生成命令后，你可以选择以下操作：

| 按键 | 操作         | 说明                                           |
| ---- | ------------ | ---------------------------------------------- |
| `y`  | 执行         | 直接在当前 shell 中执行该命令                  |
| `N`  | 取消         | 放弃本次命令（默认选项，直接回车即取消）       |
| `e`  | 编辑         | 修改命令后再决定是否执行                       |
| `r`  | 重新生成     | 让 AI 换一种方案重新生成命令                   |
| `d`  | 更详细       | 让 AI 对命令的每个参数做更详细的解释           |
| `b`  | 加入黑名单   | 将该命令模式加入黑名单，以后自动拦截           |

### 3.4 多轮对话

aict 会记住对话上下文。你可以基于前一次的结果继续追问：

```text
> 查找当前目录下最大的5个文件
（AI 生成命令并执行后...）

> 把刚才结果里最大的那个文件删除
（AI 会基于上下文理解"刚才的结果"）
```

输入 `:reset` 可清空对话历史，开始新的对话。

---

## 4. 元命令

在交互界面中，以 `:` 开头的输入会被识别为元命令（不会发送给 AI）：

```text
> :help

可用元命令：
  :exit / :quit       退出 aict
  :reset              清空当前对话历史
  :provider <name>    切换 AI provider（claude/openai/ollama 等）
  :model <name>       切换当前 provider 的模型
  :blacklist          列出所有黑名单规则
  :config dir         显示配置目录及运行模式
  :help               显示此帮助
```

### 4.1 切换 Provider 和模型

```text
> :provider openai
已切换到 provider: openai

> :model gpt-4o-mini
已切换模型: gpt-4o-mini
```

切换仅在当前会话生效，不会修改配置文件。

### 4.2 查看配置信息

```text
> :config dir
配置目录：C:\Users\elio\.aict
运行模式：installed
```

---

## 5. 安全黑名单

### 5.1 工作原理

aict 内置了一组高危命令拦截规则（正则匹配），包括：

- `rm -rf /` — 递归删除根目录
- `mkfs.*` — 格式化文件系统
- `format C:` — Windows 分区格式化
- `dd of=/dev/sd*` — 直接写入磁盘设备
- `curl ... | bash` — 远程脚本直接执行
- 等等

当 AI 生成的命令命中黑名单时，aict 会显示警告并要求逐字输入命令才能执行：

```text
[黑名单] 命令命中规则: mkfs
  原因：格式化文件系统会清除所有数据
  规则模式：\bmkfs(\.\\w+)?\b

如确认执行，请逐字输入命令（直接回车取消）：
>
```

### 5.2 查看规则

```text
> :blacklist
[builtin] rm-rf-root
  模式：^\s*rm\s+(-[rRfF]+\s+)+/\s*$
  原因：递归删除根目录会破坏整个系统
[builtin] mkfs
  模式：\bmkfs(\.\\w+)?\b
  原因：格式化文件系统会清除所有数据
...
```

### 5.3 添加自定义规则

在确认菜单中按 `b` 可将当前命令加入黑名单。程序会引导你确认正则模式和原因，规则保存在配置目录下的 `blacklist.toml` 中。

你也可以直接编辑 `blacklist.toml`：

```toml
[[rules]]
id      = "my-rule"
pattern = "\\bRemove-Item\\s+-Recurse"
reason  = "防止误删整个目录树"
source  = "user"
```

---

## 6. 命令日志

每次通过 aict 执行的命令都会自动记录在配置目录下的 `log/` 文件夹，按天归档：

```text
C:\Users\你的用户名\.aict\
  └── log\
      ├── 2026-04-29.log
      ├── 2026-04-30.log
      └── 2026-05-01.log
```

日志格式：

```text
[14:32:05] [OK] Get-ChildItem -Recurse -Filter "*.go" | Measure-Object -Line
[14:33:12] [OK] Get-Process | Sort-Object CPU -Descending | Select-Object -First 10
[14:35:40] [FAIL] ping 192.168.1.999
```

每条记录包含时间、执行结果（OK/FAIL）和完整命令。

---

## 7. 便携版

便携版将配置、日志等所有数据保存在可执行文件同目录的 `.aict/` 文件夹下，适合 U 盘携带使用。

### 7.1 制作便携版

```powershell
.\build.ps1 -Target portable-windows
```

会在 `dist\` 下生成 `aict-portable-windows-amd64.zip`，解压后结构如下：

```text
portable-windows\
  ├── aict.exe
  └── .aict\          ← 空目录，程序检测到它就进入便携模式
```

### 7.2 手动制作

将 `aict.exe` 和一个空的 `.aict\` 文件夹放在同一目录即可：

```powershell
mkdir D:\my-aict\.aict
copy dist\aict.exe D:\my-aict\aict.exe
```

运行后通过 `:config dir` 确认模式：

```text
> :config dir
配置目录：D:\my-aict\.aict
运行模式：portable
```

---

## 8. CLI 子命令

除交互式 REPL 外，aict 还支持以下子命令：

```powershell
# 显示版本
.\aict.exe version

# 显示帮助
.\aict.exe help

# 重新运行配置向导
.\aict.exe init

# 查看当前配置（API Key 脱敏显示）
.\aict.exe config show
```

`config show` 输出示例：

```text
运行模式：installed
配置目录：C:\Users\elio\.aict
默认 Provider：claude
  [claude] base_url=https://api.anthropic.com model=claude-sonnet-4-6 api_key=sk-a...xxxx
```

---

## 9. 常见使用场景

### 查找文件

```text
> 找出 src 目录下所有超过 1MB 的文件
```

### 进程管理

```text
> 查看占用 CPU 最高的 10 个进程
```

### 网络诊断

```text
> 检查 443 端口是否被占用，显示占用的进程
```

### Git 操作

```text
> 查看最近一周内谁提交了哪些文件
```

### 磁盘清理

```text
> 统计 C 盘各文件夹的大小，按从大到小排序
```

### 批量重命名

```text
> 把当前目录下所有 .jpeg 文件的扩展名改成 .jpg
```

---

## 10. 常见问题

### Q: 提示 "provider xxx 未配置" 怎么办？

检查配置文件中是否有对应的 `[providers.xxx]` 段落，或用 `aict init` 重新配置。

### Q: 如何使用本地 Ollama？

1. 确保 Ollama 已启动（默认监听 `localhost:11434`）
2. 运行 `aict init`，选择 `4) Ollama (本地)`
3. 模型名称填你已下载的模型，如 `qwen2.5-coder:7b`

### Q: API 请求报错 401 怎么办？

API Key 无效或过期。检查方式：

- 运行 `aict config show` 查看当前 Key 是否正确
- 用环境变量临时覆盖测试：`$env:AICT_CLAUDE_API_KEY = "新的key"`

### Q: 如何清空对话历史？

在 REPL 中输入 `:reset`，或退出后重新启动。

### Q: 日志文件在哪里？

运行 `:config dir` 查看配置目录，日志位于其下的 `log\` 子目录。
