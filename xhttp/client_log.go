package xhttp

import (
	"fmt"
	"net/http"
	"net/http/httputil"
)

func (c *Client) logRequest(req *http.Request, res *http.Response, reqBody []byte, resBody []byte, err error, cost string) {
	if c.logger == nil {
		return
	}

	headers := "-"
	if c.logHeader {
		headers = SerializeHeader(req.Header)
	}
	reqBy := string(reqBody)
	resBy := string(resBody)

	if err == nil {
		c.logInfo(logFormat1, res.StatusCode, req.Method, req.URL, headers, reqBy, resBy, cost)
	} else {
		c.logWarn(logFormat2, res.StatusCode, req.Method, req.URL, headers, reqBy, resBy, err, cost)
	}
}

func (c *Client) logInfo(format string, args ...interface{}) {
	if c.logger != nil {
		c.logger.Info(c.logPretty(fmt.Sprintf(format, args...)))
	}
}

func (c *Client) logWarn(format string, args ...interface{}) {
	if c.logger != nil {
		c.logger.Warn(c.logPretty(fmt.Sprintf(format, args...)))
	}
}

func (c *Client) logPretty(msg string) string {
	if c.logLength <= 0 {
		return ""
	}
	if len(msg) > c.logLength {
		msg = msg[:c.logLength]
	}
	if c.logEscape {
		msg = c.replacer.Replace(msg)
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
	fmt.Println("\n---------- Request ----------")
	fmt.Printf(string(bs))
	fmt.Println("\n---------- Request End----------")
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
	fmt.Println("\n---------- Response ----------")
	fmt.Printf(string(bs))
	fmt.Println("\n---------- Response End----------")
}
