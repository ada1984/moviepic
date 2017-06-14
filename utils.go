package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

//Md5WithFile ...
func Md5WithFile(path string) string {
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		log.Fatal(err)
	}

	// fmt.Sprintf(fileMd5, "%x", h.Sum(nil)[:16])
	return hex.EncodeToString(h.Sum(nil)[:16])
}

func Md5Name(name string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(name)))
}

func RandIndex() int {
	return rand.Intn(2147483647)
}

func ExecCmd(cmd string) {
	parts := strings.Fields(cmd)
	fmt.Println(parts)

	runCmd := exec.Command(parts[0], parts[1:]...)
	err, _ := runCmd.StderrPipe()
	runCmd.Start()
	maya, _ := ioutil.ReadAll(err)
	fmt.Println(string(maya))
	runCmd.Wait()
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
