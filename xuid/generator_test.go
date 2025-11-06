package xuid

import (
	"testing"
	"time"
)

func TestGenerator(t *testing.T) {
	t1 := time.Now()
	for i := 0; i < 1000000; i++ {
		_ = Get()
	}
	t.Log(time.Since(t1).Milliseconds())
	t.Log(Get())
}
