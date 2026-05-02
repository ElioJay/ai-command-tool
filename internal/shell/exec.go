package shell

import (
	"io"
	"os"
	"os/exec"
)

func Execute(sh Shell, cmd string, extraEnv []string) error {
	var c *exec.Cmd
	switch sh.Kind {
	case Pwsh:
		c = exec.Command(sh.Path, "-NoProfile", "-Command", cmd)
	case Cmd:
		c = exec.Command(sh.Path, "/C", cmd)
	default:
		c = exec.Command(sh.Path, "-c", cmd)
	}

	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	if extraEnv != nil {
		c.Env = append(os.Environ(), extraEnv...)
	}
	return c.Run()
}

func ExecuteCapture(sh Shell, cmd string, stdout io.Writer) error {
	var c *exec.Cmd
	switch sh.Kind {
	case Pwsh:
		c = exec.Command(sh.Path, "-NoProfile", "-Command", cmd)
	case Cmd:
		c = exec.Command(sh.Path, "/C", cmd)
	default:
		c = exec.Command(sh.Path, "-c", cmd)
	}
	c.Stdout = stdout
	c.Stderr = stdout
	return c.Run()
}
