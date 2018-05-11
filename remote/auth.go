package remote

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os/exec"
	"strconv"
	"time"
)

var client = &http.Client{}

type Auth struct {
	ID     string
	Secret string
}

type AuthResponse struct {
	AccessToken     string
	ExpiresIn       string
	ExpiresAt       time.Time
	TokenType       string
	RefreshToken    string
	AccountUsername string
	AccountID       string
}

func Authorization(auth *Auth) *AuthResponse {
	var authResponse AuthResponse
	var received = make(chan int, 1)

	localServer := func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/html")

		if req.URL.Path == "/" {
			io.WriteString(w, fmt.Sprintf("Capturing token<script>%v</script>", indexJS))
			return
		} else if req.URL.Path == "/catchtoken" {
			query := req.URL.Query()
			authResponse = AuthResponse{
				AccessToken:     query.Get("access_token"),
				ExpiresIn:       query.Get("expires_in"),
				TokenType:       query.Get("token_type"),
				RefreshToken:    query.Get("refresh_token"),
				AccountUsername: query.Get("account_username"),
				AccountID:       query.Get("account_id"),
			}
			expiresIn, _ := strconv.Atoi(authResponse.ExpiresIn)
			t := time.Now().Add(time.Duration(expiresIn) * time.Second)
			authResponse.ExpiresAt = t
			received <- 0
			io.WriteString(w, "OK")
		}
	}

	go http.ListenAndServe(":1024", http.HandlerFunc(localServer))

	authURL := "https://api.imgur.com/oauth2/authorize?response_type=token&client_id=" + auth.ID
	exec.Command("open", authURL).Run()
	// wait redirect
	<-received

	return &authResponse
}

func GetTokenFromRefreshToken(token string, auth *Auth) (*AuthResponse, error) {
	// API have bug: Invalid grant_type parameter or parameter missing
	url := "https://api.imgur.com/oauth2/token?" + url.Values{
		"client_id":     {auth.ID},
		"client_secret": {auth.Secret},
		"grant_type":    {"refresh_token"},
		"refresh_token": {token},
	}.Encode()
	resp, err := http.Get(url)

	// API have bug: These actions are forbidden.
	// resp, err := http.PostForm("https://api.imgur.com/oauth2/token", url.Values{
	// 	"client_id":     {auth.ID},
	// 	"client_secret": {auth.Secret},
	// 	"grant_type":    {"refresh_token"},
	// 	"refresh_token": {token},
	// })

	if err != nil {
		return nil, err
	}

	return getSessionFromResponse(resp)
}

func decodeJSON(resp *http.Response) (map[string]interface{}, error) {
	defer resp.Body.Close()
	result := make(map[string]interface{})
	err := json.NewDecoder(resp.Body).Decode(&result)
	return result, err
}

func getSessionFromResponse(resp *http.Response) (*AuthResponse, error) {
	result, err := decodeJSON(resp)
	if err != nil {
		return nil, err
	}

	success := result["success"].(bool)
	if success {

		t := time.Now().Add(time.Duration(result["expires_in"].(float64)) * time.Second)

		var sess = AuthResponse{
			AccessToken:  result["access_token"].(string),
			RefreshToken: result["refresh_token"].(string),
			ExpiresAt:    t,
		}

		return &sess, nil
	} else {
		return nil, errors.New(result["data"].(map[string]interface{})["error"].(string))
	}
}
