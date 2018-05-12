package main

import (
	"flag"
	"log"
	"os"
	"strings"
	"sync"

	libfs "github.com/weaming/golib/fs"
	"github.com/weaming/imgurUpload/command"
)

var (
	path      string
	anonymous = true
	logFile   = "uploaded.json"
)

func init() {
	flag.StringVar(&path, "p", path, "target photo path/directory/url to upload")
	flag.BoolVar(&anonymous, "a", anonymous, "anonymous mode will not upload to your album")
	flag.StringVar(&logFile, "log", logFile, "log file path to save upload results in json format, will be cover by $IMGUR_LOG_FILE")
	flag.Parse()

	envLogFile := os.Getenv("IMGUR_LOG_FILE")
	if envLogFile != "" {
		logFile = envLogFile
	}
}

func main() {
	upload(path)
}

func upload(path string) {
	if strings.HasPrefix(path, "http") {
		result, e := command.UploadImageFromUrl(path, anonymous)
		printResult(path, result, e)
	} else {
		if libfs.IsDir(path) {
			DIR := libfs.NewDir(path)
			var wg sync.WaitGroup
			for _, p := range DIR.AbsImages {
				wg.Add(1)
				go func(p string) {
					result, e := command.UploadImageFromPath(p, anonymous)
					printResult(p, result, e)
					wg.Done()
				}(p)
			}
			wg.Wait()
		} else {
			result, e := command.UploadImageFromPath(path, anonymous)
			printResult(path, result, e)
		}
	}
}

func printResult(path string, result *command.UploadResponse, e error) {
	if e != nil {
		log.Println(e)
	} else {
		writeLog(logFile, path, result)
		log.Println(path, result.Data.Link)
	}
}
