---
description: Find known issues for Docker Desktop
keywords: mac, troubleshooting, known issues, Docker Desktop
title: Known issues
tags: [ Troubleshooting ]
weight: 20
aliases:
 - /desktop/troubleshoot/known-issues/
---

{{< tabs >}}
{{< tab name="For Mac with Intel chip" >}}
- The Mac Activity Monitor reports that Docker is using twice the amount of memory it's actually using. This is due to a [bug in macOS].(https://docs.google.com/document/d/17ZiQC1Tp9iH320K-uqVLyiJmk4DHJ3c4zgQetJiKYQM/edit?usp=sharing) on this.

- **"Docker.app is damaged" dialog**: If you see a "Docker.app is damaged and can't be opened" dialog during installation or updates, this is typically caused by non-atomic copy operations when other applications are using the Docker CLI. See [Fix "Docker.app is damaged" on macOS](mac-damaged-dialog.md) for resolution steps.

- Force-ejecting the `.dmg` after running `Docker.app` from it can cause the
  whale icon to become unresponsive, Docker tasks to show as not responding in the Activity Monitor, and for some processes to consume a large amount of CPU resources. Reboot and restart Docker to resolve these issues.

- Docker Desktop uses the `HyperKit` hypervisor
  (https://github.com/docker/hyperkit) in macOS 10.10 Yosemite and higher. If
  you are developing with tools that have conflicts with `HyperKit`, such as
  [Intel Hardware Accelerated Execution Manager
  (HAXM)](https://software.intel.com/en-us/android/articles/intel-hardware-accelerated-execution-manager/),
  the current workaround is not to run them at the same time. You can pause
  `HyperKit` by quitting Docker Desktop temporarily while you work with HAXM.
  This allows you to continue work with the other tools and prevent `HyperKit`
  from interfering.

- If you are working with applications like [Apache
  Maven](https://maven.apache.org/) that expect settings for `DOCKER_HOST` and
  `DOCKER_CERT_PATH` environment variables, specify these to connect to Docker
  instances through Unix sockets. For example:

  ```console
  $ export DOCKER_HOST=unix:///var/run/docker.sock
  ```

{{< /tab >}}
{{< tab name="For Mac with Apple silicon" >}}

- Some command line tools do not work when Rosetta 2 is not installed.
  - The old version 1.x of `docker-compose`. Use Compose V2 instead - type `docker compose`.
  - The `docker-credential-ecr-login` credential helper.
- Some images do not support the ARM64 architecture. You can add `--platform linux/amd64` to run (or build) an Intel image using emulation.

   However, attempts to run Intel-based containers on Apple silicon machines under emulation can crash as QEMU sometimes fails to run the container. In addition, filesystem change notification APIs (`inotify`) do not work under QEMU emulation. Even when the containers do run correctly under emulation, they will be slower and use more memory than the native equivalent.

   In summary, running Intel-based containers on Arm-based machines should be regarded as "best effort" only. We recommend running `arm64` containers on Apple silicon machines whenever possible, and encouraging container authors to produce `arm64`, or multi-arch, versions of their containers. This issue should become less common over time, as more and more images are rebuilt [supporting multiple architectures](https://www.docker.com/blog/multi-arch-build-and-images-the-simple-way/).
- Users may occasionally experience data drop when a TCP stream is half-closed.

{{< /tab >}}
{{< /tabs >}}
