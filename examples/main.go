package main

import (
	"fmt"

	"github.com/iTrellis/tgorm"
	"gorm.io/gorm"
)

var (
	dbs map[string]*gorm.DB
)

type Sample struct {
	ID   string `gorm:"id"`
	Name string `gorm:"name"`
}

func (*Sample) TableName() string {
	return "test"
}

func main() {
	var err error
	dbs, err = tgorm.NewDBsFromFile("mysql.yaml")
	if err != nil {
		panic(err)
	}

	db := dbs[tgorm.DefaultDatabase]
	if db == nil {
		panic(tgorm.ErrNotFoundGormDB)
	}

	var ss []Sample
	newDB := db.Find(&ss)
	if newDB.Error != nil {
		panic(newDB.Error)
	}
	fmt.Println(ss)
}
