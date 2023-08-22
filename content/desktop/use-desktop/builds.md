---
title: Explore Builds (Beta)
description: Understand how to use the Builds view in Docker Desktop
keywords: Docker Dashboard, manage, gui, dashboard, builders, builds
---

The **Builds** view is a simple interface that lets you inspect your build
history and manage builders using Docker Desktop. By default, it
displays a list of all your ongoing and completed builds. 

> **Beta feature**
>
> The **Builds** view is currently in [Beta](../../release-lifecycle.md/#beta).
> This feature may change or be removed from future releases.
{ .experimental }


The **Builds** view displays metadata about the build, such as:

- Build name
- Target platforms
- Builder name
- Build duration
- Progress bar
- Cache usage
- Completion date

The **Active builds** section displays builds that are currently running on
builders that you're using.

The **Completed builds** section lists build records for past builds for your
active builders. The list doesn't include builds for inactive builders.

## Turn on the Builds view

1. Navigate to **Settings**.
2. Select **Features in development**.
3. In the **Beta features** tab, select the **Display Builds view** checkbox.
4. Select **Apply & restart** for the changes to take effect.

After the restart, the **Builds** view and the **Builders** settings menu
appear.

## Inspect a build

To inspect a build, select the build that you want to view in the list.

The **Info** tab displays details about the build job. The details include
information such as target stage for multi-stage builds, target platforms, and
version control information, if available.

The **Source** tab shows the [frontend](../../build/dockerfile/frontend.md)
used to create the build.

The **Error** tab appears if the build finished with an error. It displays the
[frontend](../../build/dockerfile/frontend.md) used to create the build, and
the build error displays inline in the frontend source.

The **Logs** tab displays the build logs. If the build is currently running,
the logs are updated in real-time.

The **Stats** tab displays statistics data about completed builds. Analyze the
build stats to get a better understanding of how your build gets executed, and
find ways to optimize it.

## Manage builders

To inspect your builders, and change your default builder, select
**Builder settings** to open the settings menu. For more information, see:

- [Change settings, Windows](../settings/windows.md#builders)
- [Change settings, Mac](../settings/mac.md#builders)
- [Change settings, Linux](../settings/linux.md#builders)