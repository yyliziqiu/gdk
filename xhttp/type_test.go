package xhttp

import (
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"
)

func TestResponseError_Status(t *testing.T) {
	e := newResponseError(404, "not found")
	if e.Status() != 404 {
		t.Errorf("got %d, want 404", e.Status())
	}
}

func TestResponseError_Error(t *testing.T) {
	tests := []struct {
		status int
		msg    string
		want   string
	}{
		{200, "ok", "#200 ok"},
		{404, "not found", "#404 not found"},
		{500, "internal error", "#500 internal error"},
		{0, "", "#0 "},
	}
	for _, tt := range tests {
		e := newResponseError(tt.status, tt.msg)
		if e.Error() != tt.want {
			t.Errorf("newResponseError(%d, %q).Error() = %q, want %q", tt.status, tt.msg, e.Error(), tt.want)
		}
	}
}

func TestDefaultClient(t *testing.T) {
	c := DefaultClient()
	if c == nil {
		t.Fatal("DefaultClient returned nil")
	}
	if c.Timeout != 5*time.Second {
		t.Errorf("expected 5s timeout, got %v", c.Timeout)
	}
}

func TestJoinUrl(t *testing.T) {
	tests := []struct {
		prefix string
		suffix string
		want   string
	}{
		{"http://a.com", "", "http://a.com"},
		{"http://a.com/", "", "http://a.com/"},
		// neither has /
		{"http://a.com", "path", "http://a.com/path"},
		// only prefix has /
		{"http://a.com/", "path", "http://a.com/path"},
		// only suffix has /
		{"http://a.com", "/path", "http://a.com/path"},
		// both have /
		{"http://a.com/", "/path", "http://a.com/path"},
		// suffix = "/" (len < 2), both have /
		{"http://a.com/", "/", "http://a.com/"},
		// nested paths
		{"http://a.com/api", "v1/users", "http://a.com/api/v1/users"},
		{"http://a.com/api/", "/v1/users", "http://a.com/api/v1/users"},
	}
	for _, tt := range tests {
		got := JoinUrl(tt.prefix, tt.suffix)
		if got != tt.want {
			t.Errorf("JoinUrl(%q, %q) = %q, want %q", tt.prefix, tt.suffix, got, tt.want)
		}
	}
}

func TestJoinUrls(t *testing.T) {
	tests := []struct {
		segments []string
		want     string
	}{
		{nil, ""},
		{[]string{}, ""},
		{[]string{"http://a.com"}, "http://a.com"},
		{[]string{"http://a.com", ""}, "http://a.com"},
		{[]string{"http://a.com", "api", "v1"}, "http://a.com/api/v1"},
		{[]string{"http://a.com/", "/api/", "/v1"}, "http://a.com/api/v1"},
		{[]string{"http://a.com/", "api", "v1"}, "http://a.com/api/v1"},
		{[]string{"http://a.com", "/api", "/v1"}, "http://a.com/api/v1"},
		// empty segment in middle is skipped
		{[]string{"http://a.com", "", "api"}, "http://a.com/api"},
		// single slash segments with both ends having /
		{[]string{"http://a.com/", "/"}, "http://a.com/"},
	}
	for _, tt := range tests {
		got := JoinUrls(tt.segments...)
		if got != tt.want {
			t.Errorf("JoinUrls(%v) = %q, want %q", tt.segments, got, tt.want)
		}
	}
}

func TestAppendQuery(t *testing.T) {
	t.Run("valid url no existing query", func(t *testing.T) {
		q := url.Values{"k": {"v"}}
		got := AppendQuery("http://a.com/path", q)
		if !strings.Contains(got, "k=v") {
			t.Errorf("expected k=v in %q", got)
		}
		if strings.Contains(got, "&&") {
			t.Errorf("unexpected && in %q", got)
		}
	})

	t.Run("valid url with existing query", func(t *testing.T) {
		q := url.Values{"k": {"v"}}
		got := AppendQuery("http://a.com/path?x=1", q)
		if !strings.Contains(got, "k=v") {
			t.Errorf("expected k=v in %q", got)
		}
		if !strings.Contains(got, "x=1") {
			t.Errorf("expected x=1 in %q", got)
		}
	})

	t.Run("multiple values", func(t *testing.T) {
		q := url.Values{"a": {"1", "2"}}
		got := AppendQuery("http://a.com", q)
		if !strings.Contains(got, "a=1") || !strings.Contains(got, "a=2") {
			t.Errorf("expected a=1 and a=2 in %q", got)
		}
	})

	t.Run("merges existing and new query params", func(t *testing.T) {
		q := url.Values{"new": {"val"}}
		got := AppendQuery("http://a.com/?old=1", q)
		if !strings.Contains(got, "old=1") || !strings.Contains(got, "new=val") {
			t.Errorf("expected old=1 and new=val in %q", got)
		}
	})
}

func TestSerialHeader(t *testing.T) {
	t.Run("empty header returns empty string", func(t *testing.T) {
		got := SerialHeader(http.Header{}, nil)
		if got != "" {
			t.Errorf("expected empty, got %q", got)
		}
	})

	t.Run("nil header returns empty string", func(t *testing.T) {
		got := SerialHeader(nil, nil)
		if got != "" {
			t.Errorf("expected empty, got %q", got)
		}
	})

	t.Run("single value header", func(t *testing.T) {
		h := http.Header{"Content-Type": {"application/json"}}
		got := SerialHeader(h, nil)
		if !strings.Contains(got, "Content-Type=application/json") {
			t.Errorf("expected Content-Type=application/json in %q", got)
		}
	})

	t.Run("multi-value header", func(t *testing.T) {
		h := http.Header{"Accept": {"text/html", "application/json"}}
		got := SerialHeader(h, nil)
		if !strings.Contains(got, "Accept=text/html") {
			t.Errorf("expected Accept=text/html in %q", got)
		}
		if !strings.Contains(got, "Accept=application/json") {
			t.Errorf("expected Accept=application/json in %q", got)
		}
	})

	t.Run("forbid filters out header", func(t *testing.T) {
		h := http.Header{
			"Authorization": {"Bearer token"},
			"Accept":        {"*/*"},
		}
		got := SerialHeader(h, []string{"Authorization"})
		if strings.Contains(got, "Authorization") {
			t.Errorf("Authorization should be forbidden, got %q", got)
		}
		if !strings.Contains(got, "Accept") {
			t.Errorf("Accept should be present, got %q", got)
		}
	})

	t.Run("empty forbid list includes all headers", func(t *testing.T) {
		h := http.Header{
			"Authorization": {"token"},
			"Accept":        {"*/*"},
		}
		got := SerialHeader(h, []string{})
		if !strings.Contains(got, "Authorization") {
			t.Errorf("Authorization should be present, got %q", got)
		}
	})

	t.Run("header with zero values", func(t *testing.T) {
		h := http.Header{"X-Empty": {}}
		got := SerialHeader(h, nil)
		if !strings.Contains(got, "X-Empty=") {
			t.Errorf("expected X-Empty= in %q", got)
		}
	})
}

func TestEscape(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"hello", "hello"},
		{`say "hi"`, `say \"hi\"`},
		{`a\b`, `a\\b`},
		{`a\"b`, `a\\\"b`},
		{"", ""},
		{`no special chars`, `no special chars`},
	}
	for _, tt := range tests {
		got := Escape(tt.input)
		if got != tt.want {
			t.Errorf("Escape(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestIsTextType(t *testing.T) {
	tests := []struct {
		contentType string
		want        bool
	}{
		{"", false},
		{"text/plain", true},
		{"text/html", true},
		{"text/html; charset=utf-8", true},
		{"TEXT/PLAIN", true},
		{"application/json", true},
		{"application/json; charset=utf-8", true},
		{"application/xml", true},
		{"application/xml; charset=utf-8", true},
		{"image/png", false},
		{"image/jpeg", false},
		{"application/octet-stream", false},
		{"multipart/form-data", false},
		{"application/x-www-form-urlencoded", false},
		{"video/mp4", false},
	}
	for _, tt := range tests {
		got := IsTextType(tt.contentType)
		if got != tt.want {
			t.Errorf("IsTextType(%q) = %v, want %v", tt.contentType, got, tt.want)
		}
	}
}
