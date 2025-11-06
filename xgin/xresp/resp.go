package xresp

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/yyliziqiu/gdk/xerr"
)

// ============ Response ============

func Response(ctx *gin.Context, statusCode int, data interface{}) {
	ctx.JSON(statusCode, data)
}

func ResponseError(ctx *gin.Context, statusCode int, code string, message string) {
	ctx.JSON(statusCode, NewErrorResult(code, message))
}

func Ok(ctx *gin.Context) {
	ctx.String(http.StatusOK, "")
}

func Result(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusOK, data)
}

func Error(ctx *gin.Context, err error) {
	ctx.JSON(NewErrorResult2(err, false))
}

func ErrorVerbose(ctx *gin.Context, err error) {
	ctx.JSON(NewErrorResult2(err, true))
}

func ErrorString(ctx *gin.Context, message string) {
	ctx.JSON(http.StatusBadRequest, NewErrorResult(xerr.BadRequest.Code, message))
}

// ============ Abort ============

func AbortOk(ctx *gin.Context) {
	ctx.AbortWithStatus(http.StatusOK)
}

func AbortResult(ctx *gin.Context, data interface{}) {
	ctx.AbortWithStatusJSON(http.StatusOK, data)
}

func AbortError(ctx *gin.Context, err error) {
	ctx.AbortWithStatusJSON(NewErrorResult2(err, false))
}

func AbortErrorVerbose(ctx *gin.Context, err error) {
	ctx.AbortWithStatusJSON(NewErrorResult2(err, true))
}

func AbortErrorString(ctx *gin.Context, message string) {
	ctx.AbortWithStatusJSON(http.StatusBadRequest, NewErrorResult(xerr.BadRequest.Code, message))
}

func AbortBadRequest(ctx *gin.Context) {
	ctx.AbortWithStatusJSON(NewErrorResult2(xerr.BadRequest, false))
}

func AbortUnauthorized(ctx *gin.Context) {
	ctx.AbortWithStatusJSON(NewErrorResult2(xerr.Unauthorized, false))
}

func AbortForbidden(ctx *gin.Context) {
	ctx.AbortWithStatusJSON(NewErrorResult2(xerr.Forbidden, false))
}

func AbortNotFound(ctx *gin.Context) {
	ctx.AbortWithStatusJSON(NewErrorResult2(xerr.NotFound, false))
}

func AbortMethodNotAllowed(ctx *gin.Context) {
	ctx.AbortWithStatusJSON(NewErrorResult2(xerr.MethodNotAllowed, false))
}

func AbortInternalServerError(ctx *gin.Context) {
	ctx.AbortWithStatusJSON(NewErrorResult2(xerr.InternalServerError, false))
}
