package xdb

import (
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	ormlog "gorm.io/gorm/logger"

	"github.com/yyliziqiu/gdk/xlog"
)

const (
	DefaultId = "default"
)

const (
	TypeMysql = "mysql"
	TypePgsql = "postgres"
)

type Config struct {
	Id              string        // optional
	Dsn             string        // must
	Type            string        // optional
	MaxOpenConns    int           // optional
	MaxIdleConns    int           // optional
	ConnMaxLifetime time.Duration // optional
	ConnMaxIdleTime time.Duration // optional

	// only valid when use gorm
	EnableOrm                       bool           // optional
	OrmLogger                       *logrus.Logger `json:"-"` // optional
	OrmLogLevel                     int            // optional
	OrmLogSlowThreshold             time.Duration  // optional
	OrmLogParameterizedQueries      bool           // optional
	OrmLogIgnoreRecordNotFoundError bool           // optional
}

func (c Config) Default() Config {
	if c.Id == "" {
		c.Id = DefaultId
	}
	if c.Type == "" {
		c.Type = TypeMysql
	}
	if c.MaxOpenConns == 0 {
		c.MaxOpenConns = 10
	}
	if c.MaxIdleConns == 0 {
		c.MaxIdleConns = 10
	}
	if c.ConnMaxLifetime == 0 {
		c.ConnMaxLifetime = time.Hour
	}
	if c.ConnMaxIdleTime == 0 {
		c.ConnMaxLifetime = 30 * time.Minute
	}
	if c.OrmLogLevel == 0 {
		c.OrmLogLevel = 1
	}
	if c.OrmLogSlowThreshold == 0 {
		c.OrmLogSlowThreshold = 5 * time.Second
	}
	if c.OrmLogLevel > 1 && c.OrmLogger == nil {
		c.OrmLogger = xlog.Default
	}
	return c
}

func (c Config) OrmConfig() *gorm.Config {
	config := ormlog.Config{
		LogLevel:                  ormlog.LogLevel(c.OrmLogLevel),
		SlowThreshold:             c.OrmLogSlowThreshold,
		ParameterizedQueries:      c.OrmLogParameterizedQueries,
		IgnoreRecordNotFoundError: c.OrmLogIgnoreRecordNotFoundError,
	}
	return &gorm.Config{Logger: ormlog.New(c.OrmLogger, config)}
}
