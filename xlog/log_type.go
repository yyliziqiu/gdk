package xlog

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	TextFormat  = "text"
	JsonFormat  = "json"
	Text2Format = "text2"
	Json2Format = "json2"
)

type LevelDispatch map[string][]logrus.Level

type Config struct {
	Console      bool
	Path         string
	Name         string
	Timezone     string
	Location     *time.Location `json:"-" yaml:"-"`
	Level        string
	ShowCaller   bool
	DataFormat   string
	DateFormat   string
	RotateMaxAge time.Duration
	RotateTime   time.Duration
	RotateLevel  int
}

func (c Config) Default() Config {
	if c.Path == "" {
		c.Console = true
	}

	if c.Name == "" {
		c.Name = "app"
	}

	if c.Timezone == "" {
		c.Timezone = "Asia/Shanghai"
	}
	if loc, err := time.LoadLocation(c.Timezone); err == nil {
		c.Location = loc
	} else {
		c.Location = time.Local
		fmt.Printf("Parse timezone failed when init log, error: %v.\n", err)
	}

	if c.Level == "" {
		c.Level = "debug"
	}

	if c.DataFormat == "" {
		c.DataFormat = TextFormat
	}

	if c.DateFormat == "" {
		c.DateFormat = time.DateTime
	}

	if c.RotateMaxAge == 0 {
		c.RotateMaxAge = 7 * 24 * time.Hour
	}

	if c.RotateTime == 0 {
		c.RotateTime = 24 * time.Hour
	}

	if c.RotateLevel == 0 {
		c.RotateLevel = 2
	}

	return c
}

// Text2Formatter 修复日志中的时区问题
type Text2Formatter struct {
	logrus.TextFormatter
	Location *time.Location
}

func (f *Text2Formatter) Format(entry *logrus.Entry) ([]byte, error) {
	if f.Location != nil {
		entry.Time = entry.Time.In(f.Location)
	}
	return f.TextFormatter.Format(entry)
}

// Json2Formatter 修复日志中的时区问题
type Json2Formatter struct {
	logrus.JSONFormatter
	Location *time.Location
}

func (f *Json2Formatter) Format(entry *logrus.Entry) ([]byte, error) {
	if f.Location != nil {
		entry.Time = entry.Time.In(f.Location)
	}
	return f.JSONFormatter.Format(entry)
}
