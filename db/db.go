package db

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" //
	"github.com/nadoo/convtrad"
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

//FindPicByKeyword ...
func FindSubtitlesByKeyword(keyword string, limit int) []Subtitle {
	subs := []Subtitle{}
	if keyword == "" {
		return subs
	}
	db.Preload("Pic")
	db.Where("text LIKE ? OR text LIKE ?", "%"+convtrad.ToTrad(keyword)+"%", "%"+convtrad.ToSimp(keyword)+"%").Order("RAND()").Limit(limit).Find(&subs)
	for i := range subs {
		db.Model(&subs[i]).Related(&subs[i].Pic)
	}
	return subs
}
