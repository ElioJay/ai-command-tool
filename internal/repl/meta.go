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
可用命令：
  /exit / /quit       退出 aict
  /reset              清空当前对话历史
  /provider           列出所有已配置的 provider
  /provider <name>    切换 AI provider（claude/openai/ollama 等）
  /model              显示当前 provider 和模型
  /model <name>       切换当前 provider 的模型
  /blacklist          列出所有黑名单规则
  /config dir         显示配置目录及运行模式
  /help               显示此帮助

提示：直接输入中文描述需求，aict 会生成对应命令并解释。
`
	fmt.Print(help)
}
