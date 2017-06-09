package db

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" //
)

var db *gorm.DB = nil

func init() {
	db, _ = gorm.Open("mysql", "root:1111@/moviepic?charset=utf8&parseTime=True&loc=Local")
	if db == nil {
		panic("db open fail")
	}
	db.SingularTable(true)
	db.LogMode(true)
	// db.AutoMigrate(&Movie{})
}

func Close() {
	db.Close()
}
