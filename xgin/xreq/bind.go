package xreq

import (
	"github.com/gin-gonic/gin"

	"github.com/yyliziqiu/gdk/xerr"
	"github.com/yyliziqiu/gdk/xgin"
	"github.com/yyliziqiu/gdk/xgin/xresp"
)

func bind(ctx *gin.Context, form interface{}, verbose bool) bool {
	if err := ctx.ShouldBind(form); err != nil {
		xgin.GetLogger().Warnf("Bind failed, path: %s, error: %v.", ctx.FullPath(), err)
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

func bindQuery(ctx *gin.Context, query interface{}, verbose bool) bool {
	if err := ctx.BindQuery(query); err != nil {
		xgin.GetLogger().Warnf("Bind query failed, path: %s, error: %v.", ctx.FullPath(), err)
		if verbose {
			xresp.Error(ctx, xerr.ParametersError.Wrap(err))
		} else {
			xresp.Error(ctx, xerr.ParametersError)
		}
		return false
	}
	return true
}

func BindQuery(ctx *gin.Context, query interface{}) bool {
	return bindQuery(ctx, query, false)
}

func BindQueryVerbose(ctx *gin.Context, query interface{}) bool {
	return bindQuery(ctx, query, true)
}
