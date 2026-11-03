---
title: Run tasks on a schedule
linkTitle: Schedules
description: Schedule recurring agent tasks in Docker Agentic Platform, with a fresh sandbox for each run, MCP tools for external services, and a time limit for cleanup.
keywords: docker agentic platform, schedules, cron, sandboxes, automation, mcp
weight: 25
---

Schedules run a prompt or command in a sandbox at a recurring time. Use them
for work such as a daily Jira digest, an email summary, or a weekly dependency
review, without starting an agent manually or keeping a sandbox running
between tasks.

Each run creates a fresh sandbox with the kit, compute size, credentials, and
tools saved in the schedule. When the task finishes, the sandbox is deleted.
The time limit also deletes the sandbox if the task hasn't finished.

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
2. Enter a **Name** and, optionally, a **Description**.
3. Under **Agent**, choose a kit and check its credential status. Add any
   missing credentials through **Secrets**.
4. Under **Sandbox**, choose the network **Policy**, **Tools**, and **Compute**
   size for each run. Connect or reconnect selected tools if prompted.
5. Choose a **Time limit** for each sandbox. This time to live (TTL) limits
   how long a task can run. The available choices are 1, 4, 8, and 24 hours.
   A task still running at expiry is stopped, its sandbox is deleted, and the
   run counts as failed.
6. Under **Schedule**, select a **Runs** preset and **Timezone**. For a cron
   expression, choose **Custom**. Review **Next runs** to check the timing.
7. In **Prompt**, enter the prompt for an agent kit, or the command for a
   command kit such as Shell. Specify where the task should save or send its
   results, because files left in the sandbox are deleted with it.
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

## Example: Daily Jira digest

To prepare a daily summary of your team's work, connect an MCP server that
can search Jira and a server that can send email. These can be custom remote
servers. Authorize access to the project and mailbox the task needs, then
select both servers under **Tools**.

Choose an agent kit, set the cadence to `0 9 * * 1-5`, and select your team's
time zone. Use a prompt such as:

```text
Use the connected Jira tool to find issues in project EXAMPLE updated in
the past 24 hours. Summarize progress, blockers, and issues awaiting review.
Include issue keys and links. Use the connected email tool to send the digest
to team@example.com with the subject "Daily EXAMPLE digest".
If either service fails, report the error instead of inventing a summary.
```

Replace the project and recipient with your own values. The prompt describes
the task; the connected servers must provide the corresponding search and
email tools. Test the task with **Run now** before relying on its cadence.

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
