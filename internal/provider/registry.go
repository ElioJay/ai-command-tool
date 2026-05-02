package provider

import (
	"fmt"

	"github.com/aict-tool/aict/internal/config"
)

func Build(name string, pc config.ProviderConfig) (Provider, error) {
	switch name {
	case "claude":
		baseURL := pc.BaseURL
		if baseURL == "" {
			baseURL = "https://api.anthropic.com"
		}
		return &ClaudeProvider{BaseURL: baseURL, APIKey: pc.APIKey, Model: pc.Model}, nil
	case "openai":
		baseURL := pc.BaseURL
		if baseURL == "" {
			baseURL = "https://api.openai.com"
		}
		return &OpenAIProvider{BaseURL: baseURL, APIKey: pc.APIKey, Model: pc.Model, ProviderID: "openai"}, nil
	default:
		if pc.BaseURL == "" {
			return nil, fmt.Errorf("provider %q 需要配置 base_url", name)
		}
		return &OpenAIProvider{
			BaseURL:    pc.BaseURL,
			APIKey:     pc.APIKey,
			Model:      pc.Model,
			ProviderID: name,
		}, nil
	}
}
