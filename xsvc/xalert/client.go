package xalert

import (
	"github.com/sirupsen/logrus"

	"github.com/yyliziqiu/gdk/xgin/xresp"
	"github.com/yyliziqiu/gdk/xhttp"
	"github.com/yyliziqiu/gdk/xlog"
)

type Request struct {
	KeyType string `json:"key_type"`
	Key     string `json:"key"`
	Args    Args   `json:"args"`
}

func req(keyType string, key string, args Args) Request {
	if args == nil {
		args = Args{}
	}

	return Request{
		KeyType: keyType,
		Key:     key,
		Args:    args,
	}
}

type Client struct {
	dev bool
	cli *xhttp.Client
}

func NewClient(domain string, dev bool, logger *logrus.Logger) *Client {
	prefix := domain + "/alerts/v0/command"

	return &Client{
		dev: dev,
		cli: xhttp.New(xhttp.Logger(logger), xhttp.Prefix(prefix), xhttp.Error(xresp.ErrorResult{})),
	}
}

func (t *Client) Alert(kty string, key string, args Args) error {
	if t.dev {
		xlog.Infof("[Client.Alert] Simulated, type: %s, args: %v.", kty, args)
		return nil
	}

	return t.cli.Post("/alert", nil, req(kty, key, args), nil)
}

func (t *Client) ClearTimer(kty string, key string, args Args) error {
	if t.dev {
		xlog.Infof("[Client.ClearTimer] Simulated, type: %s, args: %v.", kty, args)
		return nil
	}

	return t.cli.Post("/clear-timer", nil, req(kty, key, args), nil)
}
