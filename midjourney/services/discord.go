package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	config "midjourney/initialization"
	"net/http"
)

const url = "https://discord.com/api/v9/interactions"

func GenerateImage(prompt string) error {
	requestBody := ReqTriggerDiscord{
		Type:          2,
		GuildID:       config.GetConfig().DISCORD_SERVER_ID,
		ChannelID:     config.GetConfig().DISCORD_CHANNEL_ID,
		ApplicationId: "936929561302675456",
		SessionId:     "cb06f61453064c0983f2adae2a88c223",
		Data: DSCommand{
			Version: "1077969938624553050",
			Id:      "938956540159881230",
			Name:    "imagine",
			Type:    1,
			Options: []DSOption{{Type: 3, Name: "prompt", Value: prompt}},
			ApplicationCommand: DSApplicationCommand{
				Id:                       "938956540159881230",
				ApplicationId:            "936929561302675456",
				Version:                  "1077969938624553050",
				DefaultPermission:        true,
				DefaultMemberPermissions: nil,
				Type:                     1,
				Nsfw:                     false,
				Name:                     "imagine",
				Description:              "Lucky you!",
				DmPermission:             true,
				Options:                  []DSCommandOption{{Type: 3, Name: "prompt", Description: "The prompt to imagine", Required: true}},
			},
			Attachments: []interface{}{},
		},
	}
	err := request(requestBody)
	return err
}

func Upscale(index int64, messageId string, messageHash string) error {
	requestBody := ReqUpscaleDiscord{
		Type:          3,
		GuildId:       config.GetConfig().DISCORD_SERVER_ID,
		ChannelId:     config.GetConfig().DISCORD_CHANNEL_ID,
		MessageFlags:  0,
		MessageId:     messageId,
		ApplicationId: "936929561302675456",
		SessionId:     "45bc04dd4da37141a5f73dfbfaf5bdcf",
		Data: UpscaleData{
			ComponentType: 2,
			CustomId:      fmt.Sprintf("MJ::JOB::upsample::%d::%s", index, messageHash),
		},
	}
	err := request(requestBody)
	return err
}

func MaxUpscale(messageId string, messageHash string) error {
	requestBody := ReqUpscaleDiscord{
		Type:          3,
		GuildId:       config.GetConfig().DISCORD_SERVER_ID,
		ChannelId:     config.GetConfig().DISCORD_CHANNEL_ID,
		MessageFlags:  0,
		MessageId:     messageId,
		ApplicationId: "936929561302675456",
		SessionId:     "1f3dbdf09efdf93d81a3a6420882c92c",
		Data: UpscaleData{
			ComponentType: 2,
			CustomId:      fmt.Sprintf("MJ::JOB::variation::1::%s::SOLO", messageHash),
		},
	}

	data, _ := json.Marshal(requestBody)

	fmt.Println("max upscale request body: ", string(data))

	err := request(requestBody)
	return err
}

func Variate(index int64, messageId string, messageHash string) error {
	requestBody := ReqVariationDiscord{
		Type:          3,
		GuildId:       config.GetConfig().DISCORD_SERVER_ID,
		ChannelId:     config.GetConfig().DISCORD_CHANNEL_ID,
		MessageFlags:  0,
		MessageId:     messageId,
		ApplicationId: "936929561302675456",
		SessionId:     "45bc04dd4da37141a5f73dfbfaf5bdcf",
		Data: UpscaleData{
			ComponentType: 2,
			CustomId:      fmt.Sprintf("MJ::JOB::variation::%d::%s", index, messageHash),
		},
	}
	err := request(requestBody)
	return err
}

func Reset(messageId string, messageHash string) error {
	requestBody := ReqResetDiscord{
		Type:          3,
		GuildId:       config.GetConfig().DISCORD_SERVER_ID,
		ChannelId:     config.GetConfig().DISCORD_CHANNEL_ID,
		MessageFlags:  0,
		MessageId:     messageId,
		ApplicationId: "936929561302675456",
		SessionId:     "45bc04dd4da37141a5f73dfbfaf5bdcf",
		Data: UpscaleData{
			ComponentType: 2,
			CustomId:      fmt.Sprintf("MJ::JOB::reroll::0::%s::SOLO", messageHash),
		},
	}
	err := request(requestBody)
	return err
}

func request(params interface{}) error {
	requestData, err := json.Marshal(params)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestData))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", config.GetConfig().DISCORD_USER_TOKEN)
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	bod, respErr := ioutil.ReadAll(response.Body)
	fmt.Println("upscale response: ", string(bod), respErr, response.Status)
	return respErr
}
