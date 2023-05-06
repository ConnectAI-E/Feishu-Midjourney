package handlers

import (
	"bytes"
	"fmt"
	"midjourney/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ReqUploadFile struct {
	ImgData []byte `json:"imgData"`
	Name    string `json:"name"`
	Size    int64  `json:"size"`
}

func UploadFile(c *gin.Context) {
	var body ReqUploadFile
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	data, err := services.Attachments(body.Name, body.Size)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error saving image: %s", err)
		return
	}
	if len(data.Attachments) == 0 {
		c.String(http.StatusInternalServerError, "上传图片失败: %s", err)
		return
	}
	payload := bytes.NewReader(body.ImgData)
	client := &http.Client{}
	req, err := http.NewRequest("PUT", data.Attachments[0].UploadUrl, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Content-Type", "image/png")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()
	c.JSON(200, gin.H{"name": data.Attachments[0].UploadFilename})
}
