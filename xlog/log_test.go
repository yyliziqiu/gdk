package xlog

import (
	"testing"
)

func TestLog(t *testing.T) {
	config := Config{
		Path:       "/private/ws/self/gdk/logs",
		Timezone:   "UTC",
		DataFormat: "text2",
	}

	logger, _ := New(config.Default())
	logger.Infof("hello world")
}
