package controllers

import (
	"fmt"
	"go-blog/app/models/article"
	"go-blog/app/models/category"
	"go-blog/app/policies"
	"go-blog/app/requests"
	"go-blog/pkg/auth"
	"go-blog/pkg/route"
	"go-blog/pkg/types"
	"go-blog/pkg/view"
	"net/http"
)

type ArticlesController struct {
	BaseController
}

func (ac *ArticlesController) Show(w http.ResponseWriter, r *http.Request) {

	id := route.GetRouteVariable("id", r)

	_article, err := article.Get(id)

	// 3.如果出现错误
	if err != nil {
		ac.ResponseForSQLError(w, err)
	} else {
		view.Render(w, view.D{
			"Article":          _article,
			"CanModifyArticle": policies.CanModifyArticle(_article),
		}, "articles.show", "articles._article_meta")
	}
}

func (ac *ArticlesController) Index(w http.ResponseWriter, r *http.Request) {

	articles, pagerData, err := article.GetAll(r, 5)
	if err != nil {
		ac.ResponseForSQLError(w, err)
	} else {
		view.Render(w, view.D{
			"Articles":  articles,
			"PagerData": pagerData,
		}, "articles.index", "articles._article_meta")
	}
}

func (*ArticlesController) Create(w http.ResponseWriter, r *http.Request) {
	// 读取分类
	categories, _ := category.All()

	view.Render(w, view.D{
		"Title":      "",
		"Body":       "",
		"Article":    article.Article{},
		"Errors":     make(map[string]string),
		"Categories": categories,
	}, "articles.create", "articles._form_field")
}

func (*ArticlesController) Store(w http.ResponseWriter, r *http.Request) {
	// 1. 初始化数据
	currentUser := auth.User()
	_article := article.Article{
		Title:      r.PostFormValue("title"),
		Body:       r.PostFormValue("body"),
		UserID:     currentUser.ID,
		CategoryID: types.StringToUint64(r.PostFormValue("category_id")),
	}

	// 2. 表单验证
	errors := requests.ValidateArticleForm(_article)

	// 3. 检测错误
	if len(errors) == 0 {
		// 创建文章
		_article.Create()
		if _article.ID > 0 {
			indexURL := route.Name2URL("articles.show", "id", _article.GetStringID())
			http.Redirect(w, r, indexURL, http.StatusFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "创建文章失败，请联系管理员")
		}
	} else {
		view.Render(w, view.D{
			"Article": _article,
			"Errors":  errors,
		}, "articles.create", "articles._form_field")
	}
}

func (ac *ArticlesController) Edit(w http.ResponseWriter, r *http.Request) {
	id := route.GetRouteVariable("id", r)

	_article, err := article.Get(id)
	categories, _ := category.All()

	if err != nil {
		ac.ResponseForSQLError(w, err)
	} else {
		// 检查权限
		if !policies.CanModifyArticle(_article) {
			ac.ResponseForUnauthorized(w, r)
		}
		view.Render(w, view.D{
			"Article":    _article,
			"Errors":     view.D{},
			"Categories": categories,
		}, "articles.edit", "articles._form_field")
	}
}

func (ac *ArticlesController) Update(w http.ResponseWriter, r *http.Request) {
	id := route.GetRouteVariable("id", r)

	_article, err := article.Get(id)

	if err != nil {
		ac.ResponseForSQLError(w, err)
	} else {
		// 检查权限
		if !policies.CanModifyArticle(_article) {
			ac.ResponseForUnauthorized(w, r)
		} else {

			cid := r.PostFormValue("category_id")
			_category, err := category.Get(cid)
			if err != nil {
				ac.ResponseForSQLError(w, err)
			}

			_article.CategoryID = types.StringToUint64(cid)
			_article.Category = _category
			_article.Title = r.PostFormValue("title")
			_article.Body = r.PostFormValue("body")

			errors := requests.ValidateArticleForm(_article)

			if len(errors) == 0 {

				rowAffected, err := _article.Update()

				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					fmt.Fprint(w, "500 服务器内部错误")
					return
				}

				if rowAffected > 0 {
					showURL := route.Name2URL("articles.show", "id", id)
					http.Redirect(w, r, showURL, http.StatusFound)
				} else {
					fmt.Fprint(w, "您没有做任何更改！")
				}
			} else {
				view.Render(w, view.D{
					"Article": _article,
					"Errors":  errors,
				}, "articles.edit", "articles._form_field")
			}
		}
	}
}

func (ac *ArticlesController) Delete(w http.ResponseWriter, r *http.Request) {
	id := route.GetRouteVariable("id", r)

	_article, err := article.Get(id)

	if err != nil {
		ac.ResponseForSQLError(w, err)
	} else {
		// 检查权限
		if !policies.CanModifyArticle(_article) {
			ac.ResponseForUnauthorized(w, r)
		} else {
			rowsAffected, err := _article.Delete()
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprint(w, "500 服务器内部错误")
			} else {
				if rowsAffected > 0 {
					indexURL := route.Name2URL("articles.index")
					http.Redirect(w, r, indexURL, http.StatusFound)
				} else {
					w.WriteHeader(http.StatusNotFound)
					fmt.Fprint(w, "404 文章未找到")
				}
			}
		}
	}
}
