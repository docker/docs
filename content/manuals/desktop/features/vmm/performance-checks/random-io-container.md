---
title: Measure random reads and writes on container storage
linkTitle: Random I/O, container
description: Measure random read and write throughput on a container's own storage with stress-ng.
keywords: docker desktop, performance, benchmark, stress-ng, container storage, io
weight: 140
---

This check puts the same load as [random reads and writes in a shared
folder](random-io-shared.md) on the container's own storage. Two workers each
write 2 GB at random locations and read it back at random locations for
5 seconds. It represents an application in a container that writes logs,
temporary files, and data while also reading them.

The result is a rate reported by `stress-ng`. Higher is better. Read the
`bogo-ops-per-second-real-time` value.

Use `stress-ng` 0.21.02, the version Alpine 3.24 installs. Other versions count
the same work differently, so their numbers aren't comparable.

For how to repeat a check, see [Performance checks](/manuals/desktop/features/vmm/performance-checks/).

## Run the check

1. Run the measurement and print the rate:

   ```console
   $ docker run --rm -w /bench alpine:3.24.1 sh -c '
     apk add -q --no-cache stress-ng && stress-ng --hdd 2 --hdd-bytes 2G --hdd-opts wr-rnd,rd-rnd --timeout 5s --metrics-brief --yaml /tmp/metrics.yaml &&
     grep -E "stressor:|bogo-ops-per-second-real-time:" /tmp/metrics.yaml
   '
   ```

   The output has this shape, with your own number:

   ```text
       - stressor: hdd
         bogo-ops-per-second-real-time: 64821.661549
   ```

2. Check that the run is valid. The command must exit without an error and
   print exactly one `stressor: hdd` line. There's nothing to clean up.
