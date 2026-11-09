---
title: Run tasks on a schedule
linkTitle: Schedules
description: Schedule recurring agent tasks in Docker Agentic Platform, with a fresh sandbox for each run, MCP tools for external services, and a time limit for cleanup.
keywords: docker agentic platform, schedules, cron, sandboxes, automation, mcp
weight: 25
---

Schedules run a prompt or command in a sandbox at a recurring time. Use them
for recurring development work such as reviewing dependencies, triaging issues,
or summarizing pull requests, without starting an agent manually or keeping a
sandbox running between tasks.

Each run creates a fresh sandbox with the kit, compute size, credentials, and
tools saved in the schedule. The sandbox is deleted when the task finishes or
its time to live (TTL) expires, whichever happens first.

## Prepare tools and credentials

Choose a [kit](kits.md) that can run your task without interactive input. Save
the credentials you need in [Secrets](secrets.md), including model provider
API keys, service tokens, or custom secrets.

If you need to connect to external services, connect and authorize the relevant
[MCP servers](mcp.md) before creating the schedule. You can select these tools
for scheduled runs, including custom servers connected through the MCP gateway.
If the predefined servers don't cover your service, create your own MCP server
and [add its remote endpoint](mcp.md#add-a-server-by-url). The endpoint must
support streamable HTTP.

If you need credentials for a custom CLI or API integration, use the schedule
API to select [custom secrets](secrets.md#add-a-custom-secret). The Console
schedule form doesn't offer custom-secret selection.
Use an MCP server for integrations you configure through this form. When you
configure a schedule through the API, include the required CLI in your kit and use the
secret's placeholder where the client expects a credential. Check that your
[network policy](policies.md) permits access to the service. Don't put secret
values in the task's prompt.

## Create a schedule

1. Open **Schedules** in the Console and select **New schedule**.
2. Give the schedule a name that identifies its task, and optionally add a
   description.
3. Choose a kit with the agent or command-line tools your task needs. Make
   sure its required credentials are available so each run can start without
   interactive setup.
4. Configure the sandbox for the task:
   - Choose a compute size with enough resources to run your agent and workload.
   - Select the MCP tools the task needs to access external services. Connect
     or reconnect them if prompted so scheduled runs have authorized access.
   - Review the network policy. Scheduled sandboxes default to **Open**, which
     permits outbound access to any host. For tighter access, select a more
     restrictive [network policy](policies.md) that permits the services your
     task needs.
5. Set a time limit long enough for the task to finish. This time to live
   (TTL) limits how long each sandbox can run. The available choices are 1, 4,
   8, and 24 hours. A task still running at expiry is stopped, its sandbox is
   deleted, and the run counts as failed.
6. Set how often the task should run and choose the time zone for that
   schedule. Use a preset, or scroll to the bottom of the schedule menu and
   select **Custom** to enter a cron expression. Review the next run times to
   confirm the timing matches your intent.
7. Write the instructions for each run: a prompt for an agent kit, or a
   command for a command kit such as Shell. Specify where the task should
   save or send its results, because files left in the sandbox are deleted
   with it.
8. Select **Save schedule**.

Custom cron expressions have five fields: minute, hour, day of month, month,
and day of week. For example, `0 9 * * 1-5` runs at 09:00 every weekday in the
selected time zone.

## Use a custom kit

If you use a custom agent kit, configure its prompt template to run the agent
without interactive input. Docker Agentic Platform uses this template to build
the command for each scheduled run, substituting the prompt you entered into
the template's arguments. The image's `ENTRYPOINT` or `CMD` alone doesn't define
how a scheduled prompt is passed to the agent.

## Example: Weekly issue triage

To reduce time spent asking for missing details in bug reports, schedule an
agent to review issues and request reproduction steps. Connect an MCP server
that can read GitHub issues and comments and add comments to issues. Authorize
access to the repository, then select the server under **Tools**. If the
predefined servers don't provide these tools, add a custom remote server.

Choose an agent kit, set the cadence to `0 9 * * 1`, and select your team's
time zone to run every Monday at 09:00. Use a prompt such as:

```text
Use the connected GitHub tools to review open issues labeled "bug" in
OWNER/REPOSITORY. For each issue, read its description and comments and
check whether it includes reproduction steps, expected behavior, actual
behavior, and version information.
If details are missing and no comment has already requested them, add one
comment asking for the specific missing details. Skip issues awaiting a
response to an earlier request. Don't close issues or change their labels.
If a tool fails, stop and report the error instead of guessing.
```

Replace `OWNER/REPOSITORY` with your repository and `bug` with the label your
team uses for bug reports. The comments stay with the issues after the sandbox
is deleted, so reporters can respond and maintainers can follow up. The
connected server must provide tools to read issues and their comments and
add comments. Test the task with **Run now** before relying on its cadence.

## Manage schedules and review runs

Open a schedule to review its settings in **Overview** and its run history in
**Runs**. The history shows when each run started, its result, and duration.
Use **Run now** to trigger a run without waiting for the next scheduled time.

Use **Edit schedule** to change the task or its sandbox settings. Pause a
schedule to suspend future scheduled runs, and resume it when needed. You can
also duplicate a schedule or delete it when you no longer need it.

The schedule saves the kit and selected resources when you save its settings.
After adding credentials, updating a kit, or changing a custom MCP server's
URL, edit and save the schedule to use those changes. If a run fails, check
its reported error and confirm that its credentials, tools, and network access
are still valid.
