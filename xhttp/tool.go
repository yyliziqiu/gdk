package xhttp

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func DefaultClient() *http.Client {
	return &http.Client{
		Timeout: 5 * time.Second,
	}
}

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

var _quoteReplacer = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")

func EscapeQuotes(s string) string {
	return _quoteReplacer.Replace(s)
}
