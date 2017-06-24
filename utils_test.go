package main

import (
	"fmt"
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
