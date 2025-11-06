package xboot

import (
	"context"

	"github.com/yyliziqiu/gdk/xlog"
	"github.com/yyliziqiu/gdk/xutil"
)

type InitFunc func() error

type InitFuncs []InitFunc

func (list InitFuncs) Init() error {
	for _, fun := range list {
		xlog.Infof("Init moudle: %s", xutil.ReflectFuncName(fun))
		err := fun()
		if err != nil {
			return err
		}
	}
	return nil
}

type BootFunc func(context.Context) error

type BootFuncs []BootFunc

func (list BootFuncs) Boot(ctx context.Context) error {
	for _, fun := range list {
		xlog.Infof("Boot moudle: %s", xutil.ReflectFuncName(fun))
		err := fun(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

// Check 检查配置是否正确
type Check interface {
	Check() error
}

// Default 为配置设置默认值
type Default interface {
	Default()
}
