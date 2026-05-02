package cmdlog

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestLog_CreatesFileAndWritesEntry(t *testing.T) {
	dir := t.TempDir()
	l, err := New(dir)
	if err != nil {
		t.Fatal(err)
	}

	l.Log("ls -la", true)
	l.Log("rm -rf /tmp/test", false)

	today := time.Now().Format("2006-01-02")
	data, err := os.ReadFile(filepath.Join(dir, "log", today+".log"))
	if err != nil {
		t.Fatal(err)
	}

	content := string(data)
	if !strings.Contains(content, "[OK] ls -la") {
		t.Errorf("missing OK entry, got:\n%s", content)
	}
	if !strings.Contains(content, "[FAIL] rm -rf /tmp/test") {
		t.Errorf("missing FAIL entry, got:\n%s", content)
	}

	lines := strings.Split(strings.TrimSpace(content), "\n")
	if len(lines) != 2 {
		t.Errorf("expected 2 lines, got %d", len(lines))
	}
}

func TestNew_CreatesLogDir(t *testing.T) {
	dir := t.TempDir()
	l, err := New(dir)
	if err != nil {
		t.Fatal(err)
	}

	info, err := os.Stat(l.Dir())
	if err != nil {
		t.Fatal(err)
	}
	if !info.IsDir() {
		t.Error("log dir is not a directory")
	}
}
