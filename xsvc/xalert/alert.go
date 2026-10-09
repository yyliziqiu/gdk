package xalert

import (
	"github.com/sirupsen/logrus"
)

var _cli *Client

func Init(domain string, dev bool, logger *logrus.Logger) {
	_cli = NewClient(domain, dev, logger)
}

func Alert(kty string, key string, args Args) error {
	return _cli.Alert(kty, key, args)
}

func Alert2(kty string, args Args) error {
	return _cli.Alert(kty, "", args)
}

func Alert3(kty string, kvs ...any) error {
	return _cli.Alert(kty, "", NA2(kvs...))
}

func ClearTimer(kty string, key string, args Args) error {
	return _cli.ClearTimer(kty, key, args)
}

func ClearTimer2(kty string, key string) error {
	return _cli.ClearTimer(kty, key, nil)
}

func ClearTimer3(kty string, args Args) error {
	return _cli.ClearTimer(kty, "", args)
}
