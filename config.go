// GNU GPL v3 License
// Copyright (c) 2019 github.com:go-trellis

package tgorm

import (
	"fmt"
	"sync"

	"github.com/go-trellis/config"
	"github.com/jinzhu/gorm"
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
		fmt.Println("_db 0", databaseName)
		_db, err := gorm.Open("mysql", GetMysqlDSNFromConfig(databaseName, cfg.GetValuesConfig(databaseName)))
		if err != nil {
			return nil, err
		}
		fmt.Println("_db 1", _db)

		_db.DB().SetMaxIdleConns(cfg.GetInt(databaseName+".max_idle_conns", 10))

		_db.DB().SetMaxOpenConns(cfg.GetInt(databaseName+".max_open_conns", 100))

		_db.LogMode(cfg.GetBoolean(databaseName + ".show_sql"))

		if _isD := cfg.GetBoolean(databaseName + ".is_default"); _isD {
			dbs[DefaultDatabase] = _db
		}

		fmt.Println("_db 2", _db)
		dbs[databaseName] = _db
	}

	fmt.Println(dbs)

	return dbs, nil
}
