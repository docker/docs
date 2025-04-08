---
description: Understand what Docker Desktop Resource Saver mode is and how to configure it
keywords: Docker Dashboard, resource saver, manage, containers, gui, dashboard, user manual
title: Docker Desktop's Resource Saver mode
linkTitle: Resource Saver mode
weight: 50
---

Resource Saver mode significantly reduces Docker
Desktop's CPU and memory utilization on the host by 2 GBs or more, by
automatically stopping the Docker Desktop Linux VM when no containers are
running for a period of time. The default time is set to 5 minutes, but this can be adjusted to suit your needs.

With Resource Saver mode, Docker Desktop uses minimal system resources when it's idle, thereby
allowing you to save battery life on your laptop and improve your multi-tasking
experience.

## Configure Resource Saver 

Resource Saver is enabled by default but can be disabled by navigating to the **Resources** tab, in **Settings**. You can also configure the idle
timer as shown below.

![Resource Saver Settings](../images/resource-saver-settings.png)

If the values available aren't sufficient for your
needs, you can reconfigure it to any value, as long as the value is larger than 30 seconds, by
changing `autoPauseTimeoutSeconds` in the Docker Desktop `settings-store.json` file (or `settings.json` for Docker Desktop versions 4.34 and earlier): 

  - Mac: `~/Library/Group Containers/group.com.docker/settings-store.json`
  - Windows: `C:\Users\[USERNAME]\AppData\Roaming\Docker\settings-store.json`
  - Linux: `~/.docker/desktop/settings-store.json`

There's no need to restart Docker Desktop after reconfiguring. 

When Docker Desktop enters Resource Saver mode: 
- A leaf icon displays on the
Docker Desktop status bar as well as on the Docker icon in
the system tray. The following image shows the Linux VM CPU and memory utilization reduced
to zero when Resource Saver mode is on. 

   ![Resource Saver Status Bar](../images/resource-saver-status-bar.png)

- Docker commands that don't run containers, for example listing container images or volumes, don't necessarily trigger an exit from Resource Saver mode as Docker Desktop can serve such commands without unnecessarily waking up the Linux VM.

> [!NOTE]
>
> Docker Desktop exits the Resource Saver mode automatically when it needs to.
> Commands that cause an exit from Resource Saver take a little longer to execute
> (about 3 to 10 seconds) as Docker Desktop restarts the Linux VM.
> It's generally faster on Mac and Linux, and slower on Windows with Hyper-V.
> Once the Linux VM is restarted, subsequent container runs occur immediately as usual.

## Resource Saver mode versus Pause

Resource Saver has higher precedence than the older [Pause](pause.md) feature,
meaning that while Docker Desktop is in Resource Saver mode, manually pausing
Docker Desktop is not possible (nor does it make sense since Resource Saver
actually stops the Docker Desktop Linux VM). In general, we recommend keeping
Resource Saver enabled as opposed to disabling it and using the manual Pause
feature, as it results in much better CPU and memory savings.

## Resource Saver mode on Windows

Resource Saver works a bit differently on Windows with WSL. Instead of
stopping the WSL VM, it only pauses the Docker Engine inside the
`docker-desktop` WSL distribution. That's because in WSL there's a single Linux VM
shared by all WSL distributions, so Docker Desktop can't stop the Linux VM (i.e.,
the WSL Linux VM is not owned by Docker Desktop). As a result, Resource Saver
reduces CPU utilization on WSL, but it does not reduce Docker's memory
utilization. 

To reduce memory utilization on WSL, we instead recommend that
users enable WSL's `autoMemoryReclaim` feature as described in the
[Docker  Desktop WSL docs](/manuals/desktop/features/wsl/_index.md). Finally, since Docker Desktop does not
stop the Linux VM on WSL, exit from Resource Saver mode is immediate (there's
no exit delay).
