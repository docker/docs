---
title: Troubleshooting
weight: 60
description: Resolve common issues when using Docker Sandboxes.
---

{{< summary-bar feature_name="Docker Sandboxes sbx" >}}

## Resetting sandboxes

If you hit persistent issues or corrupted state, run
[`sbx reset`](/reference/cli/sbx/reset/) to stop all VMs and delete all sandbox
data. Create fresh sandboxes afterwards.

## Agent can't install packages or reach an API

Sandboxes use a [deny-by-default network policy](security/policy.md).
If the agent fails to install packages or call an external API, the target
domain is likely not in the allow list. Check which requests are being blocked:

```console
$ sbx policy log
```

Then allow the domains your workflow needs:

```console
$ sbx policy allow network "*.npmjs.org,*.pypi.org,files.pythonhosted.org"
```

To allow all outbound traffic instead:

```console
$ sbx policy allow network "**"
```

## Docker authentication failure

If you see a message like `You are not authenticated to Docker`, your login
session has expired. In an interactive terminal, the CLI prompts you to sign in
again. In non-interactive environments such as scripts or CI, run `sbx login`
to re-authenticate.

## Agent authentication failure

If the agent can't reach its model provider or you see API key errors, the key
is likely invalid, expired, or not configured. Verify it's set in your shell
configuration file and that you sourced it or opened a new terminal.

For agents that use the [credential proxy](security/credentials.md), make sure
you haven't set the API key to an invalid value inside the sandbox — the proxy
injects credentials automatically on outbound requests.

If credentials are configured correctly but API calls still fail, check
`sbx policy log` and look at the **PROXY** column. Requests routed through
the `transparent` proxy don't get credential injection. This can happen when a
client inside the sandbox (such as a process in a Docker container) isn't
configured to use the forward proxy. See
[Monitoring network activity](security/policy.md#monitoring)
for details.

## Docker not available inside the sandbox on Windows

On Windows, sandboxes use non-docker template variants by default, so `docker`
commands aren't available inside the sandbox. To use Docker inside a sandbox on
Windows, specify a `-docker` template:

```console
$ sbx run --template docker.io/docker/sandbox-templates:claude-code-docker claude
```

The `-docker` variants work on Windows but have slower startup times. See
[Base images](agents/custom-environments.md#base-images) for details and the
full list of templates.

## Docker build export fails with "lchown: operation not permitted"

Running `docker build` with the local exporter (`--output=type=local` or `-o
<path>`) inside a sandbox fails because the exporter tries to `lchown` output
files to preserve ownership from the build. Processes inside the sandbox run as
an unprivileged user without `CAP_CHOWN`, so the operation is denied.

Use the tar exporter and extract the archive instead:

```console
$ mkdir -p ./result
$ docker build --output type=tar,dest=- . | tar xf - -C ./result
```

Extracting the tar archive as the current user avoids the `chown` call.

## Stale Git worktree after removing a sandbox

If you used `--branch`, worktree cleanup during `sbx rm` is best-effort. If
it fails, the sandbox is removed but the branch and worktree are left behind.
If `git worktree list` shows a stale worktree in `.sbx/` after removing a
sandbox, clean it up manually:

```console
$ git worktree remove .sbx/<sandbox-name>-worktrees/<branch-name>
$ git branch -D <branch-name>
```

## Clock drift after sleep/wake

If your laptop sleeps and wakes while a sandbox is running, the VM clock can
fall behind the host clock. This causes problems such as:

- External API calls failing because of timestamp validation.
- Git commits with incorrect timestamps.
- TLS certificate errors due to time mismatches.

To fix the issue, stop and restart the sandbox:

```console
$ sbx stop <sandbox-name>
$ sbx run <sandbox-name>
```

Restarting the sandbox re-syncs the VM clock with the host.

## Removing all state

As a last resort, if `sbx reset` doesn't resolve your issue, you can remove the
`sbx` state directory entirely. This deletes all sandbox data, configuration, and
cached images. Stop all running sandboxes first with `sbx reset`.

macOS:

```console
$ rm -rf ~/Library/Application\ Support/com.docker.sandboxes/
```

Windows:

```powershell
> Remove-Item -Recurse -Force "$env:LOCALAPPDATA\DockerSandboxes"
```
