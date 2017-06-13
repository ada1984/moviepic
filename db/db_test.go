package db

import (
	"fmt"
	"strings"
	"testing"
)

func TestInit(t *testing.T) {
	movie := &Movie{Name: "速度与激情5", Md5: "dddd", FileName: ""}
	// db.First(movie)
	// db.Create(movie)
	db.Where(movie).First(movie)
	fmt.Println(movie)
	Close()
}

func TestDemo(t *testing.T) {
	texts := []string{"a", "b", "c"}
	fmt.Println(strings.Join(texts, ""))
}

func TestFindSubtitlesByKeyword(t *testing.T) {
	fmt.Println(FindSubtitlesByKeyword("吹牛"))
}

// func TestPicExist(t *testing.T){
// 	pic := Pic{SubtitleID:}
// 	fmt.Println()
// }
