package xhttp

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"net/url"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/yyliziqiu/gdk/xtime"
	"github.com/yyliziqiu/gdk/xutil"
)

const (
	FormatJson = "json"
	FormatText = "text"
)

type Client struct {
	client        *http.Client
	logger        *logrus.Logger                 // 如果为 nil，则不记录日志
	format        string                         // 响应报文格式
	prefix        string                         // URL 前缀
	error         error                          // 响应失败时的 JSON 结构。在响应成功和失败时 JSON 结构不一致时设置，不能是指针
	dumps         bool                           // 将 HTTP 报文打印到控制台
	logLength     int                            // 最大日志长度
	logEscape     bool                           // 是否转换日志中的特殊字符
	requestBefore func(req *http.Request)        // 在发送请求前调用
	responseAfter func(res *http.Response) error // 在接收响应后调用
}

func New(options ...Option) *Client {
	cli := &http.Client{
		Timeout: 5 * time.Second,
	}

	client := &Client{
		client:        cli,
		logger:        nil,
		format:        FormatJson,
		prefix:        "",
		error:         nil,
		dumps:         false,
		logLength:     1024,
		logEscape:     false,
		requestBefore: nil,
		responseAfter: nil,
	}

	for _, option := range options {
		option(client)
	}

	return client
}

func (cli *Client) get(method string, path string, query url.Values, header http.Header, out interface{}) error {
	req, err := cli.newRequest(method, path, query, header, nil)
	if err != nil {
		return err
	}

	timer := xtime.NewTimer()

	res, err := cli.doRequest(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err := cli.handleResponse(res, out)

	cli.logRequest(req, res, nil, body, err, timer.Stops())

	return err
}

func (cli *Client) newRequest(method string, path string, query url.Values, header http.Header, body io.Reader) (*http.Request, error) {
	if !strings.HasPrefix(path, "http://") && !strings.HasPrefix(path, "https://") {
		path = JoinUrl(cli.prefix, path)
	}

	url2, err := AppendQuery(path, query)
	if err != nil {
		cli.logWarn("Append query failed, url: %s, query: %s, error: %v.", url2, query.Encode(), err)
		return nil, fmt.Errorf("append query error [%v]", err)
	}

	req, err := http.NewRequest(method, url2, body)
	if err != nil {
		cli.logWarn("New request failed, url: %s, error: %v.", url2, err)
		return nil, fmt.Errorf("new request error [%v]", err)
	}

	for key, values := range header {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	if cli.requestBefore != nil {
		cli.requestBefore(req)
	}

	return req, nil
}

func (cli *Client) doRequest(req *http.Request) (*http.Response, error) {
	cli.dumpRequest(req)

	res, err := cli.client.Do(req)
	if err != nil {
		cli.logWarn("Do request failed, url: %s, error: %v.", req.URL, err)
		return nil, err
	}

	return res, nil
}

func (cli *Client) handleResponse(res *http.Response, out interface{}) ([]byte, error) {
	cli.dumpResponse(res)

	if cli.responseAfter != nil {
		err := cli.responseAfter(res)
		if err != nil {
			return nil, err
		}
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("read response error [%v]", err)
	}

	switch cli.format {
	case FormatText:
		return body, cli.handleTextResponse(res.StatusCode, body, out)
	default:
		return body, cli.handleJsonResponse(res.StatusCode, body, out)
	}
}

func (cli *Client) handleJsonResponse(statusCode int, body []byte, out interface{}) error {
	if statusCode/100 == 2 {
		if out != nil {
			err := json.Unmarshal(body, out)
			if err != nil {
				return fmt.Errorf("unmarshal response error [%v]", err)
			}
			if jr, ok := out.(JsonResponse); ok {
				if jr.Failed() {
					err2, ok2 := out.(error)
					if ok2 {
						return err2
					}
					return newResponseError(statusCode, string(body))
				}
			}
		}
		return nil
	} else if statusCode/100 == 3 {
		return errors.New("this is a redirect response")
	} else {
		if cli.error != nil {
			ret := reflect.New(reflect.TypeOf(cli.error)).Interface()
			err := json.Unmarshal(body, ret)
			if err == nil {
				return ret.(error)
			}
		} else if out != nil {
			err := json.Unmarshal(body, out)
			if err == nil {
				err2, ok2 := out.(error)
				if ok2 {
					return err2
				}
			}
		}
		return newResponseError(statusCode, string(body))
	}
}

func (cli *Client) handleTextResponse(statusCode int, body []byte, out interface{}) error {
	if statusCode/100 != 2 {
		return newResponseError(statusCode, string(body))
	}

	if out == nil {
		return nil
	}

	bs, ok := out.(*[]byte)
	if !ok {
		return fmt.Errorf("response receiver must *[]byte type")
	}
	*bs = body

	return nil
}

func (cli *Client) post(method string, path string, query url.Values, header http.Header, in interface{}, out interface{}) error {
	if in == nil {
		in = struct{}{}
	}
	reqBody, err := json.Marshal(in)
	if err != nil {
		return fmt.Errorf("marshal request body error [%v]", err)
	}

	req, err := cli.newRequest(method, path, query, header, bytes.NewReader(reqBody))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	timer := xtime.NewTimer()

	res, err := cli.doRequest(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	resBody, err := cli.handleResponse(res, out)

	cli.logRequest(req, res, reqBody, resBody, err, timer.Stops())

	return err
}

// Get http get
//
// 若响应失败时 http 状态码为200：
//  1. 则 out 需要实现 JsonResponse 接口来判断响应是否成功
//     1-1. 若要自定义错误内容，则 out 需要实现 error 接口，否则错误信息将返回整个响应内容
//
// 若响应失败时 http 状态码为4**或5**：
//  1. 若响应成功和失败时响应的结构一至
//     1-1. 若要自定义错误内容，则 out 需要实现 error 接口，否则错误信息将返回整个响应内容
//  2. 若响应成功和失败时响应的结构不一致，则需要设置 Error(err error) 选项，err 为响应失败时的结构，注意 err 不能是指针
func (cli *Client) Get(path string, query url.Values, header http.Header, out interface{}) error {
	return cli.get(http.MethodGet, path, query, header, out)
}

// Post http post
func (cli *Client) Post(path string, query url.Values, header http.Header, in interface{}, out interface{}) error {
	return cli.post(http.MethodPost, path, query, header, in, out)
}

// Put http put
func (cli *Client) Put(path string, query url.Values, header http.Header, in interface{}, out interface{}) error {
	return cli.post(http.MethodPut, path, query, header, in, out)
}

// Patch http patch
func (cli *Client) Patch(path string, query url.Values, header http.Header, in interface{}, out interface{}) error {
	return cli.post(http.MethodPatch, path, query, header, in, out)
}

// Delete http delete
func (cli *Client) Delete(path string, query url.Values, header http.Header, out interface{}) error {
	return cli.get(http.MethodDelete, path, query, header, out)
}

// GetBinary 获取流数据
func (cli *Client) GetBinary(path string, query url.Values, header http.Header) ([]byte, string, error) {
	req, err := cli.newRequest(http.MethodGet, path, query, header, nil)
	if err != nil {
		return nil, "", err
	}

	timer := xtime.NewTimer()

	res, err := cli.doRequest(req)
	if err != nil {
		return nil, "", err
	}
	defer res.Body.Close()

	dat, err := io.ReadAll(res.Body)
	typ := res.Header.Get("Content-Type")

	cli.logRequest(req, res, nil, nil, err, timer.Stops())

	return dat, typ, err
}

// PostForm application/x-www-form-urlencoded 表单请求
func (cli *Client) PostForm(path string, query url.Values, header http.Header, in url.Values, out interface{}) error {
	reqBody := in.Encode()

	req, err := cli.newRequest(http.MethodPost, path, query, header, strings.NewReader(reqBody))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	timer := xtime.NewTimer()

	res, err := cli.doRequest(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	resBody, err := cli.handleResponse(res, out)

	reqBody, _ = url.QueryUnescape(reqBody)
	cli.logRequest(req, res, []byte(reqBody), resBody, err, timer.Stops())

	return err
}

// PostData multipart/form-data 表单请求
func (cli *Client) PostData(path string, query url.Values, header http.Header, values map[string]string, files map[string]string, out interface{}) error {
	var (
		buf    bytes.Buffer
		writer = multipart.NewWriter(&buf)
	)

	if len(values) > 0 {
		for key, value := range values {
			err := writer.WriteField(key, value)
			if err != nil {
				return err
			}
		}
	}
	if len(files) > 0 {
		for key, file := range files {
			err := cli.writeFormFile(writer, key, file)
			if err != nil {
				return err
			}
		}
	}
	err := writer.Close()
	if err != nil {
		return err
	}

	req, err := cli.newRequest(http.MethodPost, path, query, header, &buf)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	timer := xtime.NewTimer()

	res, err := cli.doRequest(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	resBody, err := cli.handleResponse(res, out)

	var cpy map[string]string
	var reqBody []byte
	if len(values) > 0 {
		cpy = values
		for key, file := range files {
			cpy[key] = file
		}
	} else {
		cpy = files
	}
	if len(cpy) == 0 {
		reqBody = []byte("{}")
	} else {
		reqBody, _ = json.Marshal(cpy)
	}
	cli.logRequest(req, res, reqBody, resBody, err, timer.Stops())

	return err
}

func (cli *Client) writeFormFile(writer *multipart.Writer, key string, path string) error {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	part, err := writer.CreateFormFile(key, file.Name())
	if err != nil {
		return err
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return err
	}

	return nil
}

// PostBinary 上传流数据
func (cli *Client) PostBinary(path string, query url.Values, header http.Header, mimeType string, in io.Reader, out interface{}) error {
	req, err := cli.newRequest(http.MethodPost, path, query, header, in)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", mimeType)

	timer := xtime.NewTimer()

	res, err := cli.doRequest(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	resBody, err := cli.handleResponse(res, out)

	cli.logRequest(req, res, nil, resBody, err, timer.Stops())

	return err
}

// PostStream 以 multipart/form-data 形式上传流数据
func (cli *Client) PostStream(path string, query url.Values, header http.Header, values map[string]string, field string, filename string, mimeType string, stream io.Reader, out interface{}) error {
	var (
		buf    bytes.Buffer
		writer = multipart.NewWriter(&buf)
	)

	if len(values) > 0 {
		for key, value := range values {
			err := writer.WriteField(key, value)
			if err != nil {
				return err
			}
		}
	}
	if mimeType == "" {
		mimeType = xutil.ParseMimeType(filename)
	}
	if stream != nil {
		h := make(textproto.MIMEHeader)
		h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, EscapeQuotes(field), EscapeQuotes(filename)))
		h.Set("Content-Type", mimeType)
		part, err := writer.CreatePart(h)
		if err != nil {
			return err
		}
		_, err = io.Copy(part, stream)
		if err != nil {
			return err
		}
	}

	err := writer.Close()
	if err != nil {
		return err
	}

	req, err := cli.newRequest(http.MethodPost, path, query, header, &buf)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	timer := xtime.NewTimer()

	res, err := cli.doRequest(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	resBody, err := cli.handleResponse(res, out)

	var reqBody []byte
	if len(values) == 0 {
		reqBody = []byte(fmt.Sprintf(`{"%s":"%s"}`, field, filename))
	} else {
		values[field] = filename
		reqBody, _ = json.Marshal(values)
	}
	cli.logRequest(req, res, reqBody, resBody, err, timer.Stops())

	return err
}
