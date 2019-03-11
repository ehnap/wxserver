package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	AccessTokenURL = "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s"
)

type AccessTokenInfo struct {
	AccessToken string `json:"access_token,omitempty"`
	ExpiresIn   int    `json:"expires_in,omitempty"`
	ErrCode     int    `json:"errcode,omitempty"`
	ErrMsg      string `json:"errmsg,omitempty"`
	lastTime    time.Time
}

// GetAccessToken 向微信服务器获取AccessToken
func GetAccessToken(appid string, appsecret string) string {
	tokenInfo := AccessTokenInfo{}
	url := fmt.Sprintf(AccessTokenURL, appid, appsecret)
	response, err := http.Get(url)

	if err != nil {
		return ""
	}

	result, responseError := ioutil.ReadAll(response.Body)
	if responseError != nil {
		return ""
	}

	resultJSON := AccessTokenInfo{}
	jsonError := json.Unmarshal(result, &resultJSON)
	tokenInfo = resultJSON
	tokenInfo.lastTime = time.Now()
	if jsonError != nil {
		return ""
	}

	if resultJSON.ErrCode == 0 && resultJSON.AccessToken != "" {
		return tokenInfo.AccessToken
	}
	return ""
}
