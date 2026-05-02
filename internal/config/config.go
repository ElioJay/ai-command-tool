package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
)

type ProviderConfig struct {
	BaseURL string `toml:"base_url"`
	APIKey  string `toml:"api_key"`
	Model   string `toml:"model"`
}

type UIConfig struct {
	Stream bool   `toml:"stream"`
	Color  string `toml:"color"`
}

type Config struct {
	DefaultProvider string                    `toml:"default_provider"`
	Providers       map[string]ProviderConfig `toml:"providers"`
	UI              UIConfig                  `toml:"ui"`
}

func Load(configDir string) (*Config, error) {
	path := filepath.Join(configDir, "config.toml")
	cfg := &Config{
		Providers: make(map[string]ProviderConfig),
		UI:        UIConfig{Stream: true, Color: "auto"},
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return cfg, nil
	}

	if _, err := toml.DecodeFile(path, cfg); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}
	applyEnvOverrides(cfg)
	return cfg, nil
}

func Save(cfg *Config, configDir string) error {
	if err := os.MkdirAll(configDir, 0o700); err != nil {
		return err
	}
	path := filepath.Join(configDir, "config.toml")
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o600)
	if err != nil {
		return err
	}
	defer f.Close()
	return toml.NewEncoder(f).Encode(cfg)
}

func applyEnvOverrides(cfg *Config) {
	if v := os.Getenv("AICT_PROVIDER"); v != "" {
		cfg.DefaultProvider = v
	}
	for name, pc := range cfg.Providers {
		upper := strings.ToUpper(name)
		if v := os.Getenv("AICT_" + upper + "_API_KEY"); v != "" {
			pc.APIKey = v
			cfg.Providers[name] = pc
		}
		if v := os.Getenv("AICT_" + upper + "_MODEL"); v != "" {
			pc.Model = v
			cfg.Providers[name] = pc
		}
		if v := os.Getenv("AICT_" + upper + "_BASE_URL"); v != "" {
			pc.BaseURL = v
			cfg.Providers[name] = pc
		}
	}
}

func (c *Config) CurrentProvider() (string, ProviderConfig, error) {
	pc, ok := c.Providers[c.DefaultProvider]
	if !ok {
		return "", ProviderConfig{}, fmt.Errorf("provider %q 未配置", c.DefaultProvider)
	}
	return c.DefaultProvider, pc, nil
}
