package xdb

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	"github.com/yyliziqiu/gdk/xlog"
)

type Migration struct {
	Db   *gorm.DB
	Poll time.Duration
	Once []schema.Tabler
	Cron []schema.Tabler
}

func Migrates(ctx context.Context, migrations []Migration) (err error) {
	for _, migration := range migrations {
		if err = Migrate(ctx, migration); err != nil {
			return err
		}
	}
	return nil
}

func Migrate(ctx context.Context, migration Migration) (err error) {
	db := migration.Db.Set("gorm:table_options", "ENGINE=InnoDB")
	if err = migrateTables(db, migration.Once); err != nil {
		return fmt.Errorf("migrate once tables failed [%v]", err)
	}

	if len(migration.Cron) == 0 {
		return nil
	}

	if err = migrateTables(db, migration.Cron); err != nil {
		return fmt.Errorf("migrate cron tables failed [%v]", err)
	}
	go runMigrateCronTables(ctx, db, migration.Cron, migration.Poll)

	return nil
}

func migrateTables(db *gorm.DB, tables []schema.Tabler) error {
	for _, table := range tables {
		name := table.TableName()
		if has := db.Table(name).Migrator().HasTable(&table); has {
			continue
		}
		if err := db.Table(name).Migrator().CreateTable(&table); err != nil {
			return fmt.Errorf("create table %s failed [%v]", name, err)
		}
		xlog.Infof("Migration create table: %s", name)
	}
	return nil
}

func runMigrateCronTables(ctx context.Context, db *gorm.DB, tables []schema.Tabler, interval time.Duration) {
	ticker := time.NewTicker(interval)
	for {
		select {
		case <-ticker.C:
			if err := migrateTables(db, tables); err != nil {
				xlog.Errorf("Migrate cron tables failed, error: %v.", err)
			}
		case <-ctx.Done():
			return
		}
	}
}
