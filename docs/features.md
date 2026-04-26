# AI Agent Feature Plan — Go + Ollama + Gemma 4

## Setup

- Language: Go
- LLM: Gemma 4 via locally running Ollama (OpenAI-compatible `/v1/` API)
- Gemma 4 supports native tool/function calling
- Purpose: Learning project — understand how agentic AI works by implementing it manually

---

## Features

### 1. Foundation

- Ollama client setup (OpenAI-compatible base URL, model config, API key placeholder)
- Basic completion: send prompt → receive response
- Streaming response display (show tokens as they arrive)
- CLI/REPL Interface
    - CLI = Command Line Interface — text-based terminal program (no GUI, no browser)
    - REPL = Read-Eval-Print Loop — reads your input, evaluates (calls LLM), prints the result, loops back. Like a chat terminal: shows a `>` prompt, you type, the agent responds, repeats until quit.
- Predefined test prompts (hardcoded inputs for repeatable manual testing)

### 2. Tool Calling

- Tool definition & registration (JSON schema describing what each tool does)
- Mock data tools (tools that return hardcoded fake data, no real backend needed)
- The agent loop (receive input → call LLM → tool_call or final answer? → if tool_call: execute tool, feed result back → repeat)
- Structured output parsing (reliably extracting JSON from model responses, even when malformed)
- Vector DB Lookup (semantic search tool backed by an embedded vector database)

### 3. Reasoning Pattern (ReAct)

Sits on top of a working agent loop.

- Thought (model reasons about what to do next)
- Action (model decides which tool to call)
- Observation (tool result fed back to model)
- Final Answer (model produces end response when done reasoning)

### 4. State & Memory & History Management

- Short-term conversation history (pass prior turns to each LLM call)
- Context Management
    - Reducing History Length
        - Sliding Window (keep only the last N messages to stay within token limits)
        - Summarization (compress older history into a summary)
        - History Lookup (retrieve relevant past turns by query)
    - Validating Output with Tool Result (cross-check model answer against tool data)

### 5. Guardrails & Sanitization

Input/Output Sanitization belongs here — it is a safety concern, not a memory concern.

- Input Sanitization (clean/validate user input before sending to LLM)
- Output Sanitization (clean/validate LLM output before displaying or acting on it)
- JSON Formatting failures (handle cases where model returns malformed JSON for tool calls)
- Domain Gating (refuse requests outside the intended domain)
- Fact Alignment (check model claims against known ground truth / tool results)

### 6. Observability

- Step/trace logging (log each thought, action, observation as it happens)
- Reasoning chain visualization (display the full ReAct trace for debugging)

### 7. Advanced Planning

- Decomposition (break complex goals into sub-tasks)
- Plan & Execute (generate a plan first, then execute steps one by one)

### 8. Subagents

- Domain-Oriented Specialization (separate agents for Contact, Policy, etc.)
- Semantic Intent Mapping & Routing ("My father died. Do I get extra holiday? Who approves that?" → mapped to correct specialist agent)

### 9. Evaluation

- Test harness with expected outputs (run predefined prompts, assert expected results — verifies changes don't break behavior)

### 10. Documentation

- Architectural Decision Records (ADRs — document why key technical choices were made)
- arc42 template (structured architecture documentation)
