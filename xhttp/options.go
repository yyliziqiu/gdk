package xhttp

import (
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"time"

	"github.com/sirupsen/logrus"
)

type Option func(cli *Client)

func WithClient(client *http.Client) Option {
	return func(cli *Client) {
		cli.client = client
	}
}

func Cookie(o *cookiejar.Options) Option {
	return func(cli *Client) {
		jar, _ := cookiejar.New(o)
		cli.client.Jar = jar
	}
}

func Timeout(timeout time.Duration) Option {
	return func(cli *Client) {
		cli.client.Timeout = timeout
	}
}

func DisableRedirect() Option {
	return func(cli *Client) {
		cli.client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
	}
}

func LimitRedirect(n int) Option {
	return func(cli *Client) {
		cli.client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			if len(via) >= n {
				return fmt.Errorf("stopped after %d redirects", n)
			}
			return nil
		}
	}
}

func Logger(logger *logrus.Logger) Option {
	return func(cli *Client) {
		cli.logger = logger
	}
}

func Format(format string) Option {
	return func(cli *Client) {
		cli.format = format
	}
}

func Prefix(prefix string) Option {
	return func(cli *Client) {
		cli.prefix = prefix
	}
}

func Error(error error) Option {
	return func(cli *Client) {
		cli.error = error
	}
}

func Dumps(enabled bool) Option {
	return func(cli *Client) {
		cli.dumps = enabled
	}
}

// BaseUrl deprecate
func BaseUrl(baseUrl string) Option {
	return func(cli *Client) {
		cli.prefix = baseUrl
	}
}

func LogLength(n int) Option {
	return func(cli *Client) {
		cli.logLength = n
	}
}

func LogEscape(enabled bool) Option {
	return func(cli *Client) {
		cli.logEscape = enabled
	}
}

func RequestBefore(f func(r *http.Request)) Option {
	return func(cli *Client) {
		cli.requestBefore = f
	}
}

func BasicAuth(username string, password string) Option {
	return func(cli *Client) {
		cli.requestBefore = func(req *http.Request) {
			req.SetBasicAuth(username, password)
		}
	}
}

func BearerToken(token string) Option {
	return func(cli *Client) {
		cli.requestBefore = func(req *http.Request) {
			req.Header.Set("Authorization", "Bearer "+token)
		}
	}
}

func ResponseAfter(f func(res *http.Response) error) Option {
	return func(cli *Client) {
		cli.responseAfter = f
	}
}
