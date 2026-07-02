---
title: Org policy recipes
linkTitle: Policy recipes
weight: 22
description: Minimal, composable network policy presets for common sandbox workflows, in Admin Console and API form.
keywords: docker sandboxes, governance, organization policy, network policy, allowlist, recipes, presets, Claude Code, Codex, package managers
---

Organization policies start deny-by-default: outbound traffic is blocked unless
an explicit allow rule matches. The recipes on this page are minimal building
blocks for the hosts common workflows depend on — look up the ones you need,
check the hostnames, and copy them into a policy.

Treat each recipe as a baseline, not a complete answer. They're intentionally
narrower than the local [Balanced preset](local.md#default-preset), which allows
a wide range of destinations so individual machines work out of the box. An
organization policy should do the opposite: allow only what your sandboxes
actually need, and remove anything you don't. To find the exact hosts a workflow
needs, run it in a sandbox and watch `sbx policy log` for blocked destinations.
See [Monitoring](monitoring.md).

## Find a recipe

- **Developer essentials**: [GitHub](#github) ·
  [Certificate validation](#certificate-validation) ·
  [Ubuntu packages](#ubuntu-packages)
- **Package registries**: [Node.js and npm](#nodejs-and-npm) ·
  [Python and pip](#python-and-pip) · [Go modules](#go-modules) ·
  [Rust and Cargo](#rust-and-cargo)
- **Container images**: [Docker Hub](#docker-hub) ·
  [Other registries](#other-container-registries)
- **Coding agents**: [Claude Code](#claude-code) · [Codex](#codex)

Most sandboxes need the developer essentials regardless of language or agent: a
place to clone source from, certificate infrastructure for TLS, and OS package
mirrors. Layer a language registry and an agent recipe on top.

## Apply a recipe

Each recipe is a single network allow rule: a name and a list of host
resources. Apply them from the Admin Console or through the
[Governance API](/reference/api/ai-governance/). Policy changes take up to five
minutes to reach developer machines.

### Admin Console

Follow [Create a policy](org.md#create-a-policy) to open the **Network access**
policy form, then paste the hostnames from any recipe below as allow rules — one
per line. Group a recipe's hosts into a single rule and give it the recipe's
name so the rule list stays readable.

### API

Run the setup block once to get a token, create a policy, and define an
`add_rule` helper. Each recipe below is then a single call to that helper. The
base URL is `https://hub.docker.com/v2`, and the examples use
[`jq`](https://jqlang.github.io/jq/) to read IDs out of the responses.

```bash
# 1. Exchange a Docker Hub PAT or OAT for a short-lived bearer token.
TOKEN=$(curl -s -X POST https://hub.docker.com/v2/auth/token \
  -H "Content-Type: application/json" \
  -d '{"identifier":"my-user","secret":"dckr_pat_xxxxxxxx"}' | jq -r .access_token)

# 2. Create an empty policy for the org and capture its ID.
POLICY_ID=$(curl -s -X POST \
  https://hub.docker.com/v2/orgs/my-org/governance/policies \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"Sandbox network policy"}' | jq -r .id)

# 3. Helper that adds one network allow rule. Each recipe calls this.
add_rule() {
  curl -s -X POST \
    "https://hub.docker.com/v2/orgs/my-org/governance/policies/$POLICY_ID/rules" \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d "{\"name\":\"$1\",\"actions\":[\"connect:tcp\",\"connect:udp\"],\"resources\":$2,\"decision\":\"allow\"}"
}
```

> [!NOTE]
> The recipes use bare hostnames, with no `:port` suffix. A rule without a port
> suffix matches the host on any port, which covers both the HTTP and HTTPS
> ports the sandbox proxy handles — so plain-HTTP services such as OS package
> mirrors and certificate endpoints work without listing `:80` separately. Add a
> `:port` suffix (for example, `example.com:443`) only when you want to restrict
> a host to a specific port. For the full matching rules behind these patterns
> (exact hostnames, single- and multi-level wildcards, port suffixes, and CIDR
> ranges), see [Policy concepts](concepts.md#network-rules).

## Developer essentials

### GitHub

Clone and fetch source hosted on GitHub.

```bash
add_rule "GitHub" \
  '["github.com","**.github.com","**.githubusercontent.com"]'
```

If your developers host code elsewhere, swap in the equivalent hosts — for
example `gitlab.com` and `**.gitlab.com` for GitLab, or `bitbucket.org` for
Bitbucket.

### Certificate validation

Revocation and chain data that TLS clients fetch from certificate authorities,
much of it served over plain HTTP. Without these, some HTTPS handshakes stall or
fail. This is a trimmed set covering the most common authorities.

```bash
add_rule "Certificate validation" \
  '["**.lencr.org","ocsp.digicert.com","cacerts.digicert.com","**.pki.goog","**.amazontrust.com"]'
```

### Ubuntu packages

`apt` mirrors for the default Ubuntu base image. `ports.ubuntu.com` is the
mirror for non-x86 architectures, so ARM-based sandboxes need it.

```bash
add_rule "Ubuntu packages" \
  '["archive.ubuntu.com","security.ubuntu.com","ports.ubuntu.com","ubuntu.com"]'
```

If you build on a Debian base instead, add `debian.org` and `**.debian.org`.

## Package registries

Add only the registries for the languages your projects use.

### Node.js and npm

npm package installs and the Node.js runtime.

```bash
add_rule "npm registry" \
  '["registry.npmjs.org","npmjs.com","npmjs.org","nodejs.org","nodesource.com"]'
```

### Python and pip

pip installs and Python build backends.

```bash
add_rule "Python packages" \
  '["pypi.org","files.pythonhosted.org","pythonhosted.org","pypa.io","bootstrap.pypa.io"]'
```

### Go modules

`go` module downloads and checksum verification.

```bash
add_rule "Go modules" \
  '["proxy.golang.org","sum.golang.org","golang.org","pkg.go.dev"]'
```

### Rust and Cargo

Cargo crate downloads and `rustup` toolchain installs.

```bash
add_rule "Rust crates" \
  '["crates.io","index.crates.io","static.crates.io","static.rust-lang.org","sh.rustup.rs"]'
```

## Container images

Add only the registries your sandboxes pull from.

### Docker Hub

Pull images from Docker Hub.

```bash
add_rule "Docker Hub" \
  '["docker.io","**.docker.io","docker.com","**.docker.com","production.cloudflare.docker.com","**.production.cloudflare.docker.com"]'
```

### Other container registries

Common public registries: GitHub, Google, Microsoft, Quay, Kubernetes, and AWS
public ECR.

```bash
add_rule "Other registries" \
  '["ghcr.io","gcr.io","**.gcr.io","mcr.microsoft.com","**.data.mcr.microsoft.com","quay.io","registry.k8s.io","public.ecr.aws"]'
```

## Coding agents

Each coding agent talks to its own provider APIs. Add the recipe for the agents
your developers run, on top of the developer essentials and the relevant
language registries.

### Claude Code

Anthropic's Claude Code authenticates and streams completions through these
hosts.

```bash
add_rule "Anthropic APIs" \
  '["api.anthropic.com","statsig.anthropic.com","platform.claude.com","downloads.claude.ai","claude.com"]'
```

### Codex

OpenAI's Codex CLI authenticates and streams completions through OpenAI hosts.

```bash
add_rule "OpenAI APIs" \
  '["**.openai.com","chatgpt.com","**.chatgpt.com","**.oaistatic.com","**.oaiusercontent.com","cdn.openaimerge.com"]'
```

## Worked example: Claude Code on a Node.js project

A complete, minimal policy for developers running Claude Code against Node.js
projects hosted on GitHub combines the developer essentials, the npm registry,
and the Claude Code agent recipe. Run the [setup block](#api), then apply those
recipes:

```bash
add_rule "GitHub" \
  '["github.com","**.github.com","**.githubusercontent.com"]'
add_rule "Certificate validation" \
  '["**.lencr.org","ocsp.digicert.com","cacerts.digicert.com","**.pki.goog","**.amazontrust.com"]'
add_rule "Ubuntu packages" \
  '["archive.ubuntu.com","security.ubuntu.com","ports.ubuntu.com","ubuntu.com"]'
add_rule "npm registry" \
  '["registry.npmjs.org","npmjs.com","npmjs.org","nodejs.org","nodesource.com"]'
add_rule "Anthropic APIs" \
  '["api.anthropic.com","statsig.anthropic.com","platform.claude.com","downloads.claude.ai","claude.com"]'
```

Verify the result with `sbx policy ls` on a developer machine (the policy shows
a `Governance: managed by <org>` header), or fetch the full policy from the API:

```console
$ curl -s -H "Authorization: Bearer $TOKEN" \
    "https://hub.docker.com/v2/orgs/my-org/governance/policies/$POLICY_ID" | jq
```

## Related pages

- [Organization policy](org.md): how org governance works and where rules are
  configured
- [Policy concepts](concepts.md): resource model, rule syntax, and evaluation
- [AI Governance API](/reference/api/ai-governance/): full API reference
- [Monitoring](monitoring.md): inspect active rules and traffic with
  `sbx policy ls` and `sbx policy log`
