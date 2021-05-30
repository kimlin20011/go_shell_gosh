package main

import (
	"path/filepath"
	"strings"
)

func expandPattern(pattern string) string {
	// 把 pattern 丟進去 Glob 讓他幫我找到符合的檔案名稱
	// filenames = ["mango", "mongo"]
	filenames, _ := filepath.Glob(pattern)

	// 再把檔名連起來
	// ["mango", "mongo"] -> "mango mongo"
	return strings.Join(filenames, " ")
}

// Usage:
// expandWildcardInCmd("ls -l m??go") == "ls -l mango mongo"
//
func expandWildcardInCmd(input string) string {
	// 把字串切成好幾個參數
	// "ls -l m??go" -> ["ls", "-l", "m??go"]
	args := strings.Split(input, " ")

	// 把每個參數跑過一遍
	for i, arg := range args {
		// 如果參數有包含 * 或 ?，就用 expandPattern 進行展開
		// 如果沒有包含，就什麼都不做
		if strings.Contains(arg, "*") || strings.Contains(arg, "?") {
			args[i] = expandPattern(arg)
		}

		// "ls" -> "ls"
		// "-l" -> "-l"
		// "m??go" -> "mango mongo"
	}

	// 把 []string 組合起來
	// ["ls", "-l", "mango mongo"] -> "ls -l mango mongo"
	return strings.Join(args, " ")
}
