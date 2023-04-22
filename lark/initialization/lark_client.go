package initialization

import (
	"time"

	lark "github.com/larksuite/oapi-sdk-go/v3"
)

var (
	larkClient           *lark.Client
	accessToken          string
	_config              Config
	accessTokenResetTime time.Time
)

func LoadLarkClient(config Config) {
	larkClient = lark.NewClient(config.LarkAppId, config.LarkAppSecret)
}

func GetLarkClient() *lark.Client {
	return larkClient
}
