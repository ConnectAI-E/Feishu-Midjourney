package chore

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"lark/initialization"
	"net/http"

	"github.com/google/uuid"
	larkcard "github.com/larksuite/oapi-sdk-go/v3/card"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
)

func ReplyCard(ctx context.Context,
	msgId *string,
	cardContent string,
) error {
	client := initialization.GetLarkClient()
	resp, err := client.Im.Message.Reply(ctx, larkim.NewReplyMessageReqBuilder().
		MessageId(*msgId).
		Body(larkim.NewReplyMessageReqBodyBuilder().
			MsgType(larkim.MsgTypeInteractive).
			Uuid(uuid.New().String()).
			Content(cardContent).
			Build()).
		Build())

	// 处理错误
	if err != nil {
		fmt.Println(err)
		return err
	}

	// 服务端错误处理
	if !resp.Success() {
		fmt.Println(resp.Code, resp.Msg, resp.RequestId())
		return err
	}
	return nil
}

func NewSendCard(
	header *larkcard.MessageCardHeader,
	elements ...larkcard.MessageCardElement) (string,
	error) {
	config := larkcard.NewMessageCardConfig().
		WideScreenMode(false).
		EnableForward(true).
		UpdateMulti(false).
		Build()
	var aElementPool []larkcard.MessageCardElement
	// for _, element := range elements {
	// 	aElementPool = append(aElementPool, element)
	// }
	aElementPool = append(aElementPool, elements...)
	// 卡片消息体
	cardContent, err := larkcard.NewMessageCard().
		Config(config).
		Header(header).
		Elements(
			aElementPool,
		).
		String()
	return cardContent, err
}

// withHeader 用于生成消息头
func WithHeader(title string, color string) *larkcard.
	MessageCardHeader {
	if title == "" {
		title = "🤖️机器人提醒"
	}
	header := larkcard.NewMessageCardHeader().
		Template(color).
		Title(larkcard.NewMessageCardPlainText().
			Content(title).
			Build()).
		Build()
	return header
}

// withNote 用于生成纯文本脚注
func WithNote(note string) larkcard.MessageCardElement {
	noteElement := larkcard.NewMessageCardNote().
		Elements([]larkcard.MessageCardNoteElement{larkcard.NewMessageCardPlainText().
			Content(note).
			Build()}).
		Build()
	return noteElement
}

// withMainText 用于生成纯文本消息体
func WithMainText(msg string) larkcard.MessageCardElement {
	msg, i := processMessage(msg)
	msg = cleanTextBlock(msg)
	if i != nil {
		return nil
	}
	mainElement := larkcard.NewMessageCardDiv().
		Fields([]*larkcard.MessageCardField{larkcard.NewMessageCardField().
			Text(larkcard.NewMessageCardPlainText().
				Content(msg).
				Build()).
			IsShort(false).
			Build()}).
		Build()
	return mainElement
}

func ReplyMsg(ctx context.Context, msg string, msgId *string) error {
	msg, i := processMessage(msg)
	if i != nil {
		return i
	}
	client := initialization.GetLarkClient()
	content := larkim.NewTextMsgBuilder().
		Text(msg).
		Build()

	resp, err := client.Im.Message.Reply(ctx, larkim.NewReplyMessageReqBuilder().
		MessageId(*msgId).
		Body(larkim.NewReplyMessageReqBodyBuilder().
			MsgType(larkim.MsgTypeText).
			Uuid(uuid.New().String()).
			Content(content).
			Build()).
		Build())

	// 处理错误
	if err != nil {
		fmt.Println(err)
		return err
	}

	// 服务端错误处理
	if !resp.Success() {
		fmt.Println(resp.Code, resp.Msg, resp.RequestId())
		return err
	}
	return nil
}

func SendMsg(ctx context.Context, msg string, chatId *string) error {
	//fmt.Println("sendMsg", msg, chatId)
	msg, i := processMessage(msg)
	if i != nil {
		return i
	}
	client := initialization.GetLarkClient()
	content := larkim.NewTextMsgBuilder().
		Text(msg).
		Build()

	//fmt.Println("content", content)

	resp, err := client.Im.Message.Create(ctx, larkim.NewCreateMessageReqBuilder().
		ReceiveIdType(larkim.ReceiveIdTypeChatId).
		Body(larkim.NewCreateMessageReqBodyBuilder().
			MsgType(larkim.MsgTypeText).
			ReceiveId(*chatId).
			Content(content).
			Build()).
		Build())

	// 处理错误
	if err != nil {
		fmt.Println(err)
		return err
	}

	// 服务端错误处理
	if !resp.Success() {
		fmt.Println(resp.Code, resp.Msg, resp.RequestId())
		return err
	}
	return nil
}

func SendNewTopicCard(ctx context.Context,
	sessionId *string, msgId *string, content string) {
	newCard, _ := NewSendCard(
		WithHeader("👻️ 已开启新的话题（点击⬆️文字进入话题）", larkcard.TemplateBlue),
		WithMainText(content),
		WithNote("提醒：在对话框参与回复，可保持话题连贯"))
	ReplyCard(
		ctx,
		msgId,
		newCard,
	)
}

func UploadImage(url string) (*string, error) {
	res, err := http.Get(url)
	if err != nil {
		fmt.Println("获取图片资源失败", err)
		return nil, err
	}
	defer res.Body.Close()
	imagesBytes, _ := ioutil.ReadAll(res.Body)

	client := initialization.GetLarkClient()
	resp, err := client.Im.Image.Create(context.Background(),
		larkim.NewCreateImageReqBuilder().
			Body(larkim.NewCreateImageReqBodyBuilder().
				ImageType(larkim.ImageTypeMessage).
				Image(bytes.NewReader(imagesBytes)).
				Build()).
			Build())

	// 处理错误
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	// 服务端错误处理
	if !resp.Success() {
		fmt.Println(resp.Code, resp.Msg, resp.RequestId())
		return nil, err
	}
	return resp.Data.ImageKey, nil
}

func ReplyImage(ctx context.Context, ImageKey *string,
	msgId *string) error {
	fmt.Println("sendMsg", ImageKey, msgId)

	msgImage := larkim.MessageImage{ImageKey: *ImageKey}
	content, err := msgImage.String()
	if err != nil {
		fmt.Println(err)
		return err
	}
	client := initialization.GetLarkClient()

	resp, err := client.Im.Message.Reply(ctx, larkim.NewReplyMessageReqBuilder().
		MessageId(*msgId).
		Body(larkim.NewReplyMessageReqBodyBuilder().
			MsgType(larkim.MsgTypeImage).
			Uuid(uuid.New().String()).
			Content(content).
			Build()).
		Build())

	// 处理错误
	if err != nil {
		fmt.Println(err)
		return err
	}

	// 服务端错误处理
	if !resp.Success() {
		fmt.Println(resp.Code, resp.Msg, resp.RequestId())
		return err
	}
	return nil

}

func ReplayImageByImages(ctx context.Context, url string, msgId *string) error {
	imageKey, err := UploadImage(url)
	if err != nil {
		return err
	}
	err = ReplyImage(ctx, imageKey, msgId)
	if err != nil {
		return err
	}
	return nil
}
func replayImageByImagesDiscord(url string, redisKey string, discordMessageId string, msgHash string, msgId string, isUpscaled bool) error {
	imageKey, err := UploadImage(url)
	if err != nil {
		return err
	}
	var card string
	if isUpscaled {
		card = discordMaxUpscaleCardModel(*imageKey, discordMessageId, redisKey, msgHash)
	} else {
		card = discordCardModel(*imageKey, discordMessageId, redisKey, msgHash)
	}
	err = ReplyCard(context.Background(), &msgId, card)
	if err != nil {
		return err
	}
	return nil
}

func SendPicCreateInstructionCard(ctx context.Context,
	sessionId *string, msgId *string, content string) {
	newCard, _ := NewSendCard(
		WithHeader("🖼️  已进入图片创作模式", larkcard.TemplateBlue),
		WithNote("请注意，这将开始一个全新的对话，您将无法利用之前话题的历史信息"))
	ReplyCard(
		ctx,
		msgId,
		newCard,
	)
}

func ReplayImageByImagesDiscord(url string, redisKey string, discordMessageId string, msgHash string, msgId string, isUpscaled bool) error {
	imageKey, err := UploadImage(url)
	if err != nil {
		return err
	}
	var card string
	if isUpscaled {
		card = discordMaxUpscaleCardModel(*imageKey, discordMessageId, redisKey, msgHash)
	} else {
		card = discordCardModel(*imageKey, discordMessageId, redisKey, msgHash)
	}
	err = ReplyCard(context.Background(), &msgId, card)
	if err != nil {
		return err
	}
	return nil
}
