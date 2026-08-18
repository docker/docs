---
description:
  Understand how to diagnose and troubleshoot Docker Desktop, and how to
  check the logs.
keywords: Linux, Mac, Windows, troubleshooting, logs, issues, Docker Desktop
toc_max: 2
title: Troubleshoot Docker Desktop
linkTitle: Troubleshoot and diagnose
aliases:
  - /desktop/troubleshoot/overview/
  - /desktop/troubleshoot/
tags: [Troubleshooting]
weight: 10
---

This page contains information on how to diagnose and troubleshoot Docker Desktop, and how to check the logs.

## Troubleshoot menu

To navigate to **Troubleshoot** either:

- Select the Docker menu Docker menu {{< inline-image src="../../images/whale-x.svg" alt="whale menu" >}} and then **Troubleshoot**.
- Select the **question mark** icon near the top-right corner of the Docker Desktop Dashboard.

The **Troubleshoot** menu contains the following options:

- **Restart Docker Desktop**.

- **Support**.

- **Reset Kubernetes cluster**. Select to delete all stacks and Kubernetes resources. For more information, see [Kubernetes](/manuals/desktop/settings-and-maintenance/settings.md#kubernetes).

- **Clean up  data**. This option resets all Docker data without a
  reset to factory defaults. Selecting this option results in the loss of existing settings.

- **Reset to factory defaults**: Choose this option to reset all options on
  Docker Desktop to their initial state, the same as when Docker Desktop was first installed.

If you are a Mac or Linux user, you also have the option to **Uninstall** Docker Desktop from your system.

## Diagnose

### Diagnose from the app

1. From **Troubleshoot**, select **Get support**. This opens the in-app Support page and starts collecting the diagnostics.
   Gathering diagnostics can take several minutes. Don't close Docker Desktop while the diagnostics are being collected.
2. When the diagnostics collection process is complete, select **Upload to get a Diagnostic ID**.
3. When the diagnostics are uploaded, Docker Desktop prints a diagnostic ID. Copy this ID.
4. Use your diagnostics ID to get help:
   - If you have a paid Docker subscription, select **Contact support**. This opens the Docker Desktop support form. Fill in the information required and add the ID you copied in step three to the **Diagnostics ID field**. Then, select **Submit ticket** to request Docker Desktop support.
     > [!NOTE]
     >
     > You must be signed in to Docker Desktop to access the support form. For information on what's covered as part of Docker Desktop support, see [Support](/manuals/support/_index.md).
   - If you don't have a paid Docker subscription, select **Report a Bug** to open a new Docker Desktop issue on GitHub. Complete the information required and ensure you add the diagnostic ID you copied in step three.

### Diagnose from an error message

1. When an error message appears, select **Gather diagnostics**. Gathering diagnostics may take several minutes. Don't close Docker Desktop while the diagnostics are being collected.
2. When the diagnostics are uploaded, Docker Desktop prints a diagnostic ID. Copy this ID.
3. Use your diagnostics ID to get help:
   - If you have a paid Docker subscription, select **Contact support**. This opens the Docker Desktop support form. Fill in the information required and add the ID you copied in step three to the **Diagnostics ID field**. Then, select **Submit ticket** to request Docker Desktop support.
     > [!NOTE]
     >
     > You must be signed in to Docker Desktop to access the support form. For information on what's covered as part of Docker Desktop support, see [Support](/manuals/support/_index.md).
   - If you don't have a paid Docker subscription, you can open a new [Docker Desktop issue on GitHub](https://github.com/docker/desktop-feedback). Complete the information required and ensure you add the diagnostic ID printed in step two.

### Diagnose from the terminal

In some cases, it's useful to run the diagnostics yourself, for instance, if the Docker Desktop Dashboard cannot start.

Run the [`docker desktop diagnose`](/manuals/desktop/features/desktop-cli.md) command:

```console
$ docker desktop diagnose
```

Gathering diagnostics may take several minutes. Wait for the process to complete before closing the terminal.

After the diagnostics have finished, the terminal displays your diagnostics ID and the path to the diagnostics file. The diagnostics ID is composed of your user ID and a timestamp. For example `BE9AFAAF-F68B-41D0-9D12-84760E6B8740/20190905152051`.

Alternatively, you can use the `com.docker.diagnose` tool:

{{< tabs group="os" >}}
{{< tab name="Windows" >}}

1. Locate the `com.docker.diagnose` tool:

   ```console
   # For all-user installations
   $ C:\Program Files\Docker\Docker\resources\com.docker.diagnose.exe

   # For per-user installations
   $ %LOCALAPPDATA%\Programs\DockerDesktop\resources\com.docker.diagnose.exe
   ```

2. Create and upload the diagnostics ID. In PowerShell, run:

   ```console
   # For all-user installations
   $ & "C:\Program Files\Docker\Docker\resources\com.docker.diagnose.exe" gather -upload

   # For per-user installations
   $ & %LOCALAPPDATA%\Programs\DockerDesktop\resources\com.docker.diagnose.exe" gather -upload
   ```

After the diagnostics have finished, the terminal displays your diagnostics ID and the path to the diagnostics file. The diagnostics ID is composed of your user ID and a timestamp. For example `BE9AFAAF-F68B-41D0-9D12-84760E6B8740/20190905152051`.

{{< /tab >}}
{{< tab name="Mac" >}}

1. Locate the `com.docker.diagnose` tool:

   ```console
   $ /Applications/Docker.app/Contents/MacOS/com.docker.diagnose
   ```

2. Create and upload the diagnostics ID. Run:

   ```console
   $ /Applications/Docker.app/Contents/MacOS/com.docker.diagnose gather -upload
   ```

After the diagnostics have finished, the terminal displays your diagnostics ID and the path to the diagnostics file. The diagnostics ID is composed of your user ID and a timestamp. For example `BE9AFAAF-F68B-41D0-9D12-84760E6B8740/20190905152051`.

{{< /tab >}}
{{< tab name="Linux" >}}

1. Locate the `com.docker.diagnose` tool:

   ```console
   $ /opt/docker-desktop/bin/com.docker.diagnose
   ```

2. Create and upload the diagnostics ID. Run:

   ```console
   $ /opt/docker-desktop/bin/com.docker.diagnose gather -upload
   ```

After the diagnostics have finished, the terminal displays your diagnostics ID and the path to the diagnostics file. The diagnostics ID is composed of your user ID and a timestamp. For example `BE9AFAAF-F68B-41D0-9D12-84760E6B8740/20190905152051`.

{{< /tab >}}
{{< /tabs >}}

To view the contents of the diagnostic file:

{{< tabs group="os" >}}
{{< tab name="Windows" >}}

1. Unzip the file. In PowerShell, copy and paste the path to the diagnostics file into the following command and then run it. It should be similar to the following example:

   ```powershell
   $ Expand-Archive -LiteralPath "C:\Users\testUser\AppData\Local\Temp\5DE9978A-3848-429E-8776-950FC869186F\20230607101602.zip" -DestinationPath "C:\Users\testuser\AppData\Local\Temp\5DE9978A-3848-429E-8776-950FC869186F\20230607101602"
   ```

2. Open the file in your preferred text editor. Run:

   ```powershell
   $ code <path-to-file>
   ```

{{< /tab >}}
{{< tab name="Mac" >}}

Run:

```console
$ open /tmp/<your-diagnostics-ID>.zip
```

{{< /tab >}}
{{< tab name="Linux" >}}

Run:

```console
$ unzip –l /tmp/<your-diagnostics-ID>.zip
```

{{< /tab >}}
{{< /tabs >}}

#### Use your diagnostics ID to get help

If you have a paid Docker subscription, [Contact support](https://app.docker.com/support/contact). Fill in the information required and add your diagnostics ID.

If you don't have a paid Docker subscription, create an issue on [GitHub](https://github.com/docker/desktop-feedback).

## Check the logs

In addition to using the diagnose option to submit logs, you can browse the logs yourself.

Run the [`docker desktop logs`](/manuals/desktop/features/desktop-cli.md) command to access Docker Desktop logs.

## View the Docker daemon logs

Refer to the [Read the daemon logs](/manuals/engine/daemon/logs.md) section
to learn how to view the Docker Daemon logs.

## Further resources

- View specific [troubleshoot topics](topics.md).
- View information on [known issues](known-issues.md)
- [Fix "Docker.app is damaged" on macOS](mac-damaged-dialog.md) - Resolve macOS installation issues
- [Get support for Docker products](/manuals/support/_index.md)
