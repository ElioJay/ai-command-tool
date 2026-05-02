package config

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type WizardResult struct {
	ConfigDir ConfigDir
	Config    *Config
}

var providerDefaults = map[string]ProviderConfig{
	"claude":   {BaseURL: "https://api.anthropic.com", Model: "claude-sonnet-4-6"},
	"openai":   {BaseURL: "https://api.openai.com/v1", Model: "gpt-4o"},
	"deepseek": {BaseURL: "https://api.deepseek.com/v1", Model: "deepseek-chat"},
	"ollama":   {BaseURL: "http://localhost:11434/v1", Model: "qwen2.5-coder:7b"},
}

func RunWizard() (*WizardResult, error) {
	scanner := bufio.NewScanner(os.Stdin)
	readLine := func(prompt, defaultVal string) string {
		if defaultVal != "" {
			fmt.Printf("%s（回车使用默认 %s）：\n> ", prompt, defaultVal)
		} else {
			fmt.Printf("%s：\n> ", prompt)
		}
		scanner.Scan()
		val := strings.TrimSpace(scanner.Text())
		if val == "" {
			return defaultVal
		}
		return val
	}

	fmt.Println("欢迎使用 aict！我们来完成首次配置。")
	fmt.Println()

	fmt.Println("[1/5] 选择默认 AI provider：")
	fmt.Println("  1) Claude (官方)")
	fmt.Println("  2) OpenAI (官方)")
	fmt.Println("  3) DeepSeek")
	fmt.Println("  4) Ollama (本地)")
	fmt.Println("  5) 自定义 OpenAI 兼容接口")
	fmt.Print("> ")
	scanner.Scan()
	choice := strings.TrimSpace(scanner.Text())

	providerName := "claude"
	switch choice {
	case "2":
		providerName = "openai"
	case "3":
		providerName = "deepseek"
	case "4":
		providerName = "ollama"
	case "5":
		providerName = readLine("Provider 名称（如 my-api）", "custom")
	}

	defaults := providerDefaults[providerName]

	baseURL := readLine(fmt.Sprintf("[2/5] Base URL"), defaults.BaseURL)

	apiKey := ""
	if providerName != "ollama" {
		apiKey = readLine("[3/5] API Key", "")
		if apiKey == "" {
			fmt.Println("警告：未设置 API Key，后续可通过 AICT_<PROVIDER>_API_KEY 环境变量设置。")
		}
	} else {
		fmt.Println("[3/5] Ollama 本地模式，无需 API Key。")
	}

	model := readLine("[4/5] 模型名称", defaults.Model)

	fmt.Println("\n[5/5] 配置文件保存位置：")
	exePath, _ := os.Executable()
	exeDir := ""
	if exePath != "" {
		exeDir = exePath + "/../.aict（便携版）"
	}
	fmt.Printf("  1) 用户目录 ~/.aict（安装版，推荐）\n")
	if exeDir != "" {
		fmt.Printf("  2) 可执行文件同目录 .aict/（便携版）\n")
	}
	fmt.Print("> ")
	scanner.Scan()
	modeChoice := strings.TrimSpace(scanner.Text())

	var cd ConfigDir
	if modeChoice == "2" && exeDir != "" {
		portablePath := exePath + "/../.aict"
		if err := os.MkdirAll(portablePath, 0o700); err != nil {
			return nil, fmt.Errorf("创建便携配置目录失败: %w", err)
		}
		cd = ConfigDir{Path: portablePath, Mode: ModePortable}
	} else {
		installed := installedPath()
		if err := os.MkdirAll(installed, 0o700); err != nil {
			return nil, fmt.Errorf("创建配置目录失败: %w", err)
		}
		cd = ConfigDir{Path: installed, Mode: ModeInstalled}
	}

	cfg := &Config{
		DefaultProvider: providerName,
		Providers: map[string]ProviderConfig{
			providerName: {
				BaseURL: baseURL,
				APIKey:  apiKey,
				Model:   model,
			},
		},
		UI: UIConfig{Stream: true, Color: "auto"},
	}

	if err := Save(cfg, cd.Path); err != nil {
		return nil, fmt.Errorf("保存配置失败: %w", err)
	}
	fmt.Printf("\n配置已保存到 %s\n", cd.Path)
	return &WizardResult{ConfigDir: cd, Config: cfg}, nil
}
