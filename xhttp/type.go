package xhttp

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// 响应格式
const (
	FormatJson = "json"
	FormatText = "text"
)

// 记录请求头和响应头
const (
	LogHeaderNone uint = 0x00
	LogHeaderReq       = 0x01
	LogHeaderRes       = 0x02
	LogHeaderBoth      = 0x03
)

// 请求日志格式
const (
	logFormat01 = "Request done(%d), method: %s, url: %s, reqbody: %s, resbody: %s, cost: %s."
	logFormat02 = "Request fail(%d), method: %s, url: %s, reqbody: %s, resbody: %s, error: %v, cost: %s."

	logFormat11 = "Request done(%d), method: %s, url: %s, header: %s, reqbody: %s, resbody: %s, cost: %s."
	logFormat12 = "Request fail(%d), method: %s, url: %s, header: %s, reqbody: %s, resbody: %s, error: %v, cost: %s."

	logFormat21 = "Request done(%d), method: %s, url: %s, reqbody: %s, header: %s, resbody: %s, cost: %s."
	logFormat22 = "Request fail(%d), method: %s, url: %s, reqbody: %s, header: %s, resbody: %s, error: %v, cost: %s."

	logFormat31 = "Request done(%d), method: %s, url: %s, reqhead: %s, reqbody: %s, reshead: %s, resbody: %s, cost: %s."
	logFormat32 = "Request fail(%d), method: %s, url: %s, reqhead: %s, reqbody: %s, reshead: %s, resbody: %s, error: %v, cost: %s."
)

// 字符替换映射
var _replacer = strings.NewReplacer(
	"\t", "\\t",
	"\r", "\\r",
	"\n", "\\n",
)

// 转义字符
var _quoteReplacer = strings.NewReplacer(
	"\\", "\\\\",
	`"`, "\\\"",
)

// Content-Type 文本类型
var _textTypes = []string{
	"application/json",
	"application/xml",
}

// HTTP 头分隔符
var _headerDivide = []byte{';', ' '}

// JsonResponse 判断响应是否失败
type JsonResponse interface {
	Failed() bool
}

// ResponseError 默认响应错误
type ResponseError struct {
	status int
	errstr string
}

func newResponseError(status int, errstr string) *ResponseError {
	return &ResponseError{
		status: status,
		errstr: errstr,
	}
}

func (e ResponseError) Status() int {
	return e.status
}

func (e ResponseError) Error() string {
	return fmt.Sprintf("#%d %s", e.status, e.errstr)
}

// DefaultClient 创建默认客户端
func DefaultClient() *http.Client {
	return &http.Client{
		Timeout: 5 * time.Second,
	}
}

// JoinUrl 拼接两段 URL
func JoinUrl(prefix string, suffix string) string {
	if suffix == "" {
		return prefix
	}

	l := strings.HasSuffix(prefix, "/")
	r := strings.HasPrefix(suffix, "/")

	if l && r {
		if len(suffix) < 2 {
			return prefix
		}
		return prefix + suffix[1:]
	}

	if l || r {
		return prefix + suffix
	}

	return prefix + "/" + suffix
}

// JoinUrls 拼接多段 URL
func JoinUrls(segments ...string) string {
	if len(segments) == 0 {
		return ""
	}

	rurl := segments[0]
	for i, seg := range segments {
		if i > 0 && seg != "" {
			l := strings.HasSuffix(rurl, "/")
			r := strings.HasPrefix(seg, "/")
			if l && r {
				if len(seg) > 1 {
					rurl += seg[1:]
				}
			} else if l || r {
				rurl += seg
			} else {
				rurl += "/" + seg
			}
		}
	}

	return rurl
}

// AppendQuery 向 URL 追加查询条件
func AppendQuery(rurl string, query url.Values) string {
	parsed, err := url.Parse(rurl)
	if err != nil {
		if strings.Contains(rurl, "?") {
			return rurl + "&" + query.Encode()
		}
		return rurl + "?" + query.Encode()
	}

	for k, v := range parsed.Query() {
		for _, s := range v {
			query.Add(k, s)
		}
	}

	parsed.RawQuery = query.Encode()

	return parsed.String()
}

// SerialHeader 序列化请求头
func SerialHeader(header http.Header) string {
	if len(header) == 0 {
		return ""
	}

	length := 0
	for k, vs := range header {
		if len(vs) == 1 {
			length += len(k) + len(vs[0]) + 3
		} else if len(vs) == 0 {
			length += len(k) + 3
		} else {
			for _, v := range vs {
				length += len(k) + len(v) + 3
			}
		}
	}

	sb := strings.Builder{}
	sb.Grow(length + 64)
	for k, vs := range header {
		if len(vs) == 1 {
			sb.WriteString(k)
			sb.WriteByte('=')
			sb.WriteString(vs[0])
			sb.Write(_headerDivide)
		} else if len(vs) == 0 {
			sb.WriteString(k)
			sb.WriteByte('=')
			sb.Write(_headerDivide)
		} else {
			for _, v := range vs {
				sb.WriteString(k)
				sb.WriteByte('=')
				sb.WriteString(v)
				sb.Write(_headerDivide)
			}
		}
	}

	return sb.String()
}

// EscapeQuotes 转义请求中的特殊字符
func EscapeQuotes(s string) string {
	return _quoteReplacer.Replace(s)
}

// IsTextType 根据 contentType 判断是否为文本类型
func IsTextType(contentType string) bool {
	if contentType == "" {
		return false
	}

	ct := strings.ToLower(contentType)
	if strings.HasPrefix(ct, "text/") {
		return true
	}

	for _, t := range _textTypes {
		if strings.HasPrefix(ct, t) {
			return true
		}
	}

	return false
}
