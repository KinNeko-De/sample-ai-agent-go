package llm

import "io"

const RoleSystem = "system"
const RoleUser = "user"

type Message struct {
	Role    string
	Content string
}

type LLM interface {
	Chat(messages []Message, writer io.Writer) (Message, error)
}
