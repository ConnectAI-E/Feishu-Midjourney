package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"lark/chore"
	"lark/db"
	"lark/initialization"
	"lark/services"
	"lark/utils"
	"strings"

	larkcard "github.com/larksuite/oapi-sdk-go/v3/card"

	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
)

type PersonalMessageHandler struct {
	gpt services.ChatGPT
}

var DiscordUpscaleRank = map[string]int64{
	"U1":    1,
	"U2":    2,
	"U3":    3,
	"U4":    4,
	"V1":    1,
	"V2":    2,
	"V3":    3,
	"V4":    4,
	"reset": 0,
}

func (p PersonalMessageHandler) cardHandler(
	_ context.Context,
	cardAction *larkcard.CardAction) (interface{}, error) {

	// if data, err := json.Marshal(cardAction); err == nil {
	// 	fmt.Println("lark card message: ", string(data))
	// }

	actionValue := cardAction.Action.Value
	discordMsgId := actionValue["discordMsgId"].(string)
	value := actionValue["value"].(string)
	index := DiscordUpscaleRank[value]
	msgHash := actionValue["msgHash"].(string)
	larkMsgId := cardAction.OpenMessageID
	redisKey := actionValue["redisKey"].(string)

	larkDiscord := db.GetCache().GetInterface(redisKey)
	if larkDiscord == nil {
		return nil, nil
	}
	var idl IDiscordLarkMap
	if err := json.Unmarshal(larkDiscord, &idl); err != nil {
		fmt.Println("cardHandler: ", "** json.Unmarshal失败 **", err)
		return nil, nil
	}

	if idl.LarkChatId == "" {
		return nil, nil
	}

	if queueErr := discordQueueCheck(idl.LarkChatId); queueErr != nil {
		chore.ReplyMsg(context.Background(), queueErr.Error(), &larkMsgId)
		return nil, nil
	}
	discordQueueSet(idl.LarkChatId)

	eventType := UpscaleEventType(value)

	/******** 保留当前larkMsgId与上一条discordMsgId的映射 ********/
	idl.LarkMsgIdMapPrevDiscordMsgId[discordMsgId] = larkMsgId
	/**
	 * 能使用该字段解决 u操作 回复不带有upscaled by的问题 取决于单个用户一次只能运行一个独立任务
	 */
	idl.From = eventType
	db.GetCache().SetInterface(redisKey, idl)

	var err error
	if eventType == "reset" {
		err = SendDiscordMessageBotReset(discordMsgId, msgHash, larkMsgId)
	} else if eventType == "maxupscale" {
		err = SendDiscordMessageMaxUpscale(discordMsgId, msgHash, larkMsgId)
	} else if eventType == "U" {
		err = SendDiscordMessageBotUpscale(index, discordMsgId, msgHash, larkMsgId)
	} else if eventType == "V" {
		err = SendDiscordMessageBotV(index, discordMsgId, msgHash, larkMsgId)
	}

	/******** 执行err 清除 ********/
	if err != nil && idl.LarkChatId != "" {
		discordQueueDel(idl.LarkChatId)
	}
	return nil, nil
}

func (p PersonalMessageHandler) handle(ctx context.Context, event *larkim.P2MessageReceiveV1) error {

	// larkInfo, _ := json.Marshal(event)
	// fmt.Println("larkInfo", string(larkInfo))

	content := event.Event.Message.Content
	msgId := event.Event.Message.MessageId
	rootId := event.Event.Message.RootId
	chatId := event.Event.Message.ChatId
	eventFlag := db.GetCache().Get(event.EventV2Base.Header.EventID)

	if eventFlag == "1" {
		fmt.Println("eventId", eventFlag, "触发相同的event")
		return nil
	}
	db.GetCache().Set(event.EventV2Base.Header.EventID, "1")

	sessionId := rootId
	if sessionId == nil || *sessionId == "" {
		sessionId = msgId
	}
	if db.GetCache().Get(*msgId) != "" {
		fmt.Println("msgId", *msgId, "processed")
		return nil
	}
	db.GetCache().Set(*msgId, "1")
	qParsed := strings.Trim(parseContent(*content), " ")

	if _, foundClear := utils.EitherTrimEqual(qParsed, "/clearDiscordQueue"); foundClear {
		discordQueueDel(*chatId)
		return nil
	}

	if instruct, foundInstruct := utils.EitherCutPrefix(qParsed,
		"/imagine"); foundInstruct {
		fmt.Println("instruct: ", instruct, "\norigin data: ", qParsed)
		SendDiscordMessageBot(*msgId, instruct, ctx, *chatId)
		return nil
	}

	if _, foundPicture := utils.EitherTrimEqual(qParsed,
		"/midjourney"); foundPicture {
		chore.SendMidjourneyTipCard(context.Background(), msgId)
		return nil
	}
	return nil
}

func NewPersonalMessageHandler(config initialization.Config) MessageHandlerInterface {
	return &PersonalMessageHandler{}
}
