package main

import (
	"fmt"
	"os"
	"strings"

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

func upload(path string) (*string, error) {
	if strings.HasPrefix(path, "http") {
		return command.UploadImageFromUrl(path)
	} else {
		return command.UploadImageFromPath(path)
	}
}

func ExitErr(e error) {
	if e != nil {
		fmt.Println(e)
		os.Exit(1)
	}
}
