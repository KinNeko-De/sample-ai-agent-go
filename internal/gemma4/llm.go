package gemma4

import (
	"github.com/kinneko-de/sample-ai-agent-go/internal/llm"
	"github.com/kinneko-de/sample-ai-agent-go/internal/ollama"
)

const modelName = "gemma4:e4b"

type Gemma4LLM struct {
}

func NewGemma4LLM() *Gemma4LLM {
	return &Gemma4LLM{}
}

func (l *Gemma4LLM) Chat(messages []llm.Message) (llm.Message, error) {
	client := ollama.NewOllamaClient()
	response, err := client.Chat(modelName, messages)
	if err != nil {
		return llm.Message{}, err
	}
	return response, nil
}
