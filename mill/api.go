package mill

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	log "github.com/sirupsen/logrus"
)

func PostRequest(resource string, headers map[string]string, queries map[string]string) *http.Request {
	req, err := http.NewRequest("POST", "https://api.millheat.com/"+resource, nil)
	if err != nil {
		log.Panic(err)
	}

	for key, value := range headers {
		req.Header.Add(key, value)
	}

	query := url.Values{}
	for key, value := range queries {
		query.Set(key, value)
	}
	req.URL.RawQuery = query.Encode()

	return req
}

type ResponseData struct {
	Data       map[string]interface{} `json:"data"`
	ErrorCode  int                    `json:"errorCode"`
	Message    string                 `json:"message"`
	StatusCode int                    `json:"statusCode"`
	Success    bool                   `json:"success"`
}

func DoRequest(req *http.Request) (map[string]interface{}, error) {
	log.Debugf("%v %v", req.Method, req.URL.Path)
	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		log.Panic(err)
	}

	var responseData *ResponseData
	err = json.NewDecoder(resp.Body).Decode(&responseData)
	resp.Body.Close()
	if err != nil {
		log.Panic(err)
	}

	if responseData.ErrorCode == 0 {
		return responseData.Data, nil
	} else {
		return nil, &ApiError{req.URL.Path, responseData.ErrorCode, responseData.Message}
	}
}

type ApiError struct {
	Interface string
	Code      int
	Message   string
}

func (e *ApiError) Error() string {
	return fmt.Sprintf("api error %v: %v when accessing %v", e.Code, e.Message, e.Interface)
}

func AuthorizationCode(accessKey string, secretToken string) (authorizationCode string, err error) {
	headers := map[string]string{"access_key": accessKey, "secret_token": secretToken}
	req := PostRequest("share/applyAuthCode", headers, nil)

	data, err := DoRequest(req)

	if err == nil {
		authorizationCode = data["authorization_code"].(string)
	}
	return
}

func GetAccessAndRefreshToken(authorizationCode string, username string, password string) (accessToken string, refreshToken string, err error) {
	headers := map[string]string{"authorization_code": authorizationCode}
	queries := map[string]string{"username": username, "password": password}
	req := PostRequest("share/applyAccessToken", headers, queries)

	data, err := DoRequest(req)

	if err == nil {
		accessToken = data["access_token"].(string)
		refreshToken = data["refresh_token"].(string)
	}
	return
}

func UpdateAccessTokenAndRefreshToken(oldRefreshToken string) (accessToken string, refreshToken string, err error) {
	queries := map[string]string{"refreshtoken": oldRefreshToken}
	req := PostRequest("share/refreshtoken", nil, queries)

	data, err := DoRequest(req)

	if err == nil {
		accessToken = data["access_token"].(string)
		refreshToken = data["refresh_token"].(string)
	}
	return
}
