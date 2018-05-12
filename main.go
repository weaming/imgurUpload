package main

import (
	"flag"
	"log"
	"strings"
	"sync"

	libfs "github.com/weaming/golib/fs"
	"github.com/weaming/imgurUpload/command"
)

var (
	path      string
	anonymous = true
	log2file  = true
)

func init() {
	flag.StringVar(&path, "p", path, "target photo path/directory/url to upload")
	flag.BoolVar(&anonymous, "a", anonymous, "anonymous mode will not upload to your album")
	flag.BoolVar(&log2file, "log", log2file, "save upload results to json file")
	flag.Parse()
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
		if log2file {
			writeLog("uploaded.json", path, result)
		}
		log.Println(path, result.Data.Link)
	}
}
