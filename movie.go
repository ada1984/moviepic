package main

import (
	"fmt"
	"io/ioutil"
	"moviepic/assreader"
	"moviepic/db"
	"moviepic/utils"
	"myutils"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Movie struct {
	DbMovie    *db.Movie
	UtilsMovie *utils.Movie
}

func GenerateMovieListByFolder(folder string, ext string) []Movie {
	movies := []Movie{}
	files, _ := ioutil.ReadDir(folder)
	for _, file := range files {
		// filePath := filepath.Join(path, file.Name())
		if !file.IsDir() && (filepath.Ext(file.Name()) == ext) {
			name := strings.Replace(file.Name(), " ", "_", -1)
			if name != file.Name() {
				os.Rename(filepath.Join(folder, file.Name()), filepath.Join(folder, name))
			}
			movieName := strings.TrimSuffix(name, ext)
			movieFileName := name
			movieFilePath := filepath.Join(folder, movieFileName)
			movie := NewMovie(movieFilePath, movieName, ext)
			movies = append(movies, *movie)

		}
	}
	return movies
}

func NewMovie(path string, name string, ext string) *Movie {
	utilsMovie := &utils.Movie{FilePath: path, Ext: ext}
	md5 := myutils.Md5PartFile(path)
	dbMovie := &db.Movie{Name: name, FileName: name + ext, Md5: md5}
	dbMovie.Init()
	movie := &Movie{DbMovie: dbMovie, UtilsMovie: utilsMovie}
	return movie
}

func (movie *Movie) IsValidAssStream() bool {
	if movie.DbMovie.AssStream <= 0 {
		return false
	}
	return true
}

func (movie *Movie) PrintAVInfo() {
	_, info := myutils.ExecCmd(fmt.Sprintf(FFmpegCmdShowInfo, movie.UtilsMovie.FilePath))
	fmt.Println(info)
	fmt.Println("根据上面信息选择要生成的字幕类型编号:")
}

func (movie *Movie) StartDeal() {
	assFile, assFilePath := GenAssNameAndPath(movie.DbMovie.Name, ResFolder)
	os.Remove(assFilePath)
	if FileExist(assFilePath) {
		panic("在res文件夹下无法删除陈旧的字幕文件: " + assFile)
	}
	fmt.Println("在res文件夹下生成字幕文件: " + assFile)
	fmt.Println(myutils.ExecCmd(fmt.Sprintf(FFmpegCmdGenerateSub, movie.UtilsMovie.FilePath, movie.DbMovie.AssStream, assFilePath)))
	if !FileExist(assFilePath) {
		panic(fmt.Sprintf("无法找到匹配的字幕流: %d", movie.DbMovie.AssStream))
	}
	assReader := assreader.NewAssReader(assFilePath)
	if assReader == nil {
		panic("cannot find valid ass file")
	}
	rootPath := filepath.Join("./pics", strconv.Itoa(movie.DbMovie.ID))
	os.MkdirAll(rootPath, os.ModeDir)
	usedIndex := make(map[int]bool, len(assReader.Subs))
	for _, sub := range assReader.Subs {
		subModel := db.Subtitle{Text: sub.Text.Text, Format: strings.Join(sub.Text.Formats, ""),
			Start: int(sub.Start * 100), End: int(sub.End * 100), MovieID: movie.DbMovie.ID}
		subModel.Save()
		//生成图片
		picTime := sub.Start + 0.5
		picModel := db.Pic{Time: int(picTime * 100), SubtitleID: subModel.ID}
		if picModel.Exist() {
			continue
		}

		rand := 0
		for {
			rand = RandIndex()
			if usedIndex[rand] == false {
				usedIndex[rand] = true
				break
			}
		}

		picName := strconv.Itoa(rand) + Prefix
		picPath := filepath.Join(rootPath, picName)
		picModel.Name = rand
		cmd := fmt.Sprintf(FFmpegCmdGeneratePic, picTime, movie.UtilsMovie.FilePath, picTime,
			strings.Replace(filepath.Join(ResFolder, assFile), "\\", "/", -1), picTime, picPath)
		// cmd = strings.Replace(cmd, "\\", "/", -1)
		myutils.ExecCmdSimple(cmd)
		if !FileExist(picPath) {
			panic("cannot generate the pic")
		}
		picModel.Save()
	}
	os.Remove(assFilePath)
	fmt.Println(os.Rename(movie.UtilsMovie.FilePath, filepath.Join(DoneFolder, movie.DbMovie.FileName)))
}
