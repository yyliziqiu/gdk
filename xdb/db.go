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
	_configs map[string]Config

	_rawDbs map[string]*sql.DB
	_ormDbS map[string]*gorm.DB
)

func Init(configs ...Config) error {
	_configs = make(map[string]Config, 8)
	for _, config := range configs {
		conf := config.Default()
		_configs[conf.Id] = conf
	}

	_rawDbs = make(map[string]*sql.DB, 8)
	_ormDbS = make(map[string]*gorm.DB, 8)
	for _, config := range _configs {
		raw, err := New(config)
		if err != nil {
			Finally()
			return err
		}
		_rawDbs[config.Id] = raw

		if !config.EnableOrm {
			continue
		}
		orm, err := NewOrm(config, raw)
		if err != nil {
			Finally()
			return err
		}
		_ormDbS[config.Id] = orm
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
	for _, db := range _rawDbs {
		_ = db.Close()
	}
}

func Get(id string) *sql.DB {
	return _rawDbs[id]
}

func GetDefault() *sql.DB {
	return Get(DefaultId)
}

func GetOrm(id string) *gorm.DB {
	return _ormDbS[id]
}

func GetOrmDefault() *gorm.DB {
	return GetOrm(DefaultId)
}

func GetConfig(id string) Config {
	return _configs[id]
}

func GetConfigDefault() Config {
	return GetConfig(DefaultId)
}
