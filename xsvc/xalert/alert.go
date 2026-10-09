package xalert

import (
	"github.com/sirupsen/logrus"
)

var _cli *Client

func Init(domain string, dev bool, logger *logrus.Logger) {
	_cli = NewClient(domain, dev, logger)
}

// Alert 若 key 为空，将自动以 args 中@开头的键构建 key
func Alert(kty string, key string, args Args) error {
	return _cli.Alert(kty, key, args)
}

// Alert2 自动用 args 构建 key 时使用
func Alert2(kty string, args Args) error {
	return _cli.Alert(kty, "", args)
}

// Alert3 自动创建 args 并用 args 自动构建 key
func Alert3(kty string, kvs ...any) error {
	return _cli.Alert(kty, "", NA2(kvs...))
}

// ClearTimer 手动指定 key 时使用
func ClearTimer(kty string, key string) error {
	return _cli.ClearTimer(kty, key, nil)
}

// ClearTimer2 自动用 args 构建 key 时使用
func ClearTimer2(kty string, args Args) error {
	return _cli.ClearTimer(kty, "", args)
}
