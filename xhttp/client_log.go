package xhttp

import (
	"fmt"
	"net/http"
	"net/http/httputil"

	"github.com/yyliziqiu/gdk/xstr"
)

// 记录请求日志
const (
	logFormat1 = "Request done(%d), method: %s, url: %s, header: %s, request: %s, response: %s, cost: %s."
	logFormat2 = "Request fail(%d), method: %s, url: %s, header: %s, request: %s, response: %s, error: %v, cost: %s."
)

func (c *Client) logRequest(req *http.Request, reqbody []byte, status int, resbody []byte, err error, cost string) {
	if c.logger == nil {
		return
	}

	headers := "-"
	if c.logHeader {
		headers = SerialHeader(req.Header)
	}

	if err == nil {
		c.logInfo(logFormat1, status, req.Method, req.URL, headers, string(reqbody), string(resbody), cost)
	} else {
		c.logInfo(logFormat2, status, req.Method, req.URL, headers, string(reqbody), string(resbody), err, cost)
	}
}

// 记录 info 日志
func (c *Client) logInfo(format string, args ...any) {
	if c.logger == nil {
		return
	}
	c.logger.Info(c.logFormat(format, args...))
}

// 记录 warn 日志
func (c *Client) logWarn(format string, args ...any) {
	if c.logger == nil {
		return
	}
	c.logger.Warn(c.logFormat(format, args...))
}

// 格式化日志
func (c *Client) logFormat(format string, args ...any) string {
	msg := fmt.Sprintf(format, args...)
	if c.logLength <= 0 {
		return ""
	}

	msg = xstr.TruncateUtf8(msg, c.logLength)

	if c.logEscape {
		msg = c.logChange.Replace(msg)
	}

	return msg
}

// 向控制台打印请求内容
func (c *Client) dumpRequest(req *http.Request) {
	if !c.dumps {
		return
	}
	bs, err := httputil.DumpRequestOut(req, true)
	if err != nil {
		fmt.Printf("Dump request failed, error: %v\n", err)
		return
	}
	fmt.Println("\n========== Request Begin ==========")
	fmt.Print(string(bs))
	fmt.Println("\n========== Request End ==========")
}

// 向控制台打印响应内容
func (c *Client) dumpResponse(res *http.Response) {
	if !c.dumps {
		return
	}
	bs, err := httputil.DumpResponse(res, true)
	if err != nil {
		fmt.Printf("Dump response failed, error: %v", err)
		return
	}
	fmt.Println("\n========== Response Begin ==========")
	fmt.Print(string(bs))
	fmt.Println("\n========== Response End ==========")
}
