package xhttp

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"strings"
)

func (cli *Client) logRequest(req *http.Request, res *http.Response, reqBody []byte, resBody []byte, err error, cost string) {
	if cli.logger == nil {
		return
	}

	headers, reqSb, resSb := SerializeHeader(req.Header), "", ""
	if len(reqBody) > 0 {
		reqSb = string(reqBody)
	}
	if len(resBody) > 0 {
		resSb = string(resBody)
	}

	if err == nil {
		cli.logInfo("Request succeed(%d), method: %s, url: %s, header: %s, request: %s, response: %s, cost: %s.",
			res.StatusCode, req.Method, req.URL, headers, reqSb, resSb, cost)
	} else {
		cli.logWarn("Request failed(%d), method: %s, url: %s, header: %s, request: %s, response: %s, error: %v, cost: %s.",
			res.StatusCode, req.Method, req.URL, headers, reqSb, resSb, err, cost)
	}
}

func (cli *Client) logInfo(format string, args ...interface{}) {
	if cli.logger == nil {
		return
	}
	message := cli.logCheck(fmt.Sprintf(format, args...))
	cli.logger.Info(message)
}

func (cli *Client) logWarn(format string, args ...interface{}) {
	if cli.logger == nil {
		return
	}
	message := cli.logCheck(fmt.Sprintf(format, args...))
	cli.logger.Warn(message)
}

func (cli *Client) logCheck(log string) string {
	if cli.logLength <= 0 {
		return ""
	}
	if len(log) > cli.logLength {
		log = log[:cli.logLength]
	}
	if cli.logEscape {
		log = strings.ReplaceAll(log, "\t", "\\t")
		log = strings.ReplaceAll(log, "\r", "\\r")
		log = strings.ReplaceAll(log, "\n", "\\n")
	}
	return log
}

func (cli *Client) dumpRequest(req *http.Request) {
	if !cli.dumps {
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

func (cli *Client) dumpResponse(res *http.Response) {
	if !cli.dumps {
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
