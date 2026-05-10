package ollama

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/kinneko-de/sample-ai-agent-go/internal/llm"
)

const ollamaURL = "http://localhost:11434/api/chat"

type OllamaChat interface {
	Chat(model string, messages []llm.Message) (string, error)
}

type OllamaClient struct {
}

func NewOllamaClient() *OllamaClient {
	return &OllamaClient{}
}

func (c *OllamaClient) Chat(model string, messages []llm.Message) (llm.Message, error) {
	request := Request{
		Model:    model,
		Messages: messages,
		Think:    true,
		Stream:   false,
	}
	requestBytes, err := json.Marshal(request)
	if err != nil {
		return llm.Message{}, err
	}
	req, err := http.NewRequest("POST", ollamaURL, bytes.NewBuffer(requestBytes))
	if err != nil {
		return llm.Message{}, err
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return llm.Message{}, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return llm.Message{}, err
	}

	var respObj Response
	err2 := json.Unmarshal(bodyBytes, &respObj)
	if err2 != nil {
		return llm.Message{}, err2
	}
	answer := llm.Message{
		Role:    respObj.Message.Role,
		Content: respObj.Message.Content,
	}

	return answer, nil
}

type Request struct {
	Model    string        `json:"model"`
	Messages []llm.Message `json:"messages"`
	Think    bool          `json:"think"`
	Stream   bool          `json:"stream"`
}

type Message struct {
	Role     string `json:"role"`
	Content  string `json:"content"`
	Thinking string `json:"thinking"`
}

type Response struct {
	Model     string  `json:"model"`
	CreatedAt string  `json:"created_at"`
	Message   Message `json:"message"`
}
