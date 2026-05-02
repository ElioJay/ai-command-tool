package repl

import "testing"

func TestIsOffTopic(t *testing.T) {
	offTopic := []string{
		"你好",
		"你是谁",
		"讲个笑话",
		"帮我翻译这段英文",
		"写一篇关于AI的文章",
		"hello",
		"1+1等于几",
		"嗯",
		"ok",
	}
	for _, s := range offTopic {
		if !isOffTopic(s) {
			t.Errorf("expected off-topic: %q", s)
		}
	}

	onTopic := []string{
		"查找当前目录下所有 .go 文件",
		"查看占用 8080 端口的进程",
		"列出最近修改的 10 个文件",
		"git log 最近5条提交",
		"ping baidu.com",
		"压缩 src 目录为 zip",
		"查看磁盘空间",
		"删除 tmp 目录下的所有日志文件",
		"显示当前目录的文件大小",
		"docker ps 看看运行中的容器",
	}
	for _, s := range onTopic {
		if isOffTopic(s) {
			t.Errorf("expected on-topic: %q", s)
		}
	}
}
