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

func newScanner() (*bufio.Scanner, func(prompt, defaultVal string) string) {
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
	return scanner, readLine
}

func askProvider(scanner *bufio.Scanner, readLine func(string, string) string) (name string, pc ProviderConfig) {
	fmt.Println("选择 AI provider：")
	fmt.Println("  1) Claude (官方)")
	fmt.Println("  2) OpenAI (官方)")
	fmt.Println("  3) DeepSeek")
	fmt.Println("  4) Ollama (本地)")
	fmt.Println("  5) 自定义 OpenAI 兼容接口")
	fmt.Print("> ")
	scanner.Scan()
	choice := strings.TrimSpace(scanner.Text())

	name = "claude"
	switch choice {
	case "2":
		name = "openai"
	case "3":
		name = "deepseek"
	case "4":
		name = "ollama"
	case "5":
		name = readLine("Provider 名称（如 my-api）", "custom")
	}

	defaults := providerDefaults[name]
	pc.BaseURL = readLine("Base URL", defaults.BaseURL)

	if name != "ollama" {
		pc.APIKey = readLine("API Key", "")
		if pc.APIKey == "" {
			fmt.Println("警告：未设置 API Key，后续可通过 AICT_<PROVIDER>_API_KEY 环境变量设置。")
		}
	} else {
		fmt.Println("Ollama 本地模式，无需 API Key。")
	}

	pc.Model = readLine("模型名称", defaults.Model)
	return name, pc
}

func RunWizard() (*WizardResult, error) {
	scanner, readLine := newScanner()

	fmt.Println("欢迎使用 aict！我们来完成首次配置。")
	fmt.Println()

	fmt.Println("[1/4] ", "")
	providerName, pc := askProvider(scanner, readLine)

	fmt.Println("\n[2/4] 配置文件保存位置：")
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
			providerName: pc,
		},
		UI: UIConfig{Stream: true, Color: "auto"},
	}

	if err := Save(cfg, cd.Path); err != nil {
		return nil, fmt.Errorf("保存配置失败: %w", err)
	}
	fmt.Printf("\n配置已保存到 %s\n", cd.Path)
	return &WizardResult{ConfigDir: cd, Config: cfg}, nil
}

func AddProvider(cd ConfigDir) (*Config, error) {
	cfg, err := Load(cd.Path)
	if err != nil {
		return nil, err
	}

	scanner, readLine := newScanner()

	fmt.Println("添加新的 AI provider")
	fmt.Println()

	name, pc := askProvider(scanner, readLine)

	if _, exists := cfg.Providers[name]; exists {
		fmt.Printf("provider %q 已存在，是否覆盖？[y/N] ", name)
		scanner.Scan()
		ans := strings.TrimSpace(scanner.Text())
		if ans != "y" && ans != "Y" {
			fmt.Println("已取消。")
			return cfg, nil
		}
	}

	cfg.Providers[name] = pc

	fmt.Printf("是否将 %s 设为默认 provider？[y/N] ", name)
	scanner.Scan()
	ans := strings.TrimSpace(scanner.Text())
	if ans == "y" || ans == "Y" {
		cfg.DefaultProvider = name
	}

	if err := Save(cfg, cd.Path); err != nil {
		return nil, fmt.Errorf("保存配置失败: %w", err)
	}
	fmt.Printf("\nprovider %q 已添加。\n", name)
	printProviders(cfg)
	return cfg, nil
}

func pickProvider(cfg *Config, scanner *bufio.Scanner, prompt string) string {
	fmt.Println(prompt)
	names := make([]string, 0, len(cfg.Providers))
	i := 1
	for n := range cfg.Providers {
		mark := ""
		if n == cfg.DefaultProvider {
			mark = "（默认）"
		}
		fmt.Printf("  %d) %s%s  [模型: %s]\n", i, n, mark, cfg.Providers[n].Model)
		names = append(names, n)
		i++
	}
	fmt.Print("> ")
	scanner.Scan()
	choice := strings.TrimSpace(scanner.Text())

	for idx, n := range names {
		if choice == fmt.Sprintf("%d", idx+1) || choice == n {
			return n
		}
	}
	if len(names) == 1 {
		fmt.Printf("仅有一个 provider，已自动选择：%s\n", names[0])
		return names[0]
	}
	return ""
}

func AddModel(cd ConfigDir) (*Config, error) {
	cfg, err := Load(cd.Path)
	if err != nil {
		return nil, err
	}
	if len(cfg.Providers) == 0 {
		return nil, fmt.Errorf("尚未配置任何 provider，请先运行 aict init 或 aict add provider")
	}

	scanner, readLine := newScanner()

	fmt.Println("为 provider 设置模型")
	fmt.Println()

	targetName := pickProvider(cfg, scanner, "选择要修改模型的 provider：")
	if targetName == "" {
		fmt.Println("已取消。")
		return cfg, nil
	}

	pc := cfg.Providers[targetName]
	newModel := readLine(fmt.Sprintf("新模型名称（当前: %s）", pc.Model), "")
	if newModel == "" {
		fmt.Println("未输入模型名称，已取消。")
		return cfg, nil
	}

	pc.Model = newModel
	cfg.Providers[targetName] = pc

	if err := Save(cfg, cd.Path); err != nil {
		return nil, fmt.Errorf("保存配置失败: %w", err)
	}
	fmt.Printf("\n%s 的模型已更新为 %s\n", targetName, newModel)
	printProviders(cfg)
	return cfg, nil
}

func EditProvider(cd ConfigDir) (*Config, error) {
	cfg, err := Load(cd.Path)
	if err != nil {
		return nil, err
	}
	if len(cfg.Providers) == 0 {
		return nil, fmt.Errorf("尚未配置任何 provider，请先运行 aict init 或 aict add provider")
	}

	scanner, readLine := newScanner()

	targetName := pickProvider(cfg, scanner, "选择要修改的 provider：")
	if targetName == "" {
		fmt.Println("已取消。")
		return cfg, nil
	}

	pc := cfg.Providers[targetName]
	fmt.Printf("\n正在编辑 provider: %s\n", targetName)
	fmt.Println("（直接回车保持当前值不变）")

	pc.BaseURL = readLine(fmt.Sprintf("Base URL（当前: %s）", pc.BaseURL), pc.BaseURL)

	maskedKey := "（未设置）"
	if len(pc.APIKey) > 8 {
		maskedKey = pc.APIKey[:4] + "..." + pc.APIKey[len(pc.APIKey)-4:]
	} else if pc.APIKey != "" {
		maskedKey = "***"
	}
	newKey := readLine(fmt.Sprintf("API Key（当前: %s）", maskedKey), "")
	if newKey != "" {
		pc.APIKey = newKey
	}

	pc.Model = readLine(fmt.Sprintf("模型名称（当前: %s）", pc.Model), pc.Model)

	cfg.Providers[targetName] = pc

	if err := Save(cfg, cd.Path); err != nil {
		return nil, fmt.Errorf("保存配置失败: %w", err)
	}
	fmt.Printf("\nprovider %q 已更新。\n", targetName)
	printProviders(cfg)
	return cfg, nil
}

func DeleteProvider(cd ConfigDir) (*Config, error) {
	cfg, err := Load(cd.Path)
	if err != nil {
		return nil, err
	}
	if len(cfg.Providers) == 0 {
		return nil, fmt.Errorf("尚未配置任何 provider")
	}
	if len(cfg.Providers) == 1 {
		return nil, fmt.Errorf("仅剩一个 provider，无法删除。请先添加其他 provider")
	}

	scanner, _ := newScanner()

	targetName := pickProvider(cfg, scanner, "选择要删除的 provider：")
	if targetName == "" {
		fmt.Println("已取消。")
		return cfg, nil
	}

	fmt.Printf("确定要删除 provider %q？[y/N] ", targetName)
	scanner.Scan()
	ans := strings.TrimSpace(scanner.Text())
	if ans != "y" && ans != "Y" {
		fmt.Println("已取消。")
		return cfg, nil
	}

	delete(cfg.Providers, targetName)

	if cfg.DefaultProvider == targetName {
		for n := range cfg.Providers {
			cfg.DefaultProvider = n
			fmt.Printf("默认 provider 已自动切换为：%s\n", n)
			break
		}
	}

	if err := Save(cfg, cd.Path); err != nil {
		return nil, fmt.Errorf("保存配置失败: %w", err)
	}
	fmt.Printf("provider %q 已删除。\n", targetName)
	printProviders(cfg)
	return cfg, nil
}

func printProviders(cfg *Config) {
	fmt.Print("当前已配置的 provider：")
	for n, pc := range cfg.Providers {
		mark := ""
		if n == cfg.DefaultProvider {
			mark = "*"
		}
		fmt.Printf(" %s%s(%s)", n, mark, pc.Model)
	}
	fmt.Println()
}
