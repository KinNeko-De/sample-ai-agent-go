package agent

import (
	"io"

	"github.com/kinneko-de/sample-ai-agent-go/internal/gemma4"
	"github.com/kinneko-de/sample-ai-agent-go/internal/llm"
	"github.com/kinneko-de/sample-ai-agent-go/internal/memory"
)

const systemPrompt = "You are a helpful company assistant. Answer employee questions about the company."

type Agent struct {
	memory *memory.ShortTerm
	client llm.LLM
}

func New() *Agent {
	return &Agent{
		memory: memory.NewShortTerm(),
		client: gemma4.NewGemma4LLM(),
	}
}

func (agent *Agent) Chat(input string, writer io.Writer) error {
	userMessage := llm.Message{
		Role:    llm.RoleUser,
		Content: input,
	}
	agent.memory.Add(userMessage)

	builder := &gemma4.Gemma4SystemPrompt{}
	systemMessage := builder.Generate(systemPrompt)

	response, err := agent.client.Chat(agent.memory.Messages(systemMessage), writer)
	if err != nil {
		return err
	}

	agent.memory.Add(response)
	agent.memory.Log()
	return nil
}
