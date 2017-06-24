package main

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
)

func Md5Name(name string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(name)))
}

func RandIndex() int {
	return rand.Intn(2147483647)
}

func ReplaceBlankWithFilename([]string) {

}

func FileExist(path string) bool {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func ListFilesWithExt(path string, ext string) []string {
	extFiles := []string{}
	files, _ := ioutil.ReadDir(path)
	for _, file := range files {
		// filePath := filepath.Join(path, file.Name())
		if !file.IsDir() && (filepath.Ext(file.Name()) == ext) {
			newName := strings.Replace(file.Name(), " ", "_", -1)
			if newName != file.Name() {
				os.Rename(filepath.Join(path, file.Name()), filepath.Join(path, newName))
			}
			extFiles = append(extFiles, newName)
		}
	}
	return extFiles
}

func GenAssNameAndPath(movieFile string, ResFolder string) (string, string) {
	assFile := Md5Name(movieFile) + ".ass"
	assFilePath := filepath.Join(ResFolder, assFile)
	return assFile, assFilePath
}
