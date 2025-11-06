package xreq

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

type TestForm struct {
	Name string `form:"name" binding:"required"`
}

func TestBind(t *testing.T) {
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctx.Request, _ = http.NewRequest("GET", "http://example.com?name=matrix", nil)

	var form TestForm
	ok := Bind(ctx, &form)
	if !ok {
		t.Fatal("bind failed")
	}

	t.Log(form)
}
