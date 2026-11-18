### Mode 1: Transparent AI Gateway (Protocol-Level Proxy)

*Purpose: Act as a drop-in, OpenAI-compatible middleware layer for applications (Laravel, n8n, external scripts) to centralize authentication, usage tracking, resiliency, and provider routing without altering client-side prompt logic.*

#### 1. API & Protocol Compatibility
* **Standard Specification Emulation:** Expose endpoints adhering to the standard OpenAI schema (`/v1/chat/completions`, `/v1/embeddings`, `/v1/models`).
* **Full Streaming Support:** Support real-time Server-Sent Events (SSE) streaming while continuing to intercept and calculate token metrics on completion.
* **Dynamic Model Discovery:** Expose an aggregated `/v1/models` endpoint reflecting active providers and custom virtual aliases to populate client UIs (e.g., n8n dropdowns).

#### 2. Authentication & Secrets Isolation
* **Zero Client-Side Provider Keys:** Clients authenticate using an internal gateway credential; upstream provider keys (OpenAI, Anthropic, Bedrock, etc.) are injected strictly by the proxy.
* **Client & Workflow Attribution:** Identify the calling source via API keys or custom headers to separate metrics (e.g., distinguishing production Laravel traffic from internal n8n automation).

#### 3. Routing & Model Aliasing
* **Virtual Model Aliases:** Support semantic model identifiers (e.g., `default`, `fast`, `smart`) that can be remapped dynamically to concrete upstream models via configuration without code changes.
* **Default Fallbacks:** If a client provides an empty or default model parameter, the gateway automatically resolves and applies the current system default.

#### 4. Resiliency & Traffic Management
* **Automated Retries & Backoff:** Automatically intercept transient errors (e.g., rate limits, 5xx gateway errors) and retry with exponential backoff before surfacing errors to the client.
* **Cross-Provider Failover:** Transparently reroute failed requests to alternative providers or models when the primary upstream provider is unresponsive or rate-limited.
* **Enforced Timeouts:** Configurable request deadlines and circuit breakers to prevent hung processes in consumer applications.

#### 5. Telemetry & Governance
* **Asynchronous Audit Logging:** Log prompt tokens, completion tokens, latency, status codes, and calculated monetary costs without adding blocking latency to the request.
* **Spend & Rate Limits:** Enforce budget thresholds and rate limits per client or per virtual model alias.

---

### Mode 2: Managed Task Orchestrator & Prompt Engine (Semantic-Level Engine)

*Purpose: Act as an end-to-end task execution engine where prompts, structured output schemas, and operational parameters are fully decoupled from consumer applications. Consuming apps send only raw dynamic variables, while the gateway handles orchestration, testing, and governance.*

#### 1. Minimalist Client Interface
* **Task-Based Execution:** Expose a simplified endpoint where clients submit only an operational task identifier and dynamic input variables, requiring zero knowledge of providers, models, or prompt syntax.
* **Direct Structured Response:** Return validated, typed JSON payloads directly to the client, ready for immediate business consumption.

#### 2. Prompt Management & Dynamic Injection
* **Centralized Prompt Repository:** Store and manage system and user prompts independently of application codebases.
* **Variable Interpolation:** Support expressive dynamic variable substitution (handling strings, markdown blocks, arrays, and conditional logic) into prompt templates.
* **Prompt Versioning:** Maintain version history for prompts with environment lifecycle stages (e.g., draft, staging, production) and instant rollback capabilities.

#### 3. Structured Output & Schema Enforcement
* **Centralized JSON Schemas:** Store and maintain strict output schemas alongside task definitions.
* **Provider-Agnostic Schema Translation:** Translate stored JSON Schemas into the appropriate provider-specific format (e.g., strict JSON schema modes, tool/function calling) based on the routed model.
* **Upstream Validation:** Validate model responses against the task’s defined schema before responding to the consumer, catching malformed outputs internally.

#### 4. Task-Level Operational Controls
* **Granular Model Assignment:** Configure target models and providers independently per task, adjustable at runtime without client deploys.
* **Hyperparameter Tuning:** Manage task-specific parameters (e.g., max output tokens, temperature, reasoning effort, top-p) within the gateway.
* **Dedicated Task Resilience:** Assign task-specific timeouts, retry attempts, and prioritized fallback chains based on task criticality.

#### 5. "The Lab" / Playground UI
* **Interactive Sandbox:** An internal administrative interface to select tasks, provide test variables, and execute runs against different models.
* **Side-by-Side Model Benchmarking:** Execute identical dynamic inputs across multiple models simultaneously to compare output quality, schema adherence, token consumption, and execution time.
* **No-Deploy Tuning:** Edit and test prompt templates against live or historical test data, promoting revisions to production instantly.

#### 6. Architectural Shared-Core
* **Internal Delegation:** The task orchestration layer delegates its compiled requests directly into the Mode 1 proxy layer in-memory, automatically inheriting all logging, key management, retry handling, and cost tracking.
