package parser

import (
	"strings"

	"github.com/aict-tool/aict/internal/provider"
)

type Result struct {
	Explanation string
	Command     string
	Err         error
}

type RenderFunc func(delta string)

type state int

const (
	stateText state = iota
	stateCode
)

func Parse(chunks <-chan provider.Chunk, renderFn RenderFunc) Result {
	var fullText strings.Builder

	for chunk := range chunks {
		if chunk.Err != nil {
			return Result{Err: chunk.Err}
		}
		if chunk.Done {
			break
		}

		delta := chunk.Delta
		if renderFn != nil {
			renderFn(delta)
		}
		fullText.WriteString(delta)
	}

	lines := strings.Split(fullText.String(), "\n")
	curState := stateText
	inCommandSection := false
	var explanationLines []string
	var commandBuf strings.Builder

	for i := 0; i < len(lines); i++ {
		line := lines[i]
		trimmed := strings.TrimSpace(line)

		switch curState {
		case stateText:
			if strings.HasPrefix(trimmed, "## 命令") {
				inCommandSection = true
				continue
			}
			if strings.HasPrefix(trimmed, "## ") {
				inCommandSection = false
			}
			if strings.HasPrefix(trimmed, "```") && inCommandSection {
				curState = stateCode
				continue
			}
			if !inCommandSection && trimmed != "" && !strings.HasPrefix(trimmed, "## ") {
				explanationLines = append(explanationLines, line)
			}
		case stateCode:
			if strings.HasPrefix(trimmed, "```") {
				curState = stateText
				continue
			}
			commandBuf.WriteString(line)
			commandBuf.WriteString("\n")
		}
	}

	return Result{
		Explanation: strings.TrimSpace(strings.Join(explanationLines, "\n")),
		Command:     strings.TrimSpace(commandBuf.String()),
	}
}
