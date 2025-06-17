---
title: Secure Docker Desktop
description: Use Settings Management to control Docker Desktop settings and learn best practices for hardened environments.
weight: 40
---

In secure environments, Docker Desktop must be tightly controlled to prevent
unauthorized behavior, unvetted updates, and data leakage. This module shows
you how to enforce strict settings using `admin-settings.json` and the Admin
Console, based on centralized security policy.

For a full list of supported settings and values, see the
[Settings reference](https://docs.docker.com/security/for-admins/hardened-desktop/settings-management/settings-reference/).

## Prerequisites

- Docker Desktop installed on user machines
- Organization owner access to your Docker organization
- A way to distribute `admin-settings.json` file or manage settings via MDM

## Step one: Disable telemetry and crash reporting

To prevent unapproved outbound data from Docker Desktop:

```json
{
  "sendUsageStatistics": false,
  "sendErrorReports": false
}
```

Lock both settings to ensure users can’t re-enable them. This prevents Docker
Desktop from sending anonymized usage data or crash logs to Docker servers.

## Step two: Disable update checks

Prevent Docker Desktop from checking for new versions:

```json
{
  "disableUpdates": true
}
```

You should distribute updates internally through a vetted package source. Lock
this setting to ensure consistent versions across your organization.

## Step three: Restrict Docker Desktop features

Set the following to prevent installation of unvetted code or execution of
risky commands:

```json
{
  "allowExtensions": false,
  "blockDockerLoad": true,
  "hideOnboarding": true,
  "exposeDockerAPIOnWindows": false
}
```

These settings disable extensions, block loading local image tarballs,
suppress UI prompts, and prevent unsafe remote access to the Docker daemon.

## Step four: Lock down file sharing and emulation

Limit which host paths can be mounted into containers, enforce safe drivers,
and restrict CPU architecture emulation:

```json

{
  "allowedFileSharingPaths": ["/Users/dev/projects", "/Volumes/code"],
  "shareFilesOnStart": true,
  "useVirtioFS": true,
  "useGrpcFUSE": false,
  "useRosetta": false
  }
```

Ensure all of the above are locked. This ensures users can only mount secure,
pre-approved directories.

## Step five: Enforce SBOM indexing with Docker Scout

Enable Software Bill of Materials (SBOM) generation for all images:

```json
{
  "enableSBOMIndexing": true,
  "enableBackgroundSBOMIndexing": true
  }
```

Lock both settings to support compliance and vulnerability management workflows.

## Step six: Enforce proxy configuration

Prevent users from overriding corporate proxy settings:

```json
{
  "proxySettings": {
    "mode": "manual",
    "httpProxy": "http://proxy.corp.example.com:3128",
    "httpsProxy": "http://proxy.corp.example.com:3128",
    "noProxy": "*.internal.example.com,localhost",
    "authenticationMethod": "kerberos"
	  },
  "allowUserToEditProxySettings": false
  }
```

Lock the entire proxy config to ensure consistent and auditable network routing.

## Step seven: Standardize the Linux VM and WSL settings

Enforce low-level Linux backend configuration:

```json
{
  "useWSL2Engine": true,
  "daemonConfigFile": {
    "log-driver": "json-file",
    "storage-driver": "overlay2"
  },
  "vpnKitCIDR": "10.255.0.0/24"
}
```

Lock all of these settings to prevent drift from your approved baseline.

## Step eight: Disable embedded Kubernetes

To prevent local Kubernetes environments from drifting from production
standards:

```json
{
  "allowKubernetes": false,
  "showKubernetesSystemContainers": false,
  "kubernetesImageRepository": "registry.corp.example.com/k8s"
}
```

Lock these settings to remove ambiguity around where clusters are running and
how they behave.

## Step nine: Enforce Enhanced Container Isolation (ECI)

ECI provides a hardened boundary around Docker containers, especially important
in untrusted environments.

```json
{
  "enableECI": true,
  "allowECIConfiguration": false,
  "eciAllowedImages": ["registry.corp.example.com/secure-build-runner:latest"],
  "eciAllowDerivedImages": false,
  "eciAllowedCommands": ["ps", "pull"]
}
```

Use ECI to restrict which containers can mount the Docker socket and exactly
what they can do with it.

## Best practices

- Lock all critical security settings to prevent end-user overrides.
- Distribute `admin-settings.json` centrally, ideally via MDM or login script.
- Use [desktop settings reporting](https://docs.docker.com/security/for-admins/hardened-desktop/settings-management/compliance-reporting/) to audit compliance.
- Enforce a single approved Desktop version across your org. See [disableUpdates](https://docs.docker.com/security/for-admins/hardened-desktop/settings-management/settings-reference/#disableupdates).
- Avoid optional features unless explicitly vetted for your use case.
