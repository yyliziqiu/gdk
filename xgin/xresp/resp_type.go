package xresp

import (
	"fmt"
	"net/http"

	"github.com/yyliziqiu/gdk/xerr"
)

type ErrorResult struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e ErrorResult) Error() string {
	return fmt.Sprintf("#%s %s", e.Code, e.Message)
}

func NewErrorResult(code string, message string) ErrorResult {
	return ErrorResult{Code: code, Message: message}
}

func NewErrorResult2(err error, verbose bool) (int, ErrorResult) {
	var (
		status  = http.StatusBadRequest
		code    = xerr.BadRequest.Code
		message = xerr.BadRequest.Message
	)

	zerr, ok := err.(*xerr.Error)
	if ok {
		status, code, message = zerr.Http()
	} else if verbose {
		message = err.Error()
	}

	return status, NewErrorResult(code, message)
}
