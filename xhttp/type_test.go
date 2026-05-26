package xhttp

import (
	"fmt"
	"net/http"
	"testing"
)

func TestSerialHeader(t *testing.T) {
	str := SerialHeader(http.Header{
		"Content-Type": []string{"application/json"},
		"Set-Cookie":   []string{"test=test", "name=name"},
		"X-Custom1":    []string{},
		"X-Custom2":    []string{"\"Agent\""},
		"X-Custom3":    []string{"xxxxxxxxxxxxxxxxx"},
	}, []string{"X-Custom3"})

	fmt.Println("Header String: ", str)
	fmt.Println("Header Length: ", len(str))
}
