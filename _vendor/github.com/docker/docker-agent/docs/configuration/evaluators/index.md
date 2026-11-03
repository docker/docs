---
title: "Evaluators"
description: "Reusable provider-backed assessments, separate from chat models and decision policies."
keywords: docker agent, evaluators, typesafe, jev, openai, decisions, gateway, classification, tool guards
weight: 65
canonical: https://docs.docker.com/ai/docker-agent/configuration/evaluators/
---

Evaluators assess input against fixed criteria and return typed results. They do
not generate chat messages or execute tools. Consumers decide what an assessment
means: an evaluator reports a probability; a tool guard decides whether to ask
for approval.

Providers include [TypeSafe's Jev](https://docs.typesafe.ai/) and
[OpenAI Decisions](https://developers.openai.com/api/docs/guides/decisions). Named evaluators
are shared across agents in a loaded team. Imported agents keep their source
configuration’s evaluator bindings, even when the parent uses the same names. Configurations using this feature
require version `16` or an omitted version (latest).

## Define an evaluator

```yaml
evaluators:
  credential_exposure:
    provider: typesafe
    model: jev-latest
    type: boolean
    instructions: Does the operation in tool_input disclose credentials outside a trusted boundary?
    timeout: 3s
```

| Field | Meaning |
| --- | --- |
| `provider` | Required backend type (`typesafe` or `openai`) or named entry in `providers`. |
| `model` | Required provider model ID, such as `jev-latest` or `gpt-6-luna`. Pin a versioned ID for reproducible evaluations. |
| `type` | Required: `boolean`, `choice`, or `score`. |
| `instructions` | Required assessment question or rubric instructions. |
| `choices` | For `choice`: map of 2–255 outcome keys to descriptions. |
| `levels` | For `score`: 2–10 descriptions ordered from lowest to highest. |
| `base_url` | Optional API base URL. TypeSafe defaults to `https://api.typesafe.ai` with `/v1/systemone`; OpenAI defaults to `https://api.openai.com/v1` with `/decisions`. |
| `endpoint` | Optional exact HTTP(S) request URL. Overrides `base_url`, including a named provider's default; no path is appended. Credentials, query strings, and fragments are not allowed. |
| `token_key` | Environment variable containing the API key; defaults to `TYPESAFE_API_KEY` for TypeSafe or `OPENAI_API_KEY` for OpenAI. Used only for direct calls. |
| `bypass_models_gateway` | Connect directly even when a gateway is configured; defaults to `false`. |
| `timeout` | Request timeout as a duration such as `3s`; defaults to `10s`. |
| `cost` | Optional USD prices per million tokens: `input` and `output`. Overrides automatic pricing; `cost: {}` explicitly declares free evaluations. |

Connection defaults can be shared through a named provider:

```yaml
providers:
  assessments:
    provider: typesafe
    token_key: CORPORATE_TYPESAFE_KEY

evaluators:
  complexity:
    provider: assessments
    model: jev-latest
    type: choice
    instructions: Classify the reasoning needed to answer this request.
    choices:
      simple: Routine lookup, extraction, or localized editing.
      complex: Multi-step reasoning, architectural analysis, or difficult debugging.
      unknown: Not enough context to assess.
```

Evaluator-level `base_url` and `token_key` override provider defaults. Chat-only
provider settings do not apply; `api_type` and `auth` are rejected for evaluator
providers. Credentials come from the normal environment provider, including
configured secret sources when calling directly.

## Gateway routing

A configured models gateway (`--models-gateway`, `DOCKER_AGENT_MODELS_GATEWAY`,
or user configuration) handles **all evaluator calls**: tool guards, agent routing,
and evaluation judges, including imported agents. No evaluator API keys are needed
on the client for gateway-routed calls. Trusted Docker gateways use the Docker
login token, refreshed per request; a rejected token can be refreshed and retried
once. Docker tokens and encrypted agent configurations are never forwarded to
third-party gateways. Requests carry the current session and install identifiers,
provider/model metadata, and gateway query parameters through the shared transport.

The gateway must support the native `POST /v1/systemone` (TypeSafe) or
`POST /v1/decisions` (OpenAI) protocol and have its own upstream credentials.
The Docker AI gateway additionally requires the corresponding `typesafe` or
`openai-decisions` capability to be enabled, with applicable model/prompt policies.
A gateway rejection is an error; clients never silently fall back to a direct call.
Without an upstream gateway, recording captures the built-in provider endpoints;
private evaluator origins stay direct rather than failing because the recorder
does not recognize them. Replay (`--fake`) always stays offline, including
bypassed and private evaluators. Missing cassette interactions fail rather than
connecting to a provider, and replay does not look up evaluator API keys.

Custom `base_url` and exact `endpoint` settings do **not** implicitly bypass the
gateway. Their upstream target/path is forwarded through it, so the gateway must
support that destination. For a local evaluator or a private deployment not served
by the gateway, set `bypass_models_gateway: true`; its API key is then required
and is included in credential preflight and sandbox forwarding.

## OpenAI Decisions

Use the same types, policies, and hooks with `provider: openai`:

```yaml
evaluators:
  credential_exposure:
    provider: openai
    model: gpt-6-luna
    type: boolean
    instructions: Does this operation disclose credentials outside a trusted boundary?
```

Boolean questions use `predicate`, choices use a deterministic array of string
values/descriptions, and scores use zero-indexed levels. String state is sent as
text; objects and arrays are serialized to JSON text in `input`, retaining the
consumer's evidence structure without treating it as native conversation messages.
Answers are normalized to the same evaluator result types as TypeSafe. Refusals,
missing answers, invalid distributions, and mismatched types are errors, not safe
or zero-valued assessments.

See the [OpenAI evaluator example](https://github.com/docker/docker-agent/blob/main/examples/evaluators-openai.yaml).

In HCL, use `evaluator "name" { ... }` for a top-level named evaluator.

## Evaluation judges

Use an evaluator as an alternative to the chat judge in `docker agent eval`:

```bash
$ docker agent eval agent.yaml --judge-type evaluator
$ docker agent eval agent.yaml --judge-type evaluator --judge-model relevance_judge
```

The first command defaults to `typesafe/jev-latest`. The second selects a named
boolean evaluator from the agent configuration. Its instructions receive
`transcript` and `criterion` fields and must assess criterion satisfaction with
positive polarity. Probabilities of at least 0.5 pass. See
[choosing a judge](../../features/evaluation/index.md#choosing-a-judge)
for configuration, reporting, and limitations.

## Routing agents

A `choice` evaluator can also select the next agent through a `routing_policy` on
a `before_agent_run` or `after_agent_complete` hook. Your `instructions` and `choices`
stay authoritative: routes are never added to the assessment. See
[agent routing hooks](../hooks/index.md#agent-routing-hooks) for the selector
contract, fallback rules, and [`examples/hook_routing.yaml`](https://github.com/docker/docker-agent/blob/main/examples/hook_routing.yaml).

## Compatible endpoints: Laya on Baseten

[Laya](https://huggingface.co/convaiinnovations/laya) can use the `typesafe`
evaluator backend because it returns Jev-compatible assessments. For a Baseten
predict deployment, set `endpoint` on the evaluator so `/v1/systemone` is not
appended. Replace the example URL with your deployment's endpoint:

```yaml
evaluators:
  tool_risk:
    provider: typesafe
    model: english
    endpoint: https://model-YOUR_MODEL_ID.api.baseten.co/development/predict
    token_key: BASETEN_API_KEY
    type: choice
    instructions: Classify the operation in tool_name and tool_input.
    choices:
      read_only: Clearly read-only, with no credential access or external data transfer.
      risky: Changes data, executes code, accesses credentials, or transfers data externally.
      unknown: Insufficient information to establish the operation's effects.
    timeout: 30s
```

The client sends `Authorization: Bearer`, which Baseten accepts; no username is
needed. Laya Router deployments accept checkpoint names such as `english`, not
necessarily the Hugging Face repository ID `convaiinnovations/laya`. Confirm the
accepted names and payload format with your deployment: custom Baseten handlers
can expose a different API. A standard `laya-serve` server at `/v1/systemone` can
instead use `base_url` without an endpoint override.

Existing tool-guard policies and accounting work unchanged. Baseten deployments
do not inherit TypeSafe pricing; cost stays unknown unless explicit pricing and
valid usage are available. Set an appropriate timeout for cold starts. The English
checkpoint has a small default input budget, and probabilities and guard thresholds
need validation on your own data; protocol compatibility is not equivalent safety.

See the [Laya tool-guard example](https://github.com/docker/docker-agent/blob/main/examples/evaluators-laya.yaml).

## Result types

- **Boolean:** `probability` is the probability the statement is true. Mapped to
  TypeSafe's `noul` or OpenAI's `predicate` primitive; no confidence value is invented.
- **Choice:** `choice` identifies a highest-probability outcome, and
  `probabilities` contains the distribution across configured choices.
- **Score:** `score` is the expected zero-based level index. For three levels its
  range is 0–2, not 0–1. Probability keys are `"0"`, `"1"`, and `"2"`.

Results also carry the returned model ID, token usage, an optional estimated USD
`cost`, and optional provider confidence. Confidence is distinct from outcome
probability, and neither should be assumed comparable across providers. Missing
or invalid required answer fields are errors, not zero-valued assessments.

The Go API is `evaluator.Evaluator.Evaluate(ctx, state)`. State can be a string,
JSON object, or array. Loaded teams expose named clients through `Team.Evaluator`.
Each evaluation sends one question. Provider failures are not automatically
retried; gateway authentication may be refreshed and retried once on HTTP 401.

## Pricing and accounting

For the official TypeSafe endpoint, the returned model ID `jev-1.13.0` has
[documented pricing](https://docs.typesafe.ai/models) of **$0.042 per million input
tokens**, with output tokens free. Automatic pricing uses that exact returned ID,
not the requested alias: `jev-latest` is priced only if it resolves to a known
version. OpenAI Decisions `gpt-6-luna` uses the published base input rate of
**$0.10 per million tokens**, with output free; it never uses chat pricing. These
rates also apply when the official upstream is reached through a gateway. Unknown
or future models and custom endpoints have no assumed price. Regional processing
and long-context modifiers may increase the OpenAI charge; use a `cost` override
for applicable rates.

Set an evaluator-level override for private deployments, negotiated rates, or
models without built-in pricing:

```yaml
evaluators:
  credential_exposure:
    provider: typesafe
    model: jev-1.13.0
    type: boolean
    instructions: Does this disclose credentials outside a trusted boundary?
    cost:
      input: 0.042
      output: 0
```

All prices must be finite and nonnegative. Omitted rates in a supplied `cost`
object are zero; `cost: {}` means explicitly free, not unknown. The shared cost
configuration also accepts `cache_read` and `cache_write`, but evaluators do not
currently report cached tokens, so those rates are unused. Clients snapshot their
pricing overrides when constructed.

`Result.Cost` is a `*float64` in USD: `nil` means the charge is unknown; a pointer
to zero means a known zero charge. Estimates require both valid reported token
counts and known pricing. They are not invoices: discounts, minimum charges,
unreported work, and provider billing adjustments may differ. Even free pricing
cannot establish a charge when usage is missing.

Go consumers can attach a request-scoped callback with
`evaluator.WithUsageObserver(ctx, func(evaluator.UsageRecord))`. The callback runs
synchronously once for each attempted HTTP request, including responses whose answers fail
validation or whose HTTP status is an error. Records carry the returned model ID
(or the requested ID when no usable ID is returned), `Usage *Usage`, and
`Cost *float64`. Malformed, missing, null, or incomplete usage is represented by
`Usage == nil`; explicitly reported zero input and output counts remain non-nil.
Transport failures and unreadable or malformed responses produce an unknown-usage
record. Local validation and credential failures before sending a request do not. Callbacks shared
across concurrent evaluations must synchronize their own state; a child context
replaces its inherited observer rather than adding another callback. All callbacks
must finish before `Evaluate` returns. Custom evaluators that issue multiple
requests must report each request separately; all records count toward budgets.

`Result.Usage` remains a value for compatibility, so use observation to distinguish
missing usage from explicit zeros. Observation also captures billable tokens when
`Evaluate` returns an error instead of a result. Consumers must not count the
same request again from its result.

Tool-guard evaluations are recorded as separate session items, including calls
that allow, ask, deny, or return an invalid answer. Known costs contribute to
session totals and the cost display; `/cost` includes a **By Evaluator** breakdown
and explicitly marks unknown costs. Records survive save/reload, branching, and
session export. Runtime events expose each assessment as `evaluation_usage`;
reported token counts also reach telemetry.

Evaluator input/output tokens and known costs count against run and named budgets
for the calling agent. They do not change chat context-window usage or trigger
chat compaction. Budgets are checked before an evaluation and after accounting;
a reached limit blocks the guarded tool and stops the run, never bypassing the
guard. Unknown spend produces a warning and marks cost budgets incomplete rather
than pretending the call was free. Invalid accounting from custom evaluators
(negative or non-finite costs, negative tokens, or overflowing token totals)
is recorded as unknown for the invalid fields and blocks the guarded tool.
Consumption counters saturate at their numeric limits instead of wrapping;
cost budgets whose shared totals overflow are marked incomplete.

These limits are best-effort, not billing caps: a single evaluation can cross a
limit, concurrently admitted requests can finish after it is reached, and a
provider may charge for a timed-out request without reporting usage. Such
unreported charges cannot be added to the numerical cost or token totals.

## Tool guards

The first built-in consumer is a `type: evaluator` hook on `tool_guard`.
Boolean and choice evaluators are supported; score results are currently
available to Go consumers only. There is no evaluator-based model routing yet.

```yaml
evaluators:
  credential_exposure:
    provider: typesafe
    model: jev-latest
    type: boolean
    instructions: Does the operation in tool_input disclose credentials outside a trusted boundary?

agents:
  root:
    model: openai/gpt-5-mini
    instruction: Help with the project and respect tool policy decisions.
    toolsets:
      - type: shell
    hooks:
      tool_guard:
        - matcher: shell
          hooks:
            - type: evaluator
              evaluator: credential_exposure
              evaluator_policy:
                decisions:
                  "true": ask
                  "false": allow
                min_probability: 0.95
                fallback: ask
```

The guard sends only `tool_name`, `tool_input`, and `tool_category`, after input
transforms. It does not send conversation history or session identifiers.

For a choice, the policy looks up the selected outcome's probability. For a
boolean, it selects `true` when the probability is at least 0.5, otherwise `false`
with probability `1 - p`. A mapped outcome at or above `min_probability` uses its
configured decision. Unmapped or uncertain outcomes use `fallback`, which must
be `ask` or `deny`. The threshold above is illustrative: tune it on representative
and adversarial cases rather than treating it as a safety guarantee.

- **`allow` is advisory:** existing permissions and safety-mode checks still run.
- **`ask` requires fresh confirmation**, even under autonomous mode or an earlier
  “always allow” grant. Explicit denials still win. Non-interactive sessions deny.
- **`deny` blocks the call.**
- Provider errors, timeouts, or malformed answers **block**, regardless of the
  uncertainty fallback or hook `on_error` setting.

Confirmation metadata includes the evaluator name, selected outcome, probability,
and returned model. Existing [tool-guard semantics](../hooks/index.md#tool-phases-transform-guard-approve)
apply, including rechecking arguments rewritten by a later legacy hook.

## Security and limitations

This feature is opt-in and sends tool arguments to an external service, including
calls that permission rules may subsequently deny. Use trusted endpoints, minimize
input, and review redaction requirements; automatic redaction cannot detect every
secret. Redirects are disabled, and provider errors omit input and response bodies.

Jev's [documented adversarial-input limitations](https://docs.typesafe.ai/model-jaggedness/jev-1.13)
mean that an assessment is not a security boundary. A low-risk label must not
replace sandboxing, deterministic restrictions, or human review.

Embedders using strict team loading must explicitly enable `config.FeatureEvaluators`
and, for evaluator hooks, `config.FeatureHooks`.

See the runnable [evaluator tool-guard example](https://github.com/docker/docker-agent/blob/main/examples/evaluators.yaml).
