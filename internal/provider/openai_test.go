package provider

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestOpenAIStream(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.Header.Get("Authorization"), "Bearer ") {
			t.Error("missing Bearer token")
		}
		w.Header().Set("Content-Type", "text/event-stream")
		chunks := []string{"## 解释\n", "这是测试", "\n## 命令\n```bash\necho hi\n```"}
		for _, c := range chunks {
			payload := fmt.Sprintf(`{"choices":[{"delta":{"content":%q}}]}`, c)
			fmt.Fprintf(w, "data: %s\n\n", payload)
		}
		fmt.Fprintln(w, "data: [DONE]")
	}))
	defer srv.Close()

	p := &OpenAIProvider{BaseURL: srv.URL, APIKey: "test", Model: "gpt-test"}
	ch, err := p.Stream(context.Background(), []Message{
		{Role: "system", Content: "sys"},
		{Role: "user", Content: "hello"},
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
	if !strings.Contains(sb.String(), "echo hi") {
		t.Errorf("unexpected output: %s", sb.String())
	}
}
