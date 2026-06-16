---
description: Understand what you can do with the Logs view on Docker Dashboard
keywords: Docker Dashboard, manage, logs, gui, build logs, container logs, debugging, dashboard
title: Explore the Logs view in Docker Desktop
linkTitle: Logs
weight: 60
---

{{< summary-bar feature_name="Desktop logs" >}}

The **Logs** view provides a unified, real-time log stream from all containers and recent builds in Docker Desktop. Unlike the logs accessible from the [**Containers** view](container.md), the **Logs** view lets you monitor and search log output (up to a maximum of 100 000 entries) across your entire environment from a single interface. 

## Log entries

Each log entry in the table view shows:

| Column        | Description                                                                    |
| ------------- | ------------------------------------------------------------------------------ |
| **Timestamp** | The date and time the log line was emitted, for example `2026-02-26 11:18:53`. |
| **Object**    | The container or build that produced the log line.                             |
| **Message**   | The full log message, including any status codes such as `[ OK ]`.             |

Selecting the expand arrow to the right of a row reveals the full message for that entry.

## Search, filter, and export logs

Use the **Search** field at the top of the Logs view to find specific entries. The search bar supports:

- Plain-text terms for exact match searches
- Regular expressions (for example, `/error|warn/`)

You can save your current filters as a preset for easy access later. Presets capture your container selection, build log visibility, and case sensitivity settings, as well as any active search terms. If no containers are selected, the preset is named all; otherwise it is named after the first selected container.

To refine the log stream further, select the **Filter** icon in the toolbar to open the container filter panel. From here you can:

- Check individual containers to show only their output
- Check Compose stacks to show or hide entire groups
- Toggle off **View build logs** to exclude build-related log output in the stream
- Use **Select all** or **Clear container filters** to quickly toggle every container at once

Use the **Export** button in the top-right corner (available with Docker Desktop version 4.77 and later) to export all logs or only the logs that match your filters.

## Display settings

Select the **Display settings** icon in the toolbar to toggle the following:
- **Wrap lines**
- **Show timestamps**

You can also choose **Clear logs** to remove log entries from the view (available with Docker Dekstop 4.79 and later). A dialog lets you choose between two options:

- **Clear all logs**: Immediately hides all current log entries.
- **Clear logs before**: Activates a date and time picker. Only entries at or before the selected timestamp are hidden.

The cleared state persists across Docker Desktop restarts. Once entries are cleared, they cannot be restored.

## Feedback

Select **Give feedback** at the top of the view to share suggestions or report issues.
