---
title: Continue local work in the cloud
linkTitle: Local-to-cloud handoff
description: Prepare a bounded agent task in a local sandbox, move its filesystem to the cloud, and retrieve the results.
keywords: docker agentic platform, cloud sandboxes, workflows, local-to-cloud handoff
weight: 30
---

Prepare a sample review locally, then let an agent finish it in the cloud while
your computer is disconnected. Retrieve the report before removing either
sandbox. Provisional duration: 25 minutes; this estimate needs measurement
during review.

## Before you begin

You need the [CLI and a supported local runtime](/manuals/ai/sandboxes/install.md),
version 0.45.0 or later, an active
[Docker Agentic Platform subscription](/manuals/agentic-platform/signup.md#activate-cloud-access),
and an Anthropic API key for Claude Code in the cloud. You can use a Claude
subscription locally, but cloud sandboxes require an API key. This guide uses
an API key in both environments.

Docker bills cloud compute while the destination runs, even when disconnected.
Anthropic bills inference separately. The public sample needs no GitHub token.
Use the managed credential stores so keys stay out of the filesystem snapshot:

```console
$ sbx login
$ sbx --cloud diagnose
$ sbx secret set anthropic
$ sbx --cloud secret set anthropic
```

Local and cloud secret stores are separate. See
[Cloud credentials](/manuals/ai/sandboxes/cloud/credentials.md).

## Prepare files inside the local sandbox {#1-prepare-files-inside-the-local-sandbox}

On your computer, create a mountless sandbox by omitting the workspace argument.
Clone the project inside it so the files transfer to the cloud. Directories
mounted from your computer and clone-mode volumes are outside the sandbox
filesystem and aren't included in a move:

```console
$ sbx create --name handoff-local --cpus 2 --memory 4g claude
$ sbx policy allow network github.com:443 --sandbox handoff-local
$ sbx exec handoff-local git clone https://github.com/docker/welcome-to-docker.git /home/agent/workspace/project
$ sbx run --name handoff-local
```

Ask the local agent to prepare the task without carrying it out:

```text
Read /home/agent/workspace/project. Don't modify it or install dependencies.
Record the full git HEAD revision and confirm the working tree is clean.
Write /home/agent/workspace/HANDOFF.md with that baseline, the project path,
and this remaining task:
Review package.json and the Dockerfile. Explain the build and startup commands,
identify the application's listening port, and recommend one small improvement
with a source path and reason. Don't modify the application or install anything.
Save the review to /home/agent/workspace/review.md. Include the baseline revision,
files read, and limitations. Don't claim to have run the application.
Don't start this review yet.
```

Exit the local agent when preparation is complete. Before snapshotting, remove
any credentials saved in filesystem files, including interactive sign-in files.
Managed secrets aren't copied, but ordinary credential files are.

## Move the filesystem {#2-move-the-filesystem}

From your computer, run:

```console
$ sbx move handoff-local --to cloud --ttl 2h --on-timeout stop
```

The move stops the local source for capture and creates a separate destination.
Use its printed `sbx_`-prefixed ID in place of `<CLOUD_ID>`. The move also creates
a template: the filesystem image used to create the destination, including
`HANDOFF.md`. It remains after you remove the sandbox.

Cloud sandboxes support both `linux/amd64` and `linux/arm64`. The destination
matches your local sandbox's architecture; a move doesn't convert between them.

The explicit timeout requests that the destination stop after two hours and
retain its files; the command
fails if stopping isn't supported. Don't substitute deletion without first
arranging to save your work.

A move transfers a filesystem snapshot. Running processes, memory, open sockets,
workspace mounts, clone-mode volumes, and managed secrets don't transfer.
Local network rules don't transfer either. Published TCP sandbox ports can be
republished with cloud URLs; rejected ports produce warnings. See
[Move a sandbox](/manuals/ai/sandboxes/cloud/move.md) for transfer details.

## Continue in the Console {#3-continue-in-the-console}

Open the [Console](https://agentic-platform.docker.com) and sign in with the
same Docker account you used for the CLI. Open **Sandboxes** and select the
destination using the ID printed by the move. Use its browser terminal to
continue the task; the local agent process wasn't transferred.

Before starting work, confirm the destination's architecture, compute size,
and expiration match your expectations. Check that your cloud Anthropic
credential is set. This task needs no published ports; remove any unintended
exposure. Verify the
[cloud network policy](/manuals/ai/sandboxes/cloud/network-policy.md) permits
the model provider; local network rules don't transfer.

In an additional shell, check the transferred handoff:

```console
$ cat /home/agent/workspace/HANDOFF.md
```

Expect the recorded baseline, project path, and remaining task. Return to the
agent's terminal to restart the work.

Give it this prompt:

```text
Read /home/agent/workspace/HANDOFF.md. Confirm the baseline revision and clean
working tree before doing the remaining review. Stop and report a mismatch.
Complete only the recorded task and save /home/agent/workspace/review.md.
```

## Disconnect and retrieve the result {#4-disconnect-and-retrieve-the-result}

Once the agent starts, close the browser tab to disconnect. The cloud agent
can continue while your computer is disconnected. The timer still applies;
no further work runs after the sandbox stops until you resume it.

Return to **Sandboxes** in the Console and reopen the destination. If it is
paused, select **Restore to connect**. Resuming starts a fresh timer. Review
the agent's response in the browser terminal.

> [!TIP]
> If the terminal appears stuck after reconnecting, try resizing it by resizing
> your browser window. This can make the agent interface redraw.

When the report is ready, return to your computer's terminal and copy both files:

```console
$ sbx --cloud ttl <CLOUD_ID>
$ sbx --cloud cp <CLOUD_ID>:/home/agent/workspace/review.md ./review.md
$ sbx --cloud cp <CLOUD_ID>:/home/agent/workspace/HANDOFF.md ./HANDOFF.md
```

Read the local copies. Expect the report to identify the recorded baseline,
explain the build and startup commands with file references, and state that the
application wasn't run. Compare it with the source rather than treating the
agent's claims as verified. Files don't synchronize after a move: copying these
results is what brings them back to your computer.

## Clean up both sandboxes {#5-clean-up-both-sandboxes}

After verifying your saved files, remove the destination, its template, and the
local source. First inspect the destination to identify its backing template,
then list the templates. Find the matching entry, typically named
`moved-handoff-local-<SUFFIX>`, and use its ID as `<TEMPLATE_ID>`. Check the ID
against the inspected sandbox so you remove only this move's snapshot:

```console
$ sbx --cloud inspect <CLOUD_ID>
$ sbx --cloud template ls
```

Delete the cloud sandbox before removing its backing template:

```console
$ sbx --cloud rm <CLOUD_ID>
$ sbx --cloud template rm <TEMPLATE_ID>
$ sbx rm handoff-local
$ sbx --cloud ls
$ sbx --cloud template ls
$ sbx ls
```

Confirm both sandboxes and the template are absent. The template contains the
captured filesystem, so removing only the sandbox leaves that snapshot stored.
Removal deletes files stored only in those resources.
If retaining the cloud destination instead, stop it and confirm its state with
`sbx --cloud ls` to stop compute charges. Check the expiration again after any
resume. See [Stop or remove a sandbox](/manuals/ai/sandboxes/cloud/usage.md#stop-or-remove-a-sandbox).
