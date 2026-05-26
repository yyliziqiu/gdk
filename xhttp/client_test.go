package xhttp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"
)

// ---- shared test types ----

// testResp implements JsonResponse and error.
type testResp struct {
	Code int    `json:"code"`
	Data string `json:"data"`
}

func (r *testResp) Failed() bool  { return r.Code != 0 }
func (r *testResp) Error() string { return fmt.Sprintf("code=%d data=%s", r.Code, r.Data) }

// simpleResp implements JsonResponse but NOT error.
type simpleResp struct {
	OK   bool   `json:"ok"`
	Data string `json:"data"`
}

func (r *simpleResp) Failed() bool { return !r.OK }

// testErrBody is used as the Error() option type (value receiver so it implements error).
type testErrBody struct {
	ErrCode int    `json:"err_code"`
	ErrMsg  string `json:"err_msg"`
}

func (e testErrBody) Error() string {
	return fmt.Sprintf("err_code=%d err_msg=%s", e.ErrCode, e.ErrMsg)
}

// ---- helpers ----

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeText(w http.ResponseWriter, status int, body string) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(status)
	w.Write([]byte(body))
}

// ================================================================
// New / options
// ================================================================

func TestNew_defaults(t *testing.T) {
	c := New()
	if c.client == nil {
		t.Error("expected non-nil http.Client")
	}
	if c.format != FormatJson {
		t.Errorf("expected format json, got %s", c.format)
	}
	if !c.logAlways {
		t.Error("expected logAlways=true")
	}
	if c.logHeader != LogHeaderNone {
		t.Errorf("expected logHeader=0, got %d", c.logHeader)
	}
	if c.logEscape {
		t.Error("expected logEscape=false")
	}
}

func TestNew_withOptions(t *testing.T) {
	c := New(
		Prefix("http://example.com"),
		Format(FormatText),
		LogAlways(false),
		LogLength(512, 512),
		LogHeader(LogHeaderBoth),
		LogForbid([]string{"Authorization"}),
		LogEscape(true),
		Dumps(true),
	)
	if c.prefix != "http://example.com" {
		t.Errorf("unexpected prefix: %s", c.prefix)
	}
	if c.format != FormatText {
		t.Errorf("unexpected format: %s", c.format)
	}
	if c.logAlways {
		t.Error("expected logAlways=false")
	}
	if c.logLength != [2]int{512, 512} {
		t.Errorf("unexpected logLength: %v", c.logLength)
	}
	if c.logHeader != LogHeaderBoth {
		t.Errorf("unexpected logHeader: %d", c.logHeader)
	}
	if len(c.logForbid) != 1 || c.logForbid[0] != "Authorization" {
		t.Errorf("unexpected logForbid: %v", c.logForbid)
	}
	if !c.logEscape {
		t.Error("expected logEscape=true")
	}
	if !c.dumps {
		t.Error("expected dumps=true")
	}
}

func TestOption_Timeout(t *testing.T) {
	c := New(Timeout(10 * time.Second))
	if c.client.Timeout != 10*time.Second {
		t.Errorf("expected 10s, got %v", c.client.Timeout)
	}
}

func TestOption_BaseUrl(t *testing.T) {
	c := New(BaseUrl("http://base.com"))
	if c.prefix != "http://base.com" {
		t.Errorf("expected http://base.com, got %s", c.prefix)
	}
}

func TestOption_Cookie(t *testing.T) {
	c := New(Cookie(nil))
	if c.client.Jar == nil {
		t.Error("expected non-nil cookie jar")
	}
}

func TestOption_WithClient(t *testing.T) {
	custom := &http.Client{Timeout: 30 * time.Second}
	c := New(WithClient(custom))
	if c.client != custom {
		t.Error("expected the custom http.Client to be set")
	}
}

func TestOption_Error(t *testing.T) {
	c := New(Error(testErrBody{}))
	if c.errtyp == nil {
		t.Error("expected errtyp to be set")
	}
}

func TestOption_LogChange(t *testing.T) {
	r := strings.NewReplacer("a", "b")
	c := New(LogChange(r))
	if c.logChange == nil {
		t.Error("expected logChange to be set")
	}
}

// ================================================================
// Get
// ================================================================

func TestGet_success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		writeJSON(w, 200, testResp{Code: 0, Data: "ok"})
	}))
	defer srv.Close()

	c := New()
	var out testResp
	if err := c.Get(srv.URL, nil, nil, &out); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Data != "ok" {
		t.Errorf("unexpected data: %s", out.Data)
	}
}

func TestGet_withQuery(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("foo") != "bar" {
			t.Errorf("expected foo=bar, got %s", r.URL.RawQuery)
		}
		writeJSON(w, 200, testResp{})
	}))
	defer srv.Close()

	c := New()
	if err := c.Get(srv.URL, url.Values{"foo": {"bar"}}, nil, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGet_withHeader(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Custom") != "value" {
			t.Errorf("expected X-Custom=value, got %s", r.Header.Get("X-Custom"))
		}
		writeJSON(w, 200, testResp{})
	}))
	defer srv.Close()

	c := New()
	if err := c.Get(srv.URL, nil, http.Header{"X-Custom": {"value"}}, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGet_withPrefix(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/users" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		writeJSON(w, 200, testResp{})
	}))
	defer srv.Close()

	c := New(Prefix(srv.URL + "/api"))
	if err := c.Get("/users", nil, nil, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGet_absoluteURLIgnoresPrefix(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, 200, testResp{Data: "direct"})
	}))
	defer srv.Close()

	c := New(Prefix("http://other.com/api"))
	var out testResp
	if err := c.Get(srv.URL, nil, nil, &out); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Data != "direct" {
		t.Errorf("unexpected data: %s", out.Data)
	}
}

func TestGet_4xx(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, 404, map[string]string{"message": "not found"})
	}))
	defer srv.Close()

	c := New()
	err := c.Get(srv.URL, nil, nil, nil)
	if err == nil {
		t.Fatal("expected error for 404 response")
	}
	re, ok := err.(*ResponseError)
	if !ok {
		t.Fatalf("expected *ResponseError, got %T: %v", err, err)
	}
	if re.Status() != 404 {
		t.Errorf("expected status 404, got %d", re.Status())
	}
}

func TestGet_5xx(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, 500, map[string]string{"error": "server error"})
	}))
	defer srv.Close()

	c := New()
	err := c.Get(srv.URL, nil, nil, nil)
	if err == nil {
		t.Fatal("expected error for 500 response")
	}
}

func TestGet_jsonResponseFailed_withError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, 200, testResp{Code: 1, Data: "fail"})
	}))
	defer srv.Close()

	c := New()
	var out testResp
	err := c.Get(srv.URL, nil, nil, &out)
	if err == nil {
		t.Fatal("expected error from Failed()")
	}
	if err.Error() != "code=1 data=fail" {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestGet_jsonResponseFailed_withoutError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, 200, simpleResp{OK: false, Data: "fail"})
	}))
	defer srv.Close()

	c := New()
	var out simpleResp
	err := c.Get(srv.URL, nil, nil, &out)
	if err == nil {
		t.Fatal("expected error from Failed()")
	}
	// Should be *ResponseError since simpleResp does not implement error
	if _, ok := err.(*ResponseError); !ok {
		t.Errorf("expected *ResponseError, got %T", err)
	}
}

func TestGet_withErrType(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, 500, testErrBody{ErrCode: 42, ErrMsg: "custom error"})
	}))
	defer srv.Close()

	c := New(Error(testErrBody{}))
	err := c.Get(srv.URL, nil, nil, nil)
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "42") {
		t.Errorf("expected error code 42 in error string, got: %v", err)
	}
}

func TestGet_outNil(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, 200, testResp{Data: "ignored"})
	}))
	defer srv.Close()

	c := New()
	if err := c.Get(srv.URL, nil, nil, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// ================================================================
// Post
// ================================================================

func TestPost_success(t *testing.T) {
	type reqBody struct {
		Name string `json:"name"`
	}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("unexpected content-type: %s", r.Header.Get("Content-Type"))
		}
		var body reqBody
		json.NewDecoder(r.Body).Decode(&body)
		if body.Name != "test" {
			t.Errorf("expected name=test, got %s", body.Name)
		}
		writeJSON(w, 200, testResp{Data: "posted"})
	}))
	defer srv.Close()

	c := New()
	var out testResp
	if err := c.Post(srv.URL, nil, reqBody{Name: "test"}, &out); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Data != "posted" {
		t.Errorf("unexpected data: %s", out.Data)
	}
}

func TestPost_nilBody(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, 200, testResp{})
	}))
	defer srv.Close()

	c := New()
	if err := c.Post(srv.URL, nil, nil, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestPost_errorResponse(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, 400, map[string]string{"msg": "bad request"})
	}))
	defer srv.Close()

	c := New()
	err := c.Post(srv.URL, nil, map[string]string{"x": "1"}, nil)
	if err == nil {
		t.Fatal("expected error for 400 response")
	}
}

// ================================================================
// Put / Patch / Delete
// ================================================================

func TestPut_success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		writeJSON(w, 200, testResp{Data: "updated"})
	}))
	defer srv.Close()

	c := New()
	var out testResp
	if err := c.Put(srv.URL, nil, map[string]string{"key": "val"}, &out); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Data != "updated" {
		t.Errorf("unexpected data: %s", out.Data)
	}
}

func TestPatch_success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("expected PATCH, got %s", r.Method)
		}
		writeJSON(w, 200, testResp{})
	}))
	defer srv.Close()

	c := New()
	if err := c.Patch(srv.URL, nil, map[string]string{"x": "1"}, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDelete_success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Query().Get("id") != "1" {
			t.Errorf("expected id=1, got %s", r.URL.RawQuery)
		}
		writeJSON(w, 200, testResp{})
	}))
	defer srv.Close()

	c := New()
	if err := c.Delete(srv.URL, url.Values{"id": {"1"}}, nil, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// ================================================================
// Format: text
// ================================================================

func TestGet_textFormat_success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeText(w, 200, "hello text")
	}))
	defer srv.Close()

	c := New(Format(FormatText))
	var out []byte
	if err := c.Get(srv.URL, nil, nil, &out); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(out) != "hello text" {
		t.Errorf("unexpected text: %s", out)
	}
}

func TestGet_textFormat_error(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeText(w, 404, "not found")
	}))
	defer srv.Close()

	c := New(Format(FormatText))
	err := c.Get(srv.URL, nil, nil, nil)
	if err == nil {
		t.Fatal("expected error for 404 text response")
	}
	re, ok := err.(*ResponseError)
	if !ok {
		t.Fatalf("expected *ResponseError, got %T", err)
	}
	if re.Status() != 404 {
		t.Errorf("expected status 404, got %d", re.Status())
	}
}

func TestGet_textFormat_wrongReceiverType(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeText(w, 200, "hello")
	}))
	defer srv.Close()

	c := New(Format(FormatText))
	var out string // not *[]byte
	err := c.Get(srv.URL, nil, nil, &out)
	if err == nil {
		t.Fatal("expected error when out is not *[]byte")
	}
	if !strings.Contains(err.Error(), "*[]byte") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestGet_textFormat_nilOut(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeText(w, 200, "ignored")
	}))
	defer srv.Close()

	c := New(Format(FormatText))
	if err := c.Get(srv.URL, nil, nil, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// ================================================================
// GetBinary
// ================================================================

func TestGetBinary_success(t *testing.T) {
	data := []byte{0x01, 0x02, 0x03}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/png")
		w.Write(data)
	}))
	defer srv.Close()

	c := New()
	body, ct, err := c.GetBinary(srv.URL, nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ct != "image/png" {
		t.Errorf("expected image/png, got %s", ct)
	}
	if !bytes.Equal(body, data) {
		t.Errorf("body mismatch: got %v, want %v", body, data)
	}
}

func TestGetBinary_withQuery(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("size") != "large" {
			t.Errorf("expected size=large, got %s", r.URL.RawQuery)
		}
		w.Header().Set("Content-Type", "image/jpeg")
		w.Write([]byte("fake-image"))
	}))
	defer srv.Close()

	c := New()
	_, _, err := c.GetBinary(srv.URL, url.Values{"size": {"large"}}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// ================================================================
// PostJson
// ================================================================

func TestPostJson_success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("expected application/json, got %s", r.Header.Get("Content-Type"))
		}
		body, _ := io.ReadAll(r.Body)
		var m map[string]string
		json.Unmarshal(body, &m)
		if m["hello"] != "world" {
			t.Errorf("unexpected body: %s", body)
		}
		writeJSON(w, 200, testResp{Data: "json-posted"})
	}))
	defer srv.Close()

	c := New()
	in := []byte(`{"hello":"world"}`)
	var out testResp
	if err := c.PostJson(srv.URL, nil, in, &out); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Data != "json-posted" {
		t.Errorf("unexpected data: %s", out.Data)
	}
}

// ================================================================
// PostForm
// ================================================================

func TestPostForm_success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ct := r.Header.Get("Content-Type")
		if ct != "application/x-www-form-urlencoded" {
			t.Errorf("unexpected content-type: %s", ct)
		}
		r.ParseForm()
		if r.FormValue("name") != "test" {
			t.Errorf("expected name=test, got %s", r.FormValue("name"))
		}
		if r.FormValue("age") != "30" {
			t.Errorf("expected age=30, got %s", r.FormValue("age"))
		}
		writeJSON(w, 200, testResp{Data: "form-posted"})
	}))
	defer srv.Close()

	c := New()
	var out testResp
	if err := c.PostForm(srv.URL, nil, url.Values{"name": {"test"}, "age": {"30"}}, &out); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Data != "form-posted" {
		t.Errorf("unexpected data: %s", out.Data)
	}
}

// ================================================================
// PostData
// ================================================================

func TestPostData_withValues(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.Header.Get("Content-Type"), "multipart/form-data") {
			t.Errorf("expected multipart/form-data, got %s", r.Header.Get("Content-Type"))
		}
		r.ParseMultipartForm(1 << 20)
		if r.FormValue("key") != "val" {
			t.Errorf("expected key=val, got %s", r.FormValue("key"))
		}
		writeJSON(w, 200, testResp{Data: "data-posted"})
	}))
	defer srv.Close()

	c := New()
	var out testResp
	if err := c.PostData(srv.URL, nil, map[string]string{"key": "val"}, nil, &out); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Data != "data-posted" {
		t.Errorf("unexpected data: %s", out.Data)
	}
}

func TestPostData_emptyValuesAndFiles(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, 200, testResp{Data: "empty"})
	}))
	defer srv.Close()

	c := New()
	if err := c.PostData(srv.URL, nil, nil, nil, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// ================================================================
// PostFile
// ================================================================

func TestPostFile_success(t *testing.T) {
	f, err := os.CreateTemp("", "xhttp-test-*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())
	f.WriteString("file content")
	f.Close()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseMultipartForm(1 << 20)
		file, _, fileErr := r.FormFile("upload")
		if fileErr != nil {
			t.Errorf("expected file upload: %v", fileErr)
		} else {
			defer file.Close()
			got, _ := io.ReadAll(file)
			if string(got) != "file content" {
				t.Errorf("unexpected file content: %s", got)
			}
		}
		if r.FormValue("field1") != "v1" {
			t.Errorf("expected field1=v1, got %s", r.FormValue("field1"))
		}
		writeJSON(w, 200, testResp{Data: "file-uploaded"})
	}))
	defer srv.Close()

	c := New()
	var out testResp
	if err := c.PostFile(srv.URL, nil, map[string]string{"field1": "v1"}, "upload", f.Name(), &out); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Data != "file-uploaded" {
		t.Errorf("unexpected data: %s", out.Data)
	}
}

// ================================================================
// PostBinary
// ================================================================

func TestPostBinary_success(t *testing.T) {
	data := []byte("binary content")
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/octet-stream" {
			t.Errorf("unexpected content-type: %s", r.Header.Get("Content-Type"))
		}
		body, _ := io.ReadAll(r.Body)
		if !bytes.Equal(body, data) {
			t.Errorf("body mismatch: got %v, want %v", body, data)
		}
		writeJSON(w, 200, testResp{Data: "binary-posted"})
	}))
	defer srv.Close()

	c := New()
	var out testResp
	if err := c.PostBinary(srv.URL, nil, "application/octet-stream", bytes.NewReader(data), &out); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Data != "binary-posted" {
		t.Errorf("unexpected data: %s", out.Data)
	}
}

// ================================================================
// PostStream
// ================================================================

func TestPostStream_success(t *testing.T) {
	content := []byte("stream content")
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.Header.Get("Content-Type"), "multipart/form-data") {
			t.Errorf("unexpected content-type: %s", r.Header.Get("Content-Type"))
		}
		r.ParseMultipartForm(1 << 20)
		f, _, err := r.FormFile("file")
		if err != nil {
			t.Errorf("expected file in form: %v", err)
			writeJSON(w, 500, testResp{Code: 1})
			return
		}
		defer f.Close()
		got, _ := io.ReadAll(f)
		if !bytes.Equal(got, content) {
			t.Errorf("unexpected file content: %s", got)
		}
		writeJSON(w, 200, testResp{Data: "streamed"})
	}))
	defer srv.Close()

	c := New()
	var out testResp
	if err := c.PostStream(srv.URL, nil, nil, "file", "test.txt", "text/plain", bytes.NewReader(content), &out); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Data != "streamed" {
		t.Errorf("unexpected data: %s", out.Data)
	}
}

func TestPostStream_withValues(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseMultipartForm(1 << 20)
		if r.FormValue("category") != "docs" {
			t.Errorf("expected category=docs, got %s", r.FormValue("category"))
		}
		writeJSON(w, 200, testResp{})
	}))
	defer srv.Close()

	c := New()
	vals := map[string]string{"category": "docs"}
	if err := c.PostStream(srv.URL, nil, vals, "file", "doc.pdf", "", bytes.NewReader([]byte("data")), nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestPostStream_nilStream(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, 200, testResp{})
	}))
	defer srv.Close()

	c := New()
	if err := c.PostStream(srv.URL, nil, nil, "file", "empty.txt", "text/plain", nil, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// ================================================================
// ForwardBinary
// ================================================================

func TestForwardBinary(t *testing.T) {
	imgData := []byte{0xFF, 0xD8, 0xFF}

	src := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/jpeg")
		w.Write(imgData)
	}))
	defer src.Close()

	var received []byte
	dst := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		received, _ = io.ReadAll(r.Body)
		writeJSON(w, 200, testResp{Data: "forwarded"})
	}))
	defer dst.Close()

	c := New()
	var out testResp
	if err := c.ForwardBinary(dst.URL, nil, src.URL, &out); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !bytes.Equal(received, imgData) {
		t.Errorf("data mismatch: got %v, want %v", received, imgData)
	}
	if out.Data != "forwarded" {
		t.Errorf("unexpected data: %s", out.Data)
	}
}

// ================================================================
// ForwardStream
// ================================================================

func TestForwardStream(t *testing.T) {
	fileData := []byte("stream data")

	src := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.Write(fileData)
	}))
	defer src.Close()

	dst := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseMultipartForm(1 << 20)
		f, _, err := r.FormFile("data")
		if err != nil {
			t.Errorf("form file error: %v", err)
		} else {
			defer f.Close()
			got, _ := io.ReadAll(f)
			if !bytes.Equal(got, fileData) {
				t.Errorf("data mismatch: got %s, want %s", got, fileData)
			}
		}
		writeJSON(w, 200, testResp{Data: "stream-forwarded"})
	}))
	defer dst.Close()

	c := New()
	if err := c.ForwardStream(dst.URL, nil, nil, "data", "text/plain", src.URL+"/file.txt", nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// ================================================================
// Get2 / Post2 / Put2 / Patch2 / Delete2
// ================================================================

func TestGet2_success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		writeJSON(w, 200, testResp{Data: "get2"})
	}))
	defer srv.Close()

	c := New()
	var out testResp
	body, err := c.Get2(srv.URL, nil, &out)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Data != "get2" {
		t.Errorf("unexpected data: %s", out.Data)
	}
	if len(body) == 0 {
		t.Error("expected non-empty body bytes")
	}
}

func TestPost2_success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		writeJSON(w, 200, testResp{Data: "post2"})
	}))
	defer srv.Close()

	c := New()
	var out testResp
	body, err := c.Post2(srv.URL, nil, map[string]string{"x": "1"}, &out)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Data != "post2" {
		t.Errorf("unexpected data: %s", out.Data)
	}
	if len(body) == 0 {
		t.Error("expected non-empty body bytes")
	}
}

func TestPut2_success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		writeJSON(w, 200, testResp{Data: "put2"})
	}))
	defer srv.Close()

	c := New()
	body, err := c.Put2(srv.URL, nil, nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(body) == 0 {
		t.Error("expected non-empty body bytes")
	}
}

func TestPatch2_success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("expected PATCH, got %s", r.Method)
		}
		writeJSON(w, 200, testResp{})
	}))
	defer srv.Close()

	c := New()
	_, err := c.Patch2(srv.URL, nil, nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDelete2_success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		writeJSON(w, 200, testResp{})
	}))
	defer srv.Close()

	c := New()
	_, err := c.Delete2(srv.URL, nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// ================================================================
// Post3 / Put3 / Patch3
// ================================================================

func TestPost3_success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("unexpected content-type: %s", r.Header.Get("Content-Type"))
		}
		body, _ := io.ReadAll(r.Body)
		var m map[string]string
		json.Unmarshal(body, &m)
		if m["raw"] != "json" {
			t.Errorf("unexpected body: %s", body)
		}
		writeJSON(w, 200, testResp{Data: "post3"})
	}))
	defer srv.Close()

	c := New()
	var out testResp
	body, err := c.Post3(srv.URL, nil, []byte(`{"raw":"json"}`), &out)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Data != "post3" {
		t.Errorf("unexpected data: %s", out.Data)
	}
	if len(body) == 0 {
		t.Error("expected non-empty body bytes")
	}
}

func TestPut3_success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		writeJSON(w, 200, testResp{})
	}))
	defer srv.Close()

	c := New()
	_, err := c.Put3(srv.URL, nil, []byte(`{}`), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestPatch3_success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("expected PATCH, got %s", r.Method)
		}
		writeJSON(w, 200, testResp{})
	}))
	defer srv.Close()

	c := New()
	_, err := c.Patch3(srv.URL, nil, []byte(`{}`), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// ================================================================
// Options: redirect / auth / callbacks
// ================================================================

func TestDisableRedirect(t *testing.T) {
	target := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, 200, testResp{Data: "target"})
	}))
	defer target.Close()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, target.URL, http.StatusFound)
	}))
	defer srv.Close()

	c := New(DisableRedirect())
	err := c.Get(srv.URL, nil, nil, nil)
	// DisableRedirect returns the 302 response which triggers "this is a redirect response"
	if err == nil {
		t.Fatal("expected error from 302 response with DisableRedirect")
	}
	if !strings.Contains(err.Error(), "redirect") {
		t.Errorf("expected redirect error, got: %v", err)
	}
}

func TestLimitRedirect(t *testing.T) {
	callCount := 0
	var srvURL string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		http.Redirect(w, r, srvURL, http.StatusFound)
	}))
	defer srv.Close()
	srvURL = srv.URL

	c := New(LimitRedirect(2))
	err := c.Get(srv.URL, nil, nil, nil)
	if err == nil {
		t.Fatal("expected error after redirect limit")
	}
}

func TestBasicAuth(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()
		if !ok || user != "admin" || pass != "secret" {
			t.Errorf("expected basic auth admin:secret, got user=%q ok=%v", user, ok)
		}
		writeJSON(w, 200, testResp{})
	}))
	defer srv.Close()

	c := New(BasicAuth("admin", "secret"))
	if err := c.Get(srv.URL, nil, nil, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestBearerToken(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth != "Bearer mytoken" {
			t.Errorf("expected 'Bearer mytoken', got %s", auth)
		}
		writeJSON(w, 200, testResp{})
	}))
	defer srv.Close()

	c := New(BearerToken("mytoken"))
	if err := c.Get(srv.URL, nil, nil, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRequestBefore(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Before") != "yes" {
			t.Errorf("expected X-Before=yes, got %s", r.Header.Get("X-Before"))
		}
		writeJSON(w, 200, testResp{})
	}))
	defer srv.Close()

	c := New(RequestBefore(func(req *http.Request) {
		req.Header.Set("X-Before", "yes")
	}))
	if err := c.Get(srv.URL, nil, nil, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestResponseAfter_called(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, 200, testResp{Data: "ok"})
	}))
	defer srv.Close()

	called := false
	c := New(ResponseAfter(func(res *http.Response) error {
		called = true
		return nil
	}))
	if err := c.Get(srv.URL, nil, nil, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Error("responseAfter hook was not called")
	}
}

func TestResponseAfter_returnsError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, 200, testResp{})
	}))
	defer srv.Close()

	c := New(ResponseAfter(func(res *http.Response) error {
		return fmt.Errorf("after-hook error")
	}))
	err := c.Get(srv.URL, nil, nil, nil)
	if err == nil {
		t.Fatal("expected error from responseAfter")
	}
	if err.Error() != "after-hook error" {
		t.Errorf("unexpected error: %v", err)
	}
}

// ================================================================
// Logging options (smoke tests – ensure no panics)
// ================================================================

func TestLogHeader_logHeaderReq(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, 200, testResp{})
	}))
	defer srv.Close()

	c := New(LogHeader(LogHeaderReq))
	if err := c.Get(srv.URL, nil, http.Header{"X-Test": {"1"}}, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestLogHeader_logHeaderRes(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Response", "1")
		writeJSON(w, 200, testResp{})
	}))
	defer srv.Close()

	c := New(LogHeader(LogHeaderRes))
	if err := c.Get(srv.URL, nil, nil, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestLogHeader_logHeaderBoth(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Response", "1")
		writeJSON(w, 200, testResp{})
	}))
	defer srv.Close()

	c := New(LogHeader(LogHeaderBoth), LogForbid([]string{"User-Agent"}))
	if err := c.Get(srv.URL, nil, http.Header{"X-Test": {"1"}}, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestLogEscape_withSpecialChars(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, 200, testResp{Data: "line1\nline2"})
	}))
	defer srv.Close()

	c := New(LogEscape(true))
	var out testResp
	if err := c.Get(srv.URL, nil, nil, &out); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestLogAlways_false(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Write([]byte{0x01, 0x02})
	}))
	defer srv.Close()

	c := New(LogAlways(false))
	_, _, err := c.GetBinary(srv.URL, nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
