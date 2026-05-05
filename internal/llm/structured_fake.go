package llm

type StructuredFakeLLM struct{}

func NewStructuredFakeLLM() *StructuredFakeLLM {
	return &StructuredFakeLLM{}
}

func (f *StructuredFakeLLM) Chat(messages []Message) (string, error) {
	return jsonResponseAnswer, nil
}

const jsonResponseAnswer = `{
  "thought" : "I coud not find information about this topic",
  "action" : null,
  "input" : null,
  "answer" : "I can not help with that."
}`

const jsonResponseTool = `{
  "thought" : "I need to fetch that information",
  "action" : "tool_name",
  "input" : "tool_input",
  "answer" : null
}`
