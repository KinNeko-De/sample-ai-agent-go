package memory

import (
	"fmt"
	"log/slog"

	"github.com/kinneko-de/sample-ai-agent-go/internal/llm"
)

const roleSystem = "system"
const roleUser = "user"
const roleAssistant = "assistant"

type ShortTerm struct {
	messages []llm.Message
}

func NewShortTerm() *ShortTerm {
	return &ShortTerm{}
}

func (m *ShortTerm) AddUserInput(content string) {
	m.add(roleUser, content)
}

func (m *ShortTerm) AddAssistantResponse(content string) {
	m.add(roleAssistant, content)
}

func (m *ShortTerm) add(role, content string) {
	m.messages = append(m.messages, llm.Message{Role: role, Content: content})
}

// Messages returns the system prompt followed by all conversation turns,
// ready to be sent to the LLM.
func (m *ShortTerm) Messages() []llm.Message {
	messages := make([]llm.Message, 0, 1+len(m.messages))
	messages = append(messages, llm.Message{Role: roleSystem, Content: systemPrompt})
	messages = append(messages, m.messages...)
	return messages
}

func (m *ShortTerm) Log() {
	entries := make([]any, len(m.messages))
	for i, msg := range m.messages {
		entries[i] = slog.Group(fmt.Sprintf("%d", i), "role", msg.Role, "content", msg.Content)
	}
	slog.Debug("short-term memory", slog.Group("messages", entries...))
}
