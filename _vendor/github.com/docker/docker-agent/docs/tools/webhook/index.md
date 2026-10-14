---
title: "Webhook Tool"
description: "Reliable outbound notifications to Slack, Discord, Telegram, IFTTT, and more."
keywords: docker agent, ai agents, tools, toolsets, webhook, slack, discord, telegram, ifttt, notifications
linkTitle: "Webhook"
weight: 145
canonical: https://docs.docker.com/ai/docker-agent/tools/webhook/
---

_Reliable outbound notifications to Slack, Discord, Telegram, IFTTT, and more._

## Overview

The webhook toolset delivers a notification to a destination **you configure**. The
agent supplies only the message text: it never sees or chooses the URL, because a
webhook URL is itself a credential (Slack and Mattermost embed a secret path,
Discord a token, IFTTT a key, Telegram a bot token).

This is not a general HTTP client — that is the [`api`](../api/index.md) toolset.
The webhook toolset owns *delivery*:

- **At-least-once delivery.** Transient failures (`429`, `5xx`, network errors) are
  retried with exponential backoff, honouring the server's `Retry-After`. A `4xx`
  is permanent and fails immediately without wasting retries.
- **Non-blocking.** The call returns as soon as the notification is queued, so a
  slow or retrying endpoint never stalls the agent's turn. The agent is messaged
  back **only if delivery ultimately fails**.
- **Storm protection.** An identical message to the same destination inside a short
  window is suppressed, and notifications are rate limited, so a looping agent
  cannot flood a channel.
- **Provider-shaped payloads.** Each service's wire format is applied for you.

## Configuration

The destination lives in `webhook_config`. Use `${env.VAR}` for anything secret —
values are expanded at call time and never stored in the config file.

```yaml
toolsets:
  - type: webhook
    webhook_config:
      provider: slack
      url: ${env.SLACK_WEBHOOK_URL}
```

| Field | Required | Description |
| --- | --- | --- |
| `url` | Yes | Webhook endpoint. Usually embeds a secret — prefer `${env.VAR}`. |
| `provider` | No | Payload shape (default `generic`). |
| `headers` | No | Extra headers, for endpoints authenticating with a token. |
| `chat_id` | No | Destination chat — required for `provider: telegram`. |

`timeout` on the toolset (seconds) overrides the per-request HTTP timeout.

## Providers

| Provider | Payload sent | Where the secret lives |
| --- | --- | --- |
| `slack`, `mattermost`, `rocketchat`, `googlechat`, `teams`, `generic` | `{"text": message}` | secret webhook URL |
| `discord` | `{"content": message}` | token in the webhook URL |
| `ifttt` | `{"value1": message, "value2": …, "value3": …}` | key in the webhook URL |
| `telegram` | `{"chat_id": …, "text": message}` | bot token in the URL, plus `chat_id` |

Aliases are accepted: `msteams`/`microsoft_teams` → `teams`, `google_chat`/`gchat`
→ `googlechat`, `rocket.chat` → `rocketchat`.

### Per-service examples

```yaml
# Slack / Mattermost / Rocket.Chat — the URL is the credential
toolsets:
  - type: webhook
    webhook_config:
      provider: slack
      url: ${env.SLACK_WEBHOOK_URL}
```

```yaml
# Discord — the token is part of the webhook URL
toolsets:
  - type: webhook
    webhook_config:
      provider: discord
      url: ${env.DISCORD_WEBHOOK_URL}
```

```yaml
# Telegram — bot token in the URL, chat_id selects the destination chat
toolsets:
  - type: webhook
    webhook_config:
      provider: telegram
      url: https://api.telegram.org/bot${env.TELEGRAM_BOT_TOKEN}/sendMessage
      chat_id: "123456789"
```

```yaml
# IFTTT — the key is part of the trigger URL
toolsets:
  - type: webhook
    webhook_config:
      provider: ifttt
      url: https://maker.ifttt.com/trigger/build_failed/with/key/${env.IFTTT_KEY}
```

```yaml
# Generic endpoint authenticating with a bearer token
toolsets:
  - type: webhook
    webhook_config:
      provider: generic
      url: https://alerts.example.com/notify
      headers:
        Authorization: Bearer ${env.ALERTS_TOKEN}
```

## `send_webhook`

| Parameter | Required | Description |
| --- | --- | --- |
| `message` | Yes | The message text to deliver. |
| `value2`, `value3` | No | Extra IFTTT data fields (`provider: ifttt`). |

Returns immediately once queued. On success nothing further happens; if delivery
ultimately fails, the agent receives a message saying so.

## Example

```yaml
agents:
  root:
    model: openai/gpt-5-mini
    instruction: If a check fails, notify the team with send_webhook.
    toolsets:
      - type: webhook
        webhook_config:
          provider: slack
          url: ${env.SLACK_WEBHOOK_URL}
```

> [!NOTE]
> Requests to non-public addresses are refused (the SSRF-safe HTTP client), and the
> configured URL is never echoed back to the model or into error messages.
