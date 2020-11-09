// GNU GPL v3 License
// Copyright (c) 2019 github.com:go-trellis

package tgorm

import (
	"sync"

	"gorm.io/gorm/logger"

	"github.com/go-trellis/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var locker = &sync.Mutex{}

// NewDBsFromFile initial gorm dbs from file
func NewDBsFromFile(file string) (map[string]*gorm.DB, error) {
	conf, err := config.NewConfigOptions(config.OptionFile(file))
	if err != nil {
		return nil, err
	}
	return NewDBsFromConfig(conf, "mysql")
}

// NewDBsFromConfig initial gorm engine from config
func NewDBsFromConfig(conf config.Config, name string) (map[string]*gorm.DB, error) {

	dbs := make(map[string]*gorm.DB)

	locker.Lock()
	defer locker.Unlock()

	cfg := conf.GetValuesConfig(name)

	for _, databaseName := range cfg.GetKeys() {
		logger.Default.LogMode(logger.LogLevel(cfg.GetInt(databaseName + ".log_level")))
		_db, err := gorm.Open(mysql.New(mysql.Config{
			DSN: GetMysqlDSNFromConfig(databaseName, cfg.GetValuesConfig(databaseName)),
		}), &gorm.Config{
			Logger: logger.Default,
		})
		if err != nil {
			return nil, err
		}

		tempDB, err := _db.DB()
		if err != nil {
			return nil, err
		}
		tempDB.SetMaxIdleConns(cfg.GetInt(databaseName+".max_idle_conns", 10))
		tempDB.SetMaxOpenConns(cfg.GetInt(databaseName+".max_open_conns", 100))

		if _isD := cfg.GetBoolean(databaseName + ".is_default"); _isD {
			dbs[DefaultDatabase] = _db
		}

		dbs[databaseName] = _db
	}

	return dbs, nil
}
