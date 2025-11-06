package xtemplate

import (
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
)

var _default *Manager

func Default() *Manager {
	return _default
}

func Init(bases []string, htmls []string, funcs template.FuncMap) {
	_default = New(bases, htmls, funcs)
}

func InitGlob(basePattern string, htmlPattern string, funcs template.FuncMap) {
	_default = NewGlob(basePattern, htmlPattern, funcs)
}

func SetDebug(debug bool) *Manager {
	return _default.SetDebug(debug)
}

func SetErrorTemplateName(name string) *Manager {
	return _default.SetErrorTemplateName(name)
}

func Reload() *Manager {
	return _default.Reload()
}

func Html(wr http.ResponseWriter, name string, data any) error {
	return _default.Html(wr, name, data)
}

func HtmlGin(ctx *gin.Context, code int, name string, data any) {
	_default.HtmlGin(ctx, code, name, data)
}

func PrintDefinedTemplates() {
	_default.PrintDefinedTemplates()
}

func DefinedTemplates() []string {
	return _default.DefinedTemplates()
}
