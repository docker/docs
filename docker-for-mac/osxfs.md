---
description: OSXFS
keywords: mac, osxfs
redirect_from:
- /mackit/osxfs/
title: File system sharing (osxfs)
---

`osxfs` is a new shared file system solution, exclusive to Docker for Mac.
`osxfs` provides a close-to-native user experience for bind mounting OS X file
system trees into Docker containers. To this end, `osxfs` features a number of
unique capabilities as well as differences from a classical Linux file system.

### Case sensitivity

With Docker for Mac, file systems are shared from OS X into containers in the
same way as they operate in OS X. As a result, if a file system on OS X is
case-insensitive that behavior is shared by any bind mount from OS X into a
container. The default OS X file system is HFS+ and, during installation, it is
installed as case-insensitive by default. To get case-sensitive behavior from
your bind mounts, you must either create and format a ramdisk or external volume
as HFS+ with case-sensitivity or reformat your OS root partition with HFS+ with
case-sensitivity. We do not recommend reformatting your root partition as some
Mac software dubiously relies on case-insensitivity to function.


### Access control

`osxfs`, and therefore Docker, can access only those file system resources that
the Docker for Mac user has access to. `osxfs` does not run as `root`. If the OS
X user is an administrator, `osxfs` inherits those administrator privileges. We
are still evaluating which privileges to drop in the file system process to
balance security and ease-of-use. `osxfs` performs no additional permissions
checks and enforces no extra access control on accesses made through it. All
processes in containers can access the same objects in the same way as the
Docker user who started the containers.

### Namespaces

Much of the OS X file system that is accessible to the user is also available to
containers using the `-v` bind mount syntax. By default, you can share files in
`/Users`, `/Volumes`, `/private`, and `/tmp` directly. To add or remove
directory trees that are exported to Docker, use the **File sharing** tab in
Docker preferences <img src="../images/whale-x.png"> -> **Preferences** ->
**File sharing**. (See [Preferences](index.md#preferences).) All other paths
used in `-v` bind mounts are sourced from the Moby Linux VM running the Docker
containers, so arguments such as `-v /var/run/docker.sock:/var/run/docker.sock`
should work as expected. If an OS X path is not shared and does not exist in the
VM, an attempt to bind mount it will fail rather than create it in the VM. Paths
that already exist in the VM and contain files are reserved by Docker and cannot
be exported from OS X.

### Ownership

Initially, any containerized process that requests ownership metadata of
an object is told that its `uid` and `gid` own the object. When any
containerized process changes the ownership of a shared file system
object, e.g. with `chown`, the new ownership information is persisted in
the `com.docker.owner` extended attribute of the object. Subsequent
requests for ownership metadata will return the previously set
values. Ownership-based permissions are only enforced at the OS X file
system level with all accessing processes behaving as the user running
Docker. If the user does not have permission to read extended attributes
on an object (such as when that object's permissions are `0000`), `osxfs`
will attempt to add an access control list (ACL) entry that allows
the user to read and write extended attributes. If this attempt
fails, the object will appear to be owned by the process accessing
it until the extended attribute is readable again.

### File system events

Most `inotify` events are supported in bind mounts, and likely `dnotify` and
`fanotify` (though they have not been tested) are also supported. This means
that file system events from OS X are sent into containers and trigger any
listening processes there.

The following are **supported file system events**:

* Creation
* Modification
* Attribute changes
* Deletion
* Directory changes

The following are **partially supported file system events**:

* Move events trigger `IN_DELETE` on the source of the rename and
  `IN_MODIFY` on the destination of the rename

The following are **unsupported file system events**:

* Open
* Access
* Close events
* Unmount events (see <a href="osxfs.md#mounts">Mounts</a>)

Some events may be delivered multiple times. These limitations do not apply to
events between containers, only to those events originating in OS X.

### Mounts

The OS X mount structure is not visible in the shared volume, but volume
contents are visible. Volume contents appear in the same file system as the rest
of the shared file system. Mounting/unmounting OS X volumes that are also bind
mounted into containers may result in unexpected behavior in those containers.
Unmount events are not supported. Mount export support is planned but is still
under development.

### Symlinks

Symlinks are shared unmodified. This may cause issues when symlinks contain
paths that rely on the default case-insensitivity of the default OS X file
system, HFS+.

### File types

Symlinks, hardlinks, socket files, named pipes, regular files, and directories
are supported. Socket files and named pipes only transmit between containers and
between OS X processes -- no transmission across the hypervisor is supported,
yet. Character and block device files are not supported.

### Extended attributes

Extended attributes are not yet supported.

### Technology

`osxfs` does not use OSXFUSE. `osxfs` does not run under, inside, or
between OS X userspace processes and the OS X kernel.

### Performance issues, solutions, and roadmap

With regard to reported performance issues ([GitHub issue 77: File access in
mounted volumes extremely slow](https://github.com/docker/for-mac/issues/77)),
and a similar thread on [Docker for Mac forums on topic: File access in mounted
volumes extremely
slow](https://forums.docker.com/t/file-access-in-mounted-volumes-extremely-slow-cpu-bound/),
this topic provides an explanation of the issues, what we are doing to
address them, how the community can help us, and what you can expect in the
future. This explanation is a slightly re-worked version of an [understanding
performance
post](https://forums.docker.com/t/file-access-in-mounted-volumes-extremely-slow-cpu-bound/8076/158?u=orangesnap)
from David Sheets (@dsheets) on the [Docker development
team](https://forums.docker.com/groups/Docker) to the forum topic just
mentioned. We want to surface it in the documentation for wider reach.

#### Understanding performance

Perhaps the most important thing to understand is that shared file system
performance is multi-dimensional. This means that, depending on your workload,
you may experience exceptional, adequate, or poor performance with `osxfs`, the
file system server in Docker for Mac. File system APIs are very wide (20-40
message types) with many intricate semantics  involving on-disk state, in-memory
cache state, and concurrent access by multiple  processes. Additionally, `osxfs`
integrates   a mapping between OS X's FSEvents API and Linux's  inotify API
which is implemented inside of the file system itself complicating matters
further (cache behavior in particular).

At the highest level, there are two dimensions to file system performance:
throughput (read/write IO) and latency (roundtrip time). In a traditional file
system on a modern SSD, applications can generally expect throughput of a few
GB/s. With large sequential IO operations, `osxfs` can achieve throughput of
around 250 MB/s which, while not native speed, will not be the bottleneck for
most applications which perform acceptably on HDDs.

Latency is the time it takes for a file system call to complete. For instance,
the time between a thread issuing write in a container and resuming with the
number of bytes written. With a classical block-based file system, this latency
is typically under 10μs (microseconds). With `osxfs`, latency is presently
around 200μs for most operations or 20x slower. For workloads which demand many
sequential roundtrips, this results in significant observable slowdown. To
reduce the latency, we need to shorten the data path from a Linux system call to
OS X and back again. This requires tuning each component in the data path in
turn -- some of which require significant engineering effort. Even if we achieve
a huge latency reduction of 100μs/roundtrip, we will still "only" see a doubling
of performance. This is typical of performance engineering, which requires
significant effort to analyze slowdowns and develop optimized components. We
know how we can likely halve the roundtrip time but we haven't implemented those
improvements yet (more on this below in
[What you can do](osxfs.md#what-you-can-do)).

There is hope for significant performance improvement in the near term despite
these fundamental communication channel properties, which are difficult to
overcome (latency in particular). This hope comes in the form of increased
caching (storing "recent" values closer to their use to prevent roundtrips
completely). The Linux kernel's VFS layer contains a number of caches which can
be used to greatly improve performance by reducing the required communication
with the file system. Using this caching comes with a number of trade-offs:

* It requires understanding the cache behavior in detail in order to write
correct, stateful functionality on top of those caches.

* It harms the coherence or consistency of the file system as observed
from Linux containers and the OS X file system interfaces.

#### What we are doing

We are actively working on both increasing caching while mitigating the
associated issues and on reducing the file system data path latency. This
requires significant analysis of file system traces and speculative development
of system improvements to try to address specific performance issues. Perhaps
surprisingly, application workload can have a huge effect on performance. As an
example, here are two different use cases contributed on the [forum topic](https://forums.docker.com/t/file-access-in-mounted-volumes-extremely-slow-cpu-bound/)
and how their performance differs and suffers due to latency, caching, and
coherence:

1. A rake example (see below) appears to attempt to access 37000+
different files that don't exist on the shared volume. We can work very hard to
speed up all use cases by 2x via latency reduction but this use case will still
seem "slow". The ultimate solution for rake is to use a "negative dcache" that
keeps track of, in the Linux kernel itself, the files that do not exist.
Unfortunately, even this is not sufficient for the first time rake is run on a
shared directory. To handle that case, we actually need to develop a Linux
kernel patch which negatively caches all directory entries not in a
specified set -- and this cache must be kept up-to-date in real-time with the OS
X file system state even in the presence of missing OS X FSEvents messages and
so must be invalidated if OS X ever reports an event delivery failure.

2. Running ember build in a shared file system results in ember creating many
different temporary directories and performing lots of intermediate activity
within them. An empty ember project is over 300MB. This usage pattern does not
require coherence between Linux and OS X but, because we cannot distinguish this
fact at run-time, we maintain coherence during its hundreds of thousands of file
system accesses to manipulate temporary state. There is no "correct" solution in
this case. Either ember needs to change, the volume mount needs to have
coherence properties specified on it somehow, some heuristic needs to be
introduced to detect this access pattern and compensate, or the behavior needs
to be indicated via, e.g., extended attributes in the OS X file system.

These two examples come from performance use cases contributed by users and they
are incredibly helpful in prioritizing aspects of file system performance to
improve. We are developing statistical file system trace analysis tools
to characterize slow-performing workloads more easily in order to decide what to
work on next.

Under development, we have:

1. A Linux kernel patch to reduce data path latency by 2/7 copies and 2/5
context switches

2. Increased OS X integration to reduce the latency between the hypervisor and
the file system server

3. A server-side directory read cache to speed up traversal of large directories

4. User-facing file system tracing capabilities so that you can send us
recordings of slow workloads for analysis

5. A growing performance test suite of real world use cases (more on this below
in What you can do)

6. Experimental support for using Linux's inode, writeback, and page caches

7. End-user controls to configure the coherence of subsets of cross-OS bind
mounts without exposing all of the underlying complexity

#### What you can do

When you report shared file system performance issues, it is most helpful to
include a minimal Real World reproduction test case that demonstrates poor
performance.

Without a reproduction, it is very difficult for us to analyze your use case and
determine what improvements would speed it up. When you don't provide a
reproduction, one of us has to take the time to figure out the specific software
you are using and guess and hope that we have configured it in a typical way or
a way that has poor performance. That usually takes 1-4 hours depending on your
use case and once it is done, we must then determine what regular performance is
like and what kind of slow-down your use case is experiencing. In some cases, it
is not obvious what operation is even slow in your specific development
workflow. The additional set-up to reproduce the problem means we have less time
to fix bugs, develop analysis tools, or improve performance. So, please include
simple, immediate performance issue reproduction test cases. The [rake
reproduction
case](https://forums.docker.com/t/file-access-in-mounted-volumes-extremely-slow-cpu-bound/8076/103)
by @hirowatari shown in the forums thread is a great example.

This example originally provided:

1. A version-controlled repository so any changes/improvements to the test case
can be easily tracked.

2. A Dockerfile which constructs the exact image to run

3. A command-line invocation of how to start the container

4. A straight-forward way to measure the performance of the use case

5. A clear explanation (README) of how to run the test case

#### What you can expect

We will continue to work toward an optimized shared file system implementation
on the Beta channel of Docker for Mac.

You can expect some of the performance improvement work mentioned above to reach
the Beta channel in the coming release cycles.

In due course, we will open source all of our shared file system components. At
that time, we would be very happy to collaborate with you on improving the
implementation of osxfs and related software.

We still have on the slate to write up and publish details of shared file system
performance analysis and improvement on the Docker blog. Look for or nudge
@dsheets about those articles, which should serve as a jumping off point for
understanding the system, measuring it, or contributing to it.

#### Wrapping Up

We hope this gives you a rough idea of where `osxfs` performance is and where
it's going. We are treating good performance as a top priority feature of the
file system sharing component and we are actively working on improving it
through a number of different avenues. The osxfs project started in December
2015. Since the first integration into Docker for Mac in February 2016, we've
improved performance by 50x or more for many workloads while achieving nearly
complete POSIX compliance and without compromising coherence (it is shared and
not simply synced). Of course, in the beginning there was lots of low-hanging
fruit and now many of the remaining performance improvements require significant
engineering work on custom low-level components.

We appreciate your understanding as we continue development of the product and
work on all dimensions of performance. We want to continue to work with the
community on this, so please continue to report issues as you find them. We look
forward to collaborating with you on ideas and on the source code itself.
