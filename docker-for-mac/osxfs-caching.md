---
description: Osxfs caching
keywords: mac, osxfs, volumes
title: Performance tuning for volume mounts (shared filesystems)
---

[Docker 17.04 CE Edge](https://github.com/docker/docker.github.io/blob/v17.03/edge/index.md#docker-ce-edge-new-features) adds support
for two new flags to the [docker run `-v`, `--volume`](../engine/reference/run.md#volume-shared-filesystems)
option, `cached` and `delegated`, that can significantly improve the performance
of mounted volume access on Docker Desktop for Mac. These options begin to solve some of
the challenges discussed in 
[Performance issues, solutions, and roadmap](osxfs.md#performance-issues-solutions-and-roadmap).

> **Tip:** Release notes for Docker CE Edge 17.04 are
[here](https://github.com/moby/moby/releases/tag/v17.04.0-ce), and the
associated pull request for the additional `docker run -v` flags is
[here](https://github.com/moby/moby/pull/31047).

The following topics describe the challenges of bind-mounted volumes on `osxfs`,
and the caching options provided to optimize performance.

This blog post on [Docker on Mac
Performance](https://stories.amazee.io/docker-on-mac-performance-docker-machine-vs-docker-for-mac-4c64c0afdf99)
gives a nice, quick summary.

For information on how to configure these options in a Compose file, see
[Caching options for volume mounts](../compose/compose-file/index.md#caching-options-for-volume-mounts-docker-desktop-for-mac)
the Docker Compose topics.

## Performance implications of host-container file system consistency

With Docker distributions now available for an increasing number of
platforms, including macOS and Windows, generalizing mount semantics
during container run is a necessity to enable workload optimizations.

The current implementations of mounts on Linux provide a consistent
view of a host directory tree inside a container: reads and writes
performed either on the host or in the container are immediately
reflected in the other environment, and file system events (`inotify`,
`FSEvents`) are consistently propagated in both directions.

On Linux, these guarantees carry no overhead, since the underlying VFS is
shared directly between host and container. However, on macOS (and
other non-Linux platforms) there are significant overheads to
guaranteeing perfect consistency, since messages describing file system
actions must be passed synchronously between container and host. The
current implementation is sufficiently efficient for most tasks, but
with certain types of workloads the overhead of maintaining perfect
consistency can result in significantly worse performance than a
native (non-Docker) environment. For example,

 * running `go list ./...` in the bind-mounted `docker/docker` source tree
   takes around 26 seconds

 * writing 100MB in 1k blocks into a bind-mounted directory takes
   around 23 seconds

 * running `ember build` on a freshly created (empty) application
   involves around 70000 sequential syscalls, each of which translates
   into a request and response passed between container and host.

Optimizations to reduce latency throughout the stack have brought
significant improvements to these workloads, and a few further
optimization opportunities remain. However, even when latency is
minimized, the constraints of maintaining consistency mean that these
workloads remain unacceptably slow for some use cases.

## Tuning with consistent, cached, and delegated configurations

**_Fortunately, in many cases where the performance degradation is most
severe, perfect consistency between container and host is unnecessary._**
In particular, in many cases there is no need for writes performed in a
container to be immediately reflected on the host. For example, while
interactive development requires that writes to a bind-mounted directory
on the host immediately generate file system events within a container,
there is no need for writes to build artifacts within the container to
be immediately reflected on the host file system. Distinguishing between
these two cases makes it possible to significantly improve performance.

There are three broad scenarios to consider, based on which you can dial in the
level of consistency you need. In each case, the container has an
internally-consistent view of bind-mounted directories, but in two cases
temporary discrepancies are allowed between container and host.

 * `consistent`: perfect consistency
   (host and container have an identical view of the mount at all times)

 * `cached`: the host's view is authoritative
   (permit delays before updates on the host appear in the container)

 * `delegated`: the container's view is authoritative
   (permit delays before updates on the container appear in the host)

## Examples

Each of these configurations (`consistent`, `cached`, `delegated`) can be
specified as a suffix to the
[`-v`](../engine/reference/run.md#volume-shared-filesystems)
option of [`docker run`](../engine/reference/commandline/run.md). For
example, to bind-mount `/Users/yallop/project` in a container under the path
`/project`, you might run the following command:

```bash
docker run -v /Users/yallop/project:/project:cached alpine command
```

The caching configuration can be varied independently for each bind mount,
so you can mount each directory in a different mode:

```bash
docker run -v /Users/yallop/project:/project:cached \
 -v /host/another-path:/mount/another-point:consistent
 alpine command
```

## Semantics

The semantics of each configuration is described as a set of guarantees
relating to the observable effects of file system operations. In this
specification, "host" refers to the file system of the user's Docker
client.

### delegated

The `delegated` configuration provides the weakest set of guarantees.
For directories mounted with `delegated` the container's view of the
file system is authoritative, and writes performed by containers may not
be immediately reflected on the host file system. In situations such as NFS
asynchronous mode, if a running container with a `delegated` bind mount
crashes, then writes may be lost.

However, by relinquishing consistency, `delegated` mounts offer
significantly better performance than the other configurations. Where
the data written is ephemeral or readily reproducible, such as from scratch
space or build artifacts, `delegated` may be the right choice.

A `delegated` mount offers the following guarantees, which are presented
as constraints on the container run-time:

1.  If the implementation offers file system events, the container state
as it relates to a specific event **_must_** reflect the host file system
state at the time the event was generated if no container modifications
pertain to related file system state.

2.  If flush or sync operations are performed, relevant data **_must_** be
written back to the host file system.Between flush or sync
operations containers **_may_** cache data written, metadata modifications,
and directory structure changes.

3.  All containers hosted by the same runtime **_must_** share a consistent
cache of the mount.

4.  When any container sharing a `delegated` mount terminates, changes
to the mount **_must_** be written back to the host file system. If this
writeback fails, the container's execution **_must_** fail via exit code
and/or Docker event channels.

5.  If a `delegated` mount is shared with a `cached` or a `consistent`
mount, those portions that overlap **_must_** obey `cached` or `consistent`
mount semantics, respectively.

    Besides these constraints, the `delegated` configuration offers the
container runtime a degree of flexibility:

6. Containers **_may_** retain file data and metadata (including directory
structure, existence of nodes, etc) indefinitely and this cache **_may_**
desynchronize from the file system state of the host. Implementors should expire
caches when host file system changes occur, but this may be difficult to do on
a guaranteed timeframe due to platform limitations.

7. If changes to the mount source directory are present on the host
file system, those changes **_may_** be lost when the `delegated` mount
synchronizes with the host source directory.

8. Behaviors 6-7 **do not** apply to the file types of socket, pipe, or device.

### cached

The `cached` configuration provides all the guarantees of the `delegated`
configuration, and some additional guarantees around the visibility of writes
performed by containers. As such, `cached` typically improves the performance
of read-heavy workloads, at the cost of some temporary inconsistency between the
host and the container.

For directories mounted with `cached`, the host's view of
the file system is authoritative; writes performed by containers are immediately
visible to the host, but there may be a delay before writes performed on the
host are visible within containers.

>**Tip:** To learn more about `cached`, see the article on
[User-guided caching in Docker Desktop for Mac](https://blog.docker.com/2017/05/user-guided-caching-in-docker-for-mac/).

1. Implementations **_must_** obey `delegated` Semantics 1-5.

2. If the implementation offers file system events, the container state
as it relates to a specific event **_must_** reflect the host file system
state at the time the event was generated.

3. Container mounts **_must_** perform metadata modifications, directory
structure changes, and data writes consistently with the host file
system, and **_must not_** cache data written, metadata modifications, or
directory structure changes.

4.  If a `cached` mount is shared with a `consistent` mount, those portions
that overlap **_must_** obey `consistent` mount semantics.

    Some of the flexibility of the `delegated` configuration is retained,
namely:

5. Implementations **_may_** permit `delegated` Semantics 6.

### consistent

The `consistent` configuration places the most severe restrictions on
the container run-time. For directories mounted with `consistent` the
container and host views are always synchronized: writes performed
within the container are immediately visible on the host, and writes
performed on the host are immediately visible within the container.

The `consistent` configuration most closely reflects the behavior of
bind mounts on Linux. However, the overheads of providing strong
consistency guarantees make it unsuitable for a few use cases, where
performance is a priority and maintaining perfect consistency has low
priority.

1. Implementations **_must_** obey `cached` Semantics 1-4.

2. Container mounts **_must_** reflect metadata modifications, directory
structure changes, and data writes on the host file system immediately.

### default

The `default` configuration is identical to the `consistent`
configuration except for its name. Crucially, this means that `cached`
Semantics 4 and `delegated` Semantics 5 that require strengthening
overlapping directories do not apply to `default` mounts. This is the
default configuration if no `state` flags are supplied.
