package services

import (
	"encoding/json"
	"fmt"
	"lark/initialization"
	"net/http"
	"strings"
)

type RequestTrigger struct {
	Type         string `json:"type"`
	DiscordMsgId string `json:"discordMsgId,omitempty"`
	MsgHash      string `json:"msgHash,omitempty"`
	Prompt       string `json:"prompt,omitempty"`
	Index        int64  `json:"index,omitempty"`
}

func ReqMidjourney(params RequestTrigger) error {
	data, err := json.Marshal(params)
	if err != nil {
		fmt.Println("json marshal error: ", err)
		return err
	}
	req, err := http.NewRequest("POST", initialization.GetConfig().DISCORD_MIDJOURNEY_URL, strings.NewReader(string(data)))
	if err != nil {
		fmt.Println("http request error: ", err)
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("http request error: ", err)
		return err
	}
	defer resp.Body.Close()
	return nil
}
