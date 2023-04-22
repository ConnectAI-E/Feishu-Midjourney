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
				Content: "图片生成完成🎉",
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
			{
				Tag: "hr",
			},
			{
				Tag: "note",
				Elements: []IDiscordCardModelElement{
					{
						Tag:     "plain_text",
						Content: "* U按钮可以对图像进行升级，生成一个更大的版本，并添加更多的细节。（无法继续升级）\n* V按钮可以创建所选网格图像的轻微变化。创建一个变化会生成一个新的图像网格，与所选图像的整体风格和构图相似。\n* re-roll会重新运行一个任务。在这种情况下，它会重新运行原始提示，生成一个新的图像网格。",
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
				Content: "图片生成完成🎉",
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

func midjourneyTips() string {
	card := IDiscordCardModel{
		Header: IDiscordCardModelHeader{
			Title: IText{
				Tag:     "plain_text",
				Content: "Midjourney Tips",
			},
			Template: "purple",
		},
		Elements: []IDiscordCardModelElement{
			{
				Tag:    "img",
				ImgKey: "img_v2_760a315a-85d8-455b-bc0c-4c6b0edcc53g",
				Alt: ImageAlt{
					Tag:     "plain_text",
					Content: "图片",
				},
			},
			{
				Tag: "hr",
			},
			{
				Tag: "note",
				Elements: []IDiscordCardModelElement{
					{
						Tag:     "plain_text",
						Content: "Image Prompts: 图片URL可以添加到提示中，以影响生成结果的风格和内容。图片URL始终位于提示的开头。\nText Prompt: 对所需图像的文本描述。撰写良好的提示有助于生成令人惊叹的图像。\nParameters:  参数可以改变图像的生成方式。参数可以改变纵横比、模型、放大器等很多内容。参数位于提示的末尾。",
					},
				},
			},
			{
				Tag: "hr",
			},
			{
				Tag: "note",
				Elements: []IDiscordCardModelElement{
					{
						Tag:     "plain_text",
						Content: "Parameters: \n1. 默认采用的是version 4版本模型，如果需要使用其它版本，使用--v控制。例如：--v 5 \n2. 如果想控制生成图片的比例，可以使用--ar来控制。例如：--ar 16:9",
					},
				},
			},
			{
				Tag: "hr",
			},
			{
				Tag: "note",
				Elements: []IDiscordCardModelElement{
					{
						Tag:     "plain_text",
						Content: "更多版本参考：https://docs.midjourney.com/docs/models \n更多参数参考：https://docs.midjourney.com/docs/parameter-list",
					},
				},
			},
		},
	}
	json, _ := json.Marshal(card)

	return string(json)
}
