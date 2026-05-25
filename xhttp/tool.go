package xhttp

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"time"
)

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
		return "{}"
	}

	m := make(map[string]string, len(header))
	for key := range header {
		m[key] = header.Get(key)
	}

	bs, _ := json.Marshal(m)

	return string(bs)
}

// EscapeQuotes 转义请求中的特殊字符
var _quoteReplacer = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")

func EscapeQuotes(s string) string {
	return _quoteReplacer.Replace(s)
}

// IsTextType 根据 contentType 判断是否为文本类型
var _canLogTypes = []string{
	"application/json",
	"application/xml",
}

func IsTextType(contentType string) bool {
	if contentType == "" {
		return false
	}

	ct := strings.ToLower(contentType)
	if strings.HasPrefix(ct, "text/") {
		return true
	}

	for _, t := range _canLogTypes {
		if strings.HasPrefix(ct, t) {
			return true
		}
	}

	return false
}
