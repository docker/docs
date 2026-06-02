---
title: Org policy recipes
linkTitle: Policy recipes
weight: 22
description: Minimal, composable network policy presets for common sandbox workflows, in Admin Console and API form.
keywords: docker sandboxes, governance, organization policy, network policy, allowlist, recipes, presets, Claude Code, Codex, package managers
---

Organization policies start deny-by-default: outbound traffic is blocked unless
an explicit allow rule matches, so every host or wildcard a workflow depends on
has to be added by hand. The recipes on this page are minimal starting points
you can copy, combine, and trim to fit your organization.

These recipes are intentionally narrower than the local
[Balanced preset](local.md#default-preset). The Balanced preset is a catch-all
profile that allows a wide range of common developer destinations so individual
machines work out of the box. An organization policy should do the opposite:
allow only the destinations your organization actually wants its sandboxes to
reach. Treat each recipe as a baseline and remove anything you don't need.

To work out exactly which hosts a workflow needs, run it in a sandbox and watch
the requests it makes with `sbx policy log`. Blocked destinations show up there,
so you can add the ones you want and leave the rest denied. See
[Monitoring](monitoring.md).

For the exact matching rules behind the resource patterns used here (exact
hostnames, single- and multi-level wildcards, port suffixes, and CIDR ranges),
see
[Policy concepts](concepts.md#network-rules).

## How to apply a recipe

Each recipe is a set of network allow rules. Apply them either from the Admin
Console or through the [Governance API](/reference/api/ai-governance/).

### Admin Console

In the [Admin Console](https://app.docker.com/admin), go to
**AI governance > Network access** and add the listed hostnames as allow rules.
You can paste multiple hostnames at once, one per line. Group related hosts into
a single rule with a descriptive name so the rule list stays readable.

### API

With the API, create a policy once, then add each block as a rule. The base URL
is `https://hub.docker.com/v2` and requests use a short-lived JWT obtained by
exchanging Docker Hub credentials. The examples use
[`jq`](https://jqlang.github.io/jq/) to extract IDs from the responses.

```bash
# 1. Exchange a Docker Hub PAT or OAT for a bearer token.
TOKEN=$(curl -s -X POST https://hub.docker.com/v2/users/login \
  -H "Content-Type: application/json" \
  -d '{"username":"my-user","password":"dckr_pat_xxxxxxxx"}' | jq -r .token)

# 2. Create an empty policy for the org and capture its ID.
POLICY_ID=$(curl -s -X POST \
  https://hub.docker.com/v2/orgs/my-org/governance/policies \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"Claude Code — Node.js"}' | jq -r .id)
```

Each rule is then a single `POST` to the policy's `rules` sub-resource. Network
allow rules use the actions `connect:tcp` and `connect:udp` with
`"decision": "allow"`:

```bash
curl -s -X POST \
  "https://hub.docker.com/v2/orgs/my-org/governance/policies/$POLICY_ID/rules" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Anthropic APIs",
    "actions": ["connect:tcp", "connect:udp"],
    "resources": ["api.anthropic.com", "statsig.anthropic.com",
                  "platform.claude.com", "downloads.claude.ai", "claude.com"],
    "decision": "allow"
  }'
```

The rest of the recipes list only the rule name and resources. Substitute them
into the `name` and `resources` fields of a request like the one above. Policy
changes take up to five minutes to reach developer machines.

> [!NOTE]
> The only egress path out of a sandbox is the host HTTP/HTTPS proxy, so the
> resources below are all hostnames on ports 80 and 443. Where a host is
> reached over plain HTTP (OS package mirrors and certificate revocation
> endpoints), the port `:80` variant is listed explicitly; HTTPS-only hosts are
> listed without a port suffix.

## Developer essentials

Most sandboxes need a small baseline regardless of which agent or language they
run: a place to clone source from, the certificate infrastructure that TLS
handshakes validate against, and the OS package mirrors the base image installs
from. Start with this block and layer a language and agent recipe on top.

### Source and version control

| Rule name | Resources                                                 |
| --------- | --------------------------------------------------------- |
| GitHub    | `github.com`, `**.github.com`, `**.githubusercontent.com` |

If your developers host code elsewhere, swap GitHub for the equivalent hosts —
for example `gitlab.com` and `**.gitlab.com` for GitLab, or `bitbucket.org` for
Bitbucket.

### Certificate validation

TLS clients fetch revocation and chain data from certificate authority
endpoints, many of which are served over plain HTTP on port 80. Without these,
some HTTPS handshakes stall or fail. This is a trimmed set covering the most
common authorities:

| Rule name              | Resources                                                                                                                                                            |
| ---------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Certificate validation | `**.lencr.org`, `**.lencr.org:80`, `ocsp.digicert.com:80`, `cacerts.digicert.com:80`, `**.pki.goog`, `**.pki.goog:80`, `**.amazontrust.com`, `**.amazontrust.com:80` |

### Operating system packages

The default sandbox base image is Ubuntu, so `apt` reaches the Ubuntu mirrors
over both HTTP and HTTPS:

| Rule name       | Resources                                                                                                                                               |
| --------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Ubuntu packages | `archive.ubuntu.com`, `archive.ubuntu.com:80`, `security.ubuntu.com`, `security.ubuntu.com:80`, `ports.ubuntu.com`, `ports.ubuntu.com:80`, `ubuntu.com` |

If you build on a Debian base instead, add `debian.org` and `**.debian.org`.

## Package registries

Add only the registries for the languages your projects use.

### Node.js and npm

| Rule name    | Resources                                                                      |
| ------------ | ------------------------------------------------------------------------------ |
| npm registry | `registry.npmjs.org`, `npmjs.com`, `npmjs.org`, `nodejs.org`, `nodesource.com` |

### Python and pip

| Rule name       | Resources                                                                                |
| --------------- | ---------------------------------------------------------------------------------------- |
| Python packages | `pypi.org`, `files.pythonhosted.org`, `pythonhosted.org`, `pypa.io`, `bootstrap.pypa.io` |

### Go modules

| Rule name  | Resources                                                        |
| ---------- | ---------------------------------------------------------------- |
| Go modules | `proxy.golang.org`, `sum.golang.org`, `golang.org`, `pkg.go.dev` |

### Rust and Cargo

| Rule name   | Resources                                                                                  |
| ----------- | ------------------------------------------------------------------------------------------ |
| Rust crates | `crates.io`, `index.crates.io`, `static.crates.io`, `static.rust-lang.org`, `sh.rustup.rs` |

## Container images

If sandboxes pull container images, allow the registries you use. This block
covers Docker Hub and the most common public registries:

| Rule name        | Resources                                                                                                                             |
| ---------------- | ------------------------------------------------------------------------------------------------------------------------------------- |
| Docker Hub       | `docker.io`, `**.docker.io`, `docker.com`, `**.docker.com`, `production.cloudflare.docker.com`, `**.production.cloudflare.docker.com` |
| Other registries | `ghcr.io`, `gcr.io`, `**.gcr.io`, `mcr.microsoft.com`, `**.data.mcr.microsoft.com`, `quay.io`, `registry.k8s.io`, `public.ecr.aws`    |

## Agent rules

Each coding agent talks to its own provider APIs. Add the block for the agents
your developers run, on top of the developer essentials and the relevant
language registries.

### Claude Code

| Rule name      | Resources                                                                                                |
| -------------- | -------------------------------------------------------------------------------------------------------- |
| Anthropic APIs | `api.anthropic.com`, `statsig.anthropic.com`, `platform.claude.com`, `downloads.claude.ai`, `claude.com` |

### Codex

OpenAI's Codex CLI authenticates and streams completions through OpenAI hosts:

| Rule name   | Resources                                                                                                            |
| ----------- | -------------------------------------------------------------------------------------------------------------------- |
| OpenAI APIs | `**.openai.com`, `chatgpt.com`, `**.chatgpt.com`, `**.oaistatic.com`, `**.oaiusercontent.com`, `cdn.openaimerge.com` |

## Worked example: Claude Code on a Node.js project

The following script builds a complete, minimal policy for developers running
Claude Code against Node.js projects hosted on GitHub. It combines the developer
essentials, the npm registry, and the Claude Code agent block.

```bash
# Authenticate and create the policy.
TOKEN=$(curl -s -X POST https://hub.docker.com/v2/users/login \
  -H "Content-Type: application/json" \
  -d '{"username":"my-user","password":"dckr_pat_xxxxxxxx"}' | jq -r .token)

POLICY_ID=$(curl -s -X POST \
  https://hub.docker.com/v2/orgs/my-org/governance/policies \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"Claude Code — Node.js"}' | jq -r .id)

# Add each block as a named rule.
add_rule() {
  curl -s -X POST \
    "https://hub.docker.com/v2/orgs/my-org/governance/policies/$POLICY_ID/rules" \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d "{\"name\":\"$1\",\"actions\":[\"connect:tcp\",\"connect:udp\"],\"resources\":$2,\"decision\":\"allow\"}"
}

add_rule "GitHub" \
  '["github.com","**.github.com","**.githubusercontent.com"]'
add_rule "Certificate validation" \
  '["**.lencr.org","**.lencr.org:80","ocsp.digicert.com:80","cacerts.digicert.com:80","**.pki.goog","**.pki.goog:80","**.amazontrust.com","**.amazontrust.com:80"]'
add_rule "Ubuntu packages" \
  '["archive.ubuntu.com","archive.ubuntu.com:80","security.ubuntu.com","security.ubuntu.com:80","ports.ubuntu.com","ports.ubuntu.com:80","ubuntu.com"]'
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
  </content>
  </invoke>
