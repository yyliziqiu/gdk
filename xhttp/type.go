package xhttp

import (
	"fmt"
	"strings"
)

// 响应格式
const (
	FormatJson = "json"
	FormatText = "text"
)

// 字符替换映射
var _replacer = strings.NewReplacer(
	"\t", "\\t",
	"\r", "\\r",
	"\n", "\\n",
)

// JsonResponse 判断响应是否失败
type JsonResponse interface {
	Failed() bool
}

// ResponseError 默认响应错误
type ResponseError struct {
	status int
	errstr string
}

func newResponseError(status int, errstr string) *ResponseError {
	return &ResponseError{
		status: status,
		errstr: errstr,
	}
}

func (e ResponseError) Status() int {
	return e.status
}

func (e ResponseError) Error() string {
	return fmt.Sprintf("#%d %s", e.status, e.errstr)
}
