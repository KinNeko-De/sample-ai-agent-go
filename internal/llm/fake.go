package llm

type FakeLLM struct{}

func NewFakeLLM() *FakeLLM {
	return &FakeLLM{}
}

func (f *FakeLLM) Chat(messages []Message) (string, error) {
	return "I can not help with that.", nil
}
