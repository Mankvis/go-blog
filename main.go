package main

import (
	"embed"
	"go-blog/app/http/middlewares"
	"go-blog/bootstrap"
	"go-blog/config"
	c "go-blog/pkg/config"
	"net/http"
)

//go:embed resources/views/articles/*
//go:embed resources/views/categories/*
//go:embed resources/views/auth/*
//go:embed resources/views/layouts/*
var tplFS embed.FS

//go:embed public/*
var staticFS embed.FS

func init() {
	// 初始化配置信息
	config.Initialize()
}

func main() {
	// 初始化 SQL
	bootstrap.SetupDB()

	// 初始化模版
	bootstrap.SetupTemplate(tplFS)

	// 初始化路由绑定
	router := bootstrap.SetupRoute(staticFS)

	http.ListenAndServe(":"+c.GetString("app.port"), middlewares.RemoveTrailingSlash(router))
}
