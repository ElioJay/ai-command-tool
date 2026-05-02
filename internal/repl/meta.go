package repl

import (
	"fmt"
	"strings"
)

type MetaResult struct {
	Handled        bool
	ShouldExit     bool
	ResetHistory   bool
	SwitchProvider string
	SwitchModel    string
	ShowBlacklist  bool
	ShowConfigDir  bool
}

func HandleMeta(input string) MetaResult {
	input = strings.TrimSpace(input)
	if !strings.HasPrefix(input, "/") {
		return MetaResult{Handled: false}
	}
	parts := strings.Fields(input[1:])
	if len(parts) == 0 {
		return MetaResult{Handled: true}
	}
	cmd := strings.ToLower(parts[0])
	arg := ""
	if len(parts) > 1 {
		arg = parts[1]
	}

	switch cmd {
	case "exit", "quit":
		return MetaResult{Handled: true, ShouldExit: true}
	case "reset":
		fmt.Println("已清空对话历史。")
		return MetaResult{Handled: true, ResetHistory: true}
	case "provider":
		if arg == "" {
			return MetaResult{Handled: true, SwitchProvider: "?"}
		}
		return MetaResult{Handled: true, SwitchProvider: arg}
	case "model":
		if arg == "" {
			return MetaResult{Handled: true, SwitchModel: "?"}
		}
		return MetaResult{Handled: true, SwitchModel: arg}
	case "blacklist":
		return MetaResult{Handled: true, ShowBlacklist: true}
	case "config":
		if arg == "dir" {
			return MetaResult{Handled: true, ShowConfigDir: true}
		}
		fmt.Println("用法：/config dir")
		return MetaResult{Handled: true}
	case "help":
		printHelp()
		return MetaResult{Handled: true}
	default:
		fmt.Printf("未知命令：%s（输入 /help 查看帮助）\n", input)
		return MetaResult{Handled: true}
	}
}

func printHelp() {
	help := `
使用方式：直接输入中文描述需求，aict 会生成对应命令并解释。

REPL 命令：
  /help               显示此帮助
  /exit / /quit       退出 aict
  /reset              清空当前对话历史
  /provider           列出所有已配置的 provider
  /provider <name>    切换 AI provider（claude/openai/ollama 等）
  /model              显示当前 provider 和模型
  /model <name>       切换当前 provider 的模型
  /blacklist          列出所有黑名单规则
  /config dir         显示配置目录及运行模式

命令生成后的确认操作：
  y                   执行命令
  N                   取消（默认，直接回车即取消）
  e                   编辑命令后再决定
  r                   让 AI 重新生成命令
  d                   让 AI 给出更详细的解释
  b                   将该命令模式加入黑名单

CLI 子命令（在终端中使用）：
  aict init           重新运行配置向导
  aict add provider   添加新的 AI provider
  aict add model      为已有 provider 设置模型
  aict edit provider  修改已有 provider 的配置
  aict edit model     修改已有 provider 的模型
  aict delete provider 删除已有的 provider
  aict config show    显示当前配置（API Key 脱敏）
  aict version        显示版本号
`
	fmt.Print(help)
}
