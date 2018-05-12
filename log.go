package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/weaming/imgurUpload/command"
)

type LogData struct {
	HashMethod string `json:"hash_method"`
	Hash       string `json:"hash"` // sha256(file_path)
	URL        string `json:"url"`  // imgur link
	Width      int    `json:"width"`
	Height     int    `json:"height"`
	Size       int    `json:"size"`
}

type Log struct {
	Data LogData `json:"data"`
}

type Logs map[string]Log

var lock = sync.Mutex{}

func writeLog(path, fp string, result *command.UploadResponse) {
	lock.Lock()
	defer lock.Unlock()

	file, e := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0644)
	if e != nil {
		log.Println(e)
		return
	}

	content, e := ioutil.ReadAll(file)
	if e != nil {
		log.Println(e)
		return
	}

	logs := make(Logs)
	if len(content) != 0 {
		json.Unmarshal(content, &logs)
	}

	absfp, _ := filepath.Abs(fp)
	logs[absfp] = Log{
		Data: LogData{
			HashMethod: "sha256(filePath)",
			Hash:       Sha256([]byte(fp)),
			URL:        result.Data.Link,
			Width:      result.Data.Width,
			Height:     result.Data.Height,
			Size:       result.Data.Size,
		},
	}

	c, e := json.MarshalIndent(&logs, "", "    ")
	if e != nil {
		panic(e)
	}

	// write from begin
	file.Truncate(0)
	file.Seek(0, 0)

	written, e := file.Write(c)
	if e != nil || written != len(c) {
		log.Println(written, len(c), e)
		return
	}
}

func Sha256(content []byte) string {
	return fmt.Sprintf("%x", sha256.Sum256(content))
}
