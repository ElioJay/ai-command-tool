package shell

import (
	"testing"
)

func TestDetectFromEnv_Bash(t *testing.T) {
	t.Setenv("SHELL", "/bin/bash")
	t.Setenv("TERM_PROGRAM", "")
	sh := detectFromEnv()
	if sh.Kind != Bash {
		t.Errorf("expected Bash, got %s", sh.Kind)
	}
}

func TestDetectFromEnv_Zsh(t *testing.T) {
	t.Setenv("SHELL", "/usr/bin/zsh")
	sh := detectFromEnv()
	if sh.Kind != Zsh {
		t.Errorf("expected Zsh, got %s", sh.Kind)
	}
}

func TestShellKindString(t *testing.T) {
	cases := map[Kind]string{
		Pwsh: "pwsh", Bash: "bash", Zsh: "zsh",
		Fish: "fish", Cmd: "cmd", Sh: "sh",
	}
	for k, want := range cases {
		if got := k.String(); got != want {
			t.Errorf("Kind(%d).String() = %q, want %q", k, got, want)
		}
	}
}
