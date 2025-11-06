package xredis

import (
	"time"
)

const (
	DefaultId = "default"
)

const (
	ModeSingle          = "single"
	ModeCluster         = "cluster"
	ModeSentinel        = "sentinel"
	ModeSentinelCluster = "sentinel-cluster"
)

type Config struct {
	// 通用配置
	Id                 string        // optional
	Mode               string        // optional
	Username           string        // optional
	Password           string        // optional
	MaxRetries         int           // optional
	DialTimeout        time.Duration // optional
	ReadTimeout        time.Duration // optional
	WriteTimeout       time.Duration // optional
	PoolSize           int           // optional
	MinIdleConns       int           // optional
	MaxConnAge         time.Duration // optional
	PoolTimeout        time.Duration // optional
	IdleTimeout        time.Duration // optional
	IdleCheckFrequency time.Duration // optional

	// 单机模式配置
	Addr string // must
	Db   int    // optional

	// 集群模式配置
	Addrs          []string // must
	ReadPreference string   // must

	// 哨兵模式配置，还需要配置 Db 和 ReadPreference
	MasterName       string   // must
	SentinelAddrs    []string // must
	SentinelPassword string   // optional
}

func (c Config) Default() Config {
	if c.Id == "" {
		c.Id = DefaultId
	}
	if c.Mode == "" {
		c.Mode = ModeSingle
	}
	if c.MaxRetries == 0 {
		c.MaxRetries = 3
	}
	if c.DialTimeout == 0 {
		c.DialTimeout = time.Minute
	}
	if c.ReadTimeout == 0 {
		c.ReadTimeout = 3 * time.Second
	}
	if c.WriteTimeout == 0 {
		c.WriteTimeout = 3 * time.Second
	}
	if c.PoolSize == 0 {
		c.PoolSize = 10
	}
	if c.MinIdleConns == 0 {
		c.MinIdleConns = 10
	}
	if c.MaxConnAge == 0 {
		c.MaxConnAge = time.Hour
	}
	if c.PoolTimeout == 0 {
		c.PoolTimeout = 10 * time.Second
	}
	if c.IdleTimeout == 0 {
		c.IdleTimeout = 10 * time.Minute
	}
	if c.IdleCheckFrequency == 0 {
		c.IdleCheckFrequency = 30 * time.Second
	}
	if c.Addr == "" {
		c.Addr = "127.0.0.1:6379"
	}
	return c
}
