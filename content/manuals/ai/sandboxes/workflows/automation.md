---
title: Run sandboxes in CI
linkTitle: CI and headless
weight: 50
description: Authenticate and run Docker Sandboxes in CI systems and other headless environments.
keywords: docker sandboxes, sbx, ci, headless, automation, personal access token
---

For CI environments and scripts where a browser isn't available, authenticate
with a Docker Personal Access Token (PAT):

```console
$ echo "$DOCKER_PAT" | sbx login --username <your-docker-id> --password-stdin
```

Generate a PAT from your
[Docker account settings](https://app.docker.com/settings/personal-access-tokens)
with at least **Read** scope.

From there, the rest of the `sbx` workflow is the same as interactive use.
Create the sandbox in the background with `sbx create`, run agent tasks with
`sbx exec`, and clean up with `sbx rm`:

```console
$ sbx create --name ci-task --clone claude
$ sbx run --name ci-task  # attach and give instructions, or use sbx exec for one-off commands
$ git fetch sandbox-ci-task
$ sbx rm ci-task
```

Agent credentials (API keys, GitHub token) can be preconfigured as global
secrets so they're available to any sandbox the CI runner creates. If the
relevant environment variables are already set in the CI environment (see the
[built-in services table](../configuration/credentials.md#built-in-services) for which
variables each service reads), import them all at once:

```console
$ sbx secret import --all
```

To overwrite an existing stored entry, add `--force`. To pass a value from your
CI provider's secret store, use `-t`. For example, in a GitHub Actions step:

```yaml
- run: sbx secret set anthropic -t "${{ secrets.ANTHROPIC_API_KEY }}"
```
