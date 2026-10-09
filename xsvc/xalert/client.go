package ext

import (
	"github.com/sirupsen/logrus"

	"github.com/yyliziqiu/gdk/xgin/xresp"
	"github.com/yyliziqiu/gdk/xhttp"
	"github.com/yyliziqiu/gdk/xlog"
)

type Client struct {
	dev bool
	cli *xhttp.Client
}

func NewClient(logger *logrus.Logger, domain string, dev bool) *Client {
	return &Client{
		dev: dev,
		cli: xhttp.New(
			xhttp.Logger(logger),
			xhttp.Prefix(domain),
			xhttp.Error(xresp.ErrorResult{})),
	}
}

func (t *Client) Send(kty string, key string, args Args) error {
	if t.dev {
		xlog.Infof("[Client.Send] Simulated, type: %s, key: %s, args: %v.", kty, key, args)
		return nil
	}

	req := NewRequest(kty, key, args)

	return t.cli.Post("/alerts/v0/command/alert", nil, req, nil)
}

func (t *Client) ClearTimer(kty string, key string, args Args) error {
	if t.dev {
		xlog.Infof("[Client.ClearTimer] Simulated, type: %s, key: %s, args: %v.", kty, key, args)
		return nil
	}

	req := NewRequest(kty, key, args)

	return t.cli.Post("/alerts/v0/command/clear-timer", nil, req, nil)
}

func (t *Client) Send2(kty string, kvs ...any) error {
	args, err := NewArgs(kvs...)
	if err != nil {
		return err
	}
	return t.Send(kty, "", args)
}

func (t *Client) ClearTimer2(kty string, kvs ...any) error {
	args, err := NewArgs(kvs...)
	if err != nil {
		return err
	}
	return t.ClearTimer(kty, "", args)
}
