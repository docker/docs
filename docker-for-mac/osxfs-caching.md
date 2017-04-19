---
description: OSXFS CACHING
keywords: mac, osxfs
redirect_from:
- /mackit/osxfs-caching/
title: Behavior of the `cached` and `delegated` mount flags
---

Docker 17.04 CE Edge adds support for two new flags to the `-v`
option, `cached` and `delegated`, that can significantly improve
performance of mounted volume access in Docker for Mac.  A blog post,
[User-guided caching in Docker for Mac](link-TODO) describes the
background and motivation for the new flags.  This document gives a
more detailed description of their behavior.

## Background

With Docker distributions for an increasing number of platforms,
including macOS and Windows, generalizing mount semantics during
container run is a necessity to enable workload optimizations.

The current implementations of mounts on Linux provide a consistent
view of a host directory tree inside a container: reads and writes
performed either on the host or in the container are immediately
reflected in the other environment, and file system events (`inotify`,
`FSEvents`) are consistently propagated in both directions.

On Linux these guarantees carry no overhead, since the underlying VFS is
shared directly between host and container.  However, on macOS (and
other non-Linux platforms) there are significant overheads to
guaranteeing perfect consistency, since messages describing file system
actions must be passed synchronously between container and host.  The
current implementation is sufficiently efficient for most tasks, but
with certain types of workload the overhead of maintaining perfect
consistency can result in performance that is significantly worse than a
native (non-Docker) environment.  For example,

 * running `go list ./...` in the bind-mounted `docker/docker` source tree
   takes around 26 seconds

 * writing 100MB in 1k blocks into a bind-mounted directory takes
   around 23 seconds

 * running `ember build` on a freshly created (i.e. empty) application
   involves around 70000 sequential syscalls, each of which translates
   into a request and response passed between container and host.

Optimizations to reduce latency throughout the stack have brought
significant improvements to these workloads, and a few further
optimization opportunities remain.  However, even when latency is
minimized, the constraints of maintaining consistency mean that these
workloads remain unacceptably slow for some use cases.

**Fortunately, in many cases where the performance degradation is most
severe, perfect consistency between container and host is unnecessary.**
In particular, in many cases there is no need for writes performed in a
container to be immediately reflected on the host.  For example, while
interactive development requires that writes to a bind-mounted directory
on the host immediately generate file system events within a container,
there is no need for writes to build artifacts within the container to
be immediately reflected on the host file system.  Distinguishing between
these two cases makes it possible to significantly improve performance.

There are three broad scenarios to consider.  In each case the container
has an internally-consistent view of bind-mounted directories, but in
two cases temporary discrepancies are allowed between container and host.

 * `consistent`: perfect consistency  
   (host and container have an identical view of the mount at all times)

 * `cached`: the host's view is authoritative  
   (permit delays before updates on the host appear in the container)

 * `delegated`: the container's view is authoritative  
   (permit delays before updates on the container appear in the host)

## Semantics

The semantics of each configuration is described as a set of guarantees
relating to the observable effects of file system operations.  In this
specification, "host" refers to the file system of the user's Docker
client.

##### `delegated` Semantics

The `delegated` configuration provides the weakest set of guarantees.
For directories mounted with `delegated` the container's view of the
file system is authoritative, and writes performed by containers may not
be immediately reflected on the host file system.  As with (e.g.) NFS
asynchronous mode, if a running container with a `delegated` bind mount
crashes then writes may be lost.

However, by relinquishing consistency, `delegated` mounts can offer
significantly better performance than the other configurations.  Where
the data written is ephemeral or readily reproducible (e.g. scratch
space or build artifacts) `delegated` may be optimal for a user's
workload.

A `delegated` mount offers the following guarantees, which are presented
as constraints on the container run-time:

(1) If the implementation offers file system events, the container state
as it relates to a specific event MUST reflect the host file system
state at the time the event was generated if no container modifications
pertain to related file system state.

(2) If flush or sync operations are performed, relevant data MUST be
written back to the host file system.  Between flush or sync
operations containers MAY cache data written, metadata modifications,
and directory structure changes.

(3) All containers hosted by the same run-time MUST share a consistent
cache of the mount.

(4) When any container sharing a `delegated` mount terminates, changes
to the mount MUST be written back to the host file system. If this
writeback fails, the container's execution MUST fail via exit code
and/or Docker event channels.

(5) If a `delegated` mount is shared with a `cached` or a `consistent`
mount, those portions that overlap MUST obey `cached` or `consistent`
mount semantics respectively.

Besides these constraints, the `delegated` configuration offers the
container run-time a degree of flexibility:

(6) Containers MAY retain file data and metadata (including directory
structure, existence of nodes, etc) indefinitely and this cache MAY
desynchronize from the file system state of the host. Implementors are
encouraged to expire caches when host file system changes occur but,
due to platform limitations, may be unable to do this in any specific
time frame.

(7) If changes to the mount source directory are present on the host
file system, those changes MAY be lost when the `delegated` mount
synchronizes with the host source directory.

However,

(8) Behaviors 6-7 DO NOT apply to the file types of socket, pipe, or device.

##### `cached` Semantics

The `cached` configuration provides all the guarantees of the
`delegated` configuration and some additional guarantees around the
visibility of writes performed by containers.  For directories mounted
with `cached` the host's view of the file system is authoritative;
writes performed by containers are immediately visible to the host, but
there may be a delay before writes performed on the host are visible
within containers.

(1) Implementations MUST obey `delegated` Semantics 1-5.

Additionally,

(2) If the implementation offers file system events, the container state
as it relates to a specific event MUST reflect the host file system
state at the time the event was generated.

(3) Container mounts MUST perform metadata modifications, directory
structure changes, and data writes consistently with the host file
system, and MUST NOT cache data written, metadata modifications, or
directory structure changes.

(4) If a `cached` mount is shared with a `consistent` mount, those portions
that overlap MUST obey `consistent` mount semantics.

Some of the flexibility of the `delegated` configuration is retained,
namely:

(5) Implementations MAY permit `delegated` Semantics 6.

##### `consistent` Semantics

The `consistent` configuration places the most severe restrictions on
the container run-time.  For directories mounted with `consistent` the
container and host views are always synchronized: writes performed
within the container are immediately visible on the host, and writes
performed on the host are immediately visible within the container.

The `consistent` configuration most closely reflects the behavior of
bind mounts on Linux.  However, the overheads of providing strong
consistency guarantees make it unsuitable for a few use cases, where
performance is a priority and maintaining perfect consistency has low
priority.

(1) Implementations MUST obey `cached` Semantics 1-4.

Additionally,

(2) Container mounts MUST reflect metadata modifications, directory
structure changes, and data writes on the host file system immediately.

##### `default` Semantics

The `default` configuration is identical to the `consistent`
configuration except for its name. Crucially, this means that `cached`
Semantics 4 and `delegated` Semantics 5 that require strengthening
overlapping directories do not apply to `default` mounts. This is the
default configuration if no `state` flags are supplied.
