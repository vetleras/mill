package mill

import (
	"encoding/json"
	"os"

	log "github.com/sirupsen/logrus"
)

type Credentials struct {
	AccessKey         string `json:"access_key"`
	SecretToken       string `json:"secret_token"`
	AuthorizationCode string `json:"authorization_code"`
	AccessToken       string `json:"access_token"`
	RefreshToken      string `json:"refresh_token"`
}

type JsonParsingErr struct{}

func (e *JsonParsingErr) Error() string {
	return "JSON parsing error"
}

func CredentialsFromFile(filepath string) (credentials *Credentials, err error) {
	file, err := os.Open(filepath)
	if err != nil {
		return
	}
	defer file.Close()

	err = json.NewDecoder(file).Decode(&credentials)
	if err != nil {
		log.Error(err)
		return nil, &JsonParsingErr{}
	}

	credentials.AccessKey, credentials.RefreshToken, err = UpdateAccessTokenAndRefreshToken(credentials.RefreshToken)
	return
}

func CredentialsFromPrompt() (credentials *Credentials, err error) {
	credentials = &Credentials{}
	credentials.AccessKey = PromptInput("Enter access key")
	credentials.SecretToken = PromptInput("Enter secret token")
	credentials.AuthorizationCode, err = AuthorizationCode(credentials.AccessKey, credentials.SecretToken)
	if err != nil {
		return
	}

	username := PromptInput("Enter username")
	password := PromptPassword("Enter password")
	credentials.RefreshToken, credentials.AccessToken, err = GetAccessAndRefreshToken(credentials.AuthorizationCode, username, password)
	if err != nil {
		return
	}
	return
}

func (credentials *Credentials) ToFile(filepath string) (err error) {
	file, err := os.OpenFile(filepath, os.O_WRONLY, 0)
	if err != nil {
		return
	}
	defer file.Close()

	err = json.NewEncoder(file).Encode(&credentials)
	if err != nil {
		return
	}
	return
}
