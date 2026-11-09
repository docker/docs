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

## Prepare the task

Follow the [cloud setup instructions](_index.md) to activate cloud access. You
also need scheduling access for your account and an `sbx` CLI that provides
`sbx --cloud schedule --help`.

Sign in interactively with the Docker account that will own the schedule:

```console
$ sbx login
```

When scheduling is enabled for your account, login enrolls a credential that
lets the schedule service run tasks while you're offline. This credential
expires after 90 days; signing in interactively again renews it. PAT login with
`--password-stdin` doesn't enroll this credential. A successful login can still
report `schedule credential exchange failed`. Resolve that warning and test a
run before relying on the schedule.

Choose a published kit, such as `claude`, or an `oci://` kit reference. Local
kit paths and `git+` references aren't supported. A custom agent kit needs a
prompt template that runs without interactive input.

Configure [account-scoped cloud credentials](credentials.md) before creating
the schedule. For the example in this guide, store an Anthropic API key and a
GitHub token that can read issues and add comments in your target repository:

```console
$ sbx --cloud secret set anthropic
$ sbx --cloud secret set github
```

Each run starts without your host workspace. Prepare any required tools in
the kit, and have the task fetch its input from an external service. Keep
credentials in the secret store, outside prompts and environment variables.

## Create a paused schedule

This example asks Claude to summarize up to ten open issues and add one
comment to a tracking issue in GitHub. That comment persists after the sandbox
is deleted. Replace `<OWNER>/<REPOSITORY>` and `<ISSUE_NUMBER>` with a repository
and tracking issue you control.

Create the schedule paused so you can test it before enabling recurring runs:

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
    --prompt 'Read the REPOSITORY and TRACKING_ISSUE environment variables. Use the GitHub REST API to read up to ten open issues in REPOSITORY, excluding pull requests and TRACKING_ISSUE. Summarize their titles, links, and missing reproduction details. Add one summary comment to TRACKING_ISSUE in that repository. Do not modify other issues or code. If the same summary is already present, skip posting it. If an API call fails, stop and report the error. Do not ask for interactive input.'
```

The cron expression has five fields: minute, hour, day of month, month, and
day of week. Here, `0 9 * * 1` means every Monday at 09:00 in the specified
IANA time zone. Set `--timezone` explicitly and inspect the cadence before
enabling the schedule.

Pass exactly one of `--prompt` or `--command`. Use `--prompt` for an agent task
or `--command` for a shell command that the kit can execute without interaction.

The example sets a one-hour time limit with `--max-run-duration`. The CLI
requires at least one hour and defaults to one hour when omitted. Use this flag
instead of `--ttl`; scheduled sandboxes are deleted at timeout. CPU and memory
settings default to the kit's values. Override them with `--cpus` and `--memory`
when the workload needs different resources.

`--workdir` (`-w`, also available as `--working-dir`) selects the directory
where the prompt or command runs. The directory must exist in the sandbox;
this flag doesn't upload a host project. Use explicit `--env KEY=VALUE` pairs
for task parameters. A bare `--env KEY` is rejected.

Declare the destinations your task needs with `--allow-network`. Review
[cloud network policy](network-policy.md) and verify the task can reach its
services during the test run. Don't assume local policy transfers to scheduled
sandboxes or rely on a configured rule alone as proof of enforcement.

## Test and enable the schedule

Keep the schedule ID returned by `create`. Schedule commands take IDs, rather
than schedule names. Inspect the configuration and queue a manual run while
the schedule is paused:

```console
$ sbx --cloud schedule inspect <SCHEDULE_ID>
$ sbx --cloud schedule run <SCHEDULE_ID>
$ sbx --cloud schedule runs <SCHEDULE_ID>
$ sbx --cloud schedule runs inspect <SCHEDULE_ID> <RUN_ID>
```

Use the run ID from `run` or `runs`. Run details include the state, failure
reason, exit code when available, and timing. Queuing a run doesn't establish
that the task succeeded. Wait for completion, inspect the result, and verify
the comment appeared on the tracking issue before enabling the schedule:

```console
$ sbx --cloud schedule update <SCHEDULE_ID> --enabled=true
```

If a run is still active when another run is requested, the service cancels
the additional run as overlapping. Wait for the active run to finish before
retrying a manual test.

## Manage and update schedules

List schedules and inspect their cadence, status, kit, and next run:

```console
$ sbx --cloud schedule ls
$ sbx --cloud schedule inspect <SCHEDULE_ID>
```

Updates to `--name`, `--description`, `--cron`, `--timezone`, or `--enabled`
change only those settings. For example, pause recurring runs, change the
cadence, and resume:

```console
$ sbx --cloud schedule update <SCHEDULE_ID> --enabled=false
$ sbx --cloud schedule update <SCHEDULE_ID> --cron '0 10 * * 1' --timezone Etc/UTC
$ sbx --cloud schedule update <SCHEDULE_ID> --enabled=true
```

### Change what each run executes

Passing a kit or any run flag rebuilds the saved run configuration. Repeat
`--prompt` or `--command` and every run setting you want to keep. Omitted run
settings reset, including network rules, duration, resources, platform,
environment variables, working directory, and `--static-mcp` server selections.
The CLI warns about settings that reset; the warning doesn't preserve them.

For the walkthrough, start from the full `create` command, replace
`create claude` with `update <SCHEDULE_ID> claude`, and keep all run flags when
changing the prompt or another run setting. Test the rebuilt schedule while
paused before enabling it again.

A schedule saves the resolved kit and selected account secret names. Rebuild
the run configuration after publishing kit changes or adding credentials that
the task needs. Without an explicit kit argument, a rebuild resolves the
schedule's existing kit reference again. Changing only the cadence or enabled
state doesn't refresh the saved run configuration.

The CLI refuses to rebuild schedules that store multiple kits, policy IDs,
or a sandbox environment, because CLI flags can't preserve those settings.
Use the schedule API to preserve and update that configuration. You can still
change the five schedule settings listed earlier through the CLI.

### Recover from a failed run

Inspect the run's failure reason and exit code, and check the external result.
If the stored Docker credential expired, run `sbx login` again. Check agent
credentials and service access for authentication or connection failures.
If you add a missing secret or change the kit, rebuild the saved run
configuration with all required settings. Then queue another manual run and
verify its result before resuming recurring runs.

### Remove a schedule

Remove a schedule when you no longer need future runs:

```console
$ sbx --cloud schedule rm <SCHEDULE_ID>
```

Removal asks for confirmation; `--force` skips the prompt. Removing a schedule
prevents further runs but doesn't immediately delete an active run's sandbox.
The active sandbox can remain until its timeout. Save results externally as
part of the task rather than relying on copying files after a run finishes.
