package xhttp

import (
	"fmt"
	"net/http"
	"testing"
)

func TestSerialHeader(t *testing.T) {
	str := SerialHeader(http.Header{
		"Server":       []string{"\"Agent\""},
		"Content-Type": []string{"application/json"},
		"Set-Cookie":   []string{"test=test", "name=name"},
		"X-Custom":     []string{},
	})

	fmt.Println("Header String: ", str)
	fmt.Println("Header Length: ", len(str))
}
