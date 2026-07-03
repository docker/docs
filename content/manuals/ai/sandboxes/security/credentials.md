---
title: Credentials
weight: 20
description: How Docker Sandboxes handle API keys and authentication credentials for sandboxed agents.
keywords: docker sandboxes, credentials, api keys, authentication, proxy, ssh agent, secrets
---

Most agents need an API key for their model provider. An HTTP/HTTPS proxy on
your host intercepts outbound requests from the sandbox, looks up the matching
credential on the host, and overwrites the auth header before forwarding. The
real credential stays on the host; the sandbox sees only a sentinel value. For
the security model behind this, see
[Credential isolation](isolation.md#credential-isolation).

## How credential injection works

When a sandbox makes an outbound request, the host-side proxy decides three
things: whether the request **matches** a service the kit (or built-in agent)
declares, what **header** to write, and what **value** to inject. The kit
declares the match and the header; you provide the value on the host. The real
value never enters the sandbox — the agent sees only a sentinel like
`proxy-managed`.

There are several ways to provide that value. When more than one source has a
value for the same service, the stored secret takes precedence.

| Form                                                                        | What it is                                                   | Use it when                                                                                                      |
| --------------------------------------------------------------------------- | ------------------------------------------------------------ | ---------------------------------------------------------------------------------------------------------------- |
| [Stored secrets](#stored-secrets) (`sbx secret set`)                        | A value in your OS keychain, keyed by service                | The default for any built-in or kit-declared service                                                             |
| [Custom secrets](#custom-secrets) (`sbx secret set-custom`)                 | A value keyed to a domain and environment variable           | The service model doesn't fit — the agent validates the variable's format, or the secret rides in a request body |
| OAuth                                                                       | A host-side sign-in flow; the token never enters the sandbox | The agent supports it, such as Claude Code, Codex, Cursor, or Droid                                              |
| [Credential bindings](#credential-bindings) (`credentials.yaml`)            | Per-service mechanism and domain approval                    | Required for third-party `schemaVersion: "2"` kits                                                               |
| [Registry credentials](#registry-credentials) (`sbx secret set --registry`) | Authentication for pulling images and kits                   | Pulling templates or kits from a private registry                                                                |

For multi-provider agents (OpenCode, Docker Agent), the proxy selects
credentials based on the API endpoint being called. See individual
[agent pages](../agents/) for provider-specific details.

## Stored secrets

`sbx secret set` stores credentials in your OS keychain, keyed on a service
identifier. Built-in agents declare a fixed set of services. Custom kits can
declare their own. The same `sbx secret set` flow works for both.

### Where secrets are stored

The store backing `sbx secret set` depends on your operating system:

- macOS: the system Keychain.
- Windows: the Windows Credential Manager.
- Linux: the Secret Service exposed by your desktop keyring, such as GNOME
  Keyring or KDE Wallet.

The Ubuntu package depends on GNOME Keyring, so a standard desktop install
needs no extra setup.

On Linux hosts without a running Secret Service — headless servers and some
WSL setups — `sbx` falls back to an encrypted file under your user config
directory `$XDG_CONFIG_HOME/com.docker.sandboxes`, which defaults to
`~/.config/com.docker.sandboxes` when `$XDG_CONFIG_HOME` is unset. The fallback
is automatic and needs no configuration. When you store a secret this way,
`sbx` prints a notice:

```text
No keychain detected - this secret will be stored in an encrypted file on disk
```

The file is encrypted at rest and protected by `0700` directory permissions,
the same posture as `~/.docker/config.json`. This is weaker than an OS
keychain, which also mediates access per application. If you start a Secret
Service on the host later, `sbx` stores new secrets in the keychain again. For
more on running sandboxes without a desktop keyring, see
[Can I use Docker Sandboxes on headless Linux?](../faq.md#can-i-use-docker-sandboxes-on-headless-linux)

### Store a secret

```console
$ sbx secret set -g anthropic
```

This prompts you for the secret value interactively. The `-g` flag stores the
secret globally so it's available to all sandboxes. To scope a secret to a
specific sandbox instead:

```console
$ sbx secret set my-sandbox openai
```

> [!NOTE]
> A sandbox-scoped secret takes effect immediately, even if the sandbox is
> running. A global secret (`-g`) only applies when a sandbox is created. If
> you set or change a global secret while a sandbox is running, recreate the
> sandbox for the new value to take effect.

### Import from environment variables

If you already have API keys set in your shell, `sbx secret import` reads them
and stores them in the keychain without typing each value manually:

```console
$ sbx secret import
```

This scans your current session for the environment variables in the
[built-in services table](#built-in-services) below and prompts you to confirm
each one before writing. To import a single service:

```console
$ sbx secret import openai
```

Pass `--all` to import everything without prompting (new entries only; existing
entries are left unchanged), or `--force` to overwrite existing entries:

```console
$ sbx secret import --all
$ sbx secret import openai --force
```

Pass `--dry-run` to preview what would be imported without writing anything.
Run `sbx secret ls` afterwards to confirm what's stored. For setting up
credentials in CI, see [CI and headless use](../workflows.md#ci-and-headless-use).

### Built-in services

Each built-in service name maps to the environment variables `sbx secret import`
reads and the API domains the proxy injects credentials into:

| Service      | Environment variables              | API domains                                                                                                                   |
| ------------ | ---------------------------------- | ----------------------------------------------------------------------------------------------------------------------------- |
| `anthropic`  | `ANTHROPIC_API_KEY`                | `api.anthropic.com`, `console.anthropic.com`, `claude.ai`, `mcp-proxy.anthropic.com`                                          |
| `cursor`     | `CURSOR_API_KEY`                   | `api2.cursor.sh`, `api3.cursor.sh`, `repo42.cursor.sh`, `cursor.com`                                                          |
| `droid`      | `FACTORY_API_KEY`                  | `api.factory.ai`, `app.factory.ai`, `relay.factory.ai`                                                                        |
| `github`     | `GH_TOKEN`, `GITHUB_TOKEN`         | `api.github.com`, `github.com`, `raw.githubusercontent.com`, `gist.github.com`, `copilot.github.com`, `api.githubcopilot.com` |
| `google`     | `GEMINI_API_KEY`, `GOOGLE_API_KEY` | `generativelanguage.googleapis.com`, `oauth2.googleapis.com`, `aiplatform.googleapis.com`, `vertexai.googleapis.com`          |
| `groq`       | `GROQ_API_KEY`                     | `api.groq.com`                                                                                                                |
| `mistral`    | `MISTRAL_API_KEY`                  | `api.mistral.ai`                                                                                                              |
| `nebius`     | `NEBIUS_API_KEY`                   | `api.studio.nebius.com`, `api.tokenfactory.nebius.com`                                                                        |
| `openai`     | `OPENAI_API_KEY`                   | `api.openai.com`, `openai.com`, `chatgpt.com`, `www.chatgpt.com`                                                              |
| `openrouter` | `OPENROUTER_API_KEY`               | `openrouter.ai`                                                                                                               |
| `xai`        | `XAI_API_KEY`                      | `api.x.ai`                                                                                                                    |

When you store a secret with `sbx secret set -g <service>`, the proxy injects
it into requests to the listed API domains.

### Services declared by kits

Custom kits can declare their own service identifiers in `spec.yaml`. In
`schemaVersion: "2"`, credentials are declared under the `credentials:` list:

```yaml
credentials:
  - service: my-service
    apiKey:
      name: MY_SERVICE_TOKEN
      inject:
        - domain: api.my-service.com
          header: Authorization
          format: "Bearer %s"
```

For a known provider (such as `anthropic` or `openai`), `- service: <name>` is
enough — injection config is filled in from the provider registry. To provide
the credential value, run `sbx secret set` with the same identifier the kit
declares:

```console
$ sbx secret set -g my-service
```

There's no separate registration step; the keychain entry is keyed on the
identifier the kit already uses. See
[Authenticate to external services](../customize/kits.md#authenticate-to-external-services)
for the full kit-side wiring.

### List and remove secrets

List all stored secrets:

```console
$ sbx secret ls
SCOPE      TYPE      NAME      SECRET
(global)   service   github    gho_GCaw4o****...****43qy
```

Remove a secret:

```console
$ sbx secret rm -g github
```

> [!NOTE]
> Running `sbx reset` deletes all stored secrets along with all sandbox state.
> You'll need to re-add your secrets after a reset.

### GitHub token

The `github` service gives the agent access to the `gh` CLI inside the
sandbox. Pass your existing GitHub CLI token:

```console
$ echo "$(gh auth token)" | sbx secret set -g github
```

This is useful for agents that create pull requests, open issues, or interact
with GitHub APIs on your behalf.

### SSH agent

If your host has an SSH agent and `SSH_AUTH_SOCK` is set, Docker Sandboxes
forwards the agent into the sandbox and sets `SSH_AUTH_SOCK` there. The
private keys stay on your host. Processes inside the sandbox can request
signatures from the forwarded agent, but they can't read or copy the private
key.

Use SSH agent forwarding for Git operations over SSH and SSH-based commit
signing. The signing key must be loaded in the host SSH agent for sandboxed
commit signing to work. Outbound SSH connections are still subject to sandbox
network policy. For details, see
[Commit signing](../workflows.md#commit-signing).

## Custom secrets

> [!IMPORTANT]
> Custom secrets are experimental. Behavior, flags, and the placeholder format may
> change without notice.

For credentials that don't fit the service-identifier model — for example,
when an agent validates the environment variable format at boot, or when the
credential lands in a request body rather than a header — use
`sbx secret set-custom`. The secret is keyed on one or more target domains, an
environment variable name, and an optional placeholder string, instead of a
service identifier.

```console
$ sbx secret set-custom -g \
    --host api.example.com \
    --env API_KEY \
    --value <secret>
```

Repeat `--host` to cover multiple domains with the same secret — useful when
an API is split across related hostnames or when two unrelated endpoints share
a credential:

```console
$ sbx secret set-custom -g \
    --host api.example.com \
    --host uploads.example.com \
    --env API_KEY \
    --value <secret>
```

A `--host` value can also use wildcards, with the same syntax as
[network rules](../governance/concepts.md#network-rules): `*` matches a
single label (`*.example.com` covers `api.example.com`) and `**` matches any
number (`**.example.com` covers `api.example.com` and `v2.api.example.com`).

> [!WARNING]
> Passing the secret as `--value <secret>` records it in your shell history
> and exposes it to other processes running as your user. Avoid pasting
> real credentials inline — read the value from a variable that's already
> in your environment, and clear shell history if a real secret was passed
> on the command line.

Inside the sandbox, `API_KEY` is set to a generated placeholder (for example,
`sbx-cs-<rand>`). When a sandboxed process sends a request to any of the
configured hosts and the placeholder appears anywhere in the request, the
proxy replaces it with the real value. The agent never sees the real secret.

Prefer the [service-based flow](#stored-secrets) whenever it's an option —
the kit handles the wiring; you only provide the value.

## Credential bindings

A credential bindings file records which credential mechanisms you've approved
for each service and the domains each may be used on. It lives at
`~/.config/sbx/credentials.yaml`, or `%APPDATA%\sbx\credentials.yaml` on
Windows.

Third-party kits that declare `schemaVersion: "2"` require an approved binding
for each credential they use. `sbx` creates one interactively the first time you
run such a kit (see [First-run approval](#first-run-approval)); you can also
write entries by hand. Built-in agents are authorized by provenance and never
need a binding.

Each entry under `bindings` is keyed by a
[service identifier](#built-in-services) and approves one or both credential
mechanisms:

- `apiKey` — approves injecting the service's stored API key. The value comes
  from the [secret store](#stored-secrets) (`sbx secret set <service>`); the
  binding records approval, it doesn't hold or locate the value.
- `oauth` — approves the OAuth flow for the service. You sign in on the host,
  and the proxy handles token refresh and routing.

Each mechanism takes a `domains` list — the domains the proxy may inject the
credential into. The credential is attached only where those domains and the
ones the kit requests overlap.

```yaml
bindings:
  anthropic:
    apiKey:
      domains: [api.anthropic.com]
  github:
    apiKey:
      domains: [api.github.com, github.com]
```

A binding is only an approval record: the presence of `apiKey` or `oauth`
authorizes that mechanism. Declining a credential writes no entry at all.

### First-run approval

When a third-party kit needs a credential that has no binding, `sbx` walks you
through approving one. For each credential, you use the value already in the
secret store or enter one at the prompt, and you approve the domains it may
reach. `sbx` writes the entry to `credentials.yaml`.

In non-interactive contexts (CI or `--detached`), there's no one to answer the
prompt. `sbx` logs a warning naming the service and creates the sandbox anyway
with the credential withheld, so the agent starts unauthenticated. Pre-create
the binding — by running the kit interactively once, or by writing
`credentials.yaml` directly — before running unattended.

This makes the bindings file an allowlist of credential-to-domain approvals: a
kit can use only the credentials and domains you've approved.

### Which kits require a binding

Only third-party kits that declare `schemaVersion: "2"` require a binding.
Built-in agents also use `schemaVersion: "2"` but are authorized by provenance,
so their credentials inject automatically without a binding. Kits on
`schemaVersion: "1"` inject their declared credentials without a binding.

## Registry credentials

Registry credentials authenticate to private OCI registries when pulling
[templates](../customize/templates.md) or [kits](../customize/kits.md), and can
also let the agent pull and push images from inside the sandbox. Use
`sbx secret set --registry <host>` to store them. For Docker Hub, `sbx` reuses
your `sbx login` session — no registry secret needed. For other registries
(GitHub Container Registry, ECR, ACR, self-hosted Nexus, and so on), store
credentials with `sbx secret set --registry`.

The scope you store a credential at controls where it's used — and whether its
value enters the sandbox. The scope comes from how you target `sbx secret set`:

```text
sbx secret set [-g | SANDBOX] --registry HOST
```

- **Host-only** (no `-g`, no `SANDBOX`): the `sbx` CLI uses it to pull templates
  and kits when creating a sandbox. The credential stays on the host and is
  never available inside the sandbox.
- **Global** (`-g`): same as host-only, plus written into
  `~/.docker/config.json` in every new sandbox so the agent can pull and push
  images. The value lives inside the VM, where the agent can read it, so it's
  less isolated than the proxy-injected service credentials above. Use it when
  agents build and publish container images.
- **Sandbox-scoped** (`SANDBOX`): same in-sandbox behavior as global, but only
  for the named sandbox. Use it when only one sandbox needs registry access.

> [!NOTE]
> Registry credentials are written into a sandbox at creation time. Recreate an
> existing sandbox to pick up credentials added after it was created.

### Store registry credentials

Pipe a token from stdin and target the registry hostname:

```console
$ gh auth token | sbx secret set --registry ghcr.io --password-stdin
```

For registries that require a username (for example, ACR with an admin
account), add `--username`:

```console
$ echo "$ACR_PASSWORD" | sbx secret set \
    --registry myregistry.azurecr.io \
    --username myuser \
    --password-stdin
```

Add `-g` to store the credential globally, before you create the sandbox:

```console
$ gh auth token | sbx secret set -g --registry ghcr.io --password-stdin
$ sbx run claude                      # created with the credential in place
```

To scope the credential to a single sandbox, store it under that sandbox's name
and create the sandbox with the same name:

```console
$ gh auth token | sbx secret set my-app --registry ghcr.io --password-stdin
$ sbx run claude --name my-app
```

`sbx kit pull` also uses these credentials, with the Docker credential
store as a fallback. `sbx kit push` uses only the Docker credential store —
push targets still require a prior `docker login`.

### Remove registry credentials

Remove both the host-only and global entries for a registry:

```console
$ sbx secret rm --registry ghcr.io -f
```

To remove only the global (in-sandbox) entry and leave the
host-only credential in place, pass `-g`:

```console
$ sbx secret rm -g --registry ghcr.io -f
```

To remove a sandbox-scoped credential, pass the sandbox name:

```console
$ sbx secret rm my-sandbox --registry ghcr.io -f
```

## Best practices

- Use [stored secrets](#stored-secrets) to provide credentials. They are
  encrypted at rest in the OS keychain (or an encrypted file on Linux hosts
  without a keychain). See [Where secrets are stored](#where-secrets-are-stored).
- Don't set API keys manually inside the sandbox. Sandbox agents are
  pre-configured to use proxy-managed credentials.
- Registry credentials you make available inside a sandbox are stored in the VM
  (`~/.docker/config.json`), where the agent can read them — unlike
  proxy-injected service credentials, which never enter the sandbox. Reserve
  them for sandboxes that need registry access, and prefer sandbox scope over
  global (`-g`) to limit exposure.
- Several agents support OAuth as another secure option: the flow runs on the
  host, so the token is never exposed inside the sandbox. If you haven't stored
  a credential, the agent prompts you to authenticate — Codex prompts on the
  host from `sbx run codex`, while Claude Code, Cursor, and Droid prompt
  interactively inside the sandbox. To authenticate ahead of time, run
  `sbx secret set -g openai --oauth` for Codex or use `/login` inside Claude
  Code; Cursor and Droid have no ahead-of-time option, so their sign-in prompt
  appears when the agent starts. See the individual [agent pages](../agents/)
  for each agent's flow.
- If you store credentials in 1Password, see
  [Sourcing credentials from 1Password](../workflows.md#sourcing-credentials-from-1password)
  for how to use `op read` and `op run` with `sbx`.

## Custom templates and placeholder values

When building custom templates or installing agents manually in a shell
sandbox, some agents require environment variables like `OPENAI_API_KEY` to be
set before they start. Set these to placeholder values (e.g. `proxy-managed`)
if needed. The proxy injects actual credentials regardless of the environment
variable value.
