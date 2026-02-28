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
	"path/filepath"
	"reflect"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/yyliziqiu/gdk/xtime"
	"github.com/yyliziqiu/gdk/xutil"
)

type Client struct {
	client        *http.Client                   // 客户端
	logger        *logrus.Logger                 // 如果为 nil，则不记录日志
	format        string                         // 响应报文格式
	prefix        string                         // URL 前缀
	error         error                          // 响应失败时的 JSON 结构。在响应成功和失败时 JSON 结构不一致时设置，不能是指针
	dumps         bool                           // 将 HTTP 报文打印到控制台
	logLength     int                            // 最大日志长度
	logHeader     bool                           // 日志中是否保存请求头
	logEscape     bool                           // 是否转换日志中的特殊字符
	logChange     *strings.Replacer              // 字符替换器
	requestBefore func(req *http.Request)        // 在发送请求前调用
	responseAfter func(res *http.Response) error // 在接收响应后调用
}

func New(options ...Option) *Client {
	client := &Client{
		client:    DefaultClient(),
		logger:    nil,
		format:    FormatJson,
		prefix:    "",
		logLength: 2048,
		logHeader: false,
		logEscape: false,
		logChange: _replacer,
	}

	for _, option := range options {
		option(client)
	}

	return client
}

func (c *Client) get(method string, path string, query url.Values, header http.Header, out any) error {
	req, err := c.newRequest(method, path, query, header, nil)
	if err != nil {
		return err
	}

	if _, _, err = c.doRequest(req, nil, out, true); err != nil {
		return err
	}

	return nil
}

func (c *Client) newRequest(method string, path string, query url.Values, header http.Header, body io.Reader) (*http.Request, error) {
	if !strings.HasPrefix(path, "http://") && !strings.HasPrefix(path, "https://") {
		path = JoinUrl(c.prefix, path)
	}

	url2, err := AppendQuery(path, query)
	if err != nil {
		c.logWarn("Append query failed, url: %s, query: %s, error: %v.", url2, query.Encode(), err)
		return nil, fmt.Errorf("append query error [%v]", err)
	}

	req, err := http.NewRequest(method, url2, body)
	if err != nil {
		c.logWarn("New request failed, url: %s, error: %v.", url2, err)
		return nil, fmt.Errorf("new request error [%v]", err)
	}

	for k, vs := range header {
		for _, v := range vs {
			req.Header.Add(k, v)
		}
	}

	if c.requestBefore != nil {
		c.requestBefore(req)
	}

	return req, nil
}

func (c *Client) doRequest(req *http.Request, reqbody []byte, out any, logResbody bool) (*http.Response, []byte, error) {
	c.dumpRequest(req)

	tim := xtime.NewTimer()

	res, err := c.client.Do(req)
	if err != nil {
		c.logWarn("Do request failed, url: %s, error: %v.", req.URL, err)
		return nil, nil, err
	}

	resbody, err := c.handleResponse(res, out)
	res.Body.Close()

	if logResbody {
		c.logRequest(req, reqbody, res.StatusCode, resbody, err, tim.Stops())
	} else {
		c.logRequest(req, reqbody, res.StatusCode, nil, err, tim.Stops())
	}

	return res, resbody, err
}

func (c *Client) handleResponse(res *http.Response, out any) ([]byte, error) {
	c.dumpResponse(res)

	if c.responseAfter != nil {
		err := c.responseAfter(res)
		if err != nil {
			return nil, err
		}
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("read response error [%v]", err)
	}

	switch c.format {
	case FormatText:
		return body, c.handleTextResponse(res.StatusCode, body, out)
	default:
		return body, c.handleJsonResponse(res.StatusCode, body, out)
	}
}

func (c *Client) handleTextResponse(statusCode int, body []byte, out any) error {
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

func (c *Client) handleJsonResponse(statusCode int, body []byte, out any) error {
	if statusCode/100 == 2 {
		if out != nil {
			if err := json.Unmarshal(body, out); err != nil {
				return fmt.Errorf("unmarshal response error [%v]", err)
			}
			if jr, ok := out.(JsonResponse); ok {
				if jr.Failed() {
					if err2, ok2 := out.(error); ok2 {
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
		if c.error != nil {
			ret := reflect.New(reflect.TypeOf(c.error)).Interface()
			if err := json.Unmarshal(body, ret); err == nil {
				return ret.(error)
			}
		} else if out != nil {
			if err := json.Unmarshal(body, out); err == nil {
				if err2, ok2 := out.(error); ok2 {
					return err2
				}
			}
		}
		return newResponseError(statusCode, string(body))
	}
}

func (c *Client) post(method string, path string, query url.Values, header http.Header, in any, out any) error {
	if in == nil {
		in = struct{}{}
	}

	reqbody, err := json.Marshal(in)
	if err != nil {
		return fmt.Errorf("marshal request body error [%v]", err)
	}

	req, err := c.newRequest(method, path, query, header, bytes.NewReader(reqbody))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	if _, _, err = c.doRequest(req, reqbody, out, true); err != nil {
		return err
	}

	return nil
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
func (c *Client) Get(path string, query url.Values, header http.Header, out any) error {
	return c.get(http.MethodGet, path, query, header, out)
}

// Post http post
func (c *Client) Post(path string, query url.Values, header http.Header, in any, out any) error {
	return c.post(http.MethodPost, path, query, header, in, out)
}

// Put http put
func (c *Client) Put(path string, query url.Values, header http.Header, in any, out any) error {
	return c.post(http.MethodPut, path, query, header, in, out)
}

// Patch http patch
func (c *Client) Patch(path string, query url.Values, header http.Header, in any, out any) error {
	return c.post(http.MethodPatch, path, query, header, in, out)
}

// Delete http delete
func (c *Client) Delete(path string, query url.Values, header http.Header, out any) error {
	return c.get(http.MethodDelete, path, query, header, out)
}

// GetBinary 获取流数据
func (c *Client) GetBinary(path string, query url.Values, header http.Header) ([]byte, string, error) {
	req, err := c.newRequest(http.MethodGet, path, query, header, nil)
	if err != nil {
		return nil, "", err
	}

	res, resbody, err := c.doRequest(req, nil, nil, false)
	if err != nil {
		return nil, "", err
	}

	return resbody, res.Header.Get("Content-Type"), nil
}

// PostForm application/x-www-form-urlencoded 表单请求
func (c *Client) PostForm(path string, query url.Values, header http.Header, in url.Values, out any) error {
	reqbody := in.Encode()

	req, err := c.newRequest(http.MethodPost, path, query, header, strings.NewReader(reqbody))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	reqbody2, _ := url.QueryUnescape(reqbody)

	if _, _, err = c.doRequest(req, []byte(reqbody2), out, true); err != nil {
		return err
	}

	return nil
}

// PostData multipart/form-data 表单请求
func (c *Client) PostData(path string, query url.Values, header http.Header, values map[string]string, files map[string]string, out any) error {
	var (
		buf    bytes.Buffer
		writer = multipart.NewWriter(&buf)
	)

	if len(values) > 0 {
		for key, value := range values {
			if err := writer.WriteField(key, value); err != nil {
				return err
			}
		}
	}
	if len(files) > 0 {
		for key, file := range files {
			if err := c.writeFormFile(writer, key, file); err != nil {
				return err
			}
		}
	}
	if err := writer.Close(); err != nil {
		return err
	}

	req, err := c.newRequest(http.MethodPost, path, query, header, &buf)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	cpy := make(map[string]string)
	for key, val := range values {
		cpy[key] = val
	}
	for key, file := range files {
		cpy[key] = file
	}
	var reqbody []byte
	if len(cpy) == 0 {
		reqbody = []byte("{}")
	} else {
		reqbody, _ = json.Marshal(cpy)
	}

	if _, _, err = c.doRequest(req, reqbody, out, true); err != nil {
		return err
	}

	return nil
}

func (c *Client) writeFormFile(writer *multipart.Writer, key string, path string) error {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	part, err := writer.CreateFormFile(key, file.Name())
	if err != nil {
		return err
	}

	if _, err = io.Copy(part, file); err != nil {
		return err
	}

	return nil
}

// PostBinary 上传流数据
func (c *Client) PostBinary(path string, query url.Values, header http.Header, mimeType string, in io.Reader, out any) error {
	req, err := c.newRequest(http.MethodPost, path, query, header, in)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", mimeType)

	if _, _, err = c.doRequest(req, nil, out, true); err != nil {
		return err
	}

	return nil
}

// PostStream 以 multipart/form-data 形式上传流数据
func (c *Client) PostStream(path string, query url.Values, header http.Header, values map[string]string, field string, filename string, mimeType string, stream io.Reader, out any) error {
	var (
		buf    bytes.Buffer
		writer = multipart.NewWriter(&buf)
	)

	if len(values) > 0 {
		for key, value := range values {
			if err := writer.WriteField(key, value); err != nil {
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
		if _, err = io.Copy(part, stream); err != nil {
			return err
		}
	}
	if err := writer.Close(); err != nil {
		return err
	}

	req, err := c.newRequest(http.MethodPost, path, query, header, &buf)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	var reqbody []byte
	if len(values) == 0 {
		reqbody = []byte(fmt.Sprintf(`{"%s":"%s"}`, field, filename))
	} else {
		values[field] = filename
		reqbody, _ = json.Marshal(values)
	}

	if _, _, err = c.doRequest(req, reqbody, out, true); err != nil {
		return err
	}

	return err
}

// PostFile 上传文件
func (c *Client) PostFile(path string, query url.Values, header http.Header, values map[string]string, field string, filepath string, out any) error {
	files := map[string]string{field: filepath}
	return c.PostData(path, query, header, values, files, out)
}

// ForwardBinary 转发二进制数据
func (c *Client) ForwardBinary(path string, query url.Values, header http.Header, src string, out any) error {
	data, typ, err := c.GetBinary(src, nil, nil)
	if err != nil {
		return err
	}
	return c.PostBinary(path, query, header, typ, bytes.NewReader(data), out)
}

// ForwardStream 以 multipart/form-data 形式转发流数据
func (c *Client) ForwardStream(path string, query url.Values, header http.Header, values map[string]string, field string, mimeTyp string, src string, out any) error {
	data, _, err := c.GetBinary(src, nil, nil)
	if err != nil {
		return err
	}
	return c.PostStream(path, query, header, values, field, filepath.Base(src), mimeTyp, bytes.NewReader(data), out)
}
