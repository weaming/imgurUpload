package config

import (
	"fmt"
	"os"

	"github.com/weaming/imgurUpload/remote"
)

var Auth = remote.Auth{}

func ReadEnvs() *remote.Auth {
	id := os.Getenv("IMGUR_CLIENT_ID")
	secret := os.Getenv("IMGUR_CLIENT_SECRET")

	if id == "" || secret == "" {
		fmt.Println("In order to use imgurUpload, both IMGUR_CLIENT_ID and IMGUR_CLIENT_SECRET should be present in your environment variables.")
		fmt.Println("And your app redirect url must set to http://127.0.0.1:1024 to receive credentials.")
		os.Exit(1)
	}

	Auth.ID = id
	Auth.Secret = secret
	return &Auth
}
