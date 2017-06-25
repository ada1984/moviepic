package main

import (
	"fmt"
	"moviepic/db"
	"os"
	"testing"
)

func TestFileExist(t *testing.T) {
	fmt.Println(FileExist("res/11.ass"))
}

func TestListFilesWithExt(t *testing.T) {
	fmt.Println(ListFilesWithExt(".", ".go"))
}

func TestMd5Name(t *testing.T) {
	fmt.Println(Md5Name("fuckme"))
}

func TestMoveFolder(t *testing.T) {
	fmt.Println(os.Rename("res\\test.txt", "res\\test.txt"))
}

func TestRemoveDuplicatePics(t *testing.T) {
	condMovie := &db.Movie{}
	movies := condMovie.FindAll()
	for _, movie := range movies {
		// movie := &db.Movie{ID: 5}
		picNames := movie.FindAllPicNamesByMovie()
		for _, picName := range picNames {
			results := db.FindDuplicatePics(movie.ID, picName)
			if len(results) > 10 {
				panic("strange result, maybe sql search is wrong")
			}
			if len(results) > 1 {
				panic("fuckme")
				fmt.Println(picName, "--", movie.ID)
				fmt.Println(os.Remove(fmt.Sprintf("pics/%d/%d.jpg", movie.ID, picName)))
				for _, pic := range results {
					pic.Remove()
				}
			}
		}
	}
}
