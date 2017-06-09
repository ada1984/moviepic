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

func FileExist(path string) bool {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}
