---
title: Debug a kit
linkTitle: Debug a kit
description: Read the setup progress output, find install and startup logs inside the sandbox, and fix a kit whose install or startup command fails.
keywords: sandboxes, sbx, kits, debug, troubleshooting, install, startup, logs
weight: 36
---

{{< summary-bar feature_name="Docker Sandboxes sbx" >}}

When a kit's install or startup command fails, `sbx` shows you which command
failed, in which kit, and where the full output is. This page shows how to read
that output, find the logs inside the sandbox, and fix the common causes.

## Read the progress output

During `sbx create` and `sbx run`, the configure agent phase prints one row for
each kit operation: install commands, init files, file copies, environment
variables, and the startup commands it registers. A finished row shows ✓ with
the elapsed time, or ✗ with the exit code.

Each command row shows the command's `description` when the kit sets one.
Otherwise it shows the first 60 characters of the command. It also shows the
kit that owns the command (`kit=<name>`) and the user it runs as (`user=<uid>`).

Set `description` on every install and startup command so the row says what the
command is for. Keep it short: a row shows at most 60 characters of it.

```yaml {title="my-kit/my-kit.yaml"}
capabilities:
  - type: com.docker.sandbox/lifecycle@1
    config:
      install:
        - command: apt-get update && apt-get install -y ripgrep
          user: "0"
          description: Install ripgrep
      startup:
        - command: rg --version > /tmp/rg-version
          description: Record the ripgrep version
```

The rows for this kit look like this:

```text
AGENT
    running 1 install command
      ✓ Install ripgrep (kit=my-kit, user=0, 12.4s)
      view install logs: ls /var/log/sbx-kits
    registering 1 startup command, run on every container start
      ○ Record the ripgrep version (kit=my-kit, user=agent)
      view startup log: cat /var/log/sbx-kit-startup.log
```

## Read a failed install command

A failed install command prints the last 40 lines of its output under the row,
followed by a `docker run` line that runs the same command against the kit's
base image:

```text
  ✗ Install ripgrep (kit=my-kit, user=0, exit 100)
    ── captured output ──
    E: Failed to fetch http://deb.debian.org/debian/dists/bookworm/InRelease  403  Forbidden
    E: Unable to locate package ripgrep
    ──────────────────────
  hint: reproduce locally with:
    docker run --rm -u 0 <base-image> sh -c 'apt-get update && apt-get install -y ripgrep'
```

Copy the `docker run` line to iterate on the command without recreating the
sandbox each time. Change the command, run it again, and copy the working
version back into the kit. The reproduction runs in plain Docker, outside the sandbox, so the sandbox's
network policy doesn't apply. If the command works there but fails in the
sandbox, check the kit's network rules first.

## Find the logs inside the sandbox

The progress output shows only the end of a failed command. The sandbox keeps
the full record:

| Log | Contents |
| --- | --- |
| `/var/log/sbx-kits/<kit>/install-<n>.log` | Full output of one install command. `<n>` is the command's position in that kit's `install` list, starting at 0. |
| `/var/log/sbx-kit-startup.log` | Output of every startup command, from every container start. |

Read either log with `sbx exec`:

```console
$ sbx exec my-sandbox cat /var/log/sbx-kits/my-kit/install-0.log
$ sbx exec my-sandbox cat /var/log/sbx-kit-startup.log
```

Startup commands run on every container start. Each one writes these lines to
the startup log:

```text
> <script>
# <kit>: <description or command>
<the command's output>
ok <script>
```

A failing command ends with `fail <script> exit=<code>` instead of `ok`. It
also stops the startup commands after it, and `sbx create` prints a warning
that points at the startup log with its last lines.

After the create phase, the agent takes over the terminal. The logs in the
sandbox are the record you can come back to.

## Show everything with -D

Pass the global debug flag to `sbx create` or `sbx run`:

```console
$ sbx -D create --name my-sandbox shell --kit ./my-kit
```

With `-D`, each successful install command also streams its full output under
its row, and `sbx` prints the startup log after the sandbox starts.

## Common causes

### The install command can't reach its package source

The log shows names that don't resolve, or HTTP 403 responses from the
sandbox proxy, from `apt`, `pip`, or `npm`. The kit's install-phase network
rules don't include the package mirror. Add it to `install.allow` in the
`com.docker.sandbox/network-policy@1` capability. The `apt` mirrors, for
example, can live on a different host than the one the command names.

### Two kits run `apt` at the same time

The log shows `Could not get lock /var/lib/dpkg/lock-frontend`. Another
install command, or the base image setup, holds the dpkg lock. Retry the
package command in the kit, for example with `-o DPkg::Lock::Timeout=120`
on `apt-get`.

### A startup command needs a tool that didn't install

The startup log shows `exec: <tool>: not found` after an install command for
that tool failed. Fix the install command first. Check its
`install-<n>.log`, then rerun the sandbox.

### The command needs root

The log shows `Permission denied` when the command writes to a system path or
installs a package, and the row shows `user=1000`. Set `user: "0"` on the
install command, or write to a location the agent user owns.

## Next steps

- [Kit authoring patterns](/manuals/ai/sandboxes/customize/author/patterns.md)
  shows setup commands in context.
- [Sandbox troubleshooting](/manuals/ai/sandboxes/troubleshooting.md) covers
  other sandbox problems.
