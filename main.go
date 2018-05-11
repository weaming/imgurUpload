package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	libfs "github.com/weaming/golib/fs"
	"github.com/weaming/imgurUpload/command"
	"github.com/weaming/imgurUpload/config"
)

func main() {
	if len(os.Args) < 2 {
		if config.GetSession() == nil {
			fmt.Println("run \"imgurUpload config\" to setup config")
		} else {
			fmt.Println("missing target path (folder or path of photo)")
		}
		os.Exit(1)
	}

	path := os.Args[1]
	if path == "config" {
		command.Config()
		return
	}

	if !libfs.Exist(path) {
		fmt.Printf("target path do not exist: %v\n", path)
		os.Exit(1)
	}

	upload(path)
}

func upload(path string) {
	if strings.HasPrefix(path, "http") {
		url, e := command.UploadImageFromPath(path)
		printResult(path, url, e)
	} else {
		if libfs.IsDir(path) {
			DIR := libfs.NewDir(path)
			var wg sync.WaitGroup
			for _, p := range DIR.AbsImages {
				wg.Add(1)
				go func(p string) {
					url, e := command.UploadImageFromPath(p)
					printResult(p, url, e)
					wg.Done()
				}(p)
			}
			wg.Wait()
		} else {
			url, e := command.UploadImageFromPath(path)
			printResult(path, url, e)
		}
	}
}

func printResult(path string, url *string, e error) {
	if e != nil {
		log.Println(e)
	} else {
		log.Println(path, *url)
	}
}
