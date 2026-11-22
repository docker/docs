---
title: Troubleshoot Docker Offload
linktitle: Troubleshoot
weight: 800
description: Learn how to troubleshoot issues with Docker Offload.
keywords: cloud, troubleshooting, Docker Desktop, offload, entitlement
tags: [Troubleshooting]
---

{{< summary-bar feature_name="Docker Offload" >}}

Docker Offload requires:

- Docker Desktop 4.97 or later
- Authentication
- An active internet connection
- No restrictive proxy or firewall blocking traffic to Docker Cloud
- A Docker Agentic Platform subscription
- **Enable Docker Offload** turned off, if the setting appears in Docker
  Desktop settings

## Docker Offload commands report an error

Because there are two versions of Docker Offload and only one can be active at
a time, most access errors come from a mismatch between the **Enable Docker
Offload** setting and your subscription. Use the following table to map the
error to a fix.

| You see | Cause | Fix |
| --- | --- | --- |
| A message telling you to turn off **Enable Docker Offload** when you run `docker offload create` or `docker --cloud` | The setting is on, so your commands use Docker Offload v1 | Turn off **Enable Docker Offload** in **Settings** > **Docker Offload**, then confirm your account has a Docker Agentic Platform subscription |
| No Docker Agentic Platform subscription, with a link to purchase | Your account doesn't have a subscription that covers Docker Offload | [Activate a subscription](/manuals/agentic-platform/signup.md#activate-cloud-access) |

## Builds or containers aren't running in the cloud

1. Check whether a sandbox is connected:

   ```console
   $ docker offload status
   ```

1. Check which Docker context is selected. Docker Offload uses the `docker-cloud`
   context:

   ```console
   $ docker context ls
   ```

   If the selected context is `desktop-linux`, your commands run locally. Run
   `docker offload create` to connect, or pass `--cloud` on individual
   commands.

1. Verify authentication with `docker login`. If needed, sign out and sign in
   again:

   ```console
   $ docker logout
   $ docker login
   ```

## A Docker command reports an unknown --cloud flag

```console
$ docker run --cloud --rm hello-world
unknown flag: --cloud
```

`--cloud` is a global flag on the `docker` command, so it must come before the
subcommand:

```console
$ docker --cloud run --rm hello-world
```

A subcommand doesn't define `--cloud` itself, so it rejects the flag wherever it
appears after the subcommand name, including between other options. This applies
to every subcommand.

## Your containers are missing

If `docker ps -a` doesn't show containers you expect, check which context is
selected. Containers in your cloud sandbox and containers on your local machine
are separate, and each context only shows its own.

To list your local containers while Docker Offload is connected:

```console
$ docker --context desktop-linux ps -a
```

For more information, see [Work with local and cloud
containers](sandboxes.md#work-with-local-and-cloud-containers).

## Your containers disappeared

If `docker ps` and `docker images` list nothing, check whether your sandbox is
paused:

```console
$ docker offload status
```

A paused sandbox keeps its containers, images, and volumes, but doesn't list
them. Neither `docker ps` nor `docker images` resumes a sandbox. Run a command
that needs the sandbox, such as `docker start` or `docker logs`, and your
containers appear again. For more information, see [While your sandbox is
paused](sandboxes.md#while-your-sandbox-is-paused).

If the sandbox isn't paused, check which context is selected. See [Your
containers are missing](#your-containers-are-missing).

## Your sandbox paused while you were still working

A sandbox pauses after 5 minutes without activity by default. A job that only
waits on a remote call reads as idle, and so do long gaps between commands.

An attached session or a forwarded port only counts as activity while data is
passing through it, so an idle terminal attached to a container doesn't keep
the sandbox awake.

Raise the idle timeout, or turn idle detection off for sessions you don't want
interrupted:

```console
$ docker offload create --idle-timeout 30m
$ docker offload create --keep-alive
```

Each machine remembers the idle timeout you last used, so you don't need to
pass the flag every time. With `--keep-alive`, the sandbox runs, and bills,
until you stop it. For more information, see [Control when your sandbox
pauses](sandboxes.md#control-when-your-cloud-sandbox-pauses).

If `docker offload create` reports that your settings weren't applied, an
`idle_timeout` in a config override file takes priority over both flags.

## Your attached or exec session ended

A pause ends attached sessions and exec sessions. After the sandbox resumes,
reattach with `docker attach`, or start a new `docker exec`.

## Your sandbox stopped while your machine was asleep

Idle detection runs on the machine connected to the sandbox, so it doesn't
apply when that machine sleeps or goes offline. The sandbox's time to live
applies instead. Docker Offload renews the time to live while a machine is
connected, so it only counts down once no machine is.

To give a sandbox longer to survive a disconnect, set a longer time to live
when you create it:

```console
$ docker offload create --ttl 8h
```

## Two machines report different times until idle

This is expected. Each machine connected to a sandbox tracks its own session
and its own idle detection, and they don't share that state. Both machines
report the same account and sandbox ID. For more information, see [Connect from
more than one machine](sandboxes.md#connect-from-more-than-one-machine).

## Resuming a sandbox is slow

Pausing takes several minutes to complete. If you run a command that needs the
sandbox while it's still pausing, Docker Offload waits for the pause to finish
and then resumes the sandbox. If the pause takes longer than that wait, Docker
Offload tells you the sandbox is still pausing. Try the command again shortly.

Resuming a sandbox that has finished pausing takes under a minute.
