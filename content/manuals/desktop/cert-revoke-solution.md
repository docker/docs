---
description: Learn how to resolve issues affecting macOS users of Docker Desktop, including startup problems and false malware warnings, with upgrade, patch, and workaround solutions.
keywords: Docker desktop, fix, mac, troubleshooting, macos, false malware warning, patch, upgrade solution
title: Resolve the recent Docker Desktop issue on macOS
linkTitle: Fix startup issue for Mac
weight: 220
---

This guide provides steps to address a recent issue affecting some macOS users of Docker Desktop. The issue may prevent Docker Desktop from starting and in some cases, may also trigger inaccurate malware warnings. For more details about the incident, see the [blog post](https://www.docker.com/blog/incident-update-docker-desktop-for-mac/).

> [!NOTE]
>
> Docker Desktop versions 4.28 and earlier are not impacted by this issue. 

## Available solutions

There are a few options available depending on your situation:

### Upgrade to Docker Desktop version 4.37.2 (recommended)

The recommended way is to upgrade to the latest Docker Desktop version which is version 4.37.2. 

If possible, update directly through the app. If not, and you’re still seeing the malware pop-up, follow the steps below:

1. Kill the Docker process that cannot start properly:
   ```console
   $ sudo launchctl bootout system/com.docker.vmnetd 2>/dev/null || true
   $ sudo launchctl bootout system/com.docker.socket 2>/dev/null || true
    
   $ sudo rm /Library/PrivilegedHelperTools/com.docker.vmnetd || true
   $ sudo rm /Library/PrivilegedHelperTools/com.docker.socket || true
 
   $ ps aux | grep -i docker | awk '{print $2}' | sudo xargs kill -9 2>/dev/null
   ```
    
2. Make sure the malware pop-up is permanently closed. 

3. [Download and install version 4.37.2](/manuals/desktop/release-notes.md#4372).

4. Launch Docker Desktop. A privileged pop-up message displays after 5 to 10 seconds.

5. Enter your password.

You should now see the Docker Desktop Dashboard.

### Install a patch if you have version 4.34 - 4.36

If you can’t upgrade to the latest version and you’re seeing the malware pop-up, follow the steps below:

1. Kill the Docker process that cannot start properly:
   ```console
   $ sudo launchctl bootout system/com.docker.vmnetd 2>/dev/null || true
   $ sudo launchctl bootout system/com.docker.socket 2>/dev/null || true
    
   $ sudo rm /Library/PrivilegedHelperTools/com.docker.vmnetd || true
   $ sudo rm /Library/PrivilegedHelperTools/com.docker.socket || true
 
   $ ps aux | grep docker | awk '{print $2}' | sudo xargs kill -9 2>/dev/null
   ```

2. Make sure the malware pop-up is permanently closed.

3. [Download and install the patched installer](/manuals/desktop/release-notes.md) that matches your current base version. For example if you have version 4.36.0, install 4.36.1.

4. Launch Docker Desktop. A privileged pop-up message displays after 5 to 10 seconds.

5. Enter your password.

You should now see the Docker Desktop Dashboard.

### Wait for a patch for versions 4.32 - 4.33

For versions 4.32 - 4.33, a patch fix is in progress. If you need an immediate solution, you can use the following workaround:

1. Kill the Docker process that cannot start properly:
   ```console
   $ sudo launchctl bootout system/com.docker.vmnetd 2>/dev/null || true
   $ sudo launchctl bootout system/com.docker.socket 2>/dev/null || true
    
   $ sudo rm /Library/PrivilegedHelperTools/com.docker.vmnetd || true
   $ sudo rm /Library/PrivilegedHelperTools/com.docker.socket || true
 
   $ ps aux | grep -i docker | awk '{print $2}' | sudo xargs kill -9 2>/dev/null
   ```

2. Download and install a re-signed installer matching your exact version of Docker Desktop from the [Release notes](/manuals/desktop/release-notes.md).

3. Install new binaries:

   ```console
   $ sudo cp /Applications/Docker.app/Contents/Library/LaunchServices/com.docker.vmnetd /Library/PrivilegedHelperTools/
   $ sudo cp /Applications/Docker.app/Contents/MacOS/com.docker.socket /Library/PrivilegedHelperTools/
   ```

4. Launch Docker Desktop. A privileged pop-up message displays after 5 to 10 seconds.

5. Enter your password.

You should now see the Docker Desktop Dashboard.

## MDM script

If you are an IT administrator, you can use the following script as a workaround for your developers if they have a re-signed version of Docker Desktop version 4.35 or later.

```console
#!/bin/bash

# Stop the docker services
echo "Stopping Docker..."
sudo pkill [dD]ocker

# Stop the vmnetd service
echo "Stopping com.docker.vmnetd service..."
sudo launchctl bootout system /Library/LaunchDaemons/com.docker.vmnetd.plist

# Stop the socket service
echo "Stopping com.docker.socket service..."
sudo launchctl bootout system /Library/LaunchDaemons/com.docker.socket.plist

# Remove vmnetd binary
echo "Removing com.docker.vmnetd binary..."
sudo rm -f /Library/PrivilegedHelperTools/com.docker.vmnetd

# Remove socket binary
echo "Removing com.docker.socket binary..."
sudo rm -f /Library/PrivilegedHelperTools/com.docker.socket

# Install new binaries
echo "Install new binaries..."
sudo cp /Applications/Docker.app/Contents/Library/LaunchServices/com.docker.vmnetd /Library/PrivilegedHelperTools/
sudo cp /Applications/Docker.app/Contents/MacOS/com.docker.socket /Library/PrivilegedHelperTools/
```
