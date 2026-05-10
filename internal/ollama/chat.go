package ollama

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

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

func (c *OllamaClient) Chat(model string, messages []llm.Message, writer io.Writer) (llm.Message, error) {
	request := Request{
		Model:    model,
		Messages: messages,
		Think:    true,
		Stream:   true,
	}

	req, err := createRequet(request)
	if err != nil {
		return llm.Message{}, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return llm.Message{}, err
	}

	defer resp.Body.Close()

	fullContent, role, err := streamResponse(resp, writer)
	if err != nil {
		return llm.Message{}, err
	}

	return llm.Message{
		Role:    role,
		Content: fullContent.String(),
	}, nil
}

func streamResponse(resp *http.Response, writer io.Writer) (strings.Builder, string, error) {
	var fullContent strings.Builder
	var role string

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		var chunk Response
		if err := json.Unmarshal(scanner.Bytes(), &chunk); err != nil {
			continue
		}
		role = chunk.Message.Role
		fmt.Fprint(writer, chunk.Message.Content)
		fullContent.WriteString(chunk.Message.Content)
		if chunk.Done {
			break
		}
	}
	if err := scanner.Err(); err != nil {
		return strings.Builder{}, "", err
	}
	fmt.Fprintln(writer)
	return fullContent, role, nil
}

func createRequet(request Request) (*http.Request, error) {
	requestBytes, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", ollamaURL, bytes.NewBuffer(requestBytes))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	return req, nil
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
	Done      bool    `json:"done"`
}
