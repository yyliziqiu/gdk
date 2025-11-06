package xerr

import (
	"fmt"
	"net/http"
)

type Error struct {
	status  int
	Code    string
	Message string
}

func New(code string, message string) *Error {
	return &Error{Code: code, Message: message}
}

func New2(status int, code string, message string) *Error {
	return &Error{status: status, Code: code, Message: message}
}

func (e *Error) Error() string {
	return e.Code + " " + e.Message
}

func (e *Error) Wrap(err error) *Error {
	return e.clone(fmt.Sprintf("%s [%v]", e.Message, err))
}

func (e *Error) clone(message string) *Error {
	return &Error{status: e.status, Code: e.Code, Message: message}
}

func (e *Error) Format(message string, a ...interface{}) *Error {
	return e.clone(fmt.Sprintf(message, a...))
}

func (e *Error) With(v interface{}) *Error {
	message := ""
	switch v.(type) {
	case error:
		message = v.(error).Error()
	case string:
		message = v.(string)
	default:
		message = fmt.Sprintf("%v", v)
	}
	return e.clone(message)
}

func (e *Error) WithFields(a ...interface{}) *Error {
	return e.clone(fmt.Sprintf(e.Message, a...))
}

func (e *Error) WithStatus(status int) *Error {
	e.status = status
	return e
}

func (e *Error) GetStatus() int {
	if e.status != 0 {
		return e.status
	}
	status := http.StatusBadRequest
	if e.Code[0] != 'A' {
		status = http.StatusInternalServerError
	}
	return status
}

func (e *Error) Http() (int, string, string) {
	return e.GetStatus(), e.Code, e.Message
}
