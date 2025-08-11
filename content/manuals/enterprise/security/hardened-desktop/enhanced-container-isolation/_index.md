---
title: Enhanced Container Isolation
linkTitle: Enhanced Container Isolation
description: Enhanced Container Isolation provides additional security for Docker Desktop by preventing malicious containers from compromising the host
keywords: nhanced container isolation, container security, sysbox runtime, linux user namespaces, hardened desktop
aliases:
 - /desktop/hardened-desktop/enhanced-container-isolation/
 - /security/for-admins/hardened-desktop/enhanced-container-isolation/
weight: 20
---

{{< summary-bar feature_name="Hardened Docker Desktop" >}}

Enhanced Container Isolation (ECI) adds an extra layer of security to prevent malicious containers from compromising Docker Desktop or the host system. It uses advanced isolation techniques while maintaining full developer productivity.

ECI strengthens container isolation and locks in security configurations created by administrators, such as [Registry Access Management policies](/manuals/enterprise/security/hardened-desktop/registry-access-management.md) and [Settings Management](../settings-management/_index.md) controls.

> [!NOTE]
>
> ECI works alongside other Docker security features like reduced Linux capabilities, seccomp, and AppArmor.

## Who should use Enhanced Container Isolation?

Enhanced Container Isolation is designed for:

- Organizations that want to prevent container-based attacks and reduce security vulnerabilities in developer environments
- Security teams that need stronger container isolation without impacting developer workflows
- Enterprises that require additional protection when running untrusted or third-party container images

## How Enhanced Container Isolation works

For an overview of how ECI works, see [How Enhanced Container Isolation works](/manuals/enterprise/security/hardened-desktop/enhanced-container-isolation/how-eci-works.md).

## Enable Enhanced Container Isolation

### For developers

To turn on ECI:

1. Verify your organization has a Docker Business subscription.
1. Sign in to your organization in Docker Desktop to access ECI features.
1. Stop and remove all existing containers:

  ```console
  $ docker stop $(docker ps -q)
  $ docker rm $(docker ps -aq)
  ```

1. In Docker Desktop, go to **Settings** > **General**.
1. Select the **Use Enhanced Container Isolation** checkbox.
1. Select **Apply and restart**.

> [!IMPORTANT]
>
> ECI doesn't protect containers created before turning on the feature. Remove existing containers before turning on ECI.

### For administrators

Before you begin, [enforce sign-in](/manuals/enterprise/security/enforce-sign-in) so users authenticate with your organization
when signing into Docker Desktop.

You can configure Enhanced Container Isolation for your organization using
Settings Management.

{{< tabs >}}
{{< tab name="Admin Console" >}}

1. Sign in to [Docker Home](https:app.docker.com) and select your organization.
1. Go to **Admin Console** > **Desktop Settings Management**.
1. [Create or edit a settings policy](/manuals/enterprise/security/hardened-desktop/settings-management/configure-admin-console.md).
1. Set **Enhanced Container Isolation** to **Always enabled**.

{{< /tab >}}
{{< tab name="`admin-settings.json` file" >}}

Create an [`admin-settings.json`](/manuals/enterprise/security/hardened-desktop/settings-management/configure-json-file.md) file with:

```json
{
  "configurationFileVersion": 2,
  "enhancedContainerIsolation": {
    "value": true,
    "locked": true
  }
}
```

Configuration options:
- `"value": true`: Turns on ECI by default
- `"locked": true`: Prevents developers from turning off ECI
- `"locked": false`: Allows developers to control the setting

{{< /tab >}}
{{< /tabs >}}

#### Apply the configuration

To apply the updated Enhanced Container Isolation setting for your organization:

- New installations: Users launch Docker Desktop and sign in
- Existing installations: Users must fully quit Docker Desktop and relaunch

> [!IMPORTANT]
>
> Restarting from the Docker Desktop menu isn't sufficient. Users must completely quit and reopen Docker Desktop.

You can also configure Docker socket mount permissions for trusted images that need Docker API access.

## Verify Enhanced Container Isolation is active

When ECI is turned on, users can verify it's working by checking the user
namespace mapping:

```console
$ docker run --rm alpine cat /proc/self/uid_map
```

When ECI is turned on:

```text
0     100000      65536
```

This shows the container's root user (0) maps to an unprivileged user (100000) in the Docker Desktop VM, with a range of 64K user IDs. Each container gets an exclusive user ID range for isolation.

When ECI is turned off:

```text
0          0 4294967295
```

This shows the container root user (0) maps directly to the VM root user (0), providing less isolation.

You can also check the container runtime:

```console
$ docker inspect --format='{{.HostConfig.Runtime}}' <container_name>
```

With ECI, it returns `sysbox-runc`. Without ECI it returns `runc`.

## What developers see with enforced ECI

When administrators enforce Enhanced Container Isolation:

- The Use Enhanced Container Isolation setting appears turned on in Docker Desktop settings
- The setting is locked and can't be changed if `"locked": true` was configured
- All new containers automatically use Linux user namespaces
- Existing development workflows continue to work without modification
