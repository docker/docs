---
description: Learn how to measure running containers, and about the different metrics
keywords: docker, metrics, CPU, memory, disk, IO, run, runtime, stats
title: Runtime metrics
weight: 50
aliases:
  - /articles/runmetrics/
  - /engine/articles/run_metrics/
  - /engine/articles/runmetrics/
  - /engine/admin/runmetrics/
  - /config/containers/runmetrics/
---

## Docker stats

You can use the `docker stats` command to live stream a container's
runtime metrics. The command supports CPU, memory usage, memory limit,
and network IO metrics.

The following is a sample output from the `docker stats` command

```console
$ docker stats redis1 redis2

CONTAINER           CPU %               MEM USAGE / LIMIT     MEM %               NET I/O             BLOCK I/O
redis1              0.07%               796 KB / 64 MB        1.21%               788 B / 648 B       3.568 MB / 512 KB
redis2              0.07%               2.746 MB / 64 MB      4.29%               1.266 KB / 648 B    12.4 MB / 0 B
```

The [`docker stats`](/reference/cli/docker/container/stats/) reference
page has more details about the `docker stats` command.

## Control groups

Linux Containers rely on [control groups](https://www.kernel.org/doc/Documentation/cgroup-v1/cgroups.txt)
which not only track groups of processes, but also expose metrics about
CPU, memory, and block I/O usage. You can access those metrics and
obtain network usage metrics as well. This is relevant for "pure" LXC
containers, as well as for Docker containers.

Control groups are exposed through a pseudo-filesystem. In modern distributions, you
should find this filesystem under `/sys/fs/cgroup`. Under that directory, you
see multiple sub-directories, called `devices`, `freezer`, `blkio`, and so on.
Each sub-directory actually corresponds to a different cgroup hierarchy.

On older systems, the control groups might be mounted on `/cgroup`, without
distinct hierarchies. In that case, instead of seeing the sub-directories,
you see a bunch of files in that directory, and possibly some directories
corresponding to existing containers.

To figure out where your control groups are mounted, you can run:

```console
$ grep cgroup /proc/mounts
```

### Enumerate cgroups

The file layout of cgroups is significantly different between v1 and v2.

If `/sys/fs/cgroup/cgroup.controllers` is present on your system, you are using v2,
otherwise you are using v1.
Refer to the subsection that corresponds to your cgroup version.

cgroup v2 is used by default on the following distributions:

- Fedora (since 31)
- Debian GNU/Linux (since 11)
- Ubuntu (since 21.10)

#### cgroup v1

You can look into `/proc/cgroups` to see the different control group subsystems
known to the system, the hierarchy they belong to, and how many groups they contain.

You can also look at `/proc/<pid>/cgroup` to see which control groups a process
belongs to. The control group is shown as a path relative to the root of
the hierarchy mountpoint. `/` means the process hasn't been assigned to a
group, while `/lxc/pumpkin` indicates that the process is a member of a
container named `pumpkin`.

#### cgroup v2

On cgroup v2 hosts, the content of `/proc/cgroups` isn't meaningful.
See `/sys/fs/cgroup/cgroup.controllers` to the available controllers.

### Changing cgroup version

Changing cgroup version requires rebooting the entire system.

On systemd-based systems, cgroup v2 can be enabled by adding `systemd.unified_cgroup_hierarchy=1`
to the kernel command line.
To revert the cgroup version to v1, you need to set `systemd.unified_cgroup_hierarchy=0` instead.

If `grubby` command is available on your system (e.g. on Fedora), the command line can be modified as follows:

```console
$ sudo grubby --update-kernel=ALL --args="systemd.unified_cgroup_hierarchy=1"
```

If `grubby` command isn't available, edit the `GRUB_CMDLINE_LINUX` line in `/etc/default/grub`
and run `sudo update-grub`.

### Running Docker on cgroup v2

Docker supports cgroup v2 since Docker 20.10.
Running Docker on cgroup v2 also requires the following conditions to be satisfied:

- containerd: v1.4 or later
- runc: v1.0.0-rc91 or later
- Kernel: v4.15 or later (v5.2 or later is recommended)

Note that the cgroup v2 mode behaves slightly different from the cgroup v1 mode:

- The default cgroup driver (`dockerd --exec-opt native.cgroupdriver`) is `systemd` on v2, `cgroupfs` on v1.
- The default cgroup namespace mode (`docker run --cgroupns`) is `private` on v2, `host` on v1.
- The `docker run` flag `--oom-kill-disable` is discarded on v2.

### Find the cgroup for a given container

For each container, one cgroup is created in each hierarchy. On
older systems with older versions of the LXC userland tools, the name of
the cgroup is the name of the container. With more recent versions
of the LXC tools, the cgroup is `lxc/<container_name>.`

For Docker containers using cgroups, the cgroup name is the full
ID or long ID of the container. If a container shows up as ae836c95b4c3
in `docker ps`, its long ID might be something like
`ae836c95b4c3c9e9179e0e91015512da89fdec91612f63cebae57df9a5444c79`. You can
look it up with `docker inspect` or `docker ps --no-trunc`.

Putting everything together to look at the memory metrics for a Docker
container, take a look at the following paths:

- `/sys/fs/cgroup/memory/docker/<longid>/` on cgroup v1, `cgroupfs` driver
- `/sys/fs/cgroup/memory/system.slice/docker-<longid>.scope/` on cgroup v1, `systemd` driver
- `/sys/fs/cgroup/docker/<longid>/` on cgroup v2, `cgroupfs` driver
- `/sys/fs/cgroup/system.slice/docker-<longid>.scope/` on cgroup v2, `systemd` driver

### Metrics from cgroups: memory, CPU, block I/O

For each subsystem (memory, CPU, and block I/O), one or
more pseudo-files exist and contain statistics. The metrics
available depend on your cgroup version. The following sections
describe both cgroup v1 and v2 formats.

#### Memory metrics

##### cgroup v1: `memory.stat`

Memory metrics are found in the `memory` cgroup. The memory
control group adds a little overhead, because it does very fine-grained
accounting of the memory usage on your host. Therefore, many distributions
chose to not enable it by default. Generally, to enable it, all you have
to do is to add some kernel command-line parameters:
`cgroup_enable=memory swapaccount=1`.

The metrics are in the pseudo-file `memory.stat`.
Here is what it looks like:

    cache 11492564992
    rss 1930993664
    mapped_file 306728960
    pgpgin 406632648
    pgpgout 403355412
    swap 0
    pgfault 728281223
    pgmajfault 1724
    inactive_anon 46608384
    active_anon 1884520448
    inactive_file 7003344896
    active_file 4489052160
    unevictable 32768
    hierarchical_memory_limit 9223372036854775807
    hierarchical_memsw_limit 9223372036854775807
    total_cache 11492564992
    total_rss 1930993664
    total_mapped_file 306728960
    total_pgpgin 406632648
    total_pgpgout 403355412
    total_swap 0
    total_pgfault 728281223
    total_pgmajfault 1724
    total_inactive_anon 46608384
    total_active_anon 1884520448
    total_inactive_file 7003344896
    total_active_file 4489052160
    total_unevictable 32768

The first half (without the `total_` prefix) contains statistics relevant
to the processes within the cgroup, excluding sub-cgroups. The second half
(with the `total_` prefix) includes sub-cgroups as well.

Some metrics are "gauges", or values that can increase or decrease. For instance,
`swap` is the amount of swap space used by the members of the cgroup.
Some others are "counters", or values that can only go up, because
they represent occurrences of a specific event. For instance, `pgfault`
indicates the number of page faults since the creation of the cgroup.

`cache`
: The amount of memory used by the processes of this control group that can be
associated precisely with a block on a block device. When you read from and
write to files on disk, this amount increases. This is the case if you use
"conventional" I/O (`open`, `read`, `write` syscalls) as well as mapped files
(with `mmap`). It also accounts for the memory used by `tmpfs` mounts, though
the reasons are unclear.

`rss`
: The amount of memory that doesn't correspond to anything on disk: stacks,
heaps, and anonymous memory maps.

`mapped_file`
: Indicates the amount of memory mapped by the processes in the control group.
It doesn't give you information about how much memory is used; it rather
tells you how it's used.

`pgfault`, `pgmajfault`
: Indicate the number of times that a process of the cgroup triggered a "page
fault" and a "major fault", respectively. A page fault happens when a process
accesses a virtual memory page that is not currently mapped to a physical
memory frame. This is a normal part of memory management. For example, a page
fault occurs when the process reads from a memory zone that has been swapped
out, or that corresponds to a memory-mapped file: in that case, the kernel
loads the page from disk and lets the CPU complete the memory access. It also
happens when the process writes to a copy-on-write memory zone: the kernel
duplicates the memory page and resumes the write operation on the process's
own copy of the page. "Major" faults happen when the kernel needs to read
data from disk. When it duplicates an existing page, or allocates an empty
page, it's a regular (or "minor") fault.

`swap`
: The amount of swap currently used by the processes in this cgroup.

`active_anon`, `inactive_anon`
: The amount of anonymous memory that has been identified has respectively
_active_ and _inactive_ by the kernel. "Anonymous" memory is the memory that is
_not_ linked to disk pages. In other words, that's the equivalent of the rss
counter described above. In fact, the very definition of the rss counter is
`active_anon` + `inactive_anon` - `tmpfs` (where tmpfs is the amount of
memory used up by `tmpfs` filesystems mounted by this control group). Now,
what's the difference between "active" and "inactive"? Pages are initially
"active"; and at regular intervals, the kernel sweeps over the memory, and tags
some pages as "inactive". Whenever they're accessed again, they're
immediately re-tagged "active". When the kernel is almost out of memory, and
time comes to swap out to disk, the kernel swaps "inactive" pages.

`active_file`, `inactive_file`
: Cache memory, with _active_ and _inactive_ similar to the _anon_ memory
above. The exact formula is `cache` = `active_file` + `inactive_file` +
`tmpfs`. The exact rules used by the kernel to move memory pages between
active and inactive sets are different from the ones used for anonymous memory,
but the general principle is the same. When the kernel needs to reclaim memory,
it's cheaper to reclaim a clean (=non modified) page from this pool, since it
can be reclaimed immediately (while anonymous pages and dirty/modified pages
need to be written to disk first).

`unevictable`
: The amount of memory that cannot be reclaimed; generally, it accounts for
memory that has been "locked" with `mlock`. It's often used by crypto
frameworks to make sure that secret keys and other sensitive material never
gets swapped out to disk.

`memory_limit`, `memsw_limit`
: These aren't really metrics, but a reminder of the limits applied to this
cgroup. The first one indicates the maximum amount of physical memory that can
be used by the processes of this control group; the second one indicates the
maximum amount of RAM+swap.

Accounting for memory in the page cache is very complex. If two
processes in different control groups both read the same file
(ultimately relying on the same blocks on disk), the corresponding
memory charge is split between the control groups. This behavior is useful, but it
also means that when a cgroup is terminated, it could increase the
memory usage of another cgroup, because they're not splitting the cost
anymore for those memory pages.

##### cgroup v2: `memory.stat`, `memory.current`, `memory.max`

On cgroup v2, the memory controller interface has been redesigned.
Instead of separate controllers for memory and swap, a single unified
`memory` controller manages both. Key interface files:

- **`memory.current`** — Current total memory usage (in bytes).
- **`memory.max`** — Memory limit (in bytes, or `max` for unlimited).
  This replaces the v1 `memory.limit_in_bytes` file.
- **`memory.high`** — Throttling threshold. Memory usage above this
  limit is throttled and put under heavy reclaim pressure. Setting this
  is useful for preventing sudden spikes before reaching the hard `memory.max` limit.
- **`memory.low`** — Best-effort memory protection. If memory usage
  is below this threshold, the cgroup's memory won't be reclaimed
  unless there's no unprotected reclaimable memory available.
- **`memory.min`** — Hard memory protection. Memory usage below this
  value is never reclaimed.
- **`memory.swap.current`** — Current swap usage (in bytes).
- **`memory.swap.max`** — Swap limit (in bytes, or `max` for unlimited).
- **`memory.events`** — Event counters including `low`, `high`, `max`,
  `oom`, and `oom_kill`. These are useful for monitoring memory pressure.

The `memory.stat` file contains detailed breakdowns. The following example shows part of the `memory.stat` output:

```text
anon 1639297024
file 2166460416
kernel 939536384
kernel_stack 4882432
pagetables 46481408
sock 65536
shmem 7659520
file_mapped 621125632
file_dirty 126976
swapcached 7847936
inactive_anon 1168683008
active_anon 459481088
inactive_file 1667866624
active_file 481972224
unevictable 32133120
slab_reclaimable 864874760
slab_unreclaimable 18258176
slab 883132936
workingset_refault_anon 7578
workingset_refault_file 99782
workingset_activate_anon 2203
workingset_activate_file 54321
pgfault 728281223
pgmajfault 1724
```

Notable differences from cgroup v1:

- `anon` replaces `rss` — anonymous memory (stacks, heaps, not backed by files).
- `file` replaces `cache` — file-backed memory (page cache).
- `kernel` — total kernel memory usage (not separately tracked in v1).
- `slab` — total slab memory (sum of `slab_reclaimable` + `slab_unreclaimable`).
- `shmem` — shared memory (previously included in `cache` in v1).
- `file_mapped` replaces `mapped_file`.
- `file_dirty` — file cache pages awaiting write-back to disk.
- There is no `total_` prefix. In cgroup v2, `memory.stat` always
  includes the cgroup's entire subtree, equivalent to the v1 `total_*` counters.
- `hierarchical_memory_limit` and `hierarchical_memsw_limit` are removed.
  Use `memory.max` and `memory.swap.max` instead.

`anon`
: Anonymous memory — the amount of memory not backed by files.
  This includes process stacks, heaps, and anonymous `mmap` regions.
  This is the cgroup v2 equivalent of the v1 `rss` counter.

`file`
: File-backed memory — the amount of memory used by the page cache,
  including files read from and written to disk. This is the cgroup v2
  equivalent of the v1 `cache` counter.

`kernel`
: Total kernel memory usage by the cgroup, including kernel stacks,
  page tables, and other kernel data structures. This was not separately
  trackable in cgroup v1.

`shmem`
: Shared memory usage, including `tmpfs` mounts and shared memory segments.
  In v1, this was included within `cache`.

`inactive_anon`, `active_anon`
: Anonymous memory that the kernel has classified as _inactive_ and _active_,
  respectively. The kernel uses this classification to decide which pages
  to swap out under memory pressure (inactive pages are reclaimed first).

`inactive_file`, `active_file`
: File-backed memory classified as _inactive_ and _active_.
  The formula `file` ≈ `inactive_file` + `active_file` approximately holds.
  Under memory pressure, inactive file pages are reclaimed before active ones.

`slab_reclaimable`, `slab_unreclaimable`
: Slab allocator memory that can (`reclaimable`) or cannot (`unreclaimable`)
  be freed under memory pressure.

`pgfault`, `pgmajfault`
: Page fault counters, same meaning as in cgroup v1.
  `pgfault` counts all page faults, `pgmajfault` counts those requiring disk I/O.

`unevictable`
: Memory that cannot be reclaimed (e.g., `mlock`ed pages).

Accounting for memory in the page cache is very complex. If two
processes in different control groups both read the same file
(ultimately relying on the same blocks on disk), the corresponding
memory charge is split between the control groups. This behavior is useful, but it
also means that when a cgroup is terminated, it could increase the
memory usage of another cgroup, because they're not splitting the cost
anymore for those memory pages.

#### CPU metrics

##### cgroup v1: `cpuacct.stat`

Now that we've covered memory metrics, everything else is
simple in comparison. CPU metrics are in the
`cpuacct` controller.

For each container, a pseudo-file `cpuacct.stat` contains the CPU usage
accumulated by the processes of the container, broken down into `user` and
`system` time. The distinction is:

- `user` time is the amount of time a process has direct control of the CPU,
  executing process code.
- `system` time is the time the kernel is executing system calls on behalf of
  the process.

Those times are expressed in ticks of 1/100th of a second, also called "user
jiffies". There are `USER_HZ` _"jiffies"_ per second, and on x86 systems,
`USER_HZ` is 100. Historically, this mapped exactly to the number of scheduler
"ticks" per second, but higher frequency scheduling and
[tickless kernels](https://lwn.net/Articles/549580/) have made the number of
ticks irrelevant.

##### cgroup v2: `cpu.stat`

On cgroup v2, the `cpu` and `cpuacct` controllers are unified into a single
`cpu` controller. CPU usage statistics are available in `cpu.stat`:

```text
usage_usec 9593743878
user_usec 7111219927
system_usec 2482523950
nr_periods 0
nr_throttled 0
throttled_usec 0
```

`usage_usec`
: Total CPU time consumed by the cgroup (in microseconds).
  This is the sum of `user_usec` and `system_usec`.

`user_usec`
: Time spent in user mode (in microseconds).
  Equivalent to the v1 `user` field, but in microseconds instead of jiffies.

`system_usec`
: Time spent in kernel mode (in microseconds).
  Equivalent to the v1 `system` field, but in microseconds instead of jiffies.

`nr_periods`, `nr_throttled`, `throttled_usec`
: CFS bandwidth throttling statistics. These are only meaningful if a
  CPU limit (`cpu.max`) is set.
  - `nr_periods`: number of enforcement periods elapsed.
  - `nr_throttled`: number of periods where the cgroup was throttled.
  - `throttled_usec`: total time the cgroup was throttled (in microseconds).

Additional interface files:

- **`cpu.max`** — CPU bandwidth limit, in the format `"<quota> <period>"`.
  For example, `"10000 100000"` means 10ms per 100ms period (10% of one CPU).
  Use `"max"` for no limit (the default).
- **`cpu.weight`** — Relative weight for CPU sharing (1–10000, default 100).
  This replaces the v1 `cpu.shares` interface.

#### Block I/O metrics

##### cgroup v1: `blkio` controller

Block I/O is accounted in the `blkio` controller.
Different metrics are scattered across different files. While you can
find in-depth details in the [blkio-controller](https://www.kernel.org/doc/Documentation/cgroup-v1/blkio-controller.txt)
file in the kernel documentation, here is a short list of the most
relevant ones:

`blkio.sectors`
: Contains the number of 512-bytes sectors read and written by the processes
member of the cgroup, device by device. Reads and writes are merged in a single
counter.

`blkio.io_service_bytes`
: Indicates the number of bytes read and written by the cgroup. It has 4
counters per device, because for each device, it differentiates between
synchronous vs. asynchronous I/O, and reads vs. writes.

`blkio.io_serviced`
: The number of I/O operations performed, regardless of their size. It also has
4 counters per device.

`blkio.io_queued`
: Indicates the number of I/O operations currently queued for this cgroup. In
other words, if the cgroup isn't doing any I/O, this is zero. The opposite is
not true. In other words, if there is no I/O queued, it doesn't mean that the
cgroup is idle (I/O-wise). It could be doing purely synchronous reads on an
otherwise quiescent device, which can therefore handle them immediately,
without queuing. Also, while it's helpful to figure out which cgroup is
putting stress on the I/O subsystem, keep in mind that it's a relative
quantity. Even if a process group doesn't perform more I/O, its queue size can
increase just because the device load increases because of other devices.

##### cgroup v2: `io.stat`

On cgroup v2, the `blkio` controller is replaced by the `io` controller.
All I/O statistics are consolidated into a single `io.stat` file, with
one line per device:

```text
8:0 rbytes=17408 wbytes=0 rios=14 wios=0 dbytes=0 dios=0
8:16 rbytes=8260728320 wbytes=167597064192 rios=472549 wios=16323808 dbytes=84085252096 dios=37547
```

Each line starts with the device's `major:minor` number, followed by key-value
pairs:

`rbytes`
: Number of bytes read from the device.
  Equivalent to the v1 `blkio.io_service_bytes` read counters.

`wbytes`
: Number of bytes written to the device.
  Equivalent to the v1 `blkio.io_service_bytes` write counters.

`rios`
: Number of read I/O operations.
  Equivalent to the v1 `blkio.io_serviced` read counter.

`wios`
: Number of write I/O operations.
  Equivalent to the v1 `blkio.io_serviced` write counter.

`dbytes`
: Number of bytes discarded (trim/unmap operations).

`dios`
: Number of discard I/O operations.

Additional interface files:

- **`io.max`** — Per-device I/O rate limits. The format is
  `"<major>:<minor> rbps=<limit> wbps=<limit> riops=<limit> wiops=<limit>"`.
  For example:

  ```text
  8:0 rbps=max wbps=10485760 riops=max wiops=100
  ```

Notable differences from cgroup v1:

- All I/O stats are in a single file instead of scattered across multiple
  `blkio.*` files.
- Per-device stats are identified by `major:minor` numbers rather than
  separate stat files per device.
- Byte values are in bytes (not sectors like `blkio.sectors`).
- The `io.max` interface replaces `blkio.throttle.*` with a unified
  per-device per-limit format.

### Network metrics

Network metrics aren't exposed directly by control groups. There is a
good explanation for that: network interfaces exist within the context
of _network namespaces_. The kernel could probably accumulate metrics
about packets and bytes sent and received by a group of processes, but
those metrics wouldn't be very useful. You want per-interface metrics
(because traffic happening on the local `lo`
interface doesn't really count). But since processes in a single cgroup
can belong to multiple network namespaces, those metrics would be harder
to interpret: multiple network namespaces means multiple `lo`
interfaces, potentially multiple `eth0`
interfaces, etc.; so this is why there is no easy way to gather network
metrics with control groups.

Instead you can gather network metrics from other sources.

#### iptables

iptables (or rather, the netfilter framework for which iptables is just
an interface) can do some serious accounting.

For instance, you can setup a rule to account for the outbound HTTP
traffic on a web server:

```console
$ iptables -I OUTPUT -p tcp --sport 80
```

There is no `-j` or `-g` flag,
so the rule just counts matched packets and goes to the following
rule.

Later, you can check the values of the counters, with:

```console
$ iptables -nxvL OUTPUT
```

Technically, `-n` isn't required, but it
prevents iptables from doing DNS reverse lookups, which are probably
useless in this scenario.

Counters include packets and bytes. If you want to setup metrics for
container traffic like this, you could execute a `for`
loop to add two `iptables` rules per
container IP address (one in each direction), in the `FORWARD`
chain. This only meters traffic going through the NAT
layer; you also need to add traffic going through the userland
proxy.

Then, you need to check those counters on a regular basis. If you
happen to use `collectd`, there is a [nice plugin](https://collectd.org/wiki/index.php/Table_of_Plugins)
to automate iptables counters collection.

#### Interface-level counters

Since each container has a virtual Ethernet interface, you might want to check
directly the TX and RX counters of this interface. Each container is associated
to a virtual Ethernet interface in your host, with a name like `vethKk8Zqi`.
Figuring out which interface corresponds to which container is, unfortunately,
difficult.

But for now, the best way is to check the metrics _from within the
containers_. To accomplish this, you can run an executable from the host
environment within the network namespace of a container using **ip-netns
magic**.

The `ip-netns exec` command allows you to execute any
program (present in the host system) within any network namespace
visible to the current process. This means that your host can
enter the network namespace of your containers, but your containers
can't access the host or other peer containers.
Containers can interact with their sub-containers, though.

The exact format of the command is:

```console
$ ip netns exec <nsname> <command...>
```

For example:

```console
$ ip netns exec mycontainer netstat -i
```

`ip netns` finds the `mycontainer` container by
using namespaces pseudo-files. Each process belongs to one network
namespace, one PID namespace, one `mnt` namespace,
etc., and those namespaces are materialized under
`/proc/<pid>/ns/`. For example, the network
namespace of PID 42 is materialized by the pseudo-file
`/proc/42/ns/net`.

When you run `ip netns exec mycontainer ...`, it
expects `/var/run/netns/mycontainer` to be one of
those pseudo-files. (Symlinks are accepted.)

In other words, to execute a command within the network namespace of a
container, we need to:

- Find out the PID of any process within the container that we want to investigate;
- Create a symlink from `/var/run/netns/<somename>` to `/proc/<thepid>/ns/net`
- Execute `ip netns exec <somename> ....`

Review [Enumerate Cgroups](#enumerate-cgroups) for how to find
the cgroup of an in-container process whose network usage you want to measure.
From there, you can examine the pseudo-file named
`tasks`, which contains all the PIDs in the
cgroup (and thus, in the container). Pick any one of the PIDs.

Putting everything together, if the "short ID" of a container is held in
the environment variable `$CID`, then you can do this:

```console
$ TASKS=/sys/fs/cgroup/devices/docker/$CID*/tasks
$ PID=$(head -n 1 $TASKS)
$ mkdir -p /var/run/netns
$ ln -sf /proc/$PID/ns/net /var/run/netns/$CID
$ ip netns exec $CID netstat -i
```

## Tips for high-performance metric collection

Running a new process each time you want to update metrics is
(relatively) expensive. If you want to collect metrics at high
resolutions, and/or over a large number of containers (think 1000
containers on a single host), you don't want to fork a new process each
time.

Here is how to collect metrics from a single process. You need to
write your metric collector in C (or any language that lets you do
low-level system calls). You need to use a special system call,
`setns()`, which lets the current process enter any
arbitrary namespace. It requires, however, an open file descriptor to
the namespace pseudo-file (remember: that's the pseudo-file in
`/proc/<pid>/ns/net`).

However, there is a catch: you must not keep this file descriptor open.
If you do, when the last process of the control group exits, the
namespace isn't destroyed, and its network resources (like the
virtual interface of the container) stays around forever (or until
you close that file descriptor).

The right approach would be to keep track of the first PID of each
container, and re-open the namespace pseudo-file each time.

## Collect metrics when a container exits

Sometimes, you don't care about real time metric collection, but when a
container exits, you want to know how much CPU, memory, etc. it has
used.

Docker makes this difficult because it relies on `lxc-start`, which carefully
cleans up after itself. It is usually easier to collect metrics at regular
intervals, and this is the way the `collectd` LXC plugin works.

But, if you'd still like to gather the stats when a container stops,
here is how:

For each container, start a collection process, and move it to the
control groups that you want to monitor by writing its PID to the tasks
file of the cgroup. The collection process should periodically re-read
the tasks file to check if it's the last process of the control group.
(If you also want to collect network statistics as explained in the
previous section, you should also move the process to the appropriate
network namespace.)

When the container exits, `lxc-start` attempts to
delete the control groups. It fails, since the control group is
still in use; but that's fine. Your process should now detect that it is
the only one remaining in the group. Now is the right time to collect
all the metrics you need!

Finally, your process should move itself back to the root control group,
and remove the container control group. To remove a control group, just
`rmdir` its directory. It's counter-intuitive to
`rmdir` a directory as it still contains files; but
remember that this is a pseudo-filesystem, so usual rules don't apply.
After the cleanup is done, the collection process can exit safely.
