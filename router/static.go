package router

import "net/http"

// 静态文件
// `http.Handle("/static/", http.StripPrefix("/static/"`中的`/static/`必须是`/static/`, 而不能是`/static`
//
//	从URI的语义来说, `/static/`目录, 而`/static`是某个资源
//
// 而`http.FileServer(http.Dir("./static")`中的`./static`必须是`./static`或者`static`
//
//	因为`static`里代表的是`./static`的简写, 而`./static`是相对路径, 代表的是以当前代码文件为中心所指的文件
func regisStatic(serverMux *http.ServeMux) {
	serverMux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
}
