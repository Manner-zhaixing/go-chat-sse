package third

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	apiURL   = "https://api.deepseek.com/chat/completions"
	apiKey   = "sk-a2b35940a42e4e84a3896b2cfc5e74b5" // 替换为你的实际API密钥
	apiModel = "deepseek-chat"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream"`
}

type ChatResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Message Message `json:"message"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}
type StreamChatResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Delta struct {
			Content string `json:"content"`
		} `json:"delta"`
		Index        int    `json:"index"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
}

func StreamChatRequest(requestData ChatRequest) (*http.Response, error) {
	// 编码请求体
	requestData.Model = apiModel
	requestData.Stream = true
	requestBody, err := json.Marshal(requestData)
	if err != nil {
		return nil, fmt.Errorf("编码请求体失败: %w", err)
	}
	fmt.Println(string(requestBody))
	// 创建HTTP请求
	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %w", err)
	}
	//defer resp.Body.Close()
	return resp, nil
}
