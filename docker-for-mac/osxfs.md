---
aliases:
- /mackit/osxfs/
description: OSXFS
keywords:
- mac, osxfs
menu:
  main:
    identifier: mac-osxfs
    parent: pinata_mac_menu
    weight: 5
title: 'File system sharing '
---

# File system sharing (osxfs)

`osxfs` is a new shared file system solution, exclusive to Docker for
Mac. `osxfs` provides a close-to-native
user experience for bind mounting OS X file system trees into Docker
containers. To this end, `osxfs` features a number of unique
capabilities as well as differences from a classical Linux file system.

- [Case sensitivity](osxfs.md#case-sensitivity)
- [Access control](osxfs.md#access-control)
- [Namespaces](osxfs.md#namespaces)
- [Ownership](osxfs.md#ownership)
- [File system events](osxfs.md#file-system-events)
- [Mounts](osxfs.md#mounts)
- [Symlinks](osxfs.md#symlinks)
- [File types](osxfs.md#file-types)
- [Extended attributes](osxfs.md#extended-attributes)
- [Technology](osxfs.md#technology)

### Case sensitivity

With Docker for Mac, file systems are shared from OS X into containers
in the same way as they operate in OS X. As a result, if a file system
on OS X is case-insensitive that behavior is shared by any bind mount
from OS X into a container. The default OS X file system is HFS+ and,
during installation, it is installed as case-insensitive by default. To
get case-sensitive behavior from your bind mounts, you must either
create and format a ramdisk or external volume as HFS+ with
case-sensitivity or reformat your OS root partition with HFS+ with
case-sensitivity. We do not recommend reformatting your root partition
as some Mac software dubiously relies on case-insensitivity to function.


### Access control

`osxfs`, and therefore Docker, can access only those file system
resources that the Docker for Mac user has access to. `osxfs` does
not run as `root`. If the OS X user is an administrator, `osxfs` inherits
those administrator privileges. We are still evaluating which privileges
to drop in the file system process to balance security and
ease-of-use. `osxfs` performs no additional permissions checks and
enforces no extra access control on accesses made through it. All
processes in containers can access the same objects in the same way as
the Docker user who started the containers.

### Namespaces

Much of the OS X file system that is accessible to the user is also
available to containers using the `-v` bind mount syntax. By default,
you can share files in `/Users`, `/Volumes`, `/private`, and `/tmp`
directly. To add or remove directory trees that are exported to Docker,
use the **File sharing** tab in Docker preferences <img
src="../images/whale-x.png"> -> **Preferences** -> **File
sharing**. (See [Preferences](index.md#preferences).) All other paths
used in `-v` bind mounts are sourced from the Moby Linux VM running the
Docker containers, so arguments such as `-v
/var/run/docker.sock:/var/run/docker.sock` should work as expected. If
an OS X path is not shared and does not exist in the VM, an attempt to
bind mount it will fail rather than create it in the VM. Paths that
already exist in the VM and contain files are reserved by Docker and
cannot be exported from OS X.

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
on an object, e.g. when that object's permissions are `0000`, `osxfs`
will attempt to add an access control list entry allowing the user to
read and write extended attributes. If this is not possible or extended
attribute permissions are still denied, ownership will be reported as
the accessing process until the extended attribute is again readable.

### File system events

Most `inotify` events are supported in bind mounts, and likely `dnotify`
and `fanotify` (though they have not been tested) are also supported.
This means that file system events from OS X are sent into containers
and trigger any listening processes there.

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

Some events may be delivered multiple times. Events are not delivered for bind mounts from symlinks (notably `/tmp` will not deliver inotify events but
`/private/tmp` will). These limitations do not apply to events between
containers, only to those events originating in OS X.


### Mounts

The OS X mount structure is not visible in the shared volume, but volume
contents are visible. Volume contents appear in the same file system as the
rest of the shared file system. Mounting/unmounting OS X volumes that
are also bind mounted into containers may result in unexpected behavior
in those containers. Unmount events are not supported. Mount export
support is planned but is still under development.

### Symlinks

Symlinks are shared unmodified. This may cause issues when symlinks
contain paths that rely on the default case-insensitivity of the
default OS X file system, HFS+.

### File types

Symlinks, hardlinks, socket files, named pipes, regular files, and
directories are supported. Socket files and named pipes only transmit
between containers and between OS X processes -- no transmission across
the hypervisor is supported, yet. Character and block device files are
not supported.

### Extended attributes

Extended attributes are not yet supported.

### Technology

`osxfs` does not use OSXFUSE. `osxfs` does not run under, inside, or
between OS X userspace processes and the OS X kernel.
