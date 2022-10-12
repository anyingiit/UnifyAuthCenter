package main

import (
	"log"
	"net/http"

	"github.com/anyingiit/UnifyAuthCenter/db"
	"github.com/anyingiit/UnifyAuthCenter/models"
	"github.com/anyingiit/UnifyAuthCenter/router"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// TODO:
//
//	  登录相关:
//	  	1. 能够通过TOTP登录的界面, 该接口能够效验TOTP, 如果效验成功, 将通过Set-Cookie返回一个session
//	  	2. 提供一个查询登录状态的接口, 该接口可用于内部与外部的调用, 用于判断用户是否已经登录, 状态有效则返回200, 状态无效则返回401
//	  	3. 提供一个用于管理员端的界面, 该界面能够列出所有Session信息, 并且能够使某个Session失效
//	  工具相关:
//		1. 提供一个用于生成TOTP的工具, 该工具能够生成一个TOTP, 并将TOTP的二维码和其他相关信息通过HTML的方式展示用户浏览器(注: TOTP的恢复密码不是TOTP中的标准, 是需要用户自行定义恢复规则并生成的	)

func main() {
	// 初始化数据库
	// 初始化成功后, 全局对象`db`将可用
	initDatabase := func() error {
		// 连接数据库
		dbObj, err := gorm.Open(sqlite.Open("database/sessions.db"), &gorm.Config{})
		if err != nil {
			return err
		}

		// 创建表
		err = dbObj.AutoMigrate(&models.Session{})
		if err != nil {
			return err
		}

		db.Db = dbObj
		return nil
	}

	// 初始化路由, 将所有路由注册到给定的serverMux中
	initRouter := func(serverMux *http.ServeMux) error {
		router.RegisteRouters(serverMux)
		return nil
	}

	err := initDatabase()
	if err != nil {
		log.Fatal("init database failed, err: ", err)
		return
	}

	serverMux := http.NewServeMux()

	err = initRouter(serverMux)
	if err != nil {
		log.Fatal("regist router failed, err: ", err)
		return
	}

	serverAddress := "localhost:8066"
	log.Printf("server starting with address: %s", serverAddress)
	// 参数2如果为nil, 则使用DefaultServeMux
	// 我们可以通过http.NewServeMux()新建一个ServeMux, 然后将其作为参数传入
	// func (sh serverHandler) ServeHTTP(rw ResponseWriter, req *Request) {
	// 	handler := sh.srv.Handler
	// 	if handler == nil { // <----------------- 这里
	// 		handler = DefaultServeMux
	// 	}
	// 注: http.HandleFunc(...)和http.Handle(...)都会注册到DefaultServeMux中
	err = http.ListenAndServe(serverAddress, serverMux)
	if err != nil {
		log.Fatalf("start server failed, err: %s", err.Error())
		return
	}
}
