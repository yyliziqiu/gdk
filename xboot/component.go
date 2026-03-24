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

func InitBaseComponents(config any) InitFunc {
	return func() error {
		// database
		if v, vok := xutil.AttemptReflectValueByField(config, []string{"Db", "DB"}); vok {
			if c, cok := v.(xdb.Config); cok && c.Dsn != "" {
				xlog.Info("Init database.")
				if err := xdb.Init(c); err != nil {
					return fmt.Errorf("init database failed [%v]", err)
				}
			}
		}
		if v, vok := xutil.AttemptReflectValueByField(config, []string{"Dbs", "DBS", "DbList", "DBList"}); vok {
			if c, cok := v.([]xdb.Config); cok && len(c) > 0 {
				xlog.Info("Init database list.")
				if err := xdb.Init(c...); err != nil {
					return fmt.Errorf("init database list failed [%v]", err)
				}
			}
		}

		// elastic search
		if v, vok := xutil.AttemptReflectValueByField(config, []string{"Es", "ES", "ElasticSearch"}); vok {
			if c, cok := v.(xes.Config); cok && len(c.Hosts) > 0 {
				xlog.Info("Init elastic search.")
				if err := xes.Init(c); err != nil {
					return fmt.Errorf("init elastic search failed [%v]", err)
				}
			}
		}
		if v, vok := xutil.AttemptReflectValueByField(config, []string{"EsList", "ESList", "ElasticSearchList"}); vok {
			if c, cok := v.([]xes.Config); cok && len(c) > 0 {
				xlog.Info("Init elastic search list.")
				if err := xes.Init(c...); err != nil {
					return fmt.Errorf("init elastic search list failed [%v]", err)
				}
			}
		}

		// redis
		if v, vok := xutil.AttemptReflectValueByField(config, []string{"Redis"}); vok {
			if c, cok := v.(xredis.Config); cok && (c.Addr != "" || len(c.Addrs) > 0) {
				xlog.Info("Init redis.")
				if err := xredis.Init(c); err != nil {
					return fmt.Errorf("init redis failed [%v]", err)
				}
			}
		}
		if v, vok := xutil.AttemptReflectValueByField(config, []string{"Redises", "RedisList"}); vok {
			if c, cok := v.([]xredis.Config); cok && len(c) > 0 {
				xlog.Info("Init redis list.")
				if err := xredis.Init(c...); err != nil {
					return fmt.Errorf("init redis list failed [%v]", err)
				}
			}
		}

		// kafka
		if v, vok := xutil.AttemptReflectValueByField(config, []string{"Kafka"}); vok {
			if c, cok := v.(xkafka.Config); cok && c.Server.BootstrapServers != "" {
				xlog.Info("Init kafka.")
				if err := xkafka.Init(c); err != nil {
					return fmt.Errorf("init kafka failed [%v]", err)
				}
			}
		}
		if v, vok := xutil.AttemptReflectValueByField(config, []string{"Kafkas", "KafkaList"}); vok {
			if c, cok := v.([]xkafka.Config); cok && len(c) > 0 {
				xlog.Info("Init kafka list.")
				if err := xkafka.Init(c...); err != nil {
					return fmt.Errorf("init kafka list failed [%v]", err)
				}
			}
		}

		return nil
	}
}

func BootBaseComponents() BootFunc {
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

// InitBase deprecated
func InitBase(config any) InitFunc {
	return InitBaseComponents(config)
}

// BootBase deprecated
func BootBase() BootFunc {
	return BootBaseComponents()
}
