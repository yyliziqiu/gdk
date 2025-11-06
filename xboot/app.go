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

type App struct {
	// app 名称
	Name string

	// app 版本
	Version string

	// 默认配置文件路径，可被命令行参数覆盖
	ConfigPath string

	// 全局配置
	ConfigRoot any

	// 模块
	InitFuncs func() InitFuncs
	BootFuncs func() BootFuncs

	hasInitConfig bool
	hasInitModule bool
}

// Init app
func (app *App) Init() (err error) {
	err = app.InitConfig()
	if err != nil {
		return err
	}

	return app.InitModule()
}

func (app *App) InitConfig() (err error) {
	if app.hasInitConfig {
		return nil
	}
	app.hasInitConfig = true

	log.Println("Init config.")

	// 加载配置文件
	err = InitConfig(app.ConfigPath, app.ConfigRoot)
	if err != nil {
		return fmt.Errorf("init config failed [%v]", err)
	}

	// 检查配置是否正确
	check, ok := app.ConfigRoot.(Check)
	if ok {
		err = check.Check()
		if err != nil {
			return err
		}
	}

	// 为配置设置默认值
	def, ok := app.ConfigRoot.(Default)
	if ok {
		def.Default()
	}

	log.Println("Init log.")

	// 初始化日志
	lc := xlog.Config{Console: true}
	if lc1, ok1 := xutil.ReflectFieldValue(app.ConfigRoot, "Log"); ok1 {
		if lc2, ok2 := lc1.(xlog.Config); ok2 {
			lc = lc2
		}
	}
	err = xlog.Init(lc)
	if err != nil {
		return fmt.Errorf("init log failed [%v]", err)
	}

	return nil
}

func (app *App) InitModule() (err error) {
	if app.hasInitModule {
		return nil
	}
	app.hasInitModule = true

	err = app.InitFuncs().Init()
	if err != nil {
		xlog.Errorf("Init modules failed, error: %v", err)
		return err
	}

	return nil
}

// Run app
func (app *App) Run() (err error) {
	err = app.Init()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())

	err = app.BootFuncs().Boot(ctx)
	if err != nil {
		xlog.Errorf("Boot modules failed, error: %v", err)
		cancel()
		return err
	}

	xlog.Info("App boot successfully.")

	exitCh := make(chan os.Signal)
	signal.Notify(exitCh, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-exitCh

	cancel()

	if wait, ok := xutil.ReflectFieldValue(app.ConfigRoot, "ExitWait"); ok {
		wait2, ok2 := wait.(time.Duration)
		if ok2 && wait2 > 0 {
			time.Sleep(wait2)
		}
	}

	xlog.Info("App exit.")

	return nil
}
