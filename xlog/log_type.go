package xlog

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	TextFormat = "text"
	JsonFormat = "json"
)

type LevelDispatch map[string][]logrus.Level

type Config struct {
	Console      bool
	Path         string
	Name         string
	Timezone     string
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

func (c Config) Location() *time.Location {
	if c.Timezone == "" {
		return time.Local
	}

	loc, err := time.LoadLocation(c.Timezone)
	if err != nil {
		fmt.Printf("Parse timezone failed when init log, error: %v.\n", err)
		return time.Local
	}

	return loc
}
