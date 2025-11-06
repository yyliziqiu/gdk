package xboot

import (
	"context"
	"fmt"

	"github.com/yyliziqiu/gdk/xdb"
	"github.com/yyliziqiu/gdk/xes"
	"github.com/yyliziqiu/gdk/xkafka"
	"github.com/yyliziqiu/gdk/xlog"
	"github.com/yyliziqiu/gdk/xredis"
	"github.com/yyliziqiu/gdk/xutil"
)

func InitBase(config any) InitFunc {
	return func() (err error) {
		// db
		if val, ok := xutil.AttemptReflectFieldValue(config, []string{"Db", "DB"}); ok {
			c, ok2 := val.(xdb.Config)
			if ok2 && c.Dsn != "" {
				xlog.Info("Init database.")
				err = xdb.Init(c)
				if err != nil {
					return fmt.Errorf("init database failed [%v]", err)
				}
			}
		}
		if val, ok := xutil.AttemptReflectFieldValue(config, []string{"Dbs", "DBS", "DbList", "DBList"}); ok {
			c, ok2 := val.([]xdb.Config)
			if ok2 && len(c) > 0 {
				xlog.Info("Init database list.")
				err = xdb.Init(c...)
				if err != nil {
					return fmt.Errorf("init database list failed [%v]", err)
				}
			}
		}

		// es
		if val, ok := xutil.AttemptReflectFieldValue(config, []string{"Es", "ES", "ElasticSearch"}); ok {
			c, ok2 := val.(xes.Config)
			if ok2 && len(c.Hosts) > 0 {
				xlog.Info("Init es.")
				err = xes.Init(c)
				if err != nil {
					return fmt.Errorf("init es failed [%v]", err)
				}
			}
		}
		if val, ok := xutil.AttemptReflectFieldValue(config, []string{"EsList", "ESList", "ElasticSearchList"}); ok {
			c, ok2 := val.([]xes.Config)
			if ok2 && len(c) > 0 {
				xlog.Info("Init es list.")
				err = xes.Init(c...)
				if err != nil {
					return fmt.Errorf("init es list failed [%v]", err)
				}
			}
		}

		// redis
		if val, ok := xutil.AttemptReflectFieldValue(config, []string{"Redis"}); ok {
			c, ok2 := val.(xredis.Config)
			if ok2 && (c.Addr != "" || len(c.Addrs) > 0) {
				xlog.Info("Init redis.")
				err = xredis.Init(c)
				if err != nil {
					return fmt.Errorf("init redis failed [%v]", err)
				}
			}
		}
		if val, ok := xutil.AttemptReflectFieldValue(config, []string{"Redises", "RedisList"}); ok {
			c, ok2 := val.([]xredis.Config)
			if ok2 && len(c) > 0 {
				xlog.Info("Init redis list.")
				err = xredis.Init(c...)
				if err != nil {
					return fmt.Errorf("init redis list failed [%v]", err)
				}
			}
		}

		// kafka
		if val, ok := xutil.AttemptReflectFieldValue(config, []string{"Kafka"}); ok {
			c, ok2 := val.(xkafka.Config)
			if ok2 && c.Server.BootstrapServers != "" {
				xlog.Info("Init kafka.")
				err = xkafka.Init(c)
				if err != nil {
					return fmt.Errorf("init kafka failed [%v]", err)
				}
			}
		}
		if val, ok := xutil.AttemptReflectFieldValue(config, []string{"Kafkas", "KafkaList"}); ok {
			c, ok2 := val.([]xkafka.Config)
			if ok2 && len(c) > 0 {
				xlog.Info("Init kafka list.")
				err = xkafka.Init(c...)
				if err != nil {
					return fmt.Errorf("init kafka list failed [%v]", err)
				}
			}
		}

		return nil
	}
}

func BootBase() BootFunc {
	return func(ctx context.Context) error {
		go func() {
			<-ctx.Done()
			xdb.Finally()
			xes.Finally()
			xredis.Finally()
			xkafka.Finally()
		}()
		return nil
	}
}
