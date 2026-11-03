---
title: "Google Gemini"
description: "Use Gemini 2.5 Flash, Gemini 3.1 Pro, and other Google models with Docker Agent."
keywords: docker agent, ai agents, model providers, llm, google gemini
weight: 120
canonical: https://docs.docker.com/ai/docker-agent/providers/google/
---

_Use Gemini 2.5 Flash, Gemini 3.1 Pro, and other Google models with Docker Agent._

## Setup

Docker Agent reads the first credential it finds from these environment variables (see `pkg/model/provider/gemini/client.go`):

| Variable                    | Purpose                                                                             |
| --------------------------- | ----------------------------------------------------------------------------------- |
| `GOOGLE_API_KEY`            | Primary Gemini API key.                                                             |
| `GEMINI_API_KEY`            | Alternative name for the Gemini API key (also used by the official Google SDK).     |
| `GOOGLE_GENAI_USE_VERTEXAI` | Boolean flag (`true` or `1` enables Vertex AI); unset, empty, `false`, `0`, or invalid values do not enable it. |
| `GOOGLE_CLOUD_PROJECT`      | GCP project used when `GOOGLE_GENAI_USE_VERTEXAI` is enabled or for Vertex AI Model Garden. |
| `GOOGLE_CLOUD_LOCATION`     | GCP region for Vertex AI (defaults to the SDK default).                             |

Explicit `provider_opts.project` or `provider_opts.location` selects Vertex AI regardless of the flag.

On the Gemini Developer API, a model or [custom provider](../custom/index.md) that sets `token_key` reads its key from that variable instead of `GOOGLE_API_KEY` / `GEMINI_API_KEY`. The Vertex AI backends use Application Default Credentials and ignore `token_key`.

```bash
# Gemini Developer API
export GOOGLE_API_KEY="AI..."   # or GEMINI_API_KEY

# Vertex AI (no API key; uses Application Default Credentials)
gcloud auth application-default login
export GOOGLE_GENAI_USE_VERTEXAI=1
export GOOGLE_CLOUD_PROJECT="my-gcp-project"
export GOOGLE_CLOUD_LOCATION="us-central1"
```

## Configuration

### Inline

```yaml
agents:
  root:
    model: google/gemini-3.8-flash
```

### Named Model

```yaml
models:
  gemini:
    provider: google
    model: gemini-3.8-flash
    temperature: 0.5
```

## Available Models

| Model                     | Best For                        |
| -------------------------- | ------------------------------- |
| `gemini-3.1-pro-preview`  | Most capable Gemini model       |
| `gemini-3.8-flash`        | Fast, efficient, good balance   |
| `gemini-2.5-flash`        | Fast inference, cost-effective  |
| `gemini-2.5-pro`          | Strong reasoning, large context |

## Generated Images

Some Gemini models (e.g. `gemini-2.5-flash-image`) are designed to generate
an image directly as part of their reply, not just describe one. When the
model's image output capability resolves as enabled — see
[Output capabilities](../../configuration/models/index.md#output-capabilities)
for the `output_capabilities.image` / models.dev precedence rules — Docker
Agent requests that image output on supported Google surfaces: the models
gateway, direct Gemini API, and Vertex AI. Each eligible ordinary chat request
asks for text *and* image output. Vertex AI has deterministic guard/predicate
coverage; live image-generation validation is deferred.

```yaml
models:
  gemini-image:
    provider: google
    model: gemini-2.5-flash-image
```

`output_capabilities.image` can be omitted for models that models.dev lists
with image output; set it explicitly for custom models or to override
incorrect catalogue data. Capability is never guessed from the model name.

Session-title and compaction requests omit image response modalities and
bypass the guard even for image-output-capable models. They do not explicitly
force TEXT-only output. Ordinary image-output requests with custom function tools or
structured output are rejected locally before any request is sent. Google
Search, Maps, and code-execution built-ins remain available. When custom
tools conflict, Docker Agent uses models.dev's `tool_call` capability to clarify
whether the model cannot call tools at all or supports tools only outside an
image-output request. Unknown catalogue data keeps the conservative generic
message.

A few provider-side behaviors to know:

- **The provider decides the image format** (typically PNG). Asking for a
  `.gif` or `.svg` filename does not transcode anything — the saved file's
  extension is corrected to match the data actually returned.
- **An image is not guaranteed.** Even a correctly configured image model
  can answer with text only and generate no image; that is provider behavior,
  so reword or repeat the prompt. At a text-only stop, Docker Agent may add a
  nonfatal phrase-based warning to the turn — see
  [Generated Media](../../features/tui/index.md#generated-media) for the exact
  check and its limitations.

Docker Agent attempts to save each returned image into the session workspace,
keep a portable copy in the session database, and display it in the TUI. See
[Generated Media Files](../../features/sessions/index.md#generated-media-files)
for persistence and portability, and
[Generated Media](../../features/tui/index.md#generated-media) for file
naming, collision handling, and rendering.

## Thinking Budget

Gemini supports two approaches depending on the model version:

> [!WARNING]
> **Different thinking formats**
>
> Gemini 2.5 uses **token-based** budgets (integers). Gemini 3 uses **level-based** budgets (strings like `low`, `high`). Make sure you use the right format for your model version.

### Gemini 2.5 (Token-based)

```yaml
models:
  gemini-no-thinking:
    provider: google
    model: gemini-2.5-flash
    thinking_budget: 0 # disable thinking

  gemini-dynamic:
    provider: google
    model: gemini-2.5-flash
    thinking_budget: -1 # dynamic (model decides) — default

  gemini-fixed:
    provider: google
    model: gemini-2.5-flash
    thinking_budget: 8192 # fixed token budget
```

### Gemini 3 (Level-based)

```yaml
models:
  gemini-pro:
    provider: google
    model: gemini-3.1-pro-preview
    thinking_budget: high # default for Pro: low | high

  gemini-flash:
    provider: google
    model: gemini-3.8-flash
    thinking_budget: medium # default for Flash: minimal | low | medium | high
```

## Built-in Tools (Grounding)

Gemini models support built-in tools that let the model access Google Search, Google Maps,
and public URLs, or execute code directly during generation. Enable them via `provider_opts`:

```yaml
models:
  gemini-grounded:
    provider: google
    model: gemini-2.5-flash
    provider_opts:
      google_search: true
      google_maps: true
      code_execution: true
```

| Option           | Description                                          |
| ---------------- | ---------------------------------------------------- |
| `google_search`  | Enables Google Search grounding for up-to-date info  |
| `google_maps`    | Enables Google Maps grounding for location queries   |
| `code_execution` | Enables server-side code execution for computations  |
| `url_context`    | Lets Gemini read public URLs supplied in the prompt   |

### URL Context

Enable `url_context` to let Gemini read links supplied in your prompt without a
separate fetch tool:

```yaml
models:
  reader:
    provider: google
    model: gemini-3.8-flash
    provider_opts:
      url_context: true
      google_search: true # Optional: discover sources as well as read links
```

Then ask, for example, `Summarize https://docs.docker.com/ai/docker-agent/`.
Google fetches the content on its servers; it cannot read local files, private
network URLs, or pages requiring your browser's authentication. The option is
opt-in and accepts a YAML boolean, not the string `"true"`.

Docker Agent forwards this tool on the Gemini Developer API, models gateway,
and native Gemini Vertex AI paths. Model, backend, region, and tool-combination
support still depend on Google and the gateway. On supported models it can be
combined with Google Search. Mixing built-in and custom function tools currently
uses a flag that the Go SDK accepts only on the Gemini Developer API path
(including gateways forwarding to that API), not Vertex AI. This is not a
general browser or an alternative to the local fetch tool for authenticated
requests.

See Google's [URL Context guide](https://ai.google.dev/gemini-api/docs/url-context)
for supported content types, limits, and models, and
[`examples/gemini_url_context.yaml`](https://github.com/docker/docker-agent/blob/main/examples/gemini_url_context.yaml)
for a complete agent.

## Embeddings

Gemini embedding models can back the [RAG](../../tools/rag/index.md)
`chunked-embeddings` and `semantic-embeddings` strategies. Text is embedded
natively through the Gemini API (`batchEmbedContents`) or Vertex AI, and vectors
come back in input order. On Vertex AI, Gemini Embedding 2 models accept one
input per request, so batches are sent one text at a time there.

```yaml
models:
  gemini-embed:
    provider: google
    model: gemini-embedding-2
    provider_opts:
      output_dimensionality: 768   # optional; defaults to the model's native size

rag:
  docs:
    docs: [./docs]
    strategies:
      - type: chunked-embeddings
        embedding_model: gemini-embed
        database: ./docs.db
        vector_dimensions: 768       # must match output_dimensionality (or the native size)
```

| Option                  | Description                                                                                                           |
| ----------------------- | --------------------------------------------------------------------------------------------------------------------- |
| `output_dimensionality` | Positive integer. Truncates vectors to this size on models that support it (Matryoshka). Omit to use the native size. |

`gemini-embedding-001` works the same way. Vectors from different models (or
different `output_dimensionality` values) live in different vector spaces, so
rebuild the index after changing either.

Only text is embedded, and no `task_type` is sent: queries and documents share
one vector space with no automatic task prefixes, and the same configuration
works with Gemini Embedding 2, which rejects `task_type`. The Gemini Developer
API does not report token usage for embeddings, so RAG cost statistics stay at
zero; Vertex AI per-input token statistics are summed into the input token
count when present.

## Service Tier

Set `provider_opts.service_tier` to pick a Gemini API [Flex](https://ai.google.dev/gemini-api/docs/flex-inference)
or [Priority](https://ai.google.dev/gemini-api/docs/priority-inference) inference tier:

```yaml
models:
  gemini-flex:
    provider: google
    model: gemini-2.5-flash
    provider_opts:
      service_tier: flex
```

| Tier       | Description                                                                 |
| ---------- | --------------------------------------------------------------------------- |
| `flex`     | Discounted, queued processing; a request may wait before the model answers. |
| `standard` | The default tier; same as omitting the option.                              |
| `priority` | Premium pricing for faster, more consistent processing.                     |

The value is forwarded unchanged as the top-level `serviceTier` of every `generateContent` request on the Gemini Developer API, Vertex AI, and the Docker models gateway; the endpoint must support it. It applies to all requests using the configured model, including reranking, title generation, and compaction. Other values are passed through for the API to validate. When omitted or empty, no `serviceTier` is sent, leaving the API's default behavior unchanged. Non-string values are ignored. [Vertex AI Model Garden](#vertex-ai-model-garden) models do not use this option.

Flex requests can sit in Google's queue for up to 15 minutes before the first byte arrives. Docker Agent normally treats a stream that stays silent for 5 minutes as stalled; for `google` models with `service_tier: flex` it waits 15 minutes instead. Other tiers and providers keep the 5 minute limit. Rate-limit (`429`) and overloaded (`503`) responses are retried like any other Gemini request; Docker Agent never upgrades a request to a different tier on your behalf.

> [!WARNING]
> Docker Agent's cost estimates do not adjust for `service_tier`. They use the catalogue's standard pricing, which overestimates `flex` and underestimates `priority` charges. Set the model's [`cost` override](../../configuration/models/index.md#custom-token-pricing) to the applicable input, output, and cache token rates for your tier.

See [`examples/gemini_service_tier.yaml`](https://github.com/docker/docker-agent/blob/main/examples/gemini_service_tier.yaml) for a complete example.

## Vertex AI Model Garden

You can use non-Gemini models (e.g. Claude, Llama) hosted on Google Cloud's
[Vertex AI Model Garden](https://cloud.google.com/vertex-ai/generative-ai/docs/partner-models/use-partner-models)
through the `google` provider. When a `publisher` is specified in `provider_opts`,
requests are routed through the appropriate Vertex AI endpoint instead of the
Gemini SDK:

- **Anthropic Claude** (`publisher: anthropic`) uses the Anthropic-native
  `:rawPredict` / `:streamRawPredict` endpoints. Claude models on Vertex AI do
  not support the OpenAI `/chat/completions` path.
- **Other publishers** (e.g. `meta`, `mistral`) use Vertex AI's
  OpenAI-compatible `/chat/completions` endpoint.

For Claude, the SDK puts the model in the URL and `anthropic_version` in the
request body. ADC OAuth authenticates the request; `token_key` does not replace
ADC, and the client clears any inherited Anthropic `X-Api-Key` header.

`native_compaction` and `cache_diagnostics` are rejected by the current Vertex
client. Anthropic documents Google Cloud beta support for on-demand compaction;
that service feature is **not integrated** here. Structured output also requires
an explicit organization-policy allowance. Mythos 5/5.1 have limited-access
announcements but no dedicated Google card in the 2026-10-13 audit; this proves
neither unavailability nor account entitlement. See [Claude access paths](../anthropic/index.md#claude-access-paths)
for model IDs, evidence and policy prerequisites.

Title, sampling and compaction clones keep the selected Model Garden transport.
Claude requests remain on the Anthropic-native Vertex endpoints, with the model
in the URL and `anthropic_version: vertex-2023-10-16` in the body. Project,
location, regional/global host selection and Google OAuth authentication survive
cloning; clones never switch to the direct Claude API or send its API key.

### Authentication

Vertex AI uses Google Cloud Application Default Credentials (ADC). Make sure you
are authenticated:

```bash
gcloud auth application-default login
```

Go embedders can instead supply a request-context token source through
`vertexai.NewClientWithTokenSource` or `options.WithTokenSource`. Unlike
`token_key`, that source replaces ADC, including credential-file and metadata
lookups, and is preserved across clones. HTTP transport wrappers are preserved on both authentication paths;
ADC still applies OAuth authorization around the wrapped transport.

### Configuration

Use the location listed in the specific model's vendor documentation. Anthropic
specifies global or multi-region endpoints for models newer than Sonnet 4.6;
do not reuse an older model's single-region example.

```yaml
models:
  claude-on-vertex:
    provider: google
    model: claude-sonnet-5
    provider_opts:
      project: my-gcp-project       # GCP project ID (or set GOOGLE_CLOUD_PROJECT)
      location: global              # Use a location supported by this model (or GOOGLE_CLOUD_LOCATION)
      publisher: anthropic           # Model publisher (anthropic, meta, etc.)
```

| Option      | Description                                                                          |
| ----------- | ------------------------------------------------------------------------------------ |
| `project`   | GCP project ID. Falls back to `GOOGLE_CLOUD_PROJECT` env var                         |
| `location`  | Model-supported location (global, multi-region or region). Falls back to `GOOGLE_CLOUD_LOCATION` |
| `publisher`  | Model publisher (e.g. `anthropic`, `meta`, `mistral`). Must not be `google`          |

> [!NOTE]
> **Gemini models on Vertex AI**
>
> Setting `publisher: google` (or omitting `publisher`) uses the native Gemini SDK path. The Model Garden endpoint is only used for non-Google publishers.
