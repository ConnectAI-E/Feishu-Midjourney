package handlers

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"lark/chore"
	"lark/db"
	"lark/initialization"
	"lark/services"
	"lark/utils"
	"net/http"
	"strconv"
	"strings"
	"time"

	larkcard "github.com/larksuite/oapi-sdk-go/v3/card"

	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
)

type PersonalMessageHandler struct{}

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

	content := event.Event.Message.Content
	msgId := event.Event.Message.MessageId
	rootId := event.Event.Message.RootId
	chatId := event.Event.Message.ChatId
	eventFlag := db.GetCache().Get(event.EventV2Base.Header.EventID)

	if eventFlag == "1" {
		return nil
	}
	db.GetCache().Set(event.EventV2Base.Header.EventID, "1")

	sessionId := rootId
	if sessionId == nil || *sessionId == "" {
		sessionId = msgId
	}
	if db.GetCache().Get(*msgId) != "" {
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
		SendDiscordMessageBot(*msgId, instruct, ctx, *chatId)
		return nil
	}

	chore.ReplyMsg(ctx, "🤖️：您想进行什么操作？", msgId)
	return nil
}

type RichText struct {
	Title   string              `json:"title"`
	Content [][]RichTextContent `json:"content"`
}

type RichTextContent struct {
	Tag      string `json:"tag"`
	ImageKey string `json:"image_key"`
	Width    int    `json:"width"`
	Height   int    `json:"height"`
}

type ReqUploadFile struct {
	ImgData []byte `json:"imgData"`
	Name    string `json:"name"`
	Size    int64  `json:"size"`
}

func (p PersonalMessageHandler) handleRichText(ctx context.Context, event *larkim.P2MessageReceiveV1) error {
	content := event.Event.Message.Content
	msgId := event.Event.Message.MessageId
	rootId := event.Event.Message.RootId
	chatId := event.Event.Message.ChatId
	eventFlag := db.GetCache().Get(event.EventV2Base.Header.EventID)

	if eventFlag == "1" {
		return nil
	}
	db.GetCache().Set(event.EventV2Base.Header.EventID, "1")

	sessionId := rootId
	if sessionId == nil || *sessionId == "" {
		sessionId = msgId
	}
	if db.GetCache().Get(*msgId) != "" {
		return nil
	}
	db.GetCache().Set(*msgId, "1")
	var data RichText
	err := json.Unmarshal([]byte(*content), &data)
	if err != nil {
		chore.ReplyMsg(ctx, "🤖️：内容解析错误，请检查后重试\n错误信息", msgId)
		return nil
	}
	if len(data.Content) == 0 || len(data.Content[0]) == 0 {
		chore.ReplyMsg(ctx, "🤖️：请上传图片", msgId)
		return nil
	}
	if data.Content[0][0].Tag == "img" {
		if data.Content[0][0].ImageKey == "" {
			chore.ReplyMsg(ctx, "🤖️：请上传图片", msgId)
			return nil
		}
		imageType, size, payload, err := initialization.GetLarkMsgFile(*event.Event.Message.MessageId, data.Content[0][0].ImageKey)
		if err != nil {
			chore.ReplyMsg(ctx, fmt.Sprintf("🤖️：获取上传的图片失败，请重试\n错误信息: %v", err), msgId)
			return nil
		}
		str := *msgId + strconv.FormatInt(time.Now().UnixNano(), 10)
		hash := md5.Sum([]byte(str))
		id := hex.EncodeToString(hash[:])[:12]
		db.GetCache().SetInterface(id, IDiscordLarkMap{
			MsgId:                        *msgId,
			Count:                        0,
			LarkMsgIdMapPrevDiscordMsgId: map[string]string{},
			LarkChatId:                   *chatId,
		})
		requestBody, err := json.Marshal(ReqUploadFile{
			Size:    size,
			Name:    id + "." + imageType,
			ImgData: payload,
		})
		req, err := http.NewRequest("POST", initialization.GetConfig().DISCORD_UPLOAD_URL, bytes.NewBuffer(requestBody))
		if err != nil {
			chore.ReplyMsg(ctx, "🤖️：创建上传图片请求失败，请稍后再试", msgId)
			return nil
		}
		req.Header.Set("Content-Type", "image/jpeg")

		// 发送请求
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			chore.ReplyMsg(ctx, "🤖️：发送上传图片请求失败，请稍后再试", msgId)
			return nil
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			chore.ReplyMsg(ctx, "🤖️：上传图片失败，请重试", msgId)
			return nil
		}
		var files map[string]interface{}
		json.Unmarshal(body, &files)
		err = services.ReqMidjourney(services.RequestTrigger{
			Type:   "describe",
			Prompt: files["name"].(string),
		})
		if err != nil {
			chore.ReplyMsg(ctx, "🤖️：触发describe失败，请重试", msgId)
			return nil
		}
		return nil
	}
	chore.ReplyMsg(ctx, "🤖️：内容错误，请检查后重试", msgId)
	return nil
}

func NewPersonalMessageHandler(config initialization.Config) MessageHandlerInterface {
	return &PersonalMessageHandler{}
}
