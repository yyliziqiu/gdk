package xhttp

import (
	"fmt"
	"net/http"
	"net/http/httputil"
)

func (c *Client) logRequest(req *http.Request, reqbody []byte, status int, resbody []byte, err error, cost string) {
	if c.logger == nil {
		return
	}

	headers := "_"
	if c.logHeader {
		headers = SerializeHeader(req.Header)
	}

	if err == nil {
		c.logInfo(logFormat1, status, req.Method, req.URL, headers, string(reqbody), string(resbody), cost)
	} else {
		c.logInfo(logFormat2, status, req.Method, req.URL, headers, string(reqbody), string(resbody), err, cost)
	}
}

func (c *Client) logInfo(format string, args ...any) {
	if c.logger == nil {
		return
	}

	c.logger.Info(c.logFormat(format, args...))
}

func (c *Client) logWarn(format string, args ...any) {
	if c.logger == nil {
		return
	}

	c.logger.Warn(c.logFormat(format, args...))
}

func (c *Client) logFormat(format string, args ...any) string {
	msg := fmt.Sprintf(format, args...)

	if c.logLength <= 0 {
		return ""
	}

	if len(msg) > c.logLength {
		msg = msg[:c.logLength]
	}

	if c.logEscape {
		msg = c.logReplace.Replace(msg)
	}

	return msg
}

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
