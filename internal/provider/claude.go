package provider

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type ClaudeProvider struct {
	BaseURL string
	APIKey  string
	Model   string
}

func (p *ClaudeProvider) Name() string { return "claude" }

func (p *ClaudeProvider) Stream(ctx context.Context, msgs []Message) (<-chan Chunk, error) {
	var systemContent string
	var convMsgs []map[string]string
	for _, m := range msgs {
		if m.Role == "system" {
			systemContent = m.Content
			continue
		}
		convMsgs = append(convMsgs, map[string]string{"role": m.Role, "content": m.Content})
	}

	body := map[string]any{
		"model":      p.Model,
		"max_tokens": 4096,
		"stream":     true,
		"messages":   convMsgs,
	}
	if systemContent != "" {
		body["system"] = systemContent
	}

	data, _ := json.Marshal(body)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		p.BaseURL+"/v1/messages", bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", p.APIKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("claude API 返回 %d", resp.StatusCode)
	}

	ch := make(chan Chunk, 8)
	go func() {
		defer close(ch)
		defer resp.Body.Close()
		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			line := scanner.Text()
			if !strings.HasPrefix(line, "data: ") {
				continue
			}
			raw := strings.TrimPrefix(line, "data: ")
			var event struct {
				Type  string `json:"type"`
				Delta struct {
					Type string `json:"type"`
					Text string `json:"text"`
				} `json:"delta"`
			}
			if err := json.Unmarshal([]byte(raw), &event); err != nil {
				continue
			}
			switch event.Type {
			case "content_block_delta":
				if event.Delta.Type == "text_delta" {
					ch <- Chunk{Delta: event.Delta.Text}
				}
			case "message_stop":
				ch <- Chunk{Done: true}
				return
			}
		}
		if err := scanner.Err(); err != nil {
			ch <- Chunk{Err: err}
		}
	}()
	return ch, nil
}
