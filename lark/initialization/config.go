package initialization

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	LarkAppId                string
	LarkAppSecret            string
	LarkAppEncryptKey        string
	LarkAppVerificationToken string
	LarkBotName              string
	DISCORD_MIDJOURNEY_URL   string
	DISCORD_UPLOAD_URL       string
}

var config *Config

func LoadConfig(cfg string) *Config {
	viper.SetConfigFile(cfg)
	viper.ReadInConfig()
	viper.AutomaticEnv()
	config = &Config{
		LarkAppId:                getViperStringValue("APP_ID"),
		LarkAppSecret:            getViperStringValue("APP_SECRET"),
		LarkAppEncryptKey:        getViperStringValue("APP_ENCRYPT_KEY"),
		LarkAppVerificationToken: getViperStringValue("APP_VERIFICATION_TOKEN"),
		LarkBotName:              getViperStringValue("BOT_NAME"),
		DISCORD_MIDJOURNEY_URL:   getViperStringValue("DISCORD_MIDJOURNEY_URL"),
		DISCORD_UPLOAD_URL:       getViperStringValue("DISCORD_UPLOAD_URL"),
	}
	return config
}

func GetConfig() *Config {
	return config
}

func getViperStringValue(key string) string {
	value := viper.GetString(key)
	if value == "" {
		panic(fmt.Errorf("%s MUST be provided in environment or config.yaml file", key))
	}
	return value
}
