package repl

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/chzyer/readline"

	"github.com/aict-tool/aict/internal/config"
	"github.com/aict-tool/aict/internal/logs"
	"github.com/aict-tool/aict/internal/parser"
	"github.com/aict-tool/aict/internal/prompt"
	"github.com/aict-tool/aict/internal/provider"
	"github.com/aict-tool/aict/internal/safety"
	"github.com/aict-tool/aict/internal/shell"
)

type Session struct {
	cfg       *config.Config
	configDir config.ConfigDir
	prov      provider.Provider
	shell     shell.Shell
	blacklist *safety.Blacklist
	recorder  *logs.Recorder
	history   []provider.Message
	systemMsg string
}

func NewSession(cfg *config.Config, cd config.ConfigDir) (*Session, error) {
	name, pc, err := cfg.CurrentProvider()
	if err != nil {
		return nil, err
	}

	p, err := provider.Build(name, pc)
	if err != nil {
		return nil, err
	}

	sh := shell.Detect()
	bl, err := safety.NewBlacklist(cd.Path)
	if err != nil {
		return nil, err
	}

	rec, err := logs.New(cd.Path)
	if err != nil {
		return nil, err
	}

	return &Session{
		cfg:       cfg,
		configDir: cd,
		prov:      p,
		shell:     sh,
		blacklist: bl,
		recorder:  rec,
		systemMsg: prompt.Build(sh),
	}, nil
}

func (s *Session) Run() error {
	rl, err := readline.NewEx(&readline.Config{
		Prompt:          "> ",
		HistoryFile:     s.recorder.Dir() + "/readline",
		InterruptPrompt: "^C",
		EOFPrompt:       "/exit",
	})
	if err != nil {
		return err
	}
	defer rl.Close()

	colorMeta.Printf("aict 已启动（provider: %s，shell: %s）\n", s.prov.Name(), s.shell.Kind)
	colorMeta.Println("输入中文描述需求，或输入 /help 查看帮助")
	fmt.Println()

	for {
		line, err := rl.Readline()
		if err != nil {
			break
		}
		input := strings.TrimSpace(line)
		if input == "" {
			continue
		}

		if meta := HandleMeta(input); meta.Handled {
			if meta.ShouldExit {
				break
			}
			s.applyMeta(meta)
			continue
		}

		if isOffTopic(input) {
			colorMeta.Println("aict 仅用于生成和执行 shell 命令。请输入与命令行操作相关的需求。")
			continue
		}

		if err := s.handleQuery(input); err != nil {
			RenderError(err.Error())
		}
	}
	fmt.Println("再见！")
	return nil
}

func (s *Session) handleQuery(input string) error {
	msgs := append([]provider.Message{
		{Role: "system", Content: s.systemMsg},
	}, s.history...)
	msgs = append(msgs, provider.Message{Role: "user", Content: input})

	ctx := context.Background()
	ch, err := s.prov.Stream(ctx, msgs)
	if err != nil {
		return err
	}

	fmt.Println()
	result := parser.Parse(ch, RenderDelta)
	if result.Err != nil {
		return result.Err
	}
	RenderSeparator()

	if result.Command == "" {
		RenderError("未能提取到命令，请重新描述或使用 /reset 清空历史")
		return nil
	}

	s.recorder.LogUser(input)
	s.recorder.LogAI(result.Explanation, result.Command)

	RenderCommand(result.Command)
	executed := s.handleConfirm(result.Command, msgs, input)

	if executed {
		s.history = append(s.history,
			provider.Message{Role: "user", Content: input},
			provider.Message{Role: "assistant", Content: result.Explanation + "\n\n命令：" + result.Command},
		)
	}
	return nil
}

func (s *Session) handleConfirm(cmd string, msgs []provider.Message, userInput string) bool {
	if rule, hit := s.blacklist.Match(cmd); hit {
		RenderWarning(rule.ID, rule.Reason, rule.Pattern)
		if !AskBlacklistConfirm(cmd) {
			fmt.Println("已取消。")
			return false
		}
		return s.executeCommand(cmd)
	}

	for {
		action := AskConfirm()
		switch action {
		case ActionExecute:
			return s.executeCommand(cmd)

		case ActionCancel:
			fmt.Println("已取消。")
			return false

		case ActionEdit:
			edited := s.editCommand(cmd)
			if edited == "" {
				fmt.Println("编辑取消。")
				return false
			}
			cmd = edited
			RenderCommand(cmd)

		case ActionRegenerate:
			fmt.Println("正在重新生成...")
			regenMsgs := append(msgs, provider.Message{
				Role: "user", Content: userInput + "（上一个方案不合适，请换一种更简洁的方案）",
			})
			ch, err := s.prov.Stream(context.Background(), regenMsgs)
			if err != nil {
				RenderError(err.Error())
				return false
			}
			fmt.Println()
			result := parser.Parse(ch, RenderDelta)
			RenderSeparator()
			if result.Command != "" {
				cmd = result.Command
				RenderCommand(cmd)
			}

		case ActionDetail:
			fmt.Println("正在获取更详细的解释...")
			detailMsgs := append(msgs, provider.Message{
				Role: "user", Content: "请对上面那条命令给出更详细的解释，重点说明每个参数的含义和潜在风险，不需要重新生成命令。",
			})
			ch, err := s.prov.Stream(context.Background(), detailMsgs)
			if err != nil {
				RenderError(err.Error())
			} else {
				fmt.Println()
				parser.Parse(ch, RenderDelta)
				fmt.Println()
			}

		case ActionBlacklist:
			if inp := AskBlacklistAdd(cmd); inp != nil {
				rule := safety.Rule{
					ID:      "user-" + strings.ReplaceAll(strings.ToLower(cmd[:min(20, len(cmd))]), " ", "-"),
					Pattern: inp.Pattern,
					Reason:  inp.Reason,
				}
				if err := s.blacklist.Add(rule); err != nil {
					RenderError("保存黑名单失败：" + err.Error())
				} else {
					fmt.Println("已添加到黑名单。")
				}
			}
			return false
		}
	}
}

func (s *Session) executeCommand(cmd string) bool {
	fmt.Println()
	if err := shell.Execute(s.shell, cmd, nil); err != nil {
		s.recorder.LogExec(cmd, false)
		RenderError("执行失败：" + err.Error())
		return false
	}
	s.recorder.LogExec(cmd, true)
	fmt.Println()
	return true
}

func (s *Session) editCommand(cmd string) string {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = os.Getenv("VISUAL")
	}
	if editor != "" {
		f, err := os.CreateTemp("", "aict-cmd-*.sh")
		if err == nil {
			f.WriteString(cmd)
			f.Close()
			c := exec.Command(editor, f.Name())
			c.Stdin, c.Stdout, c.Stderr = os.Stdin, os.Stdout, os.Stderr
			c.Run()
			data, err := os.ReadFile(f.Name())
			os.Remove(f.Name())
			if err == nil {
				return strings.TrimSpace(string(data))
			}
		}
	}
	fmt.Printf("编辑命令（回车确认）：\n")
	fmt.Printf("[当前] %s\n> ", cmd)
	rl, _ := readline.New("> ")
	if rl != nil {
		defer rl.Close()
		if line, err := rl.ReadlineWithDefault(cmd); err == nil {
			return strings.TrimSpace(line)
		}
	}
	return ""
}

func (s *Session) applyMeta(meta MetaResult) {
	if meta.ResetHistory {
		s.history = nil
		s.recorder.LogReset()
	}
	if meta.SwitchProvider == "?" {
		s.listProviders()
		return
	}
	if meta.SwitchProvider != "" {
		pc, ok := s.cfg.Providers[meta.SwitchProvider]
		if !ok {
			RenderError(fmt.Sprintf("provider %q 未配置", meta.SwitchProvider))
			s.listProviders()
			return
		}
		p, err := provider.Build(meta.SwitchProvider, pc)
		if err != nil {
			RenderError(err.Error())
			return
		}
		s.prov = p
		s.cfg.DefaultProvider = meta.SwitchProvider
		fmt.Printf("已切换到 provider: %s（模型: %s）\n", meta.SwitchProvider, pc.Model)
	}
	if meta.SwitchModel == "?" {
		pc := s.cfg.Providers[s.cfg.DefaultProvider]
		fmt.Printf("当前 provider: %s，模型: %s\n", s.cfg.DefaultProvider, pc.Model)
		fmt.Println("用法：/model <模型名称>  例如 /model gpt-4o-mini")
		return
	}
	if meta.SwitchModel != "" {
		pc := s.cfg.Providers[s.cfg.DefaultProvider]
		pc.Model = meta.SwitchModel
		s.cfg.Providers[s.cfg.DefaultProvider] = pc
		p, err := provider.Build(s.cfg.DefaultProvider, pc)
		if err != nil {
			RenderError(err.Error())
			return
		}
		s.prov = p
		fmt.Printf("已切换模型: %s（provider: %s）\n", meta.SwitchModel, s.cfg.DefaultProvider)
	}
	if meta.ShowBlacklist {
		for _, r := range s.blacklist.List() {
			colorMeta.Printf("[%s] %s\n", r.Source, r.ID)
			fmt.Printf("  模式：%s\n  原因：%s\n", r.Pattern, r.Reason)
		}
	}
	if meta.ShowConfigDir {
		colorMeta.Printf("配置目录：%s\n运行模式：%s\n", s.configDir.Path, s.configDir.Mode)
	}
}

func (s *Session) listProviders() {
	fmt.Println("可用 provider：")
	for name, pc := range s.cfg.Providers {
		mark := "  "
		if name == s.cfg.DefaultProvider {
			mark = "* "
		}
		fmt.Printf("  %s%s（模型: %s）\n", mark, name, pc.Model)
	}
	fmt.Println("用法：/provider <name>  切换 provider")
}
