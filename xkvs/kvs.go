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

func (k Kvs) GetString(key string, def ...string) string {
	if val, ok := k.String(key); ok {
		return val
	}
	if len(def) > 0 {
		return def[0]
	}
	return ""
}

func (k Kvs) GetBool(key string, def ...bool) bool {
	if val, ok := k.Bool(key); ok {
		return val
	}
	if len(def) > 0 {
		return def[0]
	}
	return false
}

func (k Kvs) GetInt(key string, def ...int) int {
	if val, ok := k.Int(key); ok {
		return val
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}

func (k Kvs) GetInt64(key string, def ...int64) int64 {
	if val, ok := k.Int64(key); ok {
		return val
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}

func (k Kvs) GetFloat64(key string, def ...float64) float64 {
	if val, ok := k.Float64(key); ok {
		return val
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}

func (k Kvs) GetDuration(key string, def ...time.Duration) time.Duration {
	if val, ok := k.Duration(key); ok {
		return val
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}

// 3. 常用 key

func (k Kvs) Get(key string) string {
	return k.GetString(key)
}

func (k Kvs) Id() string {
	return k.GetString("id")
}

func (k Kvs) Type() string {
	return k.GetString("type")
}

func (k Kvs) Name() string {
	return k.GetString("name")
}

func (k Kvs) Slug() string {
	return k.GetString("slug")
}

func (k Kvs) Key() string {
	return k.GetString("key")
}

func (k Kvs) Url() string {
	return k.GetString("url")
}

func (k Kvs) Token() string {
	return k.GetString("token")
}

func (k Kvs) Enable() bool {
	return k.GetBool("enable")
}

func (k Kvs) Disable() bool {
	return k.GetBool("disable")
}
