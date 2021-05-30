# go_shell_gosh


## 筆記
* 列出所有檔案
```
ls -la
```

* 還沒有建立專案的時候compile要用
```
go run main.go color.go
```
把所有有用到的file都要compile

* 初始化 go module
```
go mod init <project-name>
```

* alias語法 （將指令設定簡寫）
`alias name='command -args'`
* 例如
`alias gst='git status'`
* 取消簡寫設定
`unalias name`

* 測資的file 要import `testing` package
    * 下`go test`指令可以直接執行測試

* 執行go，可以執行所有需要執行的function
`go run .`

* [永久設定alias的方法](https://qiita.com/yutat93/items/b5bb9c0366f21bcbea62)

### gitignore建立網站
`https://www.toptal.com/developers/gitignore`

