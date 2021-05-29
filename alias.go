package main

import "strings"

// 初始化一個空的 alias table
// 型別是 string 對應到 string
var aliasTable = map[string]string{}

func setAlias(key, value string) {
	// 把 key 跟對應的 value 丟到 alias table 裡面去
	// Ex: aliasTable["ls"] = "ls -l"
	aliasTable[key] = value
}

func unsetAlias(key string) {
	// 把 table 中的某個 key 刪掉
	// Ex: delete(aliasTable, "ls")
	delete(aliasTable, key)
}

func expandAlias(input string) string {
	// 把 input 的指令跟參數切開
	// input = "ls -h"
	// args = ["ls", "-h"]
	args := strings.SplitN(input, " ", 2)

	// 拿出指令部分，cmd = "ls"
	cmd := args[0]

	// 如果這個指令有在 alias table 裡面
	// 就會得到一個 expandedCmd（展開過後的指令）
	// 再用新的指令取代舊的指令
	if expandedCmd, ok := aliasTable[cmd]; ok {
		// cmd = "ls"  ->  expandedCmd = "ls -l"
		// 把第一個出現的 "ls" 替換成 "ls -l" 就大功告成了
		return strings.Replace(input, cmd, expandedCmd, 1)
	}

	// 如果指令不在 alias table 裡面
	// 就什麼都不做，回傳原本的 input
	return input
}
