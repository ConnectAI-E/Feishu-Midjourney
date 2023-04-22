package services

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// ChatGPTResponseBody 返回
type ChatGPTResponseBody struct {
	Meta struct {
		ErrCode    int    `json:"errCode"`
		ErrMsg     string `json:"errMsg"`
		RequestIdf string `json:"requestId"`
	} `json:"meta"`
	Data struct {
		Content string `json:"content"`
		Role    string `json:"role"`
		Model   string `json:"model"`
		Usage   struct {
			PromptTokens     int `json:"prompt_tokens"`
			CompletionTokens int `json:"completion_tokens"`
			TotalTokens      int `json:"total_tokens"`
		} `json:"usage"`
	} `json:"data"`
}

type Messages struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatGPTRequestBody
type ChatGPTRequestBody struct {
	GameId      string     `json:"gameId"`
	Messages    []Messages `json:"messages"`
	IsDowngrade bool       `json:"isDowngrade"`
}

type ChatGPT struct {
	GameId      string
	Secret      string
	IsDowngrade bool
	URL         string
}

func stringMd5(str string) string {
	hash := md5.New()
	hash.Write([]byte(str))

	return hex.EncodeToString(hash.Sum(nil))
}

func (gpt ChatGPT) GPT(msg []Messages) (resp Messages, err error) {

	requestBody := ChatGPTRequestBody{
		GameId:      gpt.GameId,
		IsDowngrade: gpt.IsDowngrade,
		Messages:    msg,
	}

	requestData, err := json.Marshal(requestBody)

	timestamp := time.Now().Unix()
	str := fmt.Sprintf("content=%ssecret=%stimestamp=%d", stringMd5(string(requestData)), gpt.Secret, timestamp)
	fmt.Println("md5 content: ", str, "\ncontent: ", string(requestData))
	sign := stringMd5(str)

	if err != nil {
		return resp, err
	}
	url := fmt.Sprintf("%s/sapi/v1/chatGPT/chat?sign=%s&timestamp=%d", gpt.URL, sign, timestamp)

	return gpt.Completions(url, requestData)
}

func (gpt ChatGPT) GPT4(msg []Messages) (resp Messages, err error) {

	requestBody := ChatGPTRequestBody{
		GameId:   gpt.GameId,
		Messages: msg,
	}

	requestData, err := json.Marshal(requestBody)

	timestamp := time.Now().Unix()
	str := fmt.Sprintf("content=%ssecret=%stimestamp=%d", stringMd5(string(requestData)), gpt.Secret, timestamp)
	fmt.Println("md5 content: ", str, "\ncontent: ", string(requestData))
	sign := stringMd5(str)

	if err != nil {
		return resp, err
	}
	url := fmt.Sprintf("%s/sapi/v1/chatGPT/chatWith4?sign=%s&timestamp=%d", gpt.URL, sign, timestamp)

	return gpt.Completions(url, requestData)
}

func (gpt ChatGPT) Completions(url string, requestData []byte) (resp Messages, err error) {

	fmt.Println("\nrequest gpt json string : ", string(requestData))
	fmt.Println("\nrequest gpt url : ", url)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestData))
	if err != nil {
		return resp, err
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{Timeout: 300 * time.Second}
	response, err := client.Do(req)
	if err != nil {
		return resp, err
	}

	defer response.Body.Close()

	if response.StatusCode/2 != 100 {
		return resp, fmt.Errorf("gpt api %s", response.Status)
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("gpt服务错误: ", string(body), err.Error())
		return resp, err
	}

	gptResponseBody := &ChatGPTResponseBody{}
	err = json.Unmarshal(body, gptResponseBody)
	if err != nil {
		fmt.Println("序列化json错误: ", string(body), err.Error())
		return resp, err
	}

	if gptResponseBody.Meta.ErrCode != 0 {
		fmt.Println("平台服务response错误: ", string(body))
		return resp, fmt.Errorf("%s", gptResponseBody.Meta.ErrMsg)
	}

	resp = Messages{
		Role:    gptResponseBody.Data.Role,
		Content: gptResponseBody.Data.Content,
	}
	return resp, nil
}

type ImageGenerationRequestBody struct {
	Prompt         string `json:"prompt"`
	N              int    `json:"n"`
	Size           string `json:"size"`
	ResponseFormat string `json:"response_format"`
}

type ImageGenerationResponseBody struct {
	Created int64 `json:"created"`
	Data    []struct {
		Url string `json:"url"`
	} `json:"data"`
}

func (gpt ChatGPT) GenerateImage(prompt string, size string,
	n int) ([]string, error) {
	requestBody := ImageGenerationRequestBody{
		Prompt:         prompt,
		N:              n,
		Size:           size,
		ResponseFormat: "url",
	}
	requestData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", ""+"images/generations", bytes.NewBuffer(requestData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	// req.Header.Set("Authorization", "Bearer "+gpt.ApiKey)
	client := &http.Client{Timeout: 110 * time.Second}
	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode/2 != 100 {
		return nil, fmt.Errorf("image generation api %s",
			response.Status)
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	imageResponseBody := &ImageGenerationResponseBody{}
	err = json.Unmarshal(body, imageResponseBody)
	if err != nil {
		return nil, err
	}

	var imageUrl []string
	for _, data := range imageResponseBody.Data {
		imageUrl = append(imageUrl, data.Url)
	}
	return imageUrl, nil

}

func FormatQuestion(question string) string {
	return "Answer:" + question
}

func (gpt ChatGPT) GenerateOneImage(prompt string, size string) (string, error) {
	urls, err := gpt.GenerateImage(prompt, size, 1)
	fmt.Println("gpt generate image urls: ", urls)
	if err != nil {
		return "", err
	}
	return urls[0], nil
}
