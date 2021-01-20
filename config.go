/*
Copyright © 2019 Henry Huang <hhh@rutcode.com>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/

package tgorm

import (
	"sync"

	"gorm.io/gorm/logger"

	"github.com/iTrellis/config"
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
