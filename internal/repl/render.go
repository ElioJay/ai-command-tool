package repl

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

var (
	colorExplain = color.New(color.FgCyan)
	colorCommand = color.New(color.FgGreen, color.Bold)
	colorWarn    = color.New(color.FgRed, color.Bold)
	colorMeta    = color.New(color.FgYellow)
	colorPrompt  = color.New(color.FgWhite, color.Bold)
)

func RenderDelta(delta string) {
	fmt.Print(delta)
}

func RenderError(msg string) {
	colorWarn.Fprintln(os.Stderr, "错误："+msg)
}

func RenderWarning(ruleID, reason, pattern string) {
	colorWarn.Println("\n[黑名单] 命令命中规则: " + ruleID)
	fmt.Println("  原因：" + reason)
	fmt.Println("  规则模式：" + pattern)
}

func RenderSeparator() {
	colorMeta.Println("\n" + "─────────────────────────────────")
}

func RenderCommand(cmd string) {
	fmt.Println()
	colorCommand.Println("$ " + cmd)
}
