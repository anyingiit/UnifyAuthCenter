package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/anyingiit/UnifyAuthCenter/db"
	"github.com/anyingiit/UnifyAuthCenter/models"
	"github.com/anyingiit/UnifyAuthCenter/utils"
	"github.com/google/uuid"
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
func generationTOTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Location", "/tool/generation_TOTP/welcome")
	w.WriteHeader(http.StatusFound)
}
func generateTOTPWelcomePage(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("./template/tool/generation_TOTP/welcome.tmpl")
	if err != nil {
		log.Printf("parse template failed, err: %s\n", err.Error())
		http.Error(w, "parse template failed", http.StatusInternalServerError)
		return
	}
	err = t.Execute(w, nil)
	if err != nil {
		log.Printf("execute template failed, err: %s\n", err.Error())
		http.Error(w, "execute template failed", http.StatusInternalServerError)
		return
	}
}

func generateTOTPFormPage(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("./template/tool/generation_TOTP/form.tmpl")
	if err != nil {
		log.Printf("parse template failed, err: %s\n", err.Error())
		http.Error(w, "parse template failed", http.StatusInternalServerError)
		return
	}
	err = t.Execute(w, nil)
	if err != nil {
		log.Printf("execute template failed, err: %s\n", err.Error())
		http.Error(w, "execute template failed", http.StatusInternalServerError)
		return
	}
}

func generateTOTPGenerationPage(w http.ResponseWriter, r *http.Request) {
	useDefaultIssureAndDefaultAccountName := false
	type Query struct {
		Issure      string
		AccountName string
	}
	query := Query{
		Issure:      r.FormValue("issure"),
		AccountName: r.FormValue("account_name"),
	}
	if query.Issure == "" || query.AccountName == "" {
		useDefaultIssureAndDefaultAccountName = true
		query.Issure = "unifyauthcenter.com"
		query.AccountName = "Jack A. Doe"
	}
	// fmt.Println(query)

	t, err := template.ParseFiles("./template/tool/generation_TOTP/generation.tmpl")
	if err != nil {
		log.Printf("parse template failed, err: %s\n", err.Error())
		http.Error(w, "parse template failed", http.StatusInternalServerError)
		return
	}
	// fmt.Println(t)

	secret, pngBase64String, err := utils.GenerationNewTOTP(query.Issure, query.AccountName, 200, 200)

	if err != nil {
		log.Printf("generation TOTP failed, err: %s\n", err.Error())
		http.Error(w, "generation TOTP failed", http.StatusInternalServerError)
		return
	}

	// fmt.Println(secret, pngBase64String)

	err = t.Execute(w, struct {
		UseDefaultIssureAndDefaultAccountName string
		Issure                                string
		AccountName                           string
		Secret                                string
		/*
			当这里类型为string时, 实际渲染出的网页会显示为`#ZgotmplZ`
			这是因为...
				"ZgotmplZ" is a special value that indicates that unsafe content reached a
				CSS or URL context at runtime. The output of the example will be
					<img src="#ZgotmplZ">
				If the data comes from a trusted source, use content types to exempt it
				from filtering: URL(`javascript:...`).
		*/
		SecretPngBase64 template.URL
	}{
		UseDefaultIssureAndDefaultAccountName: fmt.Sprintf("%t", useDefaultIssureAndDefaultAccountName),
		Issure:                                query.Issure,
		AccountName:                           query.AccountName,
		Secret:                                secret,
		SecretPngBase64:                       template.URL(`data:image/png;base64,` + pngBase64String),
	})

	if err != nil {
		log.Printf("execute template failed, err: %s\n", err.Error())
		http.Error(w, "execute template failed", http.StatusInternalServerError)
		return
	}
}

func loginPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Location", "/login/form")
	w.WriteHeader(http.StatusFound)
}

func loginFormPage(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("./template/login/form.tmpl")
	if err != nil {
		log.Printf("parse template failed, err: %s\n", err.Error())
		http.Error(w, "parse template failed", http.StatusInternalServerError)
		return
	}
	err = t.Execute(w, nil)
	if err != nil {
		log.Printf("execute template failed, err: %s\n", err.Error())
		http.Error(w, "execute template failed", http.StatusInternalServerError)
		return
	}
}
func loginHandle(w http.ResponseWriter, r *http.Request) {
	success := func(w http.ResponseWriter, setCookieString string) {
		t, err := template.ParseFiles("./template/login/success.tmpl")
		if err != nil {
			log.Printf("parse template failed, err: %s\n", err.Error())
			http.Error(w, "parse template failed", http.StatusInternalServerError)
			return
		}
		// 缓存不应存储有关客户端请求或服务器响应的任何内容，即不使用任何缓存。
		w.Header().Add("Cache-control", "no-store")
		w.Header().Add("Set-Cookie", setCookieString)
		err = t.Execute(w, nil)
		if err != nil {
			log.Printf("execute template failed, err: %s\n", err.Error())
			http.Error(w, "execute template failed", http.StatusInternalServerError)
			return
		}
	}
	failed := func(w http.ResponseWriter, reason string) {
		t, err := template.ParseFiles("./template/login/failed.tmpl")
		if err != nil {
			log.Printf("parse template failed, err: %s\n", err.Error())
			http.Error(w, "parse template failed", http.StatusInternalServerError)
			return
		}
		err = t.Execute(w, reason)
		if err != nil {
			log.Printf("execute template failed, err: %s\n", err.Error())
			http.Error(w, "execute template failed", http.StatusInternalServerError)
			return
		}
	}

	type Query struct {
		Code string
	}
	query := Query{
		Code: r.FormValue("code"),
	}
	if query.Code == "" {
		failed(w, "code is empty")
		return
	}
	validated := utils.ValidateTOTP(query.Code, "REMOVED-SEE-README")
	if !validated {
		failed(w, "code is invalid")
		return
	}

	nowTime := time.Now()
	session := &models.Session{
		UUID:      uuid.New(),
		CreatedAt: nowTime,                    // if CreatedAt value is empty time.Time obj, gorm will automatically set this to the current time
		ExpiredAt: nowTime.Add(time.Hour * 8), // Expires in 8 hours
	}
	result := session.Create()

	if result.Error != nil {
		log.Printf("create session failed, err: %s\n", result.Error.Error())
		failed(w, "cannot insert session into database")
		return
	}

	success(w, fmt.Sprintf("uuid=%s", session.UUID.String()))
}

func loginStatus(w http.ResponseWriter, r *http.Request) {
	type Cookie struct {
		uuid string
	}
	cookieStr := r.Header.Get("Cookie")
	if cookieStr == "" {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("cookie is empty"))
		return
	}
	cookie := Cookie{
		uuid: strings.Split(cookieStr, "uuid=")[1],
	}
	if cookie.uuid == "" {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("uuid is empty"))
		return
	}
	UUID, err := uuid.Parse(cookie.uuid)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("uuid type error"))
		return
	}
	session := &models.Session{
		UUID: UUID,
	}
	result := session.First()
	if result.Error != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("uuid invalid or expired"))
		return
	}

	if session.ExpiredAt.UnixNano() < time.Now().UnixNano() {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("uuid expired"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("uuid valid"))
}

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

	err := initDatabase()
	if err != nil {
		log.Fatal("init database failed, err: ", err)
		return
	}
	// 静态文件
	// `http.Handle("/static/", http.StripPrefix("/static/"`中的`/static/`必须是`/static/`, 而不能是`/static`
	//		从URI的语义来说, `/static/`目录, 而`/static`是某个资源
	// 而`http.FileServer(http.Dir("./static")`中的`./static`必须是`./static`或者`static`
	//		因为`static`里代表的是`./static`的简写, 而`./static`是相对路径, 代表的是以当前代码文件为中心所指的文件
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	// 登录相关

	// 工具相关
	// 生成TOTP
	http.HandleFunc("/tool/generation_TOTP", generationTOTP)                        // 重定向到欢迎页面
	http.HandleFunc("/tool/generation_TOTP/welcome", generateTOTPWelcomePage)       // 欢迎页面
	http.HandleFunc("/tool/generation_TOTP/form", generateTOTPFormPage)             // 表单页面
	http.HandleFunc("/tool/generation_TOTP/generation", generateTOTPGenerationPage) // 生成TOTP页面
	http.HandleFunc("/login", loginPage)
	http.HandleFunc("/login/form", loginFormPage)
	http.HandleFunc("/login/handle", loginHandle)
	// http.HandleFunc("/login/success", loginPage)
	// http.HandleFunc("/login/failed", loginPage)
	http.HandleFunc("/login/status", loginStatus)
	serverAddress := "localhost:8066"
	log.Printf("server starting with address: %s", serverAddress)
	err = http.ListenAndServe(serverAddress, nil)
	if err != nil {
		log.Fatalf("start server failed, err: %s", err.Error())
		return
	}
}
