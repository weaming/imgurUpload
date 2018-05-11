package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"time"

	libfs "github.com/weaming/golib/fs"
	"github.com/weaming/imgurUpload/remote"
)

var SessionFile string

func init() {
	usr, _ := user.Current()
	SessionFile = filepath.Join(usr.HomeDir, ".imgur")
}

type Session struct {
	AccessToken  string
	RefreshToken string
	ExpiresAt    time.Time
}

func (s *Session) StillValid() bool {
	return time.Now().Before(s.ExpiresAt)
}

func handleError(err error) {
	if err != nil {
		fmt.Printf("There was a problem reading your authentication data: %s\n", err)
		fmt.Printf("Deleting %s will force us to authenticate you and maybe fix the issue.\n", SessionFile)
		os.Exit(1)
	}
}

func GetSession() *Session {
	var err error
	var data Session

	if !libfs.Exist(SessionFile) {
		return nil
	}

	file, err := os.Open(SessionFile)
	handleError(err)
	defer file.Close()

	err = json.NewDecoder(file).Decode(&data)
	handleError(err)
	return &data
}

func SetSession(session *remote.AuthResponse) {
	bytes, _ := json.Marshal(session)
	ioutil.WriteFile(SessionFile, bytes, 0644)
}
