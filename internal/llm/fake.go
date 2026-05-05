package llm

// FakeLLM is a simple implementation of the LLM interface that always returns a fixed response.
// Outdated with sprint 01, where the LLM is forced to output json

type FakeLLM struct{}

func NewFakeLLM() *FakeLLM {
	return &FakeLLM{}
}

func (f *FakeLLM) Chat(messages []Message) (string, error) {
	return "I can not help with that.", nil
}
