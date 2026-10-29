---
title: Send local work to the cloud
description: Prepare a project in a local Docker Sandbox, move its files to the cloud, let Claude Code continue while you disconnect, and copy the finished work back.
keywords: docker sandboxes, sbx move, cloud sandbox, claude code, task handoff
summary: |
  Carry a prepared project and task notes into a cloud sandbox, leave the
  agent working, and retrieve a bundle of its results.
params:
  tags: [ai]
  time: 25 minutes
---

<!-- Publication gate: publish the companion URL and baseline revision. Run the
local-to-cloud journey with recorded CLI and agent versions and architecture.
Verify .venv and HANDOFF.md, disconnect during work, retrieve the bundle, and
capture the destination and retrieval. No sandbox transfer was run locally. -->

Prepare an engineering task locally, then let an agent continue in the cloud
while your laptop is disconnected. This guide moves a sandbox containing an
API-change reporter, its installed dependencies, and a task handoff file. You
return later to retrieve the updated source, patch, report, and test output.

The task is to extend the reporter to detect required query parameters added to
existing operations. You can start from the companion project without running
the [build guide](cloud-sandbox-api-reporter.md).

## Before you start

You need the [`sbx` CLI and a supported local sandbox runtime](/manuals/ai/sandboxes/install.md),
an Anthropic API key, and
[cloud sandbox access](/manuals/ai/sandboxes/cloud/_index.md#prerequisites).
Your cloud account must support your local sandbox's architecture. A move
preserves the platform; it doesn't convert between Arm and Intel/AMD.

The commands use a POSIX shell. Replace `REPOSITORY_URL` with the companion
repository's clone URL and `BASELINE_SHA` with its full commit SHA.

Sign in and check cloud access:

```console
$ sbx login
$ sbx --cloud diagnose
```

Configure the Anthropic key in both CLI stores:

```console
$ sbx secret set anthropic
$ sbx --cloud secret set anthropic
```

Local CLI, cloud CLI, and Console credentials are separate. Managed secrets
don't travel with a move. Using these stores also keeps credentials out of the
filesystem snapshot.

## Prepare inside the local sandbox

Create a sandbox with explicit resource limits and no workspace argument:

```console
$ sbx create --name reporter-local --cpus 2 --memory 4g claude
$ sbx policy allow network github.com:443 --sandbox reporter-local
$ sbx policy allow network pypi.org:443 --sandbox reporter-local
$ sbx policy allow network files.pythonhosted.org:443 --sandbox reporter-local
$ sbx exec reporter-local git clone REPOSITORY_URL /home/agent/workspace/project
$ sbx exec reporter-local git -C /home/agent/workspace/project checkout BASELINE_SHA
$ sbx run --name reporter-local
```

Omitting the workspace argument matters. The project lives inside the sandbox
filesystem. Ordinary host mounts and clone-mode volumes are excluded from
moves, so mounting your laptop's project would leave those files behind.

Ask the local agent to prepare the handoff:

```text
Prepare /home/agent/workspace/project/api-change-reporter for a cloud task.
Create .venv there, install requirements.txt, and run all tests. Record the
Python and dependency versions and the passing baseline. Don't implement the
extension yet. Stop if preparation fails.

Write /home/agent/workspace/project/HANDOFF.md with the baseline commit,
objective, project path, relevant files, completed preparation, and remaining
work. Record this exact test command, run from api-change-reporter/:
/home/agent/workspace/project/api-change-reporter/.venv/bin/python -m unittest discover -s tests -v

The remaining task: detect query parameters that become required on existing
operations, including optional-to-required changes. Cover inline path-level
and operation-level parameters and operation overrides. Leave reference
resolution outside the scope. Add regression tests and report the findings.
```

Read `HANDOFF.md` and the baseline results. The companion implementation has
eight tests. Exit the local agent after preparation completes, then return to
your host terminal.

## Move the prepared files

```console
$ sbx move reporter-local --to cloud --ttl 2h --on-timeout stop
```

Use the destination ID printed by the command in place of `CLOUD_ID` below.
The destination has a separate identity. The move stops the local source while
capturing its files, then creates a cloud sandbox. Running processes, memory,
and open sockets don't transfer.

The explicit timeout requests that the cloud sandbox stop after two hours and
retain its state. The move fails if that action isn't supported. Check the
destination before starting the task:

```console
$ sbx --cloud ttl CLOUD_ID
$ sbx --cloud secret ls
$ sbx --cloud policy allow network pypi.org:443 --sandbox CLOUD_ID
$ sbx --cloud policy allow network files.pythonhosted.org:443 --sandbox CLOUD_ID
$ sbx --cloud policy ls CLOUD_ID
$ sbx --cloud exec CLOUD_ID cat /home/agent/workspace/project/HANDOFF.md
$ sbx --cloud exec CLOUD_ID /home/agent/workspace/project/api-change-reporter/.venv/bin/python -m pip check
$ sbx --cloud exec CLOUD_ID curl -I https://pypi.org/simple/pyyaml/
$ sbx --cloud attach CLOUD_ID
```

Confirm that the expected cloud credential exists and that the agent responds.
If a connection fails, inspect
[cloud network policy decisions](/manuals/ai/sandboxes/cloud/network-policy.md#inspect-network-policy).
Local rules don't transfer. The destination uses cloud policy and the agent's
network rules. Rerun the baseline tests in the destination before making changes.

## Continue and disconnect

Give the cloud agent the task through the file:

```text
Read /home/agent/workspace/project/HANDOFF.md and rerun its baseline test command.
If it passes, implement the remaining task and run all regression tests.
Use the absolute interpreter path in HANDOFF.md rather than shell activation.

Save report.md with the baseline commit, changes, limitations, and results.
Save the full test output and a patch against the baseline, including any
added files. Package the updated source, requirements, examples, tests,
HANDOFF.md, patch, report, and test output into
/home/agent/workspace/result.zip. Exclude .venv, credentials, caches, and Git
metadata. Extract the bundle to a temporary directory and verify it before
declaring the task complete. If blocked, save the failure evidence instead.
```

`HANDOFF.md` carries the objective and preparation without depending on the
agent remembering an earlier conversation.

Once the agent starts working, press `Ctrl+\` to detach. Detaching leaves cloud
work running, so you can disconnect your laptop. The expiration still bounds
execution; after it stops the sandbox, the agent makes no further progress
until you resume it.

## Retrieve and keep the result

Return to inspect the result:

```console
$ sbx --cloud attach CLOUD_ID
```

Attaching starts a stopped sandbox. Review the report and test output, then
detach with `Ctrl+\`. Check the timer after resuming and copy the completed bundle
to your computer:

```console
$ sbx --cloud ttl CLOUD_ID
$ sbx --cloud cp CLOUD_ID:/home/agent/workspace/result.zip ./result.zip
```

Extract the ZIP and inspect the source, patch, report, and test output. Check
that regression tests cover added required parameters, optional parameters that
stay optional, and operation overrides. Keep both sandboxes until you've
verified the bundle, then remove them separately:

```console
$ sbx --cloud rm CLOUD_ID
$ sbx rm reporter-local
```

The source and destination don't synchronize after a move. Copying the bundle
is what brings the cloud changes back to your computer.

Optional. Before cleanup, move the cloud sandbox back when you want its files
and installed dependencies for local testing:

```console
$ sbx move CLOUD_ID --to local --name reporter-returned
```

Review the [architecture and disk-space requirements](/manuals/ai/sandboxes/cloud/move.md)
first. This creates another local sandbox; inspect its state and remove each
remaining sandbox when you're finished.
