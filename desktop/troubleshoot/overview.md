---
description: Understand how to diagnose and troubleshoot Docker Desktop, and how to check the logs.
keywords: Linux, Mac, Windows, troubleshooting, logs, issues, Docker Desktop
toc_max: 2
title: Overview
redirect_from:
- /desktop/linux/troubleshoot/
- /desktop/mac/troubleshoot/
- /desktop/windows/troubleshoot/
- /docker-for-mac/troubleshoot/
- /mackit/troubleshoot/
- /windows/troubleshoot/
- /docker-for-win/troubleshoot/
- /docker-for-windows/troubleshoot/
---

{% include upgrade-cta.html
  body="Docker Desktop offers support for developers on a paid Docker subscription (Pro, Team, or Business). Upgrade now to benefit from Docker Support. For more information, see [Support](../../support/index.md)."
  target-url="https://www.docker.com/pricing?utm_source=docker&utm_medium=webreferral&utm_campaign=docs_driven_upgrade_desktop_support"
%}

This page contains information on how to diagnose and troubleshoot Docker Desktop, and how to check the logs.

## Troubleshoot menu

To navigate to **Troubleshoot** either:

- Select the Docker menu ![whale menu](../images/whale-x.svg){: .inline} and then **Troubleshoot**
- Select the **Troubleshoot** icon near the top-right corner of Docker Dashboard

![Troubleshoot menu in Docker Desktop](../images/troubleshoot.png){:width="600px"}

The **Troubleshoot** page contains the following options:

- **Restart Docker Desktop**.

- **Get support**. Users with a paid Docker subscription can use this option to send a support request. Other users can use this option to diagnose any issues in Docker Desktop. For more information, see [Diagnose and feedback](#diagnose) and [Support](../../support/index.md).

- **Reset Kubernetes cluster**. Select to delete all stacks and Kubernetes resources. For more information, see [Kubernetes](../settings/linux.md#kubernetes).

- **Clean / Purge data**. This option resets all Docker data without a
reset to factory defaults. Selecting this option results in the loss of existing settings.

- **Reset to factory defaults**: Choose this option to reset all options on
Docker Desktop to their initial state, the same as when Docker Desktop was first installed.

If you are a Mac or Linux user, you also have the option to **Uninstall** Docker Desktop from your system.

## Diagnose

### Diagnose from the app

1. From **Troubleshoot**, select **Get support**. 
This opens the in-app **Support** page and starts collecting the diagnostics.
    ![Diagnose & Feedback](../images/diagnose-support.png){:width="600px"}
2. When the diagnostics collection process is complete, select **Upload to get a Diagnostic ID**.
3. When the diagnostics are uploaded, Docker Desktop prints a diagnostic ID. Copy this ID.
4. Use your diagnostics ID to get help:
    - If you have a paid Docker subscription, select **Contact Support**. This opens the [Docker Desktop support](https://hub.docker.com/support/desktop/){:target="_blank" rel="noopener" class="_"} form. Fill in the information required and add the ID you copied in step three to the **Diagnostics ID** field. Then, select **Submit** to request Docker Desktop support.
        > **Note**
        >
        > You must be signed in to Docker Desktop to access the support form. For information on what's covered as part of Docker Desktop support, see [Support](../../support/index.md).
    - If you don't have a paid Docker subscription, select **Report a Bug** to open a new Docker Desktop issue on GitHub. Complete the information required and ensure you add the diagnostic ID you copied in step three. 

### Diagnose from the terminal

In some cases, it's useful to run the diagnostics yourself, for instance, if
Docker Desktop cannot start.

<ul class="nav nav-tabs">
<li class="active"><a data-toggle="tab" data-target="#windows1">Windows</a></li>
<li><a data-toggle="tab" data-target="#mac1">Mac</a></li>
<li><a data-toggle="tab" data-target="#linux1">Linux</a></li>
</ul>
<div class="tab-content">
<div id="windows1" class="tab-pane fade in active" markdown="1">

1. Locate the `com.docker.diagnose` tool:

    ```console
    $ C:\Program Files\Docker\Docker\resources\com.docker.diagnose.exe
    ```

2. Create and upload the diagnostics ID. Run:

    ```console
    $ "C:\Program Files\Docker\Docker\resources\com.docker.diagnose.exe" gather -upload
    ```

After the diagnostics have finished, the terminal displays your diagnostics ID and the path to the diagnostics file. The diagnostics ID is composed of your user ID and a timestamp. For example `BE9AFAAF-F68B-41D0-9D12-84760E6B8740/20190905152051`. 

</div>
<div id="mac1" class="tab-pane fade" markdown="1">

1. Locate the `com.docker.diagnose` tool:

    ```console
    $ /Applications/Docker.app/Contents/MacOS/com.docker.diagnose
    ```

2. Create and upload the diagnostics ID. Run:

    ```console
    $ /Applications/Docker.app/Contents/MacOS/com.docker.diagnose gather -upload
    ```

After the diagnostics have finished, the terminal displays your diagnostics ID and the path to the diagnostics file. The diagnostics ID is composed of your user ID and a timestamp. For example `BE9AFAAF-F68B-41D0-9D12-84760E6B8740/20190905152051`. 

</div>
<div id="linux1" class="tab-pane fade" markdown="1">

1. Locate the `com.docker.diagnose` tool:

    ```console
    $ /opt/docker-desktop/bin/com.docker.diagnose
    ```

2. Create and upload the diagnostics ID. Run:

    ```console
    $ /opt/docker-desktop/bin/com.docker.diagnose gather -upload
    ```

After the diagnostics have finished, the terminal displays your diagnostics ID and the path to the diagnostics file. The diagnostics ID is composed of your user ID and a timestamp. For example `BE9AFAAF-F68B-41D0-9D12-84760E6B8740/20190905152051`. 
</div>
</div>

To view the contents of the diagnostic file:

<ul class="nav nav-tabs">
<li class="active"><a data-toggle="tab" data-target="#windows2">Windows</a></li>
<li><a data-toggle="tab" data-target="#mac2">Mac</a></li>
<li><a data-toggle="tab" data-target="#linux2">Linux</a></li>
</ul>
<div class="tab-content">
<div id="windows2" class="tab-pane fade in active" markdown="1">
<br>
1. Unzip the file. In PowerShell, copy and paste the path to the diagnostics file into the following command and then run it. It should be similar to the following example:

    ```powershell
    $ Expand-Archive -LiteralPath "C:\Users\testUser\AppData\Local\Temp\5DE9978A-3848-429E-8776-950FC869186F\20230607101602.zip" -DestinationPath "C:\Users\testuser\AppData\Local\Temp\5DE9978A-3848-429E-8776-950FC869186F\20230607101602"
     ```  
2. Open the file in your preferred text editor. Run:

    ```powershell
    $ code <path-to-file>
    ```

</div>
<div id="mac2" class="tab-pane fade" markdown="1">

Run:

    ```console
    $ open /tmp/<your-diagnostics-ID>.zip
    ``` 

</div>
<div id="linux2" class="tab-pane fade" markdown="1">

Run:

    ```console
    $ unzip –l /tmp/<your-diagnostics-ID>.zip
    ``` 
</div>
</div>

#### Use your diagnostics ID to get help

If you have a paid Docker subscription, open the [Docker Desktop support](https://hub.docker.com/support/desktop/){:target="_blank" rel="noopener" class="_"} form. Fill in the information required and add the ID to the Diagnostics ID field. Make sure you provide the full diagnostics ID, and not just the user ID. Select **Submit** to request Docker Desktop support.
    
If you don't have a paid Docker subscription, create an issue on GitHub:
 - [For Linux](https://github.com/docker/desktop-linux/issues){:target="_blank" rel="noopener" class="_"}
 - [For Mac](https://github.com/docker/for-mac/issues){:target="_blank" rel="noopener" class="_"}
 - [For Windows](https://github.com/docker/for-win/issues){:target="_blank" rel="noopener" class="_"}

### Self-diagnose tool

Docker Desktop contains a self-diagnose tool which can help you identify some common problems. 

 <ul class="nav nav-tabs">
<li class="active"><a data-toggle="tab" data-target="#windows3">Windows</a></li>
<li><a data-toggle="tab" data-target="#mac3">Mac</a></li>
<li><a data-toggle="tab" data-target="#linux3">Linux</a></li>
</ul>
<div class="tab-content">
<div id="windows3" class="tab-pane fade in active" markdown="1">
     
1. Locate the `com.docker.diagnose` tool. 
     
     ```console
    $ C:\Program Files\Docker\Docker\resources\com.docker.diagnose.exe
    ```

2. Run the self-diagnose tool:

    ```console
    $ "C:\Program Files\Docker\Docker\resources\com.docker.diagnose.exe" check
    ```

</div>
<div id="mac3" class="tab-pane fade" markdown="1">

1. Locate the `com.docker.diagnose` tool. 

    ```console
    $ /Applications/Docker.app/Contents/MacOS/com.docker.diagnose
     ```

2. Run the self-diagnose tool:

    ```console
    $ /Applications/Docker.app/Contents/MacOS/com.docker.diagnose check
    ```

</div>
<div id="linux3" class="tab-pane fade" markdown="1">

1. Locate the `com.docker.diagnose` tool. 

    ```console
    $ /opt/docker-desktop/bin/com.docker.diagnose
    ```

2. Run the self-diagnose tool:

    ```console
    $ /opt/docker-desktop/bin/com.docker.diagnose check
    ```
</div>
</div>

The tool runs a suite of checks and displays **PASS** or **FAIL** next to each check. If there are any failures, it highlights the most relevant at the end of the report.

You can then create an issue on GitHub:
    - [For Linux](https://github.com/docker/desktop-linux/issues){:target="_blank" rel="noopener" class="_"}
    - [For Mac](https://github.com/docker/for-mac/issues){:target="_blank" rel="noopener" class="_"}
    - [For Windows](https://github.com/docker/for-win/issues){:target="_blank" rel="noopener" class="_"}

## Check the logs

In addition to using the diagnose option to submit logs, you can browse the logs yourself.

<ul class="nav nav-tabs">
<li class="active"><a data-toggle="tab" data-target="#windows4">Windows</a></li>
<li><a data-toggle="tab" data-target="#mac4">Mac</a></li>
<li><a data-toggle="tab" data-target="#linux4">Linux</a></li>
</ul>
<div class="tab-content">
<div id="windows4" class="tab-pane fade in active" markdown="1">
<br>
In PowerShell, run:

```powershell
$ code $Env:LOCALAPPDATA\Docker\log
```

This opens up all the logs in your preferred text editor for you to explore.

</div>
<div id="mac4" class="tab-pane fade" markdown="1">

### From terminal

To watch the live flow of Docker Desktop logs in the command line, run the following script from your preferred shell.

```console
$ pred='process matches ".*(ocker|vpnkit).*" || (process in {"taskgated-helper", "launchservicesd", "kernel"} && eventMessage contains[c] "docker")'
$ /usr/bin/log stream --style syslog --level=debug --color=always --predicate "$pred"
```

Alternatively, to collect the last day of logs (`1d`) in a file, run:

```console
$ /usr/bin/log show --debug --info --style syslog --last 1d --predicate "$pred" >/tmp/logs.txt
```

### From the Console app

Mac provides a built-in log viewer, named **Console**, which you can use to check
Docker logs.

The Console lives in `/Applications/Utilities`. You can search for it with
Spotlight Search.

To read the Docker app log messages, type `docker` in the Console window search bar and press Enter. Then select `ANY` to expand the drop-down list next to your `docker` search entry, and select `Process`.

![Mac Console search for Docker app](../images/console.png)

You can use the Console Log Query to search logs, filter the results in various
ways, and create reports.

</div>
<div id="linux4" class="tab-pane fade" markdown="1">
<br>
You can access Docker Desktop logs by running the following command:

```console
$ journalctl --user --unit=docker-desktop
```

You can also find the logs for the internal components included in Docker
Desktop at `$HOME/.docker/desktop/log/`.

</div>
</div>

## View the Docker daemon logs

Refer to the [Read the daemon logs](../../config/daemon/logs.md) section
to learn how to view the Docker Daemon logs.

## Further resources

- View specific [troubleshoot topics](topics.md).
- Implement [workarounds for common problems](workarounds.md)
- View information on [known issues](known-issues.md)
