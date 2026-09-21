package xkvs

import (
	"strings"
	"time"
)

func (k Kvs) S(key string, def string) string {
	if val, ok := k.String(key); ok {
		return val
	}
	return def
}

func (k Kvs) B(key string, def bool) bool {
	if val, ok := k.Bool(key); ok {
		return val
	}
	return def
}

func (k Kvs) I(key string, def int) int {
	if val, ok := k.Int(key); ok {
		return val
	}
	return def
}

func (k Kvs) I64(key string, def int64) int64 {
	if val, ok := k.Int64(key); ok {
		return val
	}
	return def
}

func (k Kvs) F64(key string, def float64) float64 {
	if val, ok := k.Float64(key); ok {
		return val
	}
	return def
}

func (k Kvs) D(key string, def time.Duration) time.Duration {
	if val, ok := k.Duration(key); ok {
		return val
	}
	return def
}

func (k Kvs) S2(key string, def string) string {
	return k.S(strings.ToLower(key), def)
}

func (k Kvs) B2(key string, def bool) bool {
	return k.B(strings.ToLower(key), def)
}

func (k Kvs) I2(key string, def int) int {
	return k.I(strings.ToLower(key), def)
}

func (k Kvs) I642(key string, def int64) int64 {
	return k.I64(strings.ToLower(key), def)
}

func (k Kvs) F642(key string, def float64) float64 {
	return k.F64(strings.ToLower(key), def)
}

func (k Kvs) D2(key string, def time.Duration) time.Duration {
	return k.D(strings.ToLower(key), def)
}

func (k Kvs) Get2(key string) string {
	return k.Get(strings.ToLower(key))
}

func (k Kvs) Enabled() bool {
	return k.B("enabled", false)
}

func (k Kvs) Disabled() bool {
	return k.B("disabled", false)
}
