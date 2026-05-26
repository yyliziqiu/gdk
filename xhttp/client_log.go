package xhttp

import (
	"fmt"
	"net/http"
	"net/http/httputil"

	"github.com/yyliziqiu/gdk/xstr"
)

// 记录请求日志
func (c *Client) logRequest(req *http.Request, reqbody []byte, res *http.Response, resbody []byte, err error, cost string) {
	if c.logger == nil {
		return
	}

	var (
		status = res.StatusCode
		method = req.Method
		rawurl = req.URL.String()
		reqstr = xstr.TruncateUtf8(string(reqbody), c.logLength[0])
		resstr = xstr.TruncateUtf8(string(resbody), c.logLength[1])
	)

	switch c.logHeader {
	case LogHeaderNone:
		if err == nil {
			c.logInfo(logFormat01, status, method, rawurl, reqstr, resstr, cost)
		} else {
			c.logInfo(logFormat02, status, method, rawurl, reqstr, resstr, err, cost)
		}
	case LogHeaderReq:
		header := SerialHeader(req.Header, c.logForbid)
		if err == nil {
			c.logInfo(logFormat11, status, method, rawurl, header, reqstr, resstr, cost)
		} else {
			c.logInfo(logFormat12, status, method, rawurl, header, reqstr, resstr, err, cost)
		}
	case LogHeaderRes:
		header := SerialHeader(res.Header, c.logForbid)
		if err == nil {
			c.logInfo(logFormat21, status, method, rawurl, reqstr, header, resstr, cost)
		} else {
			c.logInfo(logFormat22, status, method, rawurl, reqstr, header, resstr, err, cost)
		}
	case LogHeaderBoth:
		reqhead := SerialHeader(req.Header, c.logForbid)
		reshead := SerialHeader(res.Header, c.logForbid)
		if err == nil {
			c.logInfo(logFormat31, status, method, rawurl, reqhead, reqstr, reshead, resstr, cost)
		} else {
			c.logInfo(logFormat32, status, method, rawurl, reqhead, reqstr, reshead, resstr, err, cost)
		}
	}
}

// 记录 info 日志
func (c *Client) logInfo(format string, args ...any) {
	if c.logger != nil {
		c.logger.Info(c.logPretty(fmt.Sprintf(format, args...)))
	}
}

// 记录 warn 日志
func (c *Client) logWarn(format string, args ...any) {
	if c.logger != nil {
		c.logger.Warnf(format, args...)
	}
}

// 格式化日志
func (c *Client) logPretty(msg string) string {
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
