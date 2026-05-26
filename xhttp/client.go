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

	"github.com/sirupsen/logrus"

	"github.com/yyliziqiu/gdk/xtime"
	"github.com/yyliziqiu/gdk/xutil"
)

type Client struct {
	client        *http.Client                   // 客户端
	logger        *logrus.Logger                 // 如果为 nil，则不记录日志
	format        string                         // 响应报文格式
	prefix        string                         // URL 前缀
	errtyp        reflect.Type                   // 响应失败时的 JSON 结构，必须实现 error 接口且不能设置为指针。在成功和失败响应的 JSON 结构不一致时设置
	dumps         bool                           // 将 HTTP 报文打印到控制台，调试用
	logLength     [2]int                         // 请求体或响应体的最大日志长度，此处是字符数而非字节数
	logAlways     bool                           // 是否记录所有响应体日志。默认为 true，若接口中存在二进制数据流响应需要将此值置为 false
	logHeader     uint                           // 日志中是否记录请求头或响应头
	logForbid     []string                       // 当 logHeader 开启时，设置需要屏蔽的头部字段
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
		logLength: [2]int{1024, 1024},
		logAlways: true,
		logHeader: LogHeaderNone,
		logForbid: []string{},
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

	if _, _, err = c.doRequest(req, out, nil); err != nil {
		return err
	}

	return nil
}

func (c *Client) newRequest(method string, path string, query url.Values, header http.Header, body io.Reader) (*http.Request, error) {
	if !strings.HasPrefix(path, "http://") && !strings.HasPrefix(path, "https://") {
		path = JoinUrl(c.prefix, path)
	}
	if len(query) > 0 {
		path = AppendQuery(path, query)
	}

	req, err := http.NewRequest(method, path, body)
	if err != nil {
		c.logWarn("New request failed, url: %s, error: %v.", path, err)
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

func (c *Client) doRequest(req *http.Request, out any, reqbody []byte) (*http.Response, []byte, error) {
	c.dumpRequest(req)

	tim := xtime.NewTimer()

	res, err := c.client.Do(req)
	if err != nil {
		c.logWarn("Do request failed, url: %s, error: %v.", req.URL, err)
		return nil, nil, err
	}

	resbody, err := c.handleResponse(res, out)
	res.Body.Close()

	if c.logAlways || out != nil || IsTextType(res.Header.Get("Content-Type")) {
		c.logRequest(req, reqbody, res, resbody, err, tim.Stops())
	} else {
		c.logRequest(req, reqbody, res, nil, err, tim.Stops())
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

func (c *Client) handleTextResponse(status int, body []byte, out any) error {
	if status/100 != 2 {
		return newResponseError(status, string(body))
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

func (c *Client) handleJsonResponse(status int, body []byte, out any) error {
	if status/100 == 1 {
		return nil
	}

	if status/100 == 2 {
		if out != nil {
			if err := json.Unmarshal(body, out); err != nil {
				return fmt.Errorf("unmarshal response error [%v]", err)
			}
			if jr, ok := out.(JsonResponse); ok {
				if jr.Failed() {
					if err2, ok2 := out.(error); ok2 {
						return err2
					}
					return newResponseError(status, string(body))
				}
			}
		}
		return nil
	}

	if status/100 == 3 {
		return errors.New("this is a redirect response")
	}

	if c.errtyp != nil {
		ret := reflect.New(c.errtyp).Interface()
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

	return newResponseError(status, string(body))
}

func (c *Client) post(method string, path string, query url.Values, header http.Header, in any, out any) error {
	var (
		err     error
		reqbody []byte
	)

	if in != nil {
		reqbody, err = json.Marshal(in)
		if err != nil {
			return fmt.Errorf("marshal request body error [%v]", err)
		}
	}

	req, err := c.newRequest(method, path, query, header, bytes.NewReader(reqbody))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	if _, _, err = c.doRequest(req, out, reqbody); err != nil {
		return err
	}

	return nil
}

// Get restful get
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

// Post restful post
func (c *Client) Post(path string, header http.Header, in any, out any) error {
	return c.post(http.MethodPost, path, nil, header, in, out)
}

// Put restful put
func (c *Client) Put(path string, header http.Header, in any, out any) error {
	return c.post(http.MethodPut, path, nil, header, in, out)
}

// Patch restful patch
func (c *Client) Patch(path string, header http.Header, in any, out any) error {
	return c.post(http.MethodPatch, path, nil, header, in, out)
}

// Delete restful delete
func (c *Client) Delete(path string, query url.Values, header http.Header, out any) error {
	return c.get(http.MethodDelete, path, query, header, out)
}

// GetBinary 获取二进制流数据。用来下载图片、视频等
func (c *Client) GetBinary(path string, query url.Values, header http.Header) ([]byte, string, error) {
	req, err := c.newRequest(http.MethodGet, path, query, header, nil)
	if err != nil {
		return nil, "", err
	}

	res, resbody, err := c.doRequest(req, nil, nil)
	if err != nil {
		return nil, "", err
	}

	restype := res.Header.Get("Content-Type")

	return resbody, restype, nil
}

// PostJson 发送 JSON 请求
func (c *Client) PostJson(path string, header http.Header, in []byte, out any) error {
	req, err := c.newRequest(http.MethodPost, path, nil, header, bytes.NewReader(in))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	if _, _, err = c.doRequest(req, out, in); err != nil {
		return err
	}

	return nil
}

// PostForm 发送 application/x-www-form-urlencoded 表单请求
func (c *Client) PostForm(path string, header http.Header, in url.Values, out any) error {
	reqbody := in.Encode()

	req, err := c.newRequest(http.MethodPost, path, nil, header, strings.NewReader(reqbody))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	reqbody2, _ := url.QueryUnescape(reqbody)

	if _, _, err = c.doRequest(req, out, []byte(reqbody2)); err != nil {
		return err
	}

	return nil
}

// PostData 发送 multipart/form-data 表单请求
func (c *Client) PostData(path string, header http.Header, values map[string]string, files map[string]string, out any) error {
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

	req, err := c.newRequest(http.MethodPost, path, nil, header, &buf)
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

	if _, _, err = c.doRequest(req, out, reqbody); err != nil {
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

// PostFile 发送单个文件
func (c *Client) PostFile(path string, header http.Header, values map[string]string, field string, filepath string, out any) error {
	files := map[string]string{field: filepath}
	return c.PostData(path, header, values, files, out)
}

// PostBinary 发送二进制流数据
func (c *Client) PostBinary(path string, header http.Header, mimeType string, in io.Reader, out any) error {
	req, err := c.newRequest(http.MethodPost, path, nil, header, in)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", mimeType)

	if _, _, err = c.doRequest(req, out, nil); err != nil {
		return err
	}

	return nil
}

// PostStream 以 multipart/form-data 形式发送二进制流数据。相较于 PostData 发送本地文件，此方法发送的是 io.Reader 中的数据流
func (c *Client) PostStream(path string, header http.Header, values map[string]string, field string, filename string, mimeType string, stream io.Reader, out any) error {
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
		h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, Escape(field), Escape(filename)))
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

	req, err := c.newRequest(http.MethodPost, path, nil, header, &buf)
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

	if _, _, err = c.doRequest(req, out, reqbody); err != nil {
		return err
	}

	return nil
}
