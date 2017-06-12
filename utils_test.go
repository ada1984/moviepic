package main

import (
	"fmt"
	"testing"
)

func TestMd5WithFile(t *testing.T) {
	fmt.Println(Md5WithFile("res/5.ass"))
}

func TestExecCmd(t *testing.T) {
	// cmd := fmt.Sprintf(FFmpegCmdGenerateSub, "res/【6v电影www.dy131.com】生化危机：灭绝.720p.国英双语.BD中英双字.mkv", "4")
	cmd := fmt.Sprintf(FFmpegCmdShowInfo, "res/【6v电影www.dy131.com】生化危机：灭绝.720p.国英双语.BD中英双字.mkv")
	ExecCmd(cmd)
}

func TestFileExist(t *testing.T) {
	fmt.Println(FileExist("res/11.ass"))
}

func TestListFilesWithExt(t *testing.T) {
	fmt.Println(ListFilesWithExt(".", ".go"))
}

func TestMd5Name(t *testing.T) {
	fmt.Println(Md5Name("fuckme"))
}
