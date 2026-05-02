package logs

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestRecorder(t *testing.T) {
	dir := t.TempDir()
	r, err := New(dir)
	if err != nil {
		t.Fatal(err)
	}

	r.LogUser("查看磁盘空间")
	r.LogAI("使用 df -h 查看磁盘使用情况", "df -h")
	r.LogExec("df -h", true)
	r.LogExec("bad-cmd", false)
	r.LogReset()

	today := time.Now().Format("2006-01-02")
	data, err := os.ReadFile(filepath.Join(r.Dir(), today+".log"))
	if err != nil {
		t.Fatal(err)
	}
	content := string(data)

	checks := []string{
		"[USER] 查看磁盘空间",
		"[AI] 使用 df -h 查看磁盘使用情况 | 命令：df -h",
		"[EXEC[OK]] df -h",
		"[EXEC[FAIL]] bad-cmd",
		"[SYSTEM] 对话历史已清空",
	}
	for _, c := range checks {
		if !strings.Contains(content, c) {
			t.Errorf("missing %q in log:\n%s", c, content)
		}
	}
}
