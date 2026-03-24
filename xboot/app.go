package xboot

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/yyliziqiu/gdk/xlog"
	"github.com/yyliziqiu/gdk/xutil"
)

type InitFunc func() error

type BootFunc func(context.Context) error

type ConfigChecker interface {
	Check() error
}

type SetDefaultConfig interface {
	Default()
}

type App struct {
	// app 根命令
	Command string

	// app 版本
	Version string

	// 默认配置文件路径，可被命令行参数覆盖
	ConfigFile string

	// 全局配置，必须是指针
	Config any

	// 是否自动初始化基础组件
	AutoInject bool

	// 初始化/启动模块
	InitFuncs []InitFunc
	BootFuncs []BootFunc

	hasInitConfig bool
	hasInitModule bool
	hasBootModule bool
}

// Init app
func (t *App) Init() error {
	if err := t.InitConfig(); err != nil {
		return err
	}

	xlog.Infof("App version: %s", t.Version)

	return t.InitModule()
}

func (t *App) InitConfig() error {
	if t.hasInitConfig {
		return nil
	}
	t.hasInitConfig = true

	// 加载配置文件
	log.Println("Init config.")
	if err := InitConfig(t.ConfigFile, t.Config); err != nil {
		return fmt.Errorf("init config failed [%v]", err)
	}

	// 检查配置是否正确
	if chk, ok := t.Config.(ConfigChecker); ok {
		if err := chk.Check(); err != nil {
			return err
		}
	}

	// 为配置设置默认值
	if def, ok := t.Config.(SetDefaultConfig); ok {
		def.Default()
	}

	// 初始化日志
	log.Println("Init log.")
	lc0 := xlog.Config{Console: true}
	if lc1, ok1 := t.parseConfig("Log"); ok1 {
		if lc2, ok2 := lc1.(xlog.Config); ok2 {
			lc0 = lc2
		}
	}
	if err := xlog.Init(lc0); err != nil {
		return fmt.Errorf("init log failed [%v]", err)
	}

	return nil
}

func (t *App) parseConfig(field string) (any, bool) {
	return xutil.ReflectValueByField(t.Config, field)
}

func (t *App) InitModule() error {
	if t.hasInitModule {
		return nil
	}
	t.hasInitModule = true

	// 自动注入基础组件
	if t.AutoInject {
		if err := InitBaseComponents(t.Config)(); err != nil {
			return err
		}
	}

	// 初始化模块
	for _, fun := range t.InitFuncs {
		xlog.Infof("Init module: %s", xutil.ReflectFuncName(fun))
		if err := fun(); err != nil {
			xlog.Errorf("Init modules failed, error: %v", err)
			return err
		}
	}

	return nil
}

// Run app
func (t *App) Run() error {
	// 初始化配置和模块
	if err := t.Init(); err != nil {
		return err
	}

	// 启动模块
	ctx, cancel := context.WithCancel(context.Background())
	if err := t.BootModule(ctx); err != nil {
		cancel()
		return err
	}
	xlog.Info("App boot successfully.")

	// 等待程序退出
	exitCh := make(chan os.Signal, 1)
	signal.Notify(exitCh, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-exitCh
	cancel()

	//  延时退出
	if wait1, ok1 := t.parseConfig("ExitWait"); ok1 {
		if wait2, ok2 := wait1.(time.Duration); ok2 && wait2 > 0 {
			time.Sleep(wait2)
		}
	}
	xlog.Info("App exit.")

	return nil
}

func (t *App) BootModule(ctx context.Context) error {
	if t.hasBootModule {
		return nil
	}
	t.hasBootModule = true

	for _, fun := range t.BootFuncs {
		xlog.Infof("Boot module: %s", xutil.ReflectFuncName(fun))
		if err := fun(ctx); err != nil {
			xlog.Errorf("Boot modules failed, error: %v", err)
			return err
		}
	}

	return nil
}

func (t *App) RunWithoutBlock() (func(), error) {
	// 初始化配置和模块
	if err := t.Init(); err != nil {
		return nil, err
	}

	// 启动模块
	ctx, cancel := context.WithCancel(context.Background())
	err := t.BootModule(ctx)
	if err == nil {
		xlog.Info("App boot successfully.")
	}

	return cancel, err
}
