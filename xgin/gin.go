package xgin

import (
	"fmt"
	"io"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/yyliziqiu/gdk/xgin/xresp"
	"github.com/yyliziqiu/gdk/xlog"
)

var (
	_logger1 *logrus.Logger // 记录错误日志
	_logger2 *logrus.Logger // 记录访问日志
)

func Run(conf Config, routes ...func(engine *gin.Engine)) (err error) {
	conf = conf.Default()

	// gin 全局设置
	gin.SetMode(gin.ReleaseMode)
	gin.DisableConsoleColor()

	// 错误日志设置
	gin.DefaultErrorWriter = logger1().WriterLevel(logrus.WarnLevel)

	// 访问日志设置
	gin.DefaultWriter = io.Discard
	if !conf.DisableAccessLog {
		gin.DefaultWriter = logger2().Writer()
	}

	// 创建 gin 实例
	engine := gin.New()
	engine.NoRoute(xresp.AbortNotFound)
	engine.NoMethod(xresp.AbortMethodNotAllowed)

	// 设置全局中间件
	mws := []gin.HandlerFunc{
		gin.LoggerWithFormatter(formatter),
		gin.CustomRecovery(recovery),
	}
	engine.Use(mws...)

	// 注册路由
	for _, v := range routes {
		v(engine)
	}

	if conf.TlsEnabled() {
		// https
		err = engine.RunTLS(conf.Listen, conf.CertFile, conf.KeyFile)
	} else {
		// http
		err = engine.Run(conf.Listen)
	}

	return err
}

func logger1() *logrus.Logger {
	if _logger1 == nil {
		_logger1 = xlog.New3("gin")
	}
	return _logger1
}

func logger2() *logrus.Logger {
	if _logger2 == nil {
		_logger2 = xlog.New3("gin-access")
	}
	return _logger2
}

func formatter(param gin.LogFormatterParams) string {
	if param.Latency > time.Minute {
		param.Latency = param.Latency.Truncate(time.Second)
	}

	return fmt.Sprintf("%3d | %13v | %15s |%-7s %#v\n%s",
		param.StatusCode, param.Latency, param.ClientIP, param.Method, param.Path, param.ErrorMessage)
}

func recovery(ctx *gin.Context, err interface{}) {
	_logger1.Errorf("Server panic, path: %s, error: %v", ctx.FullPath(), err)
	xresp.AbortInternalServerError(ctx)
}

func GetLogger() *logrus.Logger {
	return logger1()
}

func SetLogger(logger *logrus.Logger) {
	_logger1 = logger
}
