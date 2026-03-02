---
description: Understand what you can do with the Logs view on Docker Dashboard
keywords: Docker Dashboard, manage, logs, gui, build logs, container logs, debugging, dashboard
title: Explore the Logs view in Docker Desktop
linkTitle: Logs
weight: 60
params:
  sidebar:
    badge:
      color: blue
      text: Beta
---

{{< summary-bar feature_name="Desktop logs" >}}


The **Logs** view provides a unified, real-time log stream from all running containers and Kubernetes nodes in Docker Desktop. Unlike the logs accessible from the [**Containers** view](container.md), the **Logs** view lets you monitor and search log output across your entire environment from a single interface. You can also use Ask Gordon to analyze log output with AI-assisted prompts.

## Log entries

Each log entry in the table view shows:

| Column | Description |
|---|---|
| **Timestamp** | The date and time the log line was emitted, for example `2026-02-26 11:18:53`. |
| **Object** | The container or node that produced the log line. |
| **Message** | The full log message, including any status codes such as `[ OK ]`. |

Selecting the expand arrow to the left of a row reveals the full message for that entry.

## Search and filter logs

Use the **Search** field at the top of the Logs view to find specific entries. The search bar supports:

- Plain-text terms for exact match searches
- Regular expressions (for example, `/error|warn/`)

You can save your search terms for easy-access later.

To refine the log stream further, select the **Filter** icon in the toolbar to open the container filter panel. From here you can:

- Check individual running containers to show only their output
- Check **Running containers** or **Stopped containers** to show or hide entire groups
- Use **Select all** or **Clear all** to quickly toggle every container at once

## Ask Gordon

The Logs view integrates with [Ask Gordon](/manuals/ai/gordon/_index.md), Docker's AI assistant. Select the **Ask Gordon** button in the top-right corner of the view to open the Gordon panel, then choose a scope:

- **All visible logs**: Gordon analyzes the complete log stream currently visible on screen
- **Container**: Gordon analyzes logs from a specific container
- **Build**: Gordon analyzes build-related log output

You can type a free-form question or select one of the suggested prompts.

## Display options

Select the **Display options** icon in the toolbar to toggle the following:

- **View build logs**: Include or exclude build-related log output in the stream
- **Table view** — switch between a structured table layout and a plain log stream

The table view is useful when you need to correlate events across multiple containers because each row clearly shows which container emitted a given message and when.

## Feedback

The Logs view is in active development. Select **Give feedback** at the top of the view to share suggestions or report issues.
