package agent

import (
	"encoding/json"
	"log/slog"

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
		client: llm.NewStructuredFakeLLM(),
	}
}

func (agent *Agent) Chat(input string) string {
	agent.memory.AddUserInput(input)

	response, err := agent.client.Chat(agent.memory.Messages())
	if err != nil {
		slog.Error("LLM chat error", slog.Any("error", err))
		return "Something went wrong. Please try again."
	}

	agent.memory.AddAssistantResponse(response)
	agent.memory.Log()

	structuredResponse, err := parseResponse(response)
	if err != nil {
		slog.Error("Invalid structured response", slog.Any("error", err))
		return "The LLM returned an invalid response. Please try again."
	}

	if structuredResponse.Answer != nil {
		return *structuredResponse.Answer
	}

	return "tool calls are not implemented yet"
}

func parseResponse(response string) (*StructuredResponse, error) {
	var structuredResponse StructuredResponse
	err := json.Unmarshal([]byte(response), &structuredResponse)
	if err != nil {
		return nil, err
	}
	return &structuredResponse, nil
}

type StructuredResponse struct {
	Thought string  `json:"thought"`
	Action  *string `json:"action"`
	Input   *string `json:"input"`
	Answer  *string `json:"answer"`
}
