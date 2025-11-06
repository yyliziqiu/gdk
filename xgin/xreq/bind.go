package xreq

import (
	"github.com/gin-gonic/gin"

	"github.com/yyliziqiu/gdk/xerr"
	"github.com/yyliziqiu/gdk/xgin"
	"github.com/yyliziqiu/gdk/xgin/xresp"
)

func bind(ctx *gin.Context, form interface{}, verbose bool) bool {
	err := ctx.ShouldBind(form)
	if err != nil {
		if logger := xgin.GetLogger(); logger != nil {
			logger.Warnf("Bind failed, path: %s, error: %v.", ctx.FullPath(), err)
		}
		if verbose {
			xresp.Error(ctx, xerr.ParametersError.Wrap(err))
		} else {
			xresp.Error(ctx, xerr.ParametersError)
		}
		return false
	}
	return true
}

func Bind(ctx *gin.Context, form interface{}) bool {
	return bind(ctx, form, false)
}

func BindVerbose(ctx *gin.Context, form interface{}) bool {
	return bind(ctx, form, true)
}
