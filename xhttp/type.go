package xhttp

import (
	"fmt"
	"strings"
)

const (
	FormatJson = "json"
	FormatText = "text"
)

const (
	logFormat1 = "Request succeed(%d), method: %s, url: %s, header: %s, request: %s, response: %s, cost: %s."
	logFormat2 = "Request failed(%d), method: %s, url: %s, header: %s, request: %s, response: %s, error: %v, cost: %s."
)

var _replacer = strings.NewReplacer(
	"\t", "\\t",
	"\r", "\\r",
	"\n", "\\n",
)

type JsonResponse interface {
	Failed() bool
}

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
