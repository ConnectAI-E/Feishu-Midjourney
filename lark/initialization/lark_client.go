package initialization

import (
	lark "github.com/larksuite/oapi-sdk-go/v3"
)

var (
	larkClient *lark.Client
)

func LoadLarkClient(config Config) {
	larkClient = lark.NewClient(config.LarkAppId, config.LarkAppSecret)
}

func GetLarkClient() *lark.Client {
	return larkClient
}
