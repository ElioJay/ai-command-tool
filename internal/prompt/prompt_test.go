package prompt

import (
	"strings"
	"testing"

	"github.com/aict-tool/aict/internal/shell"
)

func TestBuild_ContainsShellInfo(t *testing.T) {
	sh := shell.Shell{Kind: shell.Bash, Path: "/bin/bash"}
	p := Build(sh)
	if !strings.Contains(p, "bash") {
		t.Error("prompt missing shell kind")
	}
	if !strings.Contains(p, "/bin/bash") {
		t.Error("prompt missing shell path")
	}
	if !strings.Contains(p, "## 解释") {
		t.Error("prompt missing format example")
	}
}
