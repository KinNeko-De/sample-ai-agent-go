package llm

type SystemPrompt interface {
	Generate(systemPrompt string) Message
}
