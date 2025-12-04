package xhttp

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/yyliziqiu/gdk/xlog"
)

func TestClient(t *testing.T) {
	type Form struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	xlog.Init(xlog.Config{Console: true})
	cli := New(
		Logger(xlog.Default),
		LogHeader(true),
		LogEscape(true),
		Dumps(true),
	)

	header := http.Header{
		"Token": []string{"1234567890"},
	}

	form := Form{
		Username: "test",
		Password: "test",
	}

	var result Form
	err := cli.Post("http://localhost:8022/logs-request", nil, header, form, &result)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", result)
}
