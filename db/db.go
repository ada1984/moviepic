package db

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" //
	"github.com/nadoo/convtrad"
)

const (
	SQLFindDuplicatePics string = "select pic.* from moviepic.pic,moviepic.subtitle where subtitle.movie_id = ? and pic.subtitle_id = subtitle.id and pic.name = ?"
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

func FindDuplicatePics(movieId, picName int) []Pic {
	pics := []Pic{}
	rows, _ := db.Raw(SQLFindDuplicatePics, movieId, picName).Rows()
	defer rows.Close()
	for rows.Next() {
		pic := Pic{}
		// fmt.Println(rows.Scan(&pic))
		db.ScanRows(rows, &pic)
		pics = append(pics, pic)
	}
	return pics
}

//FindPicByKeyword ...
func FindSubtitlesByKeyword(keyword string, limit int) []Subtitle {
	subs := []Subtitle{}
	if keyword == "" {
		return subs
	}
	db.Preload("Pic")
	db.Where("text LIKE ? OR text LIKE ?", "%"+convtrad.ToTrad(keyword)+"%", "%"+convtrad.ToSimp(keyword)+"%").Order("RAND()").Limit(limit).Preload("Movie").Preload("Pic").Find(&subs)
	// for i := range subs {
	// db.Model(&subs[i]).Related(&subs[i].Pic)
	// db.Model(&subs[i]).Related(&subs[i].Movie)
	// }
	return subs
}
