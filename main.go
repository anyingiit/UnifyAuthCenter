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
	"github.com/anyingiit/UnifyAuthCenter/myErrors"
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

func generateTOTPWelcomePage(w http.ResponseWriter, r *http.Request) error {
	t, err := template.ParseFiles("./template/tool/generation_TOTP/welcome.tmpl")
	if err != nil {
		return err
	}
	err = t.Execute(w, nil)
	if err != nil {
		return err
	}
	return nil
}

func generateTOTPFormPage(w http.ResponseWriter, r *http.Request) error {
	t, err := template.ParseFiles("./template/tool/generation_TOTP/form.tmpl")
	if err != nil {
		return err
	}
	err = t.Execute(w, nil)
	if err != nil {
		return err
	}
	return nil
}

func generateTOTPGenerationPage(w http.ResponseWriter, r *http.Request) error {
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
		return err
	}
	// fmt.Println(t)

	secret, pngBase64String, err := utils.GenerationNewTOTP(query.Issure, query.AccountName, 200, 200)

	if err != nil {
		return fmt.Errorf("generation TOTP failed, err: %s", err.Error())
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
		return err
	}

	w.WriteHeader(http.StatusOK)
	return nil
}

func authCenterLogin(w http.ResponseWriter, r *http.Request) error {
	t, err := template.ParseFiles("./template/login/form.tmpl")
	if err != nil {
		return err
	}

	// 如果没有显式调用WriteHeader, 那么在第一次Write的时候, 将自动调用WriteHeader(http.StatusOK)
	w.WriteHeader(http.StatusOK)
	err = t.Execute(w, nil)
	if err != nil {
		return err
	}

	return nil
}
func authCenterHandle(w http.ResponseWriter, r *http.Request) error {
	type Query struct {
		Code string
	}
	query := Query{
		Code: r.FormValue("code"),
	}
	if query.Code == "" {
		return myErrors.NewSimpleAuthorizationError("auth center login failed", "code is empty")
	}
	validated := utils.ValidateTOTP(query.Code, "REMOVED-SEE-README")
	if !validated {
		return myErrors.NewSimpleAuthorizationError("auth center login failed", "code is invalid")
	}

	nowTime := time.Now()
	session := &models.Session{
		UUID:      uuid.New(),
		CreatedAt: nowTime,                    // if CreatedAt value is empty time.Time obj, gorm will automatically set this to the current time
		ExpiredAt: nowTime.Add(time.Hour * 8), // Expires in 8 hours
	}
	result := session.Create()

	if result.Error != nil {
		return fmt.Errorf("create session failed, err: %s", result.Error.Error())
	}

	t, err := template.ParseFiles("./template/login/success.tmpl")
	if err != nil {
		return err
	}
	// 缓存不应存储有关客户端请求或服务器响应的任何内容，即不使用任何缓存。
	w.Header().Add("Cache-control", "no-store")
	w.Header().Add("Set-Cookie", fmt.Sprintf("uuid=%s", session.UUID.String()))
	err = t.Execute(w, nil)
	if err != nil {
		return err
	}

	return nil
}

func authCenterStatus(w http.ResponseWriter, r *http.Request) error {
	type Cookie struct {
		uuid string
	}
	cookieStr := r.Header.Get("Cookie")
	if cookieStr == "" {
		return myErrors.NewSimpleAuthorizationError("auth center unauthorized", "cookie is empty")
	}
	cookie := Cookie{
		uuid: strings.Split(cookieStr, "uuid=")[1],
	}
	if cookie.uuid == "" {
		return myErrors.NewSimpleAuthorizationError("auth center unauthorized", "uuid is empty")
	}
	UUID, err := uuid.Parse(cookie.uuid)
	if err != nil {
		return myErrors.NewSimpleAuthorizationError("auth center unauthorized", "uuid type error")
	}
	session := &models.Session{
		UUID: UUID,
	}
	result := session.First()
	if result.Error != nil {
		return myErrors.NewSimpleAuthorizationError("auth center unauthorized", "uuid invalid or expired")
	}

	if session.ExpiredAt.UnixNano() < time.Now().UnixNano() {
		return myErrors.NewSimpleAuthorizationError("auth center unauthorized", "uuid expired")
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("auth center authorized"))

	return nil
}

func adminLoginPage(w http.ResponseWriter, r *http.Request) error {
	t, err := template.ParseFiles("./template/admin/form.tmpl")
	if err != nil {
		return err
	}
	err = t.Execute(w, nil)
	if err != nil {
		return err
	}

	return nil
}

func adminLoginHandlePage(w http.ResponseWriter, r *http.Request) error {
	type Query struct {
		code string
	}
	query := Query{
		code: r.FormValue("code"),
	}
	if query.code == "" {
		//TODO
	}

	//TODO

	return nil
}

type appHandler func(http.ResponseWriter, *http.Request) error

func errWrapper(handle appHandler) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("panic: %v", r)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()

		err := handle(w, r)

		if err == nil {
			return
		}

		responseError := myErrors.NewMyError(err, http.StatusInternalServerError, "internal error", "internal error")

		if userError, ok := err.(myErrors.UserError); ok {
			responseError.StatusCode = userError.GetStatusCode()
			responseError.Event = userError.GetEvent()
			responseError.Reason = userError.GetReason()
		}

		InternalError := func(w http.ResponseWriter) {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Internal Server Error"))
		}

		if responseError.StatusCode == http.StatusInternalServerError {
			log.Printf("internal error: %s", responseError.Err.Error())
			InternalError(w)
			return
		}

		// user error
		t, err := template.ParseFiles("./template/error/user_error.tmpl")
		if err != nil {
			log.Printf("internal error: %s", responseError.Err.Error())
			InternalError(w)
			return
		}

		log.Printf("user error: %s", responseError.Err.Error())
		w.WriteHeader(responseError.StatusCode)
		err = t.Execute(w, struct {
			Message string
			Event   string
			Reason  string
		}{
			Message: http.StatusText(responseError.StatusCode),
			Event:   responseError.Event,
			Reason:  responseError.Reason,
		})
		if err != nil {
			log.Printf("internal error: %s", responseError.Err.Error())
			InternalError(w)
			return
		}
	}
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

	// 工具相关
	// 生成TOTP
	http.HandleFunc("/tool/generation_TOTP", func(w http.ResponseWriter, r *http.Request) { // 重定向到欢迎页面
		w.Header().Add("Location", "/tool/generation_TOTP/welcome")
		w.WriteHeader(http.StatusFound)
	})
	http.HandleFunc("/tool/generation_TOTP/welcome", errWrapper(generateTOTPWelcomePage))       // 欢迎页面
	http.HandleFunc("/tool/generation_TOTP/form", errWrapper(generateTOTPFormPage))             // 表单页面
	http.HandleFunc("/tool/generation_TOTP/generation", errWrapper(generateTOTPGenerationPage)) // 生成TOTP页面

	// 验证中心
	http.HandleFunc("/auth_center", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Location", "/auth_center/login")
		w.WriteHeader(http.StatusFound)
	})
	http.HandleFunc("/auth_center/login", errWrapper(authCenterLogin))
	http.HandleFunc("/auth_center/handle", errWrapper(authCenterHandle))
	// http.HandleFunc("/login/success", loginPage)
	// http.HandleFunc("/login/failed", loginPage)
	http.HandleFunc("/auth_center/status", errWrapper(authCenterStatus))

	// 管理员
	http.HandleFunc("/admin/login", errWrapper(adminLoginPage))
	http.HandleFunc("/admin/login/handle", errWrapper(adminLoginHandlePage))
	// http.HandleFunc("/admin/session_manage", errWrapper(adminLoginFormPage))
	serverAddress := "localhost:8066"
	log.Printf("server starting with address: %s", serverAddress)
	err = http.ListenAndServe(serverAddress, nil)
	if err != nil {
		log.Fatalf("start server failed, err: %s", err.Error())
		return
	}
}
