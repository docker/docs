---
linkTitle: Profiling builds
title: Profiling builds with pprof
description: |
  Learn how you can use `pprof` to analyze and debug your builds.
---

{{< introduced buildx 0.18.0 >}}

You can configure Buildx to generate [`pprof`](https://github.com/google/pprof)
memory and CPU profiles to analyze and optimize your builds. These profiles
help you identify performance bottlenecks, detect memory inefficiencies, and
ensure your builds run efficiently.

`pprof` is a profiling tool that captures and visualizes detailed data on CPU
and memory usage in Go programs. It helps solve common performance problems by
allowing developers to understand which parts of their application are
consuming the most resources, whether CPU or memory.

In the context of Buildx and BuildKit, CPU profiles show where processing time
is spent during the build, helping to optimize slow builds. Memory profiles
track where memory is allocated, allowing you to spot inefficiencies, memory
leaks, or areas that require optimization. By generating these profiles, you
can focus on making your builds faster and more resource-efficient.

## Generate profiling data

The following environment variables control whether Buildx generates profiling
data for builds:

- [`BUILDX_CPU_PROFILE`](/manuals/build/building/variables.md#buildxcpuprofile)
- [`BUILDX_MEM_PROFILE`](/manuals/build/building/variables.md#buildxmemprofile)

When set, Buildx emits profiling samples for the builds to the location
specified by the environment variable.

## How to analyze profiling samples

To analyze and visualize profiling samples, you need `pprof` from the Go
toolchain. The following example shows how to run `pprof` in a container. If
you prefer to run `pprof` directly on your system, you need to install the Go
toolchain and (optionally) GraphViz for visualization.

1. Start a `golang` container named `pprof` in the background which publishes
   port 8081 (or any other available port) to the host.

   ```console
   $ docker run --rm --name pprof -w /profiles -p 8081:8081 -dt golang:alpine

2. Install GraphViz in the container.

   ```console
   $ docker exec pprof apk add --no-cache graphviz
   ```

3. Execute your build as usual with the desired [environment
   variables](#generate-profiling-data) set.

   ```console
   $ BUILDX_CPU_PROFILE=cpu.prof docker build .
   ```

4. Copy the profiling sample into the `pprof` container.

   ```console
   $ docker cp cpu.prof pprof:/profiles
   ```

5. Run `pprof` with the sample.

   ```console
   $ docker exec -it pprof go tool pprof cpu.prof
   ```

   This opens the `pprof` interactive console. From here, you can inspect the
   profiling sample using various commands. For example, use `top 10` command
   to view the top 10 most time-consuming entries.

   ```plaintext
   (pprof) top 10
   Showing nodes accounting for 3.04s, 91.02% of 3.34s total
   Dropped 123 nodes (cum <= 0.02s)
   Showing top 10 nodes out of 159
         flat  flat%   sum%        cum   cum%
        1.14s 34.13% 34.13%      1.14s 34.13%  syscall.syscall
        0.91s 27.25% 61.38%      0.91s 27.25%  runtime.kevent
        0.35s 10.48% 71.86%      0.35s 10.48%  runtime.pthread_cond_wait
        0.22s  6.59% 78.44%      0.22s  6.59%  runtime.pthread_cond_signal
        0.15s  4.49% 82.93%      0.15s  4.49%  runtime.usleep
        0.10s  2.99% 85.93%      0.10s  2.99%  runtime.memclrNoHeapPointers
        0.10s  2.99% 88.92%      0.10s  2.99%  runtime.memmove
        0.03s   0.9% 89.82%      0.03s   0.9%  runtime.madvise
        0.02s   0.6% 90.42%      0.02s   0.6%  runtime.(*mspan).typePointersOfUnchecked
        0.02s   0.6% 91.02%      0.02s   0.6%  runtime.pcvalue
   ```

6. To view the call graph in a graphical UI, run `go tool pprof
   -http=0.0.0.0:8081 <sample.prof>` in the container.

   ```console
   $ docker exec -it pprof
   /profiles # go tool pprof -http=0.0.0.0:8081 cpu.prof
   Serving web UI on http://0.0.0.0:8081
   http://0.0.0.0:8081
   ```

For more information about using `pprof` and how to interpret the call graph,
refer to the [`pprof` README](https://github.com/google/pprof/blob/main/doc/README.md).
