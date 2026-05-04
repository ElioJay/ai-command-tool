package shell

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

type Kind int

const (
	Sh Kind = iota
	Bash
	Zsh
	Fish
	Pwsh
	Cmd
)

func (k Kind) String() string {
	return [...]string{"sh", "bash", "zsh", "fish", "pwsh", "cmd"}[k]
}

type Shell struct {
	Kind Kind
	Path string
}

func Detect() Shell {
	return DetectFromEnv()
}

func DetectFromEnv() Shell {
	if runtime.GOOS == "windows" {
		return detectWindows()
	}
	return detectUnix()
}

func detectUnix() Shell {
	shellEnv := os.Getenv("SHELL")
	if shellEnv == "" {
		return Shell{Kind: Sh, Path: "/bin/sh"}
	}
	base := strings.ToLower(filepath.Base(shellEnv))
	switch base {
	case "bash":
		return Shell{Kind: Bash, Path: shellEnv}
	case "zsh":
		return Shell{Kind: Zsh, Path: shellEnv}
	case "fish":
		return Shell{Kind: Fish, Path: shellEnv}
	case "pwsh":
		return Shell{Kind: Pwsh, Path: shellEnv}
	default:
		return Shell{Kind: Sh, Path: shellEnv}
	}
}

func detectWindows() Shell {
	if sh := os.Getenv("SHELL"); sh != "" {
		return detectUnix()
	}
	if os.Getenv("PSModulePath") != "" {
		if pwsh, err := exec.LookPath("pwsh"); err == nil {
			return Shell{Kind: Pwsh, Path: pwsh}
		}
		if pwsh, err := exec.LookPath("powershell"); err == nil {
			return Shell{Kind: Pwsh, Path: pwsh}
		}
	}
	if pwsh, err := exec.LookPath("powershell"); err == nil {
		return Shell{Kind: Pwsh, Path: pwsh}
	}
	return Shell{Kind: Cmd, Path: "cmd.exe"}
}
