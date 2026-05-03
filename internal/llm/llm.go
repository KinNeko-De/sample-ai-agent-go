package llm

type Message struct {
	Role    string
	Content string
}

type LLM interface {
	Chat(messages []Message) (string, error)
}
