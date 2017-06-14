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
	FFmpegCmdGenerateSub string = "ffmpeg -i %s -map 0:%s -y %s"
	FFmpegCmdShowInfo    string = "ffmpeg -i %s"
	ResFolder            string = "./res"
	DoneFolder           string = "./done"
	MovieExt             string = ".mkv"
)

func init() {
	if !FileExist(ResFolder) {
		panic("找不到res文件夹")
	}
	if err := os.MkdirAll(DoneFolder, os.ModeDir); err != nil {
		panic("建立done文件夹失败")
	}
}

//由于vscode不支持stdin, 为了debug需要额外处理
const (
	DEBUG bool = false
)

func main() {
	defer db.Close()

	fmt.Println("res文件夹下的mkv文件如下,请选择文件编号开始处理:")
	movies := ListFilesWithExt("res", MovieExt)

	if len(movies) <= 0 {
		panic("在res下找不到mkv文件")
	}
	for i, v := range movies {
		fmt.Printf("%d. %s\n", i+1, v)
	}

	movieIndex := -1
	if DEBUG {
		movieIndex = 1
	} else {
		fmt.Scanln(&movieIndex)
		if movieIndex < 1 || movieIndex > len(movies) {
			panic("找不到对应的编号")
		}
	}

	movieFile := movies[movieIndex-1]

	movieName := strings.TrimSuffix(movieFile, MovieExt)
	movieFilePath := filepath.Join(ResFolder, movieFile)
	if !FileExist(movieFilePath) {
		panic("在res文件夹下找不到该电影")
	}

	//....
	ExecCmd(fmt.Sprintf(FFmpegCmdShowInfo, movieFilePath))
	fmt.Println("根据上面信息选择要生成的字幕类型编号:")

	subIndex := ""
	if DEBUG {
		subIndex = "4"
	} else {
		fmt.Scanln(&subIndex)
	}

	assFile := Md5Name(movieFile) + ".ass"
	assFilePath := filepath.Join(ResFolder, assFile)
	os.Remove(assFilePath)
	if FileExist(assFilePath) {
		panic("在res文件夹下无法删除陈旧的字幕文件: " + assFile)
	}
	fmt.Println("在res文件夹下生成字幕文件: " + assFile)
	ExecCmd(fmt.Sprintf(FFmpegCmdGenerateSub, movieFilePath, subIndex, assFilePath))

	//...
	if !FileExist(assFilePath) {
		panic("无法找到匹配的字幕流:" + subIndex)
	}

	assReader := assreader.NewAssReader(assFilePath)
	if assReader == nil {
		panic("cannot find valid ass file")
	}
	movieMd5 := Md5WithFile(movieFilePath)
	movieModel := db.Movie{Name: movieName, FileName: movieFile, Md5: movieMd5}
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
	os.Remove(assFilePath)
	os.Rename(movieFilePath, filepath.Join(DoneFolder, movieFile))
}
