package controllers

import (
	"fmt"
	"go-blog/app/models/user"
	"go-blog/app/requests"
	"go-blog/pkg/auth"
	"go-blog/pkg/flash"
	"go-blog/pkg/view"
	"net/http"
)

// AuthController 处理用户认证
type AuthController struct {
}

// Register 注册页面
func (*AuthController) Register(w http.ResponseWriter, r *http.Request) {
	view.RenderSimple(w, view.D{
		"User": user.User{},
	}, "auth.register")
}

// DoRegister 处理注册逻辑
func (*AuthController) DoRegister(w http.ResponseWriter, r *http.Request) {

	// 1. 初始化数据
	_user := user.User{
		Name:            r.PostFormValue("name"),
		Email:           r.PostFormValue("email"),
		Password:        r.PostFormValue("password"),
		PasswordConfirm: r.PostFormValue("password_confirm"),
	}

	// 2. 表单验证
	errs := requests.ValidateRegistrationForm(_user)

	if len(errs) > 0 {
		// 3. 有错误发生，打印数据
		//data, _ := json.MarshalIndent(errs, "", " ")
		//fmt.Fprint(w, string(data))
		view.RenderSimple(w, view.D{
			"Errors": errs,
			"User":   _user,
		}, "auth.register")
	} else {
		// 4. 验证成功，创建数据
		_user.Create()

		if _user.ID > 0 {
			flash.Success("恭喜您注册成功!")
			auth.Login(_user)
			http.Redirect(w, r, "/", http.StatusFound)
		} else {
			fmt.Fprint(w, "创建用户失败，请联系管理员")
		}
	}

}

// Login 显示登录表单
func (*AuthController) Login(w http.ResponseWriter, r *http.Request) {
	view.RenderSimple(w, view.D{
		"Email":    "",
		"Password": "",
	}, "auth.login")
}

// DoLogin 处理登录表单提交
func (*AuthController) DoLogin(w http.ResponseWriter, r *http.Request) {
	// 1. 初始化表单数据
	email := r.PostFormValue("email")
	password := r.PostFormValue("password")

	// 2. 尝试登录
	if err := auth.Attempt(email, password); err == nil {
		// 登录成功
		flash.Success("欢迎回来!")
		http.Redirect(w, r, "/", http.StatusFound)
	} else {
		// 3. 登录失败，显示错误提示
		view.RenderSimple(w, view.D{
			"Error":    err.Error(),
			"Email":    email,
			"Passowrd": password,
		}, "auth.login")
	}
}

// Logout 用户退出
func (*AuthController) Logout(w http.ResponseWriter, r *http.Request) {
	auth.Logout()
	flash.Success("您已退出登录!")
	http.Redirect(w, r, "/", http.StatusFound)
}
