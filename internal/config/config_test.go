package config

import (
	"os"
	"path/filepath"
	"testing"
)

const testTOML = `
default_provider = "claude"

[providers.claude]
base_url = "https://api.anthropic.com"
api_key  = "sk-test"
model    = "claude-sonnet-4-6"

[ui]
stream = true
color  = "auto"
`

func TestLoadConfig(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.toml")
	os.WriteFile(path, []byte(testTOML), 0o600)

	cfg, err := Load(dir)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.DefaultProvider != "claude" {
		t.Errorf("DefaultProvider = %q", cfg.DefaultProvider)
	}
	p, ok := cfg.Providers["claude"]
	if !ok {
		t.Fatal("missing claude provider")
	}
	if p.APIKey != "sk-test" {
		t.Errorf("APIKey = %q", p.APIKey)
	}
}

func TestEnvOverride(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "config.toml"), []byte(testTOML), 0o600)
	t.Setenv("AICT_PROVIDER", "openai")
	t.Setenv("AICT_CLAUDE_API_KEY", "sk-override")

	cfg, _ := Load(dir)
	if cfg.DefaultProvider != "openai" {
		t.Errorf("env AICT_PROVIDER not applied: %s", cfg.DefaultProvider)
	}
	if cfg.Providers["claude"].APIKey != "sk-override" {
		t.Errorf("env AICT_CLAUDE_API_KEY not applied")
	}
}
