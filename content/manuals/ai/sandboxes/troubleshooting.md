---
title: Troubleshooting
weight: 60
description: Resolve common issues when using Docker Sandboxes.
keywords: docker sandboxes, sbx, troubleshooting, diagnostics, reset, network policy, git, ssh
---

{{< summary-bar feature_name="Docker Sandboxes sbx" >}}

## Run diagnostics

Before digging into a specific issue, run
[`sbx diagnose`](/reference/cli/sbx/diagnose/) to check for common problems
with your installation, such as a missing CLI binary, an unresponsive daemon,
a CLI/daemon version mismatch, missing storage directories, or broken
authentication.

```console
$ sbx diagnose
```

The command prints a summary of checks that passed, warned, or failed, along
with suggested fixes. Use `--output json` to get machine-readable output, or
`--output github-issue` to generate a Markdown snippet suitable for pasting
into a GitHub issue.

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
$ sbx policy allow network -g "*.npmjs.org,*.pypi.org,files.pythonhosted.org"
```

To allow all outbound traffic instead:

```console
$ sbx policy allow network -g "**"
```

If `sbx policy allow` doesn't unblock the request, your organization may
manage sandbox policies centrally and take precedence over local rules. See
[Organization governance](security/governance.md).

## SSH and other non-HTTP connections fail

Non-HTTP TCP connections like SSH can be allowed by adding a policy rule for
the destination IP address and port. For example, to allow SSH to a specific
host:

```console
$ sbx policy allow network -g "10.1.2.3:22"
```

Hostname-based rules (for example, `myhost:22`) don't work for non-HTTP
connections because the proxy can't resolve the hostname to an IP address in
this context. Use the IP address directly.

UDP and ICMP traffic is blocked at the network layer and can't be unblocked
with policy rules.

For Git operations over SSH, you can either add an allow rule for the Git
server's IP address or use HTTPS URLs instead:

```console
$ git clone https://github.com/owner/repo.git
```

## Can't reach a service running on the host

If a request to `127.0.0.1` or a local network IP returns "connection refused"
from inside a sandbox, the address is not reachable from within the sandbox VM.
See [Accessing host services from a sandbox](usage.md#accessing-host-services-from-a-sandbox).

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

## Docker build export fails with an ownership error

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

## Sandbox commits aren't signed

Docker Sandboxes can sign Git commits with SSH keys from your host agent.
For setup steps, see [Signed commits](usage.md#signed-commits).

If `ssh-add -L` prints `The agent has no identities.`, the sandbox can reach
the forwarded agent, but the host agent doesn't have a loaded key. Load the
signing key into your host SSH agent:

```console
$ ssh-add ~/.ssh/id_ed25519
```

If commit signing works on the host but fails in a sandbox, check whether Git
is configured to sign with a host file path such as
`/Users/me/.ssh/id_ed25519.pub`. The sandbox uses the forwarded SSH agent, not
the host key file path. Use the inline public key form instead:

```console
$ git config --global gpg.format ssh
$ git config --global user.signingkey "key::$(ssh-add -L | head -n 1)"
```

If Git reports that `ssh-keygen` is missing, use a sandbox template that
includes OpenSSH client tools.

If `git log --show-signature` reports that `gpg.ssh.allowedSignersFile` needs
to be configured, Git can't verify the SSH signature locally. This verification
config isn't required to create signed commits. GitHub uses the SSH signing
keys configured in your GitHub account to verify commits.

GPG and S/MIME signing keys aren't available inside the sandbox. If your
repository or organization requires GPG or S/MIME signatures, or if SSH signing
isn't configured, use one of these workarounds:

- Commit outside the sandbox. Let the agent make changes without committing,
  then commit and sign from your host terminal.

- Sign after the fact. Let the agent commit inside the sandbox, then re-sign
  the commits on your host:

  ```console
  $ git rebase --exec 'git commit --amend --no-edit -S' origin/main
  ```

  This replays each commit on the branch and re-signs it with your local
  signing key.

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

## Daemon fails to start after downgrading

If you downgrade `sbx` to a version older than the one that last managed your
local state, the daemon may fail to start with a database version mismatch:

```text
ERROR: failed to start backend in-process: start backend: creating containerd
server: ... database is at major version 6, but this binary only supports up
to major version 1
```

A newer version of `sbx` upgraded the local database to a schema that older
binaries don't understand. To recover, reset all sandbox state:

```console
$ sbx reset --preserve-secrets
```

This stops all VMs and deletes all sandbox data. You'll need to create new
sandboxes afterwards. The `--preserve-secrets` flag keeps any secrets you've
set so you don't have to reconfigure them.

## Removing all state

As a last resort, if `sbx reset` doesn't resolve your issue, you can remove the
`sbx` state directory entirely. This deletes all sandbox data, configuration, and
cached images. Stop all running sandboxes first with `sbx reset`.

{{< tabs >}}
{{< tab name="macOS" >}}

```console
$ rm -rf ~/Library/Application\ Support/com.docker.sandboxes/
```

{{< /tab >}}
{{< tab name="Windows" >}}

```powershell
> Remove-Item -Recurse -Force "$env:LOCALAPPDATA\DockerSandboxes"
```

{{< /tab >}}
{{< tab name="Linux" >}}

Sandbox state on Linux follows the XDG Base Directory specification and is
spread across three directories:

```console
$ rm -rf ~/.local/state/sandboxes/
$ rm -rf ~/.cache/sandboxes/
$ rm -rf ~/.config/sandboxes/
```

If you have set custom `XDG_STATE_HOME`, `XDG_CACHE_HOME`, or
`XDG_CONFIG_HOME` environment variables, replace `~/.local/state`,
`~/.cache`, and `~/.config` with the corresponding values.

{{< /tab >}}
{{< /tabs >}}

## Report an issue

If you've exhausted the steps above and the problem persists, file a GitHub
issue at [github.com/docker/sbx-releases/issues](https://github.com/docker/sbx-releases/issues).

To help Docker investigate, generate a diagnostics bundle and share it when
reporting the issue:

```console
$ sbx diagnose --upload
```

The bundle contains daemon logs, diagnostic check results, and basic system
information. When `--upload` is confirmed, the bundle is uploaded to Docker
support and the command prints a diagnostics ID. Include this ID in your
issue so the team can correlate it with the uploaded bundle.
