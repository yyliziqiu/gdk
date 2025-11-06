package xdb

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	_sql map[string]*sql.DB
	_orm map[string]*gorm.DB
)

func Init(configs ...Config) error {
	_sql = make(map[string]*sql.DB, 4)
	_orm = make(map[string]*gorm.DB, 4)

	for _, config := range configs {
		cfg := config.Default()

		raw, err := New(cfg)
		if err != nil {
			Finally()
			return err
		}
		_sql[cfg.Id] = raw

		if !cfg.EnableOrm {
			continue
		}

		orm, err := NewOrm(cfg, raw)
		if err != nil {
			Finally()
			return err
		}
		_orm[cfg.Id] = orm
	}

	return nil
}

func New(config Config) (*sql.DB, error) {
	db, err := sql.Open(config.Type, config.Dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(config.MaxOpenConns)
	db.SetMaxIdleConns(config.MaxIdleConns)
	db.SetConnMaxLifetime(config.ConnMaxLifetime)
	db.SetConnMaxIdleTime(config.ConnMaxLifetime)

	return db, nil
}

func NewOrm(config Config, db *sql.DB) (*gorm.DB, error) {
	if db == nil {
		var err error
		db, err = New(config)
		if err != nil {
			return nil, err
		}
	}

	var dial gorm.Dialector
	switch config.Type {
	case TypeMysql:
		dial = mysql.New(mysql.Config{Conn: db})
	case TypePgsql:
		dial = postgres.New(postgres.Config{Conn: db})
	default:
		return nil, fmt.Errorf("not support db type %s", config.Type)
	}

	return gorm.Open(dial, config.OrmConfig())
}

func Finally() {
	for _, db := range _sql {
		_ = db.Close()
	}
}

func Get(id string) *sql.DB {
	return _sql[id]
}

func GetDefault() *sql.DB {
	return Get(DefaultId)
}

func GetOrm(id string) *gorm.DB {
	return _orm[id]
}

func GetOrmDefault() *gorm.DB {
	return GetOrm(DefaultId)
}
