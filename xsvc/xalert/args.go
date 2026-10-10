package xalert

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/yyliziqiu/gdk/xconv"
	"github.com/yyliziqiu/gdk/xif"
	"github.com/yyliziqiu/gdk/xlog"
)

var _replacer = strings.NewReplacer(
	" ", "_",
	";", ",",
	"=", ":",
	"\r", " ",
	"\n", " ",
	"\t", " ",
)

type Args map[string]string

func NA1(kvs ...any) (Args, error) {
	args := make(Args, 16)
	return args.Add(kvs...)
}

func NA2(kvs ...any) Args {
	args := make(Args, 16)
	return args.MustAdd(kvs...)
}

func NA3(kvs ...any) (Args, string) {
	args := make(Args, 16)
	args.MustAdd(kvs...)
	return args, args.Key()
}

func (t Args) MustAdd(kvs ...any) Args {
	if _, err := t.Add(kvs...); err != nil {
		xlog.Errorf("[Args.MustAdd] Args error: %v", err)
	}
	return t
}

func (t Args) Add(kvs ...any) (Args, error) {
	if len(kvs)%2 != 0 {
		return t, errors.New("the number of args must be even")
	}

	for i := 0; i < len(kvs); i += 2 {
		k, ok := kvs[i].(string)
		if !ok {
			return t, errors.New("the key type of args must be a string")
		}

		switch v := kvs[i+1].(type) {
		case string:
			t[k] = v
		case bool:
			t[k] = xif.If(v, "是", "否")
		case int:
			t[k] = xconv.I2S(v)
		case int64:
			t[k] = xconv.I642S(v)
		case float64:
			t[k] = xconv.F642S(v, 2)
		case error:
			t[k] = v.Error()
		default:
			t[k] = fmt.Sprintf("%v", v)
		}
	}

	return t, nil
}

func (t Args) Key() string {
	ks := make([]string, 0)
	for k := range t {
		if k[0] == '@' && len(k) > 1 {
			ks = append(ks, k)
		}
	}

	sort.Strings(ks)

	labels := make([]string, 0, len(ks))
	for _, k := range ks {
		labels = append(labels, fmt.Sprintf("%s=%s", k[1:], _replacer.Replace(t[k])))
	}

	return strings.Join(labels, ";")
}
