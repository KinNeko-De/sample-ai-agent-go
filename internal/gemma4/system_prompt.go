package gemma4

import (
	"github.com/kinneko-de/sample-ai-agent-go/internal/llm"
)

const thinkingActivated = "<|think|>" // The api of ollama adds this already. It also cutes out the thinking part of the response.

type Gemma4SystemPrompt struct {
}

func (s *Gemma4SystemPrompt) Generate(systemPrompt string) llm.Message {
	return llm.Message{
		Role:    llm.RoleSystem,
		Content: systemPrompt,
	}
}
