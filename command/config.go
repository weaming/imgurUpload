package command

import (
	"fmt"
	"sync"

	"github.com/weaming/imgurUpload/config"
	"github.com/weaming/imgurUpload/remote"
)

var one = sync.Once{}

func Config() {
	one.Do(func() {
		currentConfig := config.GetSession()
		var session *remote.AuthResponse

		// refresh token
		if currentConfig != nil {
			fmt.Println("Hey! Looks like you're already authenticated!")
		}

		auth := config.ReadEnvs()
		session = remote.Authorization(auth)

		config.SetSession(session)
		fmt.Printf("Your credentials is stored to %s. You're now ready to upload gifs!\n", config.SessionFile)
	})
}
