package command

import (
	"fmt"

	"github.com/weaming/imgurUpload/config"
	"github.com/weaming/imgurUpload/remote"
)

func Config() (err error) {
	currentConfig := config.GetSession()
	var session *remote.AuthResponse
	authed := false

	if currentConfig != nil {
		fmt.Println("Hey! Looks like you're already authenticated!")
		/*
			if currentConfig.RefreshToken != "" {
				auth := config.ReadEnvs()
				session, err = remote.GetTokenFromRefreshToken(currentConfig.RefreshToken, auth)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				authed = true
			}
		*/
	}

	if !authed {
		auth := config.ReadEnvs()
		session = remote.Authorization(auth)
	}

	config.SetSession(session)
	fmt.Printf("Your credentials is stored to %s. You're now ready to upload gifs!\n", config.SessionFile)
	return nil
}
