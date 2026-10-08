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

func NewClient(baseUrl string, logger *logrus.Logger, dev bool) *Client {
	return &Client{
		dev: dev,
		cli: xhttp.New(xhttp.Prefix(baseUrl), xhttp.Logger(logger), xhttp.Error(xresp.ErrorResult{})),
	}
}

func (t *Client) Send(kty string, key string, args Args) error {
	if t.dev {
		xlog.Infof("[AlertsApi.Send] Simulated, type: %s, key: %s, args: %v.", kty, key, args)
		return nil
	}

	form := newSendFrom(kty, key, args)

	return t.cli.Post("/alerts/v0/command/alert", nil, form, nil)
}

func (t *Client) ClearTimer(kty string, key string, args Args) error {
	if t.dev {
		xlog.Infof("[AlertsApi.ClearTimer] Simulated, type: %s, key: %s, args: %v.", kty, key, args)
		return nil
	}

	form := newSendFrom(kty, key, args)

	return t.cli.Post("/alerts/v0/command/clear-timer", nil, form, nil)
}

func (t *Client) Send2(kty string, list ...any) error {
	args, err := NewArgs(list...)
	if err != nil {
		return err
	}
	return t.Send(kty, "", args)
}

func (t *Client) ClearTimer2(kty string, list ...any) error {
	args, err := NewArgs(list...)
	if err != nil {
		return err
	}
	return t.ClearTimer(kty, "", args)
}

func (t *Client) Send3(kty string, key string, list ...any) error {
	args, err := NewArgs(list...)
	if err != nil {
		return err
	}
	return t.Send(kty, key, args)
}

func (t *Client) ClearTimer3(kty string, key string) error {
	return t.ClearTimer(kty, key, Args{})
}
