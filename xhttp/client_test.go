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
)

// ----- test helpers -----

type apiResp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func (r *apiResp) Failed() bool  { return r.Code != 0 }
func (r *apiResp) Error() string { return fmt.Sprintf("[%d] %s", r.Code, r.Msg) }

type apiErrResp struct {
	ErrCode int    `json:"err_code"`
	ErrMsg  string `json:"err_msg"`
}

func (e apiErrResp) Error() string { return fmt.Sprintf("[%d] %s", e.ErrCode, e.ErrMsg) }

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}

// ----- New / options -----

func TestNew_Defaults(t *testing.T) {
	c := New()
	if c.client == nil {
		t.Fatal("client is nil")
	}
	if c.format != FormatJson {
		t.Fatalf("format: want %q, got %q", FormatJson, c.format)
	}
	if c.logLength != 2048 {
		t.Fatalf("logLength: want 2048, got %d", c.logLength)
	}
	if c.logChange == nil {
		t.Fatal("logChange is nil")
	}
}

func TestNew_WithPrefix(t *testing.T) {
	c := New(Prefix("http://example.com/api"))
	if c.prefix != "http://example.com/api" {
		t.Fatalf("prefix: want %q, got %q", "http://example.com/api", c.prefix)
	}
}

func TestNew_WithFormat(t *testing.T) {
	c := New(Format(FormatText))
	if c.format != FormatText {
		t.Fatalf("format: want %q, got %q", FormatText, c.format)
	}
}

func TestNew_WithLogOptions(t *testing.T) {
	c := New(LogLength(512), LogHeader(true), LogEscape(true))
	if c.logLength != 512 {
		t.Fatalf("logLength: want 512, got %d", c.logLength)
	}
	if !c.logHeader {
		t.Fatal("logHeader should be true")
	}
	if !c.logEscape {
		t.Fatal("logEscape should be true")
	}
}

// ----- GET -----

func TestGet_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method: want GET, got %s", r.Method)
		}
		writeJSON(w, apiResp{Code: 0, Msg: "ok"})
	}))
	defer srv.Close()

	c := New()
	var out apiResp
	if err := c.Get(srv.URL, nil, nil, &out); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Msg != "ok" {
		t.Fatalf("msg: want %q, got %q", "ok", out.Msg)
	}
}

func TestGet_WithQuery(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Get("key"); got != "val" {
			t.Errorf("query key: want %q, got %q", "val", got)
		}
		writeJSON(w, apiResp{})
	}))
	defer srv.Close()

	c := New()
	if err := c.Get(srv.URL, url.Values{"key": {"val"}}, nil, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGet_WithHeader(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("X-Token"); got != "mytoken" {
			t.Errorf("header X-Token: want %q, got %q", "mytoken", got)
		}
		writeJSON(w, apiResp{})
	}))
	defer srv.Close()

	c := New()
	if err := c.Get(srv.URL, nil, http.Header{"X-Token": {"mytoken"}}, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGet_404(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("not found"))
	}))
	defer srv.Close()

	if err := New().Get(srv.URL, nil, nil, nil); err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGet_500(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("internal error"))
	}))
	defer srv.Close()

	if err := New().Get(srv.URL, nil, nil, nil); err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGet_AbsoluteURL_IgnoresPrefix(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, apiResp{Code: 0, Msg: "direct"})
	}))
	defer srv.Close()

	c := New(Prefix("http://other-host.invalid"))
	var out apiResp
	if err := c.Get(srv.URL, nil, nil, &out); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Msg != "direct" {
		t.Fatalf("msg: want %q, got %q", "direct", out.Msg)
	}
}

// ----- POST / PUT / PATCH / DELETE -----

func TestPost_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method: want POST, got %s", r.Method)
		}
		if ct := r.Header.Get("Content-Type"); ct != "application/json" {
			t.Errorf("content-type: want application/json, got %s", ct)
		}
		var in map[string]any
		json.NewDecoder(r.Body).Decode(&in)
		if in["name"] != "alice" {
			t.Errorf("name: want alice, got %v", in["name"])
		}
		writeJSON(w, apiResp{Code: 0, Msg: "created"})
	}))
	defer srv.Close()

	c := New()
	var out apiResp
	if err := c.Post(srv.URL, nil, map[string]string{"name": "alice"}, &out); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Msg != "created" {
		t.Fatalf("msg: want %q, got %q", "created", out.Msg)
	}
}

func TestPost_NilBody(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		if string(body) != "{}" {
			t.Errorf("body: want {}, got %s", string(body))
		}
		writeJSON(w, apiResp{})
	}))
	defer srv.Close()

	if err := New().Post(srv.URL, nil, nil, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestPut_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("method: want PUT, got %s", r.Method)
		}
		writeJSON(w, apiResp{Code: 0, Msg: "updated"})
	}))
	defer srv.Close()

	var out apiResp
	if err := New().Put(srv.URL, nil, map[string]string{"key": "val"}, &out); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Msg != "updated" {
		t.Fatalf("msg: want %q, got %q", "updated", out.Msg)
	}
}

func TestPatch_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("method: want PATCH, got %s", r.Method)
		}
		writeJSON(w, apiResp{Code: 0, Msg: "patched"})
	}))
	defer srv.Close()

	var out apiResp
	if err := New().Patch(srv.URL, nil, map[string]string{"key": "val"}, &out); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Msg != "patched" {
		t.Fatalf("msg: want %q, got %q", "patched", out.Msg)
	}
}

func TestDelete_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("method: want DELETE, got %s", r.Method)
		}
		writeJSON(w, apiResp{})
	}))
	defer srv.Close()

	if err := New().Delete(srv.URL, nil, nil, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// ----- GetBinary / PostJson / PostForm -----

func TestGetBinary(t *testing.T) {
	payload := []byte{0x89, 0x50, 0x4e, 0x47}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/png")
		w.Write(payload)
	}))
	defer srv.Close()

	data, ct, err := New().GetBinary(srv.URL, nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !bytes.Equal(data, payload) {
		t.Fatal("binary data mismatch")
	}
	if !strings.HasPrefix(ct, "image/png") {
		t.Fatalf("content-type: want image/png, got %s", ct)
	}
}

func TestPostJson(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if ct := r.Header.Get("Content-Type"); ct != "application/json" {
			t.Errorf("content-type: want application/json, got %s", ct)
		}
		body, _ := io.ReadAll(r.Body)
		var in map[string]any
		json.Unmarshal(body, &in)
		if in["name"] != "bob" {
			t.Errorf("name: want bob, got %v", in["name"])
		}
		writeJSON(w, apiResp{Code: 0, Msg: "ok"})
	}))
	defer srv.Close()

	var out apiResp
	if err := New().PostJson(srv.URL, nil, []byte(`{"name":"bob"}`), &out); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestPostForm(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if ct := r.Header.Get("Content-Type"); ct != "application/x-www-form-urlencoded" {
			t.Errorf("content-type: want application/x-www-form-urlencoded, got %s", ct)
		}
		r.ParseForm()
		if got := r.FormValue("username"); got != "alice" {
			t.Errorf("username: want alice, got %s", got)
		}
		writeJSON(w, apiResp{})
	}))
	defer srv.Close()

	if err := New().PostForm(srv.URL, nil, url.Values{"username": {"alice"}}, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// ----- PostData / PostBinary / PostStream -----

func TestPostData_ValuesOnly(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if ct := r.Header.Get("Content-Type"); !strings.HasPrefix(ct, "multipart/form-data") {
			t.Errorf("content-type: want multipart/form-data prefix, got %s", ct)
		}
		r.ParseMultipartForm(1 << 20)
		if got := r.FormValue("field1"); got != "value1" {
			t.Errorf("field1: want value1, got %s", got)
		}
		writeJSON(w, apiResp{})
	}))
	defer srv.Close()

	if err := New().PostData(srv.URL, nil, map[string]string{"field1": "value1"}, nil, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestPostData_WithFile(t *testing.T) {
	f, err := os.CreateTemp("", "xhttp_test_*.txt")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	f.Write([]byte("file content"))
	f.Close()
	defer os.Remove(f.Name())

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseMultipartForm(1 << 20)
		file, _, ferr := r.FormFile("upload")
		if ferr != nil {
			t.Errorf("form file: %v", ferr)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		data, _ := io.ReadAll(file)
		if string(data) != "file content" {
			t.Errorf("file content: want %q, got %q", "file content", string(data))
		}
		writeJSON(w, apiResp{})
	}))
	defer srv.Close()

	var out apiResp
	if err := New().PostData(srv.URL, nil, nil, map[string]string{"upload": f.Name()}, &out); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestPostBinary(t *testing.T) {
	data := []byte("binary data")
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if ct := r.Header.Get("Content-Type"); ct != "application/octet-stream" {
			t.Errorf("content-type: want application/octet-stream, got %s", ct)
		}
		body, _ := io.ReadAll(r.Body)
		if !bytes.Equal(body, data) {
			t.Error("body mismatch")
		}
		writeJSON(w, apiResp{})
	}))
	defer srv.Close()

	if err := New().PostBinary(srv.URL, nil, "application/octet-stream", bytes.NewReader(data), nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestPostStream(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseMultipartForm(1 << 20)
		if got := r.FormValue("desc"); got != "test" {
			t.Errorf("desc: want test, got %s", got)
		}
		file, _, ferr := r.FormFile("stream")
		if ferr != nil {
			t.Errorf("form file: %v", ferr)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		data, _ := io.ReadAll(file)
		if string(data) != "stream content" {
			t.Errorf("stream: want %q, got %q", "stream content", string(data))
		}
		writeJSON(w, apiResp{})
	}))
	defer srv.Close()

	values := map[string]string{"desc": "test"}
	stream := strings.NewReader("stream content")
	if err := New().PostStream(srv.URL, nil, values, "stream", "data.txt", "text/plain", stream, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// ----- Get2 / Post2 / Put2 / Patch2 / Delete2 -----

func TestGet2(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, apiResp{Code: 0, Msg: "ok"})
	}))
	defer srv.Close()

	var out apiResp
	body, err := New().Get2(srv.URL, nil, &out)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(body) == 0 {
		t.Fatal("body should not be empty")
	}
	if out.Msg != "ok" {
		t.Fatalf("msg: want ok, got %s", out.Msg)
	}
}

func TestPost2(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method: want POST, got %s", r.Method)
		}
		writeJSON(w, apiResp{})
	}))
	defer srv.Close()

	body, err := New().Post2(srv.URL, nil, map[string]string{"k": "v"}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(body) == 0 {
		t.Fatal("body should not be empty")
	}
}

func TestPut2(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("method: want PUT, got %s", r.Method)
		}
		writeJSON(w, apiResp{})
	}))
	defer srv.Close()

	body, err := New().Put2(srv.URL, nil, nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(body) == 0 {
		t.Fatal("body should not be empty")
	}
}

func TestPatch2(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("method: want PATCH, got %s", r.Method)
		}
		writeJSON(w, apiResp{})
	}))
	defer srv.Close()

	body, err := New().Patch2(srv.URL, nil, nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(body) == 0 {
		t.Fatal("body should not be empty")
	}
}

func TestDelete2(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("method: want DELETE, got %s", r.Method)
		}
		writeJSON(w, apiResp{})
	}))
	defer srv.Close()

	body, err := New().Delete2(srv.URL, nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(body) == 0 {
		t.Fatal("body should not be empty")
	}
}

// ----- Post3 / Put3 / Patch3 -----

func TestPost3(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var in map[string]any
		json.Unmarshal(body, &in)
		if in["x"] != "y" {
			t.Errorf("x: want y, got %v", in["x"])
		}
		writeJSON(w, apiResp{Code: 0, Msg: "ok"})
	}))
	defer srv.Close()

	var out apiResp
	body, err := New().Post3(srv.URL, nil, []byte(`{"x":"y"}`), &out)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(body) == 0 {
		t.Fatal("body should not be empty")
	}
}

func TestPut3(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("method: want PUT, got %s", r.Method)
		}
		writeJSON(w, apiResp{})
	}))
	defer srv.Close()

	body, err := New().Put3(srv.URL, nil, []byte(`{}`), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(body) == 0 {
		t.Fatal("body should not be empty")
	}
}

func TestPatch3(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("method: want PATCH, got %s", r.Method)
		}
		writeJSON(w, apiResp{})
	}))
	defer srv.Close()

	body, err := New().Patch3(srv.URL, nil, []byte(`{}`), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(body) == 0 {
		t.Fatal("body should not be empty")
	}
}

// ----- auth / hook options -----

func TestPrefix_PrependedToPath(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/users" {
			t.Errorf("path: want /api/users, got %s", r.URL.Path)
		}
		writeJSON(w, apiResp{})
	}))
	defer srv.Close()

	c := New(Prefix(srv.URL + "/api"))
	if err := c.Get("/users", nil, nil, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestBasicAuth(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()
		if !ok || user != "admin" || pass != "secret" {
			t.Errorf("basic auth: want admin/secret, got %s/%s (ok=%v)", user, pass, ok)
		}
		writeJSON(w, apiResp{})
	}))
	defer srv.Close()

	if err := New(BasicAuth("admin", "secret")).Get(srv.URL, nil, nil, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestBearerToken(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Authorization"); got != "Bearer mytoken123" {
			t.Errorf("authorization: want %q, got %q", "Bearer mytoken123", got)
		}
		writeJSON(w, apiResp{})
	}))
	defer srv.Close()

	if err := New(BearerToken("mytoken123")).Get(srv.URL, nil, nil, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRequestBefore(t *testing.T) {
	called := false
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, apiResp{})
	}))
	defer srv.Close()

	c := New(RequestBefore(func(req *http.Request) {
		called = true
		req.Header.Set("X-Custom", "hook")
	}))
	c.Get(srv.URL, nil, nil, nil)
	if !called {
		t.Fatal("requestBefore hook was not called")
	}
}

func TestResponseAfter_Error(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, apiResp{})
	}))
	defer srv.Close()

	c := New(ResponseAfter(func(res *http.Response) error {
		return fmt.Errorf("hook error")
	}))
	err := c.Get(srv.URL, nil, nil, nil)
	if err == nil || err.Error() != "hook error" {
		t.Fatalf("expected hook error, got %v", err)
	}
}

// ----- text format -----

func TestFormat_Text_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("hello world"))
	}))
	defer srv.Close()

	c := New(Format(FormatText))
	var out []byte
	if err := c.Get(srv.URL, nil, nil, &out); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(out) != "hello world" {
		t.Fatalf("body: want %q, got %q", "hello world", string(out))
	}
}

func TestFormat_Text_4xx(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("bad request"))
	}))
	defer srv.Close()

	if err := New(Format(FormatText)).Get(srv.URL, nil, nil, nil); err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestFormat_Text_WrongReceiver(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello"))
	}))
	defer srv.Close()

	var out string
	if err := New(Format(FormatText)).Get(srv.URL, nil, nil, &out); err == nil {
		t.Fatal("expected error for non-*[]byte receiver")
	}
}

// ----- JSON response handling -----

func TestJsonResponse_Failed(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, apiResp{Code: 1001, Msg: "business error"})
	}))
	defer srv.Close()

	var out apiResp
	err := New().Get(srv.URL, nil, nil, &out)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Error() != "[1001] business error" {
		t.Fatalf("error: want %q, got %q", "[1001] business error", err.Error())
	}
}

func TestErrtyp_4xx(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(apiErrResp{ErrCode: 1001, ErrMsg: "bad request"})
	}))
	defer srv.Close()

	c := New(Error(apiErrResp{}))
	var out apiResp
	err := c.Get(srv.URL, nil, nil, &out)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	errResp, ok := err.(*apiErrResp)
	if !ok {
		t.Fatalf("expected *apiErrResp, got %T", err)
	}
	if errResp.ErrCode != 1001 {
		t.Fatalf("err_code: want 1001, got %d", errResp.ErrCode)
	}
}

func TestJsonResponse_Redirect(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/src" {
			http.Redirect(w, r, "/dst", http.StatusFound)
			return
		}
		writeJSON(w, apiResp{Code: 0, Msg: "final"})
	}))
	defer srv.Close()

	c := New(DisableRedirect())
	err := c.Get(srv.URL+"/src", nil, nil, nil)
	if err == nil {
		t.Fatal("expected redirect error, got nil")
	}
}

// ----- ResponseError -----

func TestResponseError(t *testing.T) {
	re := newResponseError(404, "not found")
	if re.Status() != 404 {
		t.Fatalf("status: want 404, got %d", re.Status())
	}
	want := "#404 not found"
	if re.Error() != want {
		t.Fatalf("error: want %q, got %q", want, re.Error())
	}
}

// ----- ForwardBinary -----

func TestForwardBinary(t *testing.T) {
	payload := []byte{0x89, 0x50, 0x4e, 0x47}

	src := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/png")
		w.Write(payload)
	}))
	defer src.Close()

	var receivedData []byte
	var receivedCT string
	dst := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedCT = r.Header.Get("Content-Type")
		receivedData, _ = io.ReadAll(r.Body)
		writeJSON(w, apiResp{})
	}))
	defer dst.Close()

	if err := New().ForwardBinary(dst.URL, nil, src.URL, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !bytes.Equal(receivedData, payload) {
		t.Fatal("forwarded data mismatch")
	}
	if !strings.HasPrefix(receivedCT, "image/png") {
		t.Fatalf("content-type: want image/png, got %s", receivedCT)
	}
}
