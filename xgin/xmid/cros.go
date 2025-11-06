package xmid

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/yyliziqiu/gdk/xgin/xresp"
)

/*
参考： https://developer.mozilla.org/zh-CN/docs/Web/HTTP/Headers

Only the 7 CORS-safelisted response headers are exposed:
Cache-Control
Content-Language
Content-Length
Content-Type
Expires
Last-Modified
Pragma

CORS-safelisted request header is one of the following HTTP headers:
Accept
Accept-Language
Content-Language
Content-Type
*/

type CrosConfig struct {
	MaxAge           string
	Origin           string
	ExposeHeaders    string
	AllowMethods     string
	AllowHeaders     string
	AllowCredentials string
}

var _crosConfig = &CrosConfig{
	MaxAge:           "86400",
	Origin:           "*",
	ExposeHeaders:    "",
	AllowMethods:     "OPTIONS, HEAD, GET, POST, PUT, PATCH, DELETE",
	AllowHeaders:     "*",
	AllowCredentials: "false",
}

func Cros(c *CrosConfig) gin.HandlerFunc {
	if c == nil {
		c = _crosConfig
	}

	return func(ctx *gin.Context) {
		ctx.Header("Access-Control-Allow-Origin", c.Origin)
		ctx.Header("Access-Control-Expose-Headers", c.ExposeHeaders)
		ctx.Header("Access-Control-Allow-Credentials", c.AllowCredentials)

		if ctx.Request.Method == http.MethodOptions {
			ctx.Header("Access-Control-Max-Age", c.MaxAge)
			ctx.Header("Access-Control-Allow-Methods", c.AllowMethods)
			ctx.Header("Access-Control-Allow-Headers", c.AllowHeaders)
		}

		if ctx.Request.Method == http.MethodOptions {
			xresp.AbortOk(ctx)
		}
	}
}
