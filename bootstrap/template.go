package bootstrap

import (
	"embed"
	"go-blog/pkg/view"
)

// SetupTemplate 模版初始化
func SetupTemplate(tmplFS embed.FS) {
	view.TplFS = tmplFS
}
