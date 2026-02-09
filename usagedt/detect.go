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
		detectLogger.Warnf("Failed to marshal payload json: %v", err)
		return nil
	}

	detectLogger.Debugf("Send request")

	req, err := http.NewRequest("POST", ApiEndpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		detectLogger.Warnf("Failed to create request: %v", err)
		return nil
	}

	req.Header.Set("Authorization", "Bearer "+viper.GetString("ai.apikey"))
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		detectLogger.Warnf("Failed to send request: %v", err)
		return nil
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		detectLogger.Warnf("Failed to read response body: %v", err)
		return nil
	}

	detectLogger.Debugf("Received response: %s", bodyText)

	var data map[string]interface{}
	err = json.Unmarshal(bodyText, &data)
	if err != nil {
		detectLogger.Warnf("Failed to unmarshal response body: %v", err)
		return nil
	}

	choicesVal, ok := data["choices"]
	if !ok {
		detectLogger.Warnf("The choice field does not exist")
		return nil
	}
	choices, ok := choicesVal.([]interface{})
	if !ok {
		detectLogger.Warnf("The choice field is not an array")
		return nil
	}
	if len(choices) == 0 {
		detectLogger.Warnf("The choice field is empty")
		return nil
	}

	choice0Val, ok := choices[0].(map[string]interface{})
	if !ok {
		detectLogger.Warnf("The choice[0] field is not a map")
		return nil
	}

	messageVal, ok := choice0Val["message"]
	if !ok {
		detectLogger.Warnf("The choice[0].message field does not exist")
		return nil
	}
	message, ok := messageVal.(map[string]interface{})
	if !ok {
		detectLogger.Warnf("The choice[0].message field is not a map")
		return nil
	}

	contentVal, ok := message["content"]
	if !ok {
		detectLogger.Warnf("The choice[0].message.content field does not exist")
		return nil
	}
	content, ok := contentVal.(string)
	if !ok {
		detectLogger.Warnf("The choice[0].message.content field is not a string")
		return nil
	}

	result := &Result{}
	err = json.Unmarshal([]byte(content), &result)
	if err != nil {
		detectLogger.Warnf("Failed to unmarshal content: %v", err)
		return nil
	}

	return result
}
