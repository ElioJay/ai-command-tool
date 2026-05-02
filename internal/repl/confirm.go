package repl

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Action int

const (
	ActionExecute    Action = iota
	ActionCancel
	ActionEdit
	ActionRegenerate
	ActionDetail
	ActionBlacklist
)

func AskConfirm() Action {
	colorPrompt.Print("\n[y] 执行  [N] 取消  [e] 编辑  [r] 重新生成  [d] 更详细  [b] 加黑名单\n> ")
	return readAction()
}

func AskBlacklistConfirm(cmd string) bool {
	fmt.Println("\n如确认执行，请逐字输入命令（直接回车取消）：")
	colorPrompt.Print("> ")
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		input := scanner.Text()
		return input == cmd
	}
	return false
}

func AskBlacklistAdd(cmd string) *blacklistInput {
	fmt.Println("\n请确认黑名单 pattern（可改为正则，留空使用逐字匹配）：")
	defaultPattern := "^\\s*" + escapeRegex(strings.TrimSpace(cmd)) + "\\b"
	colorMeta.Printf("[默认: %s]\n", defaultPattern)
	colorPrompt.Print("> ")

	scanner := bufio.NewScanner(os.Stdin)
	pattern := defaultPattern
	if scanner.Scan() && scanner.Text() != "" {
		pattern = scanner.Text()
	}

	fmt.Print("原因（可选，回车跳过）：\n> ")
	reason := ""
	if scanner.Scan() {
		reason = scanner.Text()
	}

	return &blacklistInput{Pattern: pattern, Reason: reason}
}

type blacklistInput struct {
	Pattern string
	Reason  string
}

func readAction() Action {
	buf := make([]byte, 4)
	n, _ := os.Stdin.Read(buf)
	if n == 0 {
		return ActionCancel
	}
	switch strings.ToLower(strings.TrimSpace(string(buf[:n]))) {
	case "y":
		return ActionExecute
	case "e":
		return ActionEdit
	case "r":
		return ActionRegenerate
	case "d":
		return ActionDetail
	case "b":
		return ActionBlacklist
	default:
		return ActionCancel
	}
}

func escapeRegex(s string) string {
	replacer := strings.NewReplacer(
		`.`, `\.`, `*`, `\*`, `+`, `\+`, `?`, `\?`,
		`(`, `\(`, `)`, `\)`, `[`, `\[`, `]`, `\]`,
		`{`, `\{`, `}`, `\}`, `^`, `\^`, `$`, `\$`,
		`|`, `\|`, `\`, `\\`,
	)
	return replacer.Replace(s)
}
