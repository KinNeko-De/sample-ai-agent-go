package agent

import (
	"github.com/kinneko-de/sample-ai-agent-go/internal/llm"
	"github.com/kinneko-de/sample-ai-agent-go/internal/memory"
)

type Agent struct {
	memory *memory.ShortTerm
	client llm.LLM
}

func New() *Agent {
	return &Agent{
		memory: memory.NewShortTerm(),
		client: llm.NewFakeLLM(),
	}
}

func (agent *Agent) Chat(input string) (string, error) {
	agent.memory.AddUserInput(input)

	response, err := agent.client.Chat(agent.memory.Messages())
	if err != nil {
		return "", err
	}

	agent.memory.AddAssistantResponse(response)
	agent.memory.Log()
	return response, nil
}
