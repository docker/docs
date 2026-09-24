---
title: "Run agents with kits"
linkTitle: "Run agents with kits"
description: "Discover the available kits, launch an agent sandbox, and run work inside it."
keywords: "cloud sandboxes, sandboxes api, run agents with kits"
weight: 104
params:
  sidebar:
    group: "Get started"
---

A kit supplies an agent's image and configuration, including tools and the settings they need. Launching a kit gives you a repeatable starting point without assembling those settings in every application.

Start with an [authenticated client](connect-to-cloud-with-a-bearer-token.md). Agent workloads also need credentials for their model provider. Docker authentication gives access to sandboxes; it does not sign you in to Anthropic, OpenAI, or another agent provider.

## Discover available kits {#1-discover-available-kits}

List the kits bundled with your installed SDK. Use an entry's name when launching it. Listing the catalog does not create a sandbox.

The SDK's launch-by-name helper includes these starting points:

| Kit name | Environment |
| --- | --- |
| `claude` | Claude Code |
| `codex` | Codex |
| `cursor` | Cursor |
| `devin` | Devin |
| `docker-agent` | Docker Agent |
| `gemini` | Gemini CLI |
| `opencode` | OpenCode |
| `shell` | A shell environment for your own commands |

Use the catalog result as the list for your installed version. Upgrading the SDK can update the bundled kits.

{{< tabs >}}
{{< tab name="TypeScript" >}}

```typescript
return client.kits.list();
```

<details>
<summary>Complete TypeScript example: kits/catalog.ts</summary>

```typescript
import type { Sandboxes } from '@docker/sandboxes';

export function listKits(client: Sandboxes) {
  return client.kits.list();
}
```

</details>

{{< /tab >}}
{{< /tabs >}}

## Launch a kit {#2-launch-a-kit}

To run Claude Code, choose `claude` and attach a stored `anthropic` secret containing your Anthropic API key. [Store the provider credential](get-a-stored-secret-into-a-sandbox.md) first, then pass the stored secret's resource name to the launch example. Keep credentials out of command arguments and plain environment values.

Pass a catalog name, a display name, and any stored secret names the agent needs. The example selects Small (2 vCPUs, 4 GiB) and uses the launch-and-wait helper to return a running sandbox under one deadline. Small is also the default when you omit kit resources. See [compute sizes](work-within-the-limits.md) for the other choices.

Use `launch` when you want the accepted handle immediately and will wait separately. If the wait fails, inspect the accepted sandbox retained by the error before launching another one.

`launchAndWait()` waits for the sandbox to reach the running state. A kit can still be installing tools or cloning a repository at that point. The SDK has no helper that waits for all kit setup to finish.

If your work depends on that setup, check the result it needs before starting. For example, a kit can write a completion marker after cloning a repository, or a service can expose a health check. Poll with a delay between checks and a timeout so failed setup does not leave your application waiting indefinitely.

The launch helper accepts only names returned by the bundled catalog, not community repository URLs or registry references. Account network policies still apply to kit sandboxes.

Use `shell` when you need an environment for scripts, builds, or your own executable. To use a container image directly, see [Run your own container image](run-your-own-container-image.md).

{{< tabs >}}
{{< tab name="TypeScript" >}}

```typescript
return client.kits.launchAndWait(kitName, {
  displayName,
  resources: 'small',
  storage: { secrets },
});
```

<details>
<summary>Complete TypeScript example: kits/launch.ts</summary>

```typescript
import type { Sandboxes } from '@docker/sandboxes';

export async function launchKit(
  client: Sandboxes,
  kitName: string,
  displayName: string,
  secrets: string[] = [],
) {
  return client.kits.launchAndWait(kitName, {
    displayName,
    resources: 'small',
    storage: { secrets },
  });
}
```

</details>

{{< /tab >}}
{{< /tabs >}}

## Run the agent {#3-run-the-agent}

Pass the sandbox handle returned by the launch example and your prompt to the run helper. It runs `claude -p` and waits for the response. Launching a kit alone does not submit an agent task.

The example returns captured output and an exit code. Check both when diagnosing an agent failure. For a long task, [stream output](run-something-that-produces-real-output.md); for a conversation in a terminal, [open an interactive session](run-an-interactive-shell-in-a-cloud-sandbox.md).

Save the sandbox name if you will resume the work later. [Delete the sandbox](delete-a-cloud-sandbox.md) when you no longer need its files or running processes.

{{< tabs >}}
{{< tab name="TypeScript" >}}

```typescript
return sandbox.processes.run(
  { args: ['claude', '-p', prompt] },
  { timeoutMs: 300_000 },
);
```

<details>
<summary>Complete TypeScript example: kits/run.ts</summary>

```typescript
import type { Sandbox } from '@docker/sandboxes';

export async function runKitAgent(sandbox: Sandbox, prompt: string) {
  return sandbox.processes.run(
    { args: ['claude', '-p', prompt] },
    { timeoutMs: 300_000 },
  );
}
```

</details>

{{< /tab >}}
{{< /tabs >}}

## Explore the wider kit catalog {#4-explore-the-wider-kit-catalog}

The [community kit catalog](https://github.com/docker/sbx-kits-contrib) includes agents such as Aider, Amp, Copilot, Kiro, and OpenHands, plus development tools, browser automation, source control, and security scanning. For example, Code Server adds a browser-based editor and Playwright adds browser automation.

Some kits define a complete sandbox environment; others add tools or configuration to an existing agent environment. These are not all bundled SDK launch targets. See the [Docker kits guide](https://docs.docker.com/ai/sandboxes/customize/kits/) for the catalog's usage instructions.
