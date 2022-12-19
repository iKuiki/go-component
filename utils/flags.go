package utils

import "strings"

// SplitCommandFlag 将带有子命令的args拆分为主命令的flag和子命令的args
// 比如-e pro migrate -exec grant
// 主命令的flag(mainFlags)为: -e pro
// 子命令为: migrate
// 子命令的args为: -exec grant
func SplitCommandFlag(args []string) (mainFlags []string, subCommand string, subArgs []string) {
	var previewIsFlag bool // 上一个指令是否是args
	for i := 0; i < len(args); i++ {
		currentArg := args[i]
		if strings.HasPrefix(currentArg, "-") {
			previewIsFlag = true
			continue
		}
		// 当前不是flag，则检查上一个是否是flag
		if previewIsFlag {
			previewIsFlag = false // 进入下一轮之前县把flag设置为false
			continue
		}
		// 上一个也不是flag，则这是command
		subCommand = currentArg
		if i > 0 { // 判断有没有主flag
			// 有主flag
			mainFlags = args[:i]
		}
		if i < len(args)-1 { // 判断有没有子flag
			subArgs = args[i+1:]
		}
		return
	}
	// 如果一直没找到，就全都是主flag
	mainFlags = args
	return
}
