package main

import (
	"log"
	"net/http"
)

// TODO:
//   登录相关:
//   1. 能够通过TOTP登录的界面, 该接口能够效验TOTP, 如果效验成功, 将通过Set-Cookie返回一个session
//	 2. 提供一个用于内部效验的接口, 该接口能够通过Set-Cookie中的SessionID效验用户是否处于有效会话中
//	 3. 提供一个用于管理员端的界面, 该界面能够列出所有Session信息, 并且能够使某个Session失效
//	 工具相关:
//	 1. 提供一个用于生成TOTP的工具, 该工具能够生成一个TOTP, 并将TOTP的二维码和其他相关信息通过HTML的方式展示用户浏览器

func main() {
	// `http.Handle("/static/", http.StripPrefix("/static/"`中的`/static/`必须是`/static/`, 而不能是`/static`
	//		从URI的语义来说, `/static/`目录, 而`/static`是某个资源
	// 而`http.FileServer(http.Dir("./static")`中的`./static`必须是`./static`或者`static`
	//		因为`static`里代表的是`./static`的简写, 而`./static`是相对路径, 代表的是以当前代码文件为中心所指的文件
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	serverAddress := "localhost:8066"
	log.Printf("server starting with address: %s", serverAddress)
	err := http.ListenAndServe(serverAddress, nil)
	if err != nil {
		log.Fatalf("start server failed, err: %s", err.Error())
	}
}
