package controllers

import (
	"fmt"
	"net/http"
)

type PagesController struct {
}

func (*PagesController) Home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "<h1>Hello，欢迎来到 go-blog！</h1>")
}

func (*PagesController) About(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "此博客是用以记录编程博客，如您有反馈或建议，请联系"+
		"<a href=\"https://book.mankvis.com/\"> MANKVIS ` BLOG</a>")
}

func (*PagesController) NotFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, "<h1>请求页面未找到：(<h1>"+
		"<p>如有疑惑，请联系我们。<p>")
}
