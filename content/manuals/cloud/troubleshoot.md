---
title: Troubleshoot Docker Cloud
linktitle: Troubleshoot
weight: 999
description: Learn how to troubleshoot issues with Docker Cloud.
keywords: cloud, troubleshooting, cloud mode, Docker Desktop, cloud builder, usage
tags: [Troubleshooting]
---


If you're having trouble using Docker Cloud, this guide can help you troubleshoot:

- Issues with cloud sessions and remote container runs using Docker Desktop (Cloud mode)
- Problems running cloud builds without enabling Cloud mode

## Troubleshoot Cloud mode

Docker Desktop uses Cloud mode to run both builds and containers in the cloud.
If builds or containers are failing to run, falling back to local, or reporting
session errors, use the following commands.

### 1. Check your Docker Cloud session

Docker Cloud requires an active session between Docker Desktop and the cloud.
Use the following command to check if the connection is active:

```console
$ docker cloud status
```

If you're not connected, start a new session:

```console
$ docker cloud start
```

To stop a session manually:

```console
$ docker cloud stop
```

If you're not sure what's wrong, try running diagnose to get more information
about your environment:

```console
$ docker cloud diagnose
```

This will print helpful information about your environment.

### 2. Validate cloud builder readiness

Cloud mode automatically provisions and uses a remote builder. You can confirm
the builder is initialized and ready using:


```console
$ docker cloud diagnose
```

In the output, look for a healthy session and builder status. If the builder
isn't available, try restarting the session:

```console
$ docker cloud stop
$ docker cloud start
```


### 3. Check for authentication or networking issues

Cloud mode requires:

- A valid Docker login
- An active internet connection
- No restrictive proxy or firewall blocking traffic to Docker Cloud

Use the following to verify login status:

```console
$ docker login
```

If needed, you can log out and then log in again:

```console
$ docker logout
$ docker login
```

## Troubleshoot builds only

If you're using Docker Cloud for builds only (without enabling Cloud mode in
Docker Desktop), issues may include builds running locally instead of in the
cloud, builder errors, or authentication failures.

### 1. Check your active Docker context

Cloud builds run using a remote builder, but your Docker context may still affect behavior.

List all contexts:

```console
$ docker context ls
```

You can inspect a specific context:

```console
$ docker context inspect <context-name>
```

Switch to the appropriate context if needed:

```console
$ docker context use <context-name>
```

### 2. Verify your builder configuration

Check available builders:

```console
$ docker buildx ls
```

Look for a builder with the `cloud` driver, for example `cloud-myorg-myteam`. The `*` marks the current one.

If you're not using the cloud builder, select it:

```console
$ docker buildx use cloud-myorg-myteam
```

### 3. Confirm builder readiness

Ensure the builder is initialized and ready:

```console
$ docker buildx inspect --bootstrap
```
Look for `Status: running` and a cloud driver in the output.

### 4. Check Docker login

Make sure you're authenticated with Docker:

```console
$ docker login
```

Re-authenticate if needed:

```console
$ docker logout
$ docker login
```