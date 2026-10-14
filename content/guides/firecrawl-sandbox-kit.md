---
title: Give sandboxed agents live web access with Firecrawl
linkTitle: Web access with Firecrawl
description: Add the Firecrawl kit to a Docker Sandbox so the agent can search the web, scrape pages to Markdown, and crawl sites through the sandbox proxy, with your API key kept on the host.
summary: |
  Compose the Firecrawl v3 mixin onto a sandbox workload. The agent gets the
  Firecrawl Python SDK and a network rule for the Firecrawl API, and the
  sandbox proxy adds your API key to requests, so the key never enters the
  sandbox.
keywords: ai, firecrawl, docker sandboxes, sbx, kits, mixin, web scraping, web search, crawl, credentials, agents
params:
  tags: [ai]
  time: 15 minutes
---

Agents in a Docker Sandbox run under a network policy that usually blocks most
of the web. That keeps the agent contained, but it also means the agent
can't look up documentation or read a page you point it at.
[Firecrawl](https://www.firecrawl.dev/) is a web data API that searches the
web, scrapes pages to clean Markdown, and crawls whole sites. The Firecrawl
kit gives a sandboxed agent that capability through a single API host,
without opening the rest of the web.

In this guide, you'll learn how to:

- Store your Firecrawl API key on the host as a sandbox secret
- Add the Firecrawl mixin to a v3 workload, alone or with a coding agent
- Check that the SDK, the proxy-managed credential, and the network rule work
- Call Firecrawl from inside the sandbox
- Understand the limits the sandbox network policy places on the kit

## What the kit adds

The Firecrawl kit is a v3 [mixin](/manuals/ai/sandboxes/customize/author/tool-mixins.md)
published as `docker.io/firecrawl/sbx-kit-firecrawl`. When you add it to a
workload, it contributes four things:

| Capability                       | What it does                                                                                                                         |
| -------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------ |
| `com.docker.sandbox/lifecycle@1` | An install hook that installs a pinned `firecrawl-py` release as the agent user. Sandbox creation fails if the install fails.         |
| `com.docker.sandbox/credential@1` | Requests the key stored under the `firecrawl` service. The proxy adds it as `Authorization: Bearer <key>` on requests to `api.firecrawl.dev`. |
| `com.docker.sandbox/network-policy@1` | Allows `pypi.org` and `files.pythonhosted.org` during install, and `api.firecrawl.dev` at runtime.                              |
| `com.docker.sandbox/agent-context@1` | Adds instructions that tell the agent the SDK is available, how to call it, and which sandbox limits apply.                     |

Inside the sandbox, `FIRECRAWL_API_KEY` holds the placeholder value
`proxy-managed`. The SDK sends that placeholder, and the proxy on your host
replaces it with your real key on the way out. Requests to any other host
never carry the key.

The install hook runs against the workload's own Python, so the kit works with
any workload that ships `python3` and `pip`. Docker's sandbox workloads include
both.

## Prerequisites

Before you start, make sure you have:

- [Docker Sandboxes (`sbx`) installed and signed in](../manuals/ai/sandboxes/install.md),
  version 0.46 or later
- A global network policy. If you haven't set one, run
  `sbx policy init balanced`
- A [Firecrawl API key](https://www.firecrawl.dev/app/api-keys). Keys start
  with `fc-`

## Step 1: Store the Firecrawl key on the host

The kit declares a credential for the `firecrawl` service. Store your key under
that name:

```console
$ sbx secret set firecrawl
```

`sbx` prompts for the value. Service secrets apply to every sandbox by
default. Use `--sandbox` to scope the secret to one sandbox instead.

Storing the secret doesn't give the kit access to it. The first time you run
the kit, `sbx` asks you to approve sending the `firecrawl` credential to
`api.firecrawl.dev`. For how approvals work, and how to approve in advance for
unattended runs, see
[Credential bindings](../manuals/ai/sandboxes/configuration/credentials.md#credential-bindings).

## Step 2: Run a workload with the Firecrawl mixin

V3 mixins compose onto v3 workloads. Start with Docker's shell workload, which
opens a login shell so you can try the SDK directly. From your project
directory, run:

```console
$ sbx run docker.io/docker/sbx-kit-shell:1.0.0 --name firecrawl-shell \
    --kit docker.io/firecrawl/sbx-kit-firecrawl:1.0.0
```

Approve the credential request when prompted. `sbx` pulls both kits, runs the
install hook, and opens a shell in the sandbox.

To give a coding agent web access, add the agent's mixin alongside Firecrawl.
For example, to run Claude Code on the shell workload:

```console
$ sbx run docker.io/docker/sbx-kit-shell:1.0.0 --name firecrawl-claude \
    --kit docker.io/docker/sbx-kit-claude-mixin:2.1.285 \
    --kit docker.io/firecrawl/sbx-kit-firecrawl:1.0.0
```

The Claude Code mixin also needs its own credential. See
[Build a tool mixin](/manuals/ai/sandboxes/customize/author/tool-mixins.md#try-the-mixin).
Other v3 workloads, such as `docker.io/docker/sbx-kit-codex`, accept the
Firecrawl mixin the same way. For more on combining kits, see
[Use kits](/manuals/ai/sandboxes/customize/use-kits.md#add-mixins).

> [!NOTE]
> Built-in shortcuts such as `sbx run claude` use v2 kits, and v2 agents only
> accept v2 mixins. Firecrawl also publishes a v2 form of the kit for those
> agents. See the
> [kit repository](https://github.com/firecrawl/firecrawl-docker-sandbox#2-launch-the-sandbox-with-the-kit)
> for v2 usage.

## Step 3: Check the kit

Run these checks from the shell in the `firecrawl-shell` sandbox. Each one
covers a separate part of the kit.

Check that the SDK is installed at the pinned version, in the agent user's
site-packages:

```console
$ python3 -c "import firecrawl, importlib.metadata as m; print(m.version('firecrawl-py'), firecrawl.__file__)"
4.44.0 /home/agent/.local/lib/python3.14/site-packages/firecrawl/__init__.py
```

Check that the key in the sandbox is the proxy placeholder, not your real key:

```console
$ echo "$FIRECRAWL_API_KEY"
proxy-managed
```

Then scrape a page to test the whole path: the SDK, the network rule, and
credential injection. This request uses Firecrawl credits:

```console
$ python3 - <<'EOF'
from firecrawl import Firecrawl

fc = Firecrawl()
doc = fc.scrape("https://docs.docker.com/ai/sandboxes/", formats=["markdown"])
print(doc.markdown[:300])
EOF
```

The output is the start of the page as Markdown. If the request fails with
`Unauthorized`, the proxy didn't add your key. Check that you stored the
secret under `firecrawl` and approved the kit's request. A
`PaymentRequiredError` means the request authenticated and the account is out
of credits.

## Step 4: Use Firecrawl from the agent

The kit's agent instructions point the agent at the SDK, so a coding agent can
use Firecrawl when a task needs current information or the contents of a URL.
The main calls are:

```python
from firecrawl import Firecrawl

fc = Firecrawl()  # The proxy adds the API key

fc.scrape("https://example.com", formats=["markdown"])  # One page as Markdown
fc.search("docker sandboxes kits", limit=5)             # Web search results
fc.crawl("https://docs.example.com", limit=20)          # Pages across a site
```

Each crawled page uses credits, so keep `limit` low and prefer a search
followed by a targeted scrape when the agent needs a single fact. For the full
API, including structured extraction with a schema, see the
[Firecrawl Python SDK documentation](https://docs.firecrawl.dev/sdks/python).

## How the kit interacts with network policy

Kit network rules add hosts to your sandbox policy. They don't open the web.
Keep these limits in mind, because they can look like Firecrawl errors:

- Under the `balanced` policy, most websites are blocked, so `curl` or
  `requests` to a URL from a search result fails with a proxy `403`. Scrape
  the URL through Firecrawl instead.
- Under a deny-all policy, `api.firecrawl.dev` is the only web host the agent
  can reach.
- Screenshot and other media formats return links on a Firecrawl storage host
  that isn't in the kit's allow list, so the agent can't download those files.
- If your organization manages sandbox policy centrally, kit rules can't grant
  access beyond it. An administrator must allow `api.firecrawl.dev`.

To see which rule handled a request, run this on the host:

```console
$ sbx policy log
```

## Clean up

Sandboxes persist after you exit. To remove the sandboxes from this guide:

```console
$ sbx rm firecrawl-shell firecrawl-claude
```

Files in your workspace are unaffected. To remove the stored key, run
`sbx secret rm firecrawl`.

## Learn more

- [Use kits](/manuals/ai/sandboxes/customize/use-kits.md)
- [Build a tool mixin](/manuals/ai/sandboxes/customize/author/tool-mixins.md)
- [Credentials](../manuals/ai/sandboxes/configuration/credentials.md)
- [Firecrawl kit source](https://github.com/firecrawl/firecrawl-docker-sandbox)
- [Firecrawl documentation](https://docs.firecrawl.dev/)
