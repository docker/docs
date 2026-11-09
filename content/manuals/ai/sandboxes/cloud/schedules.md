---
title: Run tasks on a schedule
linkTitle: Schedules
description: Create, test, and manage recurring agent tasks with the cloud Sandboxes CLI, using a fresh sandbox for each run and saving results outside the sandbox.
keywords: docker sandboxes, cloud sandboxes, schedules, cron, automation, sbx schedule
weight: 35
---

Schedules run unattended agent tasks at recurring times without keeping your
machine online. Each run creates a fresh cloud sandbox from a saved kit and
configuration. The sandbox is deleted when the task finishes or its time to
live (TTL) expires, whichever happens first. Instruct the task to save results
to an external destination before it exits.

Use schedules for recurring work such as dependency reviews or issue summaries.
For jobs triggered by your CI system, see [CI and headless use](../workflows/automation.md).
For interactive sessions configured from a file, see
[Cloud environment files](environment-files.md). To manage schedules in the
web interface, see [Schedules in Docker Agentic Platform](/manuals/agentic-platform/schedules.md).

## Create a schedule

Schedules require [cloud access](_index.md), scheduling access for your account,
and a CLI that provides `sbx --cloud schedule --help`. Sign in with `sbx login`
to authorize unattended runs. See [Renew schedule credentials](#renew-schedule-credentials)
for credential lifetime and login requirements.

Choose a published kit such as `claude` or an `oci://` reference. Local kit
paths and `git+` references aren't supported. The kit must include the tools
the task needs and accept a prompt or command without interactive input.

Configure [cloud credentials](credentials.md) for the agent and services it
uses. For example, an issue-summary task needs an Anthropic API key and a
GitHub token with permission to read issues and post comments:

```console
$ sbx --cloud secret set anthropic
$ sbx --cloud secret set github
```

The following schedule posts a weekly summary to a GitHub tracking issue, so
its output survives sandbox deletion. Replace the repository and issue
placeholders with a destination you control. `--enabled=false` keeps recurring
runs paused until you [test the task](#test-a-schedule).

```console
$ sbx --cloud schedule create claude \
    --name weekly-issue-summary \
    --cron '0 9 * * 1' \
    --timezone Etc/UTC \
    --enabled=false \
    --max-run-duration 1h \
    --workdir /home/agent \
    --env REPOSITORY=<OWNER>/<REPOSITORY> \
    --env TRACKING_ISSUE=<ISSUE_NUMBER> \
    --allow-network api.github.com:443 \
    --prompt 'Use the GitHub REST API to summarize up to ten open issues in REPOSITORY, excluding pull requests and TRACKING_ISSUE. Post their titles and links as one comment on TRACKING_ISSUE. Skip duplicate summaries. Stop on API errors. Do not modify code or other issues, or ask for input.'
```

Keep the returned schedule ID. Commands take IDs, rather than schedule names.
Pass exactly one of `--prompt` for an agent task or `--command` for a shell
command.

### Set the cadence

`--cron` takes a quoted five-field expression: minute, hour, day of month,
month, and day of week. For example, `'0 9 * * 1'` runs every Monday at 09:00.
Set `--timezone` explicitly to an IANA time zone such as `Etc/UTC`.

### Configure each run

| Setting | Behavior |
| --- | --- |
| `--max-run-duration` | Run time limit, at least `1h`. Defaults to `1h`. Use this instead of `--ttl`. |
| `--cpus`, `--memory` | Override the kit's resource defaults. |
| `--workdir`, `-w` | Directory in the sandbox where the task runs. `--working-dir` is an alias. The directory must exist. This doesn't upload a host project. |
| `--env KEY=VALUE` | Set task parameters. A bare `--env KEY` is rejected. Store credentials in the secret store. |
| `--allow-network` | Declare destinations the task needs. See [cloud network policy](network-policy.md) and test service access. |

Each run starts without your host workspace. Have the task fetch its input
from an external service. Don't assume local network policy transfers to
scheduled sandboxes or that a configured rule alone proves enforcement.

## Test a schedule

Queue a manual run, including while recurring runs are paused:

```console
$ sbx --cloud schedule run <SCHEDULE_ID>
$ sbx --cloud schedule runs <SCHEDULE_ID>
$ sbx --cloud schedule runs inspect <SCHEDULE_ID> <RUN_ID>
```

Use the run ID from `run` or `runs`. Inspection shows the state, failure reason,
exit code when available, and timing. Queuing a run doesn't establish success:
wait for completion and check the external result, such as the comment on the
tracking issue.

If another run is active, the service cancels the additional run as overlapping.
Wait for the active run to finish before retrying.

## List and inspect schedules

List schedules and inspect their cadence, status, kit, and next run:

```console
$ sbx --cloud schedule ls
$ sbx --cloud schedule inspect <SCHEDULE_ID>
```

## Enable or pause recurring runs

Enable a schedule after verifying a test run:

```console
$ sbx --cloud schedule update <SCHEDULE_ID> --enabled=true
```

Pause recurring runs with `--enabled=false`:

```console
$ sbx --cloud schedule update <SCHEDULE_ID> --enabled=false
```

## Change schedule settings

Updates to `--name`, `--description`, `--cron`, `--timezone`, or `--enabled`
preserve the saved run configuration. For example, move the weekly run to
10:00 UTC:

```console
$ sbx --cloud schedule update <SCHEDULE_ID> --cron '0 10 * * 1' --timezone Etc/UTC
```

## Replace the task or run configuration

Passing a kit or any run flag replaces the saved run configuration. Supply
`--prompt` or `--command` and all run settings you want to retain, even when
changing only one value. Omitted settings reset, including network rules,
duration, resources, platform, environment variables, working directory, and
`--static-mcp` server selections.

For example, this replaces the task and sets a two-hour limit. It also resets
any previously configured environment variables, network rules, and other
omitted run settings:

```console
$ sbx --cloud schedule update <SCHEDULE_ID> claude \
    --max-run-duration 2h \
    --prompt 'Check the installed tool versions and report them. Do not ask for input.'
```

To retain the issue-summary task while changing its duration, copy its full
`create` command, replace `create claude` with `update <SCHEDULE_ID> claude`,
and change `--max-run-duration` to `2h`. Keep the prompt and every other run
flag. The CLI warns about resets but doesn't preserve omitted settings.
Pause recurring runs before replacing the configuration, then test before
resuming.

A schedule saves the resolved kit and selected account secret names. Replace
the run configuration to pick up published kit changes or added credentials.
Without a kit argument, replacement resolves the existing kit reference again.
Changing only schedule settings doesn't refresh the kit or secrets.

The CLI refuses to replace run configurations that store multiple kits,
policy IDs, or a sandbox environment, because its flags can't preserve those
settings. Use the schedule API to update those configurations. You can still
change their name, description, cadence, time zone, and enabled state with
the CLI.

## Renew schedule credentials

Interactive `sbx login` enrolls a schedule credential when scheduling is enabled
for your account. The credential expires after 90 days. Sign in interactively
again to renew it. PAT login with `--password-stdin` doesn't enroll this
credential.

A successful login can still report `schedule credential exchange failed`.
Resolve that warning and test a run before relying on the schedule.

## Recover from a failed run

Inspect the run's failure reason and exit code, and check the external result.
For an expired Docker credential, run `sbx login` again. For authentication or
connection failures, check agent credentials and service access. If you add a
secret or change the kit, [replace the run configuration](#replace-the-task-or-run-configuration)
with all required settings. Test again before resuming recurring runs.

## Remove a schedule

```console
$ sbx --cloud schedule rm <SCHEDULE_ID>
```

Removal asks for confirmation. `--force` skips the prompt. Removing a schedule
prevents further runs but doesn't immediately delete an active run's sandbox.
The active sandbox can remain until its timeout.
