---
title: Docker Offload usage
linktitle: Usage
weight: 30
description: Learn how Docker Offload is billed and how to monitor your usage.
keywords: cloud, usage, offload, billing
---

{{< summary-bar feature_name="Docker Offload" >}}

Docker Offload is billed for the time your cloud sandbox runs, as part of your
Docker Agentic Platform subscription. For subscription and payment information,
see [Signup and billing](/manuals/agentic-platform/signup.md).

## What you're billed for

Billing starts when your sandbox finishes being created and stops when the
sandbox starts pausing, whether you run `docker offload stop` or the sandbox
pauses on its own after a period without activity.

You aren't billed while a sandbox is paused. A paused sandbox keeps its
containers, images, and volumes, so you can leave a workload in place between
sessions without paying for idle time.

By default, a sandbox pauses after 5 minutes without activity. You can change
that period with `--idle-timeout`.

> [!WARNING]
>
> `--keep-alive` turns off idle detection, so the sandbox keeps running, and
> billing continues, until you run `docker offload stop`. It still pauses if
> your machine sleeps past the sandbox's time to live, and after 24 hours. For
> more information, see [Control when your sandbox
> pauses](sandboxes.md#control-when-your-cloud-sandbox-pauses).

Some workloads keep a sandbox awake even when you've stopped using it, which
means it keeps billing. A container running on the host network does this, and
so does a large number of running containers. Run `docker offload status` to
see what's keeping a sandbox awake, and `docker offload stop` when you're
finished.

## Monitor your usage

To see your running sandbox and its hourly price, open the [Docker Agentic
Platform Console](https://agentic-platform.docker.com/).

### Find your sandbox in the Console

A sandbox that Docker Offload creates appears in the Console under your account
with a name that starts with `offload-`. Use that prefix to tell it apart from
sandboxes you started another way.

Your usage shows the sandbox's details, such as its size, vCPUs, and memory. It
doesn't report that Docker Offload created the sandbox, so the name prefix is
what identifies it.

## Reduce your usage

Because you're billed for the time your sandbox runs, the main ways to reduce
your spend are to stop your sandbox when you're not using it and to make your
builds faster. For build and transfer optimizations, see [Optimize Docker
Offload usage](optimize.md).
