package ext

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/yyliziqiu/gdk/xconv"
	"github.com/yyliziqiu/gdk/xif"
	"github.com/yyliziqiu/gdk/xlog"
)

type Args map[string]string

func NewArgs(list ...any) (Args, error) {
	args := make(Args, 16)
	return args.Add(list...)
}

func (t Args) Add(list ...any) (Args, error) {
	if len(list)%2 != 0 {
		return t, errors.New("the number of args error")
	}

	for i := 0; i < len(list); i += 2 {
		k, ok := list[i].(string)
		if !ok {
			return t, errors.New("the key type of args must be a string")
		}

		switch v := list[i+1].(type) {
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
		default:
			t[k] = fmt.Sprintf("%v", v)
		}
	}

	return t, nil
}

func (t Args) MustAdd(list ...any) Args {
	if _, err := t.Add(list...); err != nil {
		xlog.Errorf("[Args.MustAdd] Args error: %v", err)
	}
	return t
}

var _replacer = strings.NewReplacer(
	" ", "_",
	";", ",",
	"=", ":",
	"\r", " ",
	"\n", " ",
	"\t", " ",
)

func (t Args) Key() string {
	keys := make([]string, 0)
	for k := range t {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	labels := make([]string, 0, len(keys))
	for _, k := range keys {
		labels = append(labels, fmt.Sprintf("%s=%s", k, _replacer.Replace(t[k])))
	}

	return strings.Join(labels, ";")
}
