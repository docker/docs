---
description: Enabling seccomp in Docker
keywords: seccomp, security, docker, documentation
title: Seccomp security profiles for Docker
---

Secure computing mode (`seccomp`) is a Linux kernel feature. You can use it to
restrict the actions available within the container. The `seccomp()` system
call operates on the seccomp state of the calling process. You can use this
feature to restrict your application's access.

This feature is available only if Docker has been built with `seccomp` and the
kernel is configured with `CONFIG_SECCOMP` enabled. To check if your kernel
supports `seccomp`:

```console
$ grep CONFIG_SECCOMP= /boot/config-$(uname -r)
CONFIG_SECCOMP=y
```

> **Note**: `seccomp` profiles require seccomp 2.2.1 or later and is _disabled_ on containers running with `--privilege`.

## Pass a profile for a container

The default `seccomp` profile provides a sane default for running containers with
seccomp and unconditionally disables around 15 system calls and conditionally
gates another 50 out of 300+. It is moderately protective while providing wide 
application compatibility. The default Docker profile can be found
[here](https://github.com/moby/moby/blob/master/profiles/seccomp/default.json).

In effect, the profile is a allowlist which denies access to system calls by
default, then allowlists specific system calls. The profile works by defining a
`defaultAction` of `SCMP_ACT_ERRNO` and overriding that action only for specific
system calls. The effect of `SCMP_ACT_ERRNO` is to cause a `Permission Denied`
error. Next, the profile defines a specific list of system calls which are fully
allowed, because their `action` is overridden to be `SCMP_ACT_ALLOW`. Finally,
some specific rules are for individual system calls such as `personality`, and others, 
to allow variants of those system calls with specific arguments.

`seccomp` is instrumental for running Docker containers with least privilege. It
is not recommended to change the default `seccomp` profile.

When you run a container, it uses the default profile unless you override it
with the `--security-opt` option. For example, the following explicitly
specifies a policy:

```console
$ docker run --rm \
             -it \
             --security-opt seccomp=/path/to/seccomp/profile.json \
             hello-world
```

### Significant syscalls blocked by the default profile

Docker's default seccomp profile is an allowlist which specifies the calls that
are allowed. The table below lists the significant (but not all) syscalls that
are effectively blocked because they are not on the Allowlist. The table includes
the reason each syscall is blocked rather than white-listed.

| Syscall             | Description                                                                                                                           |
|---------------------|---------------------------------------------------------------------------------------------------------------------------------------|
| `add_key`           | Prevent containers from using the kernel keyring, which is not namespaced.                                   |
| `get_kernel_syms`   | Deny retrieval of exported kernel and module symbols. Obsolete.                                              |                              
| `keyctl`            | Prevent containers from using the kernel keyring, which is not namespaced.                                   |
| `move_pages`        | Syscall that modifies kernel memory and NUMA settings.                                                       |
| `nfsservctl`        | Deny interaction with the kernel nfs daemon. Obsolete since Linux 3.1.                                       |
| `perf_event_open`   | Tracing/profiling syscall, which could leak a lot of information on the host.                                |
| `pivot_root`        | Deny `pivot_root`, should be privileged operation.                                                           |
| `query_module`      | Deny manipulation and functions on kernel modules. Obsolete.                                                 |
| `request_key`       | Prevent containers from using the kernel keyring, which is not namespaced.                                   |
| `sysfs`             | Obsolete syscall.                                                                                            |
| `_sysctl`           | Obsolete, replaced by /proc/sys.                                                                             |
| `uselib`            | Older syscall related to shared libraries, unused for a long time.                                           |
| `userfaultfd`       | Userspace page fault handling, largely needed for process migration.                                         |
| `ustat`             | Obsolete syscall.                                                                                            |

### Significant syscalls _conditionally_ blocked by the default profile

These are a list of system calls which are gated based on the presence of a 
certain capability or argument value.

| Syscall               | Description                                                                                                     | Conditions to permit syscall    |
|-----------------------|-----------------------------------------------------------------------------------------------------------------|---------------------------------|
| `acct`                | Accounting syscall which could let containers disable their own resource limits or process accounting. | `CAP_SYS_PACCT` capability.  |
| `arch_prctl`          | Set arch-specific thread state.                                                                        | amd64 or x32 arch only.      |
| `arm_fadvise64_64`    | Predeclare an access pattern for file data.                                                            | arm or arm64 arch only.      |
| `arm_sync_file_range` | Arm-specific sync a file segment with disk.                                                            | arm or arm64 arch only.      |
| `bpf`                 | Loading potentially persistent bpf programs into kernel.                                               | `CAP_SYS_ADMIN` capability.  |
| `breakpoint`          | Arm-specific instruction.                                                                              | arm and arm64 arch only.     |
| `cacheflush`          | Flush contents of instruction and/or data cache.                                                       | arm or arm64 arch only.      |
| `chroot`              | Change root directory of calling process.                                                              | `CAP_SYS_CHROOT` capability. |
| `clock_settime`       | Time/date is not namespaced.                                                                           | `CAP_SYS_TIME` capability.   |
| `clone`               | Cloning new namespaces.  | `CAP_SYS_ADMIN` capability for CLONE_* flags, except `CLONE_NEWUSER` and specific mask values on specific arch (see [code](https://github.com/moby/moby/blob/master/profiles/seccomp/default.json)).| 
| `create_module`       | Manipulation and functions on kernel modules. Obsolete.                                                | `CAP_SYS_MODULE` capability. |
| `delete_module`       | Manipulation and functions on kernel modules.                                                          | `CAP_SYS_MODULE` capability. |
| `finit_module`        | Manipulation and functions on kernel modules.                                                          | `CAP_SYS_MODULE` capability. |
| `get_mempolicy`       | Syscall that modifies kernel memory and NUMA settings.                                                 | `CAP_SYS_NICE` capability.   |
| `init_module`         | Manipulation and functions on kernel modules.                                                          | `CAP_SYS_MODULE` capability. |
| `ioperm`              | Modifying kernel I/O privilege levels.                                                                 | `CAP_SYS_RAWIO` capability.  |
| `iopl`                | Modifying kernel I/O privilege levels.                                                                 | `CAP_SYS_RAWIO` capability.  |
| `kcmp`                | Restrict process inspection capabilities.                                                              | `CAP_SYS_PTRACE` capability. |
| `kexec_file_load`     | Sister syscall of `kexec_load` that does the same thing, slightly different arguments.                 | `CAP_SYS_BOOT` capability.   |
| `kexec_load`          | Loading a new kernel for later execution.                                                              | `CAP_SYS_BOOT` capability.   |
| `lookup_dcookie`      | Tracing/profiling syscall, which could leak a lot of information on the host.                          | `CAP_SYS_ADMIN` capability.  |
| `mbind`               | Syscall that modifies kernel memory and NUMA settings.                                                 | `CAP_SYS_NICE` capability.   |
| `modify_ldt`          | Get or set a per-process LDT entry.                                                                    | amd64 or x32 or x64 arch only. |
| `mount`               | Mounting.                                                                                              | `CAP_SYS_ADMIN` capability. |
| `name_to_handle_at`   | Sister syscall to `open_by_handle_at`.                                                                 | `CAP_DAC_READ_SEARCH` capability. |
| `open_by_handle_at`   | Cause of an old container breakout.                                                                    | `CAP_DAC_READ_SEARCH`capability. |
| `personality`         | Sets an execution domain (personality) of the caller's process.  | persona values of 0, 8, 131072, 131080, 4294967295 as defined in /include/linux/personality.h |
| `process_vm_readv`    | Process inspection capabilities.                                                     | `CAP_SYS_PTRACE` capability and kernel version 4.8+ |
| `process_vm_writev`   | Process inspection capabilities.                                                     | `CAP_SYS_PTRACE` capability and kernel version 4.8+|
| `ptrace`              | Tracing/profiling syscall. Blocked in Linux kernel versions before 4.8 to avoid seccomp bypass. Tracing/profiling arbitrary processes is already blocked by dropping `CAP_SYS_PTRACE`, because it could leak a lot of information on the host. | `CAP_SYS_PTRACE` capability and kernel version 4.8+. |
| `quotactl`            | Quota syscall which could let containers disable their own resource limits or process accounting.     | `CAP_SYS_ADMIN` capability. |
| `reboot`              | Don't let containers reboot the host unless given that capability.                                    | `CAP_SYS_BOOT` capability. |
| `s390_pci_mmio_read`  | Transfer data to/from PCI MMIO memory page.                                                           | s390 or s390x arch only. |
| `s390_pci_mmio_write` | Transfer data to/from PCI MMIO memory page.                                                           | s390 or s390x arch only. |
| `s390_runtime_instr`  | Enable/disable s390 CPU run-time instrumentation.                                                     | s390 or s390x arch only. |
| `set_mempolicy`       | Syscall that modifies kernel memory and NUMA settings.                                                | `CAP_SYS_NICE` capability. |
| `set_tls`             | Arm-specific thread local storage call.                                                               | arm or arm64 arch only. |
| `setns`               | Associating a thread with a namespace.                                                                | `CAP_SYS_ADMIN` capability. |
| `settimeofday`        | Time/date is not namespaced.                                                                          | `CAP_SYS_TIME` capability. |
| `stime`               | Time/date is not namespaced.                                                                          | `CAP_SYS_TIME` capability. |
| `swapon`              | Start/stop swapping to file/device.                                                                   | `CAP_SYS_ADMIN` capability. |
| `swapoff`             | Start/stop swapping to file/device.                                                                   | `CAP_SYS_ADMIN` capability. |
| `sync_file_range2`    | Fine control sync of a file segment with disk.                                                      | ppc64le, arm or arm64 architecture only.  |
| `syslog`              | Read and/or clear kernel message ring buffer.                                                       | `CAP_SYSLOG` capability. |
| `umount`              | Should be a privileged operation.                                                                   | `CAP_SYS_ADMIN` capability. |
| `umount2`             | Should be a privileged operation.                                                                   | `CAP_SYS_ADMIN` capability. |
| `unshare`             | Cloning new namespaces for processes.                                     | `CAP_SYS_ADMIN` capability with the exception of `unshare --user`. |
| `vhangup`             | Virtually hangup the current terminal.                                    | `CAP_SYS_TTY_CONFIG` capability. |
| `vm86`                | In kernel x86 real mode virtual machine.                                  | `CAP_SYS_ADMIN` capability. |
| `vm86old`             | In kernel x86 real mode virtual machine.                                  | `CAP_SYS_ADMIN` capability. |


## Run without the default seccomp profile

You can pass `unconfined` to run a container without the default seccomp
profile.

```console
$ docker run --rm -it --security-opt seccomp=unconfined debian:bullseye \
    unshare --map-root-user --user sh -c whoami
```
