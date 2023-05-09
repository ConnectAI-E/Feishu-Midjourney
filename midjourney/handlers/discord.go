package handlers

import (
	"encoding/json"
	"fmt"
	"midjourney/initialization"
	"net/http"
	"strings"

	discord "github.com/bwmarrin/discordgo"
	"github.com/k0kubun/pp/v3"
)

type Scene string

const (
	// FirstTrigger /** 首次触发生成 */
	FirstTrigger Scene = "FirstTrigger"
	// GenerateEnd /** 生成图片结束 */
	GenerateEnd Scene = "GenerateEnd"
	// GenerateEditError /** 发送的指令midjourney生成过程中发现错误 */
	GenerateEditError Scene = "GenerateEditError"
	/**
	 * 富文本
	 */
	RichText Scene = "RichText"
	/**
	 * 发送的指令midjourney直接报错或排队阻塞不在该项目中处理 在业务服务中处理
	 * 例如：首次触发生成多少秒后没有回调业务服务判定会指令错误或者排队阻塞
	 */
)

func DiscordMsgCreate(s *discord.Session, m *discord.MessageCreate) {
	// 过滤频道
	if m.ChannelID != initialization.GetConfig().DISCORD_CHANNEL_ID {
		return
	}

	// 过滤掉自己发送的消息
	if m.Author.ID == s.State.User.ID {
		return
	}

	/******** *********/
	pp.Println(m.Content)
	pp.Println(m.Attachments)
	/******** *********/

	if strings.Contains(m.Content, "(Waiting to start)") && !strings.Contains(m.Content, "Rerolling **") {
		trigger(m.Content, FirstTrigger)
		return
	}
	for _, attachment := range m.Attachments {
		if attachment.Width > 0 && attachment.Height > 0 {
			replay(m)
			return
		}
	}
}

func DiscordMsgUpdate(s *discord.Session, m *discord.MessageUpdate) {
	// 过滤频道
	if m.ChannelID != initialization.GetConfig().DISCORD_CHANNEL_ID {
		return
	}

	if m.Author == nil {
		return
	}

	// 过滤掉自己发送的消息
	if m.Author.ID == s.State.User.ID {
		return
	}

	if strings.Contains(m.Content, "(Stopped)") {
		trigger(m.Content, GenerateEditError)
		return
	}
	if len(m.Embeds) > 0 {
		send(m.Embeds)
		return
	}
}

type ReqCb struct {
	Embeds  []*discord.MessageEmbed `json:"embeds,omitempty"`
	Discord *discord.MessageCreate  `json:"discord,omitempty"`
	Content string                  `json:"content,omitempty"`
	Type    Scene                   `json:"type"`
}

func replay(m *discord.MessageCreate) {
	body := ReqCb{
		Discord: m,
		Type:    GenerateEnd,
	}
	request(body)
}

func send(embeds []*discord.MessageEmbed) {
	body := ReqCb{
		Embeds: embeds,
		Type:   RichText,
	}
	request(body)
}

func trigger(content string, t Scene) {
	body := ReqCb{
		Content: content,
		Type:    t,
	}
	request(body)
}

func request(params interface{}) {
	data, err := json.Marshal(params)

	// fmt.Println("请求回调接口", string(data))

	if err != nil {
		fmt.Println("json marshal error: ", err)
		return
	}
	req, err := http.NewRequest("POST", initialization.GetConfig().CB_URL, strings.NewReader(string(data)))
	if err != nil {
		fmt.Println("http request error: ", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("http request error: ", err)
		return
	}
	defer resp.Body.Close()
}
