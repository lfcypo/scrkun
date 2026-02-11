package usagedt

import (
	"bytes"
	"encoding/json"
	"image"
	"io"
	"net/http"

	"github.com/lfcypo/scrkun/capture"
	"github.com/lfcypo/scrkun/logger"
	"github.com/spf13/viper"
)

var detectLogger = logger.New("Detect")

const (
	ApiEndpoint = "https://dashscope.aliyuncs.com/compatible-mode/v1/chat/completions"
	Model       = "qwen3-vl-flash"
)

func Detect(img *image.RGBA) *Result {
	client := &http.Client{}
	requestBody := RequestBody{
		Model: Model,
		Messages: []Message{
			{
				Role: "system",
				Content: []interface{}{MessageTextContent{
					Type: "text",
					Text: SysPrompt(),
				}},
			},
			{
				Role: "user",
				Content: []interface{}{MessageImageContent{
					Type: "image_url",
					ImageURL: ImageURL{
						URL: capture.RGBAtoBase64JPEG(img, 30),
					},
				}},
			},
		},
	}
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		detectLogger.Warnf("解析 json 响应失败: %v", err)
		return nil
	}

	detectLogger.Debugf("准备发送 AI 请求")

	req, err := http.NewRequest("POST", ApiEndpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		detectLogger.Warnf("创建请求失败 请检查网络是否正常: %v", err)
		return nil
	}

	req.Header.Set("Authorization", "Bearer "+viper.GetString("ai.apikey"))
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		detectLogger.Warnf("发送请求失败 请检查网络是否正常: %v", err)
		return nil
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		detectLogger.Warnf("获取响应失败 请检查网络与配置文件是否正常: %v", err)
		return nil
	}

	//detectLogger.Debugf("接收到响应: %s", bodyText)

	var data map[string]interface{}
	err = json.Unmarshal(bodyText, &data)
	if err != nil {
		detectLogger.Warnf("解析响应失败: %v", err)
		return nil
	}

	choicesVal, ok := data["choices"]
	if !ok {
		detectLogger.Warnf("期待的块缺失")
		return nil
	}
	choices, ok := choicesVal.([]interface{})
	if !ok {
		detectLogger.Warnf("期待的块不在数组内")
		return nil
	}
	if len(choices) == 0 {
		detectLogger.Warnf("期待的块为空")
		return nil
	}

	choice0Val, ok := choices[0].(map[string]interface{})
	if !ok {
		detectLogger.Warnf("choice[0] 不是 map")
		return nil
	}

	messageVal, ok := choice0Val["message"]
	if !ok {
		detectLogger.Warnf("choice[0].message 不存在")
		return nil
	}
	message, ok := messageVal.(map[string]interface{})
	if !ok {
		detectLogger.Warnf(" choice[0].message 不是 map")
		return nil
	}

	contentVal, ok := message["content"]
	if !ok {
		detectLogger.Warnf("choice[0].message.content 不存在")
		return nil
	}
	content, ok := contentVal.(string)
	if !ok {
		detectLogger.Warnf("choice[0].message.content 不是字符串")
		return nil
	}

	result := &Result{}
	err = json.Unmarshal([]byte(content), &result)
	if err != nil {
		detectLogger.Warnf("解析响应失败: %v", err)
		return nil
	}

	return result
}
