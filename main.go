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

	// 展開通用字體Wildcard
	// "ls -l m??go"  ->  "ls -l mango mongo"
	input = expandWildcardInCmd(input)

	if strings.HasPrefix(input, `\`) {
		// 如果指令是 \ 開頭，那就把第一個字元去掉
		// \ls -> ls
		input = input[1:]
	} else { // 如果不是 \ 開頭，那就還是跑 expandAlias
		// 把使用者的輸入根據 alias 進行 expandAlias
		// Ex:  "gst"   ->  "git status"
		// Ex: "ls -h"  ->  "ls -l -h"
		input = expandAlias(input)
	}

	// 把使用者的輸入切割成 Array
	// "ps aux" -> ["ps", "aux"]
	//args := strings.Split(input, " ")
	args := parseArgs(input)

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

	// 如果使用者輸入 alias 開頭的指令
	// 譬如說 alias gst='git status'
	// args = ["alias", "gst='git status'"]
	if args[0] == "alias" {
		// 把 args[1] 切成左右兩邊，變成 key 跟 value
		// kv = ["gst", "'git status'"]
		kv := strings.Split(args[1], "=")

		// 把 value 的單引號去掉取得真正的 value
		// key = "gst", val = "git status"
		key, val := kv[0], strings.Trim(kv[1], "'")
		setAlias(key, val)

		// 沒有錯誤發生
		return nil
	}
	// 取消快捷設定
	// args = ["unalias", "gst"]
	if args[0] == "unalias" {
		// key = "gst"
		key := args[1]

		// unsetAlias("gst")
		unsetAlias(key)
		return nil
	}
	// 如果使用者下的指令是 which 開頭
	// 譬如說 which gst ls gggg
	if args[0] == "which" {
		// 那就把 gst, ls, gggg 依序傳到 lookCommand 裡面
		// 讓他跑上面的 which 流程
		for _, cmd := range args[1:] {
			lookCommand(cmd)
		}
		return nil
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

// 開頭是alias要另外切
// paresArgs 就是上圖的 parseArgs
// 這邊要根據上面的演算法來實作他
func parseArgs(input string) []string {
	// 如果 input 是 "alias" 開頭，那最多就切成兩段
	// 也就是 ["alias", "ooo xxx ooo xxx"]
	// 後面的空白不會被切到
	if strings.HasPrefix(input, "alias") {
		return strings.SplitN(input, " ", 2)
	}

	// 如果不是 "alias" 開頭
	// 那就用原本的方法，把所有空白都切開
	// "ls -l -a" -> ["ls", "-l", "-a"]
	return strings.Split(input, " ")
}

func lookCommand(cmd string) {
	// 先到 Alias Table 裡面找找看指令(gst)
	value := aliasTable[cmd]

	// 如果找到的話就輸出 gst: aliased to git status
	// 沒找到的話就繼續往下走
	if value != "" {
		fmt.Printf("%s: aliased to %s\n", cmd, value)
		return
	}

	// 到 PATH 裡面找找看指令(ls)
	value, err := exec.LookPath(cmd)

	// 找到的話就輸出 ls: /bin/ls
	// 沒找到的話就繼續往下走
	if err == nil {
		fmt.Printf("%s: %s\n", cmd, value)
		return
	}

	// Alias Table 跟 PATH 都找不到這個指令(gggg)
	// 直接輸出 gggg NOT FOUND
	fmt.Printf("%s NOT FOUND\n", cmd)
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
