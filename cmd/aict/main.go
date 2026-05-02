package main

import (
	"fmt"
	"os"

	"github.com/aict-tool/aict/internal/config"
	"github.com/aict-tool/aict/internal/repl"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "错误：%v\n", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	if len(args) > 0 {
		switch args[0] {
		case "init":
			return runInit()
		case "config":
			return runConfig(args[1:])
		case "version", "--version", "-v":
			fmt.Println("aict v0.1.0")
			return nil
		case "help", "--help", "-h":
			printUsage()
			return nil
		}
	}
	return runREPL()
}

func runREPL() error {
	cd := config.Resolve()

	cfg, err := config.Load(cd.Path)
	if err != nil {
		return err
	}
	if cfg.DefaultProvider == "" || len(cfg.Providers) == 0 {
		result, err := config.RunWizard()
		if err != nil {
			return err
		}
		cfg = result.Config
		cd = result.ConfigDir
	}

	session, err := repl.NewSession(cfg, cd)
	if err != nil {
		return err
	}
	return session.Run()
}

func runInit() error {
	result, err := config.RunWizard()
	if err != nil {
		return err
	}
	fmt.Printf("配置完成，模式：%s\n是否立即进入对话？[Y/n] ", result.ConfigDir.Mode)
	var ans string
	fmt.Scanln(&ans)
	if ans == "n" || ans == "N" {
		return nil
	}
	session, err := repl.NewSession(result.Config, result.ConfigDir)
	if err != nil {
		return err
	}
	return session.Run()
}

func runConfig(args []string) error {
	if len(args) == 0 {
		printUsage()
		return nil
	}
	switch args[0] {
	case "show":
		cd := config.Resolve()
		cfg, err := config.Load(cd.Path)
		if err != nil {
			return err
		}
		fmt.Printf("运行模式：%s\n", cd.Mode)
		fmt.Printf("配置目录：%s\n", cd.Path)
		fmt.Printf("默认 Provider：%s\n", cfg.DefaultProvider)
		for name, pc := range cfg.Providers {
			masked := maskAPIKey(pc.APIKey)
			fmt.Printf("  [%s] base_url=%s model=%s api_key=%s\n",
				name, pc.BaseURL, pc.Model, masked)
		}
	default:
		return fmt.Errorf("未知 config 子命令：%s", args[0])
	}
	return nil
}

func maskAPIKey(key string) string {
	if len(key) <= 8 {
		return "***"
	}
	return key[:4] + "..." + key[len(key)-4:]
}

func printUsage() {
	fmt.Print(`
用法：aict [命令]

命令：
  （无参数）     进入交互式 REPL
  init           重新运行配置向导
  config show    显示当前配置
  version        显示版本
  help           显示帮助

REPL 内元命令：
  :exit          退出
  :reset         清空对话历史
  :provider <x>  切换 provider
  :model <x>     切换模型
  :blacklist     列出黑名单
  :config dir    显示配置目录
  :help          REPL 帮助
`)
}
