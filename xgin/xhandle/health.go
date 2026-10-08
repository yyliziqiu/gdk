package xhandle

import (
	"github.com/gin-gonic/gin"

	"github.com/yyliziqiu/gdk/xgin/xresp"
)

func Health(ctx *gin.Context) {
	xresp.Ok(ctx)
}
