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
