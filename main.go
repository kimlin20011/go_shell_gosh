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

//在shell前面加上user資訊
func showPrompt() {
	u, _ := user.Current()   // 取得使用者資訊
	host, _ := os.Hostname() // 取得主機名稱
	wd, _ := os.Getwd()      // 取得 Working Directory

	// 把 user 跟 host 著成藍色
	userAndHost := blue(u.Username + "@" + host)

	// 把 Working Directory 變成藍底黃字
	wd = yellowWithBlueBG(wd)

	// 把字串組合起來放到 Prompt 中
	fmt.Printf("%s %s > ", userAndHost, wd)
}

func executeInput(input string) error {
	// 環境變數展開
	// 把 input 裡面的環境變數展開
	// 如果指令是 echo ${USER}_is_smart
	// 展開後就會變成 echo larry_is_smart
	input = os.ExpandEnv(input)
	// 把使用者的輸入切割成 Array
	// "ps aux" -> ["ps", "aux"]
	args := strings.Split(input, " ")

	//自己實作cd指令
	if args[0] == "cd" {
		// 如果指令是 cd dirname
		// 就跑 os.Chdir(dirname)
		err := os.Chdir(args[1])

		// 並回傳發生的錯誤（例如資料夾不存在）
		return err
	}
	//自己實作exit指令
	if args[0] == "exit" {
		// 如果指令是 cd dirname
		// 就跑 os.Chdir(dirname)
		os.Exit(1)
	}
	// 如果指令是 export 開頭
	// 就用 os.Setenv 設置環境變數
	// args = ["export", "FOO=bar"]
	if args[0] == "export" {
		// kv = ["FOO", "bar"]
		kv := strings.Split(args[1], "=")

		// key = "FOO"
		// val = "bar"
		key, val := kv[0], kv[1]

		err := os.Setenv(key, val)
		return err
	}

	// 如果指令是 unset 開頭
	// 就用 os.Unsetenv 刪除環境變數
	// args = ["unset", "FOO"]
	if args[0] == "unset" {
		err := os.Unsetenv(args[1])
		return err
	}

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
