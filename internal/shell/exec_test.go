package shell

import (
	"testing"
)

func TestExecute_Echo(t *testing.T) {
	sh := Detect()
	var cmd string
	switch sh.Kind {
	case Pwsh:
		cmd = `Write-Output "hello-aict"`
	case Cmd:
		cmd = `echo hello-aict`
	default:
		cmd = `echo hello-aict`
	}
	if err := Execute(sh, cmd, nil); err != nil {
		t.Errorf("Execute returned error: %v", err)
	}
}
