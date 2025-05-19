package third

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	apiURL   = "https://api.deepseek.com/chat/completions"
	apiKey   = "sk-a2b35940a42e4e84a3896b2cfc5e74b5" // 替换为你的实际API密钥
	apiModel = "deepseek-chat"
)

var DsDataChannel = make(chan string)

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
	} `json:"choices"`
}

func StreamChatRequest(requestData ChatRequest) error {
	// 编码请求体
	requestData.Stream = true
	requestData.Model = apiModel
	requestBody, err := json.Marshal(requestData)
	if err != nil {
		return fmt.Errorf("编码请求体失败: %w", err)
	}

	// 创建HTTP请求
	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return fmt.Errorf("创建请求失败: %w", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("发送请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API返回错误: %s, 响应体: %s", resp.Status, string(body))
	}

	// 处理流式响应
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}

		// SSE 数据行以 "data: " 开头
		if bytes.HasPrefix(line, []byte("data: ")) {
			data := line[6:] // 去掉 "data: " 前缀

			// 检查是否是结束标记
			if string(data) == "[DONE]" {
				fmt.Println("\n流式传输结束")
				break
			}

			// 解析JSON数据
			var chunk StreamChatResponse
			if err := json.Unmarshal(data, &chunk); err != nil {
				return fmt.Errorf("解析流数据失败: %w", err)
			}

			// 打印内容
			for _, choice := range chunk.Choices {
				if choice.Delta.Content != "" {
					DsDataChannel <- choice.Delta.Content
					fmt.Print(choice.Delta.Content)
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("读取流数据失败: %w", err)
	}

	return nil
}
func main() {
	// 准备请求数据
	requestData := ChatRequest{
		Messages: []Message{
			{Role: "system", Content: "You are a helpful assistant."},
			{Role: "user", Content: "Hello!介绍一下你自己"},
		},
	}
	// 发送流式请求
	streamChatRequest(requestData)
}
