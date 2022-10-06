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
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	serverAddress := "localhost:8066"
	log.Printf("server starting with address: %s", serverAddress)
	err := http.ListenAndServe(serverAddress, nil)
	if err != nil {
		log.Fatalf("start server failed, err: %s", err.Error())
	}
}
