package command

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/weaming/imgurUpload/config"
)

var client = &http.Client{}

func uploadImageWithBytes(data []byte, anonymous bool) (*string, error) {
	buffer := new(bytes.Buffer)
	m := multipart.NewWriter(buffer)
	label, err := m.CreateFormFile("image", "picture")
	if err != nil {
		return nil, err
	}
	label.Write(data)
	m.Close()
	req, err := http.NewRequest("POST", "https://api.imgur.com/3/image", buffer)
	if err != nil {
		return nil, err
	}

	if anonymous {
		auth := config.ReadEnvs()
		req.Header.Add("Authorization", "Client-ID "+auth.ID)
	} else {
		session := config.GetSession()
		if session == nil {
			Config()
			session = config.GetSession()
		}
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", session.AccessToken))
	}

	req.Header.Set("Content-Type", m.FormDataContentType())
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return checkUploadResult(&res.Body)
}

func checkUploadResult(bodyPtr *io.ReadCloser) (*string, error) {
	body := *bodyPtr
	defer body.Close()

	result := make(map[string]interface{})
	err := json.NewDecoder(body).Decode(&result)
	if err != nil {
		return nil, err
	}
	if result["status"].(float64) == 200 {
		link := result["data"].(map[string]interface{})["link"].(string)
		return &link, nil
	}
	return nil, errors.New("Invalid response from remote")
}

func UploadImageFromPath(path string, anonymous bool) (*string, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Error: Cannot open \"%s\": File does not exist.", path)
		os.Exit(1)
	}
	file, err := os.Open(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Cannot open \"%s\": %s", path, err)
		os.Exit(1)
	}
	defer file.Close()
	data, _ := ioutil.ReadAll(file)
	return uploadImageWithBytes(data, anonymous)
}

func UploadImageFromUrl(imgurl string, anonymous bool) (*string, error) {
	data := url.Values{}
	data.Set("image", imgurl)

	req, _ := http.NewRequest("POST", "https://api.imgur.com/3/image", bytes.NewBufferString(data.Encode()))

	if anonymous {
		auth := config.ReadEnvs()
		req.Header.Add("Authorization", "Client-ID "+auth.ID)
	} else {
		session := config.GetSession()
		if session == nil {
			Config()
			session = config.GetSession()
		}
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", session.AccessToken))
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return checkUploadResult(&res.Body)
}

func UploadImageFromStdin(anonymous bool) (*string, error) {
	data, _ := ioutil.ReadAll(os.Stdin)
	return uploadImageWithBytes(data, anonymous)
}
