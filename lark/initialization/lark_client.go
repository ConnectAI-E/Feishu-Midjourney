package initialization

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	lark "github.com/larksuite/oapi-sdk-go/v3"
)

var (
	larkClient           *lark.Client
	accessToken          string
	_config              Config
	accessTokenResetTime time.Time
)

type GetAccessToken struct {
	AppId     string `json:"app_id"`
	AppSecret string `json:"app_secret"`
}

type AccessToken struct {
	AppAccessToken    string `json:"app_access_token"`
	Expire            int    `json:"expire"`
	TenantAccessToken string `json:"tenant_access_token"`
	Code              int    `json:"code"`
	Msg               string `json:"msg"`
}

func LoadLarkClient(config Config) {
	larkClient = lark.NewClient(config.LarkAppId, config.LarkAppSecret)
	loadLarkAccessToken(config)
}

func GetLarkClient() *lark.Client {
	return larkClient
}

func loadLarkAccessToken(config Config) {
	_config = config
	requestBody := GetAccessToken{
		AppId:     config.LarkAppId,
		AppSecret: config.LarkAppSecret,
	}
	requestData, err := json.Marshal(requestBody)
	if err != nil {
		return
	}
	req, err := http.NewRequest("POST", "https://open.feishu.cn/open-apis/auth/v3/app_access_token/internal", bytes.NewBuffer(requestData))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	client := &http.Client{}
	resp, _ := client.Do(req)
	bod, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	var tokenObj AccessToken
	_ = json.Unmarshal(bod, &tokenObj)
	accessTokenResetTime = time.Now().Add(time.Minute)
	accessToken = tokenObj.TenantAccessToken
}

func getLarkAccessToken() string {
	if accessToken == "" || time.Now().After(accessTokenResetTime) {
		loadLarkAccessToken(_config)
	}
	return accessToken
}

func GetLarkMsgFile(msgId string, imageKey string) (string, int64, []byte, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("https://open.feishu.cn/open-apis/im/v1/messages/%s/resources/%s?type=image", msgId, imageKey), nil)
	if err != nil {
		return "", 0, nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", getLarkAccessToken()))
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return "", 0, nil, err
	}
	contentType := response.Header.Get("Content-Type")
	var imageType string
	if strings.HasPrefix(contentType, "image/jpeg") {
		imageType = "jpg"
	} else if strings.HasPrefix(contentType, "image/png") {
		imageType = "png"
	}
	imageSize := response.ContentLength
	defer response.Body.Close()
	body, _ := ioutil.ReadAll(response.Body)
	return imageType, imageSize, body, nil
}
