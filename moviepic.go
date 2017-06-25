package main

import (
	"fmt"
	"log"
	"moviepic/db"
	"myutils"
	"net/http"
	_ "net/http/pprof"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

//
const (
	Prefix               string = ".jpg"
	FFmpegCmdGeneratePic string = "ffmpeg -ss %f -i %s -vframes 1 -copyts -af asetpts=PTS-%f/TB -vf ass=%s,setpts=PTS-%f/TB -y %s"
	FFmpegCmdGenerateSub string = "ffmpeg -i %s -map 0:%d -y %s"
	FFmpegCmdShowInfo    string = "ffmpeg -i %s"
	ResFolder            string = "./res"
	DoneFolder           string = "./done"
	DropFolder           string = "./drop"
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
	runtime.GOMAXPROCS(4)
	defer db.Close()
	go func() {
		log.Println(http.ListenAndServe(":6060", nil))
	}()

	fmt.Println("res文件夹下的mkv文件如下:")
	movies := GenerateMovieListByFolder("res", MovieExt)

	if len(movies) <= 0 {
		panic("在res下找不到mkv文件")
	}

	fmt.Println("请为如下资源选取合适的字幕编号(0代表放弃):")
	newMovies := []Movie{}
	for _, movie := range movies {
		if !movie.IsValidAssStream() {
			movie.PrintAVInfo()
			subIndex := 0
			if DEBUG {
				subIndex = 4
			} else {
				fmt.Scanln(&subIndex)
			}
			if subIndex == 0 {
				os.Rename(movie.UtilsMovie.FilePath, filepath.Join(DropFolder, movie.DbMovie.FileName))
				continue
			} else {
				movie.DbMovie.AssStream = subIndex
				movie.DbMovie.Save()

			}
		}
		newMovies = append(newMovies, movie)
	}
	movies = newMovies
	fmt.Println("所有资源的字幕处理完毕")

	// IdleMovies := []string{}
	// processingMovies := []string{}
	// for i, v := range movies {
	// 	_, assFilePath := GenAssNameAndPath(movieFile, movieFilePath)
	// 	if FileExist(assFilePath) {
	// 		processingMovies = append(processingMovies, v)
	// 	} else {
	// 		IdleMovies = append(IdleMovies, v)
	// 	}
	// }

	// movies = IdleMovies
	// for i, v := range movies {
	// 	fmt.Printf("  %d. %s\n", i+1, v)
	// }
	// fmt.Println("下列电影正在被处理中:")
	// for i, v := range processingMovies {
	// 	fmt.Printf("  %d. %s\n", i+1, v)
	// }

	fmt.Println("请选择需要处理的电影的编号, 0代表全部:")
	for i, movie := range movies {
		fmt.Printf("  %d. %s\n", i+1, movie.DbMovie.Name)
	}

	movieIndex := -1
	if DEBUG {
		movieIndex = 1
	} else {
		fmt.Scanln(&movieIndex)
	}
	if movieIndex < 0 || movieIndex > len(movies) {
		panic("找不到对应的编号")
	}

	fmt.Println("完成后是否自动关机?(1代表是)")
	isAutoShutdown := ""
	if DEBUG {
		isAutoShutdown = "no"
	} else {
		fmt.Scanln(&isAutoShutdown)
	}

	if movieIndex > 0 {
		movie := movies[movieIndex-1]
		movie.StartDeal()
	} else {
		fmt.Println("选择同时处理的任务数:")
		taskPool := 1
		if DEBUG {
			taskPool = 1
		} else {
			fmt.Scanln(&taskPool)
			if taskPool < 1 {
				taskPool = 1
			}
		}
		pool := myutils.NewPool(taskPool)
		for _, movie := range movies {
			pool.Add(1)
			go func(movie Movie) {
				fmt.Printf("start task:%s\n", movie.DbMovie.Name)
				movie.StartDeal()
				pool.Done()
			}(movie)
		}
		pool.Wait()
	}

	//shudown the pc
	if isAutoShutdown == "1" {
		time.Sleep(10 * time.Second)
		myutils.Shutdown()
	}
	// movieFile := movies[movieIndex-1]

	// movieName := strings.TrimSuffix(movieFile, MovieExt)
	// movieFilePath := filepath.Join(ResFolder, movieFile)
	// if !FileExist(movieFilePath) {
	// 	panic("在res文件夹下找不到该电影")
	// }

	// movieMd5 := Md5WithFile(movieFilePath)
	// movieModel := db.Movie{Name: movieName, FileName: movieFile, Md5: movieMd5}
	// movieModel.Save()
	// if movieModel.ID <= 0 {
	// 	panic("failed to add movie to db or cannot find")
	// }

	// //....
	// subIndex := ""
	// if movieModel.AssStream > 0 {
	// 	subIndex = strconv.Itoa(movieModel.AssStream)
	// } else {
	// 	ExecCmd(fmt.Sprintf(FFmpegCmdShowInfo, movieFilePath))
	// 	fmt.Println("根据上面信息选择要生成的字幕类型编号:")

	// 	if DEBUG {
	// 		subIndex = "4"
	// 	} else {
	// 		fmt.Scanln(&subIndex)
	// 		movieModel.AssStream, _ = strconv.Atoi(subIndex)
	// 		movieModel.Save()
	// 	}
	// }

	// assFile, assFilePath := GenAssNameAndPath(movieFile, movieFilePath)
	// os.Remove(assFilePath)
	// if FileExist(assFilePath) {
	// 	panic("在res文件夹下无法删除陈旧的字幕文件: " + assFile)
	// }
	// fmt.Println("在res文件夹下生成字幕文件: " + assFile)
	// ExecCmd(fmt.Sprintf(FFmpegCmdGenerateSub, movieFilePath, subIndex, assFilePath))

	// //...
	// if !FileExist(assFilePath) {
	// 	panic("无法找到匹配的字幕流:" + subIndex)
	// }

	// assReader := assreader.NewAssReader(assFilePath)
	// if assReader == nil {
	// 	panic("cannot find valid ass file")
	// }

	// RootPath := filepath.Join("./pics", strconv.Itoa(movieModel.ID))
	// os.MkdirAll(RootPath, os.ModeDir)
	// for _, sub := range assReader.Subs {
	// 	subModel := db.Subtitle{Text: sub.Text.Text, Format: strings.Join(sub.Text.Formats, ""),
	// 		Start: int(sub.Start * 100), End: int(sub.End * 100), MovieID: movieModel.ID}
	// 	subModel.Save()

	// 	//生成图片
	// 	rand := RandIndex()
	// 	picName := strconv.Itoa(rand) + Prefix
	// 	picPath := filepath.Join(RootPath, picName)
	// 	picTime := sub.Start + 0.5

	// 	picModel := db.Pic{Name: rand, Time: int(picTime * 100), SubtitleID: subModel.ID}
	// 	if picModel.Exist() {
	// 		continue
	// 	}

	// 	cmd := fmt.Sprintf(FFmpegCmdGeneratePic, picTime, filepath.Join(ResFolder, movieFile), picTime,
	// 		strings.Replace(filepath.Join(ResFolder, assFile), "\\", "/", -1), picTime, picPath)
	// 	// cmd = strings.Replace(cmd, "\\", "/", -1)
	// 	ExecCmd(cmd)
	// 	if !FileExist(picPath) {
	// 		panic("cannot generate the pic")
	// 	}
	// 	picModel.Save()
	// }
	// os.Remove(assFilePath)
	// os.Rename(movieFilePath, filepath.Join(DoneFolder, movieFile))
}
