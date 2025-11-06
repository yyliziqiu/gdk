package xkvs

import (
	"strings"
	"time"

	"github.com/yyliziqiu/gdk/xconv"
)

type Kvs map[string]string

// 1. key 不存在时，返回 false

func (k Kvs) String(key string) (string, bool) {
	if val, ok := k[key]; ok {
		return strings.TrimSpace(val), true
	}
	return "", false
}

func (k Kvs) Bool(key string) (bool, bool) {
	if val, ok := k.String(key); ok {
		return xconv.S2B(val), true
	}
	return false, false
}

func (k Kvs) Int(key string) (int, bool) {
	if val, ok := k.String(key); ok {
		return xconv.S2I(val), true
	}
	return 0, false
}

func (k Kvs) Int64(key string) (int64, bool) {
	if val, ok := k.String(key); ok {
		return xconv.S2I64(val), true
	}
	return 0, false
}

func (k Kvs) Float64(key string) (float64, bool) {
	if val, ok := k.String(key); ok {
		return xconv.S2F64(val), true
	}
	return 0, false
}

func (k Kvs) Duration(key string) (time.Duration, bool) {
	if val, ok := k.String(key); ok {
		d, _ := time.ParseDuration(val)
		return d, true
	}
	return 0, false
}

// 2. key 不存在时，返回默认值

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

// 3. 先将 key 全部转化为小写再查询

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

// 4. 常用 key

func (k Kvs) Get(key string) string {
	return k.S(key, "")
}

func (k Kvs) Get2(key string) string {
	return k.Get(strings.ToLower(key))
}

func (k Kvs) Id() string {
	return k.S("id", "")
}

func (k Kvs) Type() string {
	return k.S("type", "")
}

func (k Kvs) Name() string {
	return k.S("name", "")
}

func (k Kvs) Slug() string {
	return k.S("slug", "")
}

func (k Kvs) Key() string {
	return k.S("key", "")
}

func (k Kvs) Url() string {
	return k.S("url", "")
}

func (k Kvs) Token() string {
	return k.S("token", "")
}

func (k Kvs) Enabled() bool {
	return k.B("enabled", false)
}

func (k Kvs) Disabled() bool {
	return k.B("disabled", false)
}
