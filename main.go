package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/user"
	"strings"
)

func main() {
	// 把使用者輸入轉換成 bufio.Reader 型別
	stdin := bufio.NewReader(os.Stdin)

	for {
		// 簡單的 prompt
		//fmt.Print("> ")
		showPrompt()

		// 逐行讀取使用者輸入
		input, _ := stdin.ReadString('\n')
		//並且去除頭尾的空白
		input = strings.TrimSpace(input)

		// 執行使用者輸入的指令
		// 如果有錯誤的話就 log 出來
		err := executeInput(input)
		if err != nil {
			log.Println(err)
		}
	}
}

//在shell前面加上user資訊
func showPrompt() {
	u, _ := user.Current()   // 取得使用者資訊
	host, _ := os.Hostname() // 取得主機名稱
	wd, _ := os.Getwd()      // 取得 Working Directory

	// 把字串組合起來放到 Prompt 中
	fmt.Printf("%s@%s %s > ", u.Username, host, wd)
}

func executeInput(input string) error {
	// 把使用者的輸入切割成 Array
	// "ps aux" -> ["ps", "aux"]
	args := strings.Split(input, " ")

	// 根據使用者的輸入建立一個指令
	// 譬如說使用者輸入 ls，就建立一個 ls 指令
	// args[0] 是指令名，放在第一個位置
	// args[1:]... 是把其他參數依序填入裡面
	// ["ls", "-l", "-a"] 即 exec.Command("ls", "-l", "-a")
	cmd := exec.Command(args[0], args[1:]...)

	// 使用 exec.Command 建立的 cmd 預設是不輸出（超怪XD）
	// 所以要把他的 Stdandard IO 重新設定成系統預設（終端機）
	// 他才能正常輸出到終端機、從終端機讀取資料
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// 啟動一個 Child Process 執行剛建立好的指令
	// 如果使用者輸入 ls 那就是執行 ls 的執行檔
	err := cmd.Run()

	// 如果有發生錯誤的話就回傳
	return err
}
