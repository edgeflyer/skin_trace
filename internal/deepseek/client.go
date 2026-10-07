package deepseek

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	baseURL    string
	apiKey     string
	model      string
	httpClient *http.Client
}

func NewClient(baseURL, apiKey, model string) (*Client, error) {
	baseURL = strings.TrimRight(baseURL, "/")

	if baseURL == "" || apiKey == "" || model == "" {
		return nil, errors.New(
			"deepseek base_url, api_key and model are required",
		)
	}

	return &Client{
		baseURL: baseURL,
		apiKey:  apiKey,
		model:   model,
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}, nil
}

func (c *Client) chat(
	ctx context.Context,
	messages []Message,
	tools []Tool,
	responseFormat *ResponseFormat,
) (Choice, error) {
	payload := ChatRequest{
		Model:          c.model,
		Messages:       messages,
		Tools:          tools,
		Thinking:       ThinkingConfig{Type: "disabled"},
		ResponseFormat: responseFormat,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return Choice{}, fmt.Errorf("encode deepseek request: %w", err)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		c.baseURL+"/chat/completions",
		bytes.NewReader(data),
	)
	if err != nil {
		return Choice{}, fmt.Errorf("build deepseek request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return Choice{}, fmt.Errorf("send deepseek request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Choice{}, fmt.Errorf("read deepseek response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return Choice{}, fmt.Errorf(
			"deepseek returned HTTP %d: %s",
			resp.StatusCode,
			string(body),
		)
	}

	var result ChatResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return Choice{}, fmt.Errorf("decode deepseek response: %w", err)
	}
	log.Printf("模型调用结果: %+v\n", result)

	if len(result.Choices) == 0 {
		return Choice{}, errors.New("deepseek returned no choices")
	}

	return result.Choices[0], nil
}

// Chat 用于普通对话，支持模型发起工具调用。
func (c *Client) Chat(
	ctx context.Context,
	messages []Message,
	tools []Tool,
) (Choice, error) {
	return c.chat(ctx, messages, tools, nil)
}

// ChatJSON 要求模型以 JSON 对象格式返回内容。
func (c *Client) ChatJSON(
	ctx context.Context,
	messages []Message,
) (Choice, error) {
	return c.chat(
		ctx,
		messages,
		nil,
		&ResponseFormat{Type: "json_object"},
	)
}
