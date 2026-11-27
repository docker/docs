---
title: Performance checks
linkTitle: Performance checks
description: Run the commands behind Docker's published performance numbers on your own machine and compare engines or settings.
keywords: docker desktop, performance, benchmark, vmm, docker vmm, file sharing, virtiofs, build, io
weight: 10
---

Every number in Docker's published performance results comes from a command you
can run yourself. Most numbers are the duration of a command. Two are a rate
reported by a test tool. These pages show you how to run each command and read
the same figure from your terminal.

Each check uses public repositories or public images pinned to an exact version,
so you time the same work as the published number. Each page is self-contained.
Pick the number you want to verify and follow that page.

## Checks

| Check | Result |
| ----- | ------ |
| [First `git status` in a shared repository](git-status.md) | Milliseconds, measured in the container |
| [First text search in a shared repository](code-search.md) | Milliseconds, measured in the container |
| [`git worktree add` in a shared folder](git-worktree.md) | Milliseconds, measured in the container |
| [Image build with a large `COPY`](build-copy-source.md) | Build time |
| [Image build with 20 steps](build-steps.md) | Milliseconds per step |
| [Rust compilation in an image build](build-rust.md) | Build time |
| [Multi-service Compose startup](compose-startup.md) | Time until a request succeeds |
| [Kubernetes install](kubernetes-startup.md) | Time until a checkout succeeds |
| [Small files in a shared folder](small-files.md) | Time for four containers |
| [Random reads and writes in a shared folder](random-io-shared.md) | Operations per second |
| [Random reads from a large file in a shared folder](random-read-shared.md) | Time for 25 000 reads |
| [Random reads from a large file in a container](random-read-container.md) | Time for 25 000 reads |
| [Random writes to a large file in a container](random-write-container.md) | Time for 25 000 writes |
| [Random reads and writes on container storage](random-io-container.md) | Operations per second |

## Prerequisites

All checks need:

- The engine you want to measure with the `docker` command working.
- `bash` or `zsh`. The timing relies on the shell's `time` keyword. On Windows,
  use Git Bash or run everything in WSL. Don't mix the two environments in one
  set of runs.
- Git v2.20 or later.
- 2-15 GB of free disk space, depending on the check.
- 10-60 minutes per check, mostly for the first download of images and source code.

Two checks need more:

- [Compose startup](compose-startup.md) needs Docker Compose 2.24.4 or later
  and `curl`. It also needs about 12 GB for images.
- [Kubernetes install](kubernetes-startup.md) needs a single-node local cluster
  from the engine you're measuring, such as Docker Desktop's Kubernetes, with
  `kubectl` pointing at it. It also needs `helm`, `curl`, and `jq`.

Work in a folder your engine can share. On Docker Desktop, that's a path inside
your home directory that's listed in **Settings** > **Resources** > **File
sharing**. On Mac, `/tmp` usually isn't shared. The instructions use
`~/perfcheck`, with one subfolder for each check.

For reliable results:

- Close builds, IDE indexing, backups, and antivirus scans.
- Connect laptops to power.
- Run one check at a time.

### Record your environment

A duration is only meaningful with the machine and settings it came from. Save
the output of these commands before you start:

```console
$ docker version
$ docker info
$ uname -a
```

On Docker Desktop, also record these settings:

- CPU and memory limits for the virtual machine
- Virtualization backend, such as Apple Virtualization, Docker VMM, WSL 2, or
  Hyper-V
- File sharing implementation, such as VirtioFS, gRPC-FUSE, or osxfs
- Whether the containerd image store is turned on
- Whether Rosetta or other emulation is turned on

Any one of these can change a result by more than a factor of two.

## How durations are taken

Checks use one of two methods, and each page says which.

### On your machine

The shell's `time` keyword measures the command:

```console
$ time docker build ...
real    0m22.964s
user    0m1.760s
sys     0m0.929s
```

The value to record is the elapsed time. In bash, that's `real`. In zsh, the same
command prints one line ending in `22.964 total`, and the elapsed time is the
number before `total`.

### Inside a container

Checks that time a command in a running container read the clock before and
after the command, inside the same container:

```console
$ docker exec some-container sh -c '
  start=$(date +%s%3N)
  ...the command...
  end=$(date +%s%3N)
  printf "ms=%s\n" "$((end - start))"
'
```

This keeps container startup and `docker exec` overhead out of the duration.
Alpine's built-in `date` can't print fractions of a second, so these checks
install GNU coreutils before timing starts.

### Rates

Two checks report a rate instead of a duration. `stress-ng` runs for 5 seconds
and reports operations per second. Higher is better.

## Repeat a check

If a result surprises you, run the check two or three times and compare the middle result. 
The first run of a session is often the slowest.

To repeat a check:

- Start from a fresh copy of the source or folder. These checks measure doing
  something for the first time. Re-running on the same files measures warm
  caches. Each page lists what to delete.
- Empty the engine's file cache between runs:

  ```console
  $ docker run --rm --privileged alpine:3.24 sh -c 'sync && echo 3 > /proc/sys/vm/drop_caches'
  ```

  This empties the cache in the engine's virtual machine, where the work
  happens, not on your computer. If your engine refuses `--privileged`, skip
  this step.

## Interpret results

Your durations differ from the published ones because they depend on your
hardware, engine, and settings. Measure two engines or two settings on one machine, and the difference should match the published difference.

- Treat differences under about 20% as noise. Internal testing showed one check drift from 235 ms to 180 ms over 20 minutes as the daemon and file cache settled.
- If the standard deviation of your durations is more than about 15% of their
  mean, something on the machine is interfering. Discard the runs.
