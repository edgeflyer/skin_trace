package deepseek

// Message 表示一条对话消息。
type Message struct {
	Role       string     `json:"role"`
	Content    any        `json:"content"`
	ToolCalls  []ToolCall `json:"tool_calls,omitempty"`
	ToolCallID string     `json:"tool_call_id,omitempty"`
}

// Tool 描述一个可供模型调用的工具。
type Tool struct {
	Type     string             `json:"type"`
	Function FunctionDefinition `json:"function"`
}

type FunctionDefinition struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Parameters  map[string]any `json:"parameters"`
}

// ToolCall 表示模型发起的一次工具调用。
type ToolCall struct {
	ID       string       `json:"id"`
	Type     string       `json:"type"`
	Function FunctionCall `json:"function"`
}

type FunctionCall struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

// ChatRequest 是发送给 DeepSeek 的请求。
type ChatRequest struct {
	Model          string          `json:"model"`
	Messages       []Message       `json:"messages"`
	Tools          []Tool          `json:"tools,omitempty"`
	Thinking       ThinkingConfig  `json:"thinking"`
	ResponseFormat *ResponseFormat `json:"response_format,omitempty"`
}

type ThinkingConfig struct {
	Type string `json:"type"`
}

// ChatResponse 是 DeepSeek 返回的响应。
type ChatResponse struct {
	Choices []Choice `json:"choices"`
}

type Choice struct {
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
}

type ContentPart struct {
	Type     string    `json:"type"`
	Text     string    `json:"text,omitempty"`
	ImageURL *ImageURL `json:"image_url,omitempty"`
}

type ImageURL struct {
	URL string `json:"url"`
}

type ResponseFormat struct {
	Type string `json:"type"`
}
