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
and an Anthropic API key for Claude Code. Your cloud account must support the
local sandbox's architecture: a move doesn't convert Arm to Intel/AMD or the reverse.

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

## 1. Prepare files inside the local sandbox

On your computer, create a sandbox without a workspace argument:

```console
$ sbx create --name handoff-local --cpus 2 --memory 4g claude
$ sbx policy allow network github.com:443 --sandbox handoff-local
$ sbx exec handoff-local git clone https://github.com/docker/welcome-to-docker.git /home/agent/workspace/project
$ sbx run --name handoff-local
```

The project belongs to the sandbox filesystem. Host workspace mounts and
clone-mode volumes aren't included in a move.

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

Read the handoff file from another terminal on your computer:

```console
$ sbx exec handoff-local cat /home/agent/workspace/HANDOFF.md
$ sbx exec handoff-local git -C /home/agent/workspace/project rev-parse HEAD
$ sbx exec handoff-local git -C /home/agent/workspace/project status --short
```

Expect the recorded revision to match and the project working tree to be clean.
Exit the local agent when preparation is complete. Before snapshotting, remove
any credentials saved in filesystem files, including interactive sign-in files.
Managed secrets aren't copied, but ordinary credential files are.

## 2. Move the filesystem

From your computer, run:

```console
$ sbx move handoff-local --to cloud --ttl 2h --on-timeout stop
```

The move stops the local source for capture and creates a separate destination.
Use its printed name or ID in place of `<CLOUD_ID>`. The explicit timeout requests
that the destination stop after two hours and retain its files; the command
fails if stopping isn't supported. Don't substitute deletion without first
arranging to save your work.

A move transfers a filesystem snapshot. Running processes, memory, open sockets,
workspace mounts, clone-mode volumes, and managed secrets don't transfer.
Local network rules don't transfer either. Published TCP sandbox ports can be
republished with cloud URLs; rejected ports produce warnings. See
[Move a sandbox](/manuals/ai/sandboxes/cloud/move.md) for transfer details.

## 3. Verify the destination and restart the agent

Inspect both states and the destination before doing more work:

```console
$ sbx ls
$ sbx --cloud ls
$ sbx --cloud inspect <CLOUD_ID>
$ sbx --cloud ttl <CLOUD_ID>
$ sbx --cloud secret ls
$ sbx --cloud ports <CLOUD_ID>
$ sbx --cloud exec <CLOUD_ID> cat /home/agent/workspace/HANDOFF.md
$ sbx --cloud exec <CLOUD_ID> git -C /home/agent/workspace/project rev-parse HEAD
$ sbx --cloud exec <CLOUD_ID> git -C /home/agent/workspace/project status --short
```

Expect the source to be stopped, a distinct destination ID, the same baseline,
a clean project tree, and the handoff file. Confirm the platform, compute size,
expiration action, and cloud Anthropic credential. This task needs no published
ports; remove any unintended exposure before continuing. Verify destination
[network policy](/manuals/ai/sandboxes/cloud/network-policy.md) permits the model
provider. The cloud uses its own policy and agent rules, rather than local rules.

Restart the agent explicitly; its local process wasn't transferred:

```console
$ sbx --cloud attach <CLOUD_ID>
```

Give it this prompt:

```text
Read /home/agent/workspace/HANDOFF.md. Confirm the baseline revision and clean
working tree before doing the remaining review. Stop and report a mismatch.
Complete only the recorded task and save /home/agent/workspace/review.md.
```

## 4. Disconnect and retrieve the result

Once the agent starts, press `Ctrl+\` to detach. The cloud agent can continue
while your computer is disconnected. The timer still applies; no further work
runs after the sandbox stops until you resume it.

Return from your computer's terminal:

```console
$ sbx --cloud attach <CLOUD_ID>
```

Review the agent's response. Attaching resumes a stopped sandbox and starts a
fresh timer. When the report is ready, detach with `Ctrl+\` and copy both files:

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

## 5. Clean up both sandboxes

After verifying your saved files, remove the destination and source separately:

```console
$ sbx --cloud rm <CLOUD_ID>
$ sbx rm handoff-local
$ sbx --cloud ls
$ sbx ls
```

Confirm both are absent. Removal deletes files stored only in those sandboxes.
If retaining the cloud destination instead, stop it and confirm its state with
`sbx --cloud ls` to stop compute charges. Check the expiration again after any
resume. See [Stop or remove a sandbox](/manuals/ai/sandboxes/cloud/usage.md#stop-or-remove-a-sandbox).
