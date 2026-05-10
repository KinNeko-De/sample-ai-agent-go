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
- System prompt (a `role: "system"` message prepended to the history on every call that sets the agent's identity, behavior rules, and output format expectations; updated as new capabilities are added)
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
- Structured step schema (the LLM fills a JSON object per ReAct cycle for internal analysis; the answer field is extracted and displayed to the user as plain text)
    - The schema is also provided to the LLM in the system prompt so it knows exactly what to fill and how
    - **Phase 1 — Simple flat schema** (used while there are no real tools / mock-only stage):
        - `thought` — the model's internal reasoning for this step (free text, required every step)
        - `action` — what the model decided to do next (free text, e.g. `"answer"` or a tool name)
        - `input` — the input or query the model is acting on (string, or `null` if not applicable)
        - `answer` — the response shown to the user; plain natural language, no JSON
    - **Phase 2 — Extended schema** (introduced together with real tool calling in feature 2):
        - `thought` — the model's internal reasoning for this step; what it knows, what it still needs, and why it is taking the next action (free text, required every step)
        - `action` — what the model decided to do next; either `"tool_call"` (needs more information) or `"final_answer"` (ready to respond to the user)
        - `tool` — name of the tool to call; only set when `action` is `"tool_call"`, otherwise `null`
        - `tool_input` — a JSON object containing the arguments for the tool; only set when `action` is `"tool_call"`, otherwise `null`; shape must match the tool's registered schema
        - `observation` — left `null` by the LLM; filled in by the agent after executing the tool, then sent back to the LLM in the next cycle
        - `final_answer` — the response shown to the user; only set when `action` is `"final_answer"`, otherwise `null`; plain natural language, no JSON
    - Enables reliable step logging and reasoning chain visualization (see feature 6)
    - Requires robust JSON extraction and repair before unmarshaling (see "JSON Formatting failures" in feature 5)

### 4. State & Memory & History Management

- Short-term conversation history (pass prior turns to each LLM call)
    - **Phase 1 — Simple message history** (introduced alongside the REPL in feature 1): each turn is stored as a message with a `role` (`user` or `assistant`) and a `content` string (the raw user input and the LLM's plain-text answer)
    - **Phase 2 — Structured message history** (introduced alongside the ReAct step schema in feature 3):
        - the assistant message's `content` stores the full JSON step object so the conversation history reflects the structured reasoning, not just the final answer
        - after a tool is executed, its result is appended as a separate message with `role: "tool"` and `content` containing the tool result; this follows the OpenAI-compatible API convention and keeps tool output semantically distinct from user and assistant messages
- Context Management
    - Reducing History Length
        - Sliding Window (keep only the last N messages to stay within token limits)
        - Summarization (compress older history into a summary)
        - History Lookup (retrieve relevant past turns by query)
    - Validating Output with Tool Result (cross-check model answer against tool data)

### 5. Guardrails & Sanitization

Input/Output Sanitization belongs here — it is a safety concern, not a memory concern.

- Input Sanitization (clean/validate user input before sending to LLM)
- Prompt Injection Defense (prevent user input from overriding the system prompt or hijacking agent behavior)
    - **Step 1 — Structural separation**: user input is always passed as a `role: "user"` message, never interpolated directly into the system prompt string; this is the baseline and must be in place from sprint 00
    - **Step 2 — Input heuristic detection**: scan user input for known injection patterns before sending to the LLM (e.g. phrases like "ignore previous instructions", "you are now", "disregard your system prompt"); reject or sanitize on match
    - **Step 3 — Output validation**: after each LLM response, check that the output still conforms to the expected schema and domain; a response that suddenly ignores the JSON schema or changes persona is a signal that injection may have succeeded
- Output Sanitization (clean/validate LLM output before displaying or acting on it)
- JSON Formatting failures (handle cases where model returns malformed JSON for tool calls)
    - **Step 1 — Blind retry**: call the LLM again with the identical prompt and history, up to a fixed number of attempts (e.g. 3), before giving up; simple to implement, no extra prompt engineering required
    - **Step 2 — Correction retry**: on a parse failure, append a follow-up message telling the LLM exactly what went wrong and what to produce instead (e.g. *"Your last response was not valid JSON. Respond only with valid JSON matching this schema: ..."*); more effective than blind retry because the model is given explicit guidance on how to fix the output
- Domain Gating (refuse requests outside the intended domain)
- Fact Alignment (check model claims against known ground truth / tool results)
- LLM-as-Judge — Runtime (a second LLM call evaluates the agent's response before it is shown to the user)
    - Judges against defined criteria: correctness, tone, safety, domain relevance, schema compliance
    - If the response fails the judgment it can be blocked, retried, or flagged
    - Trade-off: doubles LLM calls and adds latency on every turn

### 6. Advanced tool handling
Up first all tools are sent to the LLM in the system prompt. As the number of tools grow, we want a) to optimizes the context window, so that we do not have to pass all tools everytime b) to get better results, and therefore do not confuse the LLM with all the tool calls.
To solve that, we want a dynamic lookup of tools over a search tool that the LLM can call. This is also know under the term 'skills' used by Anthropic.
- Define the tool to search other tools and tell the LLM to use it
- Phase 1: The search tool should lookup tools with a simple "string contains" method
- Phase 2: The search tool should lookup tools using a semantic search (Vektor-DB/Embeddings)
- Optimize round trips: If a tool were already searched, send it to the LLM each time. It is likely that the user needs the tool a second time

### 7. Observability

- Step/trace logging (log each thought, action, observation as it happens)
- Reasoning chain visualization (display the full ReAct trace for debugging)
- Metrics (counters and measurements exported in a standard format, e.g. Prometheus)
    - Request count (total LLM calls made)
    - Token usage per call (prompt tokens, completion tokens, total)
    - Tool call count per tool (how often each tool was invoked)
    - Retry count (how many blind or correction retries occurred)
    - Error count by type (parse failures, LLM errors, tool errors)
    - Response latency (time from user input to final answer)
    - Tracing (structured spans capturing the lifecycle of a single agent turn, e.g. OpenTelemetry)
    - One root span per user turn
    - Child spans for each LLM call, tool execution, and JSON parse/repair step
    - Span attributes: model name, token counts, tool name, retry attempt number
    
### 8. Advanced Planning

- Decomposition (break complex goals into sub-tasks)
- Plan & Execute (generate a plan first, then execute steps one by one)

### 9. Voice

- Voice output
    - browser output
    - natural language
- Voice input

### 10. Subagents

- Domain-Oriented Specialization (separate agents for Contact, Policy, etc.)
- Semantic Intent Mapping & Routing ("My father died. Do I get extra holiday? Who approves that?" → mapped to correct specialist agent)

### 11. Evaluation

- Test harness with expected outputs (run predefined prompts, assert expected results — verifies changes don't break behavior)
- LLM-as-Judge — Offline (a second LLM evaluates agent responses as a batch step, not in the live request path)
    - Given a prompt, the agent's response, and optionally an expected answer, the judge scores or flags the response
    - Useful for regression testing after changes — similar to how GitHub Copilot reviews AI-generated PRs
    - No latency cost since it runs outside the live agent loop

### 12. Documentation

- Architectural Decision Records (ADRs — document why key technical choices were made)
- arc42 template (structured architecture documentation)

### 13. End-User Interface

The CLI/REPL is a developer tool. A real end-user interface will be added at a late stage once the agent behavior is stable.

- Web UI or chat-style interface for non-technical users (managers, back-office staff)
- Internal trace output (`thought`, `action`, `observation`) is shown in a separate, optional panel — hidden by default, user can opt in to open it
- Error messages are human-friendly, not technical stack traces
