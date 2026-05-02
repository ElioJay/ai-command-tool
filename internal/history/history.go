package history

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type Recorder struct {
	dir string
}

func New(configDir string) (*Recorder, error) {
	dir := filepath.Join(configDir, "history")
	if info, err := os.Stat(dir); err == nil && !info.IsDir() {
		os.Rename(dir, dir+".old")
	}
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return nil, fmt.Errorf("创建历史目录失败: %w", err)
	}
	return &Recorder{dir: dir}, nil
}

func (r *Recorder) LogUser(input string) {
	r.write("USER", input)
}

func (r *Recorder) LogAI(explanation string, command string) {
	if command != "" {
		r.write("AI", explanation+"\n命令："+command)
	} else {
		r.write("AI", explanation)
	}
}

func (r *Recorder) LogExec(cmd string, success bool) {
	status := "OK"
	if !success {
		status = "FAIL"
	}
	r.write("EXEC["+status+"]", cmd)
}

func (r *Recorder) LogReset() {
	r.write("SYSTEM", "对话历史已清空")
}

func (r *Recorder) Dir() string {
	return r.dir
}

func (r *Recorder) write(tag, content string) {
	now := time.Now()
	filename := filepath.Join(r.dir, now.Format("2006-01-02")+".log")

	f, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o600)
	if err != nil {
		return
	}
	defer f.Close()

	fmt.Fprintf(f, "[%s] [%s] %s\n", now.Format("15:04:05"), tag, content)
}
