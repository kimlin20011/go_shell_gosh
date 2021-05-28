package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	// 把使用者輸入轉換成 bufio.Reader 型別
	stdin := bufio.NewReader(os.Stdin)

	for {
		// 簡單的 prompt
		fmt.Print("> ")

		// 逐行讀取使用者輸入，並且去除頭尾的空白
		input, _ := stdin.ReadString('\n')
		input = strings.TrimSpace(input)

		// 輸出
		fmt.Println(input)
	}
}
