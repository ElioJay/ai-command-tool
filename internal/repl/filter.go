package repl

import (
	"strings"
	"unicode/utf8"
)

var chatPatterns = []string{
	"你好", "您好", "hello", "hi ",
	"你是谁", "你是什么", "你叫什么",
	"讲个笑话", "讲个故事", "写一首",
	"翻译", "帮我翻译",
	"写一篇", "写一段", "帮我写",
	"聊聊", "闲聊", "陪我聊",
	"什么意思", "是什么意思",
	"解释一下什么是",
	"你觉得", "你认为", "你怎么看",
	"推荐一", "建议一",
	"1+1", "算一下",
	"今天天气", "明天天气",
}

func IsOffTopic(input string) bool {
	lower := strings.ToLower(input)

	for _, p := range chatPatterns {
		if strings.Contains(lower, p) {
			return true
		}
	}

	if utf8.RuneCountInString(input) <= 4 && !looksLikeCommand(lower) {
		return true
	}

	return false
}

func looksLikeCommand(input string) bool {
	cmdHints := []string{
		"文件", "目录", "进程", "端口", "磁盘", "内存", "网络",
		"查找", "搜索", "删除", "复制", "移动", "重命名", "压缩", "解压",
		"安装", "卸载", "更新", "升级",
		"git", "docker", "npm", "pip", "go ", "cargo",
		"启动", "停止", "重启", "运行", "执行",
		"权限", "用户", "服务", "日志", "监控",
		"下载", "上传", "传输",
		"编译", "构建", "打包", "部署",
		"列出", "显示", "查看", "统计", "计算大小",
		"grep", "find", "ls", "cd", "mkdir", "rm ",
		"ping", "curl", "wget", "ssh", "scp",
		"cmd", "powershell", "bash", "shell", "脚本",
		"命令", "终端", "控制台",
		"环境变量", "路径", "配置",
	}
	for _, h := range cmdHints {
		if strings.Contains(input, h) {
			return true
		}
	}
	return false
}
