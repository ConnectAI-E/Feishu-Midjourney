package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	config "midjourney/initialization"
	"net/http"
	"path/filepath"
)

const (
	url             = "https://discord.com/api/v9/interactions"
	uploadUrlFormat = "https://discord.com/api/v9/channels/%s/attachments"
	appId           = "936929561302675456"
)

func GenerateImage(prompt string) error {
	requestBody := ReqTriggerDiscord{
		Type:          2,
		GuildID:       config.GetConfig().DISCORD_SERVER_ID,
		ChannelID:     config.GetConfig().DISCORD_CHANNEL_ID,
		ApplicationId: appId,
		SessionId:     "cb06f61453064c0983f2adae2a88c223",
		Data: DSCommand{
			Version: "1118961510123847772",
			Id:      "938956540159881230",
			Name:    "imagine",
			Type:    1,
			Options: []DSOption{{Type: 3, Name: "prompt", Value: prompt}},
			ApplicationCommand: DSApplicationCommand{
				Id:                       "938956540159881230",
				ApplicationId:            "936929561302675456",
				Version:                  "1118961510123847772",
				DefaultPermission:        true,
				DefaultMemberPermissions: nil,
				Type:                     1,
				Nsfw:                     false,
				Name:                     "imagine",
				Description:              "Lucky you!",
				DmPermission:             true,
				Options:                  []DSCommandOption{{Type: 3, Name: "prompt", Description: "The prompt to imagine", Required: true}},
			},
			Attachments: []ReqCommandAttachments{},
		},
	}
	_, err := request(requestBody, url)
	return err
}

func Upscale(index int64, messageId string, messageHash string) error {
	requestBody := ReqUpscaleDiscord{
		Type:          3,
		GuildId:       config.GetConfig().DISCORD_SERVER_ID,
		ChannelId:     config.GetConfig().DISCORD_CHANNEL_ID,
		MessageFlags:  0,
		MessageId:     messageId,
		ApplicationId: appId,
		SessionId:     "45bc04dd4da37141a5f73dfbfaf5bdcf",
		Data: UpscaleData{
			ComponentType: 2,
			CustomId:      fmt.Sprintf("MJ::JOB::upsample::%d::%s", index, messageHash),
		},
	}
	_, err := request(requestBody, url)
	return err
}

func MaxUpscale(messageId string, messageHash string) error {
	requestBody := ReqUpscaleDiscord{
		Type:          3,
		GuildId:       config.GetConfig().DISCORD_SERVER_ID,
		ChannelId:     config.GetConfig().DISCORD_CHANNEL_ID,
		MessageFlags:  0,
		MessageId:     messageId,
		ApplicationId: appId,
		SessionId:     "1f3dbdf09efdf93d81a3a6420882c92c",
		Data: UpscaleData{
			ComponentType: 2,
			CustomId:      fmt.Sprintf("MJ::JOB::variation::1::%s::SOLO", messageHash),
		},
	}

	data, _ := json.Marshal(requestBody)

	fmt.Println("max upscale request body: ", string(data))

	_, err := request(requestBody, url)
	return err
}

func Variate(index int64, messageId string, messageHash string) error {
	requestBody := ReqVariationDiscord{
		Type:          3,
		GuildId:       config.GetConfig().DISCORD_SERVER_ID,
		ChannelId:     config.GetConfig().DISCORD_CHANNEL_ID,
		MessageFlags:  0,
		MessageId:     messageId,
		ApplicationId: appId,
		SessionId:     "45bc04dd4da37141a5f73dfbfaf5bdcf",
		Data: UpscaleData{
			ComponentType: 2,
			CustomId:      fmt.Sprintf("MJ::JOB::variation::%d::%s", index, messageHash),
		},
	}
	_, err := request(requestBody, url)
	return err
}

func Reset(messageId string, messageHash string) error {
	requestBody := ReqResetDiscord{
		Type:          3,
		GuildId:       config.GetConfig().DISCORD_SERVER_ID,
		ChannelId:     config.GetConfig().DISCORD_CHANNEL_ID,
		MessageFlags:  0,
		MessageId:     messageId,
		ApplicationId: appId,
		SessionId:     "45bc04dd4da37141a5f73dfbfaf5bdcf",
		Data: UpscaleData{
			ComponentType: 2,
			CustomId:      fmt.Sprintf("MJ::JOB::reroll::0::%s::SOLO", messageHash),
		},
	}
	_, err := request(requestBody, url)
	return err
}

func Describe(uploadName string) error {
	requestBody := ReqTriggerDiscord{
		Type:          2,
		GuildID:       config.GetConfig().DISCORD_SERVER_ID,
		ChannelID:     config.GetConfig().DISCORD_CHANNEL_ID,
		ApplicationId: "936929561302675456",
		SessionId:     "0033db636f7ce1a951e54cdac7044de3",
		Data: DSCommand{
			Version: "1118961510123847774",
			Id:      "1092492867185950852",
			Name:    "describe",
			Type:    1,
			Options: []DSOption{{Type: 11, Name: "image", Value: 0}},
			ApplicationCommand: DSApplicationCommand{
				Id:                       "1092492867185950852",
				ApplicationId:            "936929561302675456",
				Version:                  "1118961510123847774",
				DefaultPermission:        true,
				DefaultMemberPermissions: nil,
				Type:                     1,
				Nsfw:                     false,
				Name:                     "describe",
				Description:              "Writes a prompt based on your image.",
				DmPermission:             true,
				Options:                  []DSCommandOption{{Type: 11, Name: "image", Description: "The image to describe", Required: true}},
			},
			Attachments: []ReqCommandAttachments{{
				Id:             "0",
				Filename:       filepath.Base(uploadName),
				UploadFilename: uploadName,
			}},
		},
	}
	_, err := request(requestBody, url)
	return err
}

func Attachments(name string, size int64) (ResAttachments, error) {
	requestBody := ReqAttachments{
		Files: []ReqFile{{
			Filename: name,
			FileSize: size,
			Id:       "1",
		}},
	}
	uploadUrl := fmt.Sprintf(uploadUrlFormat, config.GetConfig().DISCORD_CHANNEL_ID)
	body, err := request(requestBody, uploadUrl)
	var data ResAttachments
	json.Unmarshal(body, &data)
	return data, err
}

func request(params interface{}, url string) ([]byte, error) {
	requestData, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", config.GetConfig().DISCORD_USER_TOKEN)
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	bod, respErr := ioutil.ReadAll(response.Body)
	fmt.Println("response:", string(bod), respErr, response.Status, url)
	return bod, respErr
}
