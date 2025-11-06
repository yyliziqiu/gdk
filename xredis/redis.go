package xredis

import (
	"github.com/go-redis/redis/v8"
)

var (
	_cli map[string]*redis.Client
	_clu map[string]*redis.ClusterClient
)

func Init(configs ...Config) error {
	_cli = make(map[string]*redis.Client, 4)
	_clu = make(map[string]*redis.ClusterClient, 4)

	for _, config := range configs {
		cfg := config.Default()
		cli, clu := New(cfg)
		if cli != nil {
			_cli[cfg.Id] = cli
		}
		if clu != nil {
			_clu[cfg.Id] = clu
		}
	}

	return nil
}

func New(config Config) (*redis.Client, *redis.ClusterClient) {
	switch config.Mode {
	case ModeCluster:
		return nil, NewCluster(config)
	case ModeSentinel:
		return NewSentinel(config), nil
	case ModeSentinelCluster:
		return nil, NewSentinelCluster(config)
	default:
		return NewSingle(config), nil
	}
}

func NewSingle(config Config) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:               config.Addr,
		Username:           config.Username,
		Password:           config.Password,
		DB:                 config.Db,
		MaxRetries:         config.MaxRetries,
		DialTimeout:        config.DialTimeout,
		ReadTimeout:        config.ReadTimeout,
		WriteTimeout:       config.WriteTimeout,
		PoolSize:           config.PoolSize,
		MinIdleConns:       config.MinIdleConns,
		MaxConnAge:         config.MaxConnAge,
		PoolTimeout:        config.PoolTimeout,
		IdleTimeout:        config.IdleTimeout,
		IdleCheckFrequency: config.IdleCheckFrequency,
	})
}

func NewCluster(config Config) *redis.ClusterClient {
	ops := &redis.ClusterOptions{
		Addrs:              config.Addrs,
		Username:           config.Username,
		Password:           config.Password,
		MaxRetries:         config.MaxRetries,
		DialTimeout:        config.DialTimeout,
		ReadTimeout:        config.ReadTimeout,
		WriteTimeout:       config.WriteTimeout,
		PoolSize:           config.PoolSize,
		MinIdleConns:       config.MinIdleConns,
		MaxConnAge:         config.MaxConnAge,
		PoolTimeout:        config.PoolTimeout,
		IdleTimeout:        config.IdleTimeout,
		IdleCheckFrequency: config.IdleCheckFrequency,
	}

	switch config.ReadPreference {
	case "ReadOnly":
		ops.ReadOnly = true
	case "RouteByLatency":
		ops.RouteByLatency = true
	case "RouteRandomly":
		ops.RouteRandomly = true
	}

	return redis.NewClusterClient(ops)
}

func NewSentinel(config Config) *redis.Client {
	ops := &redis.FailoverOptions{
		MasterName:         config.MasterName,
		SentinelAddrs:      config.SentinelAddrs,
		SentinelPassword:   config.SentinelPassword,
		Username:           config.Username,
		Password:           config.Password,
		DB:                 config.Db,
		MaxRetries:         config.MaxRetries,
		DialTimeout:        config.DialTimeout,
		ReadTimeout:        config.ReadTimeout,
		WriteTimeout:       config.WriteTimeout,
		PoolSize:           config.PoolSize,
		MinIdleConns:       config.MinIdleConns,
		MaxConnAge:         config.MaxConnAge,
		PoolTimeout:        config.PoolTimeout,
		IdleTimeout:        config.IdleTimeout,
		IdleCheckFrequency: config.IdleCheckFrequency,
	}

	switch config.ReadPreference {
	case "SlaveOnly":
		ops.SlaveOnly = true
	}

	return redis.NewFailoverClient(ops)
}

func NewSentinelCluster(config Config) *redis.ClusterClient {
	ops := &redis.FailoverOptions{
		MasterName:         config.MasterName,
		SentinelAddrs:      config.SentinelAddrs,
		SentinelPassword:   config.SentinelPassword,
		Username:           config.Username,
		Password:           config.Password,
		DB:                 config.Db,
		MaxRetries:         config.MaxRetries,
		DialTimeout:        config.DialTimeout,
		ReadTimeout:        config.ReadTimeout,
		WriteTimeout:       config.WriteTimeout,
		PoolSize:           config.PoolSize,
		MinIdleConns:       config.MinIdleConns,
		MaxConnAge:         config.MaxConnAge,
		PoolTimeout:        config.PoolTimeout,
		IdleTimeout:        config.IdleTimeout,
		IdleCheckFrequency: config.IdleCheckFrequency,
	}

	switch config.ReadPreference {
	case "SlaveOnly":
		ops.SlaveOnly = true
	case "RouteByLatency":
		ops.RouteByLatency = true
	case "RouteRandomly":
		ops.RouteRandomly = true
	}

	return redis.NewFailoverClusterClient(ops)
}

func Finally() {
	for _, cli := range _cli {
		_ = cli.Close()
	}
	for _, clu := range _clu {
		_ = clu.Close()
	}
}

func GetClient(id string) *redis.Client {
	return _cli[id]
}

func GetClientDefault() *redis.Client {
	return GetClient(DefaultId)
}

func GetCluster(id string) *redis.ClusterClient {
	return _clu[id]
}

func GetClusterDefault() *redis.ClusterClient {
	return GetCluster(DefaultId)
}
