package parser

import (
	"strings"
	"testing"

	"github.com/aict-tool/aict/internal/provider"
)

func makeChunks(text string) <-chan provider.Chunk {
	ch := make(chan provider.Chunk, len(text)+1)
	for _, r := range text {
		ch <- provider.Chunk{Delta: string(r)}
	}
	ch <- provider.Chunk{Done: true}
	close(ch)
	return ch
}

const testResponse = `## 解释
这条命令列出当前目录的所有文件，包括隐藏文件。
-l 参数表示长格式输出，-a 表示显示所有文件。

## 命令
` + "```bash\n" + `ls -la
` + "```"

func TestParse_ExtractsCommandAndExplanation(t *testing.T) {
	chunks := makeChunks(testResponse)
	result := Parse(chunks, nil)

	if result.Err != nil {
		t.Fatal(result.Err)
	}
	if result.Command != "ls -la" {
		t.Errorf("Command = %q", result.Command)
	}
	if !strings.Contains(result.Explanation, "隐藏文件") {
		t.Errorf("Explanation missing expected content: %s", result.Explanation)
	}
}

func TestParse_NoCommandBlock(t *testing.T) {
	chunks := makeChunks("## 解释\n只有解释没有命令块。")
	result := Parse(chunks, nil)
	if result.Command != "" {
		t.Errorf("expected empty command, got %q", result.Command)
	}
}
