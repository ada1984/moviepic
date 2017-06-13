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

//FindPicByKeyword ...
func FindSubtitlesByKeyword(keyword string) []Subtitle {
	subs := []Subtitle{}
	if keyword == "" {
		return subs
	}
	db.Preload("Pic")
	db.Where("text LIKE ?", "%"+keyword+"%").Find(&subs)
	newSubs := []Subtitle{}
	for _, sub := range subs {
		db.Model(&sub).Related(&sub.Pic)
		newSubs = append(newSubs, sub)
	}
	return newSubs
	// 	pic := Pic{SubtitleID: sub.ID}
	// 	if pic.Exist() {
	// 		pics = append(pics, pic)
	// 	}
	// }
	// return pics

}
