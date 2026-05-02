package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolveConfigDir_Portable(t *testing.T) {
	tmp := t.TempDir()
	os.Mkdir(filepath.Join(tmp, ".aict"), 0o755)
	got := resolveFromDir(tmp)
	if got.Mode != ModePortable {
		t.Errorf("want portable, got %s", got.Mode)
	}
	if got.Path != filepath.Join(tmp, ".aict") {
		t.Errorf("path mismatch: %s", got.Path)
	}
}

func TestResolveConfigDir_Installed(t *testing.T) {
	tmp := t.TempDir()
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)
	got := resolveFromDir(tmp)
	if got.Mode != ModeInstalled {
		t.Errorf("want installed, got %s", got.Mode)
	}
}
