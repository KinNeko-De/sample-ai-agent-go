package memory

import (
	"fmt"
	"log/slog"

	"github.com/kinneko-de/sample-ai-agent-go/internal/llm"
)

type ShortTerm struct {
	messages []llm.Message
}

func NewShortTerm() *ShortTerm {
	return &ShortTerm{}
}

func (m *ShortTerm) Add(message llm.Message) {
	m.addMessage(message)
}

func (m *ShortTerm) addMessage(message llm.Message) {
	m.messages = append(m.messages, message)
}

func (m *ShortTerm) Messages(systemPrompt llm.Message) []llm.Message {
	messages := make([]llm.Message, 0, 1+len(m.messages))
	messages = append(messages, systemPrompt)
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
