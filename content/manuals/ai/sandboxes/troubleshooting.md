---
title: Troubleshooting
description: Resolve common issues when sandboxing agents locally.
weight: 50
---

{{< summary-bar feature_name="Docker Sandboxes" >}}

This guide helps you resolve common issues when using Docker Sandboxes with AI coding agents.

<!-- vale off -->

## 'sandbox' is not a docker command

<!-- vale on -->

When you run `docker sandbox`, you see an error saying the command doesn't exist.

This means the CLI plugin isn't installed or isn't in the correct location. To fix:

1. Verify the plugin exists:

   ```console
   $ ls -la ~/.docker/cli-plugins/docker-sandbox
   ```

   The file should exist and be executable.

2. If using Docker Desktop, restart it to detect the plugin.

## "Experimental Features" needs to be enabled by your administrator

You see an error about beta features being disabled when trying to use sandboxes.

This happens when your Docker Desktop installation is managed by an
administrator who has locked settings. If your organization uses [Settings Management](/enterprise/security/hardened-desktop/settings-management/),
ask your administrator to [allow beta features](/enterprise/security/hardened-desktop/settings-management/configure-json-file/#beta-features):

```json
{
  "configurationFileVersion": 2,
  "allowBetaFeatures": {
    "locked": false,
    "value": true
  }
}
```

## Authentication failure

Claude can't authenticate, or you see API key errors.

The API key is likely invalid, expired, or not configured correctly.

## Workspace contains API key configuration

You see a warning about conflicting credentials when starting a sandbox.

This happens when your workspace has a `.claude.json` file with a `primaryApiKey` field. Choose one of these approaches:

- Remove the `primaryApiKey` field from your `.claude.json`:

  ```json
  {
    "apiKeyHelper": "/path/to/script",
    "env": {
      "ANTHROPIC_BASE_URL": "https://api.anthropic.com"
    }
  }
  ```

- Or proceed with the warning - workspace credentials will be ignored in favor of sandbox credentials.

## Permission denied when accessing workspace files

Claude or commands fail with "Permission denied" errors when accessing files in the workspace.

This usually means the workspace path isn't accessible to Docker, or file permissions are too restrictive.

If using Docker Desktop:

1. Check File Sharing settings at Docker Desktop → **Settings** → **Resources** → **File Sharing**.

2. Ensure your workspace path (or a parent directory) is listed under Virtual file shares.

3. If missing, click "+" to add the directory containing your workspace.

4. Restart Docker Desktop.

For all platforms, verify file permissions:

```console
$ ls -la <workspace>
```

Ensure files are readable. If needed:

```console
$ chmod -R u+r <workspace>
```

Also verify the workspace path exists:

```console
$ cd <workspace>
$ pwd
```

## Sandbox crashes on Windows when launching multiple sandboxes

On Windows, launching too many sandboxes simultaneously can cause crashes.

If this happens, recover by closing the OpenVMM processes:

1. Open Task Manager (Ctrl+Shift+Esc).
2. Find all `docker.openvmm.exe` processes.
3. End each process.
4. Restart Docker Desktop if needed.

To avoid this issue, launch sandboxes one at a time rather than creating
multiple sandboxes concurrently.

## Persistent issues or corrupted state

If sandboxes behave unexpectedly or fail to start, reset all sandbox state:

```console
$ docker sandbox reset
```

This stops all running VMs and deletes all sandbox data. The daemon continues
running. After reset, create fresh sandboxes as needed.

Use reset when troubleshooting persistent problems or to reclaim disk space from
all sandboxes at once.
