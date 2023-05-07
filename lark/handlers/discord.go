package handlers

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"lark/chore"
	"lark/services"
	"net/http"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"lark/db"

	discord "github.com/bwmarrin/discordgo"
	"github.com/gin-gonic/gin"
)

type IDiscordLarkMap struct {
	MsgId                        string            `json:"msgId"`
	Count                        int64             `json:"count"`
	LarkMsgIdMapPrevDiscordMsgId map[string]string `json:"larkMsgIdMapPrevDiscordMsgId"`
	LarkChatId                   string            `json:"larkChatId"`
	From                         string            `json:"from"`
}

const (
	DiscordPrefix   = "<<<!"
	DiscordNextFix  = "!>>>"
	DiscordQueueKey = "**_discord_queue_**"
	DiscordLockKey  = "**_discord_lock_**"
)

type ReqCb struct {
	Embeds  []*discord.MessageEmbed `json:"embeds,omitempty"`
	Discord *discord.MessageCreate  `json:"discord,omitempty"`
	Content string                  `json:"content,omitempty"`
	Type    Scene                   `json:"type"`
}

type Scene string

const (
	FirstTrigger      Scene = "FirstTrigger"
	GenerateEnd       Scene = "GenerateEnd"
	GenerateEditError Scene = "GenerateEditError"
	Rich              Scene = "RichText"
)

func DiscordHandler(c *gin.Context) {
	var params ReqCb
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if params.Type == FirstTrigger {
		re := regexp.MustCompile(`<<<!([^!]+)!>>>`)
		match := re.FindStringSubmatch(params.Content)
		if len(match) > 0 {
			id := match[1]
			discordIteratorTag(id)
		}
		return
	}

	if params.Type == GenerateEnd {
		if id, notFound := getDiscordLardMapId(params.Discord.Content); notFound == nil {
			msgHash := generateDiscordMsgHash(params.Discord.Attachments[0].URL)
			var referenceMsgId string
			if params.Discord.MessageReference != nil {
				referenceMsgId = params.Discord.MessageReference.MessageID
			}
			discordTriggerReplayLark(
				params.Discord.Attachments[0].URL,
				params.Discord.Message.ID,
				id,
				msgHash,
				referenceMsgId,
			)
		}
		return
	}

	if params.Type == GenerateEditError {
		if id, notFound := getDiscordLardMapId(params.Content); notFound == nil {
			if idl, err := getDiscordLarkMapJson(id); err == nil {
				instructException(id, idl.LarkChatId, idl.MsgId)
			}
		}
		return
	}

	if params.Type == Rich {
		embeds := params.Embeds
		if len(embeds) > 0 {
			embed := embeds[0]
			if embed != nil && embed.Image != nil {
				if embed.Image.URL != "" {
					filename := path.Base(embed.Image.URL)
					id := strings.TrimSuffix(filename, filepath.Ext(filename))
					if data, err := getDiscordLarkMapJson(id); err == nil {
						chore.ReplyMsg(context.Background(), embed.Description, &data.MsgId)
					}
				}
			}
		}
	}
}

func SendDiscordMessageBot(msgId string, content string, ctx context.Context, larkChatId string) {
	err := discordQueueCheck(larkChatId)
	if err != nil {
		chore.ReplyMsg(ctx, err.Error(), &msgId)
		return
	}
	str := msgId + strconv.FormatInt(time.Now().UnixNano(), 10)
	hash := md5.Sum([]byte(str))
	id := hex.EncodeToString(hash[:])[:12]
	db.GetCache().SetInterface(id, IDiscordLarkMap{
		MsgId:                        msgId,
		Count:                        0,
		LarkMsgIdMapPrevDiscordMsgId: map[string]string{},
		LarkChatId:                   larkChatId,
	})
	err = services.ReqMidjourney(services.RequestTrigger{
		Type:   "generate",
		Prompt: DiscordPrefix + id + DiscordNextFix + content,
	})
	if err != nil {
		db.GetCache().Clear(id)
		chore.ReplyMsg(ctx, fmt.Sprintf("🤖️：图片生成失败，请稍后再试～\n错误信息: %v", err), &msgId)
		return
	}

	// 生成中
	discordQueueSet(larkChatId)

	// 生成中回复
	chore.ReplyMsg(context.Background(), "🤖️ ：正在生成中，请稍等......", &msgId)

	/******* 指令错误或排队中都会触发 *******/
	checkSendDiscordMessage(make(chan struct{}), id)
}

func SendDiscordMessageBotUpscale(
	index int64,
	discordMessageId string,
	msgHash string,
	larkMsgId string,
) error {
	/******* 处理同一张图片 点击同一个事件 start *******/
	flagStr := db.GetCache().Get(discordMessageId)
	if strings.Contains(flagStr, fmt.Sprint("U", index)) {
		chore.ReplyMsg(
			context.Background(),
			fmt.Sprintf("🤖️：您已经给该照片升级过: %v", fmt.Sprint("U", index)),
			&larkMsgId,
		)
		return errors.New("已经升级过")
	}
	db.GetCache().Set(discordMessageId, flagStr+fmt.Sprint("U", index))
	/******* end *******/

	err := services.ReqMidjourney(services.RequestTrigger{
		Type:         "upscale",
		DiscordMsgId: discordMessageId,
		MsgHash:      msgHash,
		Index:        index,
	})
	if err != nil {
		chore.ReplyMsg(
			context.Background(),
			fmt.Sprintf("🤖️：图片升级失败，请稍后再试～\n错误信息: %v", err),
			&larkMsgId,
		)
		return err
	}
	chore.ReplyMsg(
		context.Background(),
		fmt.Sprintf("🤖️：图片正在进行%v操作升级，请稍等......", fmt.Sprint("U", index)),
		&larkMsgId,
	)
	return nil
}

func SendDiscordMessageBotV(index int64, discordMessageId string, msgHash string, larkMsgId string) error {
	err := services.ReqMidjourney(services.RequestTrigger{
		Type:         "variation",
		DiscordMsgId: discordMessageId,
		MsgHash:      msgHash,
		Index:        index,
	})
	if err != nil {
		chore.ReplyMsg(
			context.Background(),
			fmt.Sprintf("🤖️：图片操作失败，请稍后再试~\n错误信息: %v", err),
			&larkMsgId,
		)
		return err
	}
	chore.ReplyMsg(
		context.Background(),
		fmt.Sprintf("🤖️：图片正在进行%v操作，请稍等......",
			fmt.Sprint("V", index)),
		&larkMsgId,
	)
	return nil
}

func SendDiscordMessageMaxUpscale(discordMessageId string, msgHash string, larkMsgId string) error {
	err := services.ReqMidjourney(services.RequestTrigger{
		Type:         "maxUpscale",
		DiscordMsgId: discordMessageId,
		MsgHash:      msgHash,
	})
	if err != nil {
		chore.ReplyMsg(
			context.Background(),
			fmt.Sprintf("🤖️：图片升级失败，请稍后再试～\n错误信息: %v", err),
			&larkMsgId,
		)
		return err
	}
	chore.ReplyMsg(
		context.Background(),
		"🤖️：图片正在进行最终升级，请稍等......",
		&larkMsgId,
	)
	return nil
}

func SendDiscordMessageBotReset(discordMessageId string, msgHash string, larkMsgId string) error {
	err := services.ReqMidjourney(services.RequestTrigger{
		Type:         "reset",
		DiscordMsgId: discordMessageId,
		MsgHash:      msgHash,
	})
	if err != nil {
		chore.ReplyMsg(
			context.Background(),
			fmt.Sprintf("🤖️：图片重新生成失败，请稍后再试~\n错误信息: %v", err),
			&larkMsgId,
		)
		return err
	}
	chore.ReplyMsg(
		context.Background(),
		"🤖️：图片重新生成中，请稍等......",
		&larkMsgId,
	)
	return nil
}

func checkSendDiscordMessage(done chan struct{}, id string) {
	ticker := time.NewTicker(time.Second * 15)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if idl, err := getDiscordLarkMapJson(id); err == nil {
				if idl.Count > 1 {
					fmt.Println("指令异常", id, idl.Count)
					instructException(id, idl.LarkChatId, idl.MsgId)
					close(done)
					return
				} else if idl.Count == -1 {
					close(done)
					return
				} else {
					fmt.Println("检查指令中", id, idl.Count)
					idl.Count++
					db.GetCache().SetInterface(id, idl)
				}
			}
		case <-done:
			return
		}
	}
}

func UpscaleEventType(str string) string {
	if str == "reset" {
		return "reset"
	} else if strings.Contains(str, "maxupscale") {
		return "maxupscale"
	} else if strings.Contains(str, "U") {
		return "U"
	} else if strings.Contains(str, "V") {
		return "V"
	} else {
		return ""
	}
}

func instructException(id string, larkChatId string, msgId string) {
	// 不一定是指令异常 也有可能是任务过多阻塞了
	chore.ReplyMsg(context.Background(), "🤖️ ：发送的指令存在异常，请检查后重试", &msgId)
	discordQueueDel(larkChatId)
	db.GetCache().Clear(id)
}

func getDiscordLardMapId(content string) (id string, err error) {
	re := regexp.MustCompile(`<<<!([^!]+)!>>>`)
	match := re.FindStringSubmatch(content)
	if len(match) > 0 {
		id := match[1]
		return id, nil
	}
	return "", errors.New("not found")
}

func getDiscordLarkMapJson(id string) (IDiscordLarkMap, error) {
	discordLark := db.GetCache().GetInterface(id)
	if discordLark == nil {
		return IDiscordLarkMap{}, errors.New("not found")
	}

	var idl IDiscordLarkMap
	if err := json.Unmarshal(discordLark, &idl); err != nil {
		return IDiscordLarkMap{}, errors.New("not found")
	}
	return idl, nil
}

func discordTriggerReplayLark(
	url string,
	discordMsgId string,
	key string,
	msgHash string,
	referenceMsgId string,
) {
	idl, err := getDiscordLarkMapJson(key)
	if err != nil {
		return
	}

	var msgId = idl.MsgId
	if referenceMsgId != "" {
		msgId = idl.LarkMsgIdMapPrevDiscordMsgId[referenceMsgId]
	}
	discordQueueDel(idl.LarkChatId)
	chore.ReplayImageByImagesDiscord(
		url,
		key,
		discordMsgId,
		msgHash,
		msgId,
		idl.From == "U",
	)
}

func discordIteratorTag(key string) {
	idl, err := getDiscordLarkMapJson(key)
	if err != nil {
		return
	}

	idl.Count = -1
	db.GetCache().SetInterface(key, idl)
}

func generateDiscordMsgHash(url string) string {
	_parts := strings.Split(url, "_")
	return strings.Split(_parts[len(_parts)-1], ".")[0]
}
