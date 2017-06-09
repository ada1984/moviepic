package main

import (
	"fmt"
	"moviepic/assreader"
	"moviepic/db"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

//
const (
	Prefix               string = ".jpg"
	FFmpegCmdGeneratePic string = "ffmpeg -ss %f -i %s -vframes 1 -copyts -af asetpts=PTS-%f/TB -vf ass=%s,setpts=PTS-%f/TB -y %s"
	FFmpegCmdGenerateSub string = "ffmpeg -i %s -map 0:%s -y res/1.ass"
	FFmpegCmdShowInfo    string = "ffmpeg -i %s"
	ResFolder            string = "./res"
)

func init() {
}

func main() {
	defer db.Close()
	//....
	fmt.Println("请输入电影名称(不需要文件后缀): ")
	MovieName := ""
	fmt.Scanln(&MovieName)
	// MovieName = "【6v电影www.dy131.com】生化危机：灭绝.720p.国英双语.BD中英双字"
	movieFile := MovieName + ".mkv"
	movieFilePath := filepath.Join(ResFolder, movieFile)
	if !FileExist(movieFilePath) {
		panic("在res文件夹下找不到该电影")
	}

	//....
	fmt.Println("根据下面信息选择要生成的字幕类型编号:")
	ExecCmd(fmt.Sprintf(FFmpegCmdShowInfo, movieFilePath))
	subIndex := ""
	fmt.Scanln(&subIndex)
	// subIndex = "4"

	assFile := "1.ass"
	assFilePath := filepath.Join(ResFolder, assFile)
	os.Remove(assFilePath)
	if FileExist(assFilePath) {
		panic("在res文件夹下无法删除陈旧的字幕文件1.ass")
	}
	fmt.Println("在res文件夹下生成字幕文件:1.ass")
	ExecCmd(fmt.Sprintf(FFmpegCmdGenerateSub, movieFilePath, subIndex))

	//...
	if !FileExist(assFilePath) {
		panic("无法找到匹配的字幕流:" + subIndex)
	}

	assReader := assreader.NewAssReader(assFilePath)
	if assReader == nil {
		panic("cannot find valid ass file")
	}
	movieMd5 := Md5WithFile(movieFilePath)
	movieModel := db.Movie{Name: MovieName, FileName: movieFile, Md5: movieMd5}
	movieModel.Save()
	if movieModel.ID <= 0 {
		panic("failed to add movie to db or cannot find")
	}
	RootPath := filepath.Join("./pics", strconv.Itoa(movieModel.ID))
	os.MkdirAll(RootPath, os.ModeDir)
	for _, sub := range assReader.Subs {
		subModel := db.Subtitle{Text: sub.Text.Text, Format: strings.Join(sub.Text.Formats, ""),
			Start: int(sub.Start * 100), End: int(sub.End * 100), MovieID: movieModel.ID}
		subModel.Save()

		//生成图片
		rand := RandIndex()
		picName := strconv.Itoa(rand) + Prefix
		picPath := filepath.Join(RootPath, picName)
		picTime := sub.Start + 0.5

		picModel := db.Pic{Name: rand, Time: int(picTime * 100), SubtitleID: subModel.ID}
		if picModel.Exist() {
			continue
		}

		cmd := fmt.Sprintf(FFmpegCmdGeneratePic, picTime, filepath.Join(ResFolder, movieFile), sub.Start+1,
			strings.Replace(filepath.Join(ResFolder, assFile), "\\", "/", -1), sub.Start+1, picPath)
		// cmd = strings.Replace(cmd, "\\", "/", -1)
		ExecCmd(cmd)
		if !FileExist(picPath) {
			panic("cannot generate the pic")
		}
		picModel.Save()
	}
}
