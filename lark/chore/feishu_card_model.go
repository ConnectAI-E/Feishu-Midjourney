package chore

import (
	"encoding/json"
)

type IDiscordCardModel struct {
	Header   IDiscordCardModelHeader    `json:"header"`
	Elements []IDiscordCardModelElement `json:"elements"`
}

type IDiscordCardModelHeader struct {
	Title    IText  `json:"title"`
	Template string `json:"template,omitempty"`
}

type IDiscordCardModelElement struct {
	Tag string `json:"tag"`
	// 图片相关
	ImgKey string   `json:"img_key,omitempty"`
	Alt    ImageAlt `json:"alt,omitempty"`
	// 文字相关
	Text IText `json:"text,omitempty"`
	// 按钮相关
	Actions []IDiscordCardModelAction `json:"actions,omitempty"`
	//
	Elements []IDiscordCardModelElement `json:"elements,omitempty"`
	//
	Content string `json:"content,omitempty"`
}

type ImageAlt struct {
	Tag     string `json:"tag"`
	Content string `json:"content"`
}

type IText = ImageAlt

type IButton = ImageAlt

type IDiscordCardModelAction struct {
	Tag   string  `json:"tag"`
	Text  IButton `json:"text"`
	Type  string  `json:"type"`
	Url   string  `json:"url,omitempty"`
	Value IExtra  `json:"value,omitempty"`
}

type IExtra struct {
	ChatType     string `json:"chatType"`
	Value        string `json:"value"`
	DiscordMsgId string `json:"discordMsgId"`
	RedisKey     string `json:"redisKey"`
	MsgHash      string `json:"msgHash"`
}

func discordCardModel(imgKey string, discordMsgId string, redisKey string, msgHash string) string {
	card := IDiscordCardModel{
		Header: IDiscordCardModelHeader{
			Title: IText{
				Tag:     "plain_text",
				Content: "Midjourney Bot🎉",
			},
		},
		Elements: []IDiscordCardModelElement{
			{
				Tag:    "img",
				ImgKey: imgKey,
				Alt: ImageAlt{
					Tag:     "plain_text",
					Content: "图片",
				},
			},
			{
				Tag: "action",
				Actions: []IDiscordCardModelAction{
					{
						Tag: "button",
						Text: IButton{
							Tag:     "plain_text",
							Content: "U1",
						},
						Type: "primary",
						Value: IExtra{
							ChatType:     "personal",
							Value:        "U1",
							DiscordMsgId: discordMsgId,
							RedisKey:     redisKey,
							MsgHash:      msgHash,
						},
					},
					{
						Tag: "button",
						Text: IButton{
							Tag:     "plain_text",
							Content: "U2",
						},
						Type: "primary",
						Value: IExtra{
							ChatType:     "personal",
							Value:        "U2",
							DiscordMsgId: discordMsgId,
							RedisKey:     redisKey,
							MsgHash:      msgHash,
						},
					},
					{
						Tag: "button",
						Text: IButton{
							Tag:     "plain_text",
							Content: "U3",
						},
						Type: "primary",
						Value: IExtra{
							ChatType:     "personal",
							Value:        "U3",
							DiscordMsgId: discordMsgId,
							RedisKey:     redisKey,
							MsgHash:      msgHash,
						},
					},
					{
						Tag: "button",
						Text: IButton{
							Tag:     "plain_text",
							Content: "U4",
						},
						Type: "primary",
						Value: IExtra{
							ChatType:     "personal",
							Value:        "U4",
							DiscordMsgId: discordMsgId,
							RedisKey:     redisKey,
							MsgHash:      msgHash,
						},
					},
				},
			},
			{
				Tag: "action",
				Actions: []IDiscordCardModelAction{
					{
						Tag: "button",
						Text: IButton{
							Tag:     "plain_text",
							Content: "V1",
						},
						Type: "default",
						Value: IExtra{
							ChatType:     "personal",
							Value:        "V1",
							DiscordMsgId: discordMsgId,
							RedisKey:     redisKey,
							MsgHash:      msgHash,
						},
					},
					{
						Tag: "button",
						Text: IButton{
							Tag:     "plain_text",
							Content: "V2",
						},
						Type: "default",
						Value: IExtra{
							ChatType:     "personal",
							Value:        "V2",
							DiscordMsgId: discordMsgId,
							RedisKey:     redisKey,
							MsgHash:      msgHash,
						},
					},
					{
						Tag: "button",
						Text: IButton{
							Tag:     "plain_text",
							Content: "V3",
						},
						Type: "default",
						Value: IExtra{
							ChatType:     "personal",
							Value:        "V3",
							DiscordMsgId: discordMsgId,
							RedisKey:     redisKey,
							MsgHash:      msgHash,
						},
					},
					{
						Tag: "button",
						Text: IButton{
							Tag:     "plain_text",
							Content: "V4",
						},
						Type: "default",
						Value: IExtra{
							ChatType:     "personal",
							Value:        "V4",
							DiscordMsgId: discordMsgId,
							RedisKey:     redisKey,
							MsgHash:      msgHash,
						},
					},
				},
			},
			{
				Tag: "action",
				Actions: []IDiscordCardModelAction{
					{
						Tag: "button",
						Text: IButton{
							Tag:     "plain_text",
							Content: "re-roll",
						},
						Type: "primary",
						Value: IExtra{
							ChatType:     "personal",
							Value:        "reset",
							DiscordMsgId: discordMsgId,
							RedisKey:     redisKey,
							MsgHash:      msgHash,
						},
					},
				},
			},
		},
	}

	json, _ := json.Marshal(card)

	return string(json)
}

func discordMaxUpscaleCardModel(imgKey string, discordMsgId string, redisKey string, msgHash string) string {
	card := IDiscordCardModel{
		Header: IDiscordCardModelHeader{
			Title: IText{
				Tag:     "plain_text",
				Content: "Midjourney Bot最大升级🎉",
			},
		},
		Elements: []IDiscordCardModelElement{
			{
				Tag:    "img",
				ImgKey: imgKey,
				Alt: ImageAlt{
					Tag:     "plain_text",
					Content: "图片",
				},
			},
			{
				Tag: "action",
				Actions: []IDiscordCardModelAction{
					{
						Tag: "button",
						Text: IButton{
							Tag:     "plain_text",
							Content: "Make Variations",
						},
						Type: "primary",
						Value: IExtra{
							ChatType:     "personal",
							Value:        "maxupscale",
							DiscordMsgId: discordMsgId,
							RedisKey:     redisKey,
							MsgHash:      msgHash,
						},
					},
				},
			},
		},
	}

	json, _ := json.Marshal(card)

	return string(json)
}
