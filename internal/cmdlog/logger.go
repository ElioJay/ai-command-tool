package cmdlog

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type Logger struct {
	dir string
}

func New(configDir string) (*Logger, error) {
	dir := filepath.Join(configDir, "log")
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return nil, fmt.Errorf("创建日志目录失败: %w", err)
	}
	return &Logger{dir: dir}, nil
}

func (l *Logger) Log(cmd string, success bool) {
	now := time.Now()
	filename := filepath.Join(l.dir, now.Format("2006-01-02")+".log")

	f, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o600)
	if err != nil {
		return
	}
	defer f.Close()

	status := "OK"
	if !success {
		status = "FAIL"
	}
	fmt.Fprintf(f, "[%s] [%s] %s\n", now.Format("15:04:05"), status, cmd)
}

func (l *Logger) Dir() string {
	return l.dir
}
