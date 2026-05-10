package llm

const RoleSystem = "system"
const RoleUser = "user"

type Message struct {
	Role    string
	Content string
}

type LLM interface {
	Chat(messages []Message) (Message, error)
}
