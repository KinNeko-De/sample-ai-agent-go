# sample-ai-agent-go

A learning project for building an agentic AI system from scratch in Go, using a locally running [Ollama](https://ollama.com/) instance with the Gemma 4 model.

The goal is to understand how agentic AI works by implementing each component manually — no high-level agent framework is used.

---

## Overview

- **Language:** Go
- **LLM:** Gemma 4 via Ollama (OpenAI-compatible `/v1/` API)
- **Interface:** CLI/REPL — a terminal loop that reads user input and returns a response
- **Purpose:** Learning project

---

## Features

The agent is built incrementally across sprints. Planned and implemented features include:

### Foundation
- Ollama client with OpenAI-compatible API
- System prompt for agent identity and behavior
- Basic chat completion and streaming response display
- CLI/REPL interface with conversation history
- Predefined test prompts for repeatable manual testing

### Tool Calling
- Tool definition and registration via JSON schema
- Mock data tools (hardcoded responses, no real backend needed)
- Agent loop: receive input → call LLM → tool call or final answer → execute tool → feed result back → repeat
- Structured output parsing with JSON repair

### Reasoning — ReAct Pattern
- **Thought** — the model reasons about what to do next
- **Action** — the model decides which tool to call
- **Observation** — the tool result is fed back to the model
- **Final Answer** — the model responds when done reasoning
- Structured step schema as JSON, with reasoning chain logging

### State, Memory & History
- Short-term conversation history passed on every LLM call
- Structured message history including tool results (`role: "tool"`)
- Different techniques of context management

### Guardrails & Sanitization
- Input and output sanitization
- Prompt injection defense
- JSON formatting failure handling
- Domain gating and fact alignment
- LLM-as-Judge

### Observability
- Logging for each thought, action, and observation
- Reasoning chain visualization
- OpenTelemetry metrics
- OpenTelemetry tracing with span export to Jaeger/Tempo

### Advanced Planning
- Goal decomposition into sub-tasks
- Plan-and-execute pattern

### Subagents
- Domain-oriented specialization
- Semantic intent mapping and routing

### Evaluation
- Test harness with expected outputs
- LLM-as-Judge offline batch evaluation

---

## Getting Started

### Prerequisites

- [Go](https://go.dev/) 1.26+
- [Ollama](https://ollama.com/) running locally with [Gemma 4](https://deepmind.google/models/gemma/gemma-4/) pulled

---

## License

See [LICENSE](LICENSE).
