---
title: Measure random reads and writes in a shared folder
linkTitle: Random I/O, shared folder
description: Measure random read and write throughput on a folder shared from your machine with stress-ng.
keywords: docker desktop, performance, benchmark, stress-ng, file sharing, virtiofs, bind mount, io
weight: 100
---

A build or test run that works on your source tree reads and writes across its
files at once. This check puts that load on a folder shared from your machine.
Two workers each write 2 GB at random locations and read it back at random
locations for 5 seconds, using `stress-ng`.

The result is a rate, not a duration. `stress-ng` counts the operations it
completes in 5 seconds. Higher is better. Read the
`bogo-ops-per-second-real-time` value.

Use `stress-ng` 0.21.02, the version Alpine 3.24 installs. Other versions count
the same work differently, so their numbers aren't comparable. Installing the
tool inside the container doesn't affect the result, because `stress-ng`
measures only its own 5 seconds.

For how to repeat a check, see [Performance checks](/manuals/desktop/features/vmm/performance-checks/).

> [!NOTE]
> This check fails on Docker Desktop with gRPC-FUSE file sharing, because that
> implementation doesn't let a program look at a file it already deleted. Use
> VirtioFS.

## Run the check

1. Start from an empty shared folder:

   ```console
   $ mkdir -p ~/perfcheck/stress-shared && cd ~/perfcheck/stress-shared
   $ rm -rf shared && mkdir shared
   ```

2. Run the measurement and print the rate:

   ```console
   $ docker run --rm -v "$PWD/shared:/bench" -w /bench alpine:3.24.1 sh -c '
     apk add -q --no-cache stress-ng && stress-ng --hdd 2 --hdd-bytes 2G --hdd-opts wr-rnd,rd-rnd --timeout 5s --metrics-brief --yaml /tmp/metrics.yaml &&
     grep -E "stressor:|bogo-ops-per-second-real-time:" /tmp/metrics.yaml
   '
   ```

   The output has this shape, with your own number:

   ```text
       - stressor: hdd
         bogo-ops-per-second-real-time: 64821.661549
   ```

3. Check that the run is valid. The command must exit without an error and
   print exactly one `stressor: hdd` line.

4. Clean up:

   ```console
   $ rm -rf shared
   ```

## Results that aren't valid

- Any `stress-ng` other than 0.21.02.
- Docker Desktop with gRPC-FUSE file sharing. This check fails there, because that file sharing does not let a program look at a file it has just deleted. Use VirtioFS.