package main

import (
	"lark/handlers"
	"lark/initialization"

	"github.com/gin-gonic/gin"
	"github.com/spf13/pflag"

	sdkginext "github.com/larksuite/oapi-sdk-gin"

	larkcard "github.com/larksuite/oapi-sdk-go/v3/card"
	"github.com/larksuite/oapi-sdk-go/v3/event/dispatcher"
)

func main() {
	cfg := pflag.StringP("config", "c", "./config.yaml", "apiserver config file path.")
	pflag.Parse()

	config := initialization.LoadConfig(*cfg)
	initialization.LoadLarkClient(*config)
	handlers.InitHanders(*config)

	eventHandler := dispatcher.NewEventDispatcher(
		config.LarkAppVerificationToken, config.LarkAppEncryptKey).
		OnP2MessageReceiveV1(handlers.Handler)

	cardHandler := larkcard.NewCardActionHandler(
		config.LarkAppVerificationToken, config.LarkAppEncryptKey,
		handlers.CardHandler())

	r := gin.Default()
	r.POST("/webhook/event",
		sdkginext.NewEventHandlerFunc(eventHandler))
	r.POST("/webhook/card",
		sdkginext.NewCardActionHandlerFunc(
			cardHandler))
	// discord消息回调
	r.POST("/api/discord", handlers.DiscordHandler)

	r.Run(":" + config.LarkPort)
}
