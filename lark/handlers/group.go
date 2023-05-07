package handlers

import (
	"context"
	"lark/chore"
	"lark/db"
	"lark/initialization"
	"lark/utils"
	"strings"

	larkcard "github.com/larksuite/oapi-sdk-go/v3/card"

	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
)

type GroupMessageHandler struct {
	config initialization.Config
}

func (p GroupMessageHandler) cardHandler(_ context.Context,
	cardAction *larkcard.CardAction) (interface{}, error) {
	return nil, nil
}

func (p GroupMessageHandler) handle(ctx context.Context, event *larkim.P2MessageReceiveV1) error {
	ifMention := p.judgeIfMentionMe(event)
	if !ifMention {
		return nil
	}
	content := event.Event.Message.Content
	msgId := event.Event.Message.MessageId
	rootId := event.Event.Message.RootId
	// chatId := event.Event.Message.ChatId
	sessionId := rootId
	eventFlag := db.GetCache().Get(event.EventV2Base.Header.EventID)

	if eventFlag == "1" {
		return nil
	}
	db.GetCache().Set(event.EventV2Base.Header.EventID, "1")

	if sessionId == nil || *sessionId == "" {
		sessionId = msgId
	}

	if db.GetCache().Get(*msgId) != "" {
		return nil
	}
	db.GetCache().Set(*msgId, "1")
	qParsed := strings.Trim(parseContent(*content), " ")

	if _, foundInstruct := utils.EitherCutPrefix(qParsed,
		"/imagine"); foundInstruct {
		chore.ReplyMsg(ctx, "群聊暂不支持，请私聊进行", msgId)
		return nil
	}
	return nil
}

func (p GroupMessageHandler) handleRichText(ctx context.Context, event *larkim.P2MessageReceiveV1) error {
	return nil
}

var _ MessageHandlerInterface = (*GroupMessageHandler)(nil)

func NewGroupMessageHandler(config initialization.Config) MessageHandlerInterface {
	return &GroupMessageHandler{
		config: config,
	}
}

func (p GroupMessageHandler) judgeIfMentionMe(event *larkim.P2MessageReceiveV1) bool {
	mention := event.Event.Message.Mentions
	if len(mention) != 1 {
		return false
	}
	return *mention[0].Name == p.config.LarkBotName
}
