---
redirect_from:
  - "/engine/articles/systemd/"
title: "Limit a container's resources"
description: "Limiting the system resources a container can use"
keywords: "docker, daemon, configuration"
---

By default, a container has no resource constraints and can use as much of a
given resource as the host's kernel scheduler will allow. Docker provides ways
to control how much memory, CPU, or block IO a container can use, setting runtime
configuration flags of the `docker run` command. This section provides details
on when you should set such limits and the possible implications of setting them.

## Memory

Docker can enforce hard memory limits, which allow the container to use no more
than a given amount of user or system memory, or soft limits, which allow the
container to use as much memory as it needs unless certain conditions are met,
such as when the kernel detects low memory or contention on the host machine.
Some of these options have different effects when used alone or when more than
one option is set.

Most of these options take a positive integer, followed by a suffix of `b`, `k`,
`m`, `g`, to indicate bytes, kilobytes, megabytes, or gigabytes.

| Option                | Description                 |
|-----------------------|-----------------------------|
| `-m` or `--memory=` | The maximum amount of memory the container can use. If you set this option, the minimum allowed value is `4m` (4 megabyte). |
| `--memory-swap`*    | The amount of memory this container is allowed to swap to disk. See [`--memory-swap` details](resource_constraints.md#memory-swap-details). |
| `--memory-swappiness` | By default, the host kernel can swap out a percentage of anonymous pages used by a container. You can set `--memory-swappiness` to a value between 0 and 100, to tune this percentage. See [`--memory-swappiness` details](resource_constraints.md#memory-swappiness-details). |
| `--memory-reservation` | Allows you to specify a soft limit smaller than `--memory` which is activated when Docker detects contention or low memory on the host machine. If you use `--memory-reservation`, it must be set lower than `--memory` in order for it to take precedence. Because it is a soft limit, it does not guarantee that the container will not exceed the limit. |
| `--kernel-memory` | The maximum amount of kernel memory the container can use. The minimum allowed value is `4m`. Because kernel memory cannot be swapped out, a container which is starved of kernel memory may block host machine resources, which can have side effects on the host machine and on other containers. See [`--kernel-memory` details](resource_constraints.md#kernel-memory-details). |
| `--oom-kill-disable` | By default, if an out-of-memory (OOM) error occurs, the kernel kills processes in a container. To change this behavior, use the `--oom-kill-disable` option. Only disable the OOM killer on containers where you have also set the `-m/--memory` option. If the `-m` flag is not set, the host can run out of memory and the kernel may need to kill the host system's processes to free memory. |

For more information about cgroups and memory in general, see the documentation
for [Memory Resource Controller](https://www.kernel.org/doc/Documentation/cgroup-v1/memory.txt).

### `--memory-swap` details

- If unset, and `--memory` is set, the container can use twice as much swap
      as the `--memory` setting, if the host container has swap memory configured.
      For instance, if `--memory="300m"` and `--memory-swap` is not set, the
      container can use 300m of memory and 600m of swap.
- If set to a positive integer,  and if both `--memory` and `--memory-swap`
      are set, `--memory-swap` represents the total amount of memory and swap
      that can be used, and `--memory` controls the amount used by non-swap
      memory. So if `--memory="300m"` and `--memory-swap="1g"`, the container
      can use 300m of memory and 700m (1g - 300m) swap.
- If set to `-1` (the default), the container is allowed to use unlimited swap memory.

### `--memory-swappiness` details

- A value of 0 turns off anonymous page swapping.
- A value of 100 sets all anonymous pages as swappable.
- By default, if you do not set `--memory-swappiness`, the value is
  inherited from the host machine.

### `--kernel-memory` details

Kernel memory limits are expressed in terms of the overall memory allocated to
a container. Consider the following scenarios:

- **Unlimited memory, unlimited kernel memory**: This is the default
  behavior.
- **Unlimited memory, limited kernel memory**: This is appropriate when the
  amount of memory needed by all cgroups is greater than the amount of
  memory that actually exists on the host machine. You can configure the
  kernel memory to never go over what is available on the host machine,
  and containers which need more memory need to wait for it.
- **Limited memory, umlimited kernel memory**: The overall memory is
  limited, but the kernel memory is not.
- **Limited memory, limited kernel memory**: Limiting both user and kernel
  memory can be useful for debugging memory-related problems. If a container
  is using an unexpected amount of either type of memory, it will run out
  of memory without affecting other containers or the host machine. Within
  this setting, if the kernel memory limit is lower than the user memory
  limit, running out of kernel memory will cause the container to experience
  an OOM error. If the kernel memory limit is higher than the user memory
  limit, the kernel limit will not cause the container to experience an OOM.

When you turn on any kernel memory limits, the host machine tracks "high water
mark" statistics on a per-process basis, so you can track which processes (in
this case, containers) are using excess memory. This can be seen per process
by viewing `/proc/<PID>/status` on the host machine.

## CPU

By default, each container's access to the host machine's CPU cycles is unlimited.
You can set various constraints to limit a given container's access to the host
machine's CPU cycles.

| Option                | Description                 |
|-----------------------|-----------------------------|
| `--cpu-shares` | Set this flag to a value greater or less than the default of 1024 to increase or reduce the container's weight, and give it access to a greater or lesser proportion of the host machine's CPU cycles. This is only enforced when CPU cycles are constrained. When plenty of CPU cycles are available, all containers use as much CPU as they need. In that way, this is a soft limit. `--cpu-shares` does not prevent containers from being scheduled in swarm mode. It prioritizes container CPU resources for the available CPU cycles. It does not guarantee or reserve any specific CPU access. |
| `--cpu-period` | The scheduling period of one logical CPU on a container. `--cpu-period` defaults to a time value of 100000 (100 ms). |
| `--cpu-quota` | maximum amount of time that a container can be scheduled during the period set by `--cpu-period`. |
| `--cpuset-cpus` | Use this option to pin your container to one or more CPU cores, separated by commas. |

### Example with `--cpu-period` and `--cpu-qota`

If you have 1 vCPU system and your container runs with `--cpu-period=100000` and
`--cpu-quota=50000`, the container can consume up to 50% of 1 CPU.

```bash
$ docker run -ti --cpu-period=10000 --cpu-quota=50000 busybox
```

If you have a 4 vCPU system your container runs with `--cpu-period=100000` and
`--cpu-quota=200000`, your container can consume up to 2 logical CPUs (200% of
`--cpu-period`).

```bash
$ docker run -ti --cpu-period=100000 --cpu-quota=200000
```

### Example with `--cpuset-cpus`

To give a container access to exactly 4 CPUs, issue a command like the
following:

```bash
$ docker run -ti --cpuset-cpus=4 busybox
```

## Block IO (blkio)

Two option are available for tuning a given container's access to direct block IO
devices. You can also specify bandwidth limits in terms of bytes per second or
IO operations per second.

| Option                | Description                 |
|-----------------------|-----------------------------|
| `blkio-weight` | By default, each container can use the same proportion of block IO bandwidth (blkio). The default weight is 500. To raise or lower the proportion of blkio used by a given container, set the `--blkio-weight` flag to a value between 10 and 1000. This setting affects all block IO devices equally. |
| `blkio-weight-device` | The same as `--blkio-weight`, but you can set a weight per device, using the syntax `--blkio-weight-device="DEVICE_NAME:WEIGHT"` The DEVICE_NAME:WEIGHT is a string containing a colon-separated device name and weight. |
| `--device-read-bps` and `--device-write-bps` | Limits the read or write rate to or from a device by size, using a suffix of `kb`, `mb`, or `gb`. |
| `--device-read-iops` or `--device-write-iops` | Limits the read or write rate to or from a device by IO operations per second. |


### Block IO weight examples

>**Note**: The `--blkio-weight` flag only affects direct IO and has no effect on
buffered IO.

If you specify both the `--blkio-weight` and `--blkio-weight-device`, Docker
uses `--blkio-weight` as the default weight and uses `--blkio-weight-device` to
override the default on the named device.

To set a container's device weight for `/dev/sda` to 200 and not specify a
default `blkio-weight`:

```bash
$ docker run -it \
  --blkio-weight-device "/dev/sda:200" \
  ubuntu
```

### Block bandwidth limit examples

This example limits the `ubuntu` container to a maximum write speed of 1mbps to
`/dev/sda`:

```bash
$ docker run -it --device-write-bps /dev/sda:1mb ubuntu
```

This example limits the `ubuntu` container to a maximum read rate of 1000 IO
operations per second from `/dev/sda`:

```bash
$ docker run -ti --device-read-iops /dev/sda:1000 ubuntu
```
