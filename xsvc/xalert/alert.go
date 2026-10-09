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

func ClearTimer(kty string, key string, args Args) error {
	return _cli.ClearTimer(kty, key, args)
}

func Alert2(kty string, kvs ...any) error {
	args, err := NewArgs(kvs...)
	if err != nil {
		return err
	}

	return _cli.Alert(kty, "", args)
}

func ClearTimer2(kty string, kvs ...any) error {
	args, err := NewArgs(kvs...)
	if err != nil {
		return err
	}

	return _cli.ClearTimer(kty, "", args)
}
