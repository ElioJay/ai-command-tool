package provider

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestClaudeStream(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("x-api-key") == "" {
			t.Error("missing x-api-key header")
		}
		w.Header().Set("Content-Type", "text/event-stream")
		events := []string{
			`{"type":"content_block_delta","delta":{"type":"text_delta","text":"Hello"}}`,
			`{"type":"content_block_delta","delta":{"type":"text_delta","text":" world"}}`,
			`{"type":"message_stop"}`,
		}
		for _, e := range events {
			fmt.Fprintf(w, "data: %s\n\n", e)
		}
	}))
	defer srv.Close()

	p := &ClaudeProvider{
		BaseURL: srv.URL,
		APIKey:  "test-key",
		Model:   "claude-test",
	}
	ch, err := p.Stream(context.Background(), []Message{
		{Role: "user", Content: "hi"},
	})
	if err != nil {
		t.Fatal(err)
	}
	var sb strings.Builder
	for c := range ch {
		if c.Err != nil {
			t.Fatal(c.Err)
		}
		sb.WriteString(c.Delta)
	}
	if got := sb.String(); got != "Hello world" {
		t.Errorf("got %q", got)
	}
}
