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
